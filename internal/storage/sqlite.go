package storage

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

type UsersStore struct {
	db *sql.DB
}

func New(db *sql.DB) *UsersStore {
	return &UsersStore{db: db}
}

func (s *UsersStore) Initialize() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			chatId INTEGER PRIMARY KEY,
			phone TEXT NOT NULL
		)
	`)
	return err
}

func (s *UsersStore) Add(chatID int64, phone string) error {
	_, err := s.db.Exec(
		`INSERT OR REPLACE INTO users (chatId, phone) VALUES (?, ?)`,
		chatID, phone,
	)
	return err
}

func (s *UsersStore) Get(chatID int64) (string, error) {
	var phone string
	err := s.db.QueryRow(
		`SELECT phone FROM users WHERE chatId = ?`,
		chatID,
	).Scan(&phone)
	return phone, err
}
