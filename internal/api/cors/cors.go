// Package cors decides which browser origins may call the API.
//
// It is its own package rather than one more middleware: which origins are
// allowed is configuration, it is the one piece of the request path that can
// hand another site the ability to act as a signed-in administrator, and a
// reader looking for that answer should find one file with one job. The
// handler it returns is passed to middleware.Chain, which puts it in order
// with the rest.
package cors

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

// allowedMethods and allowedHeaders are what a browser is told it may send.
// They cover the whole API: every endpoint is one of these methods, and sends
// either JSON or nothing.
const (
	allowedMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	allowedHeaders = "Authorization, Content-Type"
)

// New lets the listed browser origins call the API, and answers the browser's
// preflight request itself.
//
// Origins are compared exactly — no wildcard, no matching on suffix. The
// panel sends its session cookie, and a browser only sends credentials to an
// origin the server named, so being loose here would be giving any site that
// asked the ability to act as whoever is signed in.
//
// The provider endpoints a browser app calls directly — the token, userinfo,
// revocation and introspection endpoints, discovery and the keys — are open
// to every origin instead, without credentials. They never read a cookie: a
// caller proves itself with a client secret, a PKCE verifier or a bearer
// token it already holds, so which site the script came from grants nothing.
// A single-page app on any domain has to be able to reach them.
func New(origins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if public(c.Request.URL.Path) {
			if origin != "" {
				h := c.Writer.Header()
				h.Set("Access-Control-Allow-Origin", "*")
				h.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				h.Set("Access-Control-Allow-Headers", allowedHeaders)
			}

			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}

			c.Next()
			return
		}

		// A request with no Origin header is not a cross-origin one: curl,
		// another server, the panel's own server-side rendering. It gets no
		// headers rather than matching an empty entry in the list.
		if origin != "" && slices.Contains(origins, origin) {
			h := c.Writer.Header()
			h.Set("Access-Control-Allow-Origin", origin)
			h.Set("Access-Control-Allow-Credentials", "true")
			h.Set("Access-Control-Allow-Methods", allowedMethods)
			h.Set("Access-Control-Allow-Headers", allowedHeaders)

			// The answer depends on which origin asked, so a cache must not
			// hand one origin's answer to another.
			h.Add("Vary", "Origin")
		}

		// A preflight is answered here and never reaches a route. It is
		// answered whoever asked: a browser that was not given the headers
		// above stops at its own check, which is where that belongs.
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// publicPaths are the provider endpoints any origin may call. The authorization
// and logout endpoints are not among them: a browser is sent to those, it does
// not call them from script, and they do read the session cookie.
var publicPaths = []string{
	"/oauth2/token",
	"/oauth2/userinfo",
	"/oauth2/revoke",
	"/oauth2/introspect",
}

func public(path string) bool {
	return slices.Contains(publicPaths, path) || strings.HasPrefix(path, "/.well-known/")
}
