import type { Reroute } from '@sveltejs/kit';
import { deLocalizeUrl } from '$lib/paraglide/runtime';

// This hook is needed for Paraglide's URL-based locale routing
// It strips the locale prefix from the URL so SvelteKit can find the right route
export const reroute: Reroute = (request) => {
	return deLocalizeUrl(request.url).pathname;
};
