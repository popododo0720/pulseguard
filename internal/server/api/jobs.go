package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	pulseguardv1 "github.com/pulseguard/pulseguard/gen/pulseguard/v1"
	"github.com/pulseguard/pulseguard/internal/models"
	servergrpc "github.com/pulseguard/pulseguard/internal/server/grpc"
	"github.com/pulseguard/pulseguard/internal/server/store"
)

type createJobRequest struct {
	AgentID           string                   `json:"agent_id"`
	Name              string                   `json:"name"`
	Schedule          string                   `json:"schedule"`
	Command           string                   `json:"command"`
	WorkingDir        string                   `json:"working_dir"`
	TimeoutSeconds    int                      `json:"timeout_seconds"`
	SuccessConditions models.SuccessConditions `json:"success_conditions"`
	FailurePolicy     models.FailurePolicy     `json:"failure_policy"`
	Env               map[string]string        `json:"env"`
	Enabled           bool                     `json:"enabled"`
}

func listJobsHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		source := r.URL.Query().Get("source")
		var jobs []*models.Job
		var err error
		if source != "" {
			jobs, err = s.ListJobsBySource(source)
		} else {
			jobs, err = s.ListJobs()
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, jobs)
	}
}

func createJobHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createJobRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.TimeoutSeconds == 0 {
			req.TimeoutSeconds = 3600
		}
		if req.Env == nil {
			req.Env = map[string]string{}
		}

		job := &models.Job{
			ID:                uuid.New().String(),
			AgentID:           req.AgentID,
			Name:              req.Name,
			Schedule:          req.Schedule,
			Command:           req.Command,
			WorkingDir:        req.WorkingDir,
			TimeoutSeconds:    req.TimeoutSeconds,
			SuccessConditions: req.SuccessConditions,
			FailurePolicy:     req.FailurePolicy,
			Env:               req.Env,
			Enabled:           req.Enabled,
			Source:            "manual",
			LastStatus:        "unknown",
		}

		if err := s.CreateJob(job); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusCreated, job)
	}
}

func getJobHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		job, err := s.GetJob(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "job not found")
			return
		}
		writeJSON(w, http.StatusOK, job)
	}
}

func updateJobHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		existing, err := s.GetJob(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "job not found")
			return
		}

		var req createJobRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		existing.Name = req.Name
		existing.Schedule = req.Schedule
		existing.Command = req.Command
		existing.WorkingDir = req.WorkingDir
		existing.TimeoutSeconds = req.TimeoutSeconds
		existing.SuccessConditions = req.SuccessConditions
		existing.FailurePolicy = req.FailurePolicy
		existing.Env = req.Env
		existing.Enabled = req.Enabled

		if err := s.UpdateJob(existing); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, existing)
	}
}

func deleteJobHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := s.DeleteJob(id); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

func rerunJobHandler(s *store.Store, dispatcher *servergrpc.CommandDispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		job, err := s.GetJob(id)
		if err != nil {
			writeError(w, http.StatusNotFound, "job not found")
			return
		}

		execID := uuid.New().String()
		cmd := &pulseguardv1.Command{
			CommandId: uuid.New().String(),
			Payload: &pulseguardv1.Command_RunJob{
				RunJob: &pulseguardv1.RunJobCommand{
					JobId:       job.ID,
					ExecutionId: execID,
					Trigger:     "manual",
				},
			},
		}

		if err := dispatcher.Send(job.AgentID, cmd); err != nil {
			writeError(w, http.StatusServiceUnavailable, "agent not connected: "+err.Error())
			return
		}

		writeJSON(w, http.StatusAccepted, map[string]string{
			"execution_id": execID,
			"status":       "dispatched",
		})
	}
}

func listJobExecutionsHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		execs, err := s.ListJobExecutions(id, 100)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, execs)
	}
}

type jobReportRequest struct {
	ExitCode   int    `json:"exit_code"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
	DurationMs int64  `json:"duration_ms"`
	Trigger    string `json:"trigger"`
	Error      string `json:"error"`
}

func reportJobResultHandler(s *store.Store, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if token != "" {
			auth := r.Header.Get("Authorization")
			if auth != "Bearer "+token {
				writeError(w, http.StatusUnauthorized, "invalid or missing token")
				return
			}
		}

		jobID := chi.URLParam(r, "id")
		job, err := s.GetJob(jobID)
		if err != nil {
			writeError(w, http.StatusNotFound, "job not found")
			return
		}

		var req jobReportRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		startedAt, _ := time.Parse(time.RFC3339, req.StartedAt)
		var finishedAt *time.Time
		if req.FinishedAt != "" {
			t, _ := time.Parse(time.RFC3339, req.FinishedAt)
			finishedAt = &t
		}

		resultStatus := evaluateResult(job.SuccessConditions, req)

		execID := uuid.New().String()
		execution := &models.JobExecution{
			ID:         execID,
			JobID:      jobID,
			AgentID:    job.AgentID,
			Status:     resultStatus,
			ExitCode:   &req.ExitCode,
			Stdout:     req.Stdout,
			Stderr:     req.Stderr,
			Error:      req.Error,
			StartedAt:  startedAt,
			FinishedAt: finishedAt,
			DurationMs: &req.DurationMs,
			Trigger:    req.Trigger,
		}

		if err := s.CreateJobExecution(execution); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		_ = s.UpdateJobStatus(jobID, resultStatus)

		shouldRetry := false
		retryDelay := 0
		if resultStatus == "failure" && job.FailurePolicy.MaxRetries > 0 {
			shouldRetry = true
			retryDelay = job.FailurePolicy.RetryDelaySeconds
			if retryDelay == 0 {
				retryDelay = 60
			}
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"execution_id": execID,
			"status":       resultStatus,
			"should_retry": shouldRetry,
			"retry_delay":  retryDelay,
		})
	}
}

func evaluateResult(sc models.SuccessConditions, req jobReportRequest) string {
	if req.Error != "" {
		return "failure"
	}
	if req.ExitCode != sc.ExpectedExitCode {
		return "failure"
	}
	if sc.StdoutContains != "" && !strings.Contains(req.Stdout, sc.StdoutContains) {
		return "failure"
	}
	if sc.StdoutEndswith != "" && !strings.HasSuffix(strings.TrimSpace(req.Stdout), sc.StdoutEndswith) {
		return "failure"
	}
	if sc.StderrEmpty && strings.TrimSpace(req.Stderr) != "" {
		return "failure"
	}
	return "success"
}

func decodeJSON(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
