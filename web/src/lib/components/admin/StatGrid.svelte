<script lang="ts">
	import type { ComponentType } from 'svelte';
	import { Card, Icon } from '$lib/components/ui';

	type Props = {
		stats: { label: string; value: number; icon: ComponentType }[];
	};

	let { stats }: Props = $props();
</script>

<div class="grid">
	{#each stats as stat (stat.label)}
		<Card padded>
			<div class="stat">
				<span class="label">
					<Icon icon={stat.icon} />
					{stat.label}
				</span>
				<strong>{stat.value}</strong>
			</div>
		</Card>
	{/each}
</div>

<style>
	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(9rem, 1fr));
		gap: var(--space-3);
	}

	.stat {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}

	.label {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		color: var(--color-text-hint);
		font-size: var(--text-sm);
	}

	.stat strong {
		font-size: var(--text-2xl);
		font-weight: 600;
	}
</style>
