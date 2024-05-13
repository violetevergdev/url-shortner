package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := stmt.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(alias, urlToSave string) (int64, error) {
	stmt, err := s.db.Prepare(`INSERT INTO url(alias, url) VALUES(?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("storage.sqlite.SaveURL: %w", err)
	}

	res, err := stmt.Exec(alias, urlToSave)
	if err != nil {
		return 0, fmt.Errorf("storage.sqlite.SaveURL: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("storage.sqlite.SaveURL: failed to get last inserted id: %w", err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE alias =?`)
	if err != nil {
		return "", fmt.Errorf("storage.sqlite.GetURL: %w", err)
	}

	var urlResult string
	err = stmt.QueryRow(alias).Scan(&urlResult)
	if err != nil {
		return "", fmt.Errorf("\nstorage.sqlite.GetURL: \nfailed to get url: %w", err)
	}

	return urlResult, nil
}

func (s *Storage) DeleteURL(alias string) error {
	stmt, err := s.db.Prepare(`DELETE FROM url WHERE alias = ?`)
	if err != nil {
		return fmt.Errorf("storage.sqlite.DeleteURL: %w", err)
	}

	if _, err := stmt.Exec(alias); err != nil {
		return fmt.Errorf("\nstorage.sqlite.DeleteURL: \nfailed to delete url: %w", err)
	}

	return nil
}
