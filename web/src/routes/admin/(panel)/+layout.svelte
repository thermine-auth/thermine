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

	/* Pages sit in a centred column of a readable width. */
	main {
		max-width: var(--content-width);
		margin-inline: auto;
		padding: var(--space-5) var(--space-4);
	}

	/* The dashboard brings its own sidebar and centres its own content, so it
	   takes the full width instead. .dashboard belongs to a child component,
	   which is why :has has to match it globally. */
	main:has(:global(.dashboard)) {
		max-width: none;
		margin-inline: 0;
		padding: 0;
	}

	@media (max-width: 40rem) {
		main {
			padding: var(--space-4) var(--space-3);
		}
	}
</style>
