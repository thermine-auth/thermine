import type { LogEntry } from '$lib/api';
import { apiGet } from '$lib/server/api';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
	const { logs } = await apiGet<{ logs: LogEntry[] }>('/admin/logs?limit=100', cookies, fetch);

	return { logs };
};
