import { redirect } from '@sveltejs/kit';
import { resolve } from '$app/paths';
import type { PageLoad } from './$types';

/** There is nothing at the root yet: the admin panel is the app. */
export const load: PageLoad = () => {
	redirect(307, resolve('/admin/overview'));
};
