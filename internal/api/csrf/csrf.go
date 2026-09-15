// Package csrf stops another site from changing anything with a signed-in
// user's or administrator's cookie.
//
// The session cookies are SameSite=Lax, which keeps other sites' requests from
// carrying them — but "site" means the registrable domain, so any subdomain of
// it counts as the same site: a blog on blog.mywebsite.com, a user's page on
// a shared host. And the JSON endpoints bind any body as JSON, so a plain
// HTML form posting text/plain would reach them without a CORS preflight.
//
// Two checks close that, on every request that can change something:
//
//   - a body has to be application/json, which a form cannot send;
//   - a browser has to say it came from an origin that is allowed: the Origin
//     header, or failing that Sec-Fetch-Site. A request with neither is not a
//     browser's — a script, curl, another server — and carries no victim's
//     cookie, so it is let through to authenticate like any other.
package csrf

import (
	"mime"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

// New guards the routes it is mounted on. `origins` are the exact origins
// allowed to make changes: the app served on them, and any CORS origin
// configured.
func New(origins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		}

		if hasBody(c.Request) {
			media, _, _ := mime.ParseMediaType(c.GetHeader("Content-Type"))
			if media != "application/json" {
				c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"error": "the body must be application/json"})
				return
			}
		}

		if origin := c.GetHeader("Origin"); origin != "" {
			if !slices.Contains(origins, origin) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "requests from " + origin + " are not allowed"})
				return
			}
		} else if site := c.GetHeader("Sec-Fetch-Site"); site != "" && site != "same-origin" && site != "none" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "cross-site requests are not allowed"})
			return
		}

		c.Next()
	}
}

func hasBody(r *http.Request) bool {
	return r.ContentLength > 0 || len(r.TransferEncoding) > 0
}
