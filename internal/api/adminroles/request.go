package adminroles

import "strings"

// targetType is what these rows are called in the activity log.
const targetType = "admin_role"

// roleRequest is the body of the create and update endpoints. Permissions are
// names from the catalog, and replace what the role granted.
type roleRequest struct {
	Name        string   `json:"name" validate:"required,adminrole"`
	Description string   `json:"description" validate:"max=255"`
	Permissions []string `json:"permissions"`
}

// clean tidies what can be tidied, so the rules see the values that would
// actually be stored.
func (r *roleRequest) clean() {
	r.Name = strings.ToLower(strings.TrimSpace(r.Name))
	r.Description = strings.TrimSpace(r.Description)

	for i, permission := range r.Permissions {
		r.Permissions[i] = strings.TrimSpace(permission)
	}
}
