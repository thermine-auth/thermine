<script lang="ts">
	import type { Snippet } from 'svelte';
	import Sidebar from '$lib/components/admin/Sidebar.svelte';

	let { children }: { children: Snippet } = $props();
</script>

<div class="dashboard">
	<Sidebar />

	<div class="content">
		<div class="column">
			{@render children()}
		</div>
	</div>
</div>

<style>
	.dashboard {
		display: grid;
		grid-template-columns: var(--sidebar-width) 1fr;
		align-items: start;
	}

	.content {
		min-width: 0;
		padding: var(--space-5) var(--space-4);
	}

	/* Centred in the space left over by the sidebar, and no wider than any
	   other page. */
	.column {
		max-width: var(--content-width);
		margin-inline: auto;
	}

	@media (max-width: 55rem) {
		.dashboard {
			grid-template-columns: 1fr;
		}
	}

	@media (max-width: 40rem) {
		.content {
			padding: var(--space-4) var(--space-3);
		}
	}
</style>
