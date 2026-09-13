// Package activity answers the endpoints that report on what has been
// happening: the dashboard's counts, and the log of what administrators did.
package activity

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"xermess/internal/api/http/respond"
	"xermess/internal/store"
)

// How many log entries an endpoint returns.
const (
	dashboardEntries = 10
	defaultLimit     = 50
	maxLimit         = 200
)

// Handler holds what these endpoints need.
type Handler struct {
	store *store.Store
	log   *slog.Logger
}

// New returns a Handler.
func New(st *store.Store, log *slog.Logger) *Handler {
	return &Handler{store: st, log: log}
}

// Overview is the admin panel's front page: a few counts and the latest
// activity.
func (h *Handler) Overview(c *gin.Context) {
	ctx := c.Request.Context()

	counts, err := h.store.Counts(ctx, time.Now())
	if err != nil {
		respond.Failure(c, h.log, err, "counting for the overview failed")
		return
	}

	events, err := h.store.AuditLog(ctx, dashboardEntries)
	if err != nil {
		respond.Failure(c, h.log, err, "reading the activity log failed")
		return
	}

	c.JSON(http.StatusOK, newOverviewResponse(counts, events))
}

// Logs lists the activity log, newest first. The page size is capped so a
// caller cannot ask for the whole table.
func (h *Handler) Logs(c *gin.Context) {
	limit := defaultLimit
	if raw := c.Query("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = min(parsed, maxLimit)
		}
	}

	events, err := h.store.AuditLog(c.Request.Context(), limit)
	if err != nil {
		respond.Failure(c, h.log, err, "listing logs failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": newLogResponses(events)})
}
