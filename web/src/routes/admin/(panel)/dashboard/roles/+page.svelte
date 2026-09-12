<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { DataTable } from '$lib/components/ui';
	import { demoRoles } from '$lib/demo';
</script>

<svelte:head>
	<title>Roles · xermess admin</title>
</svelte:head>

<PageHeading
	title="Roles"
	description="A role is a named bundle of permissions. Administrators hold roles."
	demo
/>

<DataTable
	rows={demoRoles}
	empty="No roles."
	columns="minmax(8rem, 0.8fr) minmax(12rem, 2fr) auto auto"
>
	{#snippet row(role)}
		<span class="mono name">{role.name}</span>
		<span class="hint">{role.description}</span>
		<span class="hint count">{role.permissions} permissions</span>
		<span class="hint count">{role.admins} admins</span>
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
