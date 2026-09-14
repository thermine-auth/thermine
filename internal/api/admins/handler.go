// Package admins answers the endpoints a super admin manages the panel's
// administrators with: who they are, which admin roles they hold, and whether
// their account may sign in.
package admins

import (
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"xermess/internal/api/audit"
	"xermess/internal/api/respond"
	"xermess/internal/api/session"
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

// List returns a page of administrators, newest first.
func (h *Handler) List(c *gin.Context) {
	query := listQuery(c)

	admins, total, err := h.store.Admins(c.Request.Context(), query)
	if err != nil {
		respond.Failure(c, h.log, err, "listing administrators failed")
		return
	}

	c.JSON(http.StatusOK, newPageResponse(admins, total, query))
}

// Get returns one administrator.
func (h *Handler) Get(c *gin.Context) {
	admin, ok := h.find(c)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"admin": newAdminResponse(*admin)})
}

// Create adds an administrator, with the password and roles the super admin
// chose for them.
func (h *Handler) Create(c *gin.Context) {
	var req adminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.BadRequest(c, "the request body is not valid")
		return
	}

	admin, err := h.build(c, &req, nil)
	if err != nil {
		respond.Failure(c, h.log, err, "checking an administrator failed")
		return
	}

	if err := h.store.CreateAdmin(c.Request.Context(), admin); err != nil {
		h.respondWrite(c, err, "creating administrator failed")
		return
	}

	h.audit.Record(c, "admin.created", targetType, admin.ID.String())

	h.answer(c, http.StatusCreated, admin)
}

// Update replaces an administrator's details, status and roles, and their
// password when a new one is given. A new password or an account that may no
// longer sign in ends every session the administrator has open.
func (h *Handler) Update(c *gin.Context) {
	admin, ok := h.find(c)
	if !ok {
		return
	}

	var req adminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.BadRequest(c, "the request body is not valid")
		return
	}

	if _, err := h.build(c, &req, admin); err != nil {
		respond.Failure(c, h.log, err, "checking an administrator failed")
		return
	}

	ctx := c.Request.Context()

	if err := h.store.SaveAdmin(ctx, admin); err != nil {
		h.respondWrite(c, err, "updating administrator failed")
		return
	}

	h.audit.Record(c, "admin.updated", targetType, admin.ID.String())
	if req.Password != "" {
		h.audit.Record(c, "admin.password_changed", targetType, admin.ID.String())
	}

	// Your own sessions are left alone, or saving your own password here
	// would sign you out of the page you saved it from.
	changedAccess := req.Password != "" || admin.Status != model.StatusActive
	if changedAccess && admin.ID != session.Admin(c).ID {
		if err := h.store.RevokeSessionsFor(ctx, admin.ID, time.Now()); err != nil {
			respond.Failure(c, h.log, err, "ending the administrator's sessions failed")
			return
		}
	}

	h.answer(c, http.StatusOK, admin)
}

// answer writes an administrator as stored, read back so each assignment
// carries its role's and application's names.
func (h *Handler) answer(c *gin.Context, status int, admin *model.AdminUser) {
	stored, err := h.store.AdminByID(c.Request.Context(), admin.ID)
	if err != nil {
		respond.Failure(c, h.log, err, "reading the administrator back failed")
		return
	}

	c.JSON(status, gin.H{"admin": newAdminResponse(*stored)})
}

// Delete removes an administrator for good. Nobody can remove themselves:
// that is how a panel ends up with no one able to manage it.
func (h *Handler) Delete(c *gin.Context) {
	admin, ok := h.find(c)
	if !ok {
		return
	}

	if admin.ID == session.Admin(c).ID {
		respond.Conflict(c, "you cannot delete your own account")
		return
	}

	if err := h.store.DeleteAdmin(c.Request.Context(), admin); err != nil {
		respond.Failure(c, h.log, err, "deleting administrator failed")
		return
	}

	h.audit.Record(c, "admin.deleted", targetType, admin.ID.String())

	c.Status(http.StatusNoContent)
}

// build checks a submitted administrator and returns the account to write.
// Passing an existing one fills that one in instead of making a new one.
//
// A super admin may not suspend themselves or take away their own whole-panel
// super_admin role: the request that did it would be the last one they could
// make, and this endpoint is the only way back.
func (h *Handler) build(c *gin.Context, req *adminRequest, into *model.AdminUser) (*model.AdminUser, error) {
	if err := req.validate(into == nil); err != nil {
		return nil, err
	}

	assignments, err := h.assignments(c, req.Assignments)
	if err != nil {
		return nil, err
	}

	status := model.Status(req.Status)

	if into != nil && into.ID == session.Admin(c).ID {
		keepsSuper := slices.ContainsFunc(assignments, func(a model.AdminRoleAssignment) bool {
			return a.Global() && a.Role.IsSuperAdmin()
		})

		if status != model.StatusActive || !keepsSuper {
			return nil, respond.Fault{
				Status:  http.StatusConflict,
				Message: "you cannot suspend yourself or take away your own super_admin role",
			}
		}
	}

	admin := into
	if admin == nil {
		admin = &model.AdminUser{}
	}

	if req.Password != "" {
		// The password is never stored, only this hash of it.
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return nil, respond.Fault{Status: http.StatusBadRequest, Message: "password must be at most 72 bytes"}
		}
		if err != nil {
			return nil, err
		}

		admin.PasswordHash = string(hash)

		// A password set by a super admin is the way back into a locked
		// account, so the wrong guesses that locked it are forgotten.
		admin.FailedLoginCount = 0
		admin.LockedUntil = nil
	}

	// The address is the account, the same as for the first administrator:
	// it is what they sign in with.
	admin.Username = req.Email
	admin.Email = req.Email
	admin.FirstName = req.FirstName
	admin.LastName = req.LastName
	admin.Status = status
	admin.Assignments = assignments

	return admin, nil
}

// assignments turns the submitted role assignments into rows, checking every
// role and application exists. The same role for the same scope twice is one
// assignment, and super_admin can only be held for the whole panel.
func (h *Handler) assignments(c *gin.Context, requested []assignmentRequest) ([]model.AdminRoleAssignment, error) {
	ctx := c.Request.Context()

	roleIDs := make([]uuid.UUID, 0, len(requested))
	for _, a := range requested {
		roleIDs = append(roleIDs, a.RoleID)
	}

	roles, err := h.store.AdminRolesByID(ctx, roleIDs)
	if errors.Is(err, store.ErrNotFound) {
		return nil, respond.Fault{Status: http.StatusBadRequest, Message: "assignments: a role does not exist"}
	}
	if err != nil {
		return nil, err
	}

	byID := make(map[uuid.UUID]model.Role, len(roles))
	for _, role := range roles {
		byID[role.ID] = role
	}

	out := []model.AdminRoleAssignment{}
	seen := map[string]bool{}

	for _, a := range requested {
		role := byID[a.RoleID]

		key := a.RoleID.String()
		if a.ApplicationID != nil {
			if role.IsSuperAdmin() {
				return nil, respond.Fault{
					Status:  http.StatusBadRequest,
					Message: "assignments: super_admin can only be held for the whole panel",
				}
			}

			if _, err := h.store.Application(ctx, *a.ApplicationID); err != nil {
				if errors.Is(err, store.ErrNotFound) {
					return nil, respond.Fault{Status: http.StatusBadRequest, Message: "assignments: an application does not exist"}
				}
				return nil, err
			}

			key += "/" + a.ApplicationID.String()
		}

		if seen[key] {
			continue
		}
		seen[key] = true

		out = append(out, model.AdminRoleAssignment{RoleID: role.ID, Role: role, ApplicationID: a.ApplicationID})
	}

	return out, nil
}

// find loads the administrator named in the path, answering the request
// itself if the id is not a uuid or there is no such administrator.
func (h *Handler) find(c *gin.Context) (*model.AdminUser, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.BadRequest(c, "that is not an administrator id")
		return nil, false
	}

	admin, err := h.store.AdminByID(c.Request.Context(), id)

	switch {
	case errors.Is(err, store.ErrNotFound):
		respond.NotFound(c, "no such administrator")
		return nil, false
	case err != nil:
		respond.Failure(c, h.log, err, "loading administrator failed")
		return nil, false
	}

	return admin, true
}

// respondWrite turns a failed write into an answer, telling a taken address
// apart from anything else because that one is the writer's to fix.
func (h *Handler) respondWrite(c *gin.Context, err error, note string) {
	if errors.Is(err, store.ErrDuplicate) {
		respond.Conflict(c, "an administrator with that email already exists")
		return
	}

	respond.Failure(c, h.log, err, note)
}
