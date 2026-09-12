package server

import (
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
)

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
