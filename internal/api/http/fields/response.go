package fields

import "xermess/internal/model"

// listResponse is the fields, and the types a new one may have. Both come
// together because the panel's form needs the second to offer the first.
type listResponse struct {
	Fields []model.UserField `json:"fields"`
	Types  []model.FieldType `json:"types"`
}

func newListResponse(fields []model.UserField) listResponse {
	return listResponse{Fields: fields, Types: model.FieldTypes}
}
