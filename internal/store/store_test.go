package store

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

// TestTranslate covers the one piece of the store that is worth testing
// without a database: turning what the driver says into the errors the rest
// of the project handles.
func TestTranslate(t *testing.T) {
	tests := []struct {
		name string
		in   error
		want error
	}{
		{name: "nothing went wrong", in: nil, want: nil},
		{name: "no such row", in: gorm.ErrRecordNotFound, want: ErrNotFound},
		{
			name: "a wrapped missing row still counts",
			in:   errors.Join(errors.New("loading user"), gorm.ErrRecordNotFound),
			want: ErrNotFound,
		},
		{
			name: "a unique index says duplicate",
			in:   errors.New(`ERROR: duplicate key value violates unique constraint "idx_users_email"`),
			want: ErrDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := translate(tt.in)

			if !errors.Is(got, tt.want) {
				t.Fatalf("translate(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

// Anything the store does not recognise has to reach the caller unchanged: a
// connection that dropped is not a missing row.
func TestTranslateKeepsOtherErrors(t *testing.T) {
	original := errors.New("connection refused")

	if got := translate(original); !errors.Is(got, original) {
		t.Fatalf("translate(%v) = %v, want it unchanged", original, got)
	}
}
