import { redirect } from '@sveltejs/kit';
import { resolve } from '$app/paths';
import { adminApi, ApiError } from '$lib/api';
import type { LayoutLoad } from './$types';

/**
 * Everything in this group needs a session. Asking the server is the only way
 * to know: the cookie is HttpOnly, so no script can read it.
 *
 * Doing it here rather than in each page means a page never renders before
 * the session is known, and the administrator it returns is available to
 * every page below through `data`.
 */
export const load: LayoutLoad = async ({ fetch }) => {
	try {
		const { admin } = await adminApi.me(fetch);
		return { admin };
	} catch (error) {
		if (error instanceof ApiError && error.isUnauthorized) {
			redirect(307, resolve('/admin/login'));
		}

		throw error;
	}
};
