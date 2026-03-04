package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pulseguard/pulseguard/internal/models"
)

func runApplyWrapper(args []string) {
	fs := flag.NewFlagSet("apply-wrapper", flag.ExitOnError)
	server := fs.String("server", "", "HTTP server URL (required)")
	token := fs.String("token", os.Getenv("PULSEGUARD_TOKEN"), "Auth token")
	_ = fs.Parse(args)

	if *server == "" {
		slog.Error("--server is required")
		os.Exit(1)
	}

	agentPath, err := os.Executable()
	if err != nil {
		agentPath = "/usr/local/bin/pulseguard-agent"
	}

	jobs, err := fetchDiscoveredJobs(*server, *token)
	if err != nil {
		slog.Error("failed to fetch discovered jobs", "error", err)
		os.Exit(1)
	}

	if len(jobs) == 0 {
		slog.Info("no discovered jobs to wrap")
		return
	}

	crontabBytes, err := exec.Command("crontab", "-l").Output()
	if err != nil {
		slog.Error("failed to read crontab", "error", err)
		os.Exit(1)
	}
	crontab := string(crontabBytes)

	modified := false
	lines := strings.Split(crontab, "\n")
	for _, job := range jobs {
		// Escape single quotes in the command for safe shell quoting
		escapedCmd := strings.ReplaceAll(job.Command, "'", "'\\''")
		wrapLine := fmt.Sprintf("%s wrap --server %s --token %s --job-id %s --command '%s'",
			agentPath, *server, *token, job.ID, escapedCmd)

		cmdCore := stripExportPrefix(job.Command)
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				continue
			}
			if strings.Contains(line, "pulseguard-agent wrap") {
				continue
			}
			if strings.Contains(line, cmdCore) {
				schedule, _ := parseCronFields(trimmed)
				if schedule == "" {
					continue
				}
				lines[i] = schedule + " " + wrapLine
				modified = true
				slog.Info("wrapped job", "job_id", job.ID, "command", cmdCore)
				break
			}
		}
	}

	if !modified {
		slog.Info("no crontab lines were modified")
		return
	}

	newCrontab := strings.Join(lines, "\n")
	cmd := exec.Command("crontab", "-")
	cmd.Stdin = strings.NewReader(newCrontab)
	if err := cmd.Run(); err != nil {
		slog.Error("failed to write crontab", "error", err)
		os.Exit(1)
	}

	slog.Info("crontab updated successfully")

	for _, job := range jobs {
		if err := enableJob(*server, *token, job.ID); err != nil {
			slog.Error("failed to enable job", "job_id", job.ID, "error", err)
		}
	}
}

func stripExportPrefix(cmd string) string {
	for strings.HasPrefix(cmd, "export ") {
		idx := strings.Index(cmd, " && ")
		if idx < 0 {
			break
		}
		cmd = cmd[idx+4:]
	}
	return cmd
}

func parseCronFields(line string) (string, string) {
	if strings.HasPrefix(line, "@") {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			return "", ""
		}
		return parts[0], strings.TrimSpace(parts[1])
	}
	fields := strings.Fields(line)
	if len(fields) < 6 {
		return "", ""
	}
	return strings.Join(fields[:5], " "), strings.Join(fields[5:], " ")
}

func fetchDiscoveredJobs(server, token string) ([]*models.Job, error) {
	url := fmt.Sprintf("%s/api/jobs?source=discovered", strings.TrimRight(server, "/"))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jobs []*models.Job
	if err := json.Unmarshal(body, &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func enableJob(server, token, jobID string) error {
	url := fmt.Sprintf("%s/api/jobs/%s", strings.TrimRight(server, "/"), jobID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var job map[string]interface{}
	if err := json.Unmarshal(body, &job); err != nil {
		return err
	}
	job["enabled"] = true

	putBody, _ := json.Marshal(job)
	putReq, err := http.NewRequest("PUT", url, strings.NewReader(string(putBody)))
	if err != nil {
		return err
	}
	putReq.Header.Set("Content-Type", "application/json")
	if token != "" {
		putReq.Header.Set("Authorization", "Bearer "+token)
	}

	putResp, err := client.Do(putReq)
	if err != nil {
		return err
	}
	putResp.Body.Close()

	if putResp.StatusCode >= 400 {
		return fmt.Errorf("server returned status %d", putResp.StatusCode)
	}
	return nil
}
