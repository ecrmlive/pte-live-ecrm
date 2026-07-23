<script lang="ts" setup>
import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { getPlusCenterApi } from '#/api/core/plus';
import { openPlusPluginPage } from '#/utils/plus-navigation';
import {
  MERCHANT_AGENT_HUB_PATH,
  MERCHANT_ARTICLE_HUB_PATH,
} from '#/utils/qixi-live-menu';

import type { PlusCategory, PlusPluginItem } from './types';

import { useRouter } from 'vue-router';

import '#/assets/font/iconfont.css';

defineOptions({ name: 'MerchantPlusCenter' });

/** Tab 容器插件：始终进入 hub path，勿用 redirect_name（可能指向导出等子路由） */
const PLUS_TAB_HUB_ENTRY_PATHS = new Set<string>([
  MERCHANT_AGENT_HUB_PATH,
  MERCHANT_ARTICLE_HUB_PATH,
]);

/** 与 sql/scripts/gen_merchant_access.py CHILD_ORDER 对齐 */
const PLUGIN_PATH_ORDER = [
  '/plus/points/index',
  '/plus/sign',
  '/plus/homepush/index',
  '/plus/invitation/active/index',
  '/plus/coupon/index',
  '/plus/package/index',
  '/plus/agent/index',
  '/plus/article',
  '/plus/collection/index',
  '/plus/officia/index',
  '/plus/recommend/index',
  '/plus/fullfree/index',
  '/plus/fullreduce/index',
  '/plus/lottery/index',
  '/plus/table/event',
  '/plus/card/event',
  '/plus/fullreduce/product',
  '/plus/advance/index',
  '/plus/task/index',
  '/plus/surface/index',
  '/plus/buyactivity/index',
  '/plus/register/index',
  '/plus/newactivity/index',
  '/plus/seckill/index',
  '/plus/assemble/index',
  '/plus/bargain/index',
];

const router = useRouter();
const loading = ref(true);
const loadError = ref('');
const categories = ref<PlusCategory[]>([]);

function resolvePlugIcon(icon?: null | string) {
  if (icon && icon.length > 0) {
    return icon;
  }
  return 'icon-chajian1';
}

function sortPluginChildren(children: PlusPluginItem[]) {
  const order = new Map(PLUGIN_PATH_ORDER.map((path, index) => [path, index]));
  return [...children].sort((a, b) => {
    const ai = order.get(a.path) ?? 999;
    const bi = order.get(b.path) ?? 999;
    if (ai !== bi) return ai - bi;
    return a.name.localeCompare(b.name, 'zh-CN');
  });
}

function sortCategories(groups: PlusCategory[]) {
  return groups.map((group) => ({
    ...group,
    children: sortPluginChildren(group.children ?? []),
  }));
}

function gotoPlugin(item: PlusPluginItem) {
  const path = String(item.path || '').trim();
  const redirect = String(item.redirect_name || '').trim();
  const target =
    path && PLUS_TAB_HUB_ENTRY_PATHS.has(path) ? path : redirect || path;
  if (target) {
    void openPlusPluginPage(router, target);
  }
}

async function loadCategories() {
  loading.value = true;
  loadError.value = '';
  try {
    const data = await getPlusCenterApi();
    categories.value = sortCategories(data?.list ?? []);
  } catch (error) {
    categories.value = [];
    loadError.value =
      error instanceof Error ? error.message : '加载插件中心失败，请稍后重试';
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadCategories();
});
</script>

<template>
  <Page auto-content-height class="plus-center-page">
    <div v-loading="loading" class="plus-center-content space-y-8">
      <section
        v-for="(category, index) in categories"
        :key="category.plus_category_id ?? index"
        class="space-y-4"
      >
        <h3 class="text-base font-semibold text-foreground">
          {{ category.name }}
        </h3>

        <div
          class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
        >
          <button
            v-for="plug in category.children"
            :key="plug.access_id ?? plug.path"
            class="group flex min-h-[76px] items-start gap-3 rounded-lg border border-border bg-card p-4 text-left transition-colors hover:border-primary/40 hover:bg-accent/30"
            type="button"
            @click="gotoPlugin(plug)"
          >
            <span
              class="iconfont plug-icon flex shrink-0 items-center justify-center"
              :class="resolvePlugIcon(plug.icon)"
            />
            <div class="min-w-0 flex-1 pt-0.5">
              <h4
                class="truncate text-sm font-medium text-foreground group-hover:text-primary"
              >
                {{ plug.name }}
              </h4>
              <p class="mt-1 line-clamp-2 text-xs text-muted-foreground">
                {{ plug.remark || '—' }}
              </p>
            </div>
          </button>
        </div>
      </section>

      <div
        v-if="!loading && loadError"
        class="rounded-lg border border-dashed border-destructive/40 py-16 text-center text-sm text-destructive"
      >
        {{ loadError }}
      </div>

      <div
        v-else-if="!loading && categories.length === 0"
        class="rounded-lg border border-dashed border-border py-16 text-center text-sm text-muted-foreground"
      >
        暂无已安装插件，请联系平台管理员开通
      </div>
    </div>
  </Page>
</template>

<style scoped>
.plus-center-page :deep(.plus-center-content) {
  isolation: isolate;
}

.plug-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: hsl(var(--primary));
  color: #fff;
  font-size: 22px;
  line-height: 40px;
}
</style>
