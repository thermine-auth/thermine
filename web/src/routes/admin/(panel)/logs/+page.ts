import { adminApi } from '$lib/api';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	const { logs } = await adminApi.logs(100, fetch);

	return { logs };
};
