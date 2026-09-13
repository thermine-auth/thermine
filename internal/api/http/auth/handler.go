// Package auth answers the endpoints that sign an administrator in and out,
// and the ones that say who they are: the panel's whole idea of a session.
package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"xermess/internal/api/http/respond"
	"xermess/internal/api/http/session"
	authsvc "xermess/internal/auth"
	"xermess/internal/store"
)

// Handler holds what these endpoints need.
type Handler struct {
	auth  *authsvc.Service
	store *store.Store
	log   *slog.Logger

	// secure sets the cookie's Secure flag: true once served over HTTPS.
	secure bool
}

// New returns a Handler.
func New(service *authsvc.Service, st *store.Store, log *slog.Logger, secure bool) *Handler {
	return &Handler{auth: service, store: st, log: log, secure: secure}
}

// Login checks the credentials and sets the session cookie.
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.BadRequest(c, "the request body is not valid")
		return
	}

	if err := req.validate(); err != nil {
		respond.Failure(c, h.log, err, "validating a sign-in failed")
		return
	}

	token, admin, err := h.auth.Login(c.Request.Context(), req.Username, req.Password, requestOf(c))
	switch {
	case errors.Is(err, authsvc.ErrInvalidCredentials):
		// One message for every kind of failure: which usernames exist is not
		// something a sign-in page should reveal.
		respond.Error(c, http.StatusUnauthorized, "wrong username or password")
		return
	case err != nil:
		respond.Failure(c, h.log, err, "login failed")
		return
	}

	session.Set(c, token, h.secure)

	c.JSON(http.StatusOK, gin.H{"admin": newAdminResponse(admin)})
}

// Logout revokes the session and clears the cookie.
func (h *Handler) Logout(c *gin.Context) {
	token, _ := c.Cookie(session.Cookie)

	if err := h.auth.Logout(c.Request.Context(), token, requestOf(c)); err != nil {
		h.log.Error("logout failed", "error", err)
	}

	session.Clear(c, h.secure)

	c.JSON(http.StatusOK, gin.H{"status": "signed out"})
}

// Me returns the signed-in administrator, and is what the browser calls to
// find out whether it still has a session.
func (h *Handler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"admin": newAdminResponse(session.Admin(c))})
}

// Sessions lists the caller's own sessions, so they can see where they are
// signed in.
func (h *Handler) Sessions(c *gin.Context) {
	sessions, err := h.store.SessionsFor(c.Request.Context(), session.Admin(c).ID, 20)
	if err != nil {
		respond.Failure(c, h.log, err, "listing sessions failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": newSessionResponses(sessions)})
}

// requestOf describes where the call came from, for the session and the log.
func requestOf(c *gin.Context) authsvc.Request {
	return authsvc.Request{
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
}
