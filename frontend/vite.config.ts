import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import path from "path"

let dev_port_str = process.env.DEV_PORT ?? null;
let dev_port: number | undefined;
if (dev_port_str != null) {
  dev_port = parseInt(dev_port_str);
} 

// https://vitejs.dev/config/
export default defineConfig({
  server: {
      port: dev_port,
  },
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@common": path.resolve(__dirname, "../common/ts"),
    },
  },
})
