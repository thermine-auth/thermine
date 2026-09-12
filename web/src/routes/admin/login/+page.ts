import { redirect } from '@sveltejs/kit';
import { resolve } from '$app/paths';
import { adminApi } from '$lib/api';
import type { PageLoad } from './$types';

/** Someone who is already signed in has no use for the sign-in page. */
export const load: PageLoad = async ({ fetch }) => {
	try {
		await adminApi.me(fetch);
	} catch {
		return {};
	}

	redirect(307, resolve('/admin/dashboard'));
};
