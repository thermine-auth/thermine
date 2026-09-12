package server

import (
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"

	"xermess/internal/auth"
	"xermess/internal/model"
)

// adminKey is the Gin context key the signed-in administrator is stored under.
const adminKey = "admin"

// RequestLogger logs one line per request, after it has been handled.
func RequestLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Hand the request to the next middleware and the route handler.
		c.Next()

		log.Info("request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", time.Since(start).String(),
			"ip", c.ClientIP(),
		)
	}
}

// RequireAdmin refuses the request unless it carries a session for an
// administrator who may sign in. Handlers behind it can call adminFrom
// without checking.
func RequireAdmin(service *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(sessionCookie)

		user, err := service.Authenticate(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not signed in"})
			return
		}

		c.Set(adminKey, user)

		c.Next()
	}
}

// adminFrom returns the administrator making the request. It is only valid
// behind RequireAdmin, which is the only thing that sets it.
func adminFrom(c *gin.Context) *model.AdminUser {
	user, _ := c.Get(adminKey)
	admin, _ := user.(*model.AdminUser)
	return admin
}

// CORS lets the listed browser origins call the API, and answers the browser's
// preflight request itself. Origins are compared exactly.
func CORS(origins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if slices.Contains(origins, origin) {
			h := c.Writer.Header()
			h.Set("Access-Control-Allow-Origin", origin)
			h.Set("Access-Control-Allow-Credentials", "true")
			h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			h.Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		}

		// A preflight is answered here and never reaches a route.
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
