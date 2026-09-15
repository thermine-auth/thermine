<script lang="ts">
	import { RiShieldUserLine } from 'svelte-remixicon';
	import type { Application, UserField, UserRecord } from '$lib/api';
	import { Badge, DataTable, type Column } from '$lib/components/ui';
	import FieldValue from './FieldValue.svelte';
	import { valueOf } from './fields';
	import { fieldIcons, idIcon } from './fieldIcons';

	type Props = {
		users: UserRecord[];
		fields: UserField[];
		/** The applications the roles belong to, to name them. */
		applications: Application[];
		/** Called with the user whose row was chosen. */
		onOpen: (user: UserRecord) => void;
		/** The ids of the ticked rows, and how to say that they changed. Left
		    out, the rows cannot be ticked. */
		selected?: string[];
		onSelect?: (ids: string[]) => void;
	};

	let { users, fields, applications, onOpen, selected, onSelect }: Props = $props();

	const appNames = $derived(new Map(applications.map((app) => [app.id, app.name])));

	/** The id, then every field the API listed — the built-in ones first and
	    then whatever this organisation added — then the roles. One list, so a
	    new field needs no change here. */
	const columns = $derived<Column[]>([
		{ key: 'id', label: 'id', icon: idIcon, min: '10rem' },
		...fields.map((field) => ({
			key: field.name,
			label: field.name,
			icon: fieldIcons[field.type],
			min: field.type === 'email' ? '13rem' : field.type === 'bool' ? '8rem' : '9rem'
		})),
		{ key: 'roles', label: 'roles', icon: RiShieldUserLine, min: '10rem' }
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

		<td>
			<span class="roles">
				{#each user.roles ?? [] as role (role.id)}
					<!-- A global role reads as itself; an application role is
					     prefixed with its application, since the same name can be
					     a role in two of them. -->
					<Badge>
						{#if role.application_id}
							<span class="app">{appNames.get(role.application_id) ?? '?'}</span>
						{/if}
						{role.name}
					</Badge>
				{:else}
					<span class="empty">N/A</span>
				{/each}
			</span>
		</td>
	{/snippet}
</DataTable>

<style>
	.roles {
		display: inline-flex;
		gap: var(--space-1);
	}

	/* The application a role belongs to, quieter than the role: the same
	   name can be a role in two applications. */
	.app {
		margin-right: 4px;
		opacity: 0.65;
	}

	.app::after {
		content: ' /';
	}

	.empty {
		color: var(--color-text-disabled);
		font-size: var(--text-sm);
	}

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
