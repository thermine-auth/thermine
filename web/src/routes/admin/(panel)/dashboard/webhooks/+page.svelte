<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { Badge, DataTable } from '$lib/components/ui';
	import { demoWebhooks } from '$lib/demo';
</script>

<svelte:head>
	<title>Webhooks · xermess admin</title>
</svelte:head>

<PageHeading title="Webhooks" description="Where events are delivered as they happen." demo />

<DataTable
	rows={demoWebhooks}
	empty="No webhooks."
	columns="minmax(8rem, 1fr) minmax(12rem, 1.8fr) auto auto"
>
	{#snippet row(hook)}
		<span class="mono event">{hook.event}</span>
		<span class="hint url">{hook.url}</span>
		<span class="hint counts">
			{hook.delivered.toLocaleString()} delivered
			{#if hook.failed}
				<span class="failed">· {hook.failed} failed</span>
			{/if}
		</span>
		<Badge tone={hook.active ? 'success' : 'neutral'}>{hook.active ? 'active' : 'paused'}</Badge>
	{/snippet}
</DataTable>

<style>
	.event {
		font-weight: 600;
	}

	.hint {
		color: var(--color-text-hint);
	}

	.mono {
		font-family: var(--font-mono);
		font-size: var(--text-sm);
	}

	.url {
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.counts {
		font-size: var(--text-sm);
		white-space: nowrap;
	}

	.failed {
		color: var(--color-danger);
	}
</style>
