<script lang="ts">
	import type { ComponentType } from 'svelte';

	type Props = {
		/** A component from svelte-remixicon, e.g. RiLogoutBoxRLine. */
		icon: ComponentType;
		size?: string;
		/** Describe the icon when it carries meaning on its own; leave it out
		    when neighbouring text already says the same thing, so a screen
		    reader does not read it twice. */
		label?: string;
	};

	let { icon: Glyph, size = '1rem', label }: Props = $props();
</script>

<Glyph {size} {label} class="icon" />

<style>
	/* Icons sit on the text baseline rather than below it. */
	:global(.icon) {
		flex: none;
		vertical-align: -0.125em;
	}
</style>
