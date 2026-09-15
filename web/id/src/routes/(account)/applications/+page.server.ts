import { account } from '$lib/api';
import { asUser } from '$lib/server/session';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, url }) => {
	const { applications } = await asUser(() => account.applications(fetch), url);

	return { applications };
};
