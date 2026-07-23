<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { ProductCategoryOption, ProductIndexChooseItem } from '#/api/core/product';

import { useVbenModal } from '@vben/common-ui';
import { ElButton } from 'element-plus';
import { computed, reactive, ref, shallowRef, watch } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getProductIndexListApi } from '#/api/core/product';

const open = defineModel<boolean>('open', { default: false });

const props = withDefaults(
  defineProps<{
    excludeIds?: number[];
    isForm?: number;
    productType?: number;
  }>(),
  { excludeIds: () => [], isForm: 0, productType: 1 },
);

const emit = defineEmits<{
  closeDialog: [payload: { openDialog: boolean; params?: ProductIndexChooseItem; type: string }];
}>();

const categoryOptions = shallowRef<Array<{ label: string; value: number }>>([
  { label: '全部', value: 0 },
]);

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

function isExcluded(productId: number) {
  return props.excludeIds.includes(Number(productId));
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
  submitButtonOptions: {
    content: '查询',
  },
  submitOnChange: false,
  submitOnEnter: true,
  wrapperClass: 'grid-cols-1 md:grid-cols-3',
}));

const gridOptions = reactive<VxeGridProps<ProductIndexChooseItem>>({
  columns: [
    {
      field: 'image',
      slots: { default: 'image' },
      title: '商品图片',
      width: 80,
    },
    { field: 'product_name', minWidth: 180, showOverflow: true, title: '商品名称' },
    {
      field: 'category_name',
      slots: { default: 'category' },
      title: '商品分类',
      width: 120,
    },
    { field: 'create_time', title: '添加时间', width: 150 },
    {
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      title: '单选',
      width: 80,
    },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const res = await getProductIndexListApi({
          category_id: Number(formValues?.category_id) || undefined,
          is_form: props.isForm,
          list_rows: page.pageSize,
          page: page.currentPage,
          product_type: props.productType,
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
          noChoose: props.excludeIds.length
            ? !isExcluded(row.product_id)
            : true,
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

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridOptions,
});

function pickSingle(row: ProductIndexChooseItem) {
  if (row.noChoose === false) return;
  emit('closeDialog', {
    openDialog: false,
    params: row,
    type: 'success',
  });
  open.value = false;
}

function cancel() {
  emit('closeDialog', { openDialog: false, params: undefined, type: 'error' });
  open.value = false;
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    categoryOptions.value = [{ label: '全部', value: 0 }];
    void gridApi.reload();
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
    class="w-[900px]"
    title="选择商品"
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
      <template #category="{ row }">
        {{ row.category?.name ?? row.category_name }}
      </template>
      <template #action="{ row }">
        <ElButton
          :disabled="row.noChoose === false"
          size="small"
          @click="pickSingle(row)"
        >
          {{ row.noChoose === false ? '已选' : '选择' }}
        </ElButton>
      </template>
    </Grid>
  </Modal>
</template>
