package postgres

import "errors"

// ErrConflict is returned on unique constraint violations.
var ErrConflict = errors.New("repo conflict")

// ErrNotFound is returned when a record is not found.
var ErrNotFound = errors.New("repo not found")
