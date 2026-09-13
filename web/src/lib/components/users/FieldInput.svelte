<script lang="ts">
	import { Field } from '@ark-ui/svelte/field';
	import { Switch } from '@ark-ui/svelte/switch';
	import type { UserField } from '$lib/api';
	import { Icon } from '$lib/components/ui';
	import { fieldIcons } from './fieldIcons';

	type Props = {
		field: UserField;
		/** Booleans are held as booleans; everything else as the text in the
		    input. The form converts on the way in and out, so this component
		    holds no state of its own.
		    
		    It is undefined for the moment before the form has filled itself
		    in, so there is deliberately no fallback here: binding undefined to
		    a prop that has one is an error in Svelte. */
		value: string | boolean | undefined;
	};

	let { field, value = $bindable() }: Props = $props();

	/** The input type that suits the field, so dates and numbers get the
	    browser's own pickers rather than us building them. */
	const inputType = $derived(
		field.type === 'number' ? 'number' : field.type === 'date' ? 'date' : 'text'
	);
</script>

{#if field.type === 'bool'}
	<Switch.Root checked={value === true} onCheckedChange={(details) => (value = details.checked)}>
		<Switch.Control><Switch.Thumb /></Switch.Control>
		<Switch.Label>{field.label}</Switch.Label>
		<Switch.HiddenInput />
	</Switch.Root>
{:else}
	<Field.Root required={field.required}>
		<Field.Label>
			<Icon icon={fieldIcons[field.type]} />
			{field.label}
			{#if field.required}
				<Field.RequiredIndicator>*</Field.RequiredIndicator>
			{/if}
		</Field.Label>
		<Field.Input
			value={String(value ?? '')}
			oninput={(event) => (value = event.currentTarget.value)}
			type={inputType}
		/>
	</Field.Root>
{/if}
