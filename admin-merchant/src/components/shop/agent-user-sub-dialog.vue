<script setup lang="ts">
import type { VxeGridProps } from '#/adapter/vxe-table';

import { useVbenDrawer } from '@vben/common-ui';
import { reactive, ref, watch } from 'vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAgentUserFansApi } from '#/api/core/plus-agent';

defineOptions({ name: 'AgentUserSubDialog' });

type FanRow = Record<string, unknown> & {
  avatarUrl?: string;
  create_time?: string;
  mobile?: string;
  nickName?: string;
  user_id?: number;
};

const open = defineModel<boolean>('open', { default: false });

const props = defineProps<{
  userId?: number;
}>();

const level = ref('1');

const gridOptions = reactive<VxeGridProps<FanRow>>({
  columns: [
    { field: 'user_id', title: '用户ID', width: 90 },
    {
      field: 'avatarUrl',
      slots: { default: 'avatar' },
      title: '头像',
      width: 80,
    },
    { field: 'nickName', minWidth: 120, showOverflow: true, title: '昵称' },
    { field: 'mobile', minWidth: 120, title: '手机号' },
    { field: 'create_time', title: '加入时间', width: 150 },
  ],
  minHeight: 360,
  pagerConfig: {
    pageSize: 10,
    pageSizes: [10, 20, 50],
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }) => {
        if (!props.userId) {
          return { items: [], total: 0 };
        }
        const res = await getAgentUserFansApi({
          level: level.value,
          list_rows: page.pageSize,
          page: page.currentPage,
          user_id: props.userId,
        });
        return {
          items: res.list.data as FanRow[],
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

function onLevelChange() {
  void gridApi.reload();
}

const [Modal, modalApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  onOpenChange(isOpen) {
    open.value = isOpen;
  },
});

watch(open, (visible) => {
  if (visible) {
    level.value = '1';
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
    class="w-[800px]"
    title="下级用户"
  >
    <el-radio-group v-model="level" size="small" @change="onLevelChange">
      <el-radio-button value="1">一级</el-radio-button>
      <el-radio-button value="2">二级</el-radio-button>
    </el-radio-group>
    <Grid v-if="open" class="mt-4">
      <template #avatar="{ row }">
        <img
          v-if="row.avatarUrl"
          :src="String(row.avatarUrl)"
          alt=""
          class="size-8 rounded object-cover"
        />
      </template>
    </Grid>
  </Modal>
</template>
