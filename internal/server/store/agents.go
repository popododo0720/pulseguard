package store

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/pulseguard/pulseguard/internal/models"
)

func (s *Store) CreateAgent(a *models.Agent) error {
	labelsJSON, err := json.Marshal(a.Labels)
	if err != nil {
		return fmt.Errorf("marshal labels: %w", err)
	}
	_, err = s.db.Exec(`
		INSERT INTO agents (id, name, hostname, ip_address, os, arch, agent_version, labels, status, registered_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`,
		a.ID, a.Name, a.Hostname, a.IPAddress, a.OS, a.Arch, a.AgentVersion, string(labelsJSON), a.Status,
	)
	return err
}

func (s *Store) GetAgent(id string) (*models.Agent, error) {
	row := s.db.QueryRow(`SELECT id, name, hostname, ip_address, os, arch, agent_version, labels, status, last_heartbeat_at, registered_at, updated_at FROM agents WHERE id = ?`, id)
	return scanAgent(row)
}

func (s *Store) GetAgentByHostname(hostname string) (*models.Agent, error) {
	row := s.db.QueryRow(`SELECT id, name, hostname, ip_address, os, arch, agent_version, labels, status, last_heartbeat_at, registered_at, updated_at FROM agents WHERE hostname = ? ORDER BY registered_at DESC LIMIT 1`, hostname)
	return scanAgent(row)
}

func (s *Store) UpdateAgentInfo(id, hostname, ipAddress, os, arch, version string) error {
	_, err := s.db.Exec(`UPDATE agents SET hostname=?, ip_address=?, os=?, arch=?, agent_version=?, status='online', last_heartbeat_at=datetime('now'), updated_at=datetime('now') WHERE id=?`,
		hostname, ipAddress, os, arch, version, id)
	return err
}

func (s *Store) ListAgents() ([]*models.Agent, error) {
	rows, err := s.db.Query(`SELECT id, name, hostname, ip_address, os, arch, agent_version, labels, status, last_heartbeat_at, registered_at, updated_at FROM agents ORDER BY registered_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []*models.Agent
	for rows.Next() {
		a, err := scanAgentRows(rows)
		if err != nil {
			return nil, err
		}
		agents = append(agents, a)
	}
	return agents, rows.Err()
}

func (s *Store) UpdateAgentHeartbeat(id string) error {
	_, err := s.db.Exec(`UPDATE agents SET last_heartbeat_at = datetime('now'), status = 'online', updated_at = datetime('now') WHERE id = ?`, id)
	return err
}

func (s *Store) AgentCount() (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM agents`).Scan(&count)
	return count, err
}

func scanAgent(row *sql.Row) (*models.Agent, error) {
	var a models.Agent
	var labelsJSON string
	var lastHB sql.NullString
	var regAt, updAt string

	err := row.Scan(&a.ID, &a.Name, &a.Hostname, &a.IPAddress, &a.OS, &a.Arch, &a.AgentVersion, &labelsJSON, &a.Status, &lastHB, &regAt, &updAt)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(labelsJSON), &a.Labels); err != nil {
		a.Labels = map[string]string{}
	}
	if lastHB.Valid {
		t := parseDBTime(lastHB.String)
		if !t.IsZero() {
			a.LastHeartbeatAt = &t
		}
	}
	a.RegisteredAt = parseDBTime(regAt)
	a.UpdatedAt = parseDBTime(updAt)
	return &a, nil
}

func scanAgentRows(rows *sql.Rows) (*models.Agent, error) {
	var a models.Agent
	var labelsJSON string
	var lastHB sql.NullString
	var regAt, updAt string

	err := rows.Scan(&a.ID, &a.Name, &a.Hostname, &a.IPAddress, &a.OS, &a.Arch, &a.AgentVersion, &labelsJSON, &a.Status, &lastHB, &regAt, &updAt)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(labelsJSON), &a.Labels); err != nil {
		a.Labels = map[string]string{}
	}
	if lastHB.Valid {
		t := parseDBTime(lastHB.String)
		if !t.IsZero() {
			a.LastHeartbeatAt = &t
		}
	}
	a.RegisteredAt = parseDBTime(regAt)
	a.UpdatedAt = parseDBTime(updAt)
	return &a, nil
}
