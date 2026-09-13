// Package respond writes the answers that are not the happy path, so every
// endpoint fails the same way: one "error" field, and nothing about the
// server's insides.
package respond

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Fault is a problem with what was sent: the status to answer with, and what
// to say about it. Anything that is not a Fault is the server's own fault,
// and the caller is told nothing more than that something went wrong.
type Fault struct {
	Status  int
	Message string
}

func (f Fault) Error() string {
	return f.Message
}

// BadRequest says the request itself was wrong.
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// NotFound says there is no such thing.
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// Conflict says the request clashes with what is already stored.
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message)
}

// Error answers with a status and a message.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// Failure answers whatever the error turns out to be: a Fault is passed on to
// the caller as it stands, and anything else is logged with `note` and
// answered with a 500.
//
// Handlers call this instead of picking a status themselves, which is what
// keeps a database error from ever reaching a browser.
func Failure(c *gin.Context, log *slog.Logger, err error, note string) {
	var fault Fault
	if errors.As(err, &fault) {
		Error(c, fault.Status, fault.Message)
		return
	}

	log.Error(note, "error", err)
	Error(c, http.StatusInternalServerError, "something went wrong")
}
