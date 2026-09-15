<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		/** Something at the start of the row: usually a Thumb. */
		lead?: Snippet;
		/** The row's title. Left out, the children are the whole content. */
		title?: string;
		/** A line under the title. */
		description?: string;
		/** Things at the end of the row: tags, buttons. */
		end?: Snippet;
		/** Whatever else the row holds, under the title and description. */
		children?: Snippet;
	};

	let { lead, title, description, end, children }: Props = $props();
</script>

<li class="item">
	{#if lead}<span class="lead">{@render lead()}</span>{/if}

	<div class="content">
		{#if title}<span class="title">{title}</span>{/if}
		{#if description}<span class="description">{description}</span>{/if}
		{@render children?.()}
	</div>

	{#if end}<span class="end">{@render end()}</span>{/if}
</li>

<style>
	/* PocketBase's .list-item: at least 54px tall, ruled off above by a shadow
	   rather than a border, so the rule never doubles a box's own. */
	.item {
		display: flex;
		align-items: center;
		gap: 10px;
		min-height: 54px;
		padding: 10px var(--list-inset, var(--space-4));
		overflow-wrap: anywhere;
	}

	.item:not(:first-child) {
		box-shadow: 0 -1px 0 0 var(--color-secondary);
	}

	.lead,
	.end {
		display: inline-flex;
		flex-shrink: 0;
		align-items: center;
		gap: 5px;
	}

	.content {
		display: flex;
		flex: 1;
		flex-direction: column;
		min-width: 0;
		line-height: 1.5;
	}

	.title {
		font-size: var(--text-base);
	}

	.description {
		color: var(--color-text-hint);
		font-size: var(--text-sm);
	}

	@media (max-width: 34rem) {
		.item {
			flex-wrap: wrap;
		}

		.end {
			margin-left: auto;
		}
	}
</style>
