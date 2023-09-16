import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import { join } from "path";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': join(__dirname, "src")
    },
  }
})
