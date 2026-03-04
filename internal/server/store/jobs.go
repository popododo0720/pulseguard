package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pulseguard/pulseguard/internal/models"
)

func (s *Store) CreateJob(j *models.Job) error {
	scJSON, _ := json.Marshal(j.SuccessConditions)
	fpJSON, _ := json.Marshal(j.FailurePolicy)
	envJSON, _ := json.Marshal(j.Env)

	_, err := s.db.Exec(`
		INSERT INTO jobs (id, agent_id, name, schedule, command, working_dir, timeout_seconds, success_conditions, failure_policy, env, enabled, source, last_status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`,
		j.ID, j.AgentID, j.Name, j.Schedule, j.Command, j.WorkingDir, j.TimeoutSeconds,
		string(scJSON), string(fpJSON), string(envJSON), boolToInt(j.Enabled), j.Source, j.LastStatus,
	)
	return err
}

func (s *Store) GetJob(id string) (*models.Job, error) {
	row := s.db.QueryRow(`SELECT id, agent_id, name, schedule, command, working_dir, timeout_seconds, success_conditions, failure_policy, env, enabled, source, last_status, last_run_at, created_at, updated_at FROM jobs WHERE id = ?`, id)
	return scanJob(row)
}

func (s *Store) ListJobs() ([]*models.Job, error) {
	rows, err := s.db.Query(`SELECT id, agent_id, name, schedule, command, working_dir, timeout_seconds, success_conditions, failure_policy, env, enabled, source, last_status, last_run_at, created_at, updated_at FROM jobs ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		j, err := scanJobRows(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

func (s *Store) JobExistsByCommand(agentID, command string) (bool, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM jobs WHERE command = ? AND (agent_id = ? OR source = 'discovered')`, command, agentID).Scan(&count)
	return count > 0, err
}

func (s *Store) ListJobsBySource(source string) ([]*models.Job, error) {
	rows, err := s.db.Query(`SELECT id, agent_id, name, schedule, command, working_dir, timeout_seconds, success_conditions, failure_policy, env, enabled, source, last_status, last_run_at, created_at, updated_at FROM jobs WHERE source = ? ORDER BY created_at DESC`, source)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		j, err := scanJobRows(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

func (s *Store) ListJobsByAgent(agentID string) ([]*models.Job, error) {
	rows, err := s.db.Query(`SELECT id, agent_id, name, schedule, command, working_dir, timeout_seconds, success_conditions, failure_policy, env, enabled, source, last_status, last_run_at, created_at, updated_at FROM jobs WHERE agent_id = ? ORDER BY created_at DESC`, agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*models.Job
	for rows.Next() {
		j, err := scanJobRows(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, rows.Err()
}

func (s *Store) UpdateJob(j *models.Job) error {
	scJSON, _ := json.Marshal(j.SuccessConditions)
	fpJSON, _ := json.Marshal(j.FailurePolicy)
	envJSON, _ := json.Marshal(j.Env)

	_, err := s.db.Exec(`
		UPDATE jobs SET name=?, schedule=?, command=?, working_dir=?, timeout_seconds=?, success_conditions=?, failure_policy=?, env=?, enabled=?, updated_at=datetime('now')
		WHERE id=?`,
		j.Name, j.Schedule, j.Command, j.WorkingDir, j.TimeoutSeconds,
		string(scJSON), string(fpJSON), string(envJSON), boolToInt(j.Enabled), j.ID,
	)
	return err
}

func (s *Store) DeleteJob(id string) error {
	_, err := s.db.Exec(`DELETE FROM jobs WHERE id = ?`, id)
	return err
}

func (s *Store) UpdateJobStatus(id, status string) error {
	_, err := s.db.Exec(`UPDATE jobs SET last_status=?, last_run_at=datetime('now'), updated_at=datetime('now') WHERE id=?`, status, id)
	return err
}

func (s *Store) JobCount() (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM jobs`).Scan(&count)
	return count, err
}

// Executions

func (s *Store) CreateJobExecution(e *models.JobExecution) error {
	_, err := s.db.Exec(`
		INSERT INTO job_executions (id, job_id, agent_id, status, exit_code, stdout, stderr, error, started_at, finished_at, duration_ms, trigger, retry_count, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))`,
		e.ID, e.JobID, e.AgentID, e.Status, e.ExitCode, e.Stdout, e.Stderr, e.Error,
		e.StartedAt.UTC().Format("2006-01-02 15:04:05"),
		nullTimeStr(e.FinishedAt),
		e.DurationMs, e.Trigger, e.RetryCount,
	)
	return err
}

func (s *Store) ListJobExecutions(jobID string, limit int) ([]*models.JobExecution, error) {
	rows, err := s.db.Query(`SELECT id, job_id, agent_id, status, exit_code, stdout, stderr, error, started_at, finished_at, duration_ms, trigger, retry_count, created_at FROM job_executions WHERE job_id = ? ORDER BY started_at DESC LIMIT ?`, jobID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var execs []*models.JobExecution
	for rows.Next() {
		e, err := scanExecRows(rows)
		if err != nil {
			return nil, err
		}
		execs = append(execs, e)
	}
	return execs, rows.Err()
}

func (s *Store) RecentFailures(limit int) ([]*models.JobExecution, error) {
	rows, err := s.db.Query(`SELECT id, job_id, agent_id, status, exit_code, stdout, stderr, error, started_at, finished_at, duration_ms, trigger, retry_count, created_at FROM job_executions WHERE status = 'failure' ORDER BY started_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var execs []*models.JobExecution
	for rows.Next() {
		e, err := scanExecRows(rows)
		if err != nil {
			return nil, err
		}
		execs = append(execs, e)
	}
	return execs, rows.Err()
}

func (s *Store) SuccessRate() (float64, error) {
	var total, success int
	err := s.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(CASE WHEN status='success' THEN 1 ELSE 0 END), 0) FROM job_executions WHERE status IN ('success', 'failure')`).Scan(&total, &success)
	if err != nil || total == 0 {
		return 0, err
	}
	return float64(success) / float64(total) * 100, nil
}

// Helpers

func scanJob(row *sql.Row) (*models.Job, error) {
	var j models.Job
	var scJSON, fpJSON, envJSON string
	var enabled int
	var lastRun sql.NullString
	var createdAt, updatedAt string

	err := row.Scan(&j.ID, &j.AgentID, &j.Name, &j.Schedule, &j.Command, &j.WorkingDir, &j.TimeoutSeconds, &scJSON, &fpJSON, &envJSON, &enabled, &j.Source, &j.LastStatus, &lastRun, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	return populateJob(&j, scJSON, fpJSON, envJSON, enabled, lastRun, createdAt, updatedAt), nil
}

func scanJobRows(rows *sql.Rows) (*models.Job, error) {
	var j models.Job
	var scJSON, fpJSON, envJSON string
	var enabled int
	var lastRun sql.NullString
	var createdAt, updatedAt string

	err := rows.Scan(&j.ID, &j.AgentID, &j.Name, &j.Schedule, &j.Command, &j.WorkingDir, &j.TimeoutSeconds, &scJSON, &fpJSON, &envJSON, &enabled, &j.Source, &j.LastStatus, &lastRun, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	return populateJob(&j, scJSON, fpJSON, envJSON, enabled, lastRun, createdAt, updatedAt), nil
}

func populateJob(j *models.Job, scJSON, fpJSON, envJSON string, enabled int, lastRun sql.NullString, createdAt, updatedAt string) *models.Job {
	_ = json.Unmarshal([]byte(scJSON), &j.SuccessConditions)
	_ = json.Unmarshal([]byte(fpJSON), &j.FailurePolicy)
	if err := json.Unmarshal([]byte(envJSON), &j.Env); err != nil {
		j.Env = map[string]string{}
	}
	j.Enabled = enabled == 1
	if lastRun.Valid {
		t := parseDBTime(lastRun.String)
		if !t.IsZero() {
			j.LastRunAt = &t
		}
	}
	j.CreatedAt = parseDBTime(createdAt)
	j.UpdatedAt = parseDBTime(updatedAt)
	return j
}

func scanExecRows(rows *sql.Rows) (*models.JobExecution, error) {
	var e models.JobExecution
	var exitCode sql.NullInt64
	var startedAt, createdAt string
	var finishedAt sql.NullString
	var durationMs sql.NullInt64

	err := rows.Scan(&e.ID, &e.JobID, &e.AgentID, &e.Status, &exitCode, &e.Stdout, &e.Stderr, &e.Error, &startedAt, &finishedAt, &durationMs, &e.Trigger, &e.RetryCount, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("scan execution: %w", err)
	}

	if exitCode.Valid {
		ec := int(exitCode.Int64)
		e.ExitCode = &ec
	}
	e.StartedAt = parseDBTime(startedAt)
	e.CreatedAt = parseDBTime(createdAt)
	if finishedAt.Valid {
		t := parseDBTime(finishedAt.String)
		if !t.IsZero() {
			e.FinishedAt = &t
		}
	}
	if durationMs.Valid {
		d := durationMs.Int64
		e.DurationMs = &d
	}
	return &e, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func nullTimeStr(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.UTC().Format("2006-01-02 15:04:05")
}

var timeFormats = []string{
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05Z07:00",
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02 15:04:05+00:00",
	"2006-01-02 15:04:05-07:00",
}

func parseDBTime(s string) time.Time {
	for _, f := range timeFormats {
		if t, err := time.Parse(f, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}
