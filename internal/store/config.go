package store

import (
	"encoding/json"
	"github.com/Sylvester-Kapoko/Receipts/domain"
)


func (s *Store) SaveConfig (cfg *domain.StoreConfig) error {
	data, _ := json.Marshal(cfg)
	_ , err := s.db.Exec("INSERT OR REPLACE INTO config (id, data) VALUES (1, ?) ", string(data))
	return err

}

func (s *Store) LoadConfig() (*domain.StoreConfig, error) {
	var data string
	err := s.db.QueryRow("SELECT data FROM config WHERE id = 1").Scan(&data)
	if err != nil {
		return nil, err
	}

	var cfg domain.StoreConfig
	if err := json.Unmarshal([]byte(data), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
