package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	pulseguardv1 "github.com/pulseguard/pulseguard/gen/pulseguard/v1"
	"github.com/pulseguard/pulseguard/internal/agent/client"
	"github.com/pulseguard/pulseguard/internal/agent/discovery"
	"github.com/pulseguard/pulseguard/internal/agent/executor"
)

func main() {
	// Check for subcommands before flag.Parse()
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "wrap":
			runWrap(os.Args[2:])
			return
		case "apply-wrapper":
			runApplyWrapper(os.Args[2:])
			return
		}
	}

	// Daemon mode (original behavior)
	server := flag.String("server", "", "gRPC server address (required, e.g., localhost:9090)")
	token := flag.String("token", os.Getenv("PULSEGUARD_TOKEN"), "Auth token")
	discover := flag.Bool("discover", false, "Discover and report crontab entries on startup")
	flag.Parse()

	if *server == "" {
		slog.Error("--server flag is required")
		os.Exit(1)
	}

	// Connect to server
	c, err := client.New(*server, *token)
	if err != nil {
		slog.Error("failed to connect to server", "error", err)
		os.Exit(1)
	}
	defer c.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	hostname, _ := os.Hostname()

	// Log saved agent ID if available
	if savedID := loadSavedAgentID(); savedID != "" {
		slog.Info("loaded saved agent ID", "agent_id", savedID)
	}

	// Register with server
	regResp, err := c.Register(ctx, &pulseguardv1.RegisterRequest{
		Hostname:     hostname,
		Os:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		AgentVersion: "0.1.0",
		Labels:       map[string]string{},
	})
	if err != nil {
		slog.Error("failed to register with server", "error", err)
		os.Exit(1)
	}

	agentID := regResp.GetAgentId()
	heartbeatInterval := time.Duration(regResp.GetHeartbeatIntervalSeconds()) * time.Second
	if heartbeatInterval == 0 {
		heartbeatInterval = 30 * time.Second
	}

	// Persist server-assigned agent ID for informational purposes
	saveAgentID(agentID)

	slog.Info("registered with server",
		"agent_id", agentID,
		"heartbeat_interval", heartbeatInterval,
		"jobs", len(regResp.GetJobs()),
	)

	// Job config map for command stream lookups
	var jobsMu sync.RWMutex
	jobConfigs := make(map[string]*pulseguardv1.JobConfig)

	// executeAndReport runs a job via executor and reports the result.
	executeAndReport := func(job *pulseguardv1.JobConfig, execID, trigger string, retryCount int) {
		slog.Info("executing job", "job_id", job.GetId(), "name", job.GetName(), "trigger", trigger)

		env := job.GetEnv()
		if env == nil {
			env = map[string]string{}
		}

		result := executor.Run(
			context.Background(),
			job.GetCommand(),
			job.GetWorkingDir(),
			env,
			int(job.GetTimeoutSeconds()),
		)

		req := &pulseguardv1.ReportJobResultRequest{
			AgentId:        agentID,
			JobId:          job.GetId(),
			ExecutionId:    execID,
			ExitCode:       int32(result.ExitCode),
			Stdout:         result.Stdout,
			Stderr:         result.Stderr,
			StartedAtUnix:  result.StartedAt.Unix(),
			FinishedAtUnix: result.FinishedAt.Unix(),
			DurationMs:     result.DurationMs,
			Trigger:        trigger,
			RetryCount:     int32(retryCount),
			Error:          result.Error,
		}

		reportFn(c, req, &jobsMu, jobConfigs)
	}

	// Store initial job configs (no scheduling — cron wrapper handles execution)
	for _, job := range regResp.GetJobs() {
		jobsMu.Lock()
		jobConfigs[job.GetId()] = job
		jobsMu.Unlock()
	}

	// Discover crontab entries
	if *discover {
		go func() {
			discovered := discovery.DiscoverCrontab()
			if len(discovered) > 0 {
				resp, err := c.ReportDiscoveredJobs(context.Background(), &pulseguardv1.ReportDiscoveredJobsRequest{
					AgentId: agentID,
					Jobs:    discovered,
				})
				if err != nil {
					slog.Error("failed to report discovered jobs", "error", err)
				} else {
					slog.Info("discovered jobs reported", "imported", resp.GetImportedCount())
				}
			}
		}()
	}

	// Heartbeat loop
	go func() {
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				_, err := c.Heartbeat(ctx, &pulseguardv1.HeartbeatRequest{
					AgentId: agentID,
				})
				if err != nil {
					slog.Error("heartbeat failed", "error", err)
				}
			}
		}
	}()

	// Command stream listener
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			stream, err := c.CommandStream(ctx, &pulseguardv1.CommandStreamRequest{
				AgentId: agentID,
			})
			if err != nil {
				slog.Error("failed to open command stream", "error", err)
				time.Sleep(5 * time.Second)
				continue
			}

			slog.Info("command stream connected")

			for {
				cmd, err := stream.Recv()
				if err != nil {
					slog.Warn("command stream disconnected", "error", err)
					break
				}

				switch p := cmd.GetPayload().(type) {
				case *pulseguardv1.Command_RunJob:
					slog.Info("received run job command", "job_id", p.RunJob.GetJobId())
					jobsMu.RLock()
					job, ok := jobConfigs[p.RunJob.GetJobId()]
					jobsMu.RUnlock()
					if ok {
						go executeAndReport(job, p.RunJob.GetExecutionId(), p.RunJob.GetTrigger(), 0)
					} else {
						slog.Warn("received run command for unknown job", "job_id", p.RunJob.GetJobId())
					}

				case *pulseguardv1.Command_StopJob:
					slog.Info("received stop job command", "job_id", p.StopJob.GetJobId())
					// In v1, we don't support stopping running jobs

				case *pulseguardv1.Command_UpdateConfig:
					slog.Info("received config update", "jobs", len(p.UpdateConfig.GetJobs()))
					for _, job := range p.UpdateConfig.GetJobs() {
						jobsMu.Lock()
						jobConfigs[job.GetId()] = job
						jobsMu.Unlock()
					}
				}
			}

			// Reconnect delay
			time.Sleep(5 * time.Second)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down agent...")
}

// reportFn sends job results to server and handles retries.
func reportFn(c *client.Client, req *pulseguardv1.ReportJobResultRequest, jobsMu *sync.RWMutex, jobConfigs map[string]*pulseguardv1.JobConfig) {
	resp, err := c.ReportJobResult(context.Background(), req)
	if err != nil {
		slog.Error("failed to report job result", "error", err)
		return
	}
	slog.Info("job result reported",
		"job_id", req.GetJobId(),
		"status", resp.GetStatus(),
		"should_retry", resp.GetShouldRetry(),
	)

	// Handle retry
	if resp.GetShouldRetry() {
		delay := time.Duration(resp.GetRetryDelaySeconds()) * time.Second
		slog.Info("retrying job after delay", "job_id", req.GetJobId(), "delay", delay)
		time.Sleep(delay)

		jobsMu.RLock()
		job, ok := jobConfigs[req.GetJobId()]
		jobsMu.RUnlock()
		if ok {
			env := job.GetEnv()
			if env == nil {
				env = map[string]string{}
			}
			go func() {
				result := executor.Run(
					context.Background(),
					job.GetCommand(),
					job.GetWorkingDir(),
					env,
					int(job.GetTimeoutSeconds()),
				)
				retryReq := &pulseguardv1.ReportJobResultRequest{
					AgentId:        req.GetAgentId(),
					JobId:          job.GetId(),
					ExecutionId:    req.GetExecutionId() + "-retry",
					ExitCode:       int32(result.ExitCode),
					Stdout:         result.Stdout,
					Stderr:         result.Stderr,
					StartedAtUnix:  result.StartedAt.Unix(),
					FinishedAtUnix: result.FinishedAt.Unix(),
					DurationMs:     result.DurationMs,
					Trigger:        "retry",
					RetryCount:     req.GetRetryCount() + 1,
					Error:          result.Error,
				}
				reportFn(c, retryReq, jobsMu, jobConfigs)
			}()
		}
	}
}

const agentIDFile = "/var/lib/pulseguard/agent-id"

// loadSavedAgentID reads the previously saved agent ID from disk.
func loadSavedAgentID() string {
	data, err := os.ReadFile(agentIDFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// saveAgentID persists the server-assigned agent ID to disk.
func saveAgentID(id string) {
	if err := os.MkdirAll(filepath.Dir(agentIDFile), 0755); err != nil {
		slog.Warn("failed to create agent ID directory", "error", err)
		return
	}
	if err := os.WriteFile(agentIDFile, []byte(id), 0644); err != nil {
		slog.Warn("failed to save agent ID", "error", err)
	}
}
