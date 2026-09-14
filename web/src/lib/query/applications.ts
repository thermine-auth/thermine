import { queryOptions } from '@tanstack/svelte-query';

import { applicationsApi, type Application, type ApplicationPage } from '$lib/api';
import { keys } from './keys';

/** How many applications the pickers ask for: the most the API gives. */
export const APPLICATION_CHOICES_LIMIT = 500;

/** What the list is asked for, all of it from the URL. */
export type ApplicationListParams = { search: string; type: string };

/** The applications matching a search, seeded with what the server rendered. */
export function applicationsOptions(params: ApplicationListParams, initial: ApplicationPage) {
	return queryOptions({
		queryKey: keys.applications.list(params),
		queryFn: () => applicationsApi.list(params),
		initialData: initial
	});
}

/** Every application the administrator can see, for the pickers that offer
    them: the roles page, and the user and administrator panels. */
export function applicationChoicesOptions(initial: Application[]) {
	return queryOptions({
		queryKey: keys.applications.choices,
		queryFn: async () =>
			(await applicationsApi.list({ limit: APPLICATION_CHOICES_LIMIT })).applications,
		initialData: initial
	});
}
