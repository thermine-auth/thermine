<script lang="ts">
	import { RiComputerLine } from 'svelte-remixicon';
	import type { AdminSession } from '$lib/api';
	import { Badge, DataTable, Icon, type Column } from '$lib/components/ui';
	import { formatDateTime } from '$lib/format';

	type Props = { sessions: AdminSession[] };

	let { sessions }: Props = $props();

	const columns: Column[] = [
		{ key: 'state' },
		{ key: 'ip' },
		{ key: 'agent', min: '14rem' },
		{ key: 'when', align: 'end' }
	];
</script>

<DataTable {columns} rows={sessions} empty="No sessions recorded.">
	{#snippet row(session)}
		<td>
			<Badge tone={session.active ? 'success' : 'neutral'}>
				{session.active ? 'active' : 'ended'}
			</Badge>
		</td>
		<td class="muted mono">{session.ip}</td>
		<td class="muted agent">
			<Icon icon={RiComputerLine} />
			{session.user_agent || '—'}
		</td>
		<td class="muted when end">{formatDateTime(session.created_at)}</td>
	{/snippet}
</DataTable>

<style>
	.muted {
		color: var(--color-text-hint);
	}

	.mono {
		font-family: var(--font-mono);
		font-size: var(--text-sm);
	}

	/* A user agent is a long string nobody reads to the end: it is given a
	   share of the row and cut off there rather than widening the table. */
	.agent {
		max-width: 30rem;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.agent :global(svg) {
		vertical-align: -3px;
		margin-right: var(--space-2);
	}

	.when {
		font-size: var(--text-sm);
	}
</style>
