package api

import (
	"encoding/json"
	"io"
	"net/http"

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
		jobs, err := s.ListJobs()
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

func decodeJSON(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1MB max
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
