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

export type LogEntry = ActivityEvent & {
	user_agent: string;
	target_type: string;
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

/** The kinds of value a user field can hold. */
export type FieldType = 'text' | 'number' | 'bool' | 'email' | 'date';

/** The rules a field can put on its values. `min` and `max` bound a number's
    value or the length of text; `starts_with` is a prefix text must begin
    with. A rule that is not set is null. */
export type FieldRules = {
	required: boolean;
	unique: boolean;
	min: number | null;
	max: number | null;
	starts_with: string;
};

/** One user-defined column of the user record. */
export type UserField = FieldRules & {
	id: string;
	name: string;
	label: string;
	type: FieldType;
	position: number;
};

/** What the panel sends when adding a field. A field's name and type are
    fixed once records hold values under them, so an update sends the rules
    only. */
export type FieldInput = FieldRules & {
	name: string;
	label: string;
	type: FieldType;
};

/** A user. Everything beyond the email and its verified flag lives in `data`,
    described by the fields above. */
export type UserRecord = {
	id: string;
	email: string;
	email_verified: boolean;
	data: Record<string, unknown> | null;
	created_at: string;
	updated_at: string;
};

/** What the panel sends when creating or updating a user. */
export type UserInput = {
	email: string;
	email_verified: boolean;
	data: Record<string, unknown>;
};

export type UserPage = {
	users: UserRecord[];
	total: number;
	limit: number;
	offset: number;
};
