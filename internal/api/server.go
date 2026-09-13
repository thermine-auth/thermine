// Package api builds the HTTP server: the middleware every request passes
// through, and the table of which handler answers which path.
//
// The handlers themselves live one directory down, one package per subject:
// http/auth signs administrators in, http/users manages user records,
// http/fields the columns those records are made of, http/activity reports on
// what has happened. Each of those packages is the same four files — the
// handler, the requests it accepts, the answers it gives, and the rules it
// holds them to — so finding your way around a new one is the same as finding
// your way around the last.
package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"xermess/internal/api/http/activity"
	"xermess/internal/api/http/audit"
	apiauth "xermess/internal/api/http/auth"
	"xermess/internal/api/http/cors"
	"xermess/internal/api/http/fields"
	"xermess/internal/api/http/middleware"
	"xermess/internal/api/http/session"
	"xermess/internal/api/http/setup"
	"xermess/internal/api/http/users"
	"xermess/internal/auth"
	"xermess/internal/config"
	"xermess/internal/store"
)

// secureCookies sets the session cookie's Secure flag. The browser only sends
// such a cookie over HTTPS, so it has to stay off while the API is served
// over plain http in development.
const secureCookies = false

// New builds the server: middleware first, in the order every request passes
// through them, then the routes.
//
// gin.New starts with no middleware, unlike gin.Default, which adds Gin's own
// logger. We want the slog one instead, so the whole server logs the same way.
func New(cfg config.Config, st *store.Store, log *slog.Logger) *gin.Engine {
	r := gin.New()

	// Which origins may call the API is configuration, so the handler is
	// built here and handed to the chain that puts it in order.
	r.Use(middleware.Chain(log, cors.New(cfg.CORSOrigins))...)

	service := auth.New(st)
	recorder := audit.New(st, log)

	registerRoutes(r, service, handlers{
		auth:     apiauth.New(service, st, log, secureCookies),
		setup:    setup.New(st, log),
		users:    users.New(st, recorder, log),
		fields:   fields.New(st, recorder, log),
		activity: activity.New(st, log),
	})

	return r
}

// handlers is one of each, so the route table below reads as a table.
type handlers struct {
	auth     *apiauth.Handler
	setup    *setup.Handler
	users    *users.Handler
	fields   *fields.Handler
	activity *activity.Handler
}

// registerRoutes mounts every route. This is the whole API surface: what a
// path does is in the handler, but that a path exists is only ever here.
func registerRoutes(r *gin.Engine, service *auth.Service, h handlers) {
	r.GET("/healthz", health)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/hello", hello)

		// Setting the panel up and signing in are the routes that cannot
		// require a session: before the first there is no account, and before
		// the second no way to prove one. Creating an administrator is
		// refused as soon as there is one, which is what keeps the first of
		// those from being a way in.
		v1.GET("/admin/setup", h.setup.Status)
		v1.POST("/admin/setup", h.setup.Create)
		v1.POST("/admin/auth/login", h.auth.Login)

		signedIn := v1.Group("/admin", session.Require(service))
		{
			signedIn.POST("/auth/logout", h.auth.Logout)
			signedIn.GET("/me", h.auth.Me)
			signedIn.GET("/sessions", h.auth.Sessions)

			signedIn.GET("/overview", h.activity.Overview)
			signedIn.GET("/logs", h.activity.Logs)

			// The users an organisation manages, and the fields their
			// records are made of.
			signedIn.GET("/users", h.users.List)
			signedIn.POST("/users", h.users.Create)
			signedIn.GET("/users/:id", h.users.Get)
			signedIn.PATCH("/users/:id", h.users.Update)
			signedIn.DELETE("/users/:id", h.users.Delete)

			signedIn.GET("/user-fields", h.fields.List)
			signedIn.POST("/user-fields", h.fields.Create)
			signedIn.PATCH("/user-fields/:id", h.fields.Update)
			signedIn.DELETE("/user-fields/:id", h.fields.Delete)
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
