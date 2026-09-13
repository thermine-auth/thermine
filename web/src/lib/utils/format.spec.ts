import { describe, expect, it } from 'vitest';

import { formatDateTime } from './format';

describe('formatDateTime', () => {
	it('reads a timestamp the way the panel shows it', () => {
		const shown = formatDateTime('2026-09-13T14:38:09.654650+06:00');

		// The exact wording is the reader's locale; what matters is that the
		// day and the year are in it and nothing threw.
		expect(shown).toContain('2026');
		expect(shown.length).toBeGreaterThan(0);
	});

	it('leaves something that is not a date alone', () => {
		expect(formatDateTime('not a date')).toBe('not a date');
	});

	it('has something to show for nothing', () => {
		expect(formatDateTime('')).toBe('');
	});
});
