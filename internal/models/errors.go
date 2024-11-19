package models

import "errors"

// ErrNoRecord is returned when a database query does not return any rows.
var ErrNoRecord = errors.New("models: no matching record found")
