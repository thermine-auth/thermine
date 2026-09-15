import { signIn, type Application } from '$lib/api';
import type { PageServerLoad } from './$types';

/** The application the user signed out of, when the logout named one, so the
    page can offer a way back to it. */
export const load: PageServerLoad = async ({ url, fetch }) => {
	const clientId = url.searchParams.get('client_id');
	let application: Application | null = null;

	if (clientId) {
		try {
			({ application } = await signIn.application(clientId, fetch));
		} catch {
			// Not knowing the application only means no link back.
		}
	}

	return { application };
};
