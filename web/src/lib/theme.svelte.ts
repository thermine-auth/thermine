import { browser } from '$app/environment';

/** The panel is either light or dark. */
export type Theme = 'light' | 'dark';

/** Where the choice is kept, and the attribute the stylesheet reads. The
 *  inline script in app.html uses the same two names to apply the theme
 *  before the first paint; change them together. */
const STORAGE_KEY = 'xermess-theme';
const ATTRIBUTE = 'data-theme';

function isTheme(value: unknown): value is Theme {
	return value === 'light' || value === 'dark';
}

/** The theme to start with: what was chosen last time, or what the system
 *  prefers for someone who has not chosen yet. */
function initial(): Theme {
	if (!browser) return 'light';

	try {
		const saved = localStorage.getItem(STORAGE_KEY);
		if (isTheme(saved)) return saved;
	} catch {
		// Private browsing can make localStorage throw on read.
	}

	return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

class ThemeState {
	current = $state<Theme>(initial());

	/** Switches to the other theme. */
	toggle() {
		this.set(this.current === 'dark' ? 'light' : 'dark');
	}

	/** Applies a theme: the attribute drives the stylesheet, and the stored
	 *  value is what the inline script reads on the next load. */
	set(theme: Theme) {
		this.current = theme;

		if (!browser) return;

		document.documentElement.setAttribute(ATTRIBUTE, theme);

		try {
			localStorage.setItem(STORAGE_KEY, theme);
		} catch {
			// Not being able to remember the choice is not worth an error.
		}
	}
}

export const theme = new ThemeState();
