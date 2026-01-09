import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  base: '/auth/_/',
  plugins: [react(),tailwindcss(),],
  server: {
    proxy: {
      '/api': {
        target: 'https://localhost:8947/auth',
        changeOrigin: true,
        secure: false,
      },
      '/admin': {
        target: 'https://localhost:8947/auth',
        changeOrigin: true,
        secure: false,
      },
    },
  },
})