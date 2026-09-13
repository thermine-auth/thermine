package users

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"xermess/internal/api/http/respond"
	"xermess/internal/api/http/validate"
	"xermess/internal/model"
	"xermess/internal/store"
)

// validate checks the built-in fields, which are the record's own columns.
// What the additional ones have to keep is checked by normalise below,
// against the definitions in user_fields rather than against a tag.
func (r *userRequest) validate() error {
	r.clean()

	return validate.Struct(r)
}

// normalise checks the submitted values against the field definitions and
// returns what should be stored: the rules a field carries — its type, its
// bounds, its prefix, whether it has to be unique — are all applied here.
//
// `self` is the user being written, so a unique field does not find the
// record's own value and call it a clash. It is uuid.Nil when creating one.
func normalise(
	ctx context.Context,
	st *store.Store,
	fields []model.UserField,
	values map[string]any,
	self uuid.UUID,
) (map[string]any, error) {
	data := make(map[string]any, len(fields))

	for _, field := range fields {
		value, err := field.Normalise(values[field.Name])
		if err != nil {
			return nil, respond.Fault{Status: http.StatusBadRequest, Message: err.Error()}
		}

		// A field with no value is left out rather than stored as null, so
		// removing a field leaves nothing behind.
		if value == nil {
			continue
		}

		if field.Unique {
			taken, err := st.FieldValueTaken(ctx, field.Name, value, self)
			if err != nil {
				return nil, err
			}

			if taken {
				return nil, respond.Fault{
					Status:  http.StatusConflict,
					Message: field.Name + ": another user already has that value",
				}
			}
		}

		data[field.Name] = value
	}

	return data, nil
}
