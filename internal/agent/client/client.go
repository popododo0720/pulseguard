package client

import (
	"context"
	"fmt"
	"log/slog"

	pulseguardv1 "github.com/pulseguard/pulseguard/gen/pulseguard/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Client wraps the gRPC client with auth token handling.
type Client struct {
	conn   *grpc.ClientConn
	svc    pulseguardv1.AgentServiceClient
	token  string
	target string
}

// New creates a new gRPC client connecting to the given server address.
func New(target, token string) (*Client, error) {
	conn, err := grpc.NewClient(target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", target, err)
	}

	slog.Info("connected to server", "target", target)
	return &Client{
		conn:   conn,
		svc:    pulseguardv1.NewAgentServiceClient(conn),
		token:  token,
		target: target,
	}, nil
}

// Close closes the underlying gRPC connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// authCtx returns a context with the auth token in metadata.
func (c *Client) authCtx(ctx context.Context) context.Context {
	if c.token == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.token)
}

// Register registers this agent with the server.
func (c *Client) Register(ctx context.Context, req *pulseguardv1.RegisterRequest) (*pulseguardv1.RegisterResponse, error) {
	return c.svc.Register(c.authCtx(ctx), req)
}

// Heartbeat sends a heartbeat to the server.
func (c *Client) Heartbeat(ctx context.Context, req *pulseguardv1.HeartbeatRequest) (*pulseguardv1.HeartbeatResponse, error) {
	return c.svc.Heartbeat(c.authCtx(ctx), req)
}

// ReportJobResult reports a job execution result to the server.
func (c *Client) ReportJobResult(ctx context.Context, req *pulseguardv1.ReportJobResultRequest) (*pulseguardv1.ReportJobResultResponse, error) {
	return c.svc.ReportJobResult(c.authCtx(ctx), req)
}

// CommandStream opens a server-streaming connection for receiving commands.
func (c *Client) CommandStream(ctx context.Context, req *pulseguardv1.CommandStreamRequest) (pulseguardv1.AgentService_CommandStreamClient, error) {
	return c.svc.CommandStream(c.authCtx(ctx), req)
}

// ReportDiscoveredJobs reports discovered crontab entries to the server.
func (c *Client) ReportDiscoveredJobs(ctx context.Context, req *pulseguardv1.ReportDiscoveredJobsRequest) (*pulseguardv1.ReportDiscoveredJobsResponse, error) {
	return c.svc.ReportDiscoveredJobs(c.authCtx(ctx), req)
}
