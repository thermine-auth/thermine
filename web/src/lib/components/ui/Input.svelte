<script lang="ts">
	import type { ComponentType } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';
	import { Field as ArkField } from '@ark-ui/svelte/field';
	import Field from './Field.svelte';

	type Props = Omit<HTMLInputAttributes, 'value' | 'size'> & {
		label: string;
		value: string;
		icon?: ComponentType;
		hint?: string;
		error?: string;
		required?: boolean;
		disabled?: boolean;
		readOnly?: boolean;
	};

	let {
		label,
		value = $bindable(''),
		icon,
		hint,
		error,
		required,
		disabled,
		readOnly,
		...input
	}: Props = $props();
</script>

<Field {label} {icon} {hint} {error} {required} {disabled} {readOnly}>
	<ArkField.Input bind:value readonly={readOnly} {...input} />
</Field>
