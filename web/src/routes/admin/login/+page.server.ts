import { redirect } from '@sveltejs/kit';
import { resolve } from '$app/paths';
import { hasSession } from '$lib/server/api';
import type { PageServerLoad } from './$types';

/** Someone who is already signed in has no use for the sign-in page. */
export const load: PageServerLoad = async ({ cookies, fetch }) => {
	if (await hasSession(cookies, fetch)) {
		redirect(307, resolve('/admin/dashboard'));
	}

	return {};
};
