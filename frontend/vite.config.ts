import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import { join } from "path";
import legacy from '@vitejs/plugin-legacy'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    legacy()
  ],
  resolve: {
    alias: {
      '@': join(__dirname, "src")
    },
  }
})
