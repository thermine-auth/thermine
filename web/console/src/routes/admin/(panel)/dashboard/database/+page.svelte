<script lang="ts">
	import { Badge, type Column, DataTable, PageHeading } from '$lib/components/ui';
	import { demoConnections } from '$lib/data/demo';

	const columns: Column[] = [
		{ key: 'name', min: '10rem' },
		{ key: 'engine' },
		{ key: 'users' },
		{ key: 'applications' },
		{ key: 'status' }
	];
</script>

<svelte:head><title>Database · xermess admin</title></svelte:head>

<PageHeading
	title="Database connections"
	description="Where accounts and their credentials are stored."
	demo
/>

<DataTable {columns} rows={demoConnections} empty="No connections.">
	{#snippet row(connection)}
		<td><strong>{connection.name}</strong></td>
		<td class="hint mono">{connection.engine}</td>
		<td class="hint count">{connection.users.toLocaleString()} users</td>
		<td class="hint count">{connection.applications} apps</td>
		<td>
			<Badge tone={connection.status === 'active' ? 'success' : 'neutral'}>
				{connection.status}
			</Badge>
		</td>
	{/snippet}
</DataTable>

<style>
	.hint {
		color: var(--color-text-hint);
	}
	.mono {
		font-family: var(--font-mono);
		font-size: var(--text-sm);
	}
	.count {
		font-size: var(--text-sm);
		white-space: nowrap;
	}
</style>
