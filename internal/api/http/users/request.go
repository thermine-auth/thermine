package users

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"xermess/internal/store"
)

// targetType is what these records are called in the activity log.
const targetType = "user"

// maxPageSize caps how many users one request can ask for.
const maxPageSize = 200

// userRequest is the body of the create and update endpoints.
//
// The email is the account: it is the only column a user record has of its
// own, and what someone signs in with. Data carries the user-defined fields,
// whose rules live in user_fields and are checked in validation.go rather
// than in a tag here.
type userRequest struct {
	Email         string         `json:"email" validate:"required,email,max=255"`
	EmailVerified bool           `json:"email_verified"`
	Data          map[string]any `json:"data"`
}

// clean tidies what can be tidied, so the rules below see the value that
// would actually be stored.
func (r *userRequest) clean() {
	r.Email = trimmedLower(r.Email)
}

// listQuery reads the search box, the filter and the page out of the query
// string, holding each to something sensible.
func listQuery(c *gin.Context) store.UserQuery {
	query := store.UserQuery{
		Search: c.Query("search"),
		Limit:  intQuery(c, "limit", 50, maxPageSize),
		Offset: intQuery(c, "offset", 0, 1_000_000),
	}

	if verified := c.Query("verified"); verified == "true" || verified == "false" {
		wanted := verified == "true"
		query.Verified = &wanted
	}

	return query
}

// intQuery reads a whole number from the query string, falling back to a
// default and refusing anything past the cap.
func intQuery(c *gin.Context, name string, fallback, max int) int {
	raw := c.Query(name)
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return fallback
	}

	return min(value, max)
}

// trimmedLower is how an email is stored: no stray spaces, one case, so two
// spellings of the same address cannot both be registered.
func trimmedLower(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
