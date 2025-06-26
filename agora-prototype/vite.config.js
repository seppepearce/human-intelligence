import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],

	// Development server configuration
	server: {
		port: 3001,
		host: '0.0.0.0',
		strictPort: true,
		hmr: {
			port: 3002
		},
		// WebSocket proxy for agora real-time features
		proxy: {
			'/ws': {
				target: 'ws://localhost:8082',
				ws: true,
				changeOrigin: true
			},
			'/api': {
				target: 'http://localhost:8082',
				changeOrigin: true,
				rewrite: (path) => path.replace(/^\/api/, '')
			}
		}
	},

	// Preview server configuration
	preview: {
		port: 3003,
		host: '0.0.0.0',
		strictPort: true
	},

	// Build configuration
	build: {
		target: 'es2020',
		sourcemap: true,
		rollupOptions: {
			output: {
				manualChunks: {
					// Separate chunks for better caching
					'd3': ['d3', 'd3-selection', 'd3-force', 'd3-scale', 'd3-hierarchy'],
					'lucide': ['lucide-svelte']
				}
			}
		}
	},

	// Dependency optimization
	optimizeDeps: {
		include: [
			'd3',
			'd3-selection',
			'd3-force',
			'd3-scale',
			'd3-hierarchy',
			'lucide-svelte'
		]
	},

	// CSS configuration
	css: {
		postcss: './postcss.config.js',
		devSourcemap: true
	},

	// Define global constants
	define: {
		__APP_VERSION__: JSON.stringify(process.env.npm_package_version),
		__BUILD_TIME__: JSON.stringify(new Date().toISOString()),
		__DEV__: JSON.stringify(process.env.NODE_ENV === 'development')
	},

	// Testing configuration
	test: {
		include: ['src/**/*.{test,spec}.{js,ts}'],
		environment: 'jsdom',
		setupFiles: ['src/test-setup.js']
	},

	// Worker configuration for Web Workers if needed
	worker: {
		format: 'es'
	}
});
