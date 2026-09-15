<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Icon, ThemeToggle } from '$lib/components';

	let { children }: { children: Snippet } = $props();
</script>

<div class="shell">
	<div class="corner"><ThemeToggle /></div>

	<main>
		<div class="column">
			{@render children()}
		</div>
	</main>

	<footer>
		<Icon name="shield" size="0.9375rem" />
		<span>Secured by <strong>xermess</strong></span>
	</footer>
</div>

<style>
	.shell {
		display: grid;
		grid-template-rows: 1fr auto;
		min-height: 100dvh;
		padding: 0 var(--space-4);
		background:
			radial-gradient(
				60rem 28rem at 50% -8rem,
				color-mix(in srgb, var(--color-accent), transparent 88%),
				transparent
			),
			var(--color-body);
	}

	.corner {
		position: fixed;
		top: var(--space-3);
		right: var(--space-3);
	}

	main {
		display: grid;
		place-items: center;
		padding: var(--space-7) 0 var(--space-5);
	}

	.column {
		width: 100%;
		max-width: 27rem;
	}

	footer {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-1);
		padding: var(--space-4) 0 var(--space-5);
		color: var(--color-text-hint);
		font-size: var(--text-sm);
	}

	/* On a phone the page itself is the surface. */
	@media (max-width: 30rem) {
		.shell {
			background: var(--color-surface);
		}

		main {
			place-items: start center;
			padding-top: calc(var(--space-7) + var(--space-3));
		}
	}
</style>
