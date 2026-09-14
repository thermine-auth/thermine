<script lang="ts">
	import {
		RiAdminLine,
		RiAppsLine,
		RiComputerLine,
		RiHistoryLine,
		RiShieldUserLine
	} from 'svelte-remixicon';
	import ActivityTable from '$lib/components/activity/ActivityTable.svelte';
	import SessionTable from '$lib/components/profile/SessionTable.svelte';
	import StatGrid from '$lib/components/activity/StatGrid.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const stats = $derived(
		data.overview
			? [
					{ label: 'Administrators', value: data.overview.counts.admins, icon: RiAdminLine },
					{
						label: 'Active sessions',
						value: data.overview.counts.active_sessions,
						icon: RiComputerLine
					},
					{
						label: 'Applications',
						value: data.overview.counts.applications,
						icon: RiAppsLine
					},
					{ label: 'Admin roles', value: data.overview.counts.roles, icon: RiShieldUserLine },
					{ label: 'Logged events', value: data.overview.counts.events, icon: RiHistoryLine }
				]
			: []
	);
</script>

<svelte:head>
	<title>Activity · xermess admin</title>
</svelte:head>

<h1>Activity</h1>
<p class="subtitle">
	Signed in as {data.admin.full_name}{data.admin.roles.length > 0
		? ` (${data.admin.roles.join(', ')})`
		: ''}
</p>

{#if data.overview}
	<section>
		<StatGrid {stats} />
	</section>

	<section>
		<h2>Recent activity</h2>
		<ActivityTable events={data.overview.activity} />
	</section>
{:else}
	<p class="note">Your roles do not include reading activity, so there is nothing to show here.</p>
{/if}

<section>
	<h2>Your sessions</h2>
	<SessionTable sessions={data.sessions} />
</section>

<style>
	.subtitle {
		margin: var(--space-1) 0 var(--space-5);
		padding-inline: var(--page-gutter);
		color: var(--color-text-hint);
		font-size: var(--text-base);
	}

	.note {
		padding-inline: var(--page-gutter);
		color: var(--color-text-hint);
	}

	section + section {
		margin-top: var(--space-6);
	}

	h2 {
		margin-bottom: var(--space-3);
		padding-inline: var(--page-gutter);
	}
</style>
