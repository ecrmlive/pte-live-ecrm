import type { DirectiveBinding } from 'vue';

import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { buildShopPathAuthMap } from '#/utils/shop-path-auth';

const DEFAULT_IMG =
  'https://cos.qxkejiwl.top/pte-live/image/picture/picture_01.png';

function resolveImgUrl(binding: DirectiveBinding) {
  if (binding.value instanceof Object) {
    const jsonStr = binding.expression.split(',')[0].replace(/'/g, '');
    return checkChild(binding.value, jsonStr);
  }
  return binding.value;
}

function checkChild(obj: Record<string, unknown>, str: string): unknown {
  const arr = str.match(/\./g);
  if (!arr) {
    return obj[str];
  }
  const index = str.indexOf('.');
  const key = str.substring(0, index);
  const child = obj[key];
  if (child && typeof child === 'object') {
    return checkChild(child as Record<string, unknown>, str.substring(index + 1));
  }
  return null;
}

function applyImgUrl(el: HTMLElement, binding: DirectiveBinding) {
  const imgURL = resolveImgUrl(binding);
  const resolved = imgURL ? resolveCosMediaUrl(String(imgURL).trim()) : '';
  if (!resolved) {
    el.onerror = null;
    el.setAttribute('src', DEFAULT_IMG);
    return;
  }
  el.onerror = () => {
    el.onerror = null;
    el.setAttribute('src', DEFAULT_IMG);
  };
  el.setAttribute('src', resolved);
}

function applyAuth(el: HTMLElement, binding: DirectiveBinding) {
  const auth = buildShopPathAuthMap();
  const value = String(binding.value || '').toLowerCase();
  if (auth[value] !== true) {
    el.style.display = 'none';
  }
}

/** 商户业务指令：v-img-url · v-auth */
export function registerMerchantDirectives(app: import('vue').App) {
  app.directive('img-url', {
    mounted(el, binding) {
      applyImgUrl(el as HTMLElement, binding);
    },
    updated(el, binding) {
      applyImgUrl(el as HTMLElement, binding);
    },
  });

  app.directive('auth', {
    mounted(el, binding) {
      applyAuth(el as HTMLElement, binding);
    },
    updated(el, binding) {
      applyAuth(el as HTMLElement, binding);
    },
  });
}

/** @deprecated 使用 registerMerchantDirectives */
export const registerLegacyDirectives = registerMerchantDirectives;
