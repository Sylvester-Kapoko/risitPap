package store

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "time"

    "github.com/Sylvester-Kapoko/risitPap/domain"
    _ "modernc.org/sqlite"
)

type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        return nil, err
    }

    if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS receipts (
        id TEXT PRIMARY KEY,
        data TEXT NOT NULL,
        created_at TEXT NOT NULL
    )`); err != nil {
        return nil, fmt.Errorf("create receipts table: %w", err)
    }

    if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS config (
        id INTEGER PRIMARY KEY DEFAULT 1,
        data TEXT NOT NULL
    )`); err != nil {
        return nil, fmt.Errorf("create config table: %w", err)
    }

    if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS items (
        name TEXT PRIMARY KEY,
        last_price TEXT NOT NULL
    )`); err != nil {
        return nil, fmt.Errorf("create items table: %w", err)
    }

    return &Store{db: db}, nil
}

func (s *Store) Save(r *domain.Receipt) error {
    r.CreatedAt = time.Now()
    data, err := json.Marshal(r)
    if err != nil {
        return err
    }
    _, err = s.db.Exec("INSERT INTO receipts (id, data, created_at) VALUES (?, ?, ?)",
        r.TransactionID, string(data), r.CreatedAt.Format(time.RFC3339))
    if err != nil {
        return err
    }

    for _, item := range r.Items {
        _, err = s.db.Exec("INSERT OR REPLACE INTO items (name, last_price) VALUES (?, ?)",
            item.Name, item.UnitPrice.String())
        if err != nil {
            return err
        }
    }
    return nil
}

func (s *Store) List(since time.Time) ([]domain.Receipt, error) {
    rows, err := s.db.Query("SELECT data FROM receipts ORDER BY created_at DESC LIMIT 50")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var receipts []domain.Receipt
    for rows.Next() {
        var data string
        if err := rows.Scan(&data); err != nil {
            continue
        }
        var r domain.Receipt
        if err := json.Unmarshal([]byte(data), &r); err != nil {
            continue
        }
        receipts = append(receipts, r)
    }
    return receipts, rows.Err()
}

func (s *Store) Get(id string) (*domain.Receipt, error) {
    var data string
    err := s.db.QueryRow("SELECT data FROM receipts WHERE id = ?", id).Scan(&data)
    if err != nil {
        return nil, err
    }
    var r domain.Receipt
    if err := json.Unmarshal([]byte(data), &r); err != nil {
        return nil, err
    }
    return &r, nil
}

func (s *Store) SuggestItems(prefix string) ([]map[string]string, error) {
    rows, err := s.db.Query("SELECT name, last_price FROM items WHERE name LIKE ? LIMIT 5", prefix+"%")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var results []map[string]string
    for rows.Next() {
        var name, price string
        if err := rows.Scan(&name, &price); err != nil {
            return nil, err
        }
        results = append(results, map[string]string{"name": name, "price": price})
    }
    return results, rows.Err()
}