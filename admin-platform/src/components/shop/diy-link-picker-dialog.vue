<script setup lang="ts">
import type { DiyPageCategory, DiyPageLink, DiyPageRow } from '#/api/core/diy';

import { ElButton, ElDialog, ElMessage, ElOption, ElRadioButton, ElRadioGroup, ElSelect, ElTable, ElTableColumn } from 'element-plus';
import { computed, ref, watch } from 'vue';

import { listDiyPageCategoriesApi, listDiyPageLinksApi, listDiyPagesApi } from '#/api/core/diy';

const props = defineProps<{ is_linkset?: boolean; linkData?: Record<string, unknown> | null }>();
const emit = defineEmits<{ closeDialog: [payload?: { name: string; type: string; url: string }] }>();
const open = ref(false);
const categories = ref<DiyPageCategory[]>([]);
const links = ref<DiyPageLink[]>([]);
const microPages = ref<DiyPageRow[]>([]);
const cateID = ref<number>();
const selected = ref<{ name: string; url: string }>();
const source = ref<'link' | 'micro'>('link');
const loading = ref(false);

function flatten(list: DiyPageCategory[], out: DiyPageCategory[] = []) { for (const item of list) { out.push(item); flatten(item.children || [], out); } return out; }
const categoryOptions = computed(() => flatten(categories.value));
const filteredLinks = computed(() => cateID.value ? links.value.filter((item) => item.cate_id === cateID.value) : links.value);

async function load() {
  loading.value = true;
  try {
    const [categoryRes, linkRes, microRes] = await Promise.all([
      listDiyPageCategoriesApi('platform'), listDiyPageLinksApi('platform', { page: 1, limit: 100 }),
      listDiyPagesApi({ is_diy: 0, status: 1, page: 1, limit: 100 }),
    ]);
    categories.value = categoryRes.list || [];
    links.value = linkRes.list || [];
    microPages.value = microRes.list || [];
  } finally { loading.value = false; }
}

watch(() => props.is_linkset, async (value) => {
  open.value = !!value;
  if (value) { selected.value = undefined; cateID.value = undefined; await load(); }
}, { immediate: true });

function choose(row: DiyPageLink) { selected.value = { name: row.name, url: row.url }; }
function chooseMicro(row: DiyPageRow) { selected.value = { name: row.name, url: `/pages/diy/page?id=${row.id}` }; }
function confirm() {
  if (!selected.value) { ElMessage.warning('请选择页面链接'); return; }
  emit('closeDialog', { name: selected.value.name, type: source.value === 'micro' ? '微页面' : '页面链接', url: selected.value.url });
}
function close() { emit('closeDialog'); }
</script>

<template>
  <ElDialog v-model="open" title="选择平台页面链接" width="720px" destroy-on-close @closed="close">
    <ElRadioGroup v-model="source" class="mb-3"><ElRadioButton value="link">系统链接</ElRadioButton><ElRadioButton value="micro">微页面</ElRadioButton></ElRadioGroup>
    <div v-if="source === 'link'" class="mb-3 flex items-center gap-3"><span class="shrink-0 text-sm">链接分类</span><ElSelect v-model="cateID" clearable placeholder="全部分类" class="w-full"><ElOption v-for="item in categoryOptions" :key="item.id" :value="item.id" :label="`${'— '.repeat(Math.max(0, item.level - 1))}${item.name}`" /></ElSelect></div>
    <ElTable v-if="source === 'link'" v-loading="loading" :data="filteredLinks" highlight-current-row @current-change="choose">
      <ElTableColumn prop="name" label="页面名称" min-width="180" />
      <ElTableColumn prop="url" label="页面路径" min-width="300" show-overflow-tooltip />
      <ElTableColumn label="分类" width="130"><template #default="{ row }">{{ row.category?.name || '-' }}</template></ElTableColumn>
    </ElTable>
    <ElTable v-else v-loading="loading" :data="microPages" highlight-current-row @current-change="chooseMicro"><ElTableColumn prop="name" label="微页面名称" min-width="220" /><ElTableColumn prop="title" label="页面标题" min-width="180" /><ElTableColumn prop="update_time" label="更新时间" min-width="180" /></ElTable>
    <template #footer><ElButton @click="close">取消</ElButton><ElButton type="primary" @click="confirm">确定</ElButton></template>
  </ElDialog>
</template>
