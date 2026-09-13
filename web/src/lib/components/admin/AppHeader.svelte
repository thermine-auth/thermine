<script lang="ts">
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import { RiFileListLine, RiShieldKeyholeLine, RiDashboardLine } from 'svelte-remixicon';
	import type { Admin } from '$lib/api';
	import { Icon, ThemeToggle } from '$lib/components/ui';
	import AccountMenu from './AccountMenu.svelte';

	type Props = { admin: Admin };

	let { admin }: Props = $props();

	// The top-level areas. Everything inside the dashboard has its own sidebar,
	// and anything to do with this account lives in the menu on the right.
	const links = [
		{ href: resolve('/admin/dashboard'), label: 'Dashboard', icon: RiDashboardLine },
		{ href: resolve('/admin/logs'), label: 'Logs', icon: RiFileListLine }
	];

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

	<nav aria-label="Sections">
		{#each links as link (link.href)}
			<a
				href={link.href}
				class:current={isCurrent(link.href)}
				aria-current={isCurrent(link.href) ? 'page' : undefined}
			>
				<Icon icon={link.icon} />
				<span class="label">{link.label}</span>
			</a>
		{/each}
	</nav>

	<div class="account">
		<ThemeToggle />
		<AccountMenu {admin} />
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

	nav a {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		height: 35px;
		padding: 0 var(--space-3);
		border-radius: var(--radius-sm);
		color: var(--color-text-hint);
		font-size: var(--text-base);
		font-weight: 500;
		text-decoration: none;
		transition:
			background-color var(--speed-fast),
			color var(--speed-fast);
	}

	nav a:hover {
		background: var(--color-secondary);
		color: var(--color-text);
	}

	/* The current section is the one thing in the bar that is fully dark. */
	nav a.current {
		background: var(--color-primary);
		color: var(--color-primary-text);
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

		nav a {
			padding: 0 var(--space-2);
		}
	}

	/* Below this the labels do not fit beside the account menu, so the
	   sections become their icons. The label stays in the accessible name. */
	@media (max-width: 30rem) {
		nav a {
			width: 34px;
			justify-content: center;
			padding: 0;
		}

		nav a > .label {
			position: absolute;
			width: 1px;
			height: 1px;
			overflow: hidden;
			clip-path: inset(50%);
			white-space: nowrap;
		}
	}
</style>
