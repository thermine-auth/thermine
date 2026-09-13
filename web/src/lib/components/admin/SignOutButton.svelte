<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { RiLogoutBoxRLine } from 'svelte-remixicon';
	import { adminApi } from '$lib/api';
	import { Button, Icon } from '$lib/components/ui';

	type Props = { variant?: 'solid' | 'secondary' | 'ghost' };

	let { variant = 'secondary' }: Props = $props();

	let signingOut = $state(false);

	export async function signOut() {
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

<Button {variant} onclick={signOut} disabled={signingOut}>
	<Icon icon={RiLogoutBoxRLine} />
	{signingOut ? 'Signing out…' : 'Sign out'}
</Button>
