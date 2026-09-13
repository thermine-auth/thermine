<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { Badge, DataTable } from '$lib/components/ui';
	import { demoFlows } from '$lib/demo';
</script>

<svelte:head><title>Login flows · xermess admin</title></svelte:head>

<PageHeading
	title="Login flows"
	description="The steps someone is taken through when they sign in."
	demo
/>

<DataTable
	rows={demoFlows}
	empty="No flows defined."
	columns="minmax(8rem, 1fr) minmax(11rem, 1.6fr) auto auto"
>
	{#snippet row(flow)}
		<span class="name">
			<strong>{flow.name}</strong>
			{#if flow.isDefault}
				<Badge>default</Badge>
			{/if}
		</span>
		<span class="hint mono steps">{flow.steps}</span>
		<span class="hint count">{flow.applications} apps</span>
		<Badge tone={flow.status === 'active' ? 'success' : 'neutral'}>{flow.status}</Badge>
	{/snippet}
</DataTable>

<style>
	.name {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}
	.hint {
		color: var(--color-text-hint);
	}
	.mono {
		font-family: var(--font-mono);
		font-size: var(--text-sm);
	}
	.steps {
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}
	.count {
		font-size: var(--text-sm);
		white-space: nowrap;
	}
</style>
