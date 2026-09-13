<script lang="ts">
	import { Badge, type Column, DataTable, PageHeading } from '$lib/components/ui';
	import { demoApplications } from '$lib/data/demo';

	const columns: Column[] = [
		{ key: 'name', min: '10rem' },
		{ key: 'kind' },
		{ key: 'client' },
		{ key: 'users' },
		{ key: 'status' }
	];
</script>

<svelte:head><title>Applications · xermess admin</title></svelte:head>

<PageHeading
	title="Applications"
	description="The apps that sign people in through xermess."
	demo
/>

<DataTable {columns} rows={demoApplications} empty="No applications yet.">
	{#snippet row(app)}
		<td><strong>{app.name}</strong></td>
		<td class="hint">{app.kind}</td>
		<td class="hint mono">{app.clientId}</td>
		<td class="hint count">{app.users.toLocaleString()} users</td>
		<td>
			<Badge tone={app.status === 'active' ? 'success' : 'neutral'}>{app.status}</Badge>
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
