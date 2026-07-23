<script lang="ts" setup>
import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { ElIcon, ElMessage, ElMessageBox } from 'element-plus';
import { CloseBold, Plus } from '@element-plus/icons-vue';

import PlugsApi from '#/api/core/plugs';

import PlugAddModal from './plug-add-modal.vue';
import type { PlugCategory, PlugItem } from './types';

import '#/assets/font/iconfont.css';

const loading = ref(true);
const categories = ref<PlugCategory[]>([]);
const addModalOpen = ref(false);
const activeCategory = ref<null | PlugCategory>(null);

function resolvePlugIcon(icon?: null | string) {
  if (icon && icon.length > 0) {
    return icon;
  }
  return 'icon-chajian1';
}

async function loadCategories() {
  loading.value = true;
  try {
    const res = await PlugsApi.plugslist({}, true);
    categories.value =
      (res.data as { accessList?: PlugCategory[] })?.accessList ?? [];
  } catch {
    categories.value = [];
  } finally {
    loading.value = false;
  }
}

function openAddModal(category: PlugCategory) {
  activeCategory.value = category;
  addModalOpen.value = true;
}

async function handleDeletePlug(plug: PlugItem) {
  try {
    await ElMessageBox.confirm('删除后不可恢复，确认删除该记录吗?', '提示', {
      cancelButtonText: '取消',
      confirmButtonText: '确定',
      type: 'warning',
    });
  } catch {
    return;
  }

  loading.value = true;
  try {
    const res = await PlugsApi.deleteplugs({ plus_id: plug.access_id }, true);
    if (res.code === 1) {
      ElMessage.success(res.msg || '删除成功');
      await loadCategories();
    }
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadCategories();
});
</script>

<template>
  <Page auto-content-height>
    <div v-loading="loading" class="space-y-8">
      <section
        v-for="category in categories"
        :key="category.plus_category_id"
        class="space-y-4"
      >
        <h3 class="text-base font-semibold text-foreground">
          {{ category.name }}
        </h3>

        <div
          class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
        >
          <div
            v-for="plug in category.children"
            :key="plug.access_id"
            class="group relative rounded-lg border border-border bg-card p-4 transition-colors hover:border-primary/40"
          >
            <button
              v-access:code="'platform:plugs:delete'"
              class="absolute right-2 top-2 rounded p-1 text-destructive opacity-0 transition hover:bg-destructive/10 group-hover:opacity-100"
              title="删除插件"
              type="button"
              @click="handleDeletePlug(plug)"
            >
              <ElIcon :size="16">
                <CloseBold />
              </ElIcon>
            </button>

            <div class="flex items-start gap-3">
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
            </div>
          </div>

          <button
            v-access:code="'platform:plugs:add'"
            class="flex min-h-[76px] items-center gap-3 rounded-lg border border-dashed border-primary/40 bg-card/50 p-4 text-left transition hover:border-primary hover:bg-accent/40"
            type="button"
            @click="openAddModal(category)"
          >
            <span
              class="flex size-10 shrink-0 items-center justify-center rounded-lg border border-primary text-primary"
            >
              <ElIcon :size="18">
                <Plus />
              </ElIcon>
            </span>
            <span class="text-sm text-muted-foreground">
              添加插件到此类别下
            </span>
          </button>
        </div>
      </section>

      <div
        v-if="!loading && categories.length === 0"
        class="rounded-lg border border-dashed border-border py-16 text-center text-sm text-muted-foreground"
      >
        暂无插件分类
      </div>
    </div>

    <PlugAddModal
      v-model:open="addModalOpen"
      :category="activeCategory"
      @success="loadCategories"
    />
  </Page>
</template>

<style scoped>
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
