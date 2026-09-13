# The xermess API server.
#
# Build from the repository root, since the build needs the whole module:
#
#   docker build -f scripts/Dockerfile -t xermess .
#
# The image holds the server and the migration command; scripts/entrypoint.sh
# runs the migrations and then the server. The admin panel in web/ is a
# separate app and is not part of this image.

# ---------------------------------------------------------------------------
# Build
# ---------------------------------------------------------------------------
FROM golang:1.27-alpine AS build

WORKDIR /src

# The module files come first, so the layer that downloads dependencies is
# reused whenever only the source has changed.
COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./migrations

ARG VERSION=dev
ARG COMMIT=unknown

RUN CGO_ENABLED=0 go build -trimpath \
	-ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
	-o /out/xermess ./cmd/xermess
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/migrate ./cmd/migrate

# ---------------------------------------------------------------------------
# Run
# ---------------------------------------------------------------------------
FROM alpine:3.22

# Certificates for outgoing TLS; tzdata so a timezone such as Asia/Bishkek
# means what it says.
RUN apk add --no-cache ca-certificates tzdata \
	&& adduser -D -H -u 10001 xermess

WORKDIR /app

COPY --from=build /out/xermess /app/xermess
COPY --from=build /out/migrate /app/migrate
COPY scripts/entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

USER xermess

EXPOSE 8080

# Everything else is configuration, and is passed in as environment variables:
# XERMESS_DB_DSN above all. See .env.example.
ENV XERMESS_ADDR=:8080

ENTRYPOINT ["/app/entrypoint.sh"]
