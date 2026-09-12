// Package server builds the HTTP server: the middleware every request passes
// through, and the routes.
package server

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"xermess/internal/config"
)

// New builds the server: middleware first, in the order every request passes
// through them, then the routes.
//
// gin.New starts with no middleware, unlike gin.Default, which adds Gin's own
// logger. We want the slog one instead, so the whole server logs the same way.
func New(cfg config.Config, log *slog.Logger) *gin.Engine {
	r := gin.New()

	r.Use(
		RequestLogger(log), // one log line per request
		gin.Recovery(),     // a panic becomes a 500 instead of a dead connection
		CORS(cfg.CORSOrigins),
	)

	registerRoutes(r)

	return r
}

// registerRoutes mounts every route. This is the whole API surface: to add an
// endpoint, add it here and write its handler below.
//
// Handlers that need the database will take it as a parameter, and New will
// pass it through to here.
func registerRoutes(r *gin.Engine) {
	r.GET("/healthz", health)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/hello", hello)
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
