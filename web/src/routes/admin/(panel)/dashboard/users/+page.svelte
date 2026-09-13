<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { Badge, DataTable } from '$lib/components/ui';
	import { demoUsers } from '$lib/demo';

	const tone = (status: string) =>
		status === 'active' ? 'success' : status === 'blocked' ? 'danger' : 'neutral';
</script>

<svelte:head><title>Users · xermess admin</title></svelte:head>

<PageHeading title="Users" description="Everyone with an account in this organisation." demo />

<DataTable
	rows={demoUsers}
	empty="No users yet."
	columns="minmax(8rem, 1.1fr) minmax(9rem, 1.2fr) minmax(8rem, 1fr) auto auto"
>
	{#snippet row(user)}
		<strong>{user.name}</strong>
		<span class="hint">{user.email}</span>
		<span class="hint mono">{user.connection}</span>
		<Badge tone={tone(user.status)}>{user.status}</Badge>
		<span class="hint when">{user.lastLogin}</span>
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
	@media (max-width: 40rem) {
		.when {
			text-align: left;
		}
	}
</style>
