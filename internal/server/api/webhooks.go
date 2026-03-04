package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pulseguard/pulseguard/internal/models"
	"github.com/pulseguard/pulseguard/internal/server/store"
	"github.com/pulseguard/pulseguard/internal/server/webhook"
)

type createWebhookEndpointRequest struct {
	Name      string            `json:"name"`
	Slug      string            `json:"slug"`
	TargetURL string            `json:"target_url"`
	Headers   map[string]string `json:"headers"`
	Secret    string            `json:"secret"`
	Enabled   bool              `json:"enabled"`
}

func listWebhookEndpointsHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoints, err := s.ListWebhookEndpoints()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, endpoints)
	}
}

func createWebhookEndpointHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createWebhookEndpointRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Headers == nil {
			req.Headers = map[string]string{}
		}

		ep := &models.WebhookEndpoint{
			ID:        uuid.New().String(),
			Name:      req.Name,
			Slug:      req.Slug,
			TargetURL: req.TargetURL,
			Headers:   req.Headers,
			Secret:    req.Secret,
			Enabled:   req.Enabled,
		}

		if err := s.CreateWebhookEndpoint(ep); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusCreated, ep)
	}
}

func getWebhookEndpointHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ep, err := s.GetWebhookEndpoint(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "webhook endpoint not found")
			return
		}
		writeJSON(w, http.StatusOK, ep)
	}
}

func updateWebhookEndpointHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		existing, err := s.GetWebhookEndpoint(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "webhook endpoint not found")
			return
		}

		var req createWebhookEndpointRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		existing.Name = req.Name
		existing.Slug = req.Slug
		existing.TargetURL = req.TargetURL
		existing.Headers = req.Headers
		existing.Secret = req.Secret
		existing.Enabled = req.Enabled

		if err := s.UpdateWebhookEndpoint(existing); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, existing)
	}
}

func deleteWebhookEndpointHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := s.DeleteWebhookEndpoint(id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

func listWebhookRequestsHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		reqs, err := s.ListWebhookRequests(id, 100)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, reqs)
	}
}

func replayWebhookHandler(s *store.Store, proxy *webhook.Proxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpointID := chi.URLParam(r, "id")
		reqID := chi.URLParam(r, "reqId")

		ep, err := s.GetWebhookEndpoint(endpointID)
		if err != nil {
			writeError(w, http.StatusNotFound, "endpoint not found")
			return
		}

		whReq, err := s.GetWebhookRequest(reqID)
		if err != nil {
			writeError(w, http.StatusNotFound, "request not found")
			return
		}

		go proxy.Forward(ep, whReq)

		writeJSON(w, http.StatusAccepted, map[string]string{
			"status":     "replaying",
			"request_id": reqID,
		})
	}
}
