package validate

import (
	"errors"
	"net/http"
	"testing"

	"xermess/internal/api/http/respond"
)

// fault is the fault an error carries, and whether it carried one at all.
func fault(err error) (respond.Fault, bool) {
	var f respond.Fault
	ok := errors.As(err, &f)

	return f, ok
}

type probe struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"max=4"`
	Kind  string `json:"kind" validate:"omitempty,oneof=one two"`
}

// TestStruct covers what a handler sees when a rule is broken: a 400, and a
// sentence that names the field the way the request named it.
func TestStruct(t *testing.T) {
	tests := []struct {
		name  string
		value probe
		want  string // the message, or "" when the value is fine
	}{
		{
			name:  "everything in order",
			value: probe{Email: "a@b.com", Name: "Mira", Kind: "one"},
		},
		{
			name:  "a missing field",
			value: probe{Name: "Mira"},
			want:  "email is required",
		},
		{
			name:  "something that is not an address",
			value: probe{Email: "not-an-address"},
			want:  "email must be an email address",
		},
		{
			name:  "too long",
			value: probe{Email: "a@b.com", Name: "Mirabel"},
			want:  "name must be at most 4 characters",
		},
		{
			name:  "not one of the values allowed",
			value: probe{Email: "a@b.com", Kind: "three"},
			want:  "kind must be one of: one, two",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Struct(tt.value)

			if tt.want == "" {
				if err != nil {
					t.Fatalf("Struct() = %v, want nothing", err)
				}
				return
			}

			broken, ok := fault(err)
			if !ok {
				t.Fatalf("Struct() = %v, want a respond.Fault", err)
			}

			if broken.Status != http.StatusBadRequest {
				t.Errorf("status = %d, want 400", broken.Status)
			}
			if broken.Message != tt.want {
				t.Errorf("message = %q, want %q", broken.Message, tt.want)
			}
		})
	}
}

// A field with no json tag is named by the field itself, rather than by
// nothing at all.
func TestStructNamesAFieldWithoutATag(t *testing.T) {
	type untagged struct {
		Reference string `validate:"required"`
	}

	broken, ok := fault(Struct(untagged{}))
	if !ok {
		t.Fatal("want a respond.Fault")
	}

	if broken.Message != "Reference is required" {
		t.Errorf("message = %q, want %q", broken.Message, "Reference is required")
	}
}

// A rule of our own is asked about the value, and its message is used when it
// says no.
func TestRegister(t *testing.T) {
	Register("shouty", "must be in capitals", func(value string) bool {
		return value == "LOUD"
	})

	type probe struct {
		Word string `json:"word" validate:"shouty"`
	}

	if err := Struct(probe{Word: "LOUD"}); err != nil {
		t.Fatalf("Struct() = %v, want it accepted", err)
	}

	broken, ok := fault(Struct(probe{Word: "quiet"}))
	if !ok {
		t.Fatal("want a respond.Fault")
	}

	if broken.Message != "word must be in capitals" {
		t.Errorf("message = %q, want %q", broken.Message, "word must be in capitals")
	}
}

// Something that is not a struct is our own mistake, and must not come back
// as a 400 blaming whoever called the endpoint.
func TestStructRefusesToBlameTheCaller(t *testing.T) {
	err := Struct("not a struct")

	if err == nil {
		t.Fatal("Struct() = nothing, want an error")
	}

	if _, ok := fault(err); ok {
		t.Errorf("Struct() = %v, want it not to be a Fault", err)
	}
}
