<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { RiLogoutBoxRLine, RiShieldKeyholeLine } from 'svelte-remixicon';
	import { adminApi, type Admin } from '$lib/api';
	import { Button, Icon, ThemeToggle } from '$lib/components/ui';

	type Props = { admin: Admin };

	let { admin }: Props = $props();

	let signingOut = $state(false);

	async function signOut() {
		signingOut = true;

		try {
			await adminApi.logout();
		} finally {
			// However the server answered, this browser is done with the
			// session: drop what was loaded with it and go to the sign-in page.
			await invalidateAll();
			await goto(resolve('/admin/login'), { replaceState: true });
		}
	}
</script>

<!-- The navy bar is PocketBase's: the one strongly coloured surface in the
     app, which keeps the content area calm. -->
<header>
	<div class="brand">
		<span class="mark">
			<Icon icon={RiShieldKeyholeLine} size="1.125rem" />
		</span>
		<strong>xermess</strong>
		<span class="hint">admin</span>
	</div>

	<div class="account">
		<ThemeToggle />
		<span class="hint">{admin.username}</span>
		<Button variant="ghost" onclick={signOut} disabled={signingOut}>
			<Icon icon={RiLogoutBoxRLine} />
			{signingOut ? 'Signing out…' : 'Sign out'}
		</Button>
	</div>
</header>

<style>
	header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-4);
		height: var(--header-height);
		padding: 0 var(--space-4);
		background: var(--color-accent);
		color: var(--color-accent-text);
	}

	.brand {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		font-size: var(--text-lg);
	}

	.mark {
		display: grid;
		place-items: center;
		width: 28px;
		height: 28px;
		border-radius: var(--radius-md);
		background: color-mix(in srgb, currentcolor, transparent 88%);
	}

	.account {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		font-size: var(--text-base);
	}

	.hint {
		color: color-mix(in srgb, currentcolor, transparent 35%);
	}
</style>
