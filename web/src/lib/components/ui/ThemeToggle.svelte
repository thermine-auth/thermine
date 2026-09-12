<script lang="ts">
	import { RiMoonLine, RiSunLine } from 'svelte-remixicon';
	import { theme } from '$lib/theme.svelte';
	import Icon from './Icon.svelte';

	const isDark = $derived(theme.current === 'dark');
</script>

<!-- The icon shows the theme you would switch to, which is what a reader
     expects from a control they are about to click. -->
<button
	type="button"
	onclick={() => theme.toggle()}
	title={isDark ? 'Switch to the light theme' : 'Switch to the dark theme'}
	aria-label={isDark ? 'Switch to the light theme' : 'Switch to the dark theme'}
	aria-pressed={isDark}
>
	<Icon icon={isDark ? RiSunLine : RiMoonLine} size="1.125rem" />
</button>

<style>
	button {
		display: grid;
		place-items: center;
		width: 35px;
		height: 35px;
		border: none;
		border-radius: var(--radius-md);
		background: transparent;
		color: inherit;
		cursor: pointer;
		transition: background-color var(--speed-fast);
	}

	button:hover {
		background: color-mix(in srgb, currentcolor, transparent 88%);
	}
</style>
