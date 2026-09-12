<script lang="ts">
	import { Field } from '@ark-ui/svelte/field';
	import type { HTMLInputAttributes } from 'svelte/elements';

	type Props = HTMLInputAttributes & {
		label: string;
		value: string;
		disabled?: boolean;
		required?: boolean;
	};

	let { label, value = $bindable(''), disabled, required, ...input }: Props = $props();
</script>

<!-- The label sits inside the filled block, above the value. -->
<Field.Root {disabled} {required}>
	<Field.Label>
		{label}
		<Field.RequiredIndicator>*</Field.RequiredIndicator>
	</Field.Label>
	<Field.Input bind:value {...input} />
</Field.Root>
