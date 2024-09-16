import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from "path"

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@common": path.resolve(__dirname, "../common/ts"),
    },
  },
  define: {
      "import.meta.env.SERVER_URL": `"${process.env.SERVER_URL}"`,
  },
  server: {
    port: 3000,
  },
})

