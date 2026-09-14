// Package origin says where this server is, as the tokens it issues and the
// applications checking them name it.
package origin

import "github.com/gin-gonic/gin"

// Issuer is the iss a token from this server carries: its own origin, as the
// request reached it. The provider endpoints will make this configurable.
func Issuer(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	return scheme + "://" + c.Request.Host
}

// JWKSURI is where this server publishes the public keys its tokens are
// signed with, which is what an API checks a token's signature against.
func JWKSURI(c *gin.Context) string {
	return Issuer(c) + "/.well-known/jwks.json"
}
