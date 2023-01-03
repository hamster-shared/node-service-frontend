import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  build: {
  },
  plugins: [vue()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    proxy: {
      "/api": {
        target:
          "http://34.232.105.81:8080",
          // "http://172.16.10.100:8080",
        changeOrigin: true, //是否跨域
        // rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
});
