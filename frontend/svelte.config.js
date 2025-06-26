import adapter from "@sveltejs/adapter-auto";
import { vitePreprocess } from "@sveltejs/kit/vite";

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

    // Configure paths
    paths: {
      base: process.env.NODE_ENV === "production" ? "/hi" : "",
    },

    // Configure aliases for cleaner imports
    alias: {
      $components: "src/lib/components",
      $stores: "src/lib/stores",
      $utils: "src/lib/utils",
      $types: "src/lib/types",
      $api: "src/lib/api",
    },

    // Configure CSP for security - DISABLED temporarily for development
    // csp: {
    // 	mode: 'auto',
    // 	directives: {
    // 		'script-src': ['self', 'unsafe-inline'],
    // 		'style-src': ['self', 'unsafe-inline'],
    // 		'img-src': ['self', 'data:', 'https:'],
    // 		'font-src': ['self', 'https:'],
    // 		'connect-src': ['self', 'ws:', 'wss:']
    // 	}
    // },

    // Version handling
    version: {
      name: process.env.npm_package_version,
    },
  },
};

export default config;
