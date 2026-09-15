import { requireUser } from '$lib/server/session';
import type { LayoutServerLoad } from './$types';

/** Every account page is for the signed-in user; anyone else signs in first,
    and comes back to the page they asked for. */
export const load: LayoutServerLoad = async ({ cookies, fetch, url, setHeaders }) => {
	setHeaders({ 'cache-control': 'private, no-store' });

	return { user: await requireUser(cookies, fetch, url) };
};
