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
| Frontend    | `web-install`, `web-dev`, `web-build`                                |
| Setup       | `setup` (create `.env`), `tools`                                     |

## Layout

```
cmd/xermess/main.go            startup, in order, in one function
cmd/migrate/main.go            runs migrations by hand: up, down, status

internal/config/config.go      reads .env
internal/database/database.go  opens the connection
internal/database/migrate.go   applies migrations
internal/server/server.go      the server: middleware, routes, handlers
internal/server/middleware.go  request logging and CORS
internal/model/                one file per table, listed in model.All

migrations/                    one Go file per migration, applied in order
web/                           SvelteKit frontend
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

| Method | Path            | Description        |
| ------ | --------------- | ------------------ |
| `GET`  | `/healthz`      | The server is up   |
| `GET`  | `/api/v1/hello` | Placeholder        |

To add an endpoint: mount it in `registerRoutes` and write its handler below,
both in `internal/server/server.go`.

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
