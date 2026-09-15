import { describe, expect, it } from 'vitest';
import { safeNext } from './next';

describe('safeNext', () => {
	it('keeps a path on this site', () => {
		expect(safeNext('/security?tab=sessions')).toBe('/security?tab=sessions');
	});

	it('refuses anywhere else', () => {
		for (const next of [
			'https://evil.example',
			'//evil.example',
			'/\\evil.example',
			'security',
			''
		]) {
			expect(safeNext(next)).toBe('/');
		}
		expect(safeNext(null, '/applications')).toBe('/applications');
	});
});
