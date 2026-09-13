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

/** What the first administrator is made from. There is no username: the
    address is the account. */
export type SetupInput = {
	email: string;
	password: string;
	first_name: string;
	last_name: string;
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

/** One field of a user record.
 *
 *  `builtin` says which kind it is: a built-in field is a column of the
 *  record — every installation has it, and it cannot be changed or removed —
 *  while an additional one was added in the panel and its values live in
 *  `data`. The API returns both in one list, built-ins first. */
export type UserField = FieldRules & {
	/** The row this field is. A built-in field is a column rather than a row,
	    so it has none. */
	id?: string;
	name: string;
	label: string;
	type: FieldType;
	position: number;
	builtin: boolean;
};

/** What the panel sends when adding a field. A field's name and type are
    fixed once records hold values under them, so an update sends the rules
    only. */
export type FieldInput = FieldRules & {
	name: string;
	label: string;
	type: FieldType;
};

/** The built-in fields of a user record: the columns every installation has. */
export type UserBuiltins = {
	email: string;
	email_verified: boolean;
	first_name: string;
	last_name: string;
	is_active: boolean;
};

/** A user. The built-in fields are its own properties; everything an
    organisation added lives in `data`, keyed by field name. */
export type UserRecord = UserBuiltins & {
	id: string;
	data: Record<string, unknown> | null;
	created_at: string;
	updated_at: string;
};

/** What the panel sends when creating or updating a user. */
export type UserInput = UserBuiltins & {
	data: Record<string, unknown>;
};

export type UserPage = {
	users: UserRecord[];
	total: number;
	limit: number;
	offset: number;
};
