<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { Badge, DataTable } from '$lib/components/ui';
	import { demoConnections } from '$lib/demo';
</script>

<svelte:head><title>Database · xermess admin</title></svelte:head>

<PageHeading
	title="Database connections"
	description="Where accounts and their credentials are stored."
	demo
/>

<DataTable
	rows={demoConnections}
	empty="No connections."
	columns="minmax(9rem, 1.2fr) minmax(7rem, 1fr) auto auto auto"
>
	{#snippet row(connection)}
		<strong>{connection.name}</strong>
		<span class="hint mono">{connection.engine}</span>
		<span class="hint count">{connection.users.toLocaleString()} users</span>
		<span class="hint count">{connection.applications} apps</span>
		<Badge tone={connection.status === 'active' ? 'success' : 'neutral'}>{connection.status}</Badge>
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
