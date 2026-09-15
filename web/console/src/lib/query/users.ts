import { queryOptions } from '@tanstack/svelte-query';

import { usersApi, type UserField, type UserPage } from '$lib/api';
import { keys } from './keys';

/** What the list is asked for: the search box, the verified filter and the
    role filter, all of which live in the URL. */
export type UserListParams = { search: string; verified: string; role: string };

/**
 * The users matching a search.
 *
 * `initial` is what the server already rendered, so opening the page does not
 * fetch the same rows a second time. It belongs to this key only: a new
 * search is a new key, and a new page load, which brings its own.
 */
export function usersOptions(params: UserListParams, initial: UserPage) {
	return queryOptions({
		queryKey: keys.users.list(params),
		queryFn: () => usersApi.list(params),
		initialData: initial
	});
}

/** The fields every user record has. They change rarely, and everything on
    the page needs them, so they are one query the whole page shares. */
export function userFieldsOptions(initial: UserField[]) {
	return queryOptions({
		queryKey: keys.users.fields,
		queryFn: async () => (await usersApi.fields()).fields,
		initialData: initial
	});
}
