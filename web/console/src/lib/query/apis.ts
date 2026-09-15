import { queryOptions } from '@tanstack/svelte-query';

import { apisApi, type API } from '$lib/api';
import { keys } from './keys';

/** The APIs matching a search, seeded with what the server rendered. An empty
    search is every API, which the role panel offers scopes from. */
export function apisOptions(search: string, initial: API[]) {
	return queryOptions({
		queryKey: keys.apis.list(search),
		queryFn: async () => (await apisApi.list(search)).apis,
		initialData: initial
	});
}

/** One API, seeded with what the server rendered for its page. */
export function apiOptions(id: string, initial: API) {
	return queryOptions({
		queryKey: keys.apis.one(id),
		queryFn: async () => (await apisApi.get(id)).api,
		initialData: initial
	});
}
