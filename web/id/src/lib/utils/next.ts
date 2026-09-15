/**
 * Where to go after signing in, from a `?next=` parameter.
 *
 * Only a path on this site is accepted: `next=https://evil.example` or
 * `next=//evil.example` would turn the sign-in page into a way to send someone
 * anywhere with this site's name on the link.
 */
export function safeNext(next: string | null, fallback = '/'): string {
	if (!next || !next.startsWith('/') || next.startsWith('//') || next.startsWith('/\\')) {
		return fallback;
	}

	return next;
}
