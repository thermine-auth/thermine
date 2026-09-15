# Sourced by the scripts beside it. Leaves the shell at the repository root.

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

# bun installs itself in ~/.bun/bin, which non-interactive shells often miss.
[[ -d "$HOME/.bun/bin" ]] && PATH="$HOME/.bun/bin:$PATH"

die() { echo "error: $*" >&2; exit 1; }

# need <command> <hint>
need() { command -v "$1" >/dev/null || die "$1 is missing — $2"; }

# load_env exports .env, the file the server reads. Variables already set win.
load_env() {
	[[ -f .env ]] || return 0
	local key value
	while IFS='=' read -r key value || [[ -n $key ]]; do
		[[ $key =~ ^[A-Za-z_][A-Za-z0-9_]*$ && -z ${!key+x} ]] || continue
		value="${value%\"}"
		export "$key=${value#\"}"
	done <.env
}
