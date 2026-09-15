<script lang="ts">
	import type { ComponentType } from 'svelte';
	import Icon from './Icon.svelte';

	type Props = {
		/** An icon to show. */
		icon?: ComponentType;
		/** Or a few letters, such as someone's initials. */
		text?: string;
		size?: 'sm' | 'md';
	};

	let { icon, text, size = 'sm' }: Props = $props();
</script>

<!-- PocketBase's .thumb: a square on the faint tint, in the hint colour, for
     what leads a row or a card. -->
<span class="thumb {size}" aria-hidden="true">
	{#if icon}
		<Icon {icon} size={size === 'md' ? '1.15rem' : '1rem'} />
	{:else if text}
		{text}
	{/if}
</span>

<style>
	.thumb {
		display: inline-flex;
		flex-shrink: 0;
		align-items: center;
		justify-content: center;
		aspect-ratio: 1;
		border: 1px solid var(--color-secondary-alt);
		border-radius: var(--radius-sm);
		background: var(--color-surface-alt);
		color: var(--color-text-hint);
		font-size: 11px;
		font-weight: 600;
		letter-spacing: 0.02em;
	}

	.sm {
		width: 34px;
	}

	.md {
		width: var(--control-height);
		font-size: var(--text-sm);
	}
</style>
