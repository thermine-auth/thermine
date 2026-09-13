import { api, type Fetch } from './client';
import type {
	Admin,
	AdminSession,
	FieldInput,
	FieldRules,
	LogEntry,
	Overview,
	UserField,
	UserInput,
	UserRecord
} from './types';

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
	create: (input: UserInput) => api.post<{ user: UserRecord }>('/admin/users', input),

	update: (id: string, input: UserInput) =>
		api.patch<{ user: UserRecord }>(`/admin/users/${id}`, input),

	remove: (id: string) => api.delete<void>(`/admin/users/${id}`),

	addField: (field: FieldInput) => api.post<{ field: UserField }>('/admin/user-fields', field),

	updateField: (id: string, rules: FieldRules & { label: string }) =>
		api.patch<{ field: UserField }>(`/admin/user-fields/${id}`, rules),

	removeField: (id: string) => api.delete<void>(`/admin/user-fields/${id}`)
};
