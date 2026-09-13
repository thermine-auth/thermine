<script lang="ts">
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { RiKey2Line, RiMailLine } from 'svelte-remixicon';
	import { ApiError, usersApi, type UserField, type UserRecord } from '$lib/api';
	import { Alert, Button, Drawer, Input, Switch } from '$lib/components/ui';
	import { keys } from '$lib/query';
	import { additional } from './fields';
	import FieldInput from './FieldInput.svelte';

	type Props = {
		/** The user being edited, or null to create one. */
		user: UserRecord | null;
		fields: UserField[];
		open: boolean;
	};

	let { user, fields, open = $bindable(false) }: Props = $props();

	const queryClient = useQueryClient();

	// The built-in fields are the record's own columns, so they are named
	// here; the additional ones are whatever this organisation added.
	let email = $state('');
	let emailVerified = $state(false);
	let firstName = $state('');
	let lastName = $state('');
	let isActive = $state(true);

	let values = $state<Record<string, string | boolean>>({});
	let error = $state('');

	/** True while the record is being written. Ours rather than the
	    mutation's own isPending, so a form cannot be left saying "Saving…". */
	let saving = $state(false);

	const editing = $derived(user !== null);

	/** The added fields, split the way they are filled in: values first, then
	    the yes-or-no answers. */
	const extras = $derived(additional(fields));
	const details = $derived(extras.filter((field) => field.type !== 'bool'));
	const flags = $derived(extras.filter((field) => field.type === 'bool'));

	/** Fill the form whenever the drawer is opened for a different user.
	    Dates are stored as timestamps and edited as days. */
	$effect(() => {
		if (!open) return;

		email = user?.email ?? '';
		emailVerified = user?.email_verified ?? false;
		firstName = user?.first_name ?? '';
		lastName = user?.last_name ?? '';
		isActive = user?.is_active ?? true;
		error = '';

		values = Object.fromEntries(
			extras.map((field) => {
				const stored = user?.data?.[field.name];

				if (field.type === 'bool') return [field.name, stored === true];
				if (field.type === 'date') return [field.name, stored ? String(stored).slice(0, 10) : ''];

				return [field.name, stored == null ? '' : String(stored)];
			})
		);
	});

	/** Empty text is left out entirely, which is how "not set" is stored. */
	function payload() {
		const data: Record<string, unknown> = {};

		for (const field of extras) {
			const value = values[field.name];

			if (field.type === 'bool') {
				data[field.name] = value === true;
			} else if (String(value ?? '').trim() !== '') {
				data[field.name] = value;
			}
		}

		return {
			email: email.trim(),
			email_verified: emailVerified,
			first_name: firstName.trim(),
			last_name: lastName.trim(),
			is_active: isActive,
			data
		};
	}

	/** Writing the record. Whether that is a new one or an edit is the only
	    difference; what happens afterwards — the list refilled, the panel
	    closed — is the same either way. */
	const save = createMutation(() => ({
		mutationFn: () => (user ? usersApi.update(user.id, payload()) : usersApi.create(payload())),
		onSuccess: async () => {
			await queryClient.invalidateQueries({ queryKey: keys.users.all });
			open = false;
		},
		onError: (err: unknown) => {
			error = err instanceof ApiError ? err.message : 'Could not save this user';
		},
		onSettled: () => {
			saving = false;
		}
	}));

	function submit(event: SubmitEvent) {
		event.preventDefault();

		// Enter and a double click both submit; one save at a time is enough,
		// and two would race each other to write the same record.
		if (saving) return;

		error = '';
		saving = true;
		save.mutate();
	}
</script>

<Drawer bind:open title={editing ? 'Edit user record' : 'New user record'} onsubmit={submit}>
	{#if error}
		<div class="error"><Alert>{error}</Alert></div>
	{/if}

	<section>
		<h3>Account</h3>

		{#if editing}
			<!-- The id is what every other system refers to this record by, so
			     it is shown and can be copied, but it is not something to
			     edit. -->
			<Input label="id" icon={RiKey2Line} value={user?.id ?? ''} readOnly />
		{/if}

		<Input
			label="email"
			icon={RiMailLine}
			bind:value={email}
			type="email"
			autocomplete="off"
			placeholder="user@example.com"
			required
		/>

		<div class="names">
			<Input label="first_name" bind:value={firstName} />
			<Input label="last_name" bind:value={lastName} />
		</div>
	</section>

	<section>
		<h3>Flags</h3>

		<div class="flags">
			<Switch label="email_verified" bind:checked={emailVerified} />
			<Switch label="is_active" bind:checked={isActive} />

			{#each flags as field (field.id)}
				<FieldInput {field} bind:value={values[field.name]} />
			{/each}
		</div>
	</section>

	{#if details.length > 0}
		<section>
			<h3>Additional fields</h3>

			{#each details as field (field.id)}
				<FieldInput {field} bind:value={values[field.name]} />
			{/each}
		</section>
	{/if}

	{#snippet footer()}
		<span class="spacer"></span>

		<Button variant="subtle" onclick={() => (open = false)} disabled={saving}>Cancel</Button>

		<Button type="submit" loading={saving} disabled={saving || email.trim() === ''}>
			{saving ? 'Saving…' : editing ? 'Save changes' : 'Create user'}
		</Button>
	{/snippet}
</Drawer>

<style>
	section {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}

	.error {
		margin-bottom: var(--space-4);
	}

	section + section {
		margin-top: var(--space-5);
		padding-top: var(--space-5);
		border-top: 1px solid var(--color-border);
	}

	h3 {
		margin: 0;
		color: var(--color-text-hint);
		font-size: var(--text-xs);
		font-weight: 600;
		letter-spacing: 0.08em;
		text-transform: uppercase;
	}

	.names {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--space-3);
	}

	/* Flags are short, so they sit two to a row where there is room rather
	   than running down the panel one by one. */
	.flags {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(14rem, 1fr));
		gap: var(--space-2) var(--space-4);
	}

	.spacer {
		flex: 1;
	}

	@media (max-width: 30rem) {
		.names {
			grid-template-columns: 1fr;
		}
	}
</style>
