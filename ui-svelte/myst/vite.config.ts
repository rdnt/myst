import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import sveltePreprocess from 'svelte-preprocess';
import * as sass from 'sass';
import * as path from "path";

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    port: 8082,
  },
  plugins: [
    svelte({
      preprocess: sveltePreprocess({
        scss: true,
        typescript: true,
        sass: {
          sync: true,
          implementation: sass,
        },
      }),
    }),
  ],
  resolve: {
    alias: {
      '@': path.resolve('/src'),
    },
  },
  optimizeDeps: {exclude: ["svelte-navigator"]},
})
