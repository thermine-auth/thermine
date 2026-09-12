import { api, type Fetch } from './client';
import type { Admin, AdminSession, Overview } from './types';

/** Every call the admin panel makes. */
export const adminApi = {
	login: (username: string, password: string) =>
		api.post<{ admin: Admin }>('/admin/auth/login', { username, password }),

	logout: () => api.post<{ status: string }>('/admin/auth/logout'),

	me: (fetcher?: Fetch) => api.get<{ admin: Admin }>('/admin/me', fetcher),

	overview: (fetcher?: Fetch) => api.get<Overview>('/admin/overview', fetcher),

	sessions: (fetcher?: Fetch) => api.get<{ sessions: AdminSession[] }>('/admin/sessions', fetcher)
};
