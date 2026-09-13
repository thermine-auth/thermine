<script lang="ts">
	import { RiVerifiedBadgeLine } from 'svelte-remixicon';
	import type { UserField, UserRecord } from '$lib/api';
	import { Badge, DataTable, type Column } from '$lib/components/ui';
	import FieldValue from './FieldValue.svelte';
	import { fieldIcons, idIcon } from './fieldIcons';

	type Props = {
		users: UserRecord[];
		fields: UserField[];
		/** Called with the user whose row was chosen. */
		onOpen: (user: UserRecord) => void;
		/** The ids of the ticked rows; bind to act on them. */
		selection: string[];
	};

	let { users, fields, onOpen, selection = $bindable() }: Props = $props();

	/** id and the two built-in columns, then whatever fields are defined, so
	    a new field needs no change here. */
	const columns = $derived<Column[]>([
		{ key: 'id', label: 'id', icon: idIcon, min: '10rem' },
		{ key: 'email', label: 'email', icon: fieldIcons.email, min: '13rem' },
		{ key: 'email_verified', label: 'email_verified', icon: RiVerifiedBadgeLine, min: '10rem' },
		...fields.map((field) => ({
			key: field.name,
			label: field.name,
			icon: fieldIcons[field.type],
			min: '9rem'
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
	bind:selection
>
	{#snippet row(user)}
		<td><span class="chip">{shortId(user.id)}</span></td>
		<td>{user.email}</td>
		<td>
			<Badge tone={user.email_verified ? 'success' : 'neutral'}>
				{user.email_verified ? 'True' : 'False'}
			</Badge>
		</td>
		{#each fields as field (field.id)}
			<td><FieldValue type={field.type} value={user.data?.[field.name]} /></td>
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
