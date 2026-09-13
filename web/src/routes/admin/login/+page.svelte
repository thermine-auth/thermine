<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { RiShieldKeyholeLine } from 'svelte-remixicon';
	import { adminApi, ApiError } from '$lib/api';
	import { Alert, Button, Card, Icon, Input, PasswordInput, ThemeToggle } from '$lib/components/ui';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let submitting = $state(false);

	const canSubmit = $derived(email.trim() !== '' && password !== '' && !submitting);

	async function signIn(event: SubmitEvent) {
		event.preventDefault();
		if (!canSubmit) return;

		error = '';
		submitting = true;

		try {
			await adminApi.login(email.trim(), password);

			// The session changed, so anything already loaded is stale.
			await invalidateAll();
			await goto(resolve('/admin/dashboard'), { replaceState: true });
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Something went wrong';
			password = '';
			submitting = false;
		}
	}
</script>

<svelte:head>
	<title>Sign in · xermess admin</title>
</svelte:head>

<main>
	<div class="corner">
		<ThemeToggle />
	</div>

	<div class="panel">
		<Card padded>
			<form onsubmit={signIn}>
				<header>
					<span class="mark">
						<Icon icon={RiShieldKeyholeLine} size="1.25rem" label="xermess" />
					</span>
					<h1>xermess</h1>
					<p class="muted">Sign in to the admin panel</p>
				</header>

				{#if error}
					<Alert>{error}</Alert>
				{/if}

				<Input
					label="Email"
					bind:value={email}
					type="email"
					disabled={submitting}
					required
					autocomplete="email"
					autocapitalize="none"
					spellcheck={false}
				/>

				<PasswordInput label="Password" bind:value={password} disabled={submitting} required />

				<Button type="submit" size="lg" loading={submitting} disabled={!canSubmit}>
					{submitting ? 'Signing in…' : 'Sign in'}
				</Button>
			</form>
		</Card>
	</div>
</main>

<style>
	main {
		display: grid;
		place-items: center;
		min-height: 100dvh;
		padding: var(--space-5);
	}

	.corner {
		position: fixed;
		top: var(--space-3);
		right: var(--space-3);
		color: var(--color-text-hint);
	}

	.panel {
		width: 100%;
		max-width: 22rem;
	}

	form {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
		padding: var(--space-2);
	}

	header {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}

	.mark {
		display: grid;
		place-items: center;
		width: 40px;
		height: 40px;
		margin-bottom: var(--space-2);
		border-radius: var(--radius-md);
		background: var(--color-accent);
		color: var(--color-accent-text);
	}

	header h1 {
		font-size: var(--text-lg);
	}

	.muted {
		color: var(--color-text-hint);
		font-size: var(--text-base);
	}
</style>
