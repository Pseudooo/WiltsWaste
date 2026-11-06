/** @type {import("@sveltejs/vite-plugin-svelte").SvelteConfig} */
export default {
  // ... other config
  kit: {
    // ... other config
    alias: {
      "@/*": "./path/to/lib/*",
    },
  },
};
