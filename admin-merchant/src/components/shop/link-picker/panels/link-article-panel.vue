<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';
import type {
  ShopLinkArticleCategory,
  ShopLinkArticleRow,
  ShopLinkPickerItem,
} from '#/api/core/shop-link';

import type { CascaderInstance } from 'element-plus';
import { ElButton } from 'element-plus';
import { onMounted, reactive, ref, watch } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  getShopLinkArticleCategoryApi,
  getShopLinkArticleListApi,
} from '#/api/core/shop-link';

const emit = defineEmits<{
  change: [ShopLinkPickerItem];
}>();

const activeTab = ref<'detail' | 'type'>('type');
const loading = ref(false);
const categoryList = ref<ShopLinkArticleCategory[]>([]);
const categoryActive = ref<Array<number | string>>([]);
const cascaderRef = ref<CascaderInstance>();

const gridOptions = reactive<VxeGridProps<ShopLinkArticleRow>>({
  columns: [
    { field: 'article_title', minWidth: 180, showOverflow: true, title: '文章标题' },
    {
      field: 'category.name',
      slots: { default: 'category' },
      title: '文章分类',
      width: 100,
    },
    {
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 80,
    },
  ],
  minHeight: 280,
  pagerConfig: {
    pageSize: 5,
    pageSizes: [5, 10, 20],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const res = await getShopLinkArticleListApi({
          list_rows: page.pageSize,
          page: page.currentPage,
        });
        return {
          items: res.list.data ?? [],
          total: res.list.total ?? 0,
        };
      },
    },
  },
  rowConfig: {
    keyField: 'article_id',
  },
  toolbarConfig: {
    enabled: false,
  },
});

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

function emitCategory(item: ShopLinkArticleCategory) {
  emit('change', {
    name: item.name,
    type: '文章分类',
    url: `pages/content/article/list/list?category_id=${item.category_id}`,
  });
}

function emitArticle(row: ShopLinkArticleRow) {
  emit('change', {
    name: row.article_title,
    type: '文章详情',
    url: `pages/content/article/detail/detail?article_id=${row.article_id}`,
  });
}

function autoType() {
  categoryActive.value = [];
  const first = categoryList.value[0];
  if (!first) return;
  categoryActive.value = [first.category_id];
  emitCategory(first);
}

function changeCategory() {
  const nodes = cascaderRef.value?.getCheckedNodes(false);
  const data = nodes?.[0]?.data as ShopLinkArticleCategory | undefined;
  if (data) emitCategory(data);
}

async function loadCategory() {
  loading.value = true;
  try {
    const res = await getShopLinkArticleCategoryApi();
    categoryList.value = res.list ?? [];
    autoType();
  } finally {
    loading.value = false;
  }
}

async function loadArticles() {
  await gridApi.reload();
}

watch(activeTab, (tab) => {
  if (tab === 'type') {
    if (!categoryList.value.length) void loadCategory();
    else autoType();
  } else {
    void loadArticles();
  }
});

onMounted(() => {
  void loadCategory();
});
</script>

<template>
  <div class="article-box">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="分类" name="type" />
      <el-tab-pane label="详情" name="detail" />
    </el-tabs>

    <div v-if="activeTab === 'type'" v-loading="loading">
      <el-cascader
        ref="cascaderRef"
        v-model="categoryActive"
        class="w-full"
        :options="categoryList"
        :props="{ children: 'child', label: 'name', value: 'category_id' }"
        @change="changeCategory"
      />
    </div>

    <div v-else v-loading="loading">
      <Grid>
        <template #category="{ row }">
          {{ row.category?.name }}
        </template>
        <template #action="{ row }">
          <ElButton size="small" type="primary" @click="emitArticle(row)">选择</ElButton>
        </template>
      </Grid>
    </div>
  </div>
</template>
