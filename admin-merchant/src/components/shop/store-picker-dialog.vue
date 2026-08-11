<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { StoreListItem } from '#/api/core/store';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, reactive, ref, watch } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getStoreListApi } from '#/api/core/store';

type StoreRow = StoreListItem & { noChoose?: boolean };

const open = defineModel<boolean>('open', { default: false });

const props = withDefaults(
  defineProps<{
    excludeIds?: number[];
    multiple?: boolean;
  }>(),
  { excludeIds: () => [], multiple: false },
);

const emit = defineEmits<{
  closeDialog: [payload: { params?: unknown; type: string }];
}>();

const selectedRows = ref<StoreRow[]>([]);

const formOptions: VbenFormProps = {
  actionLayout: 'inline',
  collapsed: false,
  schema: [
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '门店名称' },
      fieldName: 'store_name',
      label: '门店搜索',
    },
  ],
  showCollapseButton: false,
  resetButtonOptions: { content: '重置' },
  submitButtonOptions: { content: '搜索' },
  submitOnChange: false,
  submitOnEnter: true,
  wrapperClass: 'grid-cols-1 md:grid-cols-2',
};

const canConfirm = computed(() => selectedRows.value.length > 0);

function isSelectable(row: StoreRow) {
  return row.noChoose !== false;
}

const gridOptions = reactive<VxeGridProps<StoreRow>>({
  checkboxConfig: {
    checkMethod: ({ row }) => isSelectable(row),
    highlight: true,
    reserve: true,
  },
  columns: [
    { type: 'checkbox', width: 48 },
    {
      field: 'logo',
      slots: { default: 'logo' },
      title: '门店图片',
      width: 80,
    },
    { field: 'store_name', minWidth: 140, showOverflow: true, title: '门店名称' },
    { field: 'linkman', title: '联系人', width: 100 },
    { field: 'phone', title: '联系电话', width: 120 },
    { field: 'create_time', title: '添加时间', width: 150 },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const res = await getStoreListApi({
          list_rows: page.pageSize,
          page: page.currentPage,
          store_name: String(formValues?.store_name ?? '').trim() || undefined,
        });
        const rows = (res.list.data ?? []).map((row) => ({
          ...row,
          noChoose: props.excludeIds.length
            ? !props.excludeIds.includes(row.store_id)
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
    keyField: 'store_id',
  },
  toolbarConfig: {
    enabled: false,
  },
});

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

function onCheckboxChange() {
  selectedRows.value = (gridApi.grid?.getCheckboxRecords?.() ?? []) as StoreRow[];
}

function confirm() {
  if (!selectedRows.value.length) {
    ElMessage.error('请至少选择一个门店');
    return;
  }
  emit('closeDialog', {
    params: props.multiple ? selectedRows.value : selectedRows.value[0],
    type: 'success',
  });
  open.value = false;
}

function cancel() {
  emit('closeDialog', { type: 'error' });
  open.value = false;
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    selectedRows.value = [];
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
    class="h-[min(76dvh,820px)] w-[min(94vw,1200px)] max-w-[94vw]"
    title="请选择店铺："
  >
    <Grid @checkbox-change="onCheckboxChange">
      <template #logo="{ row }">
        <img
          v-if="row.logo?.file_path"
          :alt="row.store_name"
          class="h-8 w-8 rounded object-cover"
          :src="row.logo.file_path"
        />
      </template>
    </Grid>

    <template #footer>
      <ElButton @click="cancel">取消</ElButton>
      <ElButton :disabled="!canConfirm" type="primary" @click="confirm">确定</ElButton>
    </template>
  </Modal>
</template>
