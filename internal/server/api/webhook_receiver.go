package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pulseguard/pulseguard/internal/models"
	"github.com/pulseguard/pulseguard/internal/server/store"
	"github.com/pulseguard/pulseguard/internal/server/webhook"
)

func webhookReceiverHandler(s *store.Store, proxy *webhook.Proxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		ep, err := s.GetWebhookEndpointBySlug(slug)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		if !ep.Enabled {
			http.Error(w, "endpoint disabled", http.StatusServiceUnavailable)
			return
		}

		body, err := io.ReadAll(io.LimitReader(r.Body, 10<<20)) // 10MB max
		if err != nil {
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}

		headersJSON, _ := json.Marshal(r.Header)

		whReq := &models.WebhookRequest{
			ID:          uuid.New().String(),
			EndpointID:  ep.ID,
			Method:      r.Method,
			Headers:     string(headersJSON),
			Body:        string(body),
			QueryParams: r.URL.RawQuery,
			SourceIP:    r.RemoteAddr,
		}

		if err := s.CreateWebhookRequest(whReq); err != nil {
			http.Error(w, "failed to store request", http.StatusInternalServerError)
			return
		}

		_ = s.IncrementWebhookRequestCount(ep.ID)

		// Forward asynchronously
		go proxy.Forward(ep, whReq)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":     "received",
			"request_id": whReq.ID,
		})
	}
}
