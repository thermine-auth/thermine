<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { Badge, DataTable } from '$lib/components/ui';
	import { demoKeys } from '$lib/demo';
</script>

<svelte:head>
	<title>API keys · xermess admin</title>
</svelte:head>

<PageHeading
	title="API keys"
	description="Keys that let an application call this server on its own behalf."
	demo
/>

<DataTable
	rows={demoKeys}
	empty="No keys."
	columns="minmax(7rem, 1fr) minmax(9rem, 1fr) minmax(9rem, 1.2fr) auto auto"
>
	{#snippet row(key)}
		<strong>{key.label}</strong>
		<span class="hint mono">{key.prefix}</span>
		<span class="hint mono">{key.scopes}</span>
		<Badge tone={key.active ? 'success' : 'neutral'}>{key.active ? 'active' : 'revoked'}</Badge>
		<span class="hint when">{key.lastUsed}</span>
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

	.when {
		font-size: var(--text-sm);
		text-align: right;
		white-space: nowrap;
	}
</style>
