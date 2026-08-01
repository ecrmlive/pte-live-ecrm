import fs from 'node:fs';
import path from 'node:path';

import type { Plugin } from 'vite';

const ASSET_PREFIX = '/src/assets/';

/** DIY 模板 `/src/assets/*` → `src/assets/*`（缺失时用 placeholder） */
export function srcAssetsPlugin(): Plugin {
  const nativeAssetsRoot = path.resolve(import.meta.dirname, 'src/assets');
  const placeholder = path.resolve(import.meta.dirname, 'public/placeholder.png');

  return {
    enforce: 'pre',
    name: 'pte-live-src-assets',
    resolveId(source) {
      if (!source.startsWith(ASSET_PREFIX)) {
        return null;
      }
      const rel = source.slice(ASSET_PREFIX.length);
      const nativeFile = path.join(nativeAssetsRoot, rel);
      if (fs.existsSync(nativeFile)) {
        return nativeFile;
      }
      return placeholder;
    },
  };
}
