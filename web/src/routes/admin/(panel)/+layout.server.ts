import type { Admin } from '$lib/api';
import { apiGet } from '$lib/server/api';
import type { LayoutServerLoad } from './$types';

/**
 * Everything in this group needs a session. Doing the check on the server
 * means a page is rendered already signed in, rather than appearing and then
 * being replaced once the browser has asked.
 */
export const load: LayoutServerLoad = async ({ cookies, fetch }) => {
	const { admin } = await apiGet<{ admin: Admin }>('/admin/me', cookies, fetch);

	return { admin };
};
