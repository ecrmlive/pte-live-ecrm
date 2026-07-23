<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { useVbenModal } from '@vben/common-ui';
import { computed, reactive, watch } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getBargainProductDetailApi, getBargainTaskDetailApi } from '#/api/core/plus-bargain';

defineOptions({ name: 'BargainTaskDetailDialog' });

type BargainDetailRow = Record<string, unknown> & {
  create_time?: string;
  cut_money?: number | string;
  user?: { avatarUrl?: string; nickName?: string };
  user_id?: number;
};

const open = defineModel<boolean>('open', { default: false });

const props = withDefaults(
  defineProps<{
    apiSource?: 'product' | 'task';
    bargainTaskId?: number;
  }>(),
  { apiSource: 'task' },
);

const dialogWidthClass = computed(() =>
  props.apiSource === 'product' ? 'w-[1000px]' : 'w-[800px]',
);

const gridOptions = reactive<VxeGridProps<BargainDetailRow>>({
  columns: [
    { field: 'user_id', title: '用户ID', width: 90 },
    {
      field: 'avatar',
      slots: { default: 'avatar' },
      title: '用户头像',
      width: 90,
    },
    {
      field: 'nickName',
      slots: { default: 'nickName' },
      minWidth: 120,
      title: '用户昵称',
    },
    {
      field: 'cut_money',
      slots: { default: 'cut_money' },
      minWidth: 100,
      title: '砍价金额',
    },
    { field: 'create_time', title: '砍价时间', width: 150 },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!props.bargainTaskId) {
          return { items: [], total: 0 };
        }
        const params = {
          bargain_task_id: props.bargainTaskId,
          list_rows: page.pageSize,
          page: page.currentPage,
        };
        const res =
          props.apiSource === 'product'
            ? await getBargainProductDetailApi(params)
            : await getBargainTaskDetailApi(params);
        return {
          items: res.list.data as BargainDetailRow[],
          total: res.list.total,
        };
      },
    },
  },
  rowConfig: {
    keyField: 'user_id',
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
    :close-on-click-modal="false"
    :destroy-on-close="true"
    :class="dialogWidthClass"
    title="砍价详情"
  >
    <Grid>
      <template #avatar="{ row }">
        <img
          v-if="row.user?.avatarUrl"
          :src="row.user.avatarUrl"
          alt=""
          class="size-8 rounded object-cover"
        />
      </template>
      <template #nickName="{ row }">
        {{ row.user?.nickName }}
      </template>
      <template #cut_money="{ row }">
        <span class="font-semibold text-red-500">{{ row.cut_money }}</span>
      </template>
    </Grid>
  </Modal>
</template>
