<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Dialog } from '@ark-ui/svelte/dialog';
	import { RiCloseLine } from 'svelte-remixicon';
	import IconButton from './IconButton.svelte';

	type Props = {
		open: boolean;
		title: string;
		/** A line under the title: what the panel is for. */
		description?: string;
		/** Shown next to the title, for a record's id and the like. */
		meta?: string;
		/** How wide the panel opens, as a CSS length. It never passes the edge
		    of a narrow screen, whatever this says. */
		width?: string;
		/** Given when the body is a form, so the footer's submit button
		    belongs to it and Enter saves. Without it the panel is a plain
		    container. */
		onsubmit?: (event: SubmitEvent) => void;
		children: Snippet;
		footer?: Snippet;
	};

	let {
		open = $bindable(false),
		title,
		description,
		meta,
		width,
		onsubmit,
		children,
		footer
	}: Props = $props();
</script>

<!-- The parts of the panel, rendered inside either a form or a plain div
     below so the markup is written once. -->
{#snippet panel()}
	<header>
		<div class="heading">
			<Dialog.Title>{title}</Dialog.Title>
			{#if meta}<span class="meta">{meta}</span>{/if}
		</div>

		{#if description}
			<Dialog.Description>{description}</Dialog.Description>
		{/if}

		<div class="close">
			<IconButton
				icon={RiCloseLine}
				label="Close"
				placement="left"
				onclick={() => (open = false)}
			/>
		</div>
	</header>

	<div class="body">
		{@render children()}
	</div>

	{#if footer}
		<footer>
			{@render footer()}
		</footer>
	{/if}
{/snippet}

<Dialog.Root bind:open lazyMount unmountOnExit>
	<Dialog.Backdrop />
	<Dialog.Positioner>
		<Dialog.Content style={width ? `--drawer-width: ${width}` : undefined}>
			{#if onsubmit}
				<form class="layout" {onsubmit}>{@render panel()}</form>
			{:else}
				<div class="layout">{@render panel()}</div>
			{/if}
		</Dialog.Content>
	</Dialog.Positioner>
</Dialog.Root>

<style>
	/* The panel is as tall as the window: only the body scrolls, so the title
	   and the buttons stay put however long the form is. */
	.layout {
		display: flex;
		flex-direction: column;
		min-height: 0;
		height: 100%;
	}

	header {
		position: relative;
		padding: var(--space-4) var(--space-5);
		padding-right: var(--space-6);
		border-bottom: 1px solid var(--color-border);
	}

	.heading {
		display: flex;
		align-items: baseline;
		gap: var(--space-2);
	}

	.meta {
		color: var(--color-text-hint);
		font-family: var(--font-mono);
		font-size: var(--text-sm);
	}

	.body {
		flex: 1;
		min-height: 0;
		padding: var(--space-5);
		overflow-y: auto;
	}

	footer {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-3) var(--space-5);
		border-top: 1px solid var(--color-border);
		background: var(--color-surface);
	}

	.close {
		position: absolute;
		top: var(--space-2);
		right: var(--space-2);
	}
</style>
