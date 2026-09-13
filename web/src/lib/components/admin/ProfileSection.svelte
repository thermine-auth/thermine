<script lang="ts">
	import type { ComponentType, Snippet } from 'svelte';
	import { Badge, Icon } from '$lib/components/ui';

	type Props = {
		icon: ComponentType;
		title: string;
		description: string;
		/** Marks a setting that cannot be changed here yet. */
		demo?: boolean;
		children: Snippet;
	};

	let { icon, title, description, demo = false, children }: Props = $props();
</script>

<section>
	<header>
		<span class="icon"><Icon {icon} size="1.0625rem" /></span>
		<div class="text">
			<h2>
				{title}
				{#if demo}
					<Badge>not wired up</Badge>
				{/if}
			</h2>
			<p>{description}</p>
		</div>
	</header>

	{@render children()}
</section>

<style>
	section + :global(section) {
		margin-top: var(--space-5);
	}

	section {
		padding-inline: var(--page-gutter);
	}

	header {
		display: flex;
		align-items: flex-start;
		gap: var(--space-3);
		margin-bottom: var(--space-3);
	}

	.icon {
		display: grid;
		place-items: center;
		width: 32px;
		height: 32px;
		border-radius: var(--radius-md);
		background: var(--color-secondary);
		color: var(--color-text-hint);
	}

	h2 {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}

	p {
		margin-top: 2px;
		color: var(--color-text-hint);
		font-size: var(--text-sm);
	}
</style>
