<script lang="ts">
	import type { Snippet } from 'svelte';
	import { RiErrorWarningLine } from 'svelte-remixicon';
	import Icon from './Icon.svelte';

	type Props = { children: Snippet };

	let { children }: Props = $props();
</script>

<!-- role="alert" so a screen reader announces the message when it appears. -->
<p class="alert" role="alert">
	<Icon icon={RiErrorWarningLine} size="1.0625rem" />
	<span>{@render children()}</span>
</p>

<style>
	.alert {
		display: flex;
		align-items: flex-start;
		gap: var(--space-2);
		padding: var(--space-2) var(--space-3);
		border-radius: var(--radius-md);
		background: var(--surface-danger);
		color: var(--color-danger);
		font-size: var(--text-base);
	}
</style>
