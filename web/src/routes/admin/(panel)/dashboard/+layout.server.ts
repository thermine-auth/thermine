import { isSidebarState, SIDEBAR_COOKIE, type SidebarState } from '$lib/sidebar';
import type { LayoutServerLoad } from './$types';

/**
 * Reads how the reader left the sidebar, so the first frame is already the
 * right width instead of a full column that folds a moment later.
 */
export const load: LayoutServerLoad = ({ cookies }) => {
	const saved = cookies.get(SIDEBAR_COOKIE);
	const sidebar: SidebarState = isSidebarState(saved) ? saved : 'wide';

	return { sidebar };
};
