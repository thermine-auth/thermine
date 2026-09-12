<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { Badge, DataTable } from '$lib/components/ui';
	import { demoAdmins } from '$lib/demo';

	const tone = (status: string) =>
		status === 'active' ? 'success' : status === 'suspended' ? 'danger' : 'neutral';
</script>

<svelte:head>
	<title>Administrators · xermess admin</title>
</svelte:head>

<PageHeading title="Administrators" description="Everyone who can sign in to this panel." demo />

<DataTable
	rows={demoAdmins}
	empty="No administrators."
	columns="minmax(9rem, 1.4fr) minmax(8rem, 1.4fr) minmax(6rem, 0.8fr) auto auto"
>
	{#snippet row(admin)}
		<span class="who">
			<strong>{admin.name}</strong>
			<span class="hint mono">{admin.username}</span>
		</span>
		<span class="hint">{admin.email}</span>
		<span class="hint mono">{admin.role}</span>
		<Badge tone={tone(admin.status)}>{admin.status}</Badge>
		<span class="hint when">{admin.lastSeen}</span>
	{/snippet}
</DataTable>

<style>
	.who {
		display: flex;
		flex-direction: column;
	}

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
