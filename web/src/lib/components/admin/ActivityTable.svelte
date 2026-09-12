<script lang="ts">
	import type { ComponentType } from 'svelte';
	import {
		RiErrorWarningLine,
		RiLockLine,
		RiLoginBoxLine,
		RiLogoutBoxRLine,
		RiPulseLine
	} from 'svelte-remixicon';
	import type { ActivityEvent } from '$lib/api';
	import { DataTable, Icon } from '$lib/components/ui';
	import { formatDateTime } from '$lib/format';

	type Props = { events: ActivityEvent[] };

	let { events }: Props = $props();

	/** The icon for an action, matching the names the server records in
	    audit_logs. Anything unrecognised gets the neutral one. */
	const icons: Record<string, ComponentType> = {
		'admin.login': RiLoginBoxLine,
		'admin.logout': RiLogoutBoxRLine,
		'admin.login_failed': RiErrorWarningLine,
		'admin.login_blocked': RiLockLine
	};

	function iconFor(action: string): ComponentType {
		return icons[action] ?? RiPulseLine;
	}

	/** Failed and blocked sign-ins are the rows worth noticing. */
	function isFailure(action: string): boolean {
		return action.endsWith('_failed') || action.endsWith('_blocked');
	}
</script>

<DataTable
	rows={events}
	empty="Nothing has happened yet."
	columns="minmax(9rem, 1fr) minmax(5rem, 1fr) minmax(6rem, 1fr) auto"
>
	{#snippet row(event)}
		<span class="action" class:failure={isFailure(event.action)}>
			<Icon icon={iconFor(event.action)} />
			{event.action}
		</span>
		<span class="muted">{event.actor || '—'}</span>
		<span class="muted mono">{event.ip}</span>
		<span class="muted when">{formatDateTime(event.created_at)}</span>
	{/snippet}
</DataTable>

<style>
	.action {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		font-family: var(--font-mono);
		font-size: var(--text-sm);
		font-weight: 500;
	}

	.failure {
		color: var(--color-danger);
	}

	.muted {
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
