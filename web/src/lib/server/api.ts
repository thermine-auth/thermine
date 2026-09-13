import { error, redirect, type Cookies } from '@sveltejs/kit';
import { resolve } from '$app/paths';
import { PUBLIC_API_URL } from '$env/static/public';

import { COOKIES } from '$lib/constants';

/**
 * Calls the API as the administrator making this request.
 *
 * A fetch from the server carries none of the browser's cookies, so the
 * session has to be passed on by hand. Everything the panel renders goes
 * through here, which is what lets the pages be rendered on the server
 * instead of assembled in the browser afterwards.
 */
async function call(path: string, cookies: Cookies, fetch: typeof globalThis.fetch) {
	const session = cookies.get(COOKIES.session);

	try {
		return await fetch(`${PUBLIC_API_URL}/api/v1${path}`, {
			headers: session ? { cookie: `${COOKIES.session}=${session}` } : {}
		});
	} catch {
		error(503, `Could not reach the API at ${PUBLIC_API_URL}. Is it running?`);
	}
}

/**
 * Reads from the API, sending anyone without a usable session to sign in.
 * Every page behind the panel needs that same treatment, so it lives here.
 */
export async function apiGet<T>(
	path: string,
	cookies: Cookies,
	fetch: typeof globalThis.fetch
): Promise<T> {
	const response = await call(path, cookies, fetch);

	if (response.status === 401) {
		redirect(307, resolve('/admin/login'));
	}

	if (!response.ok) {
		error(response.status, 'The server could not answer that');
	}

	return response.json() as Promise<T>;
}

/**
 * Whether the panel still has to be set up — that is, whether the API has no
 * administrator yet.
 *
 * The sign-in page asks before it draws itself: with no account to sign in
 * to, the only useful thing to show is the form that makes one.
 */
export async function setupRequired(
	cookies: Cookies,
	fetch: typeof globalThis.fetch
): Promise<boolean> {
	const response = await call('/admin/setup', cookies, fetch);

	if (!response.ok) return false;

	const { required } = (await response.json()) as { required: boolean };

	return required;
}

/** Whether the request carries a session the API accepts. Used by the sign-in
 *  page, which sends people who already have one back to the panel. */
export async function hasSession(
	cookies: Cookies,
	fetch: typeof globalThis.fetch
): Promise<boolean> {
	if (!cookies.get(COOKIES.session)) return false;

	const response = await call('/admin/me', cookies, fetch);

	return response.ok;
}
