package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// OpenDB opens a new database connection using the provided DSN (Data Source Name).
func OpenDB(dsn string) (*sql.DB, error) {
	// Open a connection to the database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	// Verify the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to verify database connection: %w", err)
	}

	return db, nil
}
