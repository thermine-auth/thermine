<script lang="ts">
	import Icon from './Icon.svelte';
</script>

<span class="brand">
	<span class="mark"><Icon name="shield" size="1rem" /></span>
	<span>Account</span>
</span>

<style>
	.brand {
		display: inline-flex;
		align-items: center;
		gap: var(--space-2);
		font-weight: 700;
		font-size: var(--text-lg);
		color: var(--color-text);
	}

	.mark {
		display: grid;
		place-items: center;
		width: 30px;
		height: 30px;
		border-radius: 8px;
		background: var(--color-accent);
		color: var(--color-accent-text);
	}
</style>
