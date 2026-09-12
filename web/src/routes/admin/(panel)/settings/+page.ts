import { adminApi } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	const { sessions } = await adminApi.sessions(fetch);

	return { sessions };
};
