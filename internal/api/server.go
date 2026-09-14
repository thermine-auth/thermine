// Package api builds the HTTP server: the middleware every request passes
// through, and the table of which handler answers which path.
//
// The handlers themselves live one directory down, one package per subject:
// auth signs administrators in, users manages user records, fields the
// columns those records are made of, applications the OAuth clients that sign
// users in, roles the roles users hold in each application, admins and
// adminroles the administrators and what they may do, activity reports on
// what has happened. Each of those packages is the same four files — the
// handler, the requests it accepts, the answers it gives, and the rules it
// holds them to — so finding your way around a new one is the same as finding
// your way around the last.
package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"xermess/internal/api/activity"
	"xermess/internal/api/adminroles"
	"xermess/internal/api/admins"
	"xermess/internal/api/apis"
	"xermess/internal/api/applications"
	"xermess/internal/api/audit"
	apiauth "xermess/internal/api/auth"
	"xermess/internal/api/cors"
	"xermess/internal/api/fields"
	"xermess/internal/api/middleware"
	"xermess/internal/api/roles"
	"xermess/internal/api/session"
	"xermess/internal/api/setup"
	"xermess/internal/api/users"
	"xermess/internal/auth"
	"xermess/internal/config"
	"xermess/internal/model"
	"xermess/internal/store"
)

// New builds the server: middleware first, in the order every request passes
// through them, then the routes.
//
// gin.New starts with no middleware, unlike gin.Default, which adds Gin's own
// logger. We want the slog one instead, so the whole server logs the same way.
func New(cfg config.Config, st *store.Store, log *slog.Logger) (*gin.Engine, error) {
	r := gin.New()

	// Gin believes every X-Forwarded-For header unless told otherwise, which
	// would let any caller choose the address the log records for it.
	if err := r.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		return nil, fmt.Errorf("trusted proxies: %w", err)
	}

	// Which origins may call the API is configuration, so the handler is
	// built here and handed to the chain that puts it in order.
	r.Use(middleware.Chain(log, cors.New(cfg.CORSOrigins))...)

	service := auth.New(st)
	recorder := audit.New(st, log)

	registerRoutes(r, service, handlers{
		auth:         apiauth.New(service, st, log, cfg.SecureCookies),
		setup:        setup.New(st, log),
		users:        users.New(st, recorder, log),
		fields:       fields.New(st, recorder, log),
		roles:        roles.New(st, recorder, log),
		admins:       admins.New(st, recorder, log),
		adminRoles:   adminroles.New(st, recorder, log),
		applications: applications.New(st, recorder, log),
		apis:         apis.New(st, recorder, log),
		activity:     activity.New(st, log),
	})

	return r, nil
}

// handlers is one of each, so the route table below reads as a table.
type handlers struct {
	auth         *apiauth.Handler
	setup        *setup.Handler
	users        *users.Handler
	fields       *fields.Handler
	roles        *roles.Handler
	admins       *admins.Handler
	adminRoles   *adminroles.Handler
	applications *applications.Handler
	apis         *apis.Handler
	activity     *activity.Handler
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
			// Everything about the caller's own account is open to every
			// administrator who can sign in.
			signedIn.POST("/auth/logout", h.auth.Logout)
			signedIn.GET("/me", h.auth.Me)
			signedIn.GET("/sessions", h.auth.Sessions)

			// The rest is guarded by what the administrator's roles allow.
			// The panel hides what someone cannot do; these checks are what
			// actually stops them.
			activity := signedIn.Group("", session.Can(model.PermActivityRead))
			activity.GET("/overview", h.activity.Overview)
			activity.GET("/logs", h.activity.Logs)

			// The users an organisation manages and the fields their records
			// are made of. Users are shared by every application, so these
			// are whole-panel permissions.
			readUsers := signedIn.Group("", session.Can(model.PermUsersRead))
			readUsers.GET("/users", h.users.List)
			readUsers.GET("/users/:id", h.users.Get)
			readUsers.GET("/users/:id/roles", h.users.Roles)
			readUsers.GET("/users/:id/role-mappings", h.users.RoleMappings)
			readUsers.GET("/user-fields", h.fields.List)

			writeUsers := signedIn.Group("", session.Can(model.PermUsersWrite))
			writeUsers.POST("/users", h.users.Create)
			writeUsers.PATCH("/users/:id", h.users.Update)
			writeUsers.DELETE("/users/:id", h.users.Delete)

			writeFields := signedIn.Group("", session.Can(model.PermUserFieldsWrite))
			writeFields.POST("/user-fields", h.fields.Create)
			writeFields.PATCH("/user-fields/:id", h.fields.Update)
			writeFields.DELETE("/user-fields/:id", h.fields.Delete)

			// Applications, the roles each defines, and who holds them. A
			// role can grant these for one application, so the routes only
			// check the administrator can reach some application; each
			// handler then checks the one the request is about.
			apps := signedIn.Group("", session.CanAnywhere(model.PermApplicationsRead))
			apps.GET("/applications", h.applications.List)
			apps.GET("/applications/:id", h.applications.Get)
			apps.PATCH("/applications/:id", h.applications.Update)
			apps.POST("/applications/:id/secret", h.applications.RotateSecret)

			// Roles: global ones, which belong to no application, and each
			// application's. Anyone who can see users or some application can
			// read them; each handler narrows to the scopes the administrator
			// reaches and checks writes against the role's scope.
			roles := signedIn.Group("", session.CanAnywhere(model.PermUsersRead, model.PermApplicationsRead))
			roles.GET("/user-roles", h.roles.List)
			roles.GET("/user-roles/:id", h.roles.Get)

			writeRoles := signedIn.Group("", session.CanAnywhere(model.PermUserRolesWrite))
			writeRoles.POST("/user-roles", h.roles.Create)
			writeRoles.PATCH("/user-roles/:id", h.roles.Update)
			writeRoles.DELETE("/user-roles/:id", h.roles.Delete)

			// Keycloak's role mapping: giving a user roles and taking them
			// away, each role checked against its own scope.
			assign := signedIn.Group("", session.CanAnywhere(model.PermRoleAssignmentsWrite))
			assign.POST("/users/:id/role-mappings", h.users.AssignRoles)
			assign.DELETE("/users/:id/role-mappings/:role", h.users.UnassignRole)

			// What an application may do with each API, and a preview of the
			// tokens it would get. Each handler checks the application.
			apps.GET("/applications/:id/apis", h.applications.APIAccess)
			apps.PUT("/applications/:id/apis/:api", h.applications.AuthorizeAPI)
			apps.DELETE("/applications/:id/apis/:api", h.applications.RevokeAPI)
			apps.POST("/applications/:id/token-preview", h.applications.TokenPreview)

			// APIs: the resource servers tokens are issued for. Anyone who can
			// see them or configure an application's access to them may read
			// them; changing them is a whole-panel permission.
			readAPIs := signedIn.Group("", session.CanAnywhere(model.PermAPIsRead, model.PermApplicationsWrite))
			readAPIs.GET("/apis", h.apis.List)
			readAPIs.GET("/apis/:id", h.apis.Get)
			readAPIs.GET("/apis/:id/applications", h.apis.Applications)

			// An API's log names every application given or refused access,
			// including ones an administrator of a few applications cannot
			// see, so it takes apis.read itself.
			signedIn.GET("/apis/:id/logs", session.Can(model.PermAPIsRead), h.apis.Logs)

			writeAPIs := signedIn.Group("", session.Can(model.PermAPIsWrite))
			writeAPIs.POST("/apis", h.apis.Create)
			writeAPIs.PATCH("/apis/:id", h.apis.Update)
			writeAPIs.DELETE("/apis/:id", h.apis.Delete)

			// Registering and removing applications changes what exists at
			// all, so it takes applications.write for the whole panel.
			registerApps := signedIn.Group("", session.Can(model.PermApplicationsWrite))
			registerApps.POST("/applications", h.applications.Create)
			registerApps.DELETE("/applications/:id", h.applications.Delete)

			// The administrators themselves, and their roles: a super
			// admin's alone, whatever other roles grant.
			super := signedIn.Group("", session.RequireSuperAdmin())
			super.GET("/admins", h.admins.List)
			super.POST("/admins", h.admins.Create)
			super.GET("/admins/:id", h.admins.Get)
			super.PATCH("/admins/:id", h.admins.Update)
			super.DELETE("/admins/:id", h.admins.Delete)

			super.GET("/admin-permissions", h.adminRoles.Permissions)
			super.GET("/admin-roles", h.adminRoles.List)
			super.POST("/admin-roles", h.adminRoles.Create)
			super.PATCH("/admin-roles/:id", h.adminRoles.Update)
			super.DELETE("/admin-roles/:id", h.adminRoles.Delete)
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
