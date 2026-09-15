<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		title: string;
		/** A line under the title: what the section is for. */
		description?: string;
		/** Controls on the right of the title, such as an "Add" button. */
		action?: Snippet;
		children: Snippet;
	};

	let { title, description, action, children }: Props = $props();
</script>

<!-- One titled group of a form. Every panel is built from these, so sections
     are spaced, ruled and titled the same way wherever they appear. -->
<section class="form-section">
	<header>
		<div class="text">
			<h3>{title}</h3>
			{#if description}
				<p>{description}</p>
			{/if}
		</div>

		{#if action}
			<div class="action">{@render action()}</div>
		{/if}
	</header>

	<div class="content">
		{@render children()}
	</div>
</section>

<style>
	.form-section + :global(.form-section) {
		margin-top: var(--space-5);
		padding-top: var(--space-5);
		border-top: 1px solid var(--color-border);
	}

	header {
		display: flex;
		align-items: flex-start;
		justify-content: space-between;
		gap: var(--space-3);
		margin-bottom: var(--space-3);
	}

	.text {
		min-width: 0;
	}

	h3 {
		margin: 0;
		color: var(--color-text);
		font-size: var(--text-base);
		font-weight: 600;
	}

	p {
		margin: 2px 0 0;
		color: var(--color-text-hint);
		font-size: var(--text-sm);
		line-height: 1.45;
	}

	.action {
		display: flex;
		flex-shrink: 0;
		gap: var(--space-2);
	}

	.content {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}
</style>
