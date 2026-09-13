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

	/* Pages fill the width, the way PocketBase does: a table is easier to
	   read with room for its columns than centred in a narrow column, and it
	   runs to the edges rather than sitting in a box. */
	main {
		padding: var(--space-4) 0;
	}

	/* The dashboard brings its own sidebar and pads its own content.
	   .dashboard belongs to a child component, which is why :has has to match
	   it globally. */
	main:has(:global(.dashboard)) {
		padding: 0;
	}
</style>
