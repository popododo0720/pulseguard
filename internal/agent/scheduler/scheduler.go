package scheduler

import (
	"context"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	pulseguardv1 "github.com/pulseguard/pulseguard/gen/pulseguard/v1"
	"github.com/pulseguard/pulseguard/internal/agent/executor"
	"github.com/robfig/cron/v3"
)

// ReportFunc is called when a job execution completes.
type ReportFunc func(req *pulseguardv1.ReportJobResultRequest)

// Scheduler manages cron-scheduled job executions.
type Scheduler struct {
	cron     *cron.Cron
	agentID  string
	report   ReportFunc
	mu       sync.Mutex
	entryMap map[string]cron.EntryID // jobID -> cron entry ID
}

// New creates a new scheduler.
func New(agentID string, reportFn ReportFunc) *Scheduler {
	return &Scheduler{
		cron:     cron.New(cron.WithSeconds()),
		agentID:  agentID,
		report:   reportFn,
		entryMap: make(map[string]cron.EntryID),
	}
}

// AddJob adds a job to the scheduler.
func (s *Scheduler) AddJob(job *pulseguardv1.JobConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove existing entry if any
	if entryID, ok := s.entryMap[job.GetId()]; ok {
		s.cron.Remove(entryID)
		delete(s.entryMap, job.GetId())
	}

	if !job.GetEnabled() {
		return nil
	}

	// robfig/cron expects 5-field cron expressions by default;
	// with WithSeconds() it expects 6-field. Standard crontab is 5-field,
	// so we prepend "0 " for the seconds field.
	schedule := "0 " + job.GetSchedule()

	entryID, err := s.cron.AddFunc(schedule, func() {
		s.runJob(job, "scheduled")
	})
	if err != nil {
		slog.Error("failed to add job to scheduler", "job_id", job.GetId(), "schedule", job.GetSchedule(), "error", err)
		return err
	}

	s.entryMap[job.GetId()] = entryID
	slog.Info("job scheduled", "job_id", job.GetId(), "name", job.GetName(), "schedule", job.GetSchedule())
	return nil
}

// RemoveJob removes a job from the scheduler.
func (s *Scheduler) RemoveJob(jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, ok := s.entryMap[jobID]; ok {
		s.cron.Remove(entryID)
		delete(s.entryMap, jobID)
	}
}

// RunJobNow triggers a job execution immediately (for manual/retry triggers).
func (s *Scheduler) RunJobNow(job *pulseguardv1.JobConfig, execID, trigger string) {
	go s.runJobWithID(job, execID, trigger, 0)
}

// Start begins the cron scheduler.
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop stops the cron scheduler.
func (s *Scheduler) Stop() {
	s.cron.Stop()
}

func (s *Scheduler) runJob(job *pulseguardv1.JobConfig, trigger string) {
	execID := uuid.New().String()
	s.runJobWithID(job, execID, trigger, 0)
}

func (s *Scheduler) runJobWithID(job *pulseguardv1.JobConfig, execID, trigger string, retryCount int) {
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
		AgentId:        s.agentID,
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

	s.report(req)
}

// Jobs returns the map of tracked job IDs for lookup.
func (s *Scheduler) Jobs() map[string]cron.EntryID {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make(map[string]cron.EntryID, len(s.entryMap))
	for k, v := range s.entryMap {
		cp[k] = v
	}
	return cp
}
