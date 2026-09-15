<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { account, messageOf } from '$lib/api';
	import { Alert, AppMark, Button, Icon, Panel } from '$lib/components';
	import { formatDate, timeAgo } from '$lib/utils/format';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	/** What each scope lets an application do, said plainly. */
	const scopeNames: Record<string, string> = {
		openid: 'Sign you in',
		profile: 'See your name',
		email: 'See your email address',
		offline_access: 'Stay signed in when you are away',
		roles: 'See your roles'
	};

	let confirming = $state<string | null>(null);
	let disconnecting = $state<string | null>(null);
	let error = $state('');
	let notice = $state('');

	async function disconnect(clientId: string, name: string) {
		disconnecting = clientId;
		error = '';
		notice = '';

		try {
			await account.disconnect(clientId);
			notice = `${name} can no longer act for you. You’ll be asked to sign in the next time you use it.`;
			confirming = null;
			await invalidateAll();
		} catch (err) {
			error = messageOf(err);
		} finally {
			disconnecting = null;
		}
	}
</script>

<svelte:head>
	<title>Connected apps · Account</title>
</svelte:head>

<div class="page">
	<div class="heading">
		<h1>Connected apps</h1>
		<p>Applications that can keep you signed in and act for you while you are away.</p>
	</div>

	{#if error}<Alert>{error}</Alert>{/if}
	{#if notice}<Alert tone="success">{notice}</Alert>{/if}

	{#if data.applications.length === 0}
		<Panel title="No connected apps">
			<p class="empty">
				No application is keeping you signed in right now. Applications appear here when you allow
				them to stay signed in.
			</p>
		</Panel>
	{:else}
		{#each data.applications as app (app.client_id)}
			<Panel title={app.name} description="Last used {timeAgo(app.last_used_at)}">
				{#snippet aside()}
					<AppMark name={app.name} logo={app.logo_uri} size={44} />
				{/snippet}

				<div class="access">
					<h3>It can</h3>
					<ul>
						{#each app.scopes as scope (scope)}
							<li><Icon name="check" size="1rem" /> {scopeNames[scope] ?? scope}</li>
						{/each}
					</ul>
				</div>

				<p class="meta">
					Connected since {formatDate(app.authorized_at)}
					{#if app.client_uri}
						·
						<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
						<a href={app.client_uri} target="_blank" rel="noopener noreferrer">
							Visit <Icon name="external" size="0.8rem" />
						</a>
					{/if}
				</p>

				{#snippet footer()}
					{#if confirming === app.client_id}
						<span class="confirm">Disconnect {app.name}?</span>
						<Button variant="secondary" size="sm" onclick={() => (confirming = null)}>Keep</Button>
						<Button
							variant="danger"
							size="sm"
							loading={disconnecting === app.client_id}
							onclick={() => disconnect(app.client_id, app.name)}>Disconnect</Button
						>
					{:else}
						<Button variant="danger" size="sm" onclick={() => (confirming = app.client_id)}
							>Disconnect</Button
						>
					{/if}
				{/snippet}
			</Panel>
		{/each}
	{/if}
</div>

<style>
	.page {
		display: flex;
		flex-direction: column;
		gap: var(--space-5);
	}

	h1 {
		font-size: var(--text-2xl);
	}

	.heading p {
		margin-top: var(--space-1);
		color: var(--color-text-hint);
	}

	.empty {
		color: var(--color-text-hint);
	}

	h3 {
		margin: 0 0 var(--space-2);
		font-size: var(--text-sm);
		color: var(--color-text-hint);
		font-weight: 600;
	}

	ul {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(14rem, 1fr));
		gap: var(--space-2);
		margin: 0;
		padding: 0;
		list-style: none;
	}

	li {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}

	li :global(svg) {
		color: var(--color-success);
	}

	.meta {
		color: var(--color-text-hint);
		font-size: var(--text-sm);
	}

	.meta a {
		display: inline-flex;
		align-items: center;
		gap: 2px;
	}

	.confirm {
		margin-right: auto;
		align-self: center;
		font-weight: 600;
	}
</style>
