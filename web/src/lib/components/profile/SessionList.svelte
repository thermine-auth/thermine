<script lang="ts">
	import { RiComputerLine, RiSmartphoneLine } from 'svelte-remixicon';
	import type { AdminSession } from '$lib/api';
	import { List, ListItem, Tag, Thumb } from '$lib/components/ui';
	import { describeUserAgent, formatDateTime, formatRelative } from '$lib/utils/format';

	type Props = {
		sessions: AdminSession[];
		/** How many to list, open ones first. Left out, every one. */
		limit?: number;
		/** A box of its own, rather than the body of a flush Panel. */
		bordered?: boolean;
	};

	let { sessions, limit, bordered = false }: Props = $props();

	/** Open sessions first, then the most recent. */
	const shown = $derived(
		[...sessions]
			.sort(
				(a, b) => Number(b.active) - Number(a.active) || b.created_at.localeCompare(a.created_at)
			)
			.slice(0, limit)
	);

	const mobile = (agent: string) => /iPhone|iPad|Android|Mobile/.test(agent);
</script>

{#if shown.length === 0}
	<p class="empty">No sessions recorded.</p>
{:else}
	<List {bordered} label="Sessions">
		{#each shown as session (session.id)}
			<ListItem title={describeUserAgent(session.user_agent)}>
				{#snippet lead()}
					<Thumb icon={mobile(session.user_agent) ? RiSmartphoneLine : RiComputerLine} />
				{/snippet}

				<span class="meta">
					<span class="ip">{session.ip || '—'}</span>
					<time datetime={session.created_at} title={formatDateTime(session.created_at)}>
						signed in {formatRelative(session.created_at)}
					</time>
				</span>

				{#snippet end()}
					<Tag tone={session.active ? 'success' : 'neutral'} dot strong>
						{session.active ? 'Active' : 'Ended'}
					</Tag>
				{/snippet}
			</ListItem>
		{/each}
	</List>
{/if}

<style>
	.empty {
		margin: 0;
		padding: var(--space-4);
		color: var(--color-text-hint);
	}

	.meta {
		display: flex;
		flex-wrap: wrap;
		gap: 0 var(--space-2);
		color: var(--color-text-hint);
		font-size: var(--text-xs);
	}

	.ip {
		font-family: var(--font-mono);
	}
</style>
