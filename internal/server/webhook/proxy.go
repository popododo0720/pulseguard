package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pulseguard/pulseguard/internal/models"
	"github.com/pulseguard/pulseguard/internal/server/store"
)

const maxRetries = 3

// Proxy handles forwarding webhook requests to target URLs.
type Proxy struct {
	store  *store.Store
	client *http.Client
}

// NewProxy creates a new webhook proxy.
func NewProxy(s *store.Store) *Proxy {
	return &Proxy{
		store: s,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Forward sends a webhook request to the endpoint's target URL with retry and exponential backoff.
func (p *Proxy) Forward(ep *models.WebhookEndpoint, whReq *models.WebhookRequest) {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		delivery := &models.WebhookDelivery{
			ID:            uuid.New().String(),
			RequestID:     whReq.ID,
			EndpointID:    ep.ID,
			Status:        "pending",
			AttemptNumber: attempt,
		}

		respStatus, respBody, err := p.doForward(ep, whReq)

		now := time.Now().UTC()
		delivery.DeliveredAt = &now

		if err != nil {
			delivery.Status = "failed"
			delivery.Error = err.Error()

			if attempt < maxRetries {
				delay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
				nextRetry := now.Add(delay)
				delivery.NextRetryAt = &nextRetry
			}
		} else if respStatus >= 200 && respStatus < 300 {
			delivery.Status = "delivered"
			delivery.ResponseStatus = &respStatus
			delivery.ResponseBody = respBody
		} else {
			delivery.Status = "failed"
			delivery.ResponseStatus = &respStatus
			delivery.ResponseBody = respBody
			delivery.Error = fmt.Sprintf("HTTP %d", respStatus)

			if attempt < maxRetries {
				delay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
				nextRetry := now.Add(delay)
				delivery.NextRetryAt = &nextRetry
			}
		}

		if saveErr := p.store.CreateWebhookDelivery(delivery); saveErr != nil {
			slog.Error("failed to save webhook delivery", "error", saveErr)
		}

		if delivery.Status == "delivered" {
			slog.Info("webhook delivered", "endpoint", ep.Slug, "request_id", whReq.ID, "attempt", attempt)
			return
		}

		if attempt < maxRetries {
			delay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
			slog.Warn("webhook delivery failed, retrying", "endpoint", ep.Slug, "attempt", attempt, "delay", delay)
			time.Sleep(delay)
		}
	}

	slog.Error("webhook delivery failed after all retries", "endpoint", ep.Slug, "request_id", whReq.ID)
}

func (p *Proxy) doForward(ep *models.WebhookEndpoint, whReq *models.WebhookRequest) (int, string, error) {
	req, err := http.NewRequest(whReq.Method, ep.TargetURL, bytes.NewBufferString(whReq.Body))
	if err != nil {
		return 0, "", fmt.Errorf("create request: %w", err)
	}

	// Set original headers
	var origHeaders map[string][]string
	if err := json.Unmarshal([]byte(whReq.Headers), &origHeaders); err == nil {
		for k, vals := range origHeaders {
			for _, v := range vals {
				req.Header.Add(k, v)
			}
		}
	}

	// Override with endpoint-configured headers
	for k, v := range ep.Headers {
		req.Header.Set(k, v)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1MB max
	return resp.StatusCode, string(body), nil
}
