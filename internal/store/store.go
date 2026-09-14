// Package store is the only place in the server that writes queries.
//
// Everything above it — the HTTP handlers, the auth service — asks the store
// for what it needs and gets models back. That keeps the queries in one place
// to read and change, keeps GORM out of the handlers, and means the errors
// the rest of the code handles are this package's own rather than the
// driver's.
package store

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Store holds the database connection every query runs on.
type Store struct {
	db *gorm.DB
}

// New returns a store backed by the given connection.
func New(db *gorm.DB) *Store {
	return &Store{db: db}
}

// ErrNotFound is returned when a row that was asked for is not there. It
// stands in for gorm.ErrRecordNotFound so callers do not import GORM to
// answer a 404.
var ErrNotFound = errors.New("not found")

// ErrDuplicate is returned when a write would repeat a value a unique index
// forbids — the same email twice, say. It is the writer's to fix, which is
// why it is told apart from anything else that can go wrong.
var ErrDuplicate = errors.New("already exists")

// translate turns what the driver returns into this package's errors.
// Anything it does not recognise is passed through unchanged.
func translate(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrNotFound
	case strings.Contains(strings.ToLower(err.Error()), "duplicate"):
		// Wrapped rather than replaced: errors.Is still finds the sentinel,
		// and what the database said is still there for DuplicateField to
		// read the column out of.
		return fmt.Errorf("%w: %s", ErrDuplicate, err)
	default:
		return err
	}
}

// likeEscaper escapes what LIKE treats specially, so a search for "50%" or
// "first_name" finds those characters rather than matching anything.
var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// contains is the LIKE pattern for a case-insensitive search: the value, lower
// cased and escaped, anywhere in the column. Compare it against LOWER(column).
func contains(search string) string {
	return "%" + likeEscaper.Replace(strings.ToLower(search)) + "%"
}
