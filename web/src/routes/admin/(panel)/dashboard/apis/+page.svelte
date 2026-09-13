<script lang="ts">
	import { Badge, type Column, DataTable, PageHeading } from '$lib/components/ui';
	import { demoApis } from '$lib/data/demo';

	const columns: Column[] = [
		{ key: 'name', min: '10rem' },
		{ key: 'identifier', min: '14rem' },
		{ key: 'scopes' },
		{ key: 'lifetime' },
		{ key: 'status' }
	];
</script>

<svelte:head><title>APIs · xermess admin</title></svelte:head>

<PageHeading
	title="APIs"
	description="The services applications ask for access tokens to reach."
	demo
/>

<DataTable {columns} rows={demoApis} empty="No APIs registered.">
	{#snippet row(api)}
		<td><strong>{api.name}</strong></td>
		<td class="hint mono url">{api.identifier}</td>
		<td class="hint count">{api.scopes} scopes</td>
		<td class="hint count">{api.tokenLifetime}</td>
		<td>
			<Badge tone={api.status === 'active' ? 'success' : 'neutral'}>{api.status}</Badge>
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
	.url {
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
