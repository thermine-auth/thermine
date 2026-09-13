import type { ComponentType } from 'svelte';

/** One column of a DataTable. The same list drives the header and the
 *  widths, so adding a column is one entry here and one `<td>` in the row
 *  snippet. */
export type Column = {
	/** Identifies the column. Also the heading, unless `label` says otherwise. */
	key: string;
	/** The heading. A table whose columns carry no labels gets no header row. */
	label?: string;
	/** Drawn before the label, the way PocketBase marks a column's type. */
	icon?: ComponentType;
	/** The narrowest this column may become, e.g. '11rem'. Past the point
	 *  where every column is at its narrowest, the table scrolls sideways. */
	min?: string;
	/** Right-aligned columns read better for dates and counts. */
	align?: 'start' | 'end';
};
