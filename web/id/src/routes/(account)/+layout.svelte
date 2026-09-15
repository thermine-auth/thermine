<script lang="ts">
	import type { Snippet } from 'svelte';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { signIn } from '$lib/api';
	import { Brand, Icon, ThemeToggle, type IconName } from '$lib/components';
	import { initials } from '$lib/utils/format';
	import type { LayoutProps } from './$types';

	let { data, children }: LayoutProps & { children: Snippet } = $props();

	const sections: { href: '/' | '/security' | '/applications'; label: string; icon: IconName }[] = [
		{ href: '/', label: 'Profile', icon: 'user' },
		{ href: '/security', label: 'Security', icon: 'lock' },
		{ href: '/applications', label: 'Connected apps', icon: 'apps' }
	];

	let signingOut = $state(false);

	async function signOut() {
		signingOut = true;
		try {
			await signIn.logout();
		} finally {
			// One navigation that reloads the data: invalidating first would
			// make this layout's load redirect to sign in, racing this goto.
			await goto(resolve('/login'), { invalidateAll: true, replaceState: true });
		}
	}

	const name = $derived(`${data.user.first_name} ${data.user.last_name}`.trim() || data.user.email);
</script>

<div class="frame">
	<header>
		<div class="bar">
			<a class="home" href={resolve('/')} aria-label="Your account"><Brand /></a>

			<div class="end">
				<ThemeToggle />
				<div class="who">
					<span class="avatar" aria-hidden="true">{initials(data.user)}</span>
					<span class="text">
						<strong>{name}</strong>
						<span>{data.user.email}</span>
					</span>
				</div>
				<button class="signout" type="button" onclick={signOut} disabled={signingOut}>
					<Icon name="logout" />
					<span>Sign out</span>
				</button>
			</div>
		</div>

		<nav aria-label="Account">
			{#each sections as section (section.href)}
				<a
					href={resolve(section.href)}
					aria-current={page.url.pathname === section.href ? 'page' : undefined}
				>
					<Icon name={section.icon} size="1rem" />
					{section.label}
				</a>
			{/each}
		</nav>
	</header>

	<main>
		{@render children()}
	</main>

	<footer>
		<Icon name="shield" size="0.9375rem" />
		<span>Secured by <strong>xermess</strong></span>
	</footer>
</div>

<style>
	.frame {
		display: grid;
		grid-template-rows: auto 1fr auto;
		min-height: 100dvh;
	}

	header {
		position: sticky;
		top: 0;
		z-index: 10;
		border-bottom: 1px solid var(--color-border);
		background: color-mix(in srgb, var(--color-surface), transparent 8%);
		backdrop-filter: blur(10px);
	}

	.bar,
	nav,
	main {
		width: 100%;
		max-width: 56rem;
		margin: 0 auto;
		padding-left: var(--space-5);
		padding-right: var(--space-5);
	}

	.bar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: var(--space-3);
		height: 64px;
	}

	.home:hover {
		text-decoration: none;
	}

	.end {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		min-width: 0;
	}

	.who {
		display: flex;
		align-items: center;
		gap: var(--space-2);
		min-width: 0;
		padding: 0 var(--space-2);
	}

	.avatar {
		display: grid;
		place-items: center;
		flex-shrink: 0;
		width: 34px;
		height: 34px;
		border-radius: 50%;
		background: var(--color-accent);
		color: var(--color-accent-text);
		font-size: var(--text-sm);
		font-weight: 700;
	}

	.text {
		display: flex;
		flex-direction: column;
		min-width: 0;
		line-height: 1.25;
	}

	.text strong,
	.text span {
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.text span {
		color: var(--color-text-hint);
		font-size: var(--text-sm);
	}

	.signout {
		display: inline-flex;
		align-items: center;
		gap: var(--space-2);
		height: 38px;
		padding: 0 var(--space-3);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-md);
		background: var(--color-surface);
		font-weight: 600;
		cursor: pointer;
	}

	.signout:hover {
		background: var(--color-input);
	}

	nav {
		display: flex;
		gap: var(--space-1);
		overflow-x: auto;
		scrollbar-width: none;
	}

	nav a {
		display: inline-flex;
		align-items: center;
		gap: var(--space-2);
		padding: var(--space-3) var(--space-3);
		border-bottom: 2px solid transparent;
		color: var(--color-text-hint);
		font-weight: 600;
		white-space: nowrap;
	}

	nav a:hover {
		color: var(--color-text);
		text-decoration: none;
	}

	nav a[aria-current='page'] {
		border-bottom-color: var(--color-primary);
		color: var(--color-text);
	}

	main {
		padding-top: var(--space-6);
		padding-bottom: var(--space-7);
	}

	footer {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-1);
		padding: var(--space-5) 0;
		color: var(--color-text-hint);
		font-size: var(--text-sm);
	}

	@media (max-width: 40rem) {
		.bar,
		nav,
		main {
			padding-left: var(--space-4);
			padding-right: var(--space-4);
		}

		.text {
			display: none;
		}

		.signout span {
			display: none;
		}
	}
</style>
