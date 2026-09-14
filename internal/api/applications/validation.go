package applications

import (
	"fmt"
	"net/http"

	"xermess/internal/api/respond"
	"xermess/internal/api/validate"
	"xermess/internal/model"
)

// applyTo checks the request and copies it onto an application. `creating`
// says whether the type was taken from the request; on an update the stored
// type stands.
//
// The shape of each field is checked here; what the fields mean next to each
// other — a public client with a secret, a redirect URI that is not exact —
// is the model's to say.
func (r *applicationRequest) applyTo(app *model.Application, creating bool) error {
	r.clean()

	if err := validate.Struct(r); err != nil {
		return err
	}

	if len(r.RedirectURIs) > maxURIs || len(r.PostLogoutRedirectURIs) > maxURIs {
		return badRequest(fmt.Sprintf("an application may have at most %d redirect URIs of each kind", maxURIs))
	}

	if creating && !model.ApplicationType(r.Type).Valid() {
		return badRequest("type must be one of: web, spa, native, m2m")
	}

	app.Name = r.Name
	app.Description = r.Description
	app.LogoURI = r.LogoURI
	app.ClientURI = r.ClientURI
	app.TokenEndpointAuthMethod = model.AuthMethod(r.TokenEndpointAuthMethod)
	app.GrantTypes = r.GrantTypes
	app.RedirectURIs = r.RedirectURIs
	app.PostLogoutRedirectURIs = r.PostLogoutRedirectURIs
	app.Scopes = r.Scopes
	app.RequirePKCE = validate.Flag(r.RequirePKCE, app.RequirePKCE)
	app.AccessTokenLifetime = r.AccessTokenLifetime
	app.IDTokenLifetime = r.IDTokenLifetime
	app.RefreshTokenLifetime = r.RefreshTokenLifetime
	app.AssertRoles = validate.Flag(r.AssertRoles, app.AssertRoles)
	app.RequireRoleAssignment = validate.Flag(r.RequireRoleAssignment, app.RequireRoleAssignment)
	app.Enabled = validate.Flag(r.Enabled, app.Enabled)

	app.Normalise()

	// A client moving from one secret method to the other keeps its secret;
	// one that has no secret method any more loses it.
	if !app.HasSecret() {
		app.ClientSecretHash = ""
		app.SecretHint = ""
		app.SecretCreatedAt = nil
	}

	if err := app.Validate(); err != nil {
		return badRequest(err.Error())
	}

	return nil
}

func badRequest(message string) error {
	return respond.Fault{Status: http.StatusBadRequest, Message: message}
}
