<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { CouponListItem } from '#/api/core/plus-coupon';

import { useVbenModal } from '@vben/common-ui';
import { ElButton } from 'element-plus';
import { reactive, watch } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getCouponListApi } from '#/api/core/plus-coupon';

const open = defineModel<boolean>('open', { default: false });

const emit = defineEmits<{
  select: [CouponListItem];
}>();

const gridOptions = reactive<VxeGridProps<CouponListItem>>({
  columns: [
    { field: 'name', minWidth: 160, showOverflow: true, title: '名称' },
    { field: 'min_price', title: '最低消费', width: 120 },
    {
      field: 'total_num',
      slots: { default: 'total_num' },
      title: '数量',
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
  minHeight: 360,
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        const res = await getCouponListApi({
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
    keyField: 'coupon_id',
  },
  toolbarConfig: {
    enabled: false,
  },
});

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

function selectRow(row: CouponListItem) {
  emit('select', row);
  open.value = false;
}

const [Modal, modalApi] = useVbenModal({
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
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
    class="w-[800px]"
    title="选择优惠券"
  >
    <Grid>
      <template #total_num="{ row }">
        {{ row.total_num > 0 ? row.total_num : '无限制' }}
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="selectRow(row)">选择</ElButton>
      </template>
    </Grid>
  </Modal>
</template>
