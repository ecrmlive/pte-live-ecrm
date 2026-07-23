<script setup lang="ts">
import type { ProductCategoryItem } from '#/api/core/product';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElEmpty } from 'element-plus';
import { ref, watch } from 'vue';

import { getProductCategoryListApi } from '#/api/core/product';

export interface CouponCategoryList {
  first: Array<{ category_id: number; name: string; parent?: string; parent_id?: number }>;
  second: Array<{ category_id: number; name: string; parent?: string; parent_id?: number }>;
}

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  categoryList?: CouponCategoryList;
}>();

const emit = defineEmits<{
  confirm: [CouponCategoryList];
}>();

const loading = ref(false);
const options = ref<ProductCategoryItem[]>([]);
const listData = ref<Array<number | number[]>>([]);
const panelKey = ref(0);

function normalizeCategoryTree(items: ProductCategoryItem[] = []): ProductCategoryItem[] {
  return items.map((item) => {
    const children = item.child?.length ? normalizeCategoryTree(item.child) : undefined;
    return {
      ...item,
      category_id: Number(item.category_id),
      child: children,
    };
  });
}

async function fetchCategories() {
  loading.value = true;
  try {
    const res = await getProductCategoryListApi();
    const rawList = Array.isArray(res.list) ? res.list : [];
    options.value = normalizeCategoryTree(rawList);
    panelKey.value += 1;
  } finally {
    loading.value = false;
  }
}

function syncFromProps() {
  listData.value = [];
  const list = props.categoryList;
  if (!list) return;
  list.first?.forEach((item) => {
    listData.value.push([item.category_id]);
  });
  list.second?.forEach((item) => {
    if (item.parent_id) {
      listData.value.push([item.parent_id, item.category_id]);
    }
  });
}

function confirm() {
  const result: CouponCategoryList = { first: [], second: [] };
  for (const path of listData.value) {
    const ids = Array.isArray(path) ? path : [path];
    if (ids.length === 1) {
      const id = ids[0];
      const node = findNode(options.value, id);
      if (node) {
        result.first.push({ category_id: node.category_id, name: node.name });
      }
    } else if (ids.length >= 2) {
      const parentId = ids[ids.length - 2];
      const categoryId = ids[ids.length - 1];
      const parent = findNode(options.value, parentId);
      const node = parent?.child?.find((row) => row.category_id === categoryId);
      if (node && parent) {
        result.second.push({
          category_id: node.category_id,
          name: node.name,
          parent: parent.name,
          parent_id: parent.category_id,
        });
      }
    }
  }
  emit('confirm', result);
  open.value = false;
}

function findNode(
  tree: ProductCategoryItem[],
  id: number | undefined,
): ProductCategoryItem | undefined {
  if (id === undefined) return undefined;
  for (const item of tree) {
    if (item.category_id === id) return item;
    const child = findNode(item.child ?? [], id);
    if (child) return child;
  }
  return undefined;
}

const [Modal, modalApi] = useVbenModal({
  appendToMain: true,
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    syncFromProps();
    void fetchCategories();
    modalApi.open();
    return;
  }
  modalApi.close();
});
</script>

<template>
  <Modal
    :append-to-main="true"
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="w-[620px]"
    title="选择分类"
  >
    <div v-loading="loading" class="category-picker-body">
      <el-cascader-panel
        v-if="options.length"
        :key="panelKey"
        v-model="listData"
        class="category-picker-panel"
        :options="options"
        :props="{ multiple: true, value: 'category_id', label: 'name', children: 'child' }"
        :show-all-levels="false"
      />
      <ElEmpty v-else-if="!loading" description="暂无商品分类，请先在商品分类中添加" />
    </div>
    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton type="primary" @click="confirm">确定</ElButton>
    </template>
  </Modal>
</template>

<style scoped lang="scss">
.category-picker-body {
  min-height: 320px;
}

.category-picker-panel {
  width: 100%;

  :deep(.el-cascader-panel) {
    width: 100%;
    min-height: 300px;
    border: 1px solid var(--el-border-color-light);
    border-radius: 8px;
  }

  :deep(.el-cascader-menu) {
    min-width: 180px;
    min-height: 300px;
  }

  :deep(.el-cascader-menu__wrap) {
    height: 300px;
  }
}
</style>
