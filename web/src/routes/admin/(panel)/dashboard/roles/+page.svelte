<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { DataTable, type Column } from '$lib/components/ui';
	import { demoRoles } from '$lib/demo';

	const columns: Column[] = [
		{ key: 'name', min: '9rem' },
		{ key: 'description', min: '14rem' },
		{ key: 'permissions' },
		{ key: 'members' }
	];
</script>

<svelte:head><title>Roles · xermess admin</title></svelte:head>

<PageHeading
	title="Roles"
	description="A role is a named bundle of permissions. People hold roles."
	demo
/>

<DataTable {columns} rows={demoRoles} empty="No roles.">
	{#snippet row(role)}
		<td class="mono name">{role.name}</td>
		<td class="hint">{role.description}</td>
		<td class="hint count">{role.permissions} permissions</td>
		<td class="hint count">{role.members} members</td>
	{/snippet}
</DataTable>

<style>
	.name {
		font-weight: 600;
	}
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
