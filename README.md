# xermess

An authentication server written in Go, with a SvelteKit web frontend.

> Early days — this is the project scaffold. The auth flows are not implemented yet.

## Requirements

- Go 1.27+
- PostgreSQL 14+
- [Bun](https://bun.sh) — for the `web/` frontend

## Getting started

```sh
make setup   # creates .env from .env.example
make run
```

The server listens on `:8080` by default. Check it is up:

```sh
curl localhost:8080/healthz
# {"status":"ok"}
```

To run the frontend:

```sh
make web-install
make web-dev
```

## Make targets

Run `make` on its own for the full list, grouped by what it is for:

| Group       | Targets                                                              |
| ----------- | -------------------------------------------------------------------- |
| Development | `run`, `build`, `clean`                                              |
| Quality     | `check` (fmt + vet + test), `fmt`, `vet`, `test`, `tidy`             |
| Database    | `migrate-diff name=x`, `migrate-apply`, `migrate-status`, `migrate-hash`, `migrate-lint` |
| Frontend    | `web-install`, `web-dev`, `web-build`, `web-start`                   |
| Setup       | `setup` (create `.env`), `tools`                                     |

## Layout

```
cmd/xermess/main.go            startup, in order, in one function
cmd/migrate/main.go            runs migrations by hand: up, down, status

internal/config/config.go      reads .env
internal/database/database.go  opens the connection
internal/database/migrate.go   applies migrations
internal/auth/auth.go          signs administrators in and out, records what they do
internal/server/server.go      the server: middleware, routes, handlers
internal/server/admin.go       the admin panel's endpoints
internal/server/middleware.go  request logging, CORS, the session check
internal/model/                one file per table, listed in model.All

migrations/                    one Go file per migration, applied in order

web/src/lib/api/               the typed client for this API
web/src/lib/components/ui/     the building blocks: Button, Card, TextField, Icon…
web/src/lib/components/admin/  the panel's own pieces: header, tables, stats
web/src/lib/styles/            fonts.css, tokens.css, base.css, ark.css
web/src/lib/theme.svelte.ts    the light/dark/auto choice
web/src/routes/admin/login/    the sign-in page
web/src/routes/admin/(panel)/  dashboard, logs, settings — everything behind a session
web/src/lib/demo.ts            placeholder rows for the sections with no backend
```

Four packages, each with one job: `config` reads settings, `database` talks to
Postgres, `model` describes the tables, `server` answers requests. Nothing
imports `server` except `main`, and `model` imports nothing of ours at all.

Startup is `run` in `cmd/xermess/main.go`, top to bottom: read the
configuration, open the database, apply migrations, serve.

The server is Gin with its defaults: `server.New` builds the engine and
`engine.Run(addr)` listens. Ctrl-C stops the process immediately, so a request
being handled at that moment is cut off — fine in development, worth revisiting
before this runs for real.

## API

| Method | Path                          | Needs a session | Description                    |
| ------ | ----------------------------- | --------------- | ------------------------------ |
| `GET`  | `/healthz`                    | no              | The server is up               |
| `GET`  | `/api/v1/hello`               | no              | Placeholder                    |
| `POST` | `/api/v1/admin/auth/login`    | no              | Sign in, sets the session cookie |
| `POST` | `/api/v1/admin/auth/logout`   | yes             | Sign out, revokes the session  |
| `GET`  | `/api/v1/admin/me`            | yes             | The signed-in administrator    |
| `GET`  | `/api/v1/admin/overview`      | yes             | Counts and recent activity     |
| `GET`  | `/api/v1/admin/logs`          | yes             | The activity log (`?limit=`)   |
| `GET`  | `/api/v1/admin/sessions`      | yes             | The caller's own sessions      |

Sessions are a random token in an HttpOnly cookie; the database keeps only a
SHA-256 hash of it, so a leaked database cannot be signed in with. They last
12 hours. Signing in, failing to sign in, and signing out are all written to
`audit_logs`.

To add an endpoint: mount it in `registerRoutes` and write its handler, both in
`internal/server/server.go`, or in `internal/server/admin.go` for the admin
panel. Put it behind `RequireAdmin` unless it is meant to be public.

## Admin panel

The SvelteKit app in `web/` is the admin panel.

```sh
make run          # the API, on :8080
make web-dev      # the panel, on :5173
```

`make web-start` builds the panel and serves that build on :4173, which is how
to check a production build locally. That is a different origin from the dev
server, so :4173 has to be in `XERMESS_CORS_ORIGINS` too — it is in
`.env.example`. An origin missing from that list has its responses discarded
by the browser, which the panel can only report as not being able to reach
the server. It is Vite's preview server, not a
production one: `adapter-auto` finds no known platform here, so nothing
deployable is produced. Choose an adapter — `adapter-node` for running it
yourself — when it is time to deploy.

Then open http://localhost:5173/admin/login and sign in with
`XERMESS_ADMIN_USERNAME` and `XERMESS_ADMIN_PASSWORD`. Signing in leads to
`/admin/dashboard`.

The header carries the three top-level sections — Dashboard, Logs, Settings —
and the dashboard has its own sidebar for what sits under it. Its Overview
page is real; Administrators, Roles, API keys and Webhooks render the
placeholder rows in `lib/demo.ts` and are marked as such in the interface.
Delete a block from that file as soon as its section talks to the API.

`web/.env` names the API in `PUBLIC_API_URL`, used both by the server when it
renders a page and by the browser for signing in and out — which is why
`XERMESS_CORS_ORIGINS` has to list the panel's own address.

Server-side rendering reads the session cookie the API set. That works
locally because cookies ignore port numbers, so a cookie set by `:8080` is
sent to `:5173`. Across two real domains it would not be: the API would have
to set the cookie on a domain that covers both.

**Rendering.** Pages are rendered on the server, which is what makes a reload
show the finished page rather than assembling one: the theme, the title and
the content are all in the first response. That means the data has to be
fetched on the server too, so the loads are `+page.server.ts` and
`+layout.server.ts`, and `lib/server/api.ts` forwards the session cookie to
the API — a fetch made by the server carries none of the browser's cookies on
its own.

Signing in and out still happen in the browser, followed by `invalidateAll()`
so the server loads run again with the new session.

**Components.** `lib/components/ui` holds the building blocks and
`lib/components/admin` the pieces only this panel uses. Pages compose those
and never reach for an Ark UI primitive directly, so a change to how a field
looks happens in one file.

**Layout.** Pages sit in a centred column capped at `--content-width`
(72rem); the dashboard takes the full width for its sidebar and centres its
own column beside it. Below 55rem the sidebar becomes a scrollable row above
the content, and below 30rem the header sections become their icons alone —
their labels stay in the accessible name.

**Styling.** [Ark UI](https://ark-ui.com) ships no CSS: every part it renders
carries `data-scope` and `data-part`, and `lib/styles/ark.css` styles those
attributes. Everything else refers to the tokens in `tokens.css`. There is no
CSS framework.

The palette is taken from the [PocketBase](https://pocketbase.io) admin UI —
its near-black primary, soft grey secondary, navy header and filled inputs —
with the dark theme built the way PocketBase builds it: one base colour mixed
with increasing amounts of white, so the greys stay in step.

Fields follow their settings page exactly: one filled block with the label
inside it at the top and the value below, no border anywhere, and focus shown
by the block darkening (`#e4e8ec` to `#dce0e5`) while the label goes from
`#687278` to the full text colour. The measurements — a 24px label row over a
38.5px value row, 5px radius, 13px bold label, 13px side padding — are in
`lib/styles/ark.css`.

**Theme.** Light or dark, switched by the toggle in the header — one click,
no menu — and remembered in a cookie.

The cookie rather than `localStorage` is the point: `hooks.server.ts` reads it
and writes `data-theme` straight onto the `<html>` tag, so the page arrives
already dark or light. Nothing has to be corrected after the fact, which is
what a flash on refresh actually is. A reader who has not chosen yet gets no
attribute, leaving the `prefers-color-scheme` rules in `tokens.css` to decide
— also before the first paint, because that happens in CSS rather than after
it.

The toggle's icon is chosen in CSS from that same attribute rather than in
JavaScript, so the server renders the right one and the browser has nothing
to correct when it hydrates.

Switching is animated as one cross-fade of the whole page through the View
Transitions API, rather than per-element transitions that each start at a
slightly different moment. Browsers without it simply change. Both that and
the toggle's own icon animation stop at `prefers-reduced-motion`.

**Fonts.** Product Sans for text and Consolas for code, both loaded with
`local()` only. Neither can be bundled — Product Sans is Google's corporate
typeface and is not licensed for redistribution, and Consolas ships with
Windows and Office — so a machine that has them uses them and one that does
not falls back quietly. To self-host licensed copies, put the files in
`static/fonts` and add a `url(...)` source in `lib/styles/fonts.css`.

**Icons** are [Remix Icon](https://remixicon.com), through `svelte-remixicon`.
They are components, so only the ones actually used are bundled — there is no
icon font to download. Pass one to the `Icon` wrapper rather than using it
directly, which keeps sizing and alignment in one place:

```svelte
<Icon icon={RiLogoutBoxRLine} />
```

Give `Icon` a `label` only when the icon carries meaning on its own; beside
text it stays `aria-hidden` so a screen reader does not read the same thing
twice.

## Database

Migrations are Go files in `migrations/`, run by
[goose](https://github.com/pressly/goose). The server applies pending ones on
start unless `XERMESS_DB_MIGRATE=false`.

```sh
make migrate-new name=add_admin_phone   # create an empty migration
make migrate-up                         # apply pending
make migrate-status                     # what is applied
make migrate-down                       # roll the newest one back
```

A migration registers itself with goose and gets the transaction goose opened.
`gormTx` wraps that transaction in GORM, so a migration can work with the model
structs instead of writing DDL by hand:

```go
func init() {
	goose.AddMigrationContext(upAddAdminPhone, downAddAdminPhone)
}

func upAddAdminPhone(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}
	return db.Migrator().AddColumn(&model.AdminUser{}, "Phone")
}

func downAddAdminPhone(_ context.Context, tx *sql.Tx) error {
	db, err := gormTx(tx)
	if err != nil {
		return err
	}
	return db.Migrator().DropColumn(&model.AdminUser{}, "Phone")
}
```

Everything a migration does runs inside goose's transaction, so one that fails
half way leaves nothing behind. Write the down function even when you think you
will not need it: it is what makes a bad deploy recoverable.

Because the migrations are Go, the `goose` command-line tool cannot run them —
only a binary that imports them can. That is what `cmd/migrate` is for, and
what the make targets above use.

Never edit a migration that has already run anywhere. Add a new one.

### The first administrator

`20260912113137_seed_admin.go` creates the four roles and one administrator
holding `super_admin`, from `XERMESS_ADMIN_USERNAME` and
`XERMESS_ADMIN_PASSWORD`. The password is stored as a bcrypt hash, never as
itself. With either setting empty the roles are still created and no
administrator is, so a test database does not need credentials.

Rolling this migration back deletes that administrator and the roles.

## Configuration

Every setting is an environment variable, read from `.env` first; real
environment variables win. `.env.example` lists all of them. `XERMESS_DB_DSN`
has no default, so a missing one stops the server.

The `migrations/` directory has to ship with the binary: goose reads the file
names from disk, at the path in `XERMESS_DB_MIGRATE_DIR`, and matches them to
the functions compiled in.

## License

[MIT](LICENSE)
