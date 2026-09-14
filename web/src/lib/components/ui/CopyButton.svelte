<script lang="ts">
	import { RiCheckLine, RiFileCopyLine } from 'svelte-remixicon';
	import IconButton from './IconButton.svelte';

	type Props = {
		/** What is copied. */
		value: string;
		/** What it is, for the button's name: "Copy client id". */
		label?: string;
	};

	let { value, label = 'value' }: Props = $props();

	let copied = $state(false);
	let timer: ReturnType<typeof setTimeout>;

	async function copy() {
		await navigator.clipboard.writeText(value);

		copied = true;
		clearTimeout(timer);
		timer = setTimeout(() => (copied = false), 1500);
	}
</script>

<!-- The icon turns into a tick for a moment, which is all the confirmation a
     copy needs. -->
<IconButton
	icon={copied ? RiCheckLine : RiFileCopyLine}
	label={copied ? 'Copied' : `Copy ${label}`}
	size="sm"
	placement="left"
	onclick={copy}
/>
