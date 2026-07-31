import { defineConfig } from '@vben/vite-config';

import ElementPlus from 'unplugin-element-plus/vite';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      plugins: [
        ElementPlus({
          format: 'esm',
        }),
      ],
      server: {
        port: 15124,
        host: true,
        proxy: {
          '/api/platform': {
            changeOrigin: true,
            target: 'http://127.0.0.1:18081',
          },
        },
      },
    },
  };
});
