package models

import "time"

type Agent struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Hostname        string            `json:"hostname"`
	IPAddress       string            `json:"ip_address"`
	OS              string            `json:"os"`
	Arch            string            `json:"arch"`
	AgentVersion    string            `json:"agent_version"`
	Labels          map[string]string `json:"labels"`
	Status          string            `json:"status"`
	LastHeartbeatAt *time.Time        `json:"last_heartbeat_at,omitempty"`
	RegisteredAt    time.Time         `json:"registered_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}
