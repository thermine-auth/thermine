<script lang="ts">
	import { RiComputerLine } from 'svelte-remixicon';
	import type { AdminSession } from '$lib/api';
	import { Badge, DataTable, Icon } from '$lib/components/ui';
	import { formatDateTime } from '$lib/format';

	type Props = { sessions: AdminSession[] };

	let { sessions }: Props = $props();
</script>

<DataTable
	rows={sessions}
	empty="No sessions recorded."
	columns="auto minmax(5rem, 1fr) minmax(8rem, 2fr) auto"
>
	{#snippet row(session)}
		<Badge tone={session.active ? 'success' : 'neutral'}>
			{session.active ? 'active' : 'ended'}
		</Badge>
		<span class="muted mono">{session.ip}</span>
		<span class="muted agent">
			<Icon icon={RiComputerLine} />
			{session.user_agent || '—'}
		</span>
		<span class="muted when">{formatDateTime(session.created_at)}</span>
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

	.agent {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		overflow: hidden;
		white-space: nowrap;
	}

	.agent :global(+ *),
	.agent {
		text-overflow: ellipsis;
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
