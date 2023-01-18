import { NODE_ENV, SERVER_URL, TEST_SERVER_URL } from '$env/static/private';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = () => {
	if (NODE_ENV == 'development') {
		return {
			serverUrl: TEST_SERVER_URL
		};
	}
	return {
		serverUrl: SERVER_URL
	};
};
