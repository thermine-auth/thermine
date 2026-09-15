<script lang="ts">
	import { Badge, type Column, DataTable, PageHeading } from '$lib/components/ui';
	import { demoProviders } from '$lib/data/demo';

	const columns: Column[] = [
		{ key: 'name', min: '9rem' },
		{ key: 'client', min: '14rem' },
		{ key: 'logins' },
		{ key: 'status' }
	];
</script>

<svelte:head><title>Social · xermess admin</title></svelte:head>

<PageHeading
	title="Social providers"
	description="The accounts people can sign in with instead of a password."
	demo
/>

<DataTable {columns} rows={demoProviders} empty="No providers enabled.">
	{#snippet row(provider)}
		<td><strong>{provider.name}</strong></td>
		<td class="hint mono url">{provider.clientId}</td>
		<td class="hint count">{provider.logins.toLocaleString()} logins</td>
		<td>
			<Badge tone={provider.status === 'active' ? 'success' : 'neutral'}>{provider.status}</Badge>
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
