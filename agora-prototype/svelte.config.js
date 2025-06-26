import adapter from '@sveltejs/adapter-auto';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://kit.svelte.dev/docs/integrations#preprocessors
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	kit: {
		// adapter-auto only supports some environments, see https://kit.svelte.dev/docs/adapter-auto for a list.
		// If your environment is not supported or you settled on a specific environment, switch out the adapter.
		// See https://kit.svelte.dev/docs/adapters for more information about adapters.
		adapter: adapter(),

		// Configure paths for the agora prototype
		paths: {
			base: process.env.NODE_ENV === 'production' ? '/agora' : '',
		},

		// App configuration
		appDir: '.svelte-kit',

		// Files configuration
		files: {
			assets: 'static',
			hooks: {
				client: 'src/hooks.client.js',
				server: 'src/hooks.server.js'
			},
			lib: 'src/lib',
			params: 'src/params',
			routes: 'src/routes',
			serviceWorker: 'src/service-worker.js',
			appTemplate: 'src/app.html',
			errorTemplate: 'src/error.html'
		},

		// Version configuration
		version: {
			name: process.env.npm_package_version
		},

		// CSP configuration for security
		csp: {
			mode: 'auto',
			directives: {
				'script-src': ['self', 'unsafe-inline'],
				'style-src': ['self', 'unsafe-inline', 'fonts.googleapis.com'],
				'font-src': ['self', 'fonts.gstatic.com'],
				'img-src': ['self', 'data:', 'https:'],
				'connect-src': ['self', 'ws://localhost:*', 'wss://localhost:*']
			}
		}
	}
};

export default config;
