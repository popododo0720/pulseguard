package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	pulseguardv1 "github.com/pulseguard/pulseguard/gen/pulseguard/v1"
	"github.com/pulseguard/pulseguard/internal/models"
	"github.com/pulseguard/pulseguard/internal/server/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// CommandDispatcher manages channels to connected agents for pushing commands.
type CommandDispatcher struct {
	mu       sync.RWMutex
	channels map[string]chan *pulseguardv1.Command
}

func NewCommandDispatcher() *CommandDispatcher {
	return &CommandDispatcher{
		channels: make(map[string]chan *pulseguardv1.Command),
	}
}

func (d *CommandDispatcher) Register(agentID string) chan *pulseguardv1.Command {
	d.mu.Lock()
	defer d.mu.Unlock()
	// Close existing channel if any
	if ch, ok := d.channels[agentID]; ok {
		close(ch)
	}
	ch := make(chan *pulseguardv1.Command, 16)
	d.channels[agentID] = ch
	return ch
}

func (d *CommandDispatcher) Unregister(agentID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if ch, ok := d.channels[agentID]; ok {
		close(ch)
		delete(d.channels, agentID)
	}
}

func (d *CommandDispatcher) Send(agentID string, cmd *pulseguardv1.Command) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	ch, ok := d.channels[agentID]
	if !ok {
		return fmt.Errorf("agent %s not connected", agentID)
	}
	select {
	case ch <- cmd:
		return nil
	default:
		return fmt.Errorf("agent %s command channel full", agentID)
	}
}

// Server implements the AgentServiceServer gRPC interface.
type Server struct {
	pulseguardv1.UnimplementedAgentServiceServer
	store      *store.Store
	dispatcher *CommandDispatcher
	token      string // simple auth token
}

func NewServer(s *store.Store, dispatcher *CommandDispatcher, token string) *Server {
	return &Server{
		store:      s,
		dispatcher: dispatcher,
		token:      token,
	}
}

func (s *Server) checkAuth(ctx context.Context) error {
	if s.token == "" {
		return nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}
	tokens := md.Get("authorization")
	if len(tokens) == 0 || tokens[0] != "Bearer "+s.token {
		return status.Error(codes.Unauthenticated, "invalid token")
	}
	return nil
}

func (s *Server) Register(ctx context.Context, req *pulseguardv1.RegisterRequest) (*pulseguardv1.RegisterResponse, error) {
	if err := s.checkAuth(ctx); err != nil {
		return nil, err
	}

	labelsMap := req.GetLabels()
	if labelsMap == nil {
		labelsMap = map[string]string{}
	}

	var agentID string
	existing, err := s.store.GetAgentByHostname(req.GetHostname())
	if err == nil && existing != nil {
		agentID = existing.ID
		_ = s.store.UpdateAgentInfo(agentID, req.GetHostname(), req.GetIpAddress(), req.GetOs(), req.GetArch(), req.GetAgentVersion())
		slog.Info("agent re-registered", "agent_id", agentID, "hostname", req.GetHostname())
	} else {
		agentID = uuid.New().String()
		agent := &models.Agent{
			ID:           agentID,
			Name:         req.GetHostname(),
			Hostname:     req.GetHostname(),
			IPAddress:    req.GetIpAddress(),
			OS:           req.GetOs(),
			Arch:         req.GetArch(),
			AgentVersion: req.GetAgentVersion(),
			Labels:       labelsMap,
			Status:       "online",
		}
		if err := s.store.CreateAgent(agent); err != nil {
			slog.Error("failed to create agent", "error", err)
			return nil, status.Errorf(codes.Internal, "failed to register agent: %v", err)
		}
		slog.Info("agent registered", "agent_id", agentID, "hostname", req.GetHostname())
	}

	jobs, err := s.store.ListJobsByAgent(agentID)
	if err != nil {
		slog.Error("failed to list jobs", "error", err)
	}

	jobConfigs := make([]*pulseguardv1.JobConfig, 0, len(jobs))
	for _, j := range jobs {
		jobConfigs = append(jobConfigs, jobToProto(j))
	}

	return &pulseguardv1.RegisterResponse{
		AgentId:                  agentID,
		HeartbeatIntervalSeconds: 30,
		Jobs:                     jobConfigs,
	}, nil
}

func (s *Server) Heartbeat(ctx context.Context, req *pulseguardv1.HeartbeatRequest) (*pulseguardv1.HeartbeatResponse, error) {
	if err := s.checkAuth(ctx); err != nil {
		return nil, err
	}

	if err := s.store.UpdateAgentHeartbeat(req.GetAgentId()); err != nil {
		slog.Error("failed to update heartbeat", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to update heartbeat: %v", err)
	}

	// In v1, no config changes are pushed via heartbeat
	return &pulseguardv1.HeartbeatResponse{}, nil
}

func (s *Server) ReportJobResult(ctx context.Context, req *pulseguardv1.ReportJobResultRequest) (*pulseguardv1.ReportJobResultResponse, error) {
	if err := s.checkAuth(ctx); err != nil {
		return nil, err
	}

	// Get the job to evaluate success conditions
	job, err := s.store.GetJob(req.GetJobId())
	if err != nil {
		slog.Error("failed to get job for result evaluation", "error", err)
		return nil, status.Errorf(codes.NotFound, "job not found: %v", err)
	}

	// Evaluate success conditions
	resultStatus := evaluateSuccessConditions(job.SuccessConditions, req)

	startedAt := time.Unix(req.GetStartedAtUnix(), 0).UTC()
	var finishedAt *time.Time
	if req.GetFinishedAtUnix() > 0 {
		t := time.Unix(req.GetFinishedAtUnix(), 0).UTC()
		finishedAt = &t
	}

	exitCode := int(req.GetExitCode())
	durationMs := req.GetDurationMs()

	exec := &models.JobExecution{
		ID:         req.GetExecutionId(),
		JobID:      req.GetJobId(),
		AgentID:    req.GetAgentId(),
		Status:     resultStatus,
		ExitCode:   &exitCode,
		Stdout:     req.GetStdout(),
		Stderr:     req.GetStderr(),
		Error:      req.GetError(),
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		DurationMs: &durationMs,
		Trigger:    req.GetTrigger(),
		RetryCount: int(req.GetRetryCount()),
	}

	if err := s.store.CreateJobExecution(exec); err != nil {
		slog.Error("failed to save execution", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to save execution: %v", err)
	}

	// Update job status
	_ = s.store.UpdateJobStatus(req.GetJobId(), resultStatus)

	// Determine retry
	shouldRetry := false
	retryDelay := int32(0)
	if resultStatus == "failure" && job.FailurePolicy.MaxRetries > 0 && int(req.GetRetryCount()) < job.FailurePolicy.MaxRetries {
		shouldRetry = true
		retryDelay = int32(job.FailurePolicy.RetryDelaySeconds)
		if retryDelay == 0 {
			retryDelay = 60
		}
	}

	slog.Info("job result reported", "job_id", req.GetJobId(), "status", resultStatus, "should_retry", shouldRetry)
	return &pulseguardv1.ReportJobResultResponse{
		Status:            resultStatus,
		ShouldRetry:       shouldRetry,
		RetryDelaySeconds: retryDelay,
	}, nil
}

func (s *Server) CommandStream(req *pulseguardv1.CommandStreamRequest, stream grpc.ServerStreamingServer[pulseguardv1.Command]) error {
	agentID := req.GetAgentId()
	slog.Info("agent connected to command stream", "agent_id", agentID)

	ch := s.dispatcher.Register(agentID)
	defer s.dispatcher.Unregister(agentID)

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			slog.Info("agent disconnected from command stream", "agent_id", agentID)
			return nil
		case cmd, ok := <-ch:
			if !ok {
				return nil
			}
			if err := stream.Send(cmd); err != nil {
				slog.Error("failed to send command", "agent_id", agentID, "error", err)
				return err
			}
		}
	}
}

func (s *Server) ReportDiscoveredJobs(ctx context.Context, req *pulseguardv1.ReportDiscoveredJobsRequest) (*pulseguardv1.ReportDiscoveredJobsResponse, error) {
	if err := s.checkAuth(ctx); err != nil {
		return nil, err
	}

	imported := int32(0)
	for _, dj := range req.GetJobs() {
		exists, err := s.store.JobExistsByCommand(req.GetAgentId(), dj.GetCommand())
		if err != nil {
			slog.Error("failed to check existing job", "error", err)
			continue
		}
		if exists {
			continue
		}

		job := &models.Job{
			ID:             uuid.New().String(),
			AgentID:        req.GetAgentId(),
			Name:           extractJobName(dj.GetCommand(), dj.GetSchedule()),
			Schedule:       dj.GetSchedule(),
			Command:        dj.GetCommand(),
			TimeoutSeconds: 3600,
			SuccessConditions: models.SuccessConditions{
				ExpectedExitCode: 0,
			},
			FailurePolicy: models.FailurePolicy{
				MaxRetries:        0,
				RetryDelaySeconds: 60,
			},
			Env:        map[string]string{},
			Enabled:    false,
			Source:     "discovered",
			LastStatus: "unknown",
		}

		if err := s.store.CreateJob(job); err != nil {
			slog.Error("failed to save discovered job", "error", err)
			continue
		}
		imported++
	}

	slog.Info("discovered jobs imported", "agent_id", req.GetAgentId(), "count", imported)
	return &pulseguardv1.ReportDiscoveredJobsResponse{
		ImportedCount: imported,
	}, nil
}

// Dispatcher returns the command dispatcher for use by REST API.
func (s *Server) Dispatcher() *CommandDispatcher {
	return s.dispatcher
}

// evaluateSuccessConditions checks all conditions; ALL must pass for success.
func evaluateSuccessConditions(sc models.SuccessConditions, req *pulseguardv1.ReportJobResultRequest) string {
	// If there's a system-level error, it's always a failure
	if req.GetError() != "" {
		return "failure"
	}

	// Check exit code
	if int(req.GetExitCode()) != sc.ExpectedExitCode {
		return "failure"
	}

	// Check stdout_contains
	if sc.StdoutContains != "" && !strings.Contains(req.GetStdout(), sc.StdoutContains) {
		return "failure"
	}

	// Check stdout_endswith
	if sc.StdoutEndswith != "" && !strings.HasSuffix(strings.TrimSpace(req.GetStdout()), sc.StdoutEndswith) {
		return "failure"
	}

	// Check stderr_empty
	if sc.StderrEmpty && strings.TrimSpace(req.GetStderr()) != "" {
		return "failure"
	}

	return "success"
}

func jobToProto(j *models.Job) *pulseguardv1.JobConfig {
	env := j.Env
	if env == nil {
		env = map[string]string{}
	}
	fpChannels := j.FailurePolicy.NotifyChannels
	if fpChannels == nil {
		fpChannels = []string{}
	}
	return &pulseguardv1.JobConfig{
		Id:             j.ID,
		Name:           j.Name,
		Schedule:       j.Schedule,
		Command:        j.Command,
		WorkingDir:     j.WorkingDir,
		TimeoutSeconds: int32(j.TimeoutSeconds),
		SuccessConditions: &pulseguardv1.SuccessConditions{
			ExpectedExitCode: int32(j.SuccessConditions.ExpectedExitCode),
			StdoutContains:   j.SuccessConditions.StdoutContains,
			StdoutEndswith:   j.SuccessConditions.StdoutEndswith,
			StderrEmpty:      j.SuccessConditions.StderrEmpty,
			FileExists:       j.SuccessConditions.FileExists,
		},
		FailurePolicy: &pulseguardv1.FailurePolicy{
			MaxRetries:        int32(j.FailurePolicy.MaxRetries),
			RetryDelaySeconds: int32(j.FailurePolicy.RetryDelaySeconds),
			NotifyChannels:    fpChannels,
		},
		Enabled: j.Enabled,
		Env:     env,
	}
}

// extractJobName creates a human-readable name from a cron command.
// Examples:
//
//	"cd /opt/scripts && python metrics.py" → "metrics.py"
//	"/usr/bin/backup.sh --full" → "backup.sh"
//	"SCRIPT_DIR=/opt && cd $SCRIPT_DIR && ./run.sh" → "run.sh"
func extractJobName(command, schedule string) string {
	// Strip output redirection (everything after >> or >)
	cmd := command
	for _, op := range []string{">>", ">"} {
		if idx := strings.Index(cmd, " "+op+" "); idx >= 0 {
			cmd = cmd[:idx]
		}
	}

	parts := strings.Fields(cmd)
	var candidate string
	for i := len(parts) - 1; i >= 0; i-- {
		p := parts[i]
		// Skip flags, env vars, operators, redirections
		if strings.HasPrefix(p, "-") || strings.Contains(p, "=") ||
			p == "&&" || p == "||" || p == "|" || p == ";" ||
			strings.HasPrefix(p, ">") || strings.HasPrefix(p, "<") ||
			strings.HasPrefix(p, "2>") || p == "2>&1" {
			continue
		}
		// Skip shell variable references
		if strings.HasPrefix(p, "$") {
			continue
		}
		// Skip common shell commands and interpreters
		if p == "cd" || p == "sh" || p == "bash" || p == "python" || p == "python3" ||
			p == "node" || p == "perl" || p == "ruby" || p == "uv" || p == "run" {
			continue
		}
		// Check if it looks like a script filename
		base := filepath.Base(p)
		if strings.Contains(base, ".") || strings.HasPrefix(p, "/") || strings.HasPrefix(p, "./") {
			candidate = base
			break
		}
		if candidate == "" {
			candidate = base
		}
	}
	if candidate == "" {
		candidate = truncate(command, 40)
	}
	if schedule != "" {
		return fmt.Sprintf("%s (cron: %s)", candidate, schedule)
	}
	return candidate
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// MarshalLabels is a helper exported for use by API handlers.
func MarshalLabels(m map[string]string) string {
	b, _ := json.Marshal(m)
	return string(b)
}
