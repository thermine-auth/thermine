import { account } from '$lib/api';
import { asUser } from '$lib/server/session';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, url }) => {
	const { sessions } = await asUser(() => account.sessions(fetch), url);

	return { sessions };
};
