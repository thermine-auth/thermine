import type { ComponentType } from 'svelte';

/** One thing a Select can offer.
 *
 *  Only `value` is required, and `options` takes bare strings as well: a
 *  short list of keywords needs nothing more than itself. The rest is for the
 *  lists that carry more — an icon saying what kind of thing this is, a
 *  second line saying what choosing it means, the way PocketBase's field type
 *  picker does.
 */
export type SelectOption<Value extends string = string> = {
	value: Value;
	/** The text on screen, and what typing a few letters matches against.
	 *  Defaults to the value. */
	label?: string;
	/** A quieter second line under the label. */
	description?: string;
	icon?: ComponentType;
	/** Shown, but not open to being chosen. */
	disabled?: boolean;
};
