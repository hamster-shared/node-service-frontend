import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    outDir: "../pkg/controller/dist",
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
          //"https://console-mock.apipost.cn/mock/ae73cd30-20d8-4975-b034-48b34891e956/", //接口域名 //接口域名
          // "http://172.16.10.13:8080",
          // "http://172.16.31.103:8080",
          // "http://localhost:8080",
          "http://183.66.65.207:38080/",
        changeOrigin: true, //是否跨域
        // rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
});
