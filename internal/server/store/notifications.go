package store

import "fmt"

type NotificationChannel struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Config    string `json:"config"`
	Enabled   bool   `json:"enabled"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (s *Store) ListNotificationChannels() ([]*NotificationChannel, error) {
	rows, err := s.db.Query(`SELECT id, type, name, config, enabled, created_at, updated_at FROM notification_channels ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*NotificationChannel
	for rows.Next() {
		c := &NotificationChannel{}
		var enabled int
		err := rows.Scan(&c.ID, &c.Type, &c.Name, &c.Config, &enabled, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		c.Enabled = enabled == 1
		channels = append(channels, c)
	}
	return channels, rows.Err()
}

func (s *Store) CreateNotificationChannel(c *NotificationChannel) error {
	_, err := s.db.Exec(`INSERT INTO notification_channels (id, type, name, config, enabled) VALUES (?, ?, ?, ?, ?)`,
		c.ID, c.Type, c.Name, c.Config, boolToInt(c.Enabled))
	return err
}

func (s *Store) UpdateNotificationChannel(c *NotificationChannel) error {
	result, err := s.db.Exec(`UPDATE notification_channels SET type=?, name=?, config=?, enabled=?, updated_at=datetime('now') WHERE id=?`,
		c.Type, c.Name, c.Config, boolToInt(c.Enabled), c.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("notification channel not found")
	}
	return nil
}

func (s *Store) DeleteNotificationChannel(id string) error {
	result, err := s.db.Exec(`DELETE FROM notification_channels WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("notification channel not found")
	}
	return nil
}
