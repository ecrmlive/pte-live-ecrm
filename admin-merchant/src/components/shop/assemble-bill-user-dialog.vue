<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { useVbenModal } from '@vben/common-ui';
import { reactive, watch } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAssembleRecordDetailApi } from '#/api/core/plus-assemble';

defineOptions({ name: 'AssembleBillUserDialog' });

type AssembleUserRow = Record<string, unknown> & {
  assemble_bill_user_id?: number;
  avatarUrl?: string;
  create_time?: string;
  is_virtual?: number;
  nickName?: string;
  orderM?: { order_no?: string; pay_price?: string; state_text?: string };
  user_id?: number;
};

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  assembleBillId?: number;
}>();

const gridOptions = reactive<VxeGridProps<AssembleUserRow>>({
  columns: [
    { field: 'assemble_bill_user_id', title: 'ID', width: 80 },
    {
      field: 'nickName',
      minWidth: 200,
      slots: { default: 'user' },
      title: '用户信息',
    },
    {
      field: 'avatarUrl',
      slots: { default: 'avatar' },
      title: '头像',
      width: 80,
    },
    {
      field: 'order_no',
      minWidth: 150,
      slots: { default: 'order_no' },
      title: '订单编号',
    },
    {
      field: 'pay_price',
      minWidth: 100,
      slots: { default: 'pay_price' },
      title: '订单金额',
    },
    {
      field: 'state_text',
      minWidth: 100,
      slots: { default: 'state_text' },
      title: '订单状态',
    },
    {
      field: 'is_virtual',
      minWidth: 90,
      slots: { default: 'is_virtual' },
      title: '虚拟成团',
    },
    { field: 'create_time', title: '参团时间', width: 150 },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 20,
    pageSizes: [10, 20, 50],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!props.assembleBillId) {
          return { items: [], total: 0 };
        }
        const res = await getAssembleRecordDetailApi({
          assemble_bill_id: props.assembleBillId,
          list_rows: page.pageSize,
          page: page.currentPage,
        });
        return {
          items: res.list.data as AssembleUserRow[],
          total: res.list.total,
        };
      },
    },
  },
  rowConfig: {
    keyField: 'assemble_bill_user_id',
  },
  toolbarConfig: {
    enabled: false,
  },
});

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

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
    :close-on-click-modal="true"
    :destroy-on-close="true"
    class="w-[1000px]"
    title="拼团详情"
  >
    <Grid>
      <template #user="{ row }">
        昵称：{{ row.nickName }}(ID：{{ row.user_id }})
      </template>
      <template #avatar="{ row }">
        <img
          v-if="row.avatarUrl"
          :src="String(row.avatarUrl)"
          alt=""
          class="size-8 rounded object-cover"
        />
      </template>
      <template #order_no="{ row }">
        {{ row.orderM?.order_no }}
      </template>
      <template #pay_price="{ row }">
        {{ row.orderM?.pay_price }}
      </template>
      <template #state_text="{ row }">
        {{ row.orderM?.state_text }}
      </template>
      <template #is_virtual="{ row }">
        {{ row.is_virtual === 1 ? '是' : '否' }}
      </template>
    </Grid>
  </Modal>
</template>
