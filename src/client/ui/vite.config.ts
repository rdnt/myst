import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import * as path from "path";

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    outDir: '../cmd/client/static'
  },
  server: {
    port: 8082,
  },
  plugins: [svelte()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  optimizeDeps: {exclude: ["svelte-navigator"]},
})
