import { error } from '@sveltejs/kit';
import { ApiError, signIn, type SignInRequest } from '$lib/api';
import type { LayoutServerLoad } from './$types';

/**
 * The sign-in under way, when the page was given one.
 *
 * The authorization endpoint sends the browser here with `?request=<handle>`,
 * and every sign-in page passes it on. Loading it once, here, is what lets
 * each page show the application's name, logo and links in the first response.
 * A handle that expired or was used is not an error page: the page says so.
 */
export const load: LayoutServerLoad = async ({ url, fetch, setHeaders }) => {
	// A sign-in page is never cached: it is about one sign-in, for one person.
	setHeaders({ 'cache-control': 'no-store' });

	const request = url.searchParams.get('request');

	let signInRequest: SignInRequest | null = null;
	let expired = false;

	if (request) {
		try {
			signInRequest = await signIn.request(request, fetch);
		} catch (err) {
			if (err instanceof ApiError && (err.status === 410 || err.status === 404)) {
				expired = true;
			} else if (err instanceof ApiError && err.status === 0) {
				error(503, err.message);
			} else {
				throw err;
			}
		}
	}

	return { request, signInRequest, expired };
};
