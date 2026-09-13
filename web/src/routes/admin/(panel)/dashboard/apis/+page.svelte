<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { Badge, DataTable } from '$lib/components/ui';
	import { demoApis } from '$lib/demo';
</script>

<svelte:head><title>APIs · xermess admin</title></svelte:head>

<PageHeading
	title="APIs"
	description="The services applications ask for access tokens to reach."
	demo
/>

<DataTable
	rows={demoApis}
	empty="No APIs registered."
	columns="minmax(8rem, 1fr) minmax(12rem, 1.6fr) auto auto auto"
>
	{#snippet row(api)}
		<strong>{api.name}</strong>
		<span class="hint mono url">{api.identifier}</span>
		<span class="hint count">{api.scopes} scopes</span>
		<span class="hint count">{api.tokenLifetime}</span>
		<Badge tone={api.status === 'active' ? 'success' : 'neutral'}>{api.status}</Badge>
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
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	.count {
		font-size: var(--text-sm);
		white-space: nowrap;
	}
</style>
