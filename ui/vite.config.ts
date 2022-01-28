import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      // with options
      '/api': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
    }
  },
  build: {
    outDir: '../build/static',
  }
})
