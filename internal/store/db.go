package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sylvester-Kapoko/risitPap/domain"
	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
)

var eatLocation = time.FixedZone("EAT", 3*60*60)

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

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password_hash TEXT NOT NULL
	)`); err != nil {
		return nil, fmt.Errorf("create users table: %w", err)
	}

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS audit_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		action TEXT NOT NULL,
		detail TEXT DEFAULT '',
		timestamp TEXT NOT NULL
	)`); err != nil {
		return nil, fmt.Errorf("create audit_log table: %w", err)
	}

	return &Store{db: db}, nil
}

// --- Receipt methods ---

func (s *Store) Save(r *domain.Receipt) error {
	// CreatedAt is set by the handler in EAT — do not overwrite here
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		"INSERT INTO receipts (id, data, created_at) VALUES (?, ?, ?)",
		r.TransactionID, string(data), r.CreatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	for _, item := range r.Items {
		_, err = s.db.Exec(
			"INSERT OR REPLACE INTO items (name, last_price) VALUES (?, ?)",
			item.Name, item.UnitPrice.String(),
		)
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
		r.CreatedAt = r.CreatedAt.In(eatLocation)
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
	r.CreatedAt = r.CreatedAt.In(eatLocation)
	return &r, nil
}

func (s *Store) SuggestItems(prefix string) ([]map[string]string, error) {
	rows, err := s.db.Query(
		"SELECT name, last_price FROM items WHERE name LIKE ? LIMIT 5",
		prefix+"%",
	)
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

// --- User / Auth methods ---

func (s *Store) HasUsers() (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count > 0, err
}

func (s *Store) CreateUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		"INSERT INTO users (username, password_hash) VALUES (?, ?)",
		username, string(hash),
	)
	return err
}

func (s *Store) ValidateUser(username, password string) (*domain.User, error) {
	var u domain.User
	var hash string
	err := s.db.QueryRow(
		"SELECT username, password_hash FROM users WHERE username = ?", username,
	).Scan(&u.Username, &hash)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return nil, sql.ErrNoRows
	}
	u.PasswordHash = ""
	return &u, nil
}

// --- eTIMS sync methods ---

func (s *Store) GetUnsyncedReceipts() ([]domain.Receipt, error) {
    rows, err := s.db.Query("SELECT data FROM receipts WHERE json_extract(data, '$.sync_status') = 'pending'")
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

func (s *Store) UpdateReceiptSync(r domain.Receipt) error {
    data, err := json.Marshal(r)
    if err != nil {
        return err
    }
    _, err = s.db.Exec("UPDATE receipts SET data = ? WHERE id = ?", string(data), r.TransactionID)
    return err
}