<script lang="ts">
	import ActivityTable from '$lib/components/admin/ActivityTable.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>Logs · xermess admin</title>
</svelte:head>

<header>
	<h1>Logs</h1>
	<p class="subtitle">
		Everything administrators have done, newest first.
		{#if data.logs.length}
			Showing the last {data.logs.length}.
		{/if}
	</p>
</header>

<ActivityTable events={data.logs} />

<style>
	header {
		margin-bottom: var(--space-4);
		padding-inline: var(--page-gutter);
	}

	.subtitle {
		margin-top: var(--space-1);
		color: var(--color-text-hint);
		font-size: var(--text-base);
	}
</style>
