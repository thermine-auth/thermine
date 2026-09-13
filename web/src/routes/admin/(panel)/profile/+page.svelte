<script lang="ts">
	import {
		RiGlobalLine,
		RiLockPasswordLine,
		RiLogoutBoxRLine,
		RiMailLine,
		RiPaletteLine,
		RiShieldCheckLine,
		RiUserSettingsLine
	} from 'svelte-remixicon';
	import SessionTable from '$lib/components/profile/SessionTable.svelte';
	import ProfileSection from '$lib/components/profile/ProfileSection.svelte';
	import SignOutButton from '$lib/components/profile/SignOutButton.svelte';
	import { Badge, Button, Card } from '$lib/components/ui';
	import { formatDateTime } from '$lib/utils/format';
	import { theme } from '$lib/state/theme.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	const details = $derived([
		{ label: 'Username', value: data.admin.username },
		{ label: 'Name', value: data.admin.full_name },
		{ label: 'Roles', value: data.admin.roles.join(', ') },
		{
			label: 'Last signed in',
			value: data.admin.last_login_at ? formatDateTime(data.admin.last_login_at) : 'Never'
		}
	]);
</script>

<svelte:head><title>Profile · xermess admin</title></svelte:head>

<header>
	<h1>Profile</h1>
	<p class="subtitle">Your account, how you sign in, and how this panel looks to you.</p>
</header>

<ProfileSection
	icon={RiUserSettingsLine}
	title="Profile information"
	description="Who you are in this organisation."
>
	<Card>
		<dl>
			{#each details as detail (detail.label)}
				<div class="row">
					<dt>{detail.label}</dt>
					<dd>{detail.value}</dd>
				</div>
			{/each}
			<div class="row">
				<dt>Status</dt>
				<dd>
					<Badge tone={data.admin.status === 'active' ? 'success' : 'neutral'}>
						{data.admin.status}
					</Badge>
				</dd>
			</div>
		</dl>
	</Card>
</ProfileSection>

<ProfileSection
	icon={RiShieldCheckLine}
	title="Two-factor authentication"
	description="A second step when you sign in."
	demo
>
	<Card padded>
		<div class="setting">
			<div>
				<strong>Not set up</strong>
				<p class="hint">An authenticator app or a security key can be added here.</p>
			</div>
			<Button variant="subtle" disabled>Set up</Button>
		</div>
	</Card>
</ProfileSection>

<ProfileSection
	icon={RiPaletteLine}
	title="Theme"
	description="Light or dark, remembered on this device."
>
	<Card padded>
		<div class="choices" role="group" aria-label="Theme">
			{#each ['light', 'dark'] as const as option (option)}
				<button
					type="button"
					class="control choice"
					data-size="md"
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
		</div>
	</Card>
</ProfileSection>

<ProfileSection
	icon={RiGlobalLine}
	title="Language"
	description="What this panel is shown in."
	demo
>
	<Card padded>
		<div class="setting">
			<div>
				<strong>English</strong>
				<p class="hint">Kyrgyz and Russian are translated but not offered here yet.</p>
			</div>
			<Button variant="subtle" disabled>Change</Button>
		</div>
	</Card>
</ProfileSection>

<ProfileSection icon={RiMailLine} title="Email" description="Where account notices are sent." demo>
	<Card padded>
		<div class="setting">
			<div>
				<strong>{data.admin.email}</strong>
				<p class="hint">Changing this will ask for confirmation at the new address.</p>
			</div>
			<Button variant="subtle" disabled>Change</Button>
		</div>
	</Card>
</ProfileSection>

<ProfileSection
	icon={RiLockPasswordLine}
	title="Password"
	description="Used with your username to sign in."
	demo
>
	<Card padded>
		<div class="setting">
			<div>
				<strong>Set</strong>
				<p class="hint">Changing it signs out every other session.</p>
			</div>
			<Button variant="subtle" disabled>Change</Button>
		</div>
	</Card>
</ProfileSection>

<ProfileSection
	icon={RiLogoutBoxRLine}
	title="Sessions"
	description="Where this account is signed in."
>
	<SessionTable sessions={data.sessions} />

	<div class="sign-out">
		<SignOutButton />
	</div>
</ProfileSection>

<style>
	header {
		margin-bottom: var(--space-5);
		padding-inline: var(--page-gutter);
	}

	.subtitle {
		margin-top: var(--space-1);
		color: var(--color-text-hint);
		font-size: var(--text-base);
	}

	dl {
		margin: 0;
	}

	.row {
		display: grid;
		grid-template-columns: 10rem 1fr;
		gap: var(--space-4);
		align-items: center;
		padding: var(--space-3) var(--space-4);
		font-size: var(--text-base);
	}

	.row + .row {
		border-top: 1px solid var(--color-border);
	}

	dt {
		color: var(--color-text-hint);
	}

	dd {
		margin: 0;
	}

	.setting {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-4);
	}

	.hint {
		margin-top: 2px;
		color: var(--color-text-hint);
		font-size: var(--text-sm);
	}

	.choices {
		display: flex;
		gap: var(--space-2);
	}

	/* The shape is the shared control; being the chosen one is this page's
	   own business. */
	.choice.selected {
		border-color: var(--color-text);
	}

	.swatch {
		width: 16px;
		height: 16px;
		border: 1px solid var(--color-border);
		border-radius: var(--radius-sm);
	}

	.swatch.light {
		background: #fff;
	}

	.swatch.dark {
		background: #1c1c1c;
	}

	.sign-out {
		margin-top: var(--space-3);
	}

	@media (max-width: 40rem) {
		.row {
			grid-template-columns: 1fr;
			gap: var(--space-1);
		}

		.setting {
			flex-direction: column;
			align-items: flex-start;
		}
	}
</style>
