<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { ProductSpecChooseItem } from '#/api/core/product';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, reactive, ref, watch } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { chooseProductSpecApi } from '#/api/core/product';

type SpecRow = ProductSpecChooseItem & { noChoose?: boolean };

const open = defineModel<boolean>('open', { default: false });

const props = withDefaults(
  defineProps<{
    excludeIds?: number[];
    multiple?: boolean;
    productId: number;
  }>(),
  { excludeIds: () => [], multiple: true },
);

const emit = defineEmits<{
  close: [payload: { open: false; params: ProductSpecChooseItem[] | null; type: string }];
}>();

const tableData = ref<SpecRow[]>([]);
const selectedRows = ref<SpecRow[]>([]);

const canConfirm = computed(() => selectedRows.value.length > 0);

const gridOptions = reactive<VxeGridProps<SpecRow>>({
  checkboxConfig: props.multiple
    ? {
        checkMethod: ({ row }) => row.noChoose !== false,
        highlight: true,
        reserve: true,
      }
    : undefined,
  columns: [
    { field: 'spec_name', minWidth: 200, showOverflow: true, title: '商品/规格' },
    ...(props.multiple
      ? [{ type: 'checkbox' as const, width: 48 }]
      : [
          {
            field: 'action',
            fixed: 'right' as const,
            slots: { default: 'action' },
            title: '单选',
            width: 80,
          },
        ]),
  ],
  data: tableData,
  minHeight: 280,
  radioConfig: props.multiple
    ? undefined
    : {
        checkMethod: ({ row }) => row.noChoose !== false,
        highlight: true,
        strict: true,
        trigger: 'row',
      },
  rowConfig: {
    keyField: 'product_sku_id',
  },
  toolbarConfig: {
    enabled: false,
  },
});

const [Grid, gridApi] = useVbenVxeGrid({
  gridEvents: {
    checkboxChange: onCheckboxChange,
    radioChange: onRadioChange,
  },
  gridOptions,
});

async function fetchSpecs() {
  if (!props.productId) {
    tableData.value = [];
    return;
  }
  gridApi.setLoading(true);
  try {
    const res = await chooseProductSpecApi(props.productId);
    tableData.value = (res.specList ?? []).map((item) => ({
      ...item,
      noChoose:
        props.excludeIds.length > 0
          ? !props.excludeIds.includes(item.product_sku_id)
          : true,
    }));
  } finally {
    gridApi.setLoading(false);
  }
}

function pickSingle(row: SpecRow) {
  if (row.noChoose === false) return;
  selectedRows.value = [row];
  confirm();
}

function onCheckboxChange() {
  selectedRows.value = (gridApi.grid?.getCheckboxRecords?.() ?? []) as SpecRow[];
}

function onRadioChange({ row }: { row: SpecRow }) {
  selectedRows.value = row ? [row] : [];
}

function confirm() {
  if (!selectedRows.value.length) {
    ElMessage.error('请至少选择一件商品规格');
    return;
  }
  emit('close', {
    open: false,
    params: props.multiple ? [...selectedRows.value] : [selectedRows.value[0]!],
    type: 'success',
  });
  open.value = false;
}

function cancel() {
  emit('close', { open: false, params: null, type: 'error' });
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
    void fetchSpecs();
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
    class="w-[600px]"
    title="选择规格"
  >
    <Grid>
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

    <template #footer>
      <ElButton @click="cancel">取消</ElButton>
      <ElButton v-if="multiple" :disabled="!canConfirm" type="primary" @click="confirm">
        确定
      </ElButton>
    </template>
  </Modal>
</template>
