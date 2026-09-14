<script lang="ts">
	import type { ComponentType, Snippet } from 'svelte';
	import { RiArrowRightSLine } from 'svelte-remixicon';
	import Icon from './Icon.svelte';
	import Thumb from './Thumb.svelte';

	type Props = {
		label: string;
		value: number;
		icon: ComponentType;
		/** Labels under the number that give it context. */
		children?: Snippet;
		/** Where the number leads, when the administrator may go there. */
		href?: string;
	};

	let { label, value, icon, children, href }: Props = $props();
</script>

<svelte:element this={href ? 'a' : 'div'} class="stat" class:link={href} {href}>
	<span class="lead"><Thumb {icon} size="md" /></span>

	<span class="body">
		<span class="label">{label}</span>
		<strong>{value.toLocaleString()}</strong>
		{#if children}<span class="tags">{@render children()}</span>{/if}
	</span>

	{#if href}<span class="go" aria-hidden="true"><Icon icon={RiArrowRightSLine} /></span>{/if}
</svelte:element>

<style>
	.stat {
		display: flex;
		align-items: flex-start;
		gap: var(--space-3);
		min-width: 0;
		padding: var(--space-3);
		border: 1px solid var(--color-secondary-alt);
		border-radius: var(--radius-sm);
		background: var(--color-surface);
		color: var(--color-text);
		text-decoration: none;
		transition:
			background-color var(--speed),
			border-color var(--speed);
	}

	.link:hover,
	.link:focus-visible {
		outline: 0;
		border-color: var(--color-border);
		background: var(--color-surface-alt);
	}

	.body {
		display: flex;
		flex: 1;
		flex-direction: column;
		gap: 2px;
		min-width: 0;
	}

	.label {
		color: var(--color-text-hint);
		font-size: var(--text-sm);
	}

	strong {
		font-size: 1.286rem;
		font-weight: 600;
		line-height: 24px;
		font-variant-numeric: tabular-nums;
	}

	.tags {
		display: flex;
		flex-wrap: wrap;
		gap: 5px;
		margin-top: 6px;
	}

	.go {
		align-self: center;
		color: var(--color-text-disabled);
		transition: color var(--speed);
	}

	.link:hover .go {
		color: var(--color-text);
	}

	@media (max-width: 34rem) {
		.stat {
			gap: var(--space-2);
			padding: var(--space-2);
		}

		.lead,
		.go {
			display: none;
		}
	}
</style>
