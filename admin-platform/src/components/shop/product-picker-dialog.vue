<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { PlatformProduct } from '#/api/core/platform-catalog';

import { computed, reactive, ref, watch } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElImage, ElMessage } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { listPlatformProductsApi } from '#/api/core/platform-catalog';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

const props = withDefaults(
  defineProps<{
    /** 多选模式（氛围图「指定商品」等） */
    multiple?: boolean;
  }>(),
  { multiple: false },
);

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{
  confirm: [PlatformProduct[]];
  select: [PlatformProduct];
}>();

const selectedRow = ref<PlatformProduct | null>(null);
const selectedRows = ref<PlatformProduct[]>([]);

const formOptions: VbenFormProps = {
  actionLayout: 'inline',
  collapsed: false,
  schema: [
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '商品名称 / ID',
      },
      fieldName: 'keyword',
      label: '关键字',
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
  wrapperClass: 'grid-cols-1 md:grid-cols-2',
};

const gridOptions = reactive<VxeGridProps<PlatformProduct>>({
  checkboxConfig: {
    highlight: true,
    reserve: true,
    trigger: 'row',
  },
  columns: [
    { type: 'radio', width: 48, visible: !props.multiple },
    { type: 'checkbox', width: 48, visible: !!props.multiple },
    { field: 'product_id', title: 'ID', width: 80 },
    {
      field: 'image',
      slots: { default: 'cover' },
      title: '封面',
      width: 72,
    },
    {
      field: 'store_name',
      minWidth: 220,
      showOverflow: true,
      title: '商品',
    },
    {
      field: 'mer_name',
      minWidth: 120,
      showOverflow: true,
      title: '店铺',
    },
    { field: 'price', title: '售价', width: 100 },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const res = await listPlatformProductsApi({
          limit: page.pageSize,
          page: page.currentPage,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
        });
        return {
          items: res.list || [],
          total: res.total || 0,
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
    checkboxAll({ records }: { records: PlatformProduct[] }) {
      if (props.multiple) selectedRows.value = records;
    },
    checkboxChange({ records }: { records: PlatformProduct[] }) {
      if (props.multiple) selectedRows.value = records;
    },
    radioChange: onRadioChange,
  },
  gridOptions,
});

const canConfirm = computed(() =>
  props.multiple
    ? selectedRows.value.length > 0
    : Boolean(selectedRow.value),
);

function coverOf(row: PlatformProduct) {
  return resolveCosMediaUrl(String(row.image || '').trim());
}

function onRadioChange({ row }: { row: PlatformProduct }) {
  selectedRow.value = row;
}

function confirm() {
  if (props.multiple) {
    if (!selectedRows.value.length) {
      ElMessage.warning('请选择商品');
      return;
    }
    emit('confirm', [...selectedRows.value]);
    open.value = false;
    return;
  }
  if (!selectedRow.value) {
    ElMessage.warning('请选择商品');
    return;
  }
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
    title="请选择商品："
  >
    <Grid>
      <template #cover="{ row }">
        <ElImage
          v-if="coverOf(row)"
          :src="coverOf(row)"
          fit="cover"
          class="product-picker-cover"
        >
          <template #error>
            <span class="text-xs text-gray-400">—</span>
          </template>
        </ElImage>
        <span v-else class="text-xs text-gray-400">—</span>
      </template>
    </Grid>
    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :disabled="!canConfirm" type="primary" @click="confirm">
        确定
      </ElButton>
    </template>
  </Modal>
</template>

<style scoped>
.product-picker-cover {
  display: block;
  width: 40px;
  height: 40px;
  border-radius: 4px;
}
</style>
