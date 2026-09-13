<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { Badge, DataTable } from '$lib/components/ui';
	import { demoLanguages } from '$lib/demo';
</script>

<svelte:head><title>Languages · xermess admin</title></svelte:head>

<PageHeading
	title="Languages"
	description="The languages the sign-in screens and emails are offered in."
	demo
/>

<DataTable
	rows={demoLanguages}
	empty="No languages."
	columns="minmax(7rem, 1fr) auto minmax(6rem, 1fr) auto"
>
	{#snippet row(language)}
		<span class="name">
			<strong>{language.name}</strong>
			{#if language.isDefault}
				<Badge>default</Badge>
			{/if}
		</span>
		<span class="hint mono">{language.code}</span>
		<span class="progress" title="{language.translated}% translated">
			<span class="bar" style="--filled: {language.translated}%"></span>
			<span class="hint count">{language.translated}%</span>
		</span>
		<Badge tone={language.status === 'active' ? 'success' : 'neutral'}>{language.status}</Badge>
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
	.count {
		font-size: var(--text-sm);
		white-space: nowrap;
	}

	.progress {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}

	.bar {
		flex: 1;
		height: 4px;
		min-width: 3rem;
		border-radius: var(--radius-pill);
		background: var(--color-secondary-alt);
	}

	.bar::before {
		display: block;
		width: var(--filled);
		height: 100%;
		border-radius: inherit;
		background: var(--color-success);
		content: '';
	}
</style>
