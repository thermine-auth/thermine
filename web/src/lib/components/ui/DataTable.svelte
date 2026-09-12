<script lang="ts" generics="Row extends { id: string }">
	import type { Snippet } from 'svelte';
	import Card from './Card.svelte';

	type Props = {
		rows: Row[];
		/** Shown instead of the rows when there are none. */
		empty: string;
		/** Renders one row. The grid columns are set by the caller. */
		row: Snippet<[Row]>;
		/** A CSS grid-template-columns value for the row layout. */
		columns: string;
	};

	let { rows, empty, row, columns }: Props = $props();
</script>

<Card>
	{#if rows.length}
		<div class="table" style="--columns: {columns}">
			{#each rows as item (item.id)}
				<div class="row">
					{@render row(item)}
				</div>
			{/each}
		</div>
	{:else}
		<p class="empty">{empty}</p>
	{/if}
</Card>

<style>
	.row {
		display: grid;
		grid-template-columns: var(--columns);
		gap: var(--space-4);
		align-items: center;
		padding: var(--space-3) var(--space-4);
		font-size: var(--text-base);
	}

	.row + .row {
		border-top: 1px solid var(--color-border);
	}

	/* Grid items default to min-width:auto, so a long value refuses to shrink
	   and spills into the next column. The cells come from a snippet, hence
	   the global match. */
	.row > :global(*) {
		min-width: 0;
	}

	.empty {
		padding: var(--space-5) var(--space-4);
		color: var(--color-text-hint);
		font-size: var(--text-base);
	}

	@media (max-width: 40rem) {
		.row {
			grid-template-columns: 1fr 1fr;
			gap: var(--space-2);
		}
	}
</style>
