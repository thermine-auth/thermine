package server

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"xermess/internal/model"
)

// fieldName is what a field may be called: the same shape as a column name,
// because that is what it stands in for.
var fieldName = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)

// fieldRequest is the body of the create and update endpoints. The rules a
// value has to keep — unique, bounded, prefixed — are part of the field, so
// they are set here rather than in each record.
type fieldRequest struct {
	Name       string   `json:"name"`
	Label      string   `json:"label"`
	Type       string   `json:"type"`
	Required   bool     `json:"required"`
	Unique     bool     `json:"unique"`
	Min        *float64 `json:"min"`
	Max        *float64 `json:"max"`
	StartsWith string   `json:"starts_with"`
}

// listUserFields returns the fields a user record has, in the order they are
// shown. The panel builds both its table and its form from this.
func (h *admin) listUserFields(c *gin.Context) {
	var fields []model.UserField
	if err := h.db.WithContext(c.Request.Context()).Order("position, created_at").Find(&fields).Error; err != nil {
		h.log.Error("listing user fields failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"fields": fields, "types": model.FieldTypes})
}

// createUserField adds a field to every user record. Existing users simply
// have no value for it until they are edited.
func (h *admin) createUserField(c *gin.Context) {
	var req fieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the request body is not valid"})
		return
	}

	name := strings.TrimSpace(strings.ToLower(req.Name))
	if !fieldName.MatchString(name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "a field name starts with a letter and holds only lower case letters, numbers and underscores",
		})
		return
	}

	fieldType := model.FieldType(strings.TrimSpace(req.Type))
	if !fieldType.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "that is not a field type"})
		return
	}

	label := strings.TrimSpace(req.Label)
	if label == "" {
		label = name
	}

	// New fields go to the end.
	var last model.UserField
	position := 1
	if err := h.db.WithContext(c.Request.Context()).Order("position DESC").First(&last).Error; err == nil {
		position = last.Position + 1
	}

	field := model.UserField{
		Name:       name,
		Label:      label,
		Type:       fieldType,
		Required:   req.Required,
		Unique:     req.Unique,
		Min:        req.Min,
		Max:        req.Max,
		StartsWith: strings.TrimSpace(req.StartsWith),
		Position:   position,
	}

	if err := field.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.WithContext(c.Request.Context()).Create(&field).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "a field with that name already exists"})
			return
		}

		h.log.Error("creating user field failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	h.record(c, "user_field.created", field.Name)

	c.JSON(http.StatusCreated, gin.H{"field": field})
}

// updateUserField changes what a field expects. Its name and its type stay
// as they are: records already hold values under that name and in that
// shape, and changing either here would leave them behind.
func (h *admin) updateUserField(c *gin.Context) {
	field, ok := h.findUserField(c)
	if !ok {
		return
	}

	var req fieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the request body is not valid"})
		return
	}

	label := strings.TrimSpace(req.Label)
	if label == "" {
		label = field.Name
	}

	field.Label = label
	field.Required = req.Required
	field.Unique = req.Unique
	field.Min = req.Min
	field.Max = req.Max
	field.StartsWith = strings.TrimSpace(req.StartsWith)

	if err := field.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.WithContext(c.Request.Context()).Save(field).Error; err != nil {
		h.log.Error("updating user field failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	h.record(c, "user_field.updated", field.Name)

	c.JSON(http.StatusOK, gin.H{"field": field})
}

// deleteUserField removes a field. The values already stored under its name
// stay in the user records until those are next saved, at which point they
// are dropped: nothing is destroyed by removing a column from the panel.
func (h *admin) deleteUserField(c *gin.Context) {
	field, ok := h.findUserField(c)
	if !ok {
		return
	}

	if err := h.db.WithContext(c.Request.Context()).Unscoped().Delete(field).Error; err != nil {
		h.log.Error("deleting user field failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return
	}

	h.record(c, "user_field.deleted", field.Name)

	c.Status(http.StatusNoContent)
}

// findUserField loads the field named in the path, answering the request
// itself if there is no such field.
func (h *admin) findUserField(c *gin.Context) (*model.UserField, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "that is not a field id"})
		return nil, false
	}

	var field model.UserField
	err = h.db.WithContext(c.Request.Context()).First(&field, "id = ?", id).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "no such field"})
		return nil, false
	case err != nil:
		h.log.Error("loading user field failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
		return nil, false
	}

	return &field, true
}
