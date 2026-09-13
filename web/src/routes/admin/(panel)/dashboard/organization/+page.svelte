<script lang="ts">
	import PageHeading from '$lib/components/admin/PageHeading.svelte';
	import { Badge, Card } from '$lib/components/ui';
	import { demoOrganization as org } from '$lib/demo';

	const details = [
		{ label: 'Name', value: org.name },
		{ label: 'Identifier', value: org.slug, mono: true },
		{ label: 'Primary domain', value: org.domain, mono: true },
		{ label: 'Region', value: org.region, mono: true },
		{ label: 'Support email', value: org.supportEmail },
		{ label: 'Created', value: org.created }
	];
</script>

<svelte:head><title>Organization · xermess admin</title></svelte:head>

<PageHeading
	title="Organization"
	description="Who this tenant belongs to and where it runs."
	demo
/>

<div class="gutter">
	<Card>
		<dl>
			{#each details as detail (detail.label)}
				<div class="row">
					<dt>{detail.label}</dt>
					<dd class:mono={detail.mono}>{detail.value}</dd>
				</div>
			{/each}
			<div class="row">
				<dt>Plan</dt>
				<dd><Badge tone="success">{org.plan}</Badge></dd>
			</div>
		</dl>
	</Card>
</div>

<style>
	dl {
		margin: 0;
	}

	.row {
		display: grid;
		grid-template-columns: 10rem 1fr;
		gap: var(--space-4);
		align-items: center;
		padding: var(--space-3) var(--space-4);
		font-size: var(--text-base);
	}

	.row + .row {
		border-top: 1px solid var(--color-border);
	}

	dt {
		color: var(--color-text-hint);
	}
	dd {
		margin: 0;
	}
	.mono {
		font-family: var(--font-mono);
		font-size: var(--text-sm);
	}

	@media (max-width: 40rem) {
		.row {
			grid-template-columns: 1fr;
			gap: var(--space-1);
		}
	}
</style>
