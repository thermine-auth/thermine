<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { Badge, DataTable } from '$lib/components/ui';
	import { demoSso } from '$lib/demo';
</script>

<svelte:head><title>SSO integrations · xermess admin</title></svelte:head>

<PageHeading
	title="SSO integrations"
	description="Organisations whose people sign in through their own identity provider."
	demo
/>

<DataTable
	rows={demoSso}
	empty="No integrations yet."
	columns="minmax(8rem, 1fr) auto minmax(7rem, 1fr) auto auto"
>
	{#snippet row(sso)}
		<strong>{sso.name}</strong>
		<span class="hint count">{sso.protocol}</span>
		<span class="hint mono">{sso.domain}</span>
		<span class="hint count">{sso.users.toLocaleString()} users</span>
		<Badge tone={sso.status === 'active' ? 'success' : 'neutral'}>{sso.status}</Badge>
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
