<script lang="ts">
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import type { RouteId } from '$app/types';
	import type { ComponentType } from 'svelte';
	import {
		RiDashboardLine,
		RiGroupLine,
		RiKey2Line,
		RiShieldUserLine,
		RiWebhookLine
	} from 'svelte-remixicon';
	import { Icon } from '$lib/components/ui';

	/** SvelteKit's own route ids. They carry the (panel) layout group, which
	    resolve() strips when turning them into a URL. */
	type Route = RouteId;

	type Item = {
		route: Route;
		label: string;
		icon: ComponentType;
		/** Shown on the right of the row: a count, or a word like "soon". */
		hint?: string;
		/** Marks a section that is still placeholder data. */
		demo?: boolean;
	};

	type Group = { label: string; items: Item[] };

	const groups: Group[] = [
		{
			label: 'General',
			items: [{ route: '/admin/(panel)/dashboard', label: 'Overview', icon: RiDashboardLine }]
		},
		{
			label: 'Access',
			items: [
				{
					route: '/admin/(panel)/dashboard/admins',
					label: 'Administrators',
					icon: RiGroupLine,
					hint: '5',
					demo: true
				},
				{
					route: '/admin/(panel)/dashboard/roles',
					label: 'Roles',
					icon: RiShieldUserLine,
					hint: '4',
					demo: true
				}
			]
		},
		{
			label: 'Integrations',
			items: [
				{
					route: '/admin/(panel)/dashboard/api-keys',
					label: 'API keys',
					icon: RiKey2Line,
					hint: '3',
					demo: true
				},
				{
					route: '/admin/(panel)/dashboard/webhooks',
					label: 'Webhooks',
					icon: RiWebhookLine,
					hint: '3',
					demo: true
				}
			]
		}
	];

	/** Overview is the section's own page, so it only matches exactly; the
	    others also match anything below them. */
	function isCurrent(route: Route): boolean {
		const href = resolve(route);
		const path = page.url.pathname;

		return route === '/admin/(panel)/dashboard'
			? path === href
			: path === href || path.startsWith(`${href}/`);
	}
</script>

<aside>
	<nav aria-label="Dashboard sections">
		{#each groups as group (group.label)}
			<div class="group">
				<h2>{group.label}</h2>

				{#each group.items as item (item.route)}
					<a
						href={resolve(item.route)}
						class:current={isCurrent(item.route)}
						aria-current={isCurrent(item.route) ? 'page' : undefined}
					>
						<Icon icon={item.icon} />
						<span class="label">{item.label}</span>
						{#if item.hint}
							<span class="hint">{item.hint}</span>
						{/if}
					</a>
				{/each}
			</div>
		{/each}
	</nav>

	<p class="note">Sections marked with placeholder data are not wired to the API yet.</p>
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
		padding: var(--space-4) var(--space-2);
		border-right: 1px solid var(--color-border);
		background: var(--color-surface);
		overflow-y: auto;
	}

	nav {
		display: flex;
		flex-direction: column;
		gap: var(--space-4);
	}

	.group {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	h2 {
		padding: 0 var(--space-2) var(--space-1);
		color: var(--color-text-hint);
		font-size: var(--text-xs);
		font-weight: 600;
		letter-spacing: 0.04em;
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
	}

	.hint {
		color: var(--color-text-hint);
		font-size: var(--text-sm);
		font-family: var(--font-mono);
	}

	.note {
		padding: var(--space-2);
		color: var(--color-text-hint);
		font-size: var(--text-xs);
		line-height: 1.4;
	}

	/* Narrow screens have no room for a column, so the sections become one
	   scrollable row above the content: the group labels and the footnote
	   would only take space from it. */
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

		nav {
			flex-direction: row;
			gap: var(--space-1);
		}

		.group {
			flex-direction: row;
			gap: var(--space-1);
		}

		h2,
		.note,
		.hint {
			display: none;
		}

		a {
			white-space: nowrap;
		}
	}
</style>
