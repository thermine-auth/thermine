<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import { SvelteURLSearchParams } from 'svelte/reactivity';
	import {
		RiAddLine,
		RiDeleteBinLine,
		RiDownloadLine,
		RiRefreshLine,
		RiSearchLine,
		RiSettings3Line
	} from 'svelte-remixicon';
	import { ApiError, usersApi, type UserRecord } from '$lib/api';
	import { Alert, Button, Icon, IconButton, SelectionBar } from '$lib/components/ui';
	import FieldsDrawer from '$lib/components/users/FieldsDrawer.svelte';
	import UserDrawer from '$lib/components/users/UserDrawer.svelte';
	import UserTable from '$lib/components/users/UserTable.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	// A writable derived: typing updates it, and it goes back to following the
	// URL whenever that changes, so the back button and a shared link both put
	// the right term in the box.
	let search = $derived(data.search);

	let editing = $state<UserRecord | null>(null);
	let userOpen = $state(false);
	let fieldsOpen = $state(false);

	/** The rows that are ticked, by id. The table prunes any that a search
	    takes out of view, so this only ever holds rows you can act on. */
	let selection = $state<string[]>([]);
	let confirmingDelete = $state(false);
	let busy = $state(false);
	let error = $state('');

	/** The search and the filter are the URL, so the server renders the
	    result and the back button walks through it. */
	async function apply(changes: { search?: string; verified?: string }) {
		const params = new SvelteURLSearchParams(page.url.searchParams);

		for (const [key, value] of Object.entries(changes)) {
			if (value) params.set(key, value);
			else params.delete(key);
		}

		const query = params.toString();
		const path = resolve('/admin/(panel)/dashboard/users');

		// resolve() has already applied any base path; the query is only ever
		// appended to what it returned.
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		await goto(query ? `${path}?${query}` : path, { keepFocus: true, noScroll: true });
	}

	/** Searching as you type, but only once you have paused: every keystroke
	    is a round trip to the server otherwise. */
	let timer: ReturnType<typeof setTimeout>;

	function debounced() {
		clearTimeout(timer);
		timer = setTimeout(() => apply({ search }), 250);
	}

	function openUser(user: UserRecord | null) {
		editing = user;
		userOpen = true;
	}

	let refreshing = $state(false);

	/** Asks the server for the list again, so a record someone else changed
	    shows up without leaving the page.

	    The spin is held for a moment even when the answer comes back at once:
	    a button that does something invisible in 20ms reads as a button that
	    did nothing. */
	async function refresh() {
		refreshing = true;

		try {
			await Promise.all([invalidateAll(), new Promise((done) => setTimeout(done, 400))]);
		} finally {
			refreshing = false;
		}
	}

	function reset() {
		selection = [];
		confirmingDelete = false;
		error = '';
	}

	/** The records behind the ticked ids, in the order the table shows them. */
	function chosen(): UserRecord[] {
		return data.page.users.filter((user) => selection.includes(user.id));
	}

	async function removeSelected() {
		error = '';
		busy = true;

		try {
			for (const id of selection) {
				await usersApi.remove(id);
			}

			reset();
			await invalidateAll();
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Could not delete these users';
		} finally {
			busy = false;
		}
	}

	/** Hands the chosen records to the browser as a file. Nothing leaves the
	    machine: the JSON is built here from what the page already has. */
	function download() {
		const blob = new Blob([JSON.stringify(chosen(), null, 2)], { type: 'application/json' });
		const url = URL.createObjectURL(blob);
		const link = document.createElement('a');

		link.href = url;
		link.download = `users-${new Date().toISOString().slice(0, 10)}.json`;
		link.click();

		URL.revokeObjectURL(url);
	}

	const filters = [
		{ value: '', label: 'All' },
		{ value: 'true', label: 'Verified' },
		{ value: 'false', label: 'Unverified' }
	];
</script>

<svelte:head><title>Users · xermess admin</title></svelte:head>

<header>
	<div class="title">
		<h1>Users</h1>
		<span class="total">{data.page.total} total</span>

		<IconButton icon={RiSettings3Line} label="Field settings" onclick={() => (fieldsOpen = true)} />

		<IconButton
			icon={RiRefreshLine}
			label="Refresh the data"
			onclick={refresh}
			spinning={refreshing}
			disabled={refreshing}
		/>
	</div>

	<div class="actions">
		<Button onclick={() => openUser(null)}>
			<Icon icon={RiAddLine} />
			New user
		</Button>
	</div>
</header>

<div class="toolbar">
	<form
		class="search"
		onsubmit={(event) => {
			event.preventDefault();
			apply({ search });
		}}
	>
		<Icon icon={RiSearchLine} />
		<input
			type="search"
			placeholder="Search email or any field…"
			bind:value={search}
			oninput={debounced}
			aria-label="Search users"
		/>
	</form>

	<div class="filter" role="group" aria-label="Filter by verified">
		{#each filters as filter (filter.value)}
			<button
				type="button"
				class:selected={data.verified === filter.value}
				aria-pressed={data.verified === filter.value}
				onclick={() => apply({ verified: filter.value })}
			>
				{filter.label}
			</button>
		{/each}
	</div>
</div>

{#if error}
	<div class="gutter error"><Alert>{error}</Alert></div>
{/if}

<UserTable users={data.page.users} fields={data.fields} onOpen={openUser} bind:selection />

<SelectionBar count={selection.length} onReset={reset}>
	{#if confirmingDelete}
		<Button
			variant="secondary"
			size="sm"
			onclick={() => (confirmingDelete = false)}
			disabled={busy}
		>
			Keep them
		</Button>
		<Button variant="danger" size="sm" onclick={removeSelected} disabled={busy}>
			{busy ? 'Deleting…' : `Delete ${selection.length}`}
		</Button>
	{:else}
		<Button variant="danger" size="sm" onclick={() => (confirmingDelete = true)} disabled={busy}>
			<Icon icon={RiDeleteBinLine} />
			Delete
		</Button>
		<Button size="sm" onclick={download} disabled={busy}>
			<Icon icon={RiDownloadLine} />
			JSON
		</Button>
	{/if}
</SelectionBar>

<UserDrawer user={editing} fields={data.fields} bind:open={userOpen} />
<FieldsDrawer fields={data.fields} bind:open={fieldsOpen} />

<style>
	header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-3);
		margin-bottom: var(--space-3);
		padding-inline: var(--page-gutter);
	}

	.title {
		display: flex;
		align-items: center;
		gap: var(--space-2);
	}

	.actions {
		display: flex;
		gap: var(--space-2);
	}

	.error {
		margin-bottom: var(--space-3);
	}

	.toolbar {
		display: flex;
		gap: var(--space-2);
		margin-bottom: var(--space-3);
		padding-inline: var(--page-gutter);
	}

	.search {
		flex: 1;
		display: flex;
		align-items: center;
		gap: var(--space-2);
		height: var(--control-height);
		padding: 0 13px;
		border-radius: var(--radius-sm);
		background: var(--color-input);
		color: var(--color-text-hint);
		transition: background-color var(--speed-fast);
	}

	.search:focus-within {
		background: var(--color-input-focus);
	}

	.search input {
		flex: 1;
		min-width: 0;
		border: none;
		background: transparent;
		color: var(--color-text);
		font: inherit;
		font-family: var(--font-sans);
		font-size: var(--text-base);
	}

	.search input:focus {
		outline: none;
	}

	.filter {
		display: flex;
		gap: 2px;
		padding: 2px;
		border-radius: var(--radius-sm);
		background: var(--color-input);
	}

	.filter button {
		padding: 0 var(--space-3);
		border: none;
		border-radius: var(--radius-sm);
		background: transparent;
		color: var(--color-text-hint);
		font: inherit;
		font-size: var(--text-base);
		cursor: pointer;
		transition:
			background-color var(--speed-fast),
			color var(--speed-fast);
	}

	.filter button:hover {
		color: var(--color-text);
	}

	.filter button.selected {
		background: var(--color-surface);
		color: var(--color-text);
		font-weight: 600;
	}

	@media (max-width: 40rem) {
		header,
		.error {
			margin-bottom: var(--space-3);
		}

		.toolbar {
			flex-direction: column;
			align-items: stretch;
		}
	}
</style>
