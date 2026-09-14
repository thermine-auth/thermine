<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		/** Where the page sits, outermost first; the last is the page itself. */
		crumbs: string[];
		/** Small controls right beside the title, such as a refresh button. */
		secondary?: Snippet;
		/** The page's main actions, at the far end. */
		actions?: Snippet;
	};

	let { crumbs, secondary, actions }: Props = $props();
</script>

<!-- PocketBase's page header: breadcrumbs, the secondary buttons right beside
     them, and the primary ones pushed to the far side. The last crumb is the
     page's title, so it is the heading. -->
<header class="page-header">
	<nav class="breadcrumbs" aria-label="Breadcrumb">
		{#each crumbs.slice(0, -1) as crumb, index (index)}
			<span class="crumb">{crumb}</span>
		{/each}
		<h1 class="crumb current" aria-current="page">{crumbs.at(-1)}</h1>
	</nav>

	{#if secondary}<div class="secondary">{@render secondary()}</div>{/if}
	{#if actions}<div class="actions">{@render actions()}</div>{/if}
</header>

<style>
	.page-header {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 10px var(--space-4);
		min-height: var(--control-height);
	}

	.breadcrumbs {
		display: inline-flex;
		align-items: center;
		gap: 30px;
		min-width: 0;
		color: var(--color-text-hint);
		font-size: 1.286rem;
	}

	.crumb {
		position: relative;
		margin: 0;
		font-size: inherit;
		font-weight: normal;
		white-space: nowrap;
	}

	.crumb:not(.current)::after {
		position: absolute;
		top: 0;
		right: -18px;
		height: 100%;
		align-content: center;
		color: var(--color-text-disabled);
		font-size: 0.85em;
		content: '/';
		pointer-events: none;
	}

	.current {
		overflow: hidden;
		color: var(--color-text);
		text-overflow: ellipsis;
	}

	.secondary {
		display: inline-flex;
		align-items: center;
		gap: 10px;
	}

	.actions {
		display: inline-flex;
		align-items: center;
		gap: 10px;
		margin-left: auto;
	}

	/* A narrow screen has room for the page's own name only. */
	@media (max-width: 34rem) {
		.crumb:not(.current) {
			display: none;
		}
	}
</style>
