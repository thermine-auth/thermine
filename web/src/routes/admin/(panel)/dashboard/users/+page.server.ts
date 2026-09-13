import type { UserField, UserPage } from '$lib/api';
import { apiGet } from '$lib/server/api';
import type { PageServerLoad } from './$types';

/**
 * The search and the filter live in the URL, so the server can render the
 * result directly, the browser's back button works, and a filtered list can
 * be linked to.
 */
export const load: PageServerLoad = async ({ cookies, fetch, url }) => {
	const search = url.searchParams.get('search')?.trim() ?? '';
	const verified = url.searchParams.get('verified') ?? '';

	const query = new URLSearchParams();
	if (search) query.set('search', search);
	if (verified === 'true' || verified === 'false') query.set('verified', verified);

	const [page, fields] = await Promise.all([
		apiGet<UserPage>(`/admin/users?${query}`, cookies, fetch),
		apiGet<{ fields: UserField[] }>('/admin/user-fields', cookies, fetch)
	]);

	return { page, fields: fields.fields, search, verified };
};
