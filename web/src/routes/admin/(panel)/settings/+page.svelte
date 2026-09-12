<script lang="ts">
	import SessionTable from '$lib/components/admin/SessionTable.svelte';
	import { Badge, Card } from '$lib/components/ui';
	import { formatDateTime } from '$lib/format';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const details = $derived([
		{ label: 'Username', value: data.admin.username },
		{ label: 'Email', value: data.admin.email },
		{ label: 'Name', value: data.admin.full_name },
		{ label: 'Roles', value: data.admin.roles.join(', ') },
		{
			label: 'Last signed in',
			value: data.admin.last_login_at ? formatDateTime(data.admin.last_login_at) : 'Never'
		}
	]);
</script>

<svelte:head>
	<title>Settings · xermess admin</title>
</svelte:head>

<header>
	<h1>Settings</h1>
	<p class="subtitle">Your account and where it is signed in.</p>
</header>

<section>
	<h2>Account</h2>
	<Card>
		<dl>
			{#each details as detail (detail.label)}
				<div class="row">
					<dt>{detail.label}</dt>
					<dd>{detail.value}</dd>
				</div>
			{/each}
			<div class="row">
				<dt>Status</dt>
				<dd>
					<Badge tone={data.admin.status === 'active' ? 'success' : 'neutral'}>
						{data.admin.status}
					</Badge>
				</dd>
			</div>
		</dl>
	</Card>
</section>

<section>
	<h2>Your sessions</h2>
	<SessionTable sessions={data.sessions} />
</section>

<style>
	header {
		margin-bottom: var(--space-5);
	}

	.subtitle {
		margin-top: var(--space-1);
		color: var(--color-text-hint);
		font-size: var(--text-base);
	}

	section + section {
		margin-top: var(--space-5);
	}

	h2 {
		margin-bottom: var(--space-3);
	}

	dl {
		margin: 0;
	}

	.row {
		display: grid;
		grid-template-columns: 10rem 1fr;
		gap: var(--space-4);
		align-items: center;
		padding: var(--space-3) var(--space-4);
		font-size: var(--text-base);
	}

	.row + .row {
		border-top: 1px solid var(--color-border);
	}

	dt {
		color: var(--color-text-hint);
	}

	dd {
		margin: 0;
	}

	@media (max-width: 40rem) {
		.row {
			grid-template-columns: 1fr;
			gap: var(--space-1);
		}
	}
</style>
