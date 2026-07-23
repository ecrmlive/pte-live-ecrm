<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { ProductCategoryOption } from '#/api/core/product';
import type { ShopLinkPickerItem } from '#/api/core/shop-link';
import type { VbenFormSchema } from '#/adapter/form';

import { Search } from '@element-plus/icons-vue';
import type { CascaderInstance } from 'element-plus';
import { ElButton } from 'element-plus';
import { computed, onMounted, reactive, ref, watch } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  getShopLinkProductCategoryApi,
  getShopLinkProductListsApi,
} from '#/api/core/shop-link';

const emit = defineEmits<{
  change: [ShopLinkPickerItem];
}>();

const activeTab = ref<'detail' | 'type'>('type');
const loading = ref(false);
const categoryList = ref<ProductCategoryOption[]>([]);
const categoryActive = ref<Array<number | string>>([]);
const cascaderRef = ref<CascaderInstance>();

type ProductRow = {
  image: Array<{ file_path: string }>;
  product_id: number;
  product_name: string;
  product_price: number | string;
};

const searchSchema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: { placeholder: '请输入商品名称' },
    fieldName: 'product_name',
    label: '商品名称',
  },
]);

const [SearchForm, searchFormApi] = useVbenForm(
  reactive({
    actionLayout: 'inline',
    commonConfig: {
      componentProps: { size: 'small' },
    },
    handleSubmit: async () => {
      onSearch();
    },
    layout: 'horizontal',
    resetButtonOptions: { show: false },
    schema: searchSchema,
    showDefaultActions: true,
    submitButtonOptions: {
      content: '查询',
      icon: Search,
    },
    wrapperClass: 'grid-cols-1 md:grid-cols-2',
  }),
);

const gridOptions = reactive<VxeGridProps<ProductRow>>({
  columns: [
    {
      field: 'product_name',
      minWidth: 180,
      slots: { default: 'product' },
      title: '商品',
    },
    {
      field: 'product_price',
      slots: { default: 'price' },
      title: '价格',
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
        const values = await searchFormApi.getValues();
        const res = await getShopLinkProductListsApi({
          list_rows: page.pageSize,
          page: page.currentPage,
          product_name: String(values.product_name ?? ''),
        });
        const rows = res.list.data ?? [];
        if (page.currentPage === 1 && rows[0]) {
          emitProduct(rows[0]);
        }
        return {
          items: rows,
          total: res.list.total ?? 0,
        };
      },
    },
  },
  rowConfig: {
    keyField: 'product_id',
  },
  toolbarConfig: {
    enabled: false,
  },
});

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

function emitCategory(item: ProductCategoryOption) {
  emit('change', {
    name: item.name,
    type: '商品分类',
    url: `pages/mall/product/list/list?category_id=${item.category_id}`,
  });
}

function emitProduct(row: ProductRow) {
  emit('change', {
    name: row.product_name,
    type: '商品详情',
    url: `pages/mall/product/detail/detail?product_id=${row.product_id}`,
  });
}

function autoType(index = 0) {
  categoryActive.value = [];
  const item = categoryList.value[index];
  if (!item) return;
  categoryActive.value.push(item.category_id);
  if (item.child?.length) {
    categoryActive.value.push(item.child[0]!.category_id);
    emitCategory(item.child[0]!);
    return;
  }
  autoType(index + 1);
}

function changeCategory() {
  const nodes = cascaderRef.value?.getCheckedNodes(false);
  const data = nodes?.[0]?.data as ProductCategoryOption | undefined;
  if (data) emitCategory(data);
}

async function loadCategory() {
  loading.value = true;
  try {
    const res = await getShopLinkProductCategoryApi();
    categoryList.value = res.list ?? [];
    autoType();
  } finally {
    loading.value = false;
  }
}

async function loadProducts() {
  await gridApi.reload();
}

function onSearch() {
  void loadProducts();
}

watch(activeTab, (tab) => {
  if (tab === 'type') {
    if (!categoryList.value.length) void loadCategory();
    else autoType();
  } else {
    void loadProducts();
  }
});

onMounted(() => {
  void loadCategory();
});
</script>

<template>
  <div class="marketing-box">
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
        :props="{
          checkStrictly: true,
          children: 'child',
          label: 'name',
          value: 'category_id',
        }"
        @change="changeCategory"
      />
    </div>

    <div v-else v-loading="loading">
      <SearchForm class="mb-3" />

      <Grid>
        <template #product="{ row }">
          <div class="flex items-center gap-2">
            <img
              v-if="row.image?.[0]?.file_path"
              v-img-url="row.image[0].file_path"
              alt=""
              class="h-5 w-5 object-cover"
            />
            <span class="truncate">{{ row.product_name }}</span>
          </div>
        </template>
        <template #price="{ row }">
          <span class="text-red-500">{{ row.product_price }}</span>
        </template>
        <template #action="{ row }">
          <ElButton size="small" type="primary" @click="emitProduct(row)">选择</ElButton>
        </template>
      </Grid>
    </div>
  </div>
</template>
