<script lang="ts">
	import { goto, invalidateAll } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { Menu } from '@ark-ui/svelte/menu';
	import { RiArrowDownSLine, RiLogoutBoxRLine, RiUserSettingsLine } from 'svelte-remixicon';
	import { adminApi, type Admin } from '$lib/api';
	import { Icon } from '$lib/components/ui';

	type Props = { admin: Admin };

	let { admin }: Props = $props();

	let signingOut = $state(false);

	/** The first letter of the name, which is enough to tell accounts apart. */
	const monogram = $derived((admin.full_name || admin.username).charAt(0).toUpperCase());

	async function open(value: string) {
		if (value === 'profile') {
			await goto(resolve('/admin/profile'));
			return;
		}

		await signOut();
	}

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

<Menu.Root
	positioning={{ placement: 'bottom-end', gutter: 6 }}
	onSelect={(details) => open(details.value)}
>
	<Menu.Trigger class="trigger">
		<span class="monogram" aria-hidden="true">{monogram}</span>
		<span class="name">{admin.username}</span>
		<Icon icon={RiArrowDownSLine} />
	</Menu.Trigger>

	<Menu.Positioner>
		<Menu.Content>
			<div class="identity">
				<strong>{admin.full_name}</strong>
				<span class="hint">{admin.email}</span>
				<span class="roles">{admin.roles.join(', ')}</span>
			</div>

			<Menu.Separator />

			<Menu.Item value="profile">
				<Icon icon={RiUserSettingsLine} />
				Profile
			</Menu.Item>

			<Menu.Item value="sign-out">
				<Icon icon={RiLogoutBoxRLine} />
				{signingOut ? 'Signing out…' : 'Sign out'}
			</Menu.Item>
		</Menu.Content>
	</Menu.Positioner>
</Menu.Root>

<style>
	:global([data-scope='menu'][data-part='trigger'].trigger) {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		height: 35px;
		padding: 0 var(--space-2);
		border: none;
		border-radius: var(--radius-sm);
		background: transparent;
		color: var(--color-text);
		font: inherit;
		font-size: var(--text-base);
		cursor: pointer;
		transition: background-color var(--speed-fast);
	}

	:global([data-scope='menu'][data-part='trigger'].trigger:hover),
	:global([data-scope='menu'][data-part='trigger'].trigger[data-state='open']) {
		background: var(--color-secondary);
	}

	.monogram {
		display: grid;
		place-items: center;
		width: 24px;
		height: 24px;
		border-radius: var(--radius-pill);
		background: var(--color-accent);
		color: var(--color-accent-text);
		font-size: var(--text-sm);
		font-weight: 700;
	}

	.identity {
		display: flex;
		flex-direction: column;
		gap: 2px;
		padding: var(--space-2);
	}

	.hint,
	.roles {
		color: var(--color-text-hint);
		font-size: var(--text-sm);
	}

	.roles {
		font-family: var(--font-mono);
	}

	@media (max-width: 40rem) {
		.name {
			display: none;
		}
	}
</style>
