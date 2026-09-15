<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLAnchorAttributes } from 'svelte/elements';
	import type { ControlProps } from './control';
	import Icon from './Icon.svelte';

	type Props = HTMLAnchorAttributes &
		Omit<ControlProps, 'loading'> & {
			href: string;
			children: Snippet;
		};

	let {
		href,
		size = 'md',
		variant = 'solid',
		colorPalette = 'neutral',
		disabled = false,
		icon,
		children,
		...rest
	}: Props = $props();
</script>

<!-- Somewhere to go rather than something to do, so it is an anchor: it can
     be opened in a new tab, and the browser shows where it leads. A disabled
     link is not a thing, so one that is turned off loses its href and says so
     to a screen reader instead.

     The href belongs to whoever used this component, and it is their resolve()
     that decides it — this one only passes it on. -->
<!-- eslint-disable svelte/no-navigation-without-resolve -->
<a
	href={disabled ? undefined : href}
	class="control"
	data-size={size}
	data-variant={variant}
	data-palette={colorPalette}
	aria-disabled={disabled || undefined}
	role={disabled ? 'link' : undefined}
	{...rest}
>
	{#if icon}
		<Icon {icon} />
	{/if}

	{@render children()}
</a>

<!-- eslint-enable svelte/no-navigation-without-resolve -->
