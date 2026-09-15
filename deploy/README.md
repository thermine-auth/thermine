# Deploying xermess

```
                    internet                          staff network / VPN only
                        │                                        │
             https://id.mywebsite.com              https://admin-id.mywebsite.com
                        │                                        │
                ┌───────┴────────────────── Caddy ───────────────┴───────┐
                │  /oauth2/* /.well-known/*       /api/v1/admin/*        │
                │  /api/v1/account/*    else      │            else      │
                └────┬──────────────────┬─────────┼─────────────┬────────┘
                     │                  │         │             │
               api :8080            id          api :8081       console
             (public listener)    :3000     (admin listener)    :3000
                     └────────── Postgres ────────┘
```

Two sites, two apps, one API process with two listeners:

| Site | Serves | API listener | Who can reach it |
| --- | --- | --- | --- |
| `id.mywebsite.com` | the id app, the OAuth/OIDC provider, the account API | public, `:8080` | everyone |
| `admin-id.mywebsite.com` | the console, the admin API | admin, `:8081` | staff only |

## Files

| File | What |
| --- | --- |
| `compose.yaml` | Postgres, `api`, `id`, `console` and Caddy on one internal network |
| `Caddyfile` | the two sites and which paths go to which listener |
| `docker/api.Dockerfile` | the API; it applies migrations when it starts |
| `docker/web.Dockerfile` | any app under `web/`, built with its directory as the context |
| `.env.example` | the settings compose reads from `deploy/.env` |

`make deploy-logs` follows the API; `make deploy-down` stops the stack.

## Why it is shaped like this

- **The admin API is not on the public listener at all.** `/api/v1/admin/*`
  answers only on `:8081`, which only the staff site routes to. Hiding the
  panel's page is not enough when its API is on the internet; this way there
  is nothing to find.
- **Each app is on the same origin as the API it calls.** The session cookie
  then belongs to that one host: the browser sends it without CORS, the app's
  server-side renderer can read it, and no other subdomain of mywebsite.com
  ever receives it.
- **Staff access is enforced before the panel.** The Caddyfile refuses anyone
  outside the staff networks. Stronger: no public DNS for the admin host, and
  reach it over WireGuard/Tailscale or an identity-aware proxy.

## Run it

```sh
cp deploy/.env.example deploy/.env      # fill in URLs, secret key, database password, SMTP
$EDITOR deploy/Caddyfile                # hostnames, email, staff networks
make deploy-up                          # docker compose -f deploy/compose.yaml up -d --build
```

Then open `https://admin-id.mywebsite.com/admin/login` from the staff network
to create the first administrator.

## Checklist before going live

- [ ] `PUBLIC_URL` and `ADMIN_URL` are `https` — session cookies are Secure because of it.
- [ ] Use a domain you own for the admin host (not `.local`, which is reserved
      for mDNS and cannot get a trusted certificate).
- [ ] The staff networks in the Caddyfile are right, or the admin host is VPN-only.
- [ ] `XERMESS_SECRET_KEY` is stored in your secret manager. Losing it signs everyone out.
- [ ] SMTP is configured, or reset emails only reach the log.
- [ ] Postgres is backed up.
- [ ] Nothing but Caddy publishes a port (`docker compose ps`).
- [ ] Register your applications' redirect URIs as `https`.
- [ ] `XERMESS_ADMIN_MFA` is `required` (the default), and every administrator has
      saved their recovery codes.
- [ ] Someone knows how to reset an administrator's two-factor sign-in (a super
      admin, from Administrators) — and that at least two super admins exist.

## Not in this directory

- Two-factor sign-in for administrators is on by default. It does not replace the
  network restriction; each covers what the other misses.
- Signing keys rotate on their own every 90 days. If one may have leaked, a super
  admin rotates at once with revocation: `POST /api/v1/admin/signing-keys/rotate`
  with `{"revoke_old": true}`.
- Rate limiting in xermess is per process. With several API replicas, add a
  shared limit at the edge too.
