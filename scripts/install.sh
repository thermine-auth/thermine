#!/usr/bin/env bash
#
# Set xermess up on a fresh machine: check what is needed, write a .env,
# create the database and bring it up to date.
#
#   scripts/install.sh
#
# It is safe to run twice — nothing here overwrites what is already there.
# `make setup` runs this.

set -euo pipefail

# shellcheck source=scripts/lib.sh
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"
cd "$root"

echo "==> checking what is installed"
need go "install Go 1.27 or newer from https://go.dev/dl/"
need psql "install the PostgreSQL client tools"

# bun puts itself in ~/.bun/bin, which a non-interactive shell does not always
# have on its PATH.
[[ -x "$HOME/.bun/bin/bun" ]] && PATH="$HOME/.bun/bin:$PATH"
command -v bun >/dev/null 2>&1 || echo "note: bun is missing, so the admin panel cannot be built"

echo "==> .env"
if [[ -f .env ]]; then
	echo "    .env already exists, leaving it alone"
else
	cp .env.example .env
	echo "    created .env — check the database settings in it before going on"
fi

load_env

echo "==> Go dependencies"
go mod download

echo "==> database"
./scripts/db.sh create | sed 's/^/    /'

echo "==> migrations"
go run ./cmd/migrate up

if command -v bun >/dev/null 2>&1; then
	echo "==> admin panel dependencies"
	(cd web && bun install)
fi

cat <<'DONE'

Ready.

  make run        start the server
  make web-dev    start the admin panel

DONE
