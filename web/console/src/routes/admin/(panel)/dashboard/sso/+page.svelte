<script lang="ts">
	import {
		RiBuilding2Line,
		RiGlobalLine,
		RiLinksLine,
		RiShieldKeyholeLine,
		RiTeamLine,
		RiUserFollowLine
	} from 'svelte-remixicon';
	import { Icon } from '$lib/components/ui';

	/** What the section will do, so its page says something useful until it does. */
	const planned = [
		{
			icon: RiShieldKeyholeLine,
			title: 'SAML 2.0 and OpenID Connect',
			text: "Connect an organisation's own identity provider — Okta, Microsoft Entra ID, Google Workspace — with either protocol."
		},
		{
			icon: RiGlobalLine,
			title: 'Sign-in by email domain',
			text: 'Send people to their company’s provider from the address they type, with no extra button to find.'
		},
		{
			icon: RiUserFollowLine,
			title: 'Just-in-time provisioning',
			text: 'Create a user the first time they sign in through their provider, with the fields it sends mapped onto theirs.'
		},
		{
			icon: RiTeamLine,
			title: 'Group to role mapping',
			text: 'Give users global and application roles from the groups their provider says they belong to.'
		}
	];
</script>

<svelte:head><title>SSO integrations · xermess admin</title></svelte:head>

<div class="page">
	<header>
		<div class="title">
			<h1>SSO integrations</h1>
			<span class="soon">Coming soon</span>
		</div>
		<p>Let organisations sign their people in through their own identity provider.</p>
	</header>

	<section class="hero">
		<span class="mark" aria-hidden="true">
			<Icon icon={RiBuilding2Line} size="1.6rem" />
			<span class="link"><Icon icon={RiLinksLine} size="0.9rem" /></span>
		</span>

		<div>
			<h2>Enterprise single sign-on is under development</h2>
			<p>
				Integrations are not available yet, and nothing on this page changes how anyone signs in.
				Users keep signing in with their email and password, and applications keep working as they
				do today.
			</p>
		</div>
	</section>

	<h2 class="planned-heading">What is planned</h2>

	<ul class="planned">
		{#each planned as item (item.title)}
			<li>
				<span class="icon"><Icon icon={item.icon} size="1.1rem" /></span>
				<div>
					<h3>{item.title}</h3>
					<p>{item.text}</p>
				</div>
			</li>
		{/each}
	</ul>
</div>

<style>
	.page {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
		max-width: 60rem;
		padding-inline: var(--page-gutter);
		padding-bottom: var(--space-5);
	}

	.title {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: var(--space-2);
	}

	header p {
		margin: var(--space-1) 0 0;
		color: var(--color-text-hint);
		font-size: var(--text-base);
	}

	.soon {
		display: inline-flex;
		align-items: center;
		height: 22px;
		padding: 0 9px;
		border: 1px solid color-mix(in srgb, var(--color-info), transparent 60%);
		border-radius: var(--radius-pill);
		background: var(--surface-info);
		color: color-mix(in srgb, var(--color-info) 80%, var(--color-text));
		font-size: var(--text-xs);
		font-weight: 600;
		letter-spacing: 0.02em;
	}

	.hero {
		display: flex;
		align-items: center;
		gap: var(--space-4);
		padding: var(--space-5) var(--space-4);
		border: 1px dashed var(--color-border);
		border-radius: var(--radius-lg);
		background: linear-gradient(
			135deg,
			color-mix(in srgb, var(--color-info), transparent 94%),
			transparent 60%
		);
	}

	.mark {
		position: relative;
		display: grid;
		flex-shrink: 0;
		place-items: center;
		width: 64px;
		height: 64px;
		border-radius: 18px;
		background: var(--color-primary);
		color: var(--color-primary-text);
	}

	.link {
		position: absolute;
		right: -6px;
		bottom: -6px;
		display: grid;
		place-items: center;
		width: 26px;
		height: 26px;
		border: 3px solid var(--color-surface);
		border-radius: var(--radius-pill);
		background: var(--color-info);
		color: #fff;
	}

	.hero h2 {
		margin: 0;
		font-size: var(--text-xl);
	}

	.hero p {
		margin: var(--space-1) 0 0;
		max-width: 40rem;
		color: var(--color-text-hint);
		line-height: 1.55;
	}

	.planned-heading {
		margin: var(--space-2) 0 0;
		font-size: var(--text-lg);
	}

	.planned {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: var(--space-3);
		margin: 0;
		padding: 0;
		list-style: none;
	}

	.planned li {
		display: flex;
		gap: var(--space-3);
		padding: var(--space-3) var(--space-4);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		background: var(--color-surface);
	}

	.icon {
		display: grid;
		flex-shrink: 0;
		place-items: center;
		width: 34px;
		height: 34px;
		border-radius: var(--radius-md);
		background: var(--color-secondary);
	}

	h3 {
		margin: 0;
		font-size: var(--text-base);
		font-weight: 600;
	}

	.planned p {
		margin: 3px 0 0;
		color: var(--color-text-hint);
		font-size: var(--text-sm);
		line-height: 1.5;
	}

	@media (max-width: 48rem) {
		.planned {
			grid-template-columns: minmax(0, 1fr);
		}

		.hero {
			flex-direction: column;
			align-items: flex-start;
		}
	}
</style>
