/** A date as a person reads it: "15 Sep 2026". */
export function formatDate(value: string | null | undefined): string {
	if (!value) return '—';
	return new Date(value).toLocaleDateString(undefined, {
		day: 'numeric',
		month: 'short',
		year: 'numeric'
	});
}

/** How long ago something was: "just now", "5 minutes ago", "3 days ago". */
export function timeAgo(value: string, now = Date.now()): string {
	const seconds = Math.max(0, Math.round((now - new Date(value).getTime()) / 1000));

	const units: [number, string][] = [
		[60 * 60 * 24 * 365, 'year'],
		[60 * 60 * 24 * 30, 'month'],
		[60 * 60 * 24, 'day'],
		[60 * 60, 'hour'],
		[60, 'minute']
	];

	for (const [size, name] of units) {
		const count = Math.floor(seconds / size);
		if (count >= 1) return `${count} ${name}${count === 1 ? '' : 's'} ago`;
	}

	return 'just now';
}

/** A user agent as a person would name the device: "Chrome on macOS". */
export function describeDevice(userAgent: string): string {
	const ua = userAgent || '';

	const browser =
		[
			['Edge', /Edg\//],
			['Opera', /OPR\//],
			['Firefox', /Firefox\//],
			['Chrome', /Chrome\//],
			['Safari', /Safari\//]
		].find(([, pattern]) => (pattern as RegExp).test(ua))?.[0] ?? 'Unknown browser';

	const system =
		[
			['iPhone', /iPhone/],
			['iPad', /iPad/],
			['Android', /Android/],
			['Windows', /Windows/],
			['macOS', /Mac OS X|Macintosh/],
			['Linux', /Linux/]
		].find(([, pattern]) => (pattern as RegExp).test(ua))?.[0] ?? 'an unknown device';

	return `${browser} on ${system}`;
}

/** The initials of a user for an avatar: "AL", or the email's first letter. */
export function initials(user: { first_name: string; last_name: string; email: string }): string {
	const letters = `${user.first_name.charAt(0)}${user.last_name.charAt(0)}`.trim();
	return (letters || user.email.charAt(0)).toUpperCase();
}
