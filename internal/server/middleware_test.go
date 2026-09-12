package server

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestRequestLogger checks the one line a finished request leaves behind.
func TestRequestLogger(t *testing.T) {
	var logged bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logged, nil))

	r := gin.New()
	r.Use(RequestLogger(log))
	r.GET("/thing", func(c *gin.Context) { c.Status(http.StatusTeapot) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/thing?q=1", nil))

	line := logged.String()
	for _, want := range []string{
		`msg=request`,
		`method=GET`,
		`path=/thing`,
		`status=418`,
		`duration=`,
		`ip=`,
	} {
		if !strings.Contains(line, want) {
			t.Errorf("log line is missing %q:\n%s", want, line)
		}
	}
}

// The line is written after the handler runs, so it can report the status the
// handler chose rather than the zero value.
func TestRequestLoggerReportsHandlerStatus(t *testing.T) {
	var logged bytes.Buffer
	log := slog.New(slog.NewTextHandler(&logged, nil))

	r := gin.New()
	r.Use(RequestLogger(log))
	r.GET("/missing", func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/missing", nil))

	if !strings.Contains(logged.String(), "status=404") {
		t.Errorf("want status=404 in the log, got:\n%s", logged.String())
	}
}

// A panic must become a 500, not a dropped connection.
func TestRecoveryTurnsPanicIntoError(t *testing.T) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/boom", func(*gin.Context) { panic("boom") })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/boom", nil))

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
}
