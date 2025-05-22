package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewClient(conn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("could not connect to Postgres: %w", err)
	}

	return db, nil
}
