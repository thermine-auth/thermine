<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { Badge, DataTable } from '$lib/components/ui';
	import { demoProviders } from '$lib/demo';
</script>

<svelte:head><title>Social · xermess admin</title></svelte:head>

<PageHeading
	title="Social providers"
	description="The accounts people can sign in with instead of a password."
	demo
/>

<DataTable
	rows={demoProviders}
	empty="No providers enabled."
	columns="minmax(7rem, 0.8fr) minmax(10rem, 1.6fr) auto auto"
>
	{#snippet row(provider)}
		<strong>{provider.name}</strong>
		<span class="hint mono url">{provider.clientId}</span>
		<span class="hint count">{provider.logins.toLocaleString()} logins</span>
		<Badge tone={provider.status === 'active' ? 'success' : 'neutral'}>{provider.status}</Badge>
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
