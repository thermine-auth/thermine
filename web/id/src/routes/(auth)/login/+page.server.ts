import { redirect } from '@sveltejs/kit';
import { currentUser } from '$lib/server/session';
import { safeNext } from '$lib/utils/next';
import type { PageServerLoad } from './$types';

/**
 * Someone already signed in who opens the sign-in page on its own — not sent
 * by an application — has nothing to do here, and goes to their account.
 * With a sign-in under way the page is shown regardless: the application may
 * have asked for the password again (prompt=login).
 */
export const load: PageServerLoad = async ({ url, cookies, fetch }) => {
	const next = safeNext(url.searchParams.get('next'));

	if (!url.searchParams.get('request') && (await currentUser(cookies, fetch))) {
		redirect(303, next);
	}

	return { next };
};
