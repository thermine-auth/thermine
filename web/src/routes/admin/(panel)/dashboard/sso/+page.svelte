<script lang="ts">
	import { Badge, type Column, DataTable, PageHeading } from '$lib/components/ui';
	import { demoSso } from '$lib/data/demo';

	const columns: Column[] = [
		{ key: 'name', min: '10rem' },
		{ key: 'protocol' },
		{ key: 'domain' },
		{ key: 'users' },
		{ key: 'status' }
	];
</script>

<svelte:head><title>SSO integrations · xermess admin</title></svelte:head>

<PageHeading
	title="SSO integrations"
	description="Organisations whose people sign in through their own identity provider."
	demo
/>

<DataTable {columns} rows={demoSso} empty="No integrations yet.">
	{#snippet row(sso)}
		<td><strong>{sso.name}</strong></td>
		<td class="hint count">{sso.protocol}</td>
		<td class="hint mono">{sso.domain}</td>
		<td class="hint count">{sso.users.toLocaleString()} users</td>
		<td>
			<Badge tone={sso.status === 'active' ? 'success' : 'neutral'}>{sso.status}</Badge>
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
