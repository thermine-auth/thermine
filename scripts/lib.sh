# Shared by the scripts in this directory. Source it, do not run it:
#
#   . "$(dirname "${BASH_SOURCE[0]}")/lib.sh"
#
# It leaves $root at the repository root and offers the few things every
# script here needs.

# The repository root, whatever directory the script was called from.
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# need checks for a command and explains what to do when it is missing.
need() {
	command -v "$1" >/dev/null 2>&1 || {
		echo "missing: $1 — $2" >&2
		exit 1
	}
}

# load_env reads .env into the environment, the same file the server reads.
# Anything already set wins, so DB_URL=… on the command line still overrides.
load_env() {
	[[ -f "$root/.env" ]] || return 0

	local key value
	# The `|| [[ -n "$key" ]]` is for a file whose last line has no newline:
	# read returns false at EOF even when it just read something.
	while IFS='=' read -r key value || [[ -n "$key" ]]; do
		[[ "$key" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]] || continue
		[[ -n "${!key-}" ]] && continue

		# Strip one layer of quotes, the way a shell would.
		value="${value%\"}"
		value="${value#\"}"
		export "$key=$value"
	done <"$root/.env"
}

# dsn is the database to work on: whatever was passed in, else .env.
dsn() {
	echo "${DB_URL:-${XERMESS_DB_DSN:-}}"
}

# db_name is the database a DSN points at.
db_name() {
	echo "$1" | sed -E 's|.*/([^/?]+)(\?.*)?$|\1|'
}

# db_server is the same server with the "postgres" database instead. Creating
# a database means connecting to a different one.
db_server() {
	echo "$1" | sed -E 's|/[^/?]+(\?.*)?$|/postgres\1|'
}
