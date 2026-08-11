<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { ProductCategoryOption, ProductChooseItem } from '#/api/core/product';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, nextTick, reactive, ref, shallowRef, watch } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getProductChooseListApi, resolveProductTypeText } from '#/api/core/product';

type ProductRow = ProductChooseItem & { noChoose?: boolean };

const open = defineModel<boolean>('open', { default: false });

const props = withDefaults(
  defineProps<{
    excludeIds?: number[];
  }>(),
  { excludeIds: () => [] },
);

const emit = defineEmits<{
  closeDialog: [payload: { openDialog: boolean; params?: unknown; type: string }];
}>();

const categoryOptions = shallowRef<Array<{ label: string; value: number }>>([
  { label: '全部', value: 0 },
]);

const selectedRows = ref<ProductRow[]>([]);

const canConfirm = computed(() => selectedRows.value.length > 0);

function isExcluded(productId: number | string) {
  if (!props.excludeIds?.length) return false;
  const id = Number(productId);
  return props.excludeIds.some((item) => Number(item) === id);
}

function dedupeRowsByProductId(rows: ProductRow[]) {
  const seen = new Set<number>();
  const out: ProductRow[] = [];
  for (const row of rows) {
    const id = Number(row.product_id);
    if (!Number.isFinite(id) || seen.has(id)) continue;
    seen.add(id);
    out.push(row);
  }
  return out;
}

function collectCheckedRows(): ProductRow[] {
  const grid = gridApi.grid;
  if (!grid) return [];

  const current = (grid.getCheckboxRecords?.() ?? []) as ProductRow[];
  const reserved = (grid.getCheckboxReserveRecords?.() ?? []) as ProductRow[];
  const full = (grid.getCheckboxRecords?.(true) ?? []) as ProductRow[];

  return dedupeRowsByProductId([...current, ...reserved, ...full]);
}

function flattenCategories(
  items: ProductCategoryOption[],
  depth = 0,
): Array<{ label: string; value: number }> {
  const out: Array<{ label: string; value: number }> = [];
  for (const item of items) {
    out.push({
      label: `${depth ? '　'.repeat(depth) : ''}${item.name}`,
      value: Number(item.category_id),
    });
    if (item.child?.length) {
      out.push(...flattenCategories(item.child, depth + 1));
    }
  }
  return out;
}

const formOptions = computed((): VbenFormProps => ({
  actionLayout: 'inline',
  collapsed: false,
  schema: [
    {
      component: 'Select',
      componentProps: {
        options: categoryOptions.value,
        placeholder: '商品分类',
      },
      defaultValue: 0,
      fieldName: 'category_id',
      label: '商品分类',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '商品名称',
      },
      fieldName: 'search',
      label: '商品名称',
    },
  ],
  showCollapseButton: false,
  resetButtonOptions: {
    content: '重置',
  },
  submitButtonOptions: {
    content: '搜索',
  },
  submitOnChange: false,
  submitOnEnter: true,
  wrapperClass: 'grid-cols-1 md:grid-cols-3',
}));

const gridOptions = reactive<VxeGridProps<ProductRow>>({
  checkboxConfig: {
    checkMethod: ({ row }) => !isExcluded(row.product_id),
    highlight: true,
    reserve: true,
  },
  columns: [
    {
      field: 'image',
      slots: { default: 'image' },
      title: '图片',
      width: 72,
    },
    { field: 'product_name', minWidth: 160, showOverflow: true, title: '名称' },
    {
      field: 'product_type_text',
      showOverflow: true,
      slots: { default: 'productType' },
      title: '类型',
      width: 88,
    },
    { field: 'category_name', showOverflow: true, title: '分类', width: 120 },
    {
      field: 'product_price',
      slots: { default: 'price' },
      title: '价格',
      width: 100,
    },
    { field: 'create_time', title: '创建时间', width: 150 },
    { fixed: 'right', type: 'checkbox', width: 48 },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const res = await getProductChooseListApi({
          category_id: Number(formValues?.category_id) || undefined,
          list_rows: page.pageSize,
          page: page.currentPage,
          search: String(formValues?.search ?? '').trim(),
        });
        if (res.category?.length) {
          categoryOptions.value = [
            { label: '全部', value: 0 },
            ...flattenCategories(res.category),
          ];
        }
        const rows = (res.list.data ?? []).map((row) => ({
          ...row,
          noChoose: !isExcluded(row.product_id),
        }));
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

function syncSelectionFromGrid() {
  selectedRows.value = collectCheckedRows();
}

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridEvents: {
    checkboxAll: syncSelectionFromGrid,
    checkboxChange: syncSelectionFromGrid,
  },
  gridOptions,
});

async function resetPicker() {
  selectedRows.value = [];
  categoryOptions.value = [{ label: '全部', value: 0 }];
  await nextTick();
  if (typeof gridApi.grid?.commitProxy === 'function') {
    void gridApi.reload();
    return;
  }
  void gridApi.query?.();
}

async function confirm() {
  await nextTick();
  syncSelectionFromGrid();
  if (!selectedRows.value.length) {
    ElMessage.error('请至少选择一件商品');
    return;
  }
  const rows = selectedRows.value.map((row) => ({
    ...row,
    product_image: row.image?.[0]?.file_path ?? row.product_image,
    image: row.image?.[0]?.file_path,
  }));
  emit('closeDialog', {
    openDialog: false,
    params: rows,
    type: 'success',
  });
  open.value = false;
}

function cancel() {
  emit('closeDialog', { openDialog: false, params: null, type: 'error' });
  open.value = false;
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
    if (isOpen) {
      void resetPicker();
    }
  },
});

watch(open, (visible) => {
  if (visible) {
    modalApi.open();
    return;
  }
  modalApi.close();
});
</script>

<template>
  <Modal
    :close-on-click-modal="false"
    :destroy-on-close="true"
    class="h-[min(76dvh,820px)] w-[min(94vw,1200px)] max-w-[94vw]"
    title="请选择商品："
  >
    <Grid>
      <template #image="{ row }">
        <img
          v-if="row.image?.[0]?.file_path"
          :alt="row.product_name"
          class="h-8 w-8 object-cover"
          :src="row.image[0].file_path"
        />
      </template>
      <template #price="{ row }">
        {{ row.product_price ?? '-' }}
      </template>
      <template #productType="{ row }">
        {{ resolveProductTypeText(row) }}
      </template>
    </Grid>

    <template #footer>
      <ElButton @click="cancel">取消</ElButton>
      <ElButton :disabled="!canConfirm" type="primary" @click="confirm">确定</ElButton>
    </template>
  </Modal>
</template>
