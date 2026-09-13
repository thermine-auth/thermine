import type { ComponentType } from 'svelte';

/**
 * The props every control shares.
 *
 * They map one-to-one onto the data attributes in styles/controls.css: a
 * component's job is to pass them through, not to decide what they look
 * like. Adding a size or a palette is a block of CSS there and one more
 * value here.
 */

/** How tall a control is. `md` is the default, and what a toolbar is built
 *  from; `sm` is for controls inside a row or a bar. */
export type Size = 'sm' | 'md' | 'lg';

/** How much of the palette a control uses: `solid` is filled and loud,
 *  `subtle` is filled and quiet, `outline` is a border, `ghost` is nothing
 *  until it is pointed at, `plain` is a link that happens to be a button. */
export type Variant = 'solid' | 'subtle' | 'outline' | 'ghost' | 'plain';

/** Which colours those variants use. See styles/palettes.css. */
export type ColorPalette = 'neutral' | 'danger' | 'success' | 'info' | 'warning';

export type ControlProps = {
	size?: Size;
	variant?: Variant;
	colorPalette?: ColorPalette;
	/** Waiting on something: the control shows a spinner and refuses clicks. */
	loading?: boolean;
	disabled?: boolean;
	/** Drawn before the label. */
	icon?: ComponentType;
};
