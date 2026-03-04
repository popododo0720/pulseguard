package store

import (
	"database/sql"
	"encoding/json"

	"github.com/pulseguard/pulseguard/internal/models"
)

// Webhook Endpoints

func (s *Store) CreateWebhookEndpoint(e *models.WebhookEndpoint) error {
	headersJSON, _ := json.Marshal(e.Headers)
	_, err := s.db.Exec(`
		INSERT INTO webhook_endpoints (id, name, slug, target_url, headers, secret, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`,
		e.ID, e.Name, e.Slug, e.TargetURL, string(headersJSON), e.Secret, boolToInt(e.Enabled),
	)
	return err
}

func (s *Store) GetWebhookEndpoint(id string) (*models.WebhookEndpoint, error) {
	row := s.db.QueryRow(`SELECT id, name, slug, target_url, headers, secret, enabled, request_count, last_request_at, created_at, updated_at FROM webhook_endpoints WHERE id = ?`, id)
	return scanWebhookEndpoint(row)
}

func (s *Store) GetWebhookEndpointBySlug(slug string) (*models.WebhookEndpoint, error) {
	row := s.db.QueryRow(`SELECT id, name, slug, target_url, headers, secret, enabled, request_count, last_request_at, created_at, updated_at FROM webhook_endpoints WHERE slug = ?`, slug)
	return scanWebhookEndpoint(row)
}

func (s *Store) ListWebhookEndpoints() ([]*models.WebhookEndpoint, error) {
	rows, err := s.db.Query(`SELECT id, name, slug, target_url, headers, secret, enabled, request_count, last_request_at, created_at, updated_at FROM webhook_endpoints ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []*models.WebhookEndpoint
	for rows.Next() {
		e, err := scanWebhookEndpointRows(rows)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, e)
	}
	return endpoints, rows.Err()
}

func (s *Store) UpdateWebhookEndpoint(e *models.WebhookEndpoint) error {
	headersJSON, _ := json.Marshal(e.Headers)
	_, err := s.db.Exec(`
		UPDATE webhook_endpoints SET name=?, slug=?, target_url=?, headers=?, secret=?, enabled=?, updated_at=datetime('now')
		WHERE id=?`,
		e.Name, e.Slug, e.TargetURL, string(headersJSON), e.Secret, boolToInt(e.Enabled), e.ID,
	)
	return err
}

func (s *Store) DeleteWebhookEndpoint(id string) error {
	_, err := s.db.Exec(`DELETE FROM webhook_endpoints WHERE id = ?`, id)
	return err
}

func (s *Store) IncrementWebhookRequestCount(endpointID string) error {
	_, err := s.db.Exec(`UPDATE webhook_endpoints SET request_count = request_count + 1, last_request_at = datetime('now'), updated_at = datetime('now') WHERE id = ?`, endpointID)
	return err
}

// Webhook Requests

func (s *Store) CreateWebhookRequest(r *models.WebhookRequest) error {
	_, err := s.db.Exec(`
		INSERT INTO webhook_requests (id, endpoint_id, method, headers, body, query_params, source_ip, received_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'))`,
		r.ID, r.EndpointID, r.Method, r.Headers, r.Body, r.QueryParams, r.SourceIP,
	)
	return err
}

func (s *Store) GetWebhookRequest(id string) (*models.WebhookRequest, error) {
	var r models.WebhookRequest
	var receivedAt string
	err := s.db.QueryRow(`SELECT id, endpoint_id, method, headers, body, query_params, source_ip, received_at FROM webhook_requests WHERE id = ?`, id).
		Scan(&r.ID, &r.EndpointID, &r.Method, &r.Headers, &r.Body, &r.QueryParams, &r.SourceIP, &receivedAt)
	if err != nil {
		return nil, err
	}
	r.ReceivedAt = parseDBTime(receivedAt)
	return &r, nil
}

func (s *Store) ListWebhookRequests(endpointID string, limit int) ([]*models.WebhookRequest, error) {
	rows, err := s.db.Query(`SELECT id, endpoint_id, method, headers, body, query_params, source_ip, received_at FROM webhook_requests WHERE endpoint_id = ? ORDER BY received_at DESC LIMIT ?`, endpointID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reqs []*models.WebhookRequest
	for rows.Next() {
		var r models.WebhookRequest
		var receivedAt string
		if err := rows.Scan(&r.ID, &r.EndpointID, &r.Method, &r.Headers, &r.Body, &r.QueryParams, &r.SourceIP, &receivedAt); err != nil {
			return nil, err
		}
		r.ReceivedAt = parseDBTime(receivedAt)
		reqs = append(reqs, &r)
	}
	return reqs, rows.Err()
}

// Webhook Deliveries

func (s *Store) CreateWebhookDelivery(d *models.WebhookDelivery) error {
	_, err := s.db.Exec(`
		INSERT INTO webhook_deliveries (id, request_id, endpoint_id, status, response_status, response_body, error, attempt_number, delivered_at, next_retry_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))`,
		d.ID, d.RequestID, d.EndpointID, d.Status, d.ResponseStatus, d.ResponseBody, d.Error, d.AttemptNumber,
		nullTimeStr(d.DeliveredAt), nullTimeStr(d.NextRetryAt),
	)
	return err
}

func (s *Store) UpdateWebhookDelivery(d *models.WebhookDelivery) error {
	_, err := s.db.Exec(`
		UPDATE webhook_deliveries SET status=?, response_status=?, response_body=?, error=?, delivered_at=?, next_retry_at=?
		WHERE id=?`,
		d.Status, d.ResponseStatus, d.ResponseBody, d.Error,
		nullTimeStr(d.DeliveredAt), nullTimeStr(d.NextRetryAt), d.ID,
	)
	return err
}

// Helpers

func scanWebhookEndpoint(row *sql.Row) (*models.WebhookEndpoint, error) {
	var e models.WebhookEndpoint
	var headersJSON string
	var enabled int
	var lastReq sql.NullString
	var createdAt, updatedAt string

	err := row.Scan(&e.ID, &e.Name, &e.Slug, &e.TargetURL, &headersJSON, &e.Secret, &enabled, &e.RequestCount, &lastReq, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(headersJSON), &e.Headers); err != nil {
		e.Headers = map[string]string{}
	}
	e.Enabled = enabled == 1
	if lastReq.Valid {
		t := parseDBTime(lastReq.String)
		if !t.IsZero() {
			e.LastRequestAt = &t
		}
	}
	e.CreatedAt = parseDBTime(createdAt)
	e.UpdatedAt = parseDBTime(updatedAt)
	return &e, nil
}

func scanWebhookEndpointRows(rows *sql.Rows) (*models.WebhookEndpoint, error) {
	var e models.WebhookEndpoint
	var headersJSON string
	var enabled int
	var lastReq sql.NullString
	var createdAt, updatedAt string

	err := rows.Scan(&e.ID, &e.Name, &e.Slug, &e.TargetURL, &headersJSON, &e.Secret, &enabled, &e.RequestCount, &lastReq, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(headersJSON), &e.Headers); err != nil {
		e.Headers = map[string]string{}
	}
	e.Enabled = enabled == 1
	if lastReq.Valid {
		t := parseDBTime(lastReq.String)
		if !t.IsZero() {
			e.LastRequestAt = &t
		}
	}
	e.CreatedAt = parseDBTime(createdAt)
	e.UpdatedAt = parseDBTime(updatedAt)
	return &e, nil
}
