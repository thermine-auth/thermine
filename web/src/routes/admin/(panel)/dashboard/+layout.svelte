<script lang="ts">
	import type { Snippet } from 'svelte';
	import Sidebar from '$lib/components/layout/Sidebar.svelte';
	import { useShell } from '$lib/state/shell.svelte';

	let { children }: { children: Snippet } = $props();

	const shell = useShell();
</script>

<div class="dashboard">
	<Sidebar collapsed={shell.collapsed} onToggle={shell.toggle} />

	<div class="content">
		{@render children()}
	</div>
</div>

<style>
	/* The column is as wide as the panel layout says, the same value the
	   header's logo block uses, so the two line up while it folds. */
	.dashboard {
		display: grid;
		grid-template-columns: var(--sidebar-width) 1fr;
		align-items: start;
	}

	/* No side padding: the tables inside reach the sidebar and the window
	   edge, and everything else is inset by the gutter instead. */
	.content {
		min-width: 0;
		padding: var(--space-4) 0;
	}

	@media (max-width: 55rem) {
		.dashboard {
			grid-template-columns: 1fr;
		}
	}
</style>
