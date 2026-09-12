package server

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"xermess/internal/auth"
	"xermess/internal/model"
)

// sessionCookie carries the session token between the browser and the server.
// It is HttpOnly so no script can read it, which is what keeps a cross-site
// scripting bug from turning into a stolen session.
const sessionCookie = "xermess_session"

// admin holds what the admin endpoints need.
type admin struct {
	auth   *auth.Service
	db     *gorm.DB
	log    *slog.Logger
	secure bool // set the cookie's Secure flag: true once served over HTTPS
}

// loginRequest is the body of POST /admin/auth/login.
type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

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

// login checks the credentials and sets the session cookie.
func (h *admin) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	token, user, err := h.auth.Login(c.Request.Context(), req.Username, req.Password, requestOf(c))
	switch {
	case errors.Is(err, auth.ErrInvalidCredentials):
		// One message for every kind of failure: which usernames exist is not
		// something a sign-in page should reveal.
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong username or password"})
		return
	case err != nil:
		h.log.Error("login failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	h.setSessionCookie(c, token)

	c.JSON(http.StatusOK, gin.H{"admin": newAdminResponse(user)})
}

// logout revokes the session and clears the cookie.
func (h *admin) logout(c *gin.Context) {
	token, _ := c.Cookie(sessionCookie)

	if err := h.auth.Logout(c.Request.Context(), token, requestOf(c)); err != nil {
		h.log.Error("logout failed", "error", err)
	}

	h.clearSessionCookie(c)

	c.JSON(http.StatusOK, gin.H{"status": "signed out"})
}

// me returns the signed-in administrator, and is what the browser calls to
// find out whether it still has a session.
func (h *admin) me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"admin": newAdminResponse(adminFrom(c))})
}

// overview is the admin panel's front page: a few counts and the latest
// activity.
func (h *admin) overview(c *gin.Context) {
	ctx := c.Request.Context()

	var counts struct {
		Admins         int64 `json:"admins"`
		ActiveSessions int64 `json:"active_sessions"`
		Roles          int64 `json:"roles"`
		Events         int64 `json:"events"`
	}

	if err := h.db.WithContext(ctx).Model(&model.AdminUser{}).Count(&counts.Admins).Error; err != nil {
		h.log.Error("overview failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}
	h.db.WithContext(ctx).Model(&model.AdminUserSession{}).
		Where("revoked_at IS NULL AND expires_at > ?", time.Now()).Count(&counts.ActiveSessions)
	h.db.WithContext(ctx).Model(&model.Role{}).Count(&counts.Roles)
	h.db.WithContext(ctx).Model(&model.AuditLog{}).Count(&counts.Events)

	var events []model.AuditLog
	h.db.WithContext(ctx).Order("created_at DESC").Limit(10).Find(&events)

	activity := make([]gin.H, 0, len(events))
	for _, event := range events {
		activity = append(activity, gin.H{
			"id":         event.ID.String(),
			"action":     event.Action,
			"actor":      event.ActorEmail,
			"ip":         event.IP,
			"created_at": event.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"counts": counts, "activity": activity})
}

// logs lists the activity log, newest first. The page size is capped so a
// caller cannot ask for the whole table.
func (h *admin) logs(c *gin.Context) {
	const defaultLimit, maxLimit = 50, 200

	limit := defaultLimit
	if raw := c.Query("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = min(parsed, maxLimit)
		}
	}

	var events []model.AuditLog
	if err := h.db.WithContext(c.Request.Context()).
		Order("created_at DESC").Limit(limit).Find(&events).Error; err != nil {
		h.log.Error("listing logs failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	out := make([]gin.H, 0, len(events))
	for _, event := range events {
		out = append(out, gin.H{
			"id":          event.ID.String(),
			"action":      event.Action,
			"actor":       event.ActorEmail,
			"ip":          event.IP,
			"user_agent":  event.UserAgent,
			"target_type": event.TargetType,
			"created_at":  event.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"logs": out})
}

// sessions lists the caller's own sessions, so they can see where they are
// signed in.
func (h *admin) sessions(c *gin.Context) {
	var sessions []model.AdminUserSession
	h.db.WithContext(c.Request.Context()).
		Where("admin_user_id = ?", adminFrom(c).ID).
		Order("created_at DESC").Limit(20).Find(&sessions)

	out := make([]gin.H, 0, len(sessions))
	for _, session := range sessions {
		out = append(out, gin.H{
			"id":         session.ID.String(),
			"ip":         session.IP,
			"user_agent": session.UserAgent,
			"created_at": session.CreatedAt,
			"expires_at": session.ExpiresAt,
			"active":     session.IsActive(time.Now()),
		})
	}

	c.JSON(http.StatusOK, gin.H{"sessions": out})
}

// setSessionCookie stores the token in the browser.
func (h *admin) setSessionCookie(c *gin.Context, token string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		MaxAge:   int(auth.SessionLifetime.Seconds()),
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// clearSessionCookie removes it again.
func (h *admin) clearSessionCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// requestOf describes where the call came from, for the session and the log.
func requestOf(c *gin.Context) auth.Request {
	return auth.Request{
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
}
