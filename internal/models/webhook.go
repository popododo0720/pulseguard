package models

import "time"

type WebhookEndpoint struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Slug          string            `json:"slug"`
	TargetURL     string            `json:"target_url"`
	Headers       map[string]string `json:"headers"`
	Secret        string            `json:"secret,omitempty"`
	Enabled       bool              `json:"enabled"`
	RequestCount  int               `json:"request_count"`
	LastRequestAt *time.Time        `json:"last_request_at,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type WebhookRequest struct {
	ID          string    `json:"id"`
	EndpointID  string    `json:"endpoint_id"`
	Method      string    `json:"method"`
	Headers     string    `json:"headers"`
	Body        string    `json:"body"`
	QueryParams string    `json:"query_params"`
	SourceIP    string    `json:"source_ip"`
	ReceivedAt  time.Time `json:"received_at"`
}

type WebhookDelivery struct {
	ID             string     `json:"id"`
	RequestID      string     `json:"request_id"`
	EndpointID     string     `json:"endpoint_id"`
	Status         string     `json:"status"`
	ResponseStatus *int       `json:"response_status,omitempty"`
	ResponseBody   string     `json:"response_body"`
	Error          string     `json:"error"`
	AttemptNumber  int        `json:"attempt_number"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
	NextRetryAt    *time.Time `json:"next_retry_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
