import type { Handle } from '@sveltejs/kit';

import { COOKIES } from '$lib/constants';

export { handleFetch } from '$lib/server/proxy';

/**
 * Two things every response of this app gets.
 *
 * The reader's theme goes into the HTML before it is sent, so a reader who
 * chose dark is served dark rather than shown light and corrected a moment
 * later. A reader who has not chosen gets no attribute, and the
 * prefers-color-scheme rules in tokens.css decide — also before the first
 * paint.
 *
 * And no other site may frame any page: every one of them takes a password or
 * acts for a signed-in user, and a transparent frame over a fake page is how
 * clickjacking gets either.
 */
export const handle: Handle = async ({ event, resolve }) => {
	const saved = event.cookies.get(COOKIES.theme);
	const theme = saved === 'dark' || saved === 'light' ? saved : null;

	const attributes = theme ? `data-theme="${theme}" style="color-scheme: ${theme}"` : '';

	const response = await resolve(event, {
		transformPageChunk: ({ html }) => html.replace('__THEME__', attributes)
	});

	response.headers.set('X-Frame-Options', 'DENY');
	response.headers.set('Content-Security-Policy', "frame-ancestors 'none'");
	response.headers.set('Referrer-Policy', 'no-referrer');
	response.headers.set('X-Content-Type-Options', 'nosniff');

	return response;
};
