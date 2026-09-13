package auth

import (
	"time"

	"xermess/internal/model"
)

// adminResponse is an administrator as the browser sees it. It is built by
// hand rather than returning the model, so a column added later cannot
// accidentally start being published.
type adminResponse struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
	Status    string     `json:"status"`
	Roles     []string   `json:"roles"`
	LastLogin *time.Time `json:"last_login_at,omitempty"`
}

func newAdminResponse(a *model.AdminUser) adminResponse {
	roles := make([]string, 0, len(a.Roles))
	for _, role := range a.Roles {
		roles = append(roles, string(role.Name))
	}

	return adminResponse{
		ID:        a.ID.String(),
		Username:  a.Username,
		Email:     a.Email,
		FullName:  a.FullName(),
		Status:    string(a.Status),
		Roles:     roles,
		LastLogin: a.LastLoginAt,
	}
}

// sessionResponse is one of the places an administrator is signed in.
type sessionResponse struct {
	ID        string    `json:"id"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Active    bool      `json:"active"`
}

func newSessionResponses(sessions []model.AdminUserSession) []sessionResponse {
	now := time.Now()
	out := make([]sessionResponse, 0, len(sessions))

	for _, s := range sessions {
		out = append(out, sessionResponse{
			ID:        s.ID.String(),
			IP:        s.IP,
			UserAgent: s.UserAgent,
			CreatedAt: s.CreatedAt,
			ExpiresAt: s.ExpiresAt,
			Active:    s.IsActive(now),
		})
	}

	return out
}
