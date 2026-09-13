// Package users answers the endpoints for the people an organisation manages:
// the records themselves, as opposed to the administrators who edit them.
package users

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"xermess/internal/api/http/audit"
	"xermess/internal/api/http/respond"
	"xermess/internal/model"
	"xermess/internal/store"
)

// Handler holds what these endpoints need.
type Handler struct {
	store *store.Store
	audit audit.Recorder
	log   *slog.Logger
}

// New returns a Handler.
func New(st *store.Store, recorder audit.Recorder, log *slog.Logger) *Handler {
	return &Handler{store: st, audit: recorder, log: log}
}

// List returns a page of users, newest first.
//
// `search` matches the email or any text the user-defined fields hold, so one
// box covers the whole record rather than one column.
func (h *Handler) List(c *gin.Context) {
	query := listQuery(c)

	users, total, err := h.store.Users(c.Request.Context(), query)
	if err != nil {
		respond.Failure(c, h.log, err, "listing users failed")
		return
	}

	c.JSON(http.StatusOK, newPageResponse(users, total, query))
}

// Get returns one user.
func (h *Handler) Get(c *gin.Context) {
	user, ok := h.find(c)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Create adds a user.
func (h *Handler) Create(c *gin.Context) {
	var req userRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.BadRequest(c, "the request body is not valid")
		return
	}

	user, err := h.build(c, &req, nil)
	if err != nil {
		respond.Failure(c, h.log, err, "checking a user failed")
		return
	}

	if err := h.store.CreateUser(c.Request.Context(), user); err != nil {
		h.respondWrite(c, err, "creating user failed")
		return
	}

	h.audit.Record(c, "user.created", targetType, user.ID.String())

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// Update replaces a user's email, verified flag and fields.
func (h *Handler) Update(c *gin.Context) {
	user, ok := h.find(c)
	if !ok {
		return
	}

	var req userRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.BadRequest(c, "the request body is not valid")
		return
	}

	if _, err := h.build(c, &req, user); err != nil {
		respond.Failure(c, h.log, err, "checking a user failed")
		return
	}

	if err := h.store.SaveUser(c.Request.Context(), user); err != nil {
		h.respondWrite(c, err, "updating user failed")
		return
	}

	h.audit.Record(c, "user.updated", targetType, user.ID.String())

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Delete removes a user for good.
func (h *Handler) Delete(c *gin.Context) {
	user, ok := h.find(c)
	if !ok {
		return
	}

	if err := h.store.DeleteUser(c.Request.Context(), user); err != nil {
		respond.Failure(c, h.log, err, "deleting user failed")
		return
	}

	h.audit.Record(c, "user.deleted", targetType, user.ID.String())

	c.Status(http.StatusNoContent)
}

// build checks a submitted record and returns the user to write. Passing an
// existing user fills that one in instead of making a new one, and keeps a
// unique field from clashing with the record's own value.
func (h *Handler) build(c *gin.Context, req *userRequest, into *model.User) (*model.User, error) {
	ctx := c.Request.Context()

	if err := req.validate(); err != nil {
		return nil, err
	}

	fields, err := h.store.UserFields(ctx)
	if err != nil {
		return nil, err
	}

	self := uuid.Nil
	if into != nil {
		self = into.ID
	}

	data, err := normalise(ctx, h.store, fields, req.Data, self)
	if err != nil {
		return nil, err
	}

	user := into
	if user == nil {
		user = &model.User{}
	}

	user.Email = req.Email
	user.EmailVerified = req.EmailVerified
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.IsActive = req.IsActive
	user.Data = data

	return user, nil
}

// find loads the user named in the path, answering the request itself if the
// id is not a uuid or there is no such user.
func (h *Handler) find(c *gin.Context) (*model.User, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.BadRequest(c, "that is not a user id")
		return nil, false
	}

	user, err := h.store.User(c.Request.Context(), id)

	switch {
	case errors.Is(err, store.ErrNotFound):
		respond.NotFound(c, "no such user")
		return nil, false
	case err != nil:
		respond.Failure(c, h.log, err, "loading user failed")
		return nil, false
	}

	return user, true
}

// respondWrite turns a failed write into an answer, telling a duplicate email
// apart from anything else because that one is the writer's to fix. The
// address is the only column of a user record that has to be unique; an
// additional field that has to be is checked before the write, in
// validation.go.
func (h *Handler) respondWrite(c *gin.Context, err error, note string) {
	if errors.Is(err, store.ErrDuplicate) {
		respond.Conflict(c, "a user with that email already exists")
		return
	}

	respond.Failure(c, h.log, err, note)
}
