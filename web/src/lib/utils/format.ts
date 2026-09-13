/**
 * How values are shown to a reader.
 *
 * Dates arrive from the API as RFC 3339 strings and are shown in the reader's
 * own locale and time zone. Anything that is not a date is handed back as it
 * came: "Invalid Date" in a table tells nobody anything, while the value
 * itself at least says what was stored.
 */
export function formatDateTime(value: string): string {
	const date = new Date(value);

	return Number.isNaN(date.getTime()) ? value : date.toLocaleString();
}

/** The day alone, for a value whose time of day means nothing. */
export function formatDate(value: string): string {
	const date = new Date(value);

	return Number.isNaN(date.getTime()) ? value : date.toLocaleDateString();
}
