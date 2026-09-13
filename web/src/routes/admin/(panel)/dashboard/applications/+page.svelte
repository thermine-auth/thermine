<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { Badge, DataTable } from '$lib/components/ui';
	import { demoApplications } from '$lib/demo';
</script>

<svelte:head><title>Applications · xermess admin</title></svelte:head>

<PageHeading
	title="Applications"
	description="The apps that sign people in through xermess."
	demo
/>

<DataTable
	rows={demoApplications}
	empty="No applications yet."
	columns="minmax(9rem, 1.2fr) minmax(8rem, 1fr) minmax(7rem, 1fr) auto auto"
>
	{#snippet row(app)}
		<strong>{app.name}</strong>
		<span class="hint">{app.kind}</span>
		<span class="hint mono">{app.clientId}</span>
		<span class="hint count">{app.users.toLocaleString()} users</span>
		<Badge tone={app.status === 'active' ? 'success' : 'neutral'}>{app.status}</Badge>
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
