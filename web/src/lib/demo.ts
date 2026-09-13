/**
 * Placeholder data for the sections that have no backend yet.
 *
 * Nothing here is real: it exists so the sidebar has somewhere to lead and so
 * the tables can be designed against realistic shapes. Delete a block as soon
 * as its section talks to the API.
 */

type Status = 'active' | 'disabled' | 'draft';

export type DemoApplication = {
	id: string;
	name: string;
	kind: string;
	clientId: string;
	users: number;
	status: Status;
};

export const demoApplications: DemoApplication[] = [
	{
		id: 'app1',
		name: 'Customer portal',
		kind: 'Single page app',
		clientId: 'xm_c7f2a9',
		users: 18420,
		status: 'active'
	},
	{
		id: 'app2',
		name: 'Mobile app',
		kind: 'Native',
		clientId: 'xm_2b91de',
		users: 9037,
		status: 'active'
	},
	{
		id: 'app3',
		name: 'Back office',
		kind: 'Regular web app',
		clientId: 'xm_44ac10',
		users: 126,
		status: 'active'
	},
	{
		id: 'app4',
		name: 'Partner sandbox',
		kind: 'Machine to machine',
		clientId: 'xm_9f0b77',
		users: 0,
		status: 'draft'
	}
];

export type DemoApi = {
	id: string;
	name: string;
	identifier: string;
	scopes: number;
	tokenLifetime: string;
	status: Status;
};

export const demoApis: DemoApi[] = [
	{
		id: 'api1',
		name: 'Accounts API',
		identifier: 'https://api.xermess.dev/accounts',
		scopes: 14,
		tokenLifetime: '1 hour',
		status: 'active'
	},
	{
		id: 'api2',
		name: 'Billing API',
		identifier: 'https://api.xermess.dev/billing',
		scopes: 8,
		tokenLifetime: '30 minutes',
		status: 'active'
	},
	{
		id: 'api3',
		name: 'Reporting API',
		identifier: 'https://api.xermess.dev/reports',
		scopes: 3,
		tokenLifetime: '12 hours',
		status: 'disabled'
	}
];

export type DemoSso = {
	id: string;
	name: string;
	protocol: string;
	domain: string;
	users: number;
	status: Status;
};

export const demoSso: DemoSso[] = [
	{
		id: 'sso1',
		name: 'Acme Corp',
		protocol: 'SAML 2.0',
		domain: 'acme.com',
		users: 2140,
		status: 'active'
	},
	{
		id: 'sso2',
		name: 'Globex',
		protocol: 'OIDC',
		domain: 'globex.io',
		users: 618,
		status: 'active'
	},
	{
		id: 'sso3',
		name: 'Initech',
		protocol: 'SAML 2.0',
		domain: 'initech.co',
		users: 0,
		status: 'draft'
	}
];

export type DemoConnection = {
	id: string;
	name: string;
	engine: string;
	users: number;
	applications: number;
	status: Status;
};

export const demoConnections: DemoConnection[] = [
	{
		id: 'db1',
		name: 'Username & password',
		engine: 'xermess',
		users: 27583,
		applications: 4,
		status: 'active'
	},
	{
		id: 'db2',
		name: 'Legacy accounts',
		engine: 'External MySQL',
		users: 4120,
		applications: 1,
		status: 'active'
	},
	{
		id: 'db3',
		name: 'Staff directory',
		engine: 'LDAP',
		users: 212,
		applications: 1,
		status: 'disabled'
	}
];

export type DemoProvider = {
	id: string;
	name: string;
	clientId: string;
	logins: number;
	status: Status;
};

export const demoProviders: DemoProvider[] = [
	{
		id: 'soc1',
		name: 'Google',
		clientId: '8417…apps.googleusercontent.com',
		logins: 12904,
		status: 'active'
	},
	{ id: 'soc2', name: 'GitHub', clientId: 'Iv1.4b2c…', logins: 3311, status: 'active' },
	{ id: 'soc3', name: 'Apple', clientId: 'dev.xermess.signin', logins: 1877, status: 'active' },
	{ id: 'soc4', name: 'Microsoft', clientId: 'f0c1…', logins: 0, status: 'draft' }
];

export type DemoFlow = {
	id: string;
	name: string;
	steps: string;
	applications: number;
	isDefault: boolean;
	status: Status;
};

export const demoFlows: DemoFlow[] = [
	{
		id: 'flow1',
		name: 'Standard login',
		steps: 'Identifier → Password → MFA',
		applications: 3,
		isDefault: true,
		status: 'active'
	},
	{
		id: 'flow2',
		name: 'Passwordless',
		steps: 'Identifier → Email code',
		applications: 1,
		isDefault: false,
		status: 'active'
	},
	{
		id: 'flow3',
		name: 'Staff login',
		steps: 'SSO → MFA',
		applications: 1,
		isDefault: false,
		status: 'active'
	},
	{
		id: 'flow4',
		name: 'Trial signup',
		steps: 'Identifier → Password',
		applications: 0,
		isDefault: false,
		status: 'draft'
	}
];

export type DemoUser = {
	id: string;
	name: string;
	email: string;
	connection: string;
	lastLogin: string;
	status: 'active' | 'invited' | 'blocked';
};

export const demoUsers: DemoUser[] = [
	{
		id: 'u1',
		name: 'Mira Chen',
		email: 'mira@acme.com',
		connection: 'Google',
		lastLogin: '4 minutes ago',
		status: 'active'
	},
	{
		id: 'u2',
		name: 'Tomas Neal',
		email: 'tomas@globex.io',
		connection: 'Username & password',
		lastLogin: '2 hours ago',
		status: 'active'
	},
	{
		id: 'u3',
		name: 'Priya Raman',
		email: 'priya@acme.com',
		connection: 'SAML · Acme Corp',
		lastLogin: 'Yesterday',
		status: 'active'
	},
	{
		id: 'u4',
		name: 'Jonas Weber',
		email: 'jonas@initech.co',
		connection: 'Username & password',
		lastLogin: 'Never',
		status: 'invited'
	},
	{
		id: 'u5',
		name: 'Ade Oyelaran',
		email: 'ade@globex.io',
		connection: 'GitHub',
		lastLogin: '3 weeks ago',
		status: 'blocked'
	}
];

export type DemoRole = {
	id: string;
	name: string;
	description: string;
	permissions: number;
	members: number;
};

export const demoRoles: DemoRole[] = [
	{
		id: 'r1',
		name: 'owner',
		description: 'Everything, including billing and other owners',
		permissions: 42,
		members: 2
	},
	{
		id: 'r2',
		name: 'administrator',
		description: 'Applications, connections and users',
		permissions: 31,
		members: 5
	},
	{
		id: 'r3',
		name: 'support',
		description: 'Read users and reset their passwords',
		permissions: 9,
		members: 12
	},
	{
		id: 'r4',
		name: 'auditor',
		description: 'Read only, including the activity log',
		permissions: 6,
		members: 3
	}
];

export type DemoLanguage = {
	id: string;
	name: string;
	code: string;
	translated: number;
	isDefault: boolean;
	status: Status;
};

export const demoLanguages: DemoLanguage[] = [
	{ id: 'l1', name: 'English', code: 'en', translated: 100, isDefault: true, status: 'active' },
	{ id: 'l2', name: 'Kyrgyz', code: 'ky', translated: 92, isDefault: false, status: 'active' },
	{ id: 'l3', name: 'Russian', code: 'ru', translated: 88, isDefault: false, status: 'active' },
	{ id: 'l4', name: 'Turkish', code: 'tr', translated: 41, isDefault: false, status: 'draft' }
];

/** The organisation this panel administers. */
export const demoOrganization = {
	name: 'Xermess',
	slug: 'xermess',
	domain: 'xermess.dev',
	region: 'eu-central',
	plan: 'Growth',
	created: '4 February 2026',
	supportEmail: 'support@xermess.dev'
};
