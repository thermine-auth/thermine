<script lang="ts">
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import type { RouteId } from '$app/types';
	import type { ComponentType } from 'svelte';
	import {
		RiAppsLine,
		RiBuildingLine,
		RiCodeBoxLine,
		RiDatabase2Line,
		RiGitBranchLine,
		RiGroupLine,
		RiLinksLine,
		RiPulseLine,
		RiShareLine,
		RiShieldUserLine,
		RiSidebarFoldLine,
		RiSidebarUnfoldLine,
		RiTranslate2
	} from 'svelte-remixicon';
	import { Icon, Tooltip } from '$lib/components/ui';

	type Props = {
		/** Folded to icons only. The width itself is set by the layout. */
		collapsed: boolean;
		onToggle: () => void;
	};

	let { collapsed, onToggle }: Props = $props();

	type Item = {
		route: RouteId;
		label: string;
		icon: ComponentType;
		/** Marks a section that is still placeholder data. */
		demo?: boolean;
	};

	/** A group with no label is a section on its own, shown above the rest. */
	type Group = { label?: string; items: Item[] };

	const groups: Group[] = [
		{
			items: [{ route: '/admin/(panel)/dashboard', label: 'Activity', icon: RiPulseLine }]
		},
		{
			label: 'Applications',
			items: [
				{
					route: '/admin/(panel)/dashboard/applications',
					label: 'Applications',
					icon: RiAppsLine,
					demo: true
				},
				{ route: '/admin/(panel)/dashboard/apis', label: 'APIs', icon: RiCodeBoxLine, demo: true },
				{
					route: '/admin/(panel)/dashboard/sso',
					label: 'SSO integrations',
					icon: RiLinksLine,
					demo: true
				}
			]
		},
		{
			label: 'Authentication',
			items: [
				{
					route: '/admin/(panel)/dashboard/database',
					label: 'Database',
					icon: RiDatabase2Line,
					demo: true
				},
				{
					route: '/admin/(panel)/dashboard/social',
					label: 'Social',
					icon: RiShareLine,
					demo: true
				},
				{
					route: '/admin/(panel)/dashboard/flows',
					label: 'Login flows',
					icon: RiGitBranchLine,
					demo: true
				}
			]
		},
		{
			label: 'User management',
			items: [
				{ route: '/admin/(panel)/dashboard/users', label: 'Users', icon: RiGroupLine },
				{
					route: '/admin/(panel)/dashboard/roles',
					label: 'Roles',
					icon: RiShieldUserLine,
					demo: true
				}
			]
		},
		{
			label: 'Settings',
			items: [
				{
					route: '/admin/(panel)/dashboard/organization',
					label: 'Organization',
					icon: RiBuildingLine,
					demo: true
				},
				{
					route: '/admin/(panel)/dashboard/languages',
					label: 'Languages',
					icon: RiTranslate2,
					demo: true
				}
			]
		}
	];

	/** Activity is the section's own page, so it only matches exactly; the
	    others also match anything below them. */
	function isCurrent(route: RouteId): boolean {
		const href = resolve(route);
		const path = page.url.pathname;

		return route === '/admin/(panel)/dashboard'
			? path === href
			: path === href || path.startsWith(`${href}/`);
	}
</script>

<aside class:mini={collapsed}>
	<nav aria-label="Sections">
		{#each groups as group (group.label ?? 'top')}
			<div class="group">
				{#if group.label}
					<h2>{group.label}</h2>
				{/if}

				{#each group.items as item (item.route)}
					<!-- Folded, the name is gone from the column, so the tooltip is
					     the only thing left saying what the icon leads to. Unfolded
					     it would only repeat the label, so it is switched off. -->
					<Tooltip label={item.label} placement="right" disabled={!collapsed}>
						{#snippet children(trigger)}
							<a
								{...trigger()}
								href={resolve(item.route)}
								class:current={isCurrent(item.route)}
								aria-current={isCurrent(item.route) ? 'page' : undefined}
							>
								<Icon icon={item.icon} />
								<span class="label">{item.label}</span>
								{#if item.demo}
									<span class="dot" aria-hidden="true"></span>
								{/if}
							</a>
						{/snippet}
					</Tooltip>
				{/each}
			</div>
		{/each}
	</nav>

	<div class="foot">
		<p class="note">
			<span class="dot" aria-hidden="true"></span>
			Sections marked this way show placeholder data.
		</p>

		<Tooltip label="Expand the sidebar" placement="right" disabled={!collapsed}>
			{#snippet children(trigger)}
				<button
					{...trigger()}
					type="button"
					class="fold"
					onclick={onToggle}
					aria-expanded={!collapsed}
					aria-label={collapsed ? 'Expand the sidebar' : 'Collapse the sidebar'}
				>
					<Icon icon={collapsed ? RiSidebarUnfoldLine : RiSidebarFoldLine} />
					<span class="label">Collapse</span>
				</button>
			{/snippet}
		</Tooltip>
	</div>
</aside>

<style>
	aside {
		position: sticky;
		top: var(--header-height);
		display: flex;
		flex-direction: column;
		justify-content: space-between;
		gap: var(--space-4);
		height: calc(100dvh - var(--header-height));
		padding: var(--space-3) var(--space-2) var(--space-2);
		border-right: 1px solid var(--color-border);
		background: var(--color-surface);
		overflow-x: hidden;
		overflow-y: auto;
	}

	nav {
		display: flex;
		flex-direction: column;
		gap: var(--space-3);
	}

	.group {
		display: flex;
		flex-direction: column;
		gap: 1px;
	}

	/* Nothing in the column re-wraps while it is folding: the width animates,
	   and text that no longer fits is clipped rather than reflowing into a
	   different shape on the way. */
	h2 {
		padding: var(--space-2) var(--space-2) var(--space-1);
		white-space: nowrap;
		color: var(--color-text-hint);
		font-size: var(--text-xs);
		font-weight: 600;
		letter-spacing: 0.06em;
		text-transform: uppercase;
	}

	a {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		height: 34px;
		padding: 0 var(--space-2);
		border-radius: var(--radius-sm);
		color: var(--color-text);
		font-size: var(--text-base);
		text-decoration: none;
		transition:
			background-color var(--speed-fast),
			color var(--speed-fast);
	}

	a:hover {
		background: var(--color-secondary);
	}

	a.current {
		background: var(--color-secondary-alt);
		font-weight: 600;
	}

	.label {
		flex: 1;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	/* A quiet dot rather than a word: it marks the section without competing
	   with its name. */
	.dot {
		flex: none;
		width: 5px;
		height: 5px;
		border-radius: var(--radius-pill);
		background: var(--color-text-hint);
		opacity: 0.5;
	}

	.note {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		width: 13rem;
		padding: var(--space-2);
		color: var(--color-text-hint);
		font-size: var(--text-xs);
		line-height: 1.4;
	}

	.foot {
		display: flex;
		flex-direction: column;
		gap: var(--space-1);
	}

	/* The fold button sits where a nav item would, so it reads as part of the
	   list rather than a control bolted underneath it. */
	.fold {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		height: 34px;
		padding: 0 var(--space-2);
		border: none;
		border-radius: var(--radius-sm);
		background: transparent;
		color: var(--color-text-hint);
		font: inherit;
		font-size: var(--text-base);
		text-align: left;
		cursor: pointer;
		transition:
			background-color var(--speed-fast),
			color var(--speed-fast);
	}

	.fold:hover {
		background: var(--color-secondary);
		color: var(--color-text);
	}

	/* Folded: icons only, centred, with the labels and the footnote gone. The
	   names come back as tooltips, which is what the title attributes are
	   for. */
	.mini h2,
	.mini .note,
	.mini .label {
		display: none;
	}

	.mini a,
	.mini .fold {
		justify-content: center;
		padding: 0;
	}

	/* The placeholder mark has no room beside the name any more, so it moves
	   to the corner of the icon. */
	.mini a {
		position: relative;
	}

	.mini .dot {
		position: absolute;
		top: 5px;
		right: 14px;
	}

	/* Narrow screens have no room for a column, so the sections become one
	   scrollable row above the content. */
	@media (max-width: 55rem) {
		aside {
			position: sticky;
			top: var(--header-height);
			z-index: 5;
			height: auto;
			padding: var(--space-2);
			border-right: none;
			border-bottom: 1px solid var(--color-border);
			overflow-x: auto;
			overscroll-behavior-x: contain;
		}

		nav,
		.group {
			flex-direction: row;
			gap: var(--space-1);
		}

		h2,
		.note,
		.fold,
		.dot {
			display: none;
		}

		a {
			white-space: nowrap;
		}
	}
</style>
