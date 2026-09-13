package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FieldType is the kind of value a user field holds. Adding one here means
// teaching Normalise below how to read it, and the panel how to render and
// edit it.
type FieldType string

const (
	FieldText   FieldType = "text"
	FieldNumber FieldType = "number"
	FieldBool   FieldType = "bool"
	FieldEmail  FieldType = "email"
	FieldDate   FieldType = "date"
)

// FieldTypes lists every type, for validation and for the panel's type picker.
var FieldTypes = []FieldType{FieldText, FieldNumber, FieldBool, FieldEmail, FieldDate}

// Valid reports whether t is a known type.
func (t FieldType) Valid() bool {
	for _, known := range FieldTypes {
		if t == known {
			return true
		}
	}
	return false
}

// UserField is one column of the user record, defined at runtime rather than
// in the schema: the values themselves live in users.data, so adding a field
// is a row here rather than a migration.
type UserField struct {
	Base

	Name     string    `gorm:"uniqueIndex;size:64;not null" json:"name"`
	Label    string    `gorm:"size:100;not null" json:"label"`
	Type     FieldType `gorm:"type:varchar(16);not null" json:"type"`
	Required bool      `gorm:"not null;default:false" json:"required"`

	// Unique means no two users may hold the same value. It is checked when a
	// user is written, since the values live in a JSON column rather than in
	// a column the database could index. The column is is_unique because
	// UNIQUE is a word of SQL's own.
	Unique bool `gorm:"column:is_unique;not null;default:false" json:"unique"`

	// Min and Max bound a number's value, or the length of text. Nil means
	// the field is not bounded that way.
	Min *float64 `json:"min"`
	Max *float64 `json:"max"`

	// StartsWith is a prefix text has to begin with, such as "+" for a phone
	// number. Empty means anything is accepted.
	StartsWith string `gorm:"size:64;not null;default:''" json:"starts_with"`

	// Position orders the columns in the table and the inputs in the form.
	Position int `gorm:"not null;default:0" json:"position"`
}

// Bounded reports whether this type of field can carry a Min and a Max.
// Numbers are bounded by value, text by length; the rest are not bounded.
func (t FieldType) Bounded() bool {
	return t == FieldNumber || t == FieldText || t == FieldEmail
}

// Prefixed reports whether this type of field can carry a StartsWith.
func (t FieldType) Prefixed() bool {
	return t == FieldText || t == FieldEmail
}

// Validate checks the definition itself, as opposed to a value stored under
// it: a bound the field cannot have, or a maximum below its minimum, is a
// mistake in the panel rather than in someone's record.
func (f UserField) Validate() error {
	if !f.Type.Valid() {
		return ErrFieldValue{f.Name, "is not a field type"}
	}

	if !f.Type.Bounded() && (f.Min != nil || f.Max != nil) {
		return ErrFieldValue{f.Name, "cannot have a smallest or largest value"}
	}

	if !f.Type.Prefixed() && f.StartsWith != "" {
		return ErrFieldValue{f.Name, "cannot have a prefix"}
	}

	if f.Min != nil && f.Max != nil && *f.Min > *f.Max {
		return ErrFieldValue{f.Name, "has a largest value below its smallest"}
	}

	if f.Type != FieldNumber && (negative(f.Min) || negative(f.Max)) {
		return ErrFieldValue{f.Name, "is measured in characters, which cannot be negative"}
	}

	return nil
}

func negative(value *float64) bool {
	return value != nil && *value < 0
}

// TableName pins the table name.
func (UserField) TableName() string {
	return "user_fields"
}

// ErrFieldValue is returned when a value does not suit its field. The message
// is shown to whoever is editing, so it says what was expected.
type ErrFieldValue struct {
	Field  string
	Reason string
}

func (e ErrFieldValue) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Reason)
}

// Normalise checks a value against the field and returns it in the form that
// should be stored: numbers as float64, booleans as bool, dates as RFC 3339,
// and everything else trimmed. An empty value comes back as nil, which is how
// "not set" is stored.
func (f UserField) Normalise(value any) (any, error) {
	if value == nil {
		return nil, f.requiredError()
	}

	switch f.Type {
	case FieldBool:
		switch typed := value.(type) {
		case bool:
			return typed, nil
		case string:
			parsed, err := strconv.ParseBool(strings.TrimSpace(typed))
			if err != nil {
				return nil, ErrFieldValue{f.Name, "must be true or false"}
			}
			return parsed, nil
		default:
			return nil, ErrFieldValue{f.Name, "must be true or false"}
		}

	case FieldNumber:
		switch typed := value.(type) {
		case float64:
			return typed, f.checkNumber(typed)
		case int:
			return float64(typed), f.checkNumber(float64(typed))
		case string:
			trimmed := strings.TrimSpace(typed)
			if trimmed == "" {
				return nil, f.requiredError()
			}
			parsed, err := strconv.ParseFloat(trimmed, 64)
			if err != nil {
				return nil, ErrFieldValue{f.Name, "must be a number"}
			}
			return parsed, f.checkNumber(parsed)
		default:
			return nil, ErrFieldValue{f.Name, "must be a number"}
		}

	case FieldDate:
		text, ok := value.(string)
		if !ok {
			return nil, ErrFieldValue{f.Name, "must be a date"}
		}
		trimmed := strings.TrimSpace(text)
		if trimmed == "" {
			return nil, f.requiredError()
		}
		// The panel sends a plain date; both that and a full timestamp are
		// accepted, and both are stored as a timestamp.
		for _, layout := range []string{time.RFC3339, "2006-01-02"} {
			if parsed, err := time.Parse(layout, trimmed); err == nil {
				return parsed.Format(time.RFC3339), nil
			}
		}
		return nil, ErrFieldValue{f.Name, "must be a date like 2026-09-13"}

	case FieldEmail:
		text, ok := value.(string)
		if !ok {
			return nil, ErrFieldValue{f.Name, "must be an email address"}
		}
		trimmed := strings.TrimSpace(text)
		if trimmed == "" {
			return nil, f.requiredError()
		}
		if !strings.Contains(trimmed, "@") || strings.HasPrefix(trimmed, "@") || strings.HasSuffix(trimmed, "@") {
			return nil, ErrFieldValue{f.Name, "must be an email address"}
		}
		return trimmed, f.checkText(trimmed)

	default: // FieldText
		text, ok := value.(string)
		if !ok {
			return nil, ErrFieldValue{f.Name, "must be text"}
		}
		trimmed := strings.TrimSpace(text)
		if trimmed == "" {
			return nil, f.requiredError()
		}
		return trimmed, f.checkText(trimmed)
	}
}

// checkNumber holds a number to the bounds the field was given.
func (f UserField) checkNumber(value float64) error {
	if f.Min != nil && value < *f.Min {
		return ErrFieldValue{f.Name, fmt.Sprintf("must be %s or more", number(*f.Min))}
	}

	if f.Max != nil && value > *f.Max {
		return ErrFieldValue{f.Name, fmt.Sprintf("must be %s or less", number(*f.Max))}
	}

	return nil
}

// checkText holds text to its length bounds and its prefix. Length is counted
// in characters rather than bytes, so an accented name counts as one each.
func (f UserField) checkText(value string) error {
	length := float64(len([]rune(value)))

	if f.Min != nil && length < *f.Min {
		return ErrFieldValue{f.Name, fmt.Sprintf("must be at least %s characters", number(*f.Min))}
	}

	if f.Max != nil && length > *f.Max {
		return ErrFieldValue{f.Name, fmt.Sprintf("must be at most %s characters", number(*f.Max))}
	}

	if f.StartsWith != "" && !strings.HasPrefix(value, f.StartsWith) {
		return ErrFieldValue{f.Name, fmt.Sprintf("must start with %q", f.StartsWith)}
	}

	return nil
}

// number prints a bound the way it was written: 3 rather than 3.0.
func number(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

// requiredError is nil for an optional field, so an empty value is simply not
// stored.
func (f UserField) requiredError() error {
	if f.Required {
		return ErrFieldValue{f.Name, "is required"}
	}
	return nil
}
