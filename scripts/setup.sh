#!/usr/bin/env bash
#
# Prepare a fresh checkout: .env with its own secret key, the database and its
# migrations, and the web apps' dependencies. Safe to run again.

set -euo pipefail
. "$(dirname "${BASH_SOURCE[0]}")/lib.sh"

need go "https://go.dev/dl"
need bun "https://bun.sh"
need node "https://nodejs.org"
need psql "install the PostgreSQL client tools"
need openssl "install openssl"

[[ -f .env ]] || { cp .env.example .env && echo "created .env — review the database settings"; }

# A key that already exists stays: the stored signing keys are sealed with it.
if ! grep -qE '^XERMESS_SECRET_KEY=.{32,}' .env; then
	sed -i.bak '/^XERMESS_SECRET_KEY=/d' .env && rm .env.bak
	echo "XERMESS_SECRET_KEY=$(openssl rand -base64 32)" >>.env
	echo "wrote a new XERMESS_SECRET_KEY to .env"
fi

go mod download
scripts/db.sh create
go run ./cmd/migrate up

for app in console id; do
	echo "==> web/$app"
	(cd "web/$app" && bun install --frozen-lockfile)
	[[ ! -f web/$app/.env.example || -f web/$app/.env ]] || cp "web/$app/.env.example" "web/$app/.env"
done

echo "Ready. Start everything with: make dev"
