/**
 * The names the cache knows things by.
 *
 * Keys are built here rather than written out at each call, so invalidating
 * "the users" after a write cannot miss a spelling. A key is a list read
 * left to right: everything under `users.all` goes when a user changes.
 */
export const keys = {
	users: {
		all: ['users'] as const,
		/** One page of the list, as the search box and filter describe it. */
		list: (params: { search: string; verified: string }) =>
			['users', 'list', params.search, params.verified] as const,
		/** The fields a user record is made of. */
		fields: ['users', 'fields'] as const
	},

	admin: {
		sessions: ['admin', 'sessions'] as const,
		logs: (limit: number) => ['admin', 'logs', limit] as const,
		overview: ['admin', 'overview'] as const
	}
} as const;
