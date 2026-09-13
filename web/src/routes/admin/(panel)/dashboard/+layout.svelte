<script lang="ts">
	import { untrack, type Snippet } from 'svelte';
	import Sidebar from '$lib/components/layout/Sidebar.svelte';
	import { rememberSidebar } from '$lib/state/sidebar';
	import type { LayoutData } from './$types';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();

	/** Starts as the server rendered it and is the reader's from then on, so
	    the value is read once on purpose: a later load returns the same
	    cookie, and re-reading it would undo a toggle. */
	let collapsed = $state(untrack(() => data.sidebar) === 'mini');

	function toggle() {
		collapsed = !collapsed;
		rememberSidebar(collapsed ? 'mini' : 'wide');
	}
</script>

<div class="dashboard" class:mini={collapsed}>
	<Sidebar {collapsed} onToggle={toggle} />

	<div class="content">
		{@render children()}
	</div>
</div>

<style>
	.dashboard {
		display: grid;
		grid-template-columns: var(--sidebar-width) 1fr;
		align-items: start;
		transition: grid-template-columns var(--speed);
	}

	/* Folded, the column is just wide enough for the icons. The width is a
	   variable so the sidebar's own sticky layout follows it. */
	.dashboard.mini {
		--sidebar-width: 3.5rem;
	}

	/* No side padding: the tables inside reach the sidebar and the window
	   edge, and everything else is inset by the gutter instead. */
	.content {
		min-width: 0;
		padding: var(--space-4) 0;
	}

	@media (prefers-reduced-motion: reduce) {
		.dashboard {
			transition: none;
		}
	}

	@media (max-width: 55rem) {
		.dashboard,
		.dashboard.mini {
			grid-template-columns: 1fr;
		}
	}
</style>
