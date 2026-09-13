// Package session is everything about who is signed in, over HTTP: the cookie
// the browser carries, the check every admin route passes through, and the
// way a handler asks who is making the request.
package session

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"xermess/internal/auth"
	"xermess/internal/model"
)

// Cookie carries the session token between the browser and the server. It is
// HttpOnly so no script can read it, which is what keeps a cross-site
// scripting bug from turning into a stolen session.
const Cookie = "xermess_session"

// key is the Gin context key the signed-in administrator is stored under.
const key = "admin"

// Set stores the token in the browser. `secure` adds the Secure flag, which
// can only be on once the API is served over HTTPS.
func Set(c *gin.Context, token string, secure bool) {
	write(c, token, int(auth.SessionLifetime.Seconds()), secure)
}

// Clear removes it again.
func Clear(c *gin.Context, secure bool) {
	write(c, "", -1, secure)
}

func write(c *gin.Context, value string, maxAge int, secure bool) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     Cookie,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

// Require refuses the request unless it carries a session for an
// administrator who may sign in. Handlers behind it can call Admin without
// checking.
func Require(service *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(Cookie)

		user, err := service.Authenticate(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not signed in"})
			return
		}

		c.Set(key, user)

		c.Next()
	}
}

// Admin returns the administrator making the request. It is only valid behind
// Require, which is the only thing that sets it.
func Admin(c *gin.Context) *model.AdminUser {
	user, _ := c.Get(key)
	admin, _ := user.(*model.AdminUser)

	return admin
}
