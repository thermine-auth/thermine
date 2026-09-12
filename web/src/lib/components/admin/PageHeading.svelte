<script lang="ts">
	import { Badge } from '$lib/components/ui';

	type Props = {
		title: string;
		description: string;
		/** Marks a section whose contents are still placeholder data. */
		demo?: boolean;
	};

	let { title, description, demo = false }: Props = $props();
</script>

<header>
	<div class="title">
		<h1>{title}</h1>
		{#if demo}
			<Badge>placeholder data</Badge>
		{/if}
	</div>
	<p>{description}</p>
</header>

<style>
	header {
		margin-bottom: var(--space-4);
	}

	.title {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}

	p {
		margin-top: var(--space-1);
		color: var(--color-text-hint);
		font-size: var(--text-base);
	}
</style>
