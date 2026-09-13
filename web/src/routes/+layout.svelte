<script lang="ts">
	import type { Snippet } from 'svelte';
	import favicon from '$lib/assets/favicon.svg';
	import '$lib/styles/app.css';

	let { children }: { children: Snippet } = $props();
</script>

<!-- No <title> here on purpose. A title in this layout is only applied once
     the app has booted, so it lands between the one the browser already has
     and the one the page sets, which shows as a second flicker. The default
     lives in app.html instead, where it is in the very first byte. -->
<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children()}
