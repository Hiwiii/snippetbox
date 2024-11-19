package models

import (
	"database/sql"
	"time"
	"errors"
)

// Define a Snippet type to hold the data for an individual snippet.
// The fields correspond to the fields in the MySQL snippets table.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Insert inserts a new snippet into the database.
// Insert inserts a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// Get returns a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	// Write the SQL statement to execute
	stmt := `SELECT id, title, content, created, expires FROM snippets 
			 WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Use the QueryRow() method to execute the statement and return a sql.Row object
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct
	s := &Snippet{}

	// Use row.Scan() to copy the values from the sql.Row into the Snippet struct fields
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns no rows, row.Scan() will return a sql.ErrNoRows error
		// Handle that specific error and return a custom ErrNoRecord error
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	// If everything is OK, return the Snippet object
	return s, nil
}


// Latest returns the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	// Write the SQL statement to execute. It selects the 10 most recent snippets
	// where the expiry date is still in the future, ordered by descending ID.
	stmt := `SELECT id, title, content, created, expires 
	         FROM snippets
	         WHERE expires > UTC_TIMESTAMP()
	         ORDER BY id DESC 
	         LIMIT 10`

	// Use the Query() method to execute the SQL statement. This returns a
	// *sql.Rows result set containing the result of the query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// Ensure that the result set is properly closed before the method returns.
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs.
	snippets := []*Snippet{}

	// Use rows.Next() to iterate through the rows in the result set.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &Snippet{}

		// Use rows.Scan() to copy the values from each field in the row to
		// the corresponding field in the Snippet struct.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		// Append the Snippet struct to the slice.
		snippets = append(snippets, s)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any
	// error that was encountered during iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK, return the slice of snippets.
	return snippets, nil
}

