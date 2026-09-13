import type { Handle } from '@sveltejs/kit';

import { COOKIES } from '$lib/constants';

/**
 * Puts the reader's theme into the HTML before it is sent.
 *
 * The panel renders in the browser, so for the first frames there is no
 * stylesheet and no JavaScript: whatever `<html>` carries is what gets
 * painted. Deciding it here means a reader who chose dark is served dark,
 * rather than being shown light and corrected a moment later.
 *
 * A reader who has not chosen yet gets no attribute at all, which leaves the
 * prefers-color-scheme rules in tokens.css to decide — also without a flash,
 * because that happens in CSS rather than after it.
 */
export const handle: Handle = async ({ event, resolve }) => {
	const saved = event.cookies.get(COOKIES.theme);
	const theme = saved === 'dark' || saved === 'light' ? saved : null;

	const attributes = theme
		? `data-theme="${theme}" style="background-color: ${theme === 'dark' ? '#1c1c1c' : '#fff'}; color-scheme: ${theme}"`
		: '';

	return resolve(event, {
		transformPageChunk: ({ html }) => html.replace('__THEME__', attributes)
	});
};
