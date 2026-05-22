// internal/store/audit.go
package store

import "time"

// LogAction writes an audit entry.
func (s *Store) LogAction(username, action, detail string) error {
    _, err := s.db.Exec(
        "INSERT INTO audit_log (username, action, detail, timestamp) VALUES (?, ?, ?, ?)",
        username, action, detail, time.Now().Format(time.RFC3339),
    )
    return err
}

// GetAuditLogs returns the most recent audit entries.
func (s *Store) GetAuditLogs(limit int) ([]map[string]string, error) {
    rows, err := s.db.Query("SELECT username, action, timestamp FROM audit_log ORDER BY id DESC LIMIT ?", limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var logs []map[string]string
    for rows.Next() {
        var user, action, ts string
        if err := rows.Scan(&user, &action, &ts); err != nil {
            return nil, err
        }
        logs = append(logs, map[string]string{"user": user, "action": action, "time": ts})
    }
    return logs, rows.Err()
}

// LogSystem is a convenience wrapper for system events.
func (s *Store) LogSystem(action, detail string) error {
    return s.LogAction("system", action, detail)
}