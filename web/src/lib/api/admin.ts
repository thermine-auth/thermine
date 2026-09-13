import { api, type Fetch } from './client';
import type {
	Admin,
	AdminSession,
	FieldInput,
	FieldRules,
	LogEntry,
	Overview,
	SetupInput,
	UserField,
	UserInput,
	UserPage,
	UserRecord
} from './types';

/** Setting the panel up: the two calls that work without a session, because
    before the first administrator exists there is nobody to be. */
export const setupApi = {
	status: (fetcher?: Fetch) => api.get<{ required: boolean }>('/admin/setup', fetcher),

	create: (input: SetupInput) => api.post<{ admin: { id: string } }>('/admin/setup', input)
};

/** Every call the admin panel makes. */
export const adminApi = {
	login: (username: string, password: string) =>
		api.post<{ admin: Admin }>('/admin/auth/login', { username, password }),

	logout: () => api.post<{ status: string }>('/admin/auth/logout'),

	me: (fetcher?: Fetch) => api.get<{ admin: Admin }>('/admin/me', fetcher),

	overview: (fetcher?: Fetch) => api.get<Overview>('/admin/overview', fetcher),

	sessions: (fetcher?: Fetch) => api.get<{ sessions: AdminSession[] }>('/admin/sessions', fetcher),

	logs: (limit = 50, fetcher?: Fetch) =>
		api.get<{ logs: LogEntry[] }>(`/admin/logs?limit=${limit}`, fetcher)
};

/** The users an organisation manages, and the shape of their records. */
export const usersApi = {
	/** A page of users. The search and the filter are the same ones the URL
	    carries, so a link and a query key describe the same list. */
	list: (params: { search?: string; verified?: string }, fetcher?: Fetch) => {
		const query = new URLSearchParams();
		if (params.search) query.set('search', params.search);
		if (params.verified === 'true' || params.verified === 'false') {
			query.set('verified', params.verified);
		}

		return api.get<UserPage>(`/admin/users?${query}`, fetcher);
	},

	fields: (fetcher?: Fetch) => api.get<{ fields: UserField[] }>('/admin/user-fields', fetcher),

	create: (input: UserInput) => api.post<{ user: UserRecord }>('/admin/users', input),

	update: (id: string, input: UserInput) =>
		api.patch<{ user: UserRecord }>(`/admin/users/${id}`, input),

	remove: (id: string) => api.delete<void>(`/admin/users/${id}`),

	addField: (field: FieldInput) => api.post<{ field: UserField }>('/admin/user-fields', field),

	updateField: (id: string, rules: FieldRules & { label: string }) =>
		api.patch<{ field: UserField }>(`/admin/user-fields/${id}`, rules),

	removeField: (id: string) => api.delete<void>(`/admin/user-fields/${id}`)
};
