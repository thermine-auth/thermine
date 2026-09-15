<script lang="ts">
	import { Badge, type Column, DataTable, PageHeading } from '$lib/components/ui';
	import { demoFlows } from '$lib/data/demo';

	const columns: Column[] = [
		{ key: 'name', min: '10rem' },
		{ key: 'steps', min: '14rem' },
		{ key: 'applications' },
		{ key: 'status' }
	];
</script>

<svelte:head><title>Login flows · xermess admin</title></svelte:head>

<PageHeading
	title="Login flows"
	description="The steps someone is taken through when they sign in."
	demo
/>

<DataTable {columns} rows={demoFlows} empty="No flows defined.">
	{#snippet row(flow)}
		<td>
			<span class="name">
				<strong>{flow.name}</strong>
				{#if flow.isDefault}
					<Badge>default</Badge>
				{/if}
			</span>
		</td>
		<td class="hint mono steps">{flow.steps}</td>
		<td class="hint count">{flow.applications} apps</td>
		<td>
			<Badge tone={flow.status === 'active' ? 'success' : 'neutral'}>{flow.status}</Badge>
		</td>
	{/snippet}
</DataTable>

<style>
	.name {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}
	.hint {
		color: var(--color-text-hint);
	}
	.mono {
		font-family: var(--font-mono);
		font-size: var(--text-sm);
	}
	.steps {
		max-width: 24rem;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	.count {
		font-size: var(--text-sm);
		white-space: nowrap;
	}
</style>
