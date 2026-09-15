import type { LogEntry } from '$lib/api';
import { apiGet, requirePermission } from '$lib/server/api';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies, fetch, parent }) => {
	requirePermission((await parent()).admin, 'activity.read');

	const { logs } = await apiGet<{ logs: LogEntry[] }>('/admin/logs?limit=100', cookies, fetch);

	return { logs };
};
