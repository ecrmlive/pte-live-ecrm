import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    port: 5183,
    proxy: {
      "/api": {
        // C 端唯一目标服务：api-business（Compose 对外端口 18082）。
        target: "http://127.0.0.1:18082",
        changeOrigin: true,
      },
    },
  },
});
