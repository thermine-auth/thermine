package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"xermess/internal/config"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// testEngine builds the real server, with its logging discarded.
func testEngine(origins ...string) *gin.Engine {
	cfg := config.Config{Addr: ":0", CORSOrigins: origins}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	return New(cfg, nil, log)
}

// do sends a request through the router and returns the response.
func do(r *gin.Engine, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func TestRoutes(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		wantCode int
		wantBody map[string]string
	}{
		{
			name: "health check", method: http.MethodGet, path: "/healthz",
			wantCode: http.StatusOK, wantBody: map[string]string{"status": "ok"},
		},
		{
			name: "hello", method: http.MethodGet, path: "/api/v1/hello",
			wantCode: http.StatusOK, wantBody: map[string]string{"hello": "world"},
		},
		{
			name: "unknown path", method: http.MethodGet, path: "/nope",
			wantCode: http.StatusNotFound, wantBody: map[string]string{"error": "not found"},
		},
		{
			name: "unknown path under the api group", method: http.MethodGet, path: "/api/v1/nope",
			wantCode: http.StatusNotFound, wantBody: map[string]string{"error": "not found"},
		},
	}

	r := testEngine()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := do(r, tt.method, tt.path, nil)

			if w.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d (body: %s)", w.Code, tt.wantCode, w.Body)
			}

			var got map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("body is not JSON: %s", w.Body)
			}
			for key, want := range tt.wantBody {
				if got[key] != want {
					t.Errorf("body[%q] = %q, want %q", key, got[key], want)
				}
			}
		})
	}
}

func TestCORS(t *testing.T) {
	const allowed = "http://localhost:5173"

	tests := []struct {
		name        string
		origin      string
		method      string
		wantCode    int
		wantAllowed bool
	}{
		{name: "allowed origin", origin: allowed, method: http.MethodGet, wantCode: http.StatusOK, wantAllowed: true},
		{name: "unknown origin", origin: "http://evil.test", method: http.MethodGet, wantCode: http.StatusOK},
		{name: "no origin", origin: "", method: http.MethodGet, wantCode: http.StatusOK},
		{name: "preflight from allowed origin", origin: allowed, method: http.MethodOptions, wantCode: http.StatusNoContent, wantAllowed: true},
		{name: "preflight from unknown origin", origin: "http://evil.test", method: http.MethodOptions, wantCode: http.StatusNoContent},
	}

	r := testEngine(allowed)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := map[string]string{}
			if tt.origin != "" {
				headers["Origin"] = tt.origin
			}

			w := do(r, tt.method, "/healthz", headers)

			if w.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", w.Code, tt.wantCode)
			}

			got := w.Header().Get("Access-Control-Allow-Origin")
			switch {
			case tt.wantAllowed && got != tt.origin:
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, tt.origin)
			case !tt.wantAllowed && got != "":
				t.Errorf("Access-Control-Allow-Origin = %q, want it unset", got)
			}
		})
	}
}

// An empty origin header must never match an empty entry in the list.
func TestCORSIgnoresRequestsWithoutOrigin(t *testing.T) {
	r := testEngine()

	w := do(r, http.MethodGet, "/healthz", nil)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("Access-Control-Allow-Origin = %q, want it unset", got)
	}
}
