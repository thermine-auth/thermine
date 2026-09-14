package origin

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIssuer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	plain, _ := gin.CreateTestContext(httptest.NewRecorder())
	plain.Request = httptest.NewRequest("GET", "http://auth.example.com/x", nil)

	if got := Issuer(plain); got != "http://auth.example.com" {
		t.Errorf("Issuer() = %q", got)
	}

	proxied, _ := gin.CreateTestContext(httptest.NewRecorder())
	proxied.Request = httptest.NewRequest("GET", "http://auth.example.com/x", nil)
	proxied.Request.Header.Set("X-Forwarded-Proto", "https")

	if got := JWKSURI(proxied); got != "https://auth.example.com/.well-known/jwks.json" {
		t.Errorf("JWKSURI() = %q, want https behind a TLS-terminating proxy", got)
	}
}
