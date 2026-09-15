import type { AdminSession, MfaStatus } from '$lib/api';
import { apiGet } from '$lib/server/api';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies, fetch }) => {
	const [{ sessions }, { mfa }] = await Promise.all([
		apiGet<{ sessions: AdminSession[] }>('/admin/sessions', cookies, fetch),
		apiGet<{ mfa: MfaStatus }>('/admin/mfa', cookies, fetch)
	]);

	return { sessions, mfa };
};
