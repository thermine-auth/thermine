/** Dates arrive from the API as RFC 3339 strings and are shown in the
    reader's own locale and time zone. */
export function formatDateTime(value: string): string {
	return new Date(value).toLocaleString();
}
