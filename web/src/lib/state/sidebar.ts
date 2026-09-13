import { COOKIES } from '$lib/constants';

/** How wide the dashboard sidebar is: the full column, or icons only.
 *
 *  The choice lives in a cookie rather than localStorage so the server can
 *  read it and render the panel already folded, the same reason the theme is
 *  kept in one. dashboard/+layout.server.ts reads this name. */
export type SidebarState = 'wide' | 'mini';

const ONE_YEAR = 60 * 60 * 24 * 365;

export function isSidebarState(value: unknown): value is SidebarState {
	return value === 'wide' || value === 'mini';
}

/** Writes the choice down, so the next page arrives the same way. */
export function rememberSidebar(state: SidebarState) {
	document.cookie = `${COOKIES.sidebar}=${state}; path=/; max-age=${ONE_YEAR}; samesite=lax`;
}
