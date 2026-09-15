<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { ComponentType } from 'svelte';
	import { Field } from '@ark-ui/svelte/field';
	import Icon from './Icon.svelte';

	type Props = {
		label: string;
		/** The control itself: an input, a select, a textarea. */
		children: Snippet;
		/** Drawn before the label, to say what kind of value this is. */
		icon?: ComponentType;
		/** A line under the control, for what someone needs to know before
		    filling it in. */
		hint?: string;
		/** A line under the control, for what went wrong. It replaces the
		    hint and marks the field as invalid. */
		error?: string;
		required?: boolean;
		disabled?: boolean;
		readOnly?: boolean;
	};

	let { label, children, icon, hint, error, required, disabled, readOnly }: Props = $props();
</script>

<!-- Every field in the panel is this: a label inside the filled block, the
     control under it, and one line of help or of blame. The parts are Ark's,
     so they are styled once in styles/ark.css and every field that uses this
     looks the same without saying so. -->
<Field.Root {required} {disabled} {readOnly} invalid={error !== undefined}>
	<Field.Label>
		{#if icon}
			<Icon {icon} />
		{/if}
		{label}
		{#if required}
			<Field.RequiredIndicator>*</Field.RequiredIndicator>
		{/if}
	</Field.Label>

	{@render children()}

	{#if error}
		<Field.ErrorText>{error}</Field.ErrorText>
	{:else if hint}
		<Field.HelperText>{hint}</Field.HelperText>
	{/if}
</Field.Root>
