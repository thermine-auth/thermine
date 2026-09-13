<script lang="ts">
	import type { FieldType } from '$lib/api';
	import { Badge } from '$lib/components/ui';

	type Props = {
		type: FieldType;
		value: unknown;
	};

	let { type, value }: Props = $props();

	const empty = $derived(value === null || value === undefined || value === '');

	/** Dates are stored as timestamps and read as days. */
	function asDate(raw: unknown): string {
		const date = new Date(String(raw));
		return Number.isNaN(date.getTime()) ? String(raw) : date.toLocaleDateString();
	}
</script>

{#if empty}
	<!-- An empty cell is still a cell: saying so is clearer than a blank. -->
	<span class="empty">N/A</span>
{:else if type === 'bool'}
	<Badge tone={value ? 'success' : 'neutral'}>{value ? 'True' : 'False'}</Badge>
{:else if type === 'number'}
	<span class="mono">{Number(value).toLocaleString()}</span>
{:else if type === 'date'}
	<span class="mono">{asDate(value)}</span>
{:else}
	<span class="text">{value}</span>
{/if}

<style>
	.empty {
		color: var(--color-text-disabled);
		font-size: var(--text-sm);
	}

	.mono {
		font-family: var(--font-mono);
		font-size: var(--text-sm);
	}

	.text {
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
</style>
