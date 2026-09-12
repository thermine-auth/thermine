import { adminApi } from '$lib/api';
import type { PageLoad } from './$types';

/** Both calls are independent, so they go out together. */
export const load: PageLoad = async ({ fetch }) => {
	const [overview, sessions] = await Promise.all([
		adminApi.overview(fetch),
		adminApi.sessions(fetch)
	]);

	return { overview, sessions: sessions.sessions };
};
