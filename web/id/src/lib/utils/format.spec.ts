import { describe, expect, it } from 'vitest';
import { describeDevice, initials, timeAgo } from './format';

describe('describeDevice', () => {
	it('names common browsers and systems', () => {
		const chromeMac =
			'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0 Safari/537.36';
		const safariPhone =
			'Mozilla/5.0 (iPhone; CPU iPhone OS 18_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Mobile/15E148 Safari/604.1';

		expect(describeDevice(chromeMac)).toBe('Chrome on macOS');
		expect(describeDevice(safariPhone)).toBe('Safari on iPhone');
		expect(describeDevice('')).toBe('Unknown browser on an unknown device');
	});
});

describe('timeAgo', () => {
	it('says how long ago', () => {
		const now = Date.parse('2026-09-15T12:00:00Z');
		expect(timeAgo('2026-09-15T11:59:30Z', now)).toBe('just now');
		expect(timeAgo('2026-09-15T11:55:00Z', now)).toBe('5 minutes ago');
		expect(timeAgo('2026-09-12T12:00:00Z', now)).toBe('3 days ago');
	});
});

describe('initials', () => {
	it('uses the name, or the email', () => {
		expect(initials({ first_name: 'Ada', last_name: 'Lovelace', email: 'a@x' })).toBe('AL');
		expect(initials({ first_name: '', last_name: '', email: 'grace@x' })).toBe('G');
	});
});
