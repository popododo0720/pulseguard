package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// Result holds the output of an executed command.
type Result struct {
	ExitCode   int
	Stdout     string
	Stderr     string
	Error      string
	StartedAt  time.Time
	FinishedAt time.Time
	DurationMs int64
}

// Run executes a shell command with the given timeout and working directory.
func Run(ctx context.Context, command, workingDir string, env map[string]string, timeoutSec int) *Result {
	if timeoutSec <= 0 {
		timeoutSec = 3600
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	// Set environment variables
	if len(env) > 0 {
		cmd.Env = cmd.Environ()
		for k, v := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startedAt := time.Now().UTC()
	err := cmd.Run()
	finishedAt := time.Now().UTC()

	result := &Result{
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		DurationMs: finishedAt.Sub(startedAt).Milliseconds(),
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
			result.Error = err.Error()
		}
	}

	return result
}
