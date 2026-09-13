import { redirect } from '@sveltejs/kit';
import { resolve } from '$app/paths';
import { hasSession, setupRequired } from '$lib/server/api';
import type { PageServerLoad } from './$types';

/**
 * Who should not be here: someone already signed in, and someone arriving at
 * a panel that has no administrator yet — there is nothing to sign in to, so
 * they are sent to make the first account instead.
 */
export const load: PageServerLoad = async ({ cookies, fetch }) => {
	if (await hasSession(cookies, fetch)) {
		redirect(307, resolve('/admin/dashboard'));
	}

	if (await setupRequired(cookies, fetch)) {
		redirect(307, resolve('/admin/new-super-admin'));
	}

	return {};
};
