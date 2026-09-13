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

The Makefile is the index of what can be done; anything longer than a couple
of lines lives in `scripts/` and is called from a target. `make setup` runs
`scripts/install.sh`, `make db-create` runs `scripts/db.sh create`, and so on
— so there is one implementation of each, and the container, which has no
make, can run the same scripts.

Run `make` on its own for the full list, grouped by what it is for:

| Group           | Targets                                                                |
| --------------- | ---------------------------------------------------------------------- |
| Getting started | `setup`, `env`, `deps`                                                 |
| Development     | `run`, `run-migrate`, `build`, `release`, `clean`                      |
| Database        | `migrate-up`, `migrate-down`, `migrate-status`, `migrate-new name=x`, `db-create`, `db-reset`, `db-psql` |
| Quality         | `check` (fmt + vet + test), `fmt`, `vet`, `test`, `tidy`               |
| Container       | `docker-build`, `docker-run`                                           |
| Frontend        | `web-install`, `web-dev`, `web-build`, `web-start`                     |

## Layout

```
cmd/xermess/main.go            startup, in order, in one function
cmd/migrate/main.go            runs migrations by hand: up, down, status

internal/config/config.go      reads .env
internal/database/database.go  opens the connection
internal/database/migrate.go   applies migrations
internal/store/                every query in the project, one file per subject
internal/auth/auth.go          signs administrators in and out, records what they do
internal/api/server.go         the engine, and the table of every route
internal/api/http/auth/        signing in and out, and who is signed in
internal/api/http/users/       the users an organisation manages
internal/api/http/fields/      the fields a user record is made of
internal/api/http/activity/    the dashboard counts and the log
internal/api/http/middleware/  request logging, recovery, and their order
internal/api/http/cors/        which browser origins may call the API
internal/api/http/session/     the cookie, the session check, the current admin
internal/api/http/respond/     how an error is written, once for every endpoint
internal/api/http/audit/       recording what an administrator did
internal/model/                one file per table, listed in model.All

migrations/                    one Go file per migration, applied in order

scripts/install.sh             set the project up on a fresh machine
scripts/db.sh                  create, reset or open the database
scripts/build.sh               build for release, stamping in the version
scripts/entrypoint.sh          what the container runs: migrate, then serve
scripts/lib.sh                 what those scripts share: .env, the DSN, checks
scripts/Dockerfile             the API image, built from the repository root

web/src/lib/api/               the typed client for this API
web/src/lib/components/ui/     the building blocks: Button, Card, TextField, Icon…
web/src/lib/components/admin/  the panel's own pieces: header, tables, stats
web/src/lib/styles/            fonts.css, tokens.css, base.css, ark.css
web/src/lib/theme.svelte.ts    the light/dark/auto choice
web/src/routes/admin/login/    the sign-in page
web/src/routes/admin/(panel)/  dashboard, logs, profile — everything behind a session
web/src/lib/demo.ts            placeholder rows for the sections with no backend
```

### The API packages

Handlers are grouped by subject, and every group is the same four files, so a
package you have never opened is laid out like the last one you did:

```
internal/api/http/users/handler.go     the endpoints: what happens, in order
internal/api/http/users/request.go     the bodies and query strings it accepts
internal/api/http/users/response.go    the shapes it answers with
internal/api/http/users/validation.go  the rules a request has to keep
```

A handler reads a request, asks the store, and answers. Nothing else: the
rules live in `validation.go` and return a `respond.Fault` carrying the status
to answer with, and `respond.Failure` turns that into the answer — or logs
anything that is not a Fault and says only that something went wrong.

### The store

`internal/store` is the only package that writes queries. A handler asks it
for what it needs — `store.Users`, `store.UserField`, `store.WriteAudit` — and
gets models back, so the handlers stay about HTTP and the queries stay in one
place to read and change. It has its own errors, `store.ErrNotFound` and
`store.ErrDuplicate`, which is why nothing above it imports GORM.

```
internal/store/store.go        the Store type, and the errors it returns
internal/store/users.go        listing, searching and writing users
internal/store/user_fields.go  the field definitions
internal/store/admins.go       administrators, for signing in
internal/store/sessions.go     sessions: start, find, revoke, list
internal/store/audit.go        the activity log, and the dashboard counts
```

Each package has one job: `config` reads settings, `database` opens the
connection, `model` describes the tables, `store` queries them, `auth` decides
who may sign in, `api` answers requests. Nothing imports `api` except `main`,
and `model` imports nothing of ours at all.

Startup is `run` in `cmd/xermess/main.go`, top to bottom: read the
configuration, open the database, apply migrations, build the store, serve.

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
| `GET`  | `/api/v1/admin/users`         | yes             | List users (`?search=&verified=&limit=&offset=`) |
| `POST` | `/api/v1/admin/users`         | yes             | Create a user                  |
| `GET`  | `/api/v1/admin/users/:id`     | yes             | One user                       |
| `PATCH`| `/api/v1/admin/users/:id`     | yes             | Update a user                  |
| `DELETE`| `/api/v1/admin/users/:id`    | yes             | Delete a user                  |
| `GET`  | `/api/v1/admin/user-fields`   | yes             | The fields a user record has   |
| `POST` | `/api/v1/admin/user-fields`   | yes             | Add a field                    |
| `DELETE`| `/api/v1/admin/user-fields/:id` | yes          | Remove a field                 |
| `GET`  | `/api/v1/admin/sessions`      | yes             | The caller's own sessions      |

Sessions are a random token in an HttpOnly cookie; the database keeps only a
SHA-256 hash of it, so a leaked database cannot be signed in with. They last
12 hours. Signing in, failing to sign in, and signing out are all written to
`audit_logs`.

To add an endpoint: write the handler method in the package it belongs to
under `internal/api/http/`, then mount it in `registerRoutes` in
`internal/api/server.go` — the one place that says which paths exist. Put it
behind `session.Require` unless it is meant to be public.

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

The header carries the two top-level areas — Dashboard and Logs — with the
account menu on the right, which is where Profile and Sign out live. The
dashboard has its own sidebar:

```
Activity                      metrics and recent events   (real)
Applications  Applications · APIs · SSO integrations
Authentication  Database · Social · Login flows
User management  Users · Roles
Settings  Organization · Languages
```

Activity and Users are real; everything else renders the placeholder rows in
`lib/demo.ts`, marked with a dot in the sidebar and a badge on the page, so
nothing there is mistaken for real state. Delete a block from that file as
soon as its section talks to the API.

### Users and their fields

A user has an email and a verified flag as columns. Everything else an
organisation wants to keep is defined at runtime in `user_fields` and stored
in `users.data`, so adding a first name or a "phone verified" flag is a row in
a table rather than a migration.

`internal/model/user_field.go` owns what a field means: the types on offer
(`text`, `number`, `bool`, `email`, `date`) and `Normalise`, which checks a
submitted value and returns what should be stored — the one place to change
when adding a type. The panel builds both its table and its form from the same
list, so a new field appears as a column and an input without any code
changing.

Searching matches the email or any stored value, because `data` is searched as
text; the search and the verified filter live in the URL, so the server renders
the result and a filtered list can be linked to.

The profile page collects what belongs to the signed-in account: profile
information, two-factor, theme, language, email, password and sessions. Theme
and sign out work; the rest are marked "not wired up" and their controls are
disabled rather than pretending.

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

**Layout.** Pages fill the width, the way PocketBase does: a table with room
for its columns reads better than one centred in a narrow column. Below 55rem
the sidebar becomes a scrollable row above the content, and below 30rem the
header sections become their icons alone — their labels stay in the accessible
name.

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
