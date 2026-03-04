package discovery

import (
	"bufio"
	"log/slog"
	"os/exec"
	"os/user"
	"strings"

	pulseguardv1 "github.com/pulseguard/pulseguard/gen/pulseguard/v1"
)

// DiscoverCrontab parses the current user's crontab and returns discovered jobs.
func DiscoverCrontab() []*pulseguardv1.DiscoveredCronJob {
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.Output()
	if err != nil {
		slog.Warn("failed to read crontab", "error", err)
		return nil
	}

	currentUser := "unknown"
	if u, err := user.Current(); err == nil {
		currentUser = u.Username
	}

	return parseCrontab(string(output), currentUser)
}

func parseCrontab(content, username string) []*pulseguardv1.DiscoveredCronJob {
	var jobs []*pulseguardv1.DiscoveredCronJob
	var envVars []string

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "=") && !strings.HasPrefix(line, "*") && !strings.HasPrefix(line, "@") {
			if k, _, ok := strings.Cut(line, "="); ok && !strings.Contains(k, " ") && k == strings.ToUpper(k) {
				envVars = append(envVars, line)
			}
			continue
		}

		if strings.Contains(line, "pulseguard-agent wrap") {
			continue
		}

		schedule, command := parseCronLine(line)
		if schedule == "" || command == "" {
			continue
		}

		fullCommand := command
		if len(envVars) > 0 {
			var exports []string
			for _, ev := range envVars {
				exports = append(exports, "export "+ev)
			}
			fullCommand = strings.Join(exports, " && ") + " && " + command
		}

		jobs = append(jobs, &pulseguardv1.DiscoveredCronJob{
			Schedule: schedule,
			Command:  fullCommand,
			User:     username,
			Source:   "crontab",
		})
	}

	return jobs
}

func parseCronLine(line string) (schedule, command string) {
	// Handle @-style schedules
	if strings.HasPrefix(line, "@") {
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			return "", ""
		}
		return parts[0], strings.TrimSpace(parts[1])
	}

	// Standard 5-field cron expression
	fields := strings.Fields(line)
	if len(fields) < 6 {
		return "", ""
	}

	schedule = strings.Join(fields[:5], " ")
	command = strings.Join(fields[5:], " ")
	return schedule, command
}
