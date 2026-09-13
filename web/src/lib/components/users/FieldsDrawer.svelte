<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { Field } from '@ark-ui/svelte/field';
	import { Switch } from '@ark-ui/svelte/switch';
	import { RiAddLine, RiDeleteBinLine, RiPencilLine } from 'svelte-remixicon';
	import { ApiError, usersApi, type FieldType, type UserField } from '$lib/api';
	import { Alert, Badge, Button, Drawer, Icon, IconButton } from '$lib/components/ui';
	import { fieldIcons } from './fieldIcons';

	type Props = {
		fields: UserField[];
		open: boolean;
	};

	let { fields, open = $bindable(false) }: Props = $props();

	const types: FieldType[] = ['text', 'number', 'bool', 'email', 'date'];

	/** The field being edited, or null while a new one is being written. A
	    field's name and type are fixed once records hold values under them,
	    so editing offers its rules only. */
	let editing = $state<UserField | null>(null);

	let name = $state('');
	let label = $state('');
	let type = $state<FieldType>('text');
	let required = $state(false);
	let unique = $state(false);
	let min = $state('');
	let max = $state('');
	let startsWith = $state('');

	let error = $state('');
	let busy = $state(false);

	/** Which rules this type of field can keep. Bools and dates carry none of
	    them, so the inputs are not offered. */
	const bounded = $derived(type === 'number' || type === 'text' || type === 'email');
	const prefixed = $derived(type === 'text' || type === 'email');
	const lengths = $derived(type !== 'number');

	/** A field name is a column name: the panel suggests one from the label so
	    nobody has to think about the rule. */
	function suggestName(from: string): string {
		return from
			.trim()
			.toLowerCase()
			.replace(/[^a-z0-9]+/g, '_')
			.replace(/^[^a-z]+/, '')
			.replace(/_+$/, '')
			.slice(0, 64);
	}

	function blank() {
		editing = null;
		name = '';
		label = '';
		type = 'text';
		required = false;
		unique = false;
		min = '';
		max = '';
		startsWith = '';
		error = '';
	}

	function edit(field: UserField) {
		editing = field;
		name = field.name;
		label = field.label;
		type = field.type;
		required = field.required;
		unique = field.unique;
		min = field.min === null ? '' : String(field.min);
		max = field.max === null ? '' : String(field.max);
		startsWith = field.starts_with;
		error = '';
	}

	/** An empty box means "no bound", which the API reads as null. */
	function bound(value: string): number | null {
		const trimmed = value.trim();
		if (trimmed === '') return null;

		const parsed = Number(trimmed);
		return Number.isFinite(parsed) ? parsed : null;
	}

	function rules() {
		return {
			label: label.trim() || name.trim(),
			required,
			unique,
			min: bounded ? bound(min) : null,
			max: bounded ? bound(max) : null,
			starts_with: prefixed ? startsWith.trim() : ''
		};
	}

	async function submit(event: SubmitEvent) {
		event.preventDefault();

		error = '';
		busy = true;

		try {
			if (editing) {
				await usersApi.updateField(editing.id, rules());
			} else {
				await usersApi.addField({
					...rules(),
					name: name.trim() || suggestName(label),
					type
				});
			}

			await invalidateAll();
			blank();
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Could not save this field';
		} finally {
			busy = false;
		}
	}

	async function remove(field: UserField) {
		error = '';
		busy = true;

		try {
			await usersApi.removeField(field.id);
			await invalidateAll();

			if (editing?.id === field.id) blank();
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Could not remove this field';
		} finally {
			busy = false;
		}
	}

	/** What a field expects, said in a few words for the list. */
	function summary(field: UserField): string[] {
		const said: string[] = [];

		if (field.required) said.push('required');
		if (field.unique) said.push('unique');

		const unit = field.type === 'number' ? '' : ' chars';
		if (field.min !== null && field.max !== null) said.push(`${field.min}–${field.max}${unit}`);
		else if (field.min !== null) said.push(`min ${field.min}${unit}`);
		else if (field.max !== null) said.push(`max ${field.max}${unit}`);

		if (field.starts_with) said.push(`starts with ${field.starts_with}`);

		return said;
	}
</script>

<Drawer
	bind:open
	title="User fields"
	description="The columns a user record has, and the rules their values keep. Adding one needs no migration: the values live in the record itself."
>
	{#if error}
		<div class="error"><Alert>{error}</Alert></div>
	{/if}

	<ul class="fields">
		{#each fields as field (field.id)}
			<li class:editing={editing?.id === field.id}>
				<Icon icon={fieldIcons[field.type]} />

				<span class="name">{field.name}</span>
				<span class="hint">{field.type}</span>

				<span class="rules">
					{#each summary(field) as rule (rule)}
						<Badge>{rule}</Badge>
					{/each}
				</span>

				<IconButton
					icon={RiPencilLine}
					label="Edit {field.name}"
					size="sm"
					onclick={() => edit(field)}
					disabled={busy}
				/>

				<IconButton
					icon={RiDeleteBinLine}
					label="Remove {field.name}"
					size="sm"
					tone="danger"
					onclick={() => remove(field)}
					disabled={busy}
				/>
			</li>
		{:else}
			<li class="hint">No fields yet. A user has only an email and its verified flag.</li>
		{/each}
	</ul>

	<form onsubmit={submit}>
		<h3>{editing ? `Edit ${editing.name}` : 'Add a field'}</h3>

		<div class="row">
			<Field.Root>
				<Field.Label>Label</Field.Label>
				<Field.Input
					value={label}
					oninput={(event) => {
						label = event.currentTarget.value;
						if (!editing) name = suggestName(label);
					}}
					placeholder="Phone number"
				/>
			</Field.Root>

			<Field.Root readOnly={editing !== null}>
				<Field.Label>Name</Field.Label>
				<Field.Input
					value={name}
					oninput={(event) => (name = event.currentTarget.value)}
					placeholder="phone_number"
					readonly={editing !== null}
				/>
			</Field.Root>

			<Field.Root readOnly={editing !== null}>
				<Field.Label>Type</Field.Label>
				<Field.Select
					value={type}
					disabled={editing !== null}
					onchange={(event) => (type = event.currentTarget.value as FieldType)}
				>
					{#each types as option (option)}
						<option value={option}>{option}</option>
					{/each}
				</Field.Select>
			</Field.Root>
		</div>

		{#if bounded || prefixed}
			<div class="row">
				{#if bounded}
					<Field.Root>
						<Field.Label>{lengths ? 'Least characters' : 'Smallest value'}</Field.Label>
						<Field.Input
							value={min}
							oninput={(event) => (min = event.currentTarget.value)}
							type="number"
							placeholder="any"
						/>
					</Field.Root>

					<Field.Root>
						<Field.Label>{lengths ? 'Most characters' : 'Largest value'}</Field.Label>
						<Field.Input
							value={max}
							oninput={(event) => (max = event.currentTarget.value)}
							type="number"
							placeholder="any"
						/>
					</Field.Root>
				{/if}

				{#if prefixed}
					<Field.Root>
						<Field.Label>Must start with</Field.Label>
						<Field.Input
							value={startsWith}
							oninput={(event) => (startsWith = event.currentTarget.value)}
							placeholder="+"
						/>
					</Field.Root>
				{/if}
			</div>
		{/if}

		<div class="switches">
			<Switch.Root checked={required} onCheckedChange={(details) => (required = details.checked)}>
				<Switch.Control><Switch.Thumb /></Switch.Control>
				<Switch.Label>Required</Switch.Label>
				<Switch.HiddenInput />
			</Switch.Root>

			<Switch.Root checked={unique} onCheckedChange={(details) => (unique = details.checked)}>
				<Switch.Control><Switch.Thumb /></Switch.Control>
				<Switch.Label>Unique</Switch.Label>
				<Switch.HiddenInput />
			</Switch.Root>
		</div>

		<div class="actions">
			{#if editing}
				<Button variant="secondary" onclick={blank} disabled={busy}>Cancel</Button>
			{/if}

			<Button type="submit" disabled={busy || (name.trim() === '' && label.trim() === '')}>
				{#if !editing}
					<Icon icon={RiAddLine} />
				{/if}
				{editing ? 'Save field' : 'Add field'}
			</Button>
		</div>
	</form>

	{#snippet footer()}
		<span class="spacer"></span>
		<Button variant="secondary" onclick={() => (open = false)}>Done</Button>
	{/snippet}
</Drawer>

<style>
	.error {
		margin-bottom: var(--space-4);
	}

	.fields {
		display: flex;
		flex-direction: column;
		gap: 2px;
		margin: 0;
		padding: 0;
		list-style: none;
	}

	.fields li {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-2);
		border-radius: var(--radius-sm);
		font-size: var(--text-base);
	}

	.fields li:hover {
		background: var(--row-hover);
	}

	.fields li.editing {
		background: var(--surface-info);
	}

	.name {
		font-family: var(--font-mono);
		font-size: var(--text-sm);
		font-weight: 600;
	}

	.hint {
		color: var(--color-text-hint);
		font-size: var(--text-sm);
	}

	/* The rules take the space left over, so the buttons stay in a column
	   down the right whatever a field expects. */
	.rules {
		flex: 1;
		display: flex;
		flex-wrap: wrap;
		justify-content: flex-end;
		gap: var(--space-1);
	}

	form {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
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

	.row {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(9rem, 1fr));
		gap: var(--space-2);
	}

	.switches {
		display: flex;
		gap: var(--space-5);
	}

	.actions {
		display: flex;
		justify-content: flex-end;
		gap: var(--space-2);
	}

	.spacer {
		flex: 1;
	}
</style>
