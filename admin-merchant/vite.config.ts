import path from 'node:path';

import { defineConfig } from '@vben/vite-config';

import ElementPlus from 'unplugin-element-plus/vite';

import { srcAssetsPlugin } from './vite-plugin-src-assets';

export default defineConfig(async () => {
  return {
    application: {},
    vite: {
      // 商户后台由 platform.qxkejiwl.top/shop/ 反代。使用相对资源路径，
      // 避免 /js、/jse 被同域的平台后台根应用抢占而导致白屏。
      base: './',
      plugins: [
        srcAssetsPlugin(),
        ElementPlus({
          format: 'esm',
        }),
      ],
      server: {
        port: 11525,
        host: true,
      },
    },
  };
});
