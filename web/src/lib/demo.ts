/**
 * Placeholder data for the dashboard sections that have no backend yet.
 *
 * Nothing here is real: it exists so the sidebar has somewhere to lead and so
 * the tables can be designed against realistic shapes. Delete a block as soon
 * as its section talks to the API.
 */

export type DemoAdmin = {
	id: string;
	name: string;
	username: string;
	email: string;
	role: string;
	status: 'active' | 'invited' | 'suspended';
	lastSeen: string;
};

export const demoAdmins: DemoAdmin[] = [
	{
		id: 'a1',
		name: 'Ada Lovelace',
		username: 'ada',
		email: 'ada@xermess.dev',
		role: 'super_admin',
		status: 'active',
		lastSeen: '2 minutes ago'
	},
	{
		id: 'a2',
		name: 'Grace Hopper',
		username: 'grace',
		email: 'grace@xermess.dev',
		role: 'admin',
		status: 'active',
		lastSeen: '3 hours ago'
	},
	{
		id: 'a3',
		name: 'Alan Turing',
		username: 'alan',
		email: 'alan@xermess.dev',
		role: 'support',
		status: 'invited',
		lastSeen: 'Never'
	},
	{
		id: 'a4',
		name: 'Katherine Johnson',
		username: 'katherine',
		email: 'katherine@xermess.dev',
		role: 'auditor',
		status: 'active',
		lastSeen: 'Yesterday'
	},
	{
		id: 'a5',
		name: 'Edsger Dijkstra',
		username: 'edsger',
		email: 'edsger@xermess.dev',
		role: 'support',
		status: 'suspended',
		lastSeen: '3 weeks ago'
	}
];

export type DemoRole = {
	id: string;
	name: string;
	description: string;
	permissions: number;
	admins: number;
};

export const demoRoles: DemoRole[] = [
	{
		id: 'r1',
		name: 'super_admin',
		description: 'Everything, including other administrators',
		permissions: 24,
		admins: 1
	},
	{
		id: 'r2',
		name: 'admin',
		description: 'Day to day administration',
		permissions: 17,
		admins: 1
	},
	{
		id: 'r3',
		name: 'support',
		description: 'Read accounts and reset passwords',
		permissions: 8,
		admins: 2
	},
	{
		id: 'r4',
		name: 'auditor',
		description: 'Read only, including the activity log',
		permissions: 5,
		admins: 1
	}
];

export type DemoKey = {
	id: string;
	label: string;
	prefix: string;
	scopes: string;
	created: string;
	lastUsed: string;
	active: boolean;
};

export const demoKeys: DemoKey[] = [
	{
		id: 'k1',
		label: 'Web app',
		prefix: 'xm_live_8f2a…',
		scopes: 'tokens:issue, users:read',
		created: '12 Aug 2026',
		lastUsed: '4 minutes ago',
		active: true
	},
	{
		id: 'k2',
		label: 'Mobile app',
		prefix: 'xm_live_b71c…',
		scopes: 'tokens:issue',
		created: '3 Jul 2026',
		lastUsed: '2 hours ago',
		active: true
	},
	{
		id: 'k3',
		label: 'Nightly export',
		prefix: 'xm_live_0d94…',
		scopes: 'users:read',
		created: '19 Feb 2026',
		lastUsed: '6 weeks ago',
		active: false
	}
];

export type DemoWebhook = {
	id: string;
	event: string;
	url: string;
	delivered: number;
	failed: number;
	active: boolean;
};

export const demoWebhooks: DemoWebhook[] = [
	{
		id: 'w1',
		event: 'admin.login',
		url: 'https://hooks.xermess.dev/audit',
		delivered: 1842,
		failed: 0,
		active: true
	},
	{
		id: 'w2',
		event: 'admin.login_failed',
		url: 'https://hooks.xermess.dev/alerts',
		delivered: 96,
		failed: 2,
		active: true
	},
	{
		id: 'w3',
		event: 'user.created',
		url: 'https://crm.example.com/xermess',
		delivered: 12043,
		failed: 31,
		active: false
	}
];
