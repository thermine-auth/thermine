<script lang="ts">
	import type { ComponentType } from 'svelte';
	import { Tooltip } from '@ark-ui/svelte/tooltip';
	import Icon from './Icon.svelte';

	type Props = {
		icon: ComponentType;
		/** What the button does. It is both the tooltip and the name the
		    button is read out by, so it is never left out. */
		label: string;
		/** `md` is the size of the other controls on a page; `sm` is for
		    buttons that sit inside a table row or a list. */
		size?: 'md' | 'sm';
		/** `danger` colours the hover for something that removes. */
		tone?: 'default' | 'danger';
		/** Turns the icon, for a button that is waiting on something. */
		spinning?: boolean;
		/** Which side the tooltip appears on. */
		placement?: 'top' | 'bottom' | 'left' | 'right';
		disabled?: boolean;
		onclick?: () => void;
	};

	let {
		icon,
		label,
		size = 'md',
		tone = 'default',
		spinning = false,
		placement = 'bottom',
		disabled = false,
		onclick
	}: Props = $props();

	const classes = $derived(
		['icon-button', size, tone, spinning ? 'spinning' : ''].filter(Boolean).join(' ')
	);
</script>

<!-- The trigger is the button itself rather than a wrapper around one: a
     button inside a button is not markup a browser will keep. -->
<Tooltip.Root
	openDelay={250}
	closeDelay={80}
	positioning={{ placement, gutter: 6, strategy: 'fixed' }}
>
	<Tooltip.Trigger class={classes} type="button" aria-label={label} {disabled} {onclick}>
		<Icon {icon} size={size === 'sm' ? '1rem' : '1.125rem'} />
	</Tooltip.Trigger>

	<Tooltip.Positioner>
		<Tooltip.Content>{label}</Tooltip.Content>
	</Tooltip.Positioner>
</Tooltip.Root>
