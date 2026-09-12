/** The shapes the API returns. They mirror the Go responses in
    internal/server/admin.go; change them together. */

export type Admin = {
	id: string;
	username: string;
	email: string;
	full_name: string;
	status: string;
	roles: string[];
	last_login_at?: string;
};

export type ActivityEvent = {
	id: string;
	action: string;
	actor: string;
	ip: string;
	created_at: string;
};

export type Overview = {
	counts: {
		admins: number;
		active_sessions: number;
		roles: number;
		events: number;
	};
	activity: ActivityEvent[];
};

export type AdminSession = {
	id: string;
	ip: string;
	user_agent: string;
	created_at: string;
	expires_at: string;
	active: boolean;
};
