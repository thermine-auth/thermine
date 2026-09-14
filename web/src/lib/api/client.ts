import { PUBLIC_API_URL } from '$env/static/public';

/** Where the API lives. In development it is a different port from this app. */
const BASE_URL = PUBLIC_API_URL;

/** A failed request. `status` is 0 when the server could not be reached. */
export class ApiError extends Error {
	constructor(
		readonly status: number,
		message: string
	) {
		super(message);
		this.name = 'ApiError';
	}

	/** The session is missing or has expired. */
	get isUnauthorized(): boolean {
		return this.status === 401;
	}
}

/**
 * SvelteKit hands `load` functions their own fetch, which it uses to track
 * dependencies and to replay requests on the client. Passing it in is what
 * makes `invalidate` work; anything outside a load can leave it out.
 */
export type Fetch = typeof globalThis.fetch;

type RequestOptions = {
	method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
	body?: unknown;
	fetch?: Fetch;
};

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
	const { method = 'GET', body, fetch: fetcher = globalThis.fetch } = options;

	let response: Response;
	try {
		response = await fetcher(`${BASE_URL}/api/v1${path}`, {
			method,
			// The session lives in a cookie the browser will not send across
			// origins unless asked. The server allows this origin by name.
			credentials: 'include',
			headers: { 'Content-Type': 'application/json' },
			body: body === undefined ? undefined : JSON.stringify(body)
		});
	} catch {
		// fetch rejects both when the server never answered and when the
		// browser discarded the answer, which is what a missing CORS header
		// looks like from here. The two are indistinguishable to script, so
		// the message names both.
		throw new ApiError(
			0,
			`Could not reach ${BASE_URL}. It may be down, or this origin may not be in its XERMESS_CORS_ORIGINS.`
		);
	}

	// A 204 has no body to read.
	const payload = response.status === 204 ? {} : await response.json().catch(() => ({}));

	if (!response.ok) {
		throw new ApiError(response.status, payload.error ?? 'Something went wrong');
	}

	return payload as T;
}

export const api = {
	get: <T>(path: string, fetcher?: Fetch) => request<T>(path, { fetch: fetcher }),
	post: <T>(path: string, body?: unknown, fetcher?: Fetch) =>
		request<T>(path, { method: 'POST', body, fetch: fetcher }),
	put: <T>(path: string, body?: unknown, fetcher?: Fetch) =>
		request<T>(path, { method: 'PUT', body, fetch: fetcher }),
	patch: <T>(path: string, body?: unknown, fetcher?: Fetch) =>
		request<T>(path, { method: 'PATCH', body, fetch: fetcher }),
	delete: <T>(path: string, fetcher?: Fetch) =>
		request<T>(path, { method: 'DELETE', fetch: fetcher })
};
