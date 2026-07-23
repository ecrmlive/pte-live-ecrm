<script lang="ts" setup>
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref, watch } from 'vue';

import { Page } from '@vben/common-ui';

import { ElButton, ElTabPane, ElTabs } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import LiveTrafficApi from '#/api/core/live-traffic';
import {
  PLATFORM_LIST_GRID_CLASS,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { parseApiListPage } from '#/utils/list-response';

import ShopTrafficManageModal from '../shop/shop-traffic-manage-modal.vue';

interface RankRow {
  app_id: number;
  app_name: string;
  used_gb_total: number;
  remain_gb: number;
  total_gb: number;
  operator_name: string;
}

type RankTab = 'consumption' | 'recharge';

const activeTab = ref<RankTab>('consumption');
const trafficModalOpen = ref(false);
const manageShop = ref<Pick<RankRow, 'app_id' | 'app_name'> | null>(null);

function fmt(v: number | string | undefined) {
  return Number(v || 0).toFixed(2);
}

const formOptions: VbenFormProps = {
  showCollapseButton: false,
  schema: [
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '商城 ID' },
      fieldName: 'app_id',
      formItemClass: 'pb-0',
      label: '商城 ID',
    },
  ],
};

async function fetchRankPage(
  rankType: RankTab,
  pageSize: number,
  currentPage: number,
  formValues?: Record<string, unknown>,
) {
  const params: Record<string, unknown> = {
    rank_type: rankType,
    page: currentPage,
    list_rows: pageSize,
  };
  const appId = String(formValues?.app_id ?? '').trim();
  if (appId) {
    params.app_id = Number(appId);
  }
  const res = await LiveTrafficApi.usageRank(params);
  const page = parseApiListPage<RankRow>(res.data);
  return { items: page.list, total: page.total };
}

function buildRankColumns(): VxeGridProps<RankRow>['columns'] {
  return [
    { field: 'app_id', title: '商城 ID', width: 90 },
    {
      field: 'app_name',
      minWidth: 140,
      showOverflow: true,
      title: '商城名称',
    },
    {
      field: 'used_gb_total',
      formatter: ({ cellValue }) => fmt(cellValue),
      title: '消费总流量(GB)',
      width: 130,
    },
    {
      field: 'remain_gb',
      formatter: ({ cellValue }) => fmt(cellValue),
      title: '剩余总流量(GB)',
      width: 130,
    },
    {
      field: 'total_gb',
      formatter: ({ cellValue }) => fmt(cellValue),
      title: '充值总流量(GB)',
      width: 130,
    },
    {
      field: 'operator_name',
      minWidth: 100,
      showOverflow: true,
      title: '操作员',
      width: 100,
    },
    {
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 90,
    },
  ];
}

const rankGridBase: VxeGridProps<RankRow> = {
  border: true,
  height: 'auto',
  minHeight: 0,
  columns: buildRankColumns(),
  pagerConfig: platformListPagerConfig(),
  rowConfig: { isHover: true, keyField: 'app_id' },
};

const consumptionGridOptions: VxeGridProps<RankRow> = {
  ...rankGridBase,
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) =>
        fetchRankPage('consumption', page.pageSize, page.currentPage, formValues),
    },
  },
};

const rechargeGridOptions: VxeGridProps<RankRow> = {
  ...rankGridBase,
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) =>
        fetchRankPage('recharge', page.pageSize, page.currentPage, formValues),
    },
  },
};

const [ConsumptionGrid, consumptionGridApi] = useVbenVxeGrid({
  formOptions,
  gridClass: PLATFORM_LIST_GRID_CLASS,
  gridOptions: consumptionGridOptions,
});
const [RechargeGrid, rechargeGridApi] = useVbenVxeGrid({
  formOptions,
  gridClass: PLATFORM_LIST_GRID_CLASS,
  gridOptions: rechargeGridOptions,
});

function openTrafficManage(row: RankRow) {
  manageShop.value = row;
  trafficModalOpen.value = true;
}

function reloadActiveGrid() {
  if (activeTab.value === 'consumption') {
    consumptionGridApi.reload();
    return;
  }
  rechargeGridApi.reload();
}

watch(trafficModalOpen, (open, wasOpen) => {
  if (wasOpen && !open) {
    reloadActiveGrid();
  }
});
</script>

<template>
  <Page auto-content-height content-class="traffic-rank-page flex flex-col overflow-hidden min-h-0">
    <ElTabs v-model="activeTab" class="traffic-rank-tabs shrink-0">
      <ElTabPane label="消耗榜单" name="consumption" />
      <ElTabPane label="充值榜单" name="recharge" />
    </ElTabs>

    <ConsumptionGrid v-if="activeTab === 'consumption'" class="traffic-rank-grid min-h-0 flex-1">
      <template #action="{ row }">
        <ElButton link type="primary" @click="openTrafficManage(row)">
          充值
        </ElButton>
      </template>
    </ConsumptionGrid>

    <RechargeGrid v-else class="traffic-rank-grid min-h-0 flex-1">
      <template #action="{ row }">
        <ElButton link type="primary" @click="openTrafficManage(row)">
          充值
        </ElButton>
      </template>
    </RechargeGrid>

    <ShopTrafficManageModal
      v-if="manageShop"
      v-model:open="trafficModalOpen"
      :app-id="manageShop.app_id"
      :app-name="manageShop.app_name"
    />
  </Page>
</template>

<style scoped>
.traffic-rank-tabs :deep(.el-tabs__header) {
  margin-bottom: 12px;
}

.traffic-rank-grid {
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.traffic-rank-grid :deep(.vxe-grid--toolbar-wrapper),
.traffic-rank-grid :deep(.vxe-grid--pager-wrapper) {
  flex-shrink: 0;
}
</style>
