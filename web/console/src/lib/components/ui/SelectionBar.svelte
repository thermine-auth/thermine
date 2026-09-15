<script lang="ts">
	import type { Snippet } from 'svelte';
	import { fly } from 'svelte/transition';

	import Button from './Button.svelte';

	type Props = {
		/** How many rows are ticked. The bar shows itself when this is
		    anything but zero. */
		count: number;
		/** What one row is called, e.g. 'record' or 'session'. */
		noun?: string;
		onReset: () => void;
		/** The buttons that act on the selection. */
		children: Snippet;
	};

	let { count, noun = 'record', onReset, children }: Props = $props();
</script>

{#if count > 0}
	<!-- The bar rides the bottom of the window while there is table left to
	     scroll, and comes to rest under the last row when there is not. -->
	<div class="bar" role="status" transition:fly={{ y: 12, duration: 150 }}>
		<span class="count">
			Selected <strong>{count}</strong>
			{count === 1 ? noun : `${noun}s`}
		</span>

		<Button size="sm" variant="subtle" onclick={onReset}>Reset</Button>

		<span class="actions">
			{@render children()}
		</span>
	</div>
{/if}

<style>
	.bar {
		position: sticky;
		bottom: var(--space-4);
		z-index: 20;
		display: flex;
		align-items: center;
		gap: var(--space-2);
		width: fit-content;
		max-width: calc(100% - var(--page-gutter) * 2);
		margin: var(--space-4) auto 0;
		padding: var(--space-2);
		padding-left: var(--space-4);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-pill);
		background: var(--color-surface);
		box-shadow: var(--shadow-md);
	}

	.count {
		color: var(--color-text-hint);
		white-space: nowrap;
	}

	.count strong {
		color: var(--color-text);
	}

	.actions {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		margin-left: var(--space-3);
	}

	/* Every button in the bar is rounded like the bar that holds it. The ones
	   in .actions come from the caller, hence the global match. */
	.bar :global(.control) {
		border-radius: var(--radius-pill);
	}
</style>
