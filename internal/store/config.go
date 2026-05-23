// internal/store/config.go
package store

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Sylvester-Kapoko/risitPap/domain"
	"github.com/google/uuid"
)

func (s *Store) SaveConfig(cfg *domain.StoreConfig) error {
	// Generate RegisterID exactly once
	if cfg.RegisterID == "" {
		cfg.RegisterID = "MK-" + uuid.New().String()
	}
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	_, err = s.db.Exec("INSERT OR REPLACE INTO config (id, data) VALUES (1, ?)", string(data))
	return err
}

func (s *Store) LoadConfig() (*domain.StoreConfig, error) {
	var data string
	err := s.db.QueryRow("SELECT data FROM config WHERE id = 1").Scan(&data)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(strings.NewReader(data))
	dec.DisallowUnknownFields()
	var cfg domain.StoreConfig
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
