import { browser } from '$app/environment';
import { QueryClient } from '@tanstack/svelte-query';

/**
 * A query client for one visitor.
 *
 * It is built in the root layout, which means a new one per request while
 * rendering on the server: a client is a cache, and one shared between
 * requests would be one visitor's data shown to another.
 *
 * The defaults suit an admin panel. Pages arrive already rendered by the
 * server, so what is on screen is as fresh as the request that drew it —
 * asking again straight away would be a second request for the same rows.
 * After half a minute it is worth checking, and worth checking again when
 * someone comes back to the tab, since the data is shared with whoever else
 * is signed in.
 */
export function createQueryClient(): QueryClient {
	return new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 30_000,
				refetchOnWindowFocus: browser,
				// One retry covers a dropped connection; more would just delay
				// telling someone that the server is not answering.
				retry: browser ? 1 : false
			}
		}
	});
}
