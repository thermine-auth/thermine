#!/usr/bin/env bash
#
# The database this project runs on.
#
#   scripts/db.sh create    create it, if it is not there yet
#   scripts/db.sh reset     empty it and migrate from scratch — asks first
#   scripts/db.sh psql      open a psql session on it
#
# Which database is read from .env, or from DB_URL when that is set:
#
#   DB_URL=postgres://…/xermess_test scripts/db.sh create

set -euo pipefail

# shellcheck source=scripts/lib.sh
. "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib.sh"
cd "$root"

need psql "install the PostgreSQL client tools"
load_env

url="$(dsn)"
[[ -n "$url" ]] || {
	echo "no database: set XERMESS_DB_DSN in .env, or pass DB_URL=…" >&2
	exit 1
}

name="$(db_name "$url")"
server="$(db_server "$url")"

create() {
	if psql "$server" -tAc "SELECT 1 FROM pg_database WHERE datname='$name'" | grep -q 1; then
		echo "database $name already exists"
		return
	fi

	psql "$server" -c "CREATE DATABASE \"$name\"" >/dev/null
	echo "created database $name"
}

reset() {
	printf 'This deletes every table and row in "%s". Type yes to continue: ' "$name"
	read -r answer
	[[ "$answer" == "yes" ]] || {
		echo "cancelled"
		exit 1
	}

	psql "$url" -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	XERMESS_DB_DSN="$url" go run ./cmd/migrate up
}

case "${1:-}" in
create) create ;;
reset) reset ;;
psql) exec psql "$url" ;;
*)
	echo "usage: scripts/db.sh create|reset|psql" >&2
	exit 1
	;;
esac
