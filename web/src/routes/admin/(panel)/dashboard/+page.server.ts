import type { AdminSession, Overview } from '$lib/api';
import { apiGet } from '$lib/server/api';
import type { PageServerLoad } from './$types';

/** Both calls are independent, so they go out together. */
export const load: PageServerLoad = async ({ cookies, fetch }) => {
	const [overview, sessions] = await Promise.all([
		apiGet<Overview>('/admin/overview', cookies, fetch),
		apiGet<{ sessions: AdminSession[] }>('/admin/sessions', cookies, fetch)
	]);

	return { overview, sessions: sessions.sessions };
};
