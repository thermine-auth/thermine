import { queryOptions } from '@tanstack/svelte-query';

import { rolesApi, type Role, type RolePage } from '$lib/api';
import { keys } from './keys';

/** What the list is asked for, all of it from the URL: the tab (global or
    application roles), the application chip, the search box and the default
    filter. An empty application means every application. */
export type RoleListParams = {
	scope: 'global' | 'application';
	application: string;
	search: string;
	isDefault: string;
};

/** How many roles the pickers ask for: the most the API gives in one page. */
export const ROLE_CHOICES_LIMIT = 500;

/**
 * The roles matching a search.
 *
 * `initial` is what the server already rendered, so opening the page does not
 * fetch the same rows a second time.
 */
export function rolesOptions(params: RoleListParams, initial: RolePage) {
	return queryOptions({
		queryKey: keys.roles.list(params),
		queryFn: () =>
			rolesApi.list({
				scope: params.scope,
				application: params.application || undefined,
				search: params.search,
				default: params.isDefault
			}),
		initialData: initial
	});
}

/** Every role of every application the administrator can see, whatever the
    list on screen is filtered to: the user panel offers them to hold, and the
    role panel offers an application's to inherit. */
export function roleChoicesOptions(initial: Role[]) {
	return queryOptions({
		queryKey: keys.roles.choices,
		queryFn: async () => (await rolesApi.list({ limit: ROLE_CHOICES_LIMIT })).roles,
		initialData: initial
	});
}
