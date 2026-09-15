/** A failed request. `status` is 0 when the server could not be reached. */
export class ApiError extends Error {
	constructor(
		readonly status: number,
		message: string
	) {
		super(message);
		this.name = 'ApiError';
	}
}

export type Fetch = typeof globalThis.fetch;

type Options = {
	method?: 'GET' | 'POST' | 'PATCH' | 'DELETE';
	body?: unknown;
	fetch?: Fetch;
};

/**
 * Calls the xermess API, at /api/v1 on this app's own origin. From the browser
 * the proxy routes it to the API; from a server load, pass the load's `fetch`,
 * and hooks.server.ts sends it to the API directly with the reader's cookie.
 */
export async function request<T>(path: string, options: Options = {}): Promise<T> {
	const { method = 'GET', body, fetch: fetcher = globalThis.fetch } = options;

	const headers: Record<string, string> = {};
	if (body !== undefined) headers['Content-Type'] = 'application/json';

	let response: Response;
	try {
		response = await fetcher(`/api/v1${path}`, {
			method,
			credentials: 'same-origin',
			headers,
			body: body === undefined ? undefined : JSON.stringify(body)
		});
	} catch {
		throw new ApiError(0, 'Could not reach the server. Check your connection and try again.');
	}

	const payload = response.status === 204 ? {} : await response.json().catch(() => ({}));

	if (!response.ok) {
		throw new ApiError(response.status, payload.error ?? 'Something went wrong. Try again.');
	}

	return payload as T;
}

/** The message to show for an error thrown by a call. */
export function messageOf(err: unknown): string {
	return err instanceof ApiError ? err.message : 'Something went wrong. Try again.';
}
