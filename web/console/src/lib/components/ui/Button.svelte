<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';
	import type { ControlProps } from './control';
	import Icon from './Icon.svelte';

	type Props = HTMLButtonAttributes &
		ControlProps & {
			children: Snippet;
		};

	let {
		size = 'md',
		variant = 'solid',
		colorPalette = 'neutral',
		loading = false,
		disabled = false,
		icon,
		type = 'button',
		children,
		...rest
	}: Props = $props();
</script>

<button
	{type}
	class="control"
	data-size={size}
	data-variant={variant}
	data-palette={colorPalette}
	data-loading={loading || undefined}
	disabled={disabled || loading}
	{...rest}
>
	{#if loading}
		<span class="spinner" aria-hidden="true"></span>
	{:else if icon}
		<Icon {icon} />
	{/if}

	{@render children()}
</button>
