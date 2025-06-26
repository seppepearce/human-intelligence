import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [sveltekit()],

  // Development server configuration
  server: {
    port: 3000,
    strictPort: false,
    host: true, // Listen on all addresses

    // Proxy API calls to Go backend
    proxy: {
      "/api": {
        target: "http://localhost:8081",
        changeOrigin: true,
        secure: false,
      },
      "/ws": {
        target: "ws://localhost:8081",
        ws: true,
        changeOrigin: true,
      },
    },

    // Hot reload configuration
    hmr: {
      port: 3001,
    },
  },

  // Preview server (for production builds)
  preview: {
    port: 3000,
    strictPort: false,
    host: true,
  },

  // Build configuration
  build: {
    target: "esnext",
    sourcemap: true,

    // Optimize chunks
    rollupOptions: {
      output: {
        manualChunks: {
          // Vendor chunks
          "vendor-ui": ["lucide-svelte", "@tabler/icons-svelte"],
          "vendor-utils": ["clsx", "tailwind-merge", "zod"],
          "vendor-data": ["date-fns", "marked"],
          "vendor-viz": ["d3"],
          "vendor-math": ["katex"],
          "vendor-highlight": ["highlight.js"],
          "vendor-realtime": ["socket.io-client"],
        },
      },
    },

    // Optimize dependencies
    commonjsOptions: {
      include: [/node_modules/],
    },
  },

  // Dependency optimization
  optimizeDeps: {
    include: [
      "lucide-svelte",
      "@tabler/icons-svelte",
      "clsx",
      "tailwind-merge",
      "d3",
      "socket.io-client",
      "marked",
      "highlight.js",
      "katex",
    ],
  },

  // CSS configuration
  css: {
    postcss: "./postcss.config.js",
  },

  // Define global constants
  define: {
    __APP_VERSION__: JSON.stringify(process.env.npm_package_version),
    __BUILD_TIME__: JSON.stringify(new Date().toISOString()),
  },

  // Environment variables
  envPrefix: ["VITE_", "HI_"],

  // Worker configuration for potential background tasks
  worker: {
    format: "es",
  },
});
