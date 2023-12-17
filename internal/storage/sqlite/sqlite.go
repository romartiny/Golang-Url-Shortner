package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"golang-url-shortner/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	//for errors debugging
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//simple migration
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
            alias TEXT NOT NULL UNIQUE,
            url TEXT NOT NULL UNIQUE
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)

	//catch errors
	if err != nil {
		return nil, fmt.Errorf("#{op}: #{err}")
	}

	//catch execution errors
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SaveURL to saving url
func (s *Storage) SaveURL(urlToSave string, alisa string) (int64, error) {
	//for errors debugging
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url (url, alias) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alisa)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUrlExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	//get id from storage
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	//return id and empty error message
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	//for errors debugging
	const op = "storage.sqlite.GetURL"

	//prepare query for get url via alias
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias =?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)

	//not alias in database & est error
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	//return url
	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	//for errors debugging
	const op = "storage.sqlite.DeleteURL"

	//find and delete url via alias
	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias =?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	//error
	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
