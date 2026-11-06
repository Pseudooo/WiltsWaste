import { svelte } from '@sveltejs/vite-plugin-svelte'
import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

export default defineConfig({
  plugins: [svelte(), tailwindcss()],
  resolve: {
    alias: {
      $lib: path.resolve(__dirname, './src/lib'),
    },
  },
})