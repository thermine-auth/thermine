# xermess

An authentication server written in Go, with a SvelteKit web frontend.

> Early days — this is the project scaffold. The auth flows are not implemented yet.

## Requirements

- Go 1.27+
- [Bun](https://bun.sh) (for the `web/` frontend)

## Getting started

```sh
cp .env.example .env
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

| Target         | Description                     |
| -------------- | ------------------------------- |
| `build`        | Build the server into `bin/`    |
| `run`          | Run the server                  |
| `test`         | Run Go tests                    |
| `fmt`          | Format Go code                  |
| `vet`          | Run `go vet`                    |
| `tidy`         | Tidy `go.mod`                   |
| `web-install`  | Install frontend dependencies   |
| `web-dev`      | Run the frontend dev server     |
| `web-build`    | Build the frontend              |
| `clean`        | Remove build artifacts          |

Run `make` on its own to list them.

## Layout

```
cmd/xermess/        server entrypoint
internal/config/    environment configuration
internal/server/    engine setup (server.go) and the route table (routes.go)
internal/handler/   one file per resource: health.go, auth.go
web/                SvelteKit frontend
```

Adding an endpoint is two steps: write the method on a handler in
`internal/handler/`, then mount it in `registerRoutes` in
`internal/server/routes.go`.

## API

| Method | Path                    | Description                     |
| ------ | ----------------------- | ------------------------------- |
| `GET`  | `/healthz`              | Liveness check                  |
| `GET`  | `/readyz`               | Readiness check                 |
| `POST` | `/api/v1/auth/register` | Create an account               |
| `POST` | `/api/v1/auth/login`    | Exchange credentials for tokens |
| `POST` | `/api/v1/auth/refresh`  | Exchange a refresh token        |
| `POST` | `/api/v1/auth/logout`   | Revoke the current session      |
| `GET`  | `/api/v1/auth/me`       | The authenticated account       |

The `/auth` routes validate their request bodies but return `501
Not Implemented` — the logic is not written yet.

Errors all use one envelope:

```json
{ "error": { "code": "invalid_request", "message": "..." } }
```

## Configuration

Configuration comes from the environment. See `.env.example`.

| Variable       | Default | Description                       |
| -------------- | ------- | --------------------------------- |
| `XERMESS_ADDR` | `:8080` | Address the server listens on     |

## License

[MIT](LICENSE)
