// Package server builds the HTTP server: the middleware every request passes
// through, and the routes.
package server

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"xermess/internal/auth"
	"xermess/internal/config"
)

// New builds the server: middleware first, in the order every request passes
// through them, then the routes.
//
// gin.New starts with no middleware, unlike gin.Default, which adds Gin's own
// logger. We want the slog one instead, so the whole server logs the same way.
func New(cfg config.Config, db *gorm.DB, log *slog.Logger) *gin.Engine {
	r := gin.New()

	r.Use(
		RequestLogger(log), // one log line per request
		gin.Recovery(),     // a panic becomes a 500 instead of a dead connection
		CORS(cfg.CORSOrigins),
	)

	registerRoutes(r, &admin{
		auth: auth.New(db),
		db:   db,
		log:  log,
		// The cookie is only sent over HTTPS when the browser reaches the API
		// over HTTPS; over plain http in development it has to stay off.
		secure: false,
	})

	return r
}

// registerRoutes mounts every route. This is the whole API surface.
func registerRoutes(r *gin.Engine, h *admin) {
	r.GET("/healthz", health)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/hello", hello)

		// Signing in is the one admin route that cannot require a session.
		v1.POST("/admin/auth/login", h.login)

		signedIn := v1.Group("/admin", RequireAdmin(h.auth))
		{
			signedIn.POST("/auth/logout", h.logout)
			signedIn.GET("/me", h.me)
			signedIn.GET("/overview", h.overview)
			signedIn.GET("/sessions", h.sessions)
		}
	}

	r.NoRoute(notFound)
}

// health says the server is up.
func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// hello is a placeholder, and an example of what a handler looks like.
func hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"hello": "world"})
}

// notFound answers any path that no route matched.
func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}
