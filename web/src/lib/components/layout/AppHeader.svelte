<script lang="ts">
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import { RiFileListLine, RiShieldKeyholeLine, RiDashboardLine } from 'svelte-remixicon';
	import type { Admin } from '$lib/api';
	import { Icon, LinkButton, ThemeToggle } from '$lib/components/ui';
	import { can } from '$lib/permissions';
	import AccountMenu from './AccountMenu.svelte';

	type Props = { admin: Admin };

	let { admin }: Props = $props();

	// The top-level areas. Everything inside the dashboard has its own sidebar,
	// and anything to do with this account lives in the menu on the right.
	// The logs are only offered to an administrator whose roles allow reading
	// them.
	const links = $derived([
		{ href: resolve('/admin/dashboard'), label: 'Dashboard', icon: RiDashboardLine },
		...(can(admin, 'activity.read')
			? [{ href: resolve('/admin/logs'), label: 'Logs', icon: RiFileListLine }]
			: [])
	]);

	/** A link is current when the page is it or sits below it. */
	function isCurrent(href: string): boolean {
		return page.url.pathname === href || page.url.pathname.startsWith(`${href}/`);
	}
</script>

<header>
	<a class="brand" href={resolve('/admin/dashboard')}>
		<span class="mark">
			<Icon icon={RiShieldKeyholeLine} size="1.125rem" />
		</span>
		<strong>xermess</strong>
	</a>

	<!-- The sections are links that look like buttons, so they are the same
	     control as everything else in the bar, at the same size: the bar is
	     55px tall, so its controls take the smaller step and leave air above
	     and below. The current one is the one thing in the bar that is fully
	     dark. -->
	<nav aria-label="Sections">
		{#each links as link (link.href)}
			<LinkButton
				href={link.href}
				icon={link.icon}
				size="sm"
				variant={isCurrent(link.href) ? 'solid' : 'ghost'}
				aria-current={isCurrent(link.href) ? 'page' : undefined}
			>
				<span class="label">{link.label}</span>
			</LinkButton>
		{/each}
	</nav>

	<div class="account">
		<ThemeToggle size="sm" />
		<AccountMenu {admin} size="sm" />
	</div>
</header>

<style>
	/* The bar stays put while the page scrolls, so the sections and the
	   account menu are always one click away on a long table. */
	header {
		position: sticky;
		top: 0;
		z-index: 10;
		display: flex;
		align-items: center;
		gap: var(--space-4);
		height: var(--header-height);
		padding: 0 var(--space-4);
		border-bottom: 1px solid var(--color-border);
		background: var(--color-surface);
	}

	.brand {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		font-size: var(--text-lg);
		text-decoration: none;
		color: var(--color-text);
	}

	.mark {
		display: grid;
		place-items: center;
		width: 28px;
		height: 28px;
		border-radius: var(--radius-sm);
		background: var(--color-accent);
		color: var(--color-accent-text);
	}

	/* The sections sit with the logo on the left; the account menu is pushed
	   to the far end by its own auto margin. */
	nav {
		display: flex;
		align-items: center;
		gap: var(--space-1);
	}

	.account {
		margin-left: auto;
		display: flex;
		align-items: center;
		gap: var(--space-1);
	}

	@media (max-width: 55rem) {
		header {
			gap: var(--space-2);
			padding: 0 var(--space-2);
		}

		.brand strong {
			display: none;
		}

		nav :global(.control) {
			padding: 0 var(--space-3);
		}
	}

	/* Below this the labels do not fit beside the account menu, so the
	   sections become their icons. The label stays in the accessible name. */
	@media (max-width: 30rem) {
		nav :global(.control) {
			width: var(--control-height-sm);
			padding: 0;
		}

		.label {
			position: absolute;
			width: 1px;
			height: 1px;
			overflow: hidden;
			clip-path: inset(50%);
			white-space: nowrap;
		}
	}
</style>
