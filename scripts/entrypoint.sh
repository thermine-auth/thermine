#!/usr/bin/env sh
#
# What the container runs: wait for the database, bring it up to date, then
# hand over to the server.
#
# Migrations run here rather than inside the server so a container that is
# restarted does not race another one applying the same migration; set
# XERMESS_DB_MIGRATE=false and this still works.

set -eu

echo "entrypoint: waiting for the database"

# The server itself will fail fast on a bad DSN; this only covers a database
# that is still starting beside us, which is the usual case with compose.
attempt=1
until /app/migrate status >/dev/null 2>&1; do
	if [ "$attempt" -ge 30 ]; then
		echo "entrypoint: the database did not come up" >&2
		exit 1
	fi

	attempt=$((attempt + 1))
	sleep 1
done

echo "entrypoint: applying migrations"
/app/migrate up

echo "entrypoint: starting xermess"
exec /app/xermess "$@"
