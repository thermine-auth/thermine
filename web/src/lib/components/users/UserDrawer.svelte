<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { Switch } from '@ark-ui/svelte/switch';
	import { Field } from '@ark-ui/svelte/field';
	import { RiKey2Line, RiMailLine } from 'svelte-remixicon';
	import { ApiError, usersApi, type UserField, type UserRecord } from '$lib/api';
	import { Alert, Button, Drawer, Icon } from '$lib/components/ui';
	import FieldInput from './FieldInput.svelte';

	type Props = {
		/** The user being edited, or null to create one. */
		user: UserRecord | null;
		fields: UserField[];
		open: boolean;
	};

	let { user, fields, open = $bindable(false) }: Props = $props();

	let email = $state('');
	let emailVerified = $state(false);
	let values = $state<Record<string, string | boolean>>({});
	let error = $state('');
	let saving = $state(false);

	const editing = $derived(user !== null);

	/** The record's own fields are grouped: the ones that hold a value, then
	    the flags, which read as a list of yes-or-no answers rather than as
	    inputs someone has to fill in. */
	const details = $derived(fields.filter((field) => field.type !== 'bool'));
	const flags = $derived(fields.filter((field) => field.type === 'bool'));

	/** Fill the form whenever the drawer is opened for a different user.
	    Dates are stored as timestamps and edited as days. */
	$effect(() => {
		if (!open) return;

		email = user?.email ?? '';
		emailVerified = user?.email_verified ?? false;
		error = '';

		values = Object.fromEntries(
			fields.map((field) => {
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

		for (const field of fields) {
			const value = values[field.name];

			if (field.type === 'bool') {
				data[field.name] = value === true;
			} else if (String(value ?? '').trim() !== '') {
				data[field.name] = value;
			}
		}

		return { email: email.trim(), email_verified: emailVerified, data };
	}

	async function save(event: SubmitEvent) {
		event.preventDefault();

		error = '';
		saving = true;

		try {
			if (user) {
				await usersApi.update(user.id, payload());
			} else {
				await usersApi.create(payload());
			}

			await invalidateAll();
			open = false;
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Could not save this user';
		} finally {
			saving = false;
		}
	}
</script>

<Drawer bind:open title={editing ? 'Edit user record' : 'New user record'} onsubmit={save}>
	{#if error}
		<div class="error"><Alert>{error}</Alert></div>
	{/if}

	<section>
		<h3>Account</h3>

		{#if editing}
			<!-- The id is what every other system refers to this record by, so
			     it is shown and can be copied, but it is not something to
			     edit. -->
			<Field.Root readOnly>
				<Field.Label>
					<Icon icon={RiKey2Line} />
					id
				</Field.Label>
				<Field.Input value={user?.id ?? ''} readonly />
			</Field.Root>
		{/if}

		<Field.Root required>
			<Field.Label>
				<Icon icon={RiMailLine} />
				email
				<Field.RequiredIndicator>*</Field.RequiredIndicator>
			</Field.Label>
			<Field.Input
				value={email}
				oninput={(event) => (email = event.currentTarget.value)}
				type="email"
				autocomplete="off"
				placeholder="user@example.com"
			/>
		</Field.Root>

		<Switch.Root
			checked={emailVerified}
			onCheckedChange={(details) => (emailVerified = details.checked)}
		>
			<Switch.Control><Switch.Thumb /></Switch.Control>
			<Switch.Label>email_verified</Switch.Label>
			<Switch.HiddenInput />
		</Switch.Root>
	</section>

	{#if details.length > 0}
		<section>
			<h3>Details</h3>

			{#each details as field (field.id)}
				<FieldInput {field} bind:value={values[field.name]} />
			{/each}
		</section>
	{/if}

	{#if flags.length > 0}
		<section>
			<h3>Flags</h3>

			<div class="flags">
				{#each flags as field (field.id)}
					<FieldInput {field} bind:value={values[field.name]} />
				{/each}
			</div>
		</section>
	{/if}

	{#snippet footer()}
		<span class="spacer"></span>

		<Button variant="secondary" onclick={() => (open = false)} disabled={saving}>Cancel</Button>

		<Button type="submit" disabled={saving || email.trim() === ''}>
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
</style>
