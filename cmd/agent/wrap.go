package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pulseguard/pulseguard/internal/agent/executor"
)

type reportRequest struct {
	ExitCode   int    `json:"exit_code"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
	DurationMs int64  `json:"duration_ms"`
	Trigger    string `json:"trigger"`
	Error      string `json:"error"`
}

func runWrap(args []string) {
	fs := flag.NewFlagSet("wrap", flag.ExitOnError)
	server := fs.String("server", "", "HTTP server URL (required)")
	token := fs.String("token", os.Getenv("PULSEGUARD_TOKEN"), "Auth token")
	jobID := fs.String("job-id", "", "Job ID to report (required)")
	cmdFlag := fs.String("command", "", "Command to execute (alternative to positional args after --)")
	_ = fs.Parse(args)

	if *server == "" || *jobID == "" {
		slog.Error("--server and --job-id are required")
		os.Exit(1)
	}

	var command string
	if *cmdFlag != "" {
		command = *cmdFlag
	} else {
		cmdArgs := fs.Args()
		if len(cmdArgs) == 0 {
			slog.Error("no command specified (use --command or -- args)")
			os.Exit(1)
		}
		command = strings.Join(cmdArgs, " ")
	}
	result := executor.Run(context.Background(), command, "", nil, 0)

	report := reportRequest{
		ExitCode:   result.ExitCode,
		Stdout:     result.Stdout,
		Stderr:     result.Stderr,
		StartedAt:  result.StartedAt.Format(time.RFC3339),
		FinishedAt: result.FinishedAt.Format(time.RFC3339),
		DurationMs: result.DurationMs,
		Trigger:    "cron-wrapper",
		Error:      result.Error,
	}

	body, _ := json.Marshal(report)
	url := fmt.Sprintf("%s/api/jobs/%s/report", strings.TrimRight(*server, "/"), *jobID)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		slog.Error("failed to create request", "error", err)
		os.Exit(result.ExitCode)
	}
	req.Header.Set("Content-Type", "application/json")
	if *token != "" {
		req.Header.Set("Authorization", "Bearer "+*token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("failed to report result", "error", err)
	} else {
		resp.Body.Close()
		if resp.StatusCode >= 400 {
			slog.Error("server returned error", "status", resp.StatusCode)
		}
	}

	os.Exit(result.ExitCode)
}
