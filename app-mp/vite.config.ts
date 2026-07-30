import { defineConfig } from "vite";
import uni from "@dcloudio/vite-plugin-uni";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [uni()],
  server: {
    port: 5173,
    proxy: {
      "/api": {
        // H5 与小程序统一转发到 C 端目标服务 api-business。
        target: "http://127.0.0.1:18082",
        changeOrigin: true,
      },
    },
  },
});
