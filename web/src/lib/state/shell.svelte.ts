import { getContext, setContext } from 'svelte';
import { rememberSidebar, type SidebarState } from './sidebar';

/**
 * Whether the sidebar is folded, shared by everything drawn beside it.
 *
 * The header's logo block and the dashboard's sidebar are one column of the
 * screen, so they have to agree on how wide it is and fold together. The
 * panel layout makes this once per page view and hands it down as context —
 * never as module state, which on the server would be shared by every
 * request.
 */
export type Shell = {
	readonly collapsed: boolean;
	toggle: () => void;
};

const key = Symbol('shell');

/** Makes the shared state, starting as the cookie the server read says. */
export function provideShell(initial: SidebarState): Shell {
	let collapsed = $state(initial === 'mini');

	const shell: Shell = {
		get collapsed() {
			return collapsed;
		},
		toggle() {
			collapsed = !collapsed;
			rememberSidebar(collapsed ? 'mini' : 'wide');
		}
	};

	return setContext(key, shell);
}

/** The shared state, from a component inside the panel layout. */
export function useShell(): Shell {
	const shell = getContext<Shell | undefined>(key);
	if (!shell) throw new Error('useShell() needs the panel layout above it');

	return shell;
}
