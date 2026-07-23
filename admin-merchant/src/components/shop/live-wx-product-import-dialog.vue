<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { LiveWxProductItem } from '#/api/core/plus-live-wx-product';

import { useVbenModal } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';
import { computed, reactive, ref, watch } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getLiveWxProductPickListApi } from '#/api/core/plus-live-wx-product';

defineOptions({ name: 'LiveWxProductImportDialog' });

type WxProductRow = LiveWxProductItem & { noChoose?: boolean };

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  roomId: number | string;
}>();

const emit = defineEmits<{
  confirm: [goodsIds: Array<number | string>];
}>();

const selectedRows = ref<WxProductRow[]>([]);
const canConfirm = computed(() => selectedRows.value.length > 0);

const gridOptions = reactive<VxeGridProps<WxProductRow>>({
  checkboxConfig: {
    checkMethod: ({ row }) => row.noChoose !== false,
    highlight: true,
    reserve: true,
  },
  columns: [
    { type: 'checkbox', width: 44 },
    {
      field: 'cover_img',
      slots: { default: 'image' },
      title: '商品图片',
      width: 70,
    },
    { field: 'name', minWidth: 160, showOverflow: true, title: '商品名称' },
    { field: 'price_type_text', title: '价格类型', width: 100 },
    { field: 'price_text', title: '价格', width: 160 },
    { field: 'create_time', title: '添加时间', width: 140 },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 20,
    pageSizes: [10, 20, 50],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const res = await getLiveWxProductPickListApi({
          list_rows: page.pageSize,
          page: page.currentPage,
          room_id: props.roomId,
        });
        const excludeIds = res.excludeIds ?? [];
        const rows = (res.list.data ?? []).map((item) => ({
          ...item,
          noChoose:
            excludeIds.length === 0 ||
            !excludeIds.includes(item.goods_id as number | string),
        }));
        return {
          items: rows,
          total: res.list.total ?? 0,
        };
      },
    },
  },
  rowConfig: {
    keyField: 'goods_id',
  },
  toolbarConfig: {
    enabled: false,
  },
});

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

function onCheckboxChange() {
  selectedRows.value = (gridApi.grid?.getCheckboxRecords?.() ?? []) as WxProductRow[];
}

function confirm() {
  if (!selectedRows.value.length) {
    ElMessage.error('请至少选择一件商品');
    return;
  }
  emit(
    'confirm',
    selectedRows.value.map((item) => item.goods_id as number | string),
  );
  open.value = false;
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    modalApi.open();
    selectedRows.value = [];
    void gridApi.reload();
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
    <Grid @checkbox-change="onCheckboxChange">
      <template #image="{ row }">
        <img :src="row.cover_img" alt="" class="h-[30px] w-[30px] object-cover" />
      </template>
    </Grid>

    <template #footer>
      <ElButton @click="open = false">取消</ElButton>
      <ElButton :disabled="!canConfirm" type="primary" @click="confirm">确定</ElButton>
    </template>
  </Modal>
</template>
