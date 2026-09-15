import type { ComponentType } from 'svelte';
import {
	RiCalendarLine,
	RiHashtag,
	RiKey2Line,
	RiMailLine,
	RiText,
	RiToggleLine
} from 'svelte-remixicon';
import type { FieldType } from '$lib/api';

/** The icon that marks a column's type in the table header, the way
 *  PocketBase marks its own: the header says what kind of value is below it
 *  without spending a word on it. */
export const fieldIcons: Record<FieldType, ComponentType> = {
	text: RiText,
	number: RiHashtag,
	bool: RiToggleLine,
	email: RiMailLine,
	date: RiCalendarLine
};

/** Used for the id column, which is not a user-defined field. */
export const idIcon = RiKey2Line;
