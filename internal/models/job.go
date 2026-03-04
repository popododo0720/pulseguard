package models

import "time"

type SuccessConditions struct {
	ExpectedExitCode int    `json:"expected_exit_code"`
	StdoutContains   string `json:"stdout_contains,omitempty"`
	StdoutEndswith   string `json:"stdout_endswith,omitempty"`
	StderrEmpty      bool   `json:"stderr_empty,omitempty"`
	FileExists       string `json:"file_exists,omitempty"`
}

type FailurePolicy struct {
	MaxRetries        int      `json:"max_retries"`
	RetryDelaySeconds int      `json:"retry_delay_seconds"`
	NotifyChannels    []string `json:"notify_channels,omitempty"`
}

type Job struct {
	ID                string            `json:"id"`
	AgentID           string            `json:"agent_id"`
	Name              string            `json:"name"`
	Schedule          string            `json:"schedule"`
	Command           string            `json:"command"`
	WorkingDir        string            `json:"working_dir"`
	TimeoutSeconds    int               `json:"timeout_seconds"`
	SuccessConditions SuccessConditions `json:"success_conditions"`
	FailurePolicy     FailurePolicy     `json:"failure_policy"`
	Env               map[string]string `json:"env"`
	Enabled           bool              `json:"enabled"`
	Source            string            `json:"source"`
	LastStatus        string            `json:"last_status"`
	LastRunAt         *time.Time        `json:"last_run_at,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

type JobExecution struct {
	ID         string     `json:"id"`
	JobID      string     `json:"job_id"`
	AgentID    string     `json:"agent_id"`
	Status     string     `json:"status"`
	ExitCode   *int       `json:"exit_code,omitempty"`
	Stdout     string     `json:"stdout"`
	Stderr     string     `json:"stderr"`
	Error      string     `json:"error"`
	StartedAt  time.Time  `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	DurationMs *int64     `json:"duration_ms,omitempty"`
	Trigger    string     `json:"trigger"`
	RetryCount int        `json:"retry_count"`
	CreatedAt  time.Time  `json:"created_at"`
}
