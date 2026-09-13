package model

import "testing"

// field builds a field of the given type for a test to use.
func field(t FieldType, required bool) UserField {
	return UserField{Name: "probe", Label: "Probe", Type: t, Required: required}
}

// TestNormalise covers what a submitted value becomes, which is the one place
// a field's type actually means something.
func TestNormalise(t *testing.T) {
	tests := []struct {
		name  string
		field UserField
		input any
		want  any
	}{
		{name: "text is trimmed", field: field(FieldText, false), input: "  Mira  ", want: "Mira"},
		{name: "number from a string", field: field(FieldNumber, false), input: "1250", want: 1250.0},
		{name: "number from a number", field: field(FieldNumber, false), input: 42.5, want: 42.5},
		{name: "bool from a bool", field: field(FieldBool, false), input: true, want: true},
		{name: "bool from a string", field: field(FieldBool, false), input: "true", want: true},
		{name: "email is trimmed", field: field(FieldEmail, false), input: " a@b.com ", want: "a@b.com"},
		{name: "a day becomes a timestamp", field: field(FieldDate, false), input: "2026-09-13", want: "2026-09-13T00:00:00Z"},

		// An optional field with nothing in it is stored as nothing, rather
		// than as an empty string that would have to be handled everywhere.
		{name: "empty optional text", field: field(FieldText, false), input: "", want: nil},
		{name: "missing optional text", field: field(FieldText, false), input: nil, want: nil},
		{name: "empty optional number", field: field(FieldNumber, false), input: "", want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.field.Normalise(tt.input)
			if err != nil {
				t.Fatalf("Normalise(%#v) = error %v, want %#v", tt.input, err, tt.want)
			}
			if got != tt.want {
				t.Errorf("Normalise(%#v) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormaliseRejects(t *testing.T) {
	tests := []struct {
		name  string
		field UserField
		input any
	}{
		{name: "text in a number", field: field(FieldNumber, false), input: "many"},
		{name: "a word in a bool", field: field(FieldBool, false), input: "maybe"},
		{name: "an address with no at sign", field: field(FieldEmail, false), input: "nobody"},
		{name: "a date that is not one", field: field(FieldDate, false), input: "last tuesday"},
		{name: "a missing required value", field: field(FieldText, true), input: nil},
		{name: "an empty required value", field: field(FieldText, true), input: "   "},
		{name: "an empty required number", field: field(FieldNumber, true), input: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.field.Normalise(tt.input); err == nil {
				t.Errorf("Normalise(%#v) accepted the value, want an error", tt.input)
			}
		})
	}
}

// The error names the field, because it is shown to whoever is editing.
func TestNormaliseErrorNamesTheField(t *testing.T) {
	_, err := UserField{Name: "phone_verified", Type: FieldBool}.Normalise("maybe")
	if err == nil {
		t.Fatal("want an error")
	}

	if got := err.Error(); got != "phone_verified: must be true or false" {
		t.Errorf("error = %q", got)
	}
}

func TestFieldTypeValid(t *testing.T) {
	for _, known := range FieldTypes {
		if !known.Valid() {
			t.Errorf("%s.Valid() = false, want true", known)
		}
	}

	for _, unknown := range []FieldType{"", "blob", "Text"} {
		if unknown.Valid() {
			t.Errorf("%q.Valid() = true, want false", unknown)
		}
	}
}

// bound is a pointer to a number, which is how a field carries a bound it
// might not have.
func bound(value float64) *float64 {
	return &value
}

// TestNormaliseRules covers the rules a field can put on its values: the
// checks that turn a field definition into something an admin can rely on.
func TestNormaliseRules(t *testing.T) {
	tests := []struct {
		name    string
		field   UserField
		input   any
		wantErr bool
	}{
		{
			name:  "a number inside its bounds",
			field: UserField{Name: "age", Type: FieldNumber, Min: bound(18), Max: bound(120)},
			input: "30",
		},
		{
			name:    "a number below its smallest",
			field:   UserField{Name: "age", Type: FieldNumber, Min: bound(18)},
			input:   "17",
			wantErr: true,
		},
		{
			name:    "a number above its largest",
			field:   UserField{Name: "age", Type: FieldNumber, Max: bound(120)},
			input:   "121",
			wantErr: true,
		},
		{
			name:  "text of an allowed length",
			field: UserField{Name: "nickname", Type: FieldText, Min: bound(2), Max: bound(8)},
			input: "Mira",
		},
		{
			name:    "text that is too short",
			field:   UserField{Name: "nickname", Type: FieldText, Min: bound(2)},
			input:   "M",
			wantErr: true,
		},
		{
			name:    "text that is too long",
			field:   UserField{Name: "nickname", Type: FieldText, Max: bound(3)},
			input:   "Mirabel",
			wantErr: true,
		},
		{
			name:  "text with the prefix it needs",
			field: UserField{Name: "phone", Type: FieldText, StartsWith: "+"},
			input: "+996700000001",
		},
		{
			name:    "text without the prefix it needs",
			field:   UserField{Name: "phone", Type: FieldText, StartsWith: "+"},
			input:   "996700000001",
			wantErr: true,
		},
		{
			name:  "an empty optional field skips its rules",
			field: UserField{Name: "phone", Type: FieldText, StartsWith: "+", Min: bound(5)},
			input: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.field.Normalise(tt.input)

			if tt.wantErr && err == nil {
				t.Fatalf("Normalise(%#v) = no error, want one", tt.input)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Normalise(%#v) = error %v, want none", tt.input, err)
			}
		})
	}
}

// TestFieldValidate covers the definition itself: a rule a field cannot keep
// is a mistake in the panel, and is refused there rather than confusing
// whoever edits a record later.
func TestFieldValidate(t *testing.T) {
	tests := []struct {
		name    string
		field   UserField
		wantErr bool
	}{
		{name: "text with a length", field: UserField{Name: "a", Type: FieldText, Max: bound(10)}},
		{name: "a number with bounds", field: UserField{Name: "a", Type: FieldNumber, Min: bound(0), Max: bound(9)}},
		{name: "a plain bool", field: UserField{Name: "a", Type: FieldBool}},
		{
			name:    "a bool cannot be bounded",
			field:   UserField{Name: "a", Type: FieldBool, Max: bound(10)},
			wantErr: true,
		},
		{
			name:    "a date cannot have a prefix",
			field:   UserField{Name: "a", Type: FieldDate, StartsWith: "2026"},
			wantErr: true,
		},
		{
			name:    "a largest below its smallest",
			field:   UserField{Name: "a", Type: FieldNumber, Min: bound(10), Max: bound(1)},
			wantErr: true,
		},
		{
			name:    "a negative length",
			field:   UserField{Name: "a", Type: FieldText, Min: bound(-1)},
			wantErr: true,
		},
		{
			name:    "an unknown type",
			field:   UserField{Name: "a", Type: "blob"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.field.Validate()

			if tt.wantErr && err == nil {
				t.Fatal("Validate() = no error, want one")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Validate() = %v, want none", err)
			}
		})
	}
}
