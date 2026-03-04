package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pulseguard/pulseguard/internal/server/grpc"
	"github.com/pulseguard/pulseguard/internal/server/store"
	"github.com/pulseguard/pulseguard/internal/server/webhook"
	"github.com/rs/cors"
)

// NewRouter creates and configures the chi router with all API routes.
func NewRouter(s *store.Store, dispatcher *grpc.CommandDispatcher, proxy *webhook.Proxy, devMode bool) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(Recoverer)
	r.Use(RequestLogger)

	// API routes with JSON content type
	r.Route("/api", func(r chi.Router) {
		r.Use(JSONContentType)

		// Dashboard
		r.Get("/dashboard", dashboardHandler(s))

		// Agents
		r.Get("/agents", listAgentsHandler(s))
		r.Get("/agents/{id}", getAgentHandler(s))

		// Jobs
		r.Get("/jobs", listJobsHandler(s))
		r.Post("/jobs", createJobHandler(s))
		r.Get("/jobs/{id}", getJobHandler(s))
		r.Put("/jobs/{id}", updateJobHandler(s))
		r.Delete("/jobs/{id}", deleteJobHandler(s))
		r.Post("/jobs/{id}/rerun", rerunJobHandler(s, dispatcher))
		r.Get("/jobs/{id}/executions", listJobExecutionsHandler(s))

		// Webhook endpoints
		r.Get("/webhook-endpoints", listWebhookEndpointsHandler(s))
		r.Post("/webhook-endpoints", createWebhookEndpointHandler(s))
		r.Get("/webhook-endpoints/{id}", getWebhookEndpointHandler(s))
		r.Put("/webhook-endpoints/{id}", updateWebhookEndpointHandler(s))
		r.Delete("/webhook-endpoints/{id}", deleteWebhookEndpointHandler(s))
		r.Get("/webhook-endpoints/{id}/requests", listWebhookRequestsHandler(s))
		r.Post("/webhook-endpoints/{id}/requests/{reqId}/replay", replayWebhookHandler(s, proxy))
	})

	// Webhook receiver (not under /api, no JSON middleware)
	r.Post("/wh/{slug}", webhookReceiverHandler(s, proxy))

	var handler http.Handler = r
	if devMode {
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		})
		handler = c.Handler(r)
	}

	return handler
}
