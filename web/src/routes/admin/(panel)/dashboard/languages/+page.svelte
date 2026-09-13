<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { Badge, DataTable, type Column } from '$lib/components/ui';
	import { demoLanguages } from '$lib/demo';

	const columns: Column[] = [
		{ key: 'name', min: '9rem' },
		{ key: 'code' },
		{ key: 'translated', min: '10rem' },
		{ key: 'status' }
	];
</script>

<svelte:head><title>Languages · xermess admin</title></svelte:head>

<PageHeading
	title="Languages"
	description="The languages the sign-in screens and emails are offered in."
	demo
/>

<DataTable {columns} rows={demoLanguages} empty="No languages.">
	{#snippet row(language)}
		<td>
			<span class="name">
				<strong>{language.name}</strong>
				{#if language.isDefault}
					<Badge>default</Badge>
				{/if}
			</span>
		</td>
		<td class="hint mono">{language.code}</td>
		<td>
			<span class="progress" title="{language.translated}% translated">
				<span class="bar" style="--filled: {language.translated}%"></span>
				<span class="hint count">{language.translated}%</span>
			</span>
		</td>
		<td>
			<Badge tone={language.status === 'active' ? 'success' : 'neutral'}>{language.status}</Badge>
		</td>
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
