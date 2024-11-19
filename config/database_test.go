package config

import (
	"testing"
)

func TestOpenDB(t *testing.T) {
	dsn := "web:Secure@123@tcp(localhost:3306)/snippetbox?parseTime=true"

	// Test opening the database
	db, err := OpenDB(dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test a simple query to ensure connection is valid
	var result string
	err = db.QueryRow("SELECT 'Connection successful!'").Scan(&result)
	if err != nil {
		t.Fatalf("Test query failed: %v", err)
	}

	if result != "Connection successful!" {
		t.Fatalf("Unexpected result from test query: %s", result)
	}

	t.Log("Database connection test successful!")
}
