/**
 * The names both sides of the app have to agree on.
 *
 * A cookie is read in one place and written in another — the theme by the
 * hook that renders the page and by the toggle that changes it, the session
 * by the server loads — so the names live here rather than as a string in
 * each file that happens to spell them the same way.
 */
export const COOKIES = {
	/** The API's session. Set by the API itself, HttpOnly: the panel never
	    reads it, it only passes it on from a server load. */
	session: 'xermess_session',

	/** Light or dark, read while rendering so the page arrives themed. */
	theme: 'xermess-theme',

	/** Whether the dashboard sidebar is folded, read the same way. */
	sidebar: 'xermess-sidebar'
} as const;
