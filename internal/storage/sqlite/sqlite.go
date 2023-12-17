package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" //init sqlite driver
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
