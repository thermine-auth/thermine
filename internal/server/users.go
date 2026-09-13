package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"xermess/internal/model"
)

// maxPageSize caps how many users one request can ask for.
const maxPageSize = 200

// userRequest is the body of the create and update endpoints. Data carries
// the user-defined fields, which are checked against user_fields before
// anything is written.
type userRequest struct {
	Email         string         `json:"email"`
	EmailVerified bool           `json:"email_verified"`
	Data          map[string]any `json:"data"`
}

// listUsers returns a page of users, newest first.
//
// `search` matches the email or any text the user-defined fields hold, so one
// box covers the whole record rather than one column.
func (h *admin) listUsers(c *gin.Context) {
	ctx := c.Request.Context()

	limit := intQuery(c, "limit", 50, maxPageSize)
	offset := intQuery(c, "offset", 0, 1_000_000)

	query := h.db.WithContext(ctx).Model(&model.User{})

	if search := strings.TrimSpace(c.Query("search")); search != "" {
		like := "%" + strings.ToLower(search) + "%"
		// data is stored as JSON; casting to text searches every field at
		// once, which is what a single search box should do.
		query = query.Where("LOWER(email) LIKE ? OR LOWER(data::text) LIKE ?", like, like)
	}

	if verified := c.Query("verified"); verified == "true" || verified == "false" {
		query = query.Where("email_verified = ?", verified == "true")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		h.log.Error("counting users failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	var users []model.User
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		h.log.Error("listing users failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":  users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// getUser returns one user.
func (h *admin) getUser(c *gin.Context) {
	user, ok := h.findUser(c)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// createUser adds a user.
func (h *admin) createUser(c *gin.Context) {
	var req userRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the request body is not valid"})
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	data, err := h.normalise(c, req.Data, uuid.Nil)
	if err != nil {
		return
	}

	user := model.User{Email: email, EmailVerified: req.EmailVerified, Data: data}
	if err := h.db.WithContext(c.Request.Context()).Create(&user).Error; err != nil {
		h.respondWriteError(c, err, "creating user failed")
		return
	}

	h.record(c, "user.created", user.ID.String())

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// updateUser replaces a user's email, verified flag and fields.
func (h *admin) updateUser(c *gin.Context) {
	user, ok := h.findUser(c)
	if !ok {
		return
	}

	var req userRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the request body is not valid"})
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	data, err := h.normalise(c, req.Data, user.ID)
	if err != nil {
		return
	}

	user.Email = email
	user.EmailVerified = req.EmailVerified
	user.Data = data

	if err := h.db.WithContext(c.Request.Context()).Save(user).Error; err != nil {
		h.respondWriteError(c, err, "updating user failed")
		return
	}

	h.record(c, "user.updated", user.ID.String())

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// deleteUser removes a user for good.
func (h *admin) deleteUser(c *gin.Context) {
	user, ok := h.findUser(c)
	if !ok {
		return
	}

	if err := h.db.WithContext(c.Request.Context()).Unscoped().Delete(user).Error; err != nil {
		h.log.Error("deleting user failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	h.record(c, "user.deleted", user.ID.String())

	c.Status(http.StatusNoContent)
}

// findUser loads the user named in the path, answering the request itself if
// the id is not a uuid or there is no such user.
func (h *admin) findUser(c *gin.Context) (*model.User, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "that is not a user id"})
		return nil, false
	}

	var user model.User
	err = h.db.WithContext(c.Request.Context()).First(&user, "id = ?", id).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "no such user"})
		return nil, false
	case err != nil:
		h.log.Error("loading user failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return nil, false
	}

	return &user, true
}

// normalise checks the submitted values against the field definitions and
// returns what should be stored. It answers the request itself on a bad
// value, and the caller stops.
//
// `self` is the user being written, so a field that has to be unique does not
// find the record's own value and call it a clash. It is uuid.Nil when the
// user is being created.
func (h *admin) normalise(c *gin.Context, values map[string]any, self uuid.UUID) (map[string]any, error) {
	var fields []model.UserField
	if err := h.db.WithContext(c.Request.Context()).Order("position").Find(&fields).Error; err != nil {
		h.log.Error("loading user fields failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return nil, err
	}

	data := make(map[string]any, len(fields))

	for _, field := range fields {
		normalised, err := field.Normalise(values[field.Name])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return nil, err
		}

		// A field with no value is left out rather than stored as null, so
		// removing a field leaves nothing behind.
		if normalised == nil {
			continue
		}

		if field.Unique {
			taken, err := h.valueTaken(c, field, normalised, self)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
				return nil, err
			}

			if taken {
				c.JSON(http.StatusConflict, gin.H{
					"error": fmt.Sprintf("%s: another user already has that value", field.Name),
				})
				return nil, errDuplicateValue
			}
		}

		data[field.Name] = normalised
	}

	return data, nil
}

// errDuplicateValue stops the caller once the answer has been written.
var errDuplicateValue = errors.New("value already taken")

// valueTaken reports whether another user already holds this value for a
// field that has to be unique. The values live in a JSON column, so the check
// is a query rather than an index: correct, and fine at this scale.
func (h *admin) valueTaken(c *gin.Context, field model.UserField, value any, self uuid.UUID) (bool, error) {
	text := fmt.Sprintf("%v", value)
	if number, ok := value.(float64); ok {
		text = strconv.FormatFloat(number, 'f', -1, 64)
	}

	query := h.db.WithContext(c.Request.Context()).
		Model(&model.User{}).
		Where("data::jsonb ->> ? = ?", field.Name, text)

	if self != uuid.Nil {
		query = query.Where("id <> ?", self)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		h.log.Error("checking a unique field failed", "field", field.Name, "error", err)
		return false, err
	}

	return count > 0, nil
}

// respondWriteError turns a failed write into an answer, telling a duplicate
// email apart from anything else because that one is the writer's to fix.
func (h *admin) respondWriteError(c *gin.Context, err error, message string) {
	if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		c.JSON(http.StatusConflict, gin.H{"error": "a user with that email already exists"})
		return
	}

	h.log.Error(message, "error", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
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
