<script lang="ts">
	import type { UserField, UserRecord } from '$lib/api';
	import { DataTable, type Column } from '$lib/components/ui';
	import FieldValue from './FieldValue.svelte';
	import { valueOf } from './fields';
	import { fieldIcons, idIcon } from './fieldIcons';

	type Props = {
		users: UserRecord[];
		fields: UserField[];
		/** Called with the user whose row was chosen. */
		onOpen: (user: UserRecord) => void;
		/** The ids of the ticked rows, and how to say that they changed. */
		selected: string[];
		onSelect: (ids: string[]) => void;
	};

	let { users, fields, onOpen, selected, onSelect }: Props = $props();

	/** The id, then every field the API listed: the built-in ones first and
	    then whatever this organisation added. One list, so a new field needs
	    no change here. */
	const columns = $derived<Column[]>([
		{ key: 'id', label: 'id', icon: idIcon, min: '10rem' },
		...fields.map((field) => ({
			key: field.name,
			label: field.name,
			icon: fieldIcons[field.type],
			min: field.type === 'email' ? '13rem' : field.type === 'bool' ? '8rem' : '9rem'
		}))
	]);

	/** The first characters of the id, which is all anyone reads of it. */
	function shortId(id: string): string {
		return id.replace(/-/g, '').slice(0, 15);
	}
</script>

<DataTable
	{columns}
	rows={users}
	empty="No users match this."
	{onOpen}
	label={(user) => `Edit ${user.email}`}
	{selected}
	{onSelect}
>
	{#snippet row(user)}
		<td><span class="chip">{shortId(user.id)}</span></td>

		{#each fields as field (field.name)}
			<td><FieldValue type={field.type} value={valueOf(user, field)} /></td>
		{/each}
	{/snippet}
</DataTable>

<style>
	/* PocketBase shows the id as a small chip rather than raw text, which
	   stops it competing with the values beside it. */
	.chip {
		display: inline-flex;
		align-items: center;
		height: 25px;
		padding: 0 7px;
		border-radius: var(--radius-sm);
		background: var(--color-secondary-alt);
		color: var(--color-text);
		font-family: var(--font-mono);
		font-size: var(--text-sm);
	}
</style>
