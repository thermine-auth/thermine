<script lang="ts">
	import { RiAdminLine, RiComputerLine, RiHistoryLine, RiShieldUserLine } from 'svelte-remixicon';
	import ActivityTable from '$lib/components/admin/ActivityTable.svelte';
	import SessionTable from '$lib/components/admin/SessionTable.svelte';
	import StatGrid from '$lib/components/admin/StatGrid.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const stats = $derived([
		{ label: 'Administrators', value: data.overview.counts.admins, icon: RiAdminLine },
		{
			label: 'Active sessions',
			value: data.overview.counts.active_sessions,
			icon: RiComputerLine
		},
		{ label: 'Roles', value: data.overview.counts.roles, icon: RiShieldUserLine },
		{ label: 'Logged events', value: data.overview.counts.events, icon: RiHistoryLine }
	]);
</script>

<svelte:head>
	<title>Overview · xermess admin</title>
</svelte:head>

<h1>Overview</h1>
<p class="subtitle">
	Signed in as {data.admin.full_name} ({data.admin.roles.join(', ')})
</p>

<section>
	<StatGrid {stats} />
</section>

<section>
	<h2>Recent activity</h2>
	<ActivityTable events={data.overview.activity} />
</section>

<section>
	<h2>Your sessions</h2>
	<SessionTable sessions={data.sessions} />
</section>

<style>
	.subtitle {
		margin: var(--space-1) 0 var(--space-5);
		color: var(--color-text-hint);
		font-size: var(--text-base);
	}

	section + section {
		margin-top: var(--space-6);
	}

	h2 {
		margin-bottom: var(--space-3);
	}
</style>
