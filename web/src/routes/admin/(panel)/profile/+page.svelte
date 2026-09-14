<script lang="ts">
	import {
		RiComputerLine,
		RiGlobalLine,
		RiLockPasswordLine,
		RiMailLine,
		RiPaletteLine,
		RiShieldCheckLine,
		RiShieldKeyholeLine,
		RiUserSettingsLine
	} from 'svelte-remixicon';
	import SessionList from '$lib/components/profile/SessionList.svelte';
	import SignOutButton from '$lib/components/profile/SignOutButton.svelte';
	import {
		Button,
		List,
		ListItem,
		PageContainer,
		PageHeader,
		Panel,
		Tag,
		Thumb
	} from '$lib/components/ui';
	import { theme } from '$lib/state/theme.svelte';
	import { formatDateTime, formatRelative } from '$lib/utils/format';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const admin = $derived(data.admin);

	/** Two letters for the account's thumb. */
	const initials = $derived(
		admin.full_name
			.split(/\s+/)
			.filter(Boolean)
			.slice(0, 2)
			.map((part) => part[0])
			.join('')
			.toUpperCase() || admin.email.slice(0, 2).toUpperCase()
	);

	const activeSessions = $derived(data.sessions.filter((session) => session.active).length);
</script>

<svelte:head><title>Profile · xermess admin</title></svelte:head>

<PageContainer>
	<div class="page">
		<PageHeader crumbs={['Account', 'Profile']}>
			{#snippet actions()}
				<SignOutButton />
			{/snippet}
		</PageHeader>

		<!-- Who is signed in, at a glance. -->
		<List bordered label="Account">
			<ListItem title={admin.full_name} description={admin.email}>
				{#snippet lead()}<Thumb text={initials} size="md" />{/snippet}
				{#snippet end()}
					<Tag tone={admin.status === 'active' ? 'success' : 'neutral'} dot strong>
						{admin.status}
					</Tag>
					{#if admin.is_super_admin}
						<Tag tone="info">super admin</Tag>
					{/if}
				{/snippet}
			</ListItem>
		</List>

		<Panel title="Profile information" icon={RiUserSettingsLine} flush>
			<List label="Profile information">
				<ListItem>
					<span class="detail"><span>Username</span><b>{admin.username}</b></span>
				</ListItem>
				<ListItem>
					<span class="detail"><span>Name</span><b>{admin.full_name}</b></span>
				</ListItem>
				<ListItem>
					<span class="detail">
						<span>Roles</span>
						<span class="tags">
							{#each admin.roles as role (role)}
								<Tag small>{role}</Tag>
							{:else}
								<b>—</b>
							{/each}
						</span>
					</span>
				</ListItem>
				<ListItem>
					<span class="detail">
						<span>Last signed in</span>
						{#if admin.last_login_at}
							<b title={formatDateTime(admin.last_login_at)}>
								{formatDateTime(admin.last_login_at)}
								<small>· {formatRelative(admin.last_login_at)}</small>
							</b>
						{:else}
							<b>Never</b>
						{/if}
					</span>
				</ListItem>
			</List>
		</Panel>

		<Panel title="Sign-in & security" icon={RiShieldKeyholeLine} flush>
			{#snippet meta()}<Tag small>Not available yet</Tag>{/snippet}
			<List label="Sign-in and security">
				<ListItem
					title="Email"
					description="{admin.email} · changing it asks for confirmation at the new address."
				>
					{#snippet lead()}<Thumb icon={RiMailLine} />{/snippet}
					{#snippet end()}<Button size="sm" variant="subtle" disabled>Change</Button>{/snippet}
				</ListItem>
				<ListItem title="Password" description="Set · changing it signs out every other session.">
					{#snippet lead()}<Thumb icon={RiLockPasswordLine} />{/snippet}
					{#snippet end()}<Button size="sm" variant="subtle" disabled>Change</Button>{/snippet}
				</ListItem>
				<ListItem
					title="Two-factor authentication"
					description="Not set up · an authenticator app or a security key, as a second step."
				>
					{#snippet lead()}<Thumb icon={RiShieldCheckLine} />{/snippet}
					{#snippet end()}<Button size="sm" variant="subtle" disabled>Set up</Button>{/snippet}
				</ListItem>
			</List>
		</Panel>

		<Panel title="Preferences" icon={RiPaletteLine} flush>
			<List label="Preferences">
				<ListItem title="Theme" description="Light or dark, remembered on this device.">
					{#snippet lead()}<Thumb icon={RiPaletteLine} />{/snippet}
					{#snippet end()}
						<span class="choices" role="group" aria-label="Theme">
							{#each ['light', 'dark'] as const as option (option)}
								<button
									type="button"
									class="control choice"
									data-size="sm"
									data-variant="outline"
									data-palette="neutral"
									class:selected={theme.current === option}
									aria-pressed={theme.current === option}
									onclick={() => theme.set(option)}
								>
									<span class="swatch {option}"></span>
									{option === 'light' ? 'Light' : 'Dark'}
								</button>
							{/each}
						</span>
					{/snippet}
				</ListItem>
				<ListItem
					title="Language"
					description="English · Kyrgyz and Russian are translated but not offered here yet."
				>
					{#snippet lead()}<Thumb icon={RiGlobalLine} />{/snippet}
					{#snippet end()}
						<Tag small>Soon</Tag>
						<Button size="sm" variant="subtle" disabled>Change</Button>
					{/snippet}
				</ListItem>
			</List>
		</Panel>

		<Panel title="Sessions" icon={RiComputerLine} flush>
			{#snippet meta()}
				<Tag tone={activeSessions > 0 ? 'success' : 'neutral'} dot>
					{activeSessions} active
				</Tag>
			{/snippet}
			<SessionList sessions={data.sessions} />
		</Panel>
	</div>
</PageContainer>

<style>
	.page {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}

	/* One fact about the account: its name in the hint colour on the left,
	   the value beside it. */
	.detail {
		display: grid;
		grid-template-columns: 10rem minmax(0, 1fr);
		align-items: center;
		gap: var(--space-3);
		color: var(--color-text-hint);
		font-size: var(--text-sm);
	}

	.detail b {
		color: var(--color-text);
		font-size: var(--text-base);
		font-weight: normal;
	}

	.detail small {
		color: var(--color-text-hint);
		font-size: var(--text-sm);
	}

	.tags {
		display: flex;
		flex-wrap: wrap;
		gap: 5px;
	}

	.choices {
		display: flex;
		gap: 5px;
	}

	/* The shape is the shared control; being the chosen one is this page's
	   own business. */
	.choice.selected {
		border-color: var(--color-text);
	}

	.swatch {
		width: 14px;
		height: 14px;
		border: 1px solid var(--color-border);
		border-radius: 4px;
	}

	.swatch.light {
		background: #fff;
	}

	.swatch.dark {
		background: #1c1c1c;
	}

	@media (max-width: 34rem) {
		.detail {
			grid-template-columns: 1fr;
			gap: 2px;
		}
	}
</style>
