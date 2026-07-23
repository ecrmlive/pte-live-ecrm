<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { ProductListItem } from '#/api/core/product';

import { useVbenModal } from '@vben/common-ui';
import { ElButton } from 'element-plus';
import { computed, reactive, ref, watch } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getProductListApi } from '#/api/core/product';

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{
  select: [ProductListItem];
}>();

const selectedRow = ref<ProductListItem | null>(null);

const formOptions: VbenFormProps = {
  actionLayout: 'inline',
  collapsed: false,
  schema: [
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '商品名称',
      },
      fieldName: 'product_name',
      label: '商品名称',
    },
  ],
  showCollapseButton: false,
  submitButtonOptions: {
    content: '查询',
  },
  submitOnChange: false,
  submitOnEnter: true,
  wrapperClass: 'grid-cols-1 md:grid-cols-2',
};

const gridOptions = reactive<VxeGridProps<ProductListItem>>({
  columns: [
    { type: 'radio', width: 48 },
    { field: 'product_id', title: 'ID', width: 80 },
    { field: 'product_name', minWidth: 220, showOverflow: true, title: '商品' },
    { field: 'product_price', title: '售价', width: 100 },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 15,
    pageSizes: [15, 30, 50],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const res = await getProductListApi({
          list_rows: page.pageSize,
          page: page.currentPage,
          product_name: String(formValues?.product_name ?? '').trim(),
        });
        return {
          items: res.list.data,
          total: res.list.total,
        };
      },
    },
  },
  radioConfig: {
    highlight: true,
    strict: true,
    trigger: 'row',
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
  gridEvents: {
    radioChange: onRadioChange,
  },
  gridOptions,
});

const canConfirm = computed(() => Boolean(selectedRow.value));

function onRadioChange({ row }: { row: ProductListItem }) {
  selectedRow.value = row;
}

function confirm() {
  if (!selectedRow.value) return;
  emit('select', selectedRow.value);
  open.value = false;
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    selectedRow.value = null;
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
    class="w-[860px]"
    title="选择商品"
  >
    <Grid />
    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :disabled="!canConfirm" type="primary" @click="confirm">确定</ElButton>
    </template>
  </Modal>
</template>
