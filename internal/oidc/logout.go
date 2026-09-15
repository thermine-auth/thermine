package oidc

import (
	"context"

	"errors"
	"net/url"
	"slices"

	"xermess/internal/jose"
	"xermess/internal/model"
	"xermess/internal/store"
)

// LogoutParams are the parameters of an RP-initiated logout (OpenID Connect
// RP-Initiated Logout 1.0).
type LogoutParams struct {
	IDTokenHint           string
	ClientID              string
	PostLogoutRedirectURI string
	State                 string
}

// Logout ends the user's session at this server and returns where to send the
// browser: the application's registered post-logout redirect URI when it named
// one, otherwise the sign-in app's signed-out page.
//
// It does not revoke the refresh tokens applications hold: signing out of the
// provider is not signing out of every application, and each application
// ends its own session.
func (s *Service) Logout(ctx context.Context, p LogoutParams, sessionToken string, client Client) (string, error) {
	clientID := p.ClientID

	if p.IDTokenHint != "" {
		// The hint names who asked; it may well have expired, which is fine.
		// Its signature and issuer still have to be this server's.
		var claims struct {
			Issuer   string `json:"iss"`
			Audience any    `json:"aud"`
		}
		if _, err := jose.Verify(p.IDTokenHint, s.keys.lookup(ctx), &claims); err != nil || claims.Issuer != s.issuer {
			return s.errorPage(ErrInvalidRequest, "id_token_hint is not an ID token from this server"), nil
		}

		audience, _ := claims.Audience.(string)
		if clientID != "" && clientID != audience {
			return s.errorPage(ErrInvalidRequest, "client_id does not match the id_token_hint"), nil
		}
		clientID = audience
	}

	var app *model.Application
	if clientID != "" {
		found, err := s.store.ApplicationByClientID(ctx, clientID)
		switch {
		case errors.Is(err, store.ErrNotFound):
			return s.errorPage(ErrInvalidRequest, "no application has this client_id"), nil
		case err != nil:
			return "", err
		}
		app = found
	}

	// Only a URI the application registered is followed, so a logout link
	// cannot be used to send someone anywhere.
	if p.PostLogoutRedirectURI != "" {
		if app == nil {
			return s.errorPage(ErrInvalidRequest, "post_logout_redirect_uri needs id_token_hint or client_id"), nil
		}
		if !slices.Contains(app.PostLogoutRedirectURIs, p.PostLogoutRedirectURI) {
			return s.errorPage(ErrInvalidRequest, "post_logout_redirect_uri is not registered for this application"), nil
		}
	}

	if err := s.SignOut(ctx, sessionToken, client); err != nil {
		return "", err
	}

	if p.PostLogoutRedirectURI != "" {
		return withQuery(p.PostLogoutRedirectURI, url.Values{"state": {p.State}}), nil
	}

	values := url.Values{}
	if app != nil {
		values.Set("client_id", app.ClientID)
	}

	return withQuery(s.accountURL+PageLoggedOut, values), nil
}

// PublicApplication is what anyone may know about an application: what its
// sign-in pages show.
type PublicApplication struct {
	Name              string `json:"name"`
	LogoURI           string `json:"logo_uri"`
	ClientURI         string `json:"client_uri"`
	PolicyURI         string `json:"policy_uri"`
	TosURI            string `json:"tos_uri"`
	AllowRegistration bool   `json:"allow_registration"`
}

// Public returns what an application's sign-in pages show about it.
func Public(app *model.Application) PublicApplication {
	return PublicApplication{
		Name:              app.Name,
		LogoURI:           app.LogoURI,
		ClientURI:         app.ClientURI,
		PolicyURI:         app.PolicyURI,
		TosURI:            app.TosURI,
		AllowRegistration: app.AllowRegistration && app.Enabled,
	}
}

// ApplicationByClientID returns what the signed-out page shows about an
// application, found by its client id.
func (s *Service) ApplicationByClientID(ctx context.Context, clientID string) (*PublicApplication, error) {
	app, err := s.store.ApplicationByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	public := Public(app)
	return &public, nil
}
