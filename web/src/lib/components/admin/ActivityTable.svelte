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
	import { DataTable, Icon, type Column } from '$lib/components/ui';
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

	const columns: Column[] = [
		{ key: 'action', min: '14rem' },
		{ key: 'actor' },
		{ key: 'ip' },
		{ key: 'when', align: 'end' }
	];
</script>

<DataTable {columns} rows={events} empty="Nothing has happened yet.">
	{#snippet row(event)}
		<td class="action" class:failure={isFailure(event.action)}>
			<Icon icon={iconFor(event.action)} />
			{event.action}
		</td>
		<td class="muted">{event.actor || '—'}</td>
		<td class="muted mono">{event.ip}</td>
		<td class="muted when end">{formatDateTime(event.created_at)}</td>
	{/snippet}
</DataTable>

<style>
	/* A cell is a table cell; the icon beside the action needs a box of its
	   own to line up with the text. */
	.action {
		font-family: var(--font-mono);
		font-size: var(--text-sm);
		font-weight: 500;
	}

	.action :global(svg) {
		vertical-align: -3px;
		margin-right: var(--space-2);
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
	}
</style>
