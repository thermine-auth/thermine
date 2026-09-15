import { resolve } from '$app/paths';

/** The sign-in pages that pass the sign-in handle from one to the next. */
export type AuthPage = '/login' | '/register' | '/forgot-password';

/**
 * A link to another sign-in page that keeps the sign-in under way: without the
 * handle, registering or resetting a password could not send the user back to
 * the application that asked.
 */
export function authHref(page: AuthPage, request: string | null): string {
	const path = resolve(page);
	return request ? `${path}?request=${encodeURIComponent(request)}` : path;
}

/** Sends the browser on to where the API said: usually back to an
    application, on another origin, so it is a full navigation. */
export function leaveTo(location: string): void {
	window.location.assign(location);
}
