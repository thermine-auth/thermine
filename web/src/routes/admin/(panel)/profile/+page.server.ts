import type { AdminSession } from '$lib/api';
import { apiGet } from '$lib/server/api';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
	const { sessions } = await apiGet<{ sessions: AdminSession[] }>(
		'/admin/sessions',
		cookies,
		fetch
	);

	return { sessions };
};
