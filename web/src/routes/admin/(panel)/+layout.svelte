<script lang="ts">
	import AppHeader from '$lib/components/admin/AppHeader.svelte';
	import type { LayoutData } from './$types';

	let { data, children }: { data: LayoutData; children: import('svelte').Snippet } = $props();
</script>

<div class="shell">
	<AppHeader admin={data.admin} />

	<main>
		{@render children()}
	</main>
</div>

<style>
	.shell {
		min-height: 100dvh;
	}

	main {
		max-width: 64rem;
		margin: 0 auto;
		padding: var(--space-6) var(--space-5);
	}
</style>
