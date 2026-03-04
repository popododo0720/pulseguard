package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pulseguardv1 "github.com/pulseguard/pulseguard/gen/pulseguard/v1"
	"github.com/pulseguard/pulseguard/internal/server/api"
	servergrpc "github.com/pulseguard/pulseguard/internal/server/grpc"
	"github.com/pulseguard/pulseguard/internal/server/store"
	"github.com/pulseguard/pulseguard/internal/server/webhook"
)

//go:embed migrations
var migrationsFS embed.FS

func main() {
	port := flag.Int("port", 8080, "REST API port")
	grpcPort := flag.Int("grpc-port", 9090, "gRPC server port")
	dbPath := flag.String("db", "./pulseguard.db", "SQLite database path")
	devMode := flag.Bool("dev", false, "Enable dev mode (CORS for localhost:5173)")
	token := flag.String("token", os.Getenv("PULSEGUARD_TOKEN"), "Auth token for agent connections")
	flag.Parse()

	slog.Info("starting PulseGuard server",
		"rest_port", *port,
		"grpc_port", *grpcPort,
		"db", *dbPath,
		"dev", *devMode,
	)

	// Initialize store
	s, err := store.New(*dbPath)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer s.Close()

	// Run migrations
	migrationSQL, err := migrationsFS.ReadFile("migrations/001_initial.sql")
	if err != nil {
		slog.Error("failed to read migration file", "error", err)
		os.Exit(1)
	}
	if err := s.Init(string(migrationSQL)); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Command dispatcher for agent communication
	dispatcher := servergrpc.NewCommandDispatcher()

	// Webhook proxy
	proxy := webhook.NewProxy(s)

	// gRPC server
	grpcServer := grpc.NewServer()
	agentSvc := servergrpc.NewServer(s, dispatcher, *token)
	pulseguardv1.RegisterAgentServiceServer(grpcServer, agentSvc)

	// REST server
	router := api.NewRouter(s, dispatcher, proxy, *devMode)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}

	// Graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	// Start gRPC server
	g.Go(func() error {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
		if err != nil {
			return fmt.Errorf("listen grpc: %w", err)
		}
		slog.Info("gRPC server listening", "port", *grpcPort)
		return grpcServer.Serve(lis)
	})

	// Start REST server
	g.Go(func() error {
		slog.Info("REST server listening", "port", *port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	})

	// Shutdown handler
	g.Go(func() error {
		<-ctx.Done()
		slog.Info("shutting down servers...")
		grpcServer.GracefulStop()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
