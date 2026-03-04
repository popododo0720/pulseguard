package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pulseguard/pulseguard/internal/server/grpc"
	"github.com/pulseguard/pulseguard/internal/server/store"
	"github.com/pulseguard/pulseguard/internal/server/webhook"
	"github.com/rs/cors"
)

// NewRouter creates and configures the chi router with all API routes.
func NewRouter(s *store.Store, dispatcher *grpc.CommandDispatcher, proxy *webhook.Proxy, devMode bool, webDir string, token ...string) http.Handler {
	authToken := ""
	if len(token) > 0 {
		authToken = token[0]
	}

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
		r.Delete("/agents/{id}", deleteAgentHandler(s))

		// Jobs
		r.Get("/jobs", listJobsHandler(s))
		r.Post("/jobs", createJobHandler(s))
		r.Get("/jobs/{id}", getJobHandler(s))
		r.Put("/jobs/{id}", updateJobHandler(s))
		r.Delete("/jobs/{id}", deleteJobHandler(s))
		r.Post("/jobs/{id}/rerun", rerunJobHandler(s, dispatcher))
		r.Get("/jobs/{id}/executions", listJobExecutionsHandler(s))
		r.Post("/jobs/{id}/report", reportJobResultHandler(s, authToken))

		// Notifications
		r.Get("/notifications", listNotificationChannelsHandler(s))
		r.Post("/notifications", createNotificationChannelHandler(s))
		r.Put("/notifications/{id}", updateNotificationChannelHandler(s))
		r.Delete("/notifications/{id}", deleteNotificationChannelHandler(s))

		// Settings
		r.Get("/settings", settingsHandler(authToken))

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

	// SPA static file serving (production mode)
	if !devMode && webDir != "" {
		r.Get("/*", spaHandler(webDir))
	}

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

// spaHandler serves static files from webDir, falling back to index.html for SPA routing.
func spaHandler(webDir string) http.HandlerFunc {
	fs := http.Dir(webDir)
	fileServer := http.FileServer(fs)

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		fullPath := filepath.Join(webDir, filepath.Clean(strings.TrimPrefix(path, "/")))
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
			return
		}

		fileServer.ServeHTTP(w, r)
	}
}
