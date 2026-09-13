<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { HTMLButtonAttributes } from 'svelte/elements';

	type Props = HTMLButtonAttributes & {
		/** `solid` for the one action a screen is about, `secondary` for the
		    rest, `ghost` for icon buttons in a bar, `danger` for confirming
		    something that cannot be undone. */
		variant?: 'solid' | 'secondary' | 'ghost' | 'danger';
		/** `sm` for buttons that sit inside a bar or a row. */
		size?: 'md' | 'sm';
		children: Snippet;
	};

	let { variant = 'solid', size = 'md', type = 'button', children, ...rest }: Props = $props();
</script>

<button {type} class="{variant} {size}" {...rest}>
	{@render children()}
</button>

<style>
	button {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-2);
		min-height: var(--control-height);
		padding: 0 var(--space-4);
		border: 1px solid transparent;
		border-radius: var(--radius-md);
		font-family: var(--font-sans);
		font-size: var(--text-base);
		font-weight: 600;
		white-space: nowrap;
		cursor: pointer;
		transition:
			background-color var(--speed-fast),
			border-color var(--speed-fast),
			opacity var(--speed-fast);
	}

	button:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.sm {
		min-height: 32px;
		padding: 0 var(--space-3);
		font-size: var(--text-sm);
	}

	.solid {
		background: var(--color-primary);
		color: var(--color-primary-text);
	}

	.solid:hover:not(:disabled) {
		background: color-mix(in srgb, var(--color-primary), white 12%);
	}

	.secondary {
		background: var(--color-secondary);
		color: var(--color-text);
	}

	.secondary:hover:not(:disabled) {
		background: var(--color-secondary-alt);
	}

	.danger {
		background: var(--color-danger);
		color: #fff;
	}

	.danger:hover:not(:disabled) {
		background: color-mix(in srgb, var(--color-danger), black 10%);
	}

	.ghost {
		min-height: 35px;
		padding: 0 var(--space-2);
		background: transparent;
		color: inherit;
		font-weight: 500;
	}

	.ghost:hover:not(:disabled) {
		background: color-mix(in srgb, currentcolor, transparent 88%);
	}
</style>
