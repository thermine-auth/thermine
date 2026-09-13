<script lang="ts">
	import { Tooltip } from '@ark-ui/svelte/tooltip';
	import type { ComponentType } from 'svelte';
	import type { ControlProps, Size } from './control';
	import Icon from './Icon.svelte';

	type Props = Omit<ControlProps, 'icon'> & {
		icon: ComponentType;
		/** What the button does. It is both the tooltip and the name the
		    button is read out by, so it is never left out. */
		label: string;
		/** Which side the tooltip appears on. */
		placement?: 'top' | 'bottom' | 'left' | 'right';
		onclick?: () => void;
	};

	let {
		icon,
		label,
		size = 'md',
		variant = 'ghost',
		colorPalette = 'neutral',
		loading = false,
		disabled = false,
		placement = 'bottom',
		onclick
	}: Props = $props();

	/** The glyph is a share of the button, so every size stays in proportion. */
	const glyph: Record<Size, string> = { sm: '1rem', md: '1.125rem', lg: '1.25rem' };
</script>

<!-- The trigger is the button itself rather than a wrapper around one: a
     button inside a button is not markup a browser will keep. -->
<Tooltip.Root
	openDelay={250}
	closeDelay={80}
	positioning={{ placement, gutter: 6, strategy: 'fixed' }}
>
	<Tooltip.Trigger
		type="button"
		class="control"
		data-icon="true"
		data-size={size}
		data-variant={variant}
		data-palette={colorPalette}
		data-loading={loading || undefined}
		aria-label={label}
		disabled={disabled || loading}
		{onclick}
	>
		{#if loading}
			<span class="spinner" aria-hidden="true"></span>
		{:else}
			<Icon {icon} size={glyph[size]} />
		{/if}
	</Tooltip.Trigger>

	<Tooltip.Positioner>
		<Tooltip.Content>{label}</Tooltip.Content>
	</Tooltip.Positioner>
</Tooltip.Root>
