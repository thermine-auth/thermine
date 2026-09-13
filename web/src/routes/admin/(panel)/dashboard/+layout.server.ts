import { COOKIES } from '$lib/constants';
import { isSidebarState, type SidebarState } from '$lib/state/sidebar';
import type { LayoutServerLoad } from './$types';

/**
 * Reads how the reader left the sidebar, so the first frame is already the
 * right width instead of a full column that folds a moment later.
 */
export const load: LayoutServerLoad = ({ cookies }) => {
	const saved = cookies.get(COOKIES.sidebar);
	const sidebar: SidebarState = isSidebarState(saved) ? saved : 'wide';

	return { sidebar };
};
