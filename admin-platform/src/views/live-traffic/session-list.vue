<script lang="ts" setup>
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ElTooltip } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import LiveTrafficApi from '#/api/core/live-traffic';
import { Page } from '@vben/common-ui';
import {
  formatLiveApiDateTime,
  formatSettlementStatus,
} from '#/utils/live-api-time';
import { parseApiListPage } from '#/utils/list-response';
import {
  PLATFORM_SEARCH_SELECT_PROPS
} from '#/utils/platform-list-search-form';

interface SessionRow {
  app_id: number;
  app_name: string;
  live_id: number;
  room_name: string;
  session_id: string;
  session_start: string;
  session_end: string;
  lvb_play_gb: number;
  vod_play_gb: number;
  quota_gb_total: number;
  push_flux_gb: number;
  push_flux_note?: string;
  settlement_status: string;
  settle_error?: string;
  stream_name?: string;
}

const SETTLEMENT_OPTIONS = [
  { label: '待结算', value: 'pending' },
  { label: '已完成', value: 'done' },
  { label: '失败', value: 'failed' },
];

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
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '直播间 ID' },
      fieldName: 'live_id',
      formItemClass: 'pb-0',
      label: '直播间 ID',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: SETTLEMENT_OPTIONS,
        placeholder: '全部',
        ...PLATFORM_SEARCH_SELECT_PROPS,
      },
      fieldName: 'settlement_status',
      formItemClass: 'pb-0',
      label: '状态',
    },
  ],
};

async function fetchList(
  pageSize: number,
  currentPage: number,
  formValues?: Record<string, unknown>,
) {
  const params: Record<string, unknown> = {
    list_rows: pageSize,
    page: currentPage,
  };
  const appId = String(formValues?.app_id ?? '').trim();
  const liveId = String(formValues?.live_id ?? '').trim();
  const status = String(formValues?.settlement_status ?? '').trim();
  if (appId) params.app_id = Number(appId);
  if (liveId) params.live_id = Number(liveId);
  if (status) params.settlement_status = status;
  const res = await LiveTrafficApi.sessionList(params);
  const pageData = parseApiListPage<SessionRow>(res.data);
  return { items: pageData.list, total: pageData.total };
}

const gridOptions = {
  border: true,
  columns: [
    { field: 'app_id', title: '商城', width: 80 },
    {
      field: 'app_name',
      showOverflow: true,
      title: '名称',
      width: 120,
    },
    { field: 'live_id', title: '直播间', width: 80 },
    {
      field: 'room_name',
      minWidth: 120,
      showOverflow: true,
      title: '房间名',
    },
    {
      field: 'session_id',
      minWidth: 140,
      showOverflow: true,
      title: '场次',
    },
    {
      field: 'stream_name',
      minWidth: 120,
      showOverflow: true,
      title: '流名',
    },
    {
      field: 'session_start',
      slots: { default: 'sessionStart' },
      title: '开播',
      width: 170,
    },
    {
      field: 'session_end',
      slots: { default: 'sessionEnd' },
      title: '结束',
      width: 170,
    },
    {
      field: 'lvb_play_gb',
      slots: { default: 'lvbGb' },
      title: 'LVB',
      width: 80,
    },
    {
      field: 'vod_play_gb',
      slots: { default: 'vodGb' },
      title: 'VOD',
      width: 80,
    },
    {
      field: 'quota_gb_total',
      slots: { default: 'quotaGb' },
      title: '合计',
      width: 80,
    },
    {
      field: 'push_flux_gb',
      slots: { default: 'pushFlux', header: 'pushHeader' },
      width: 80,
    },
    {
      field: 'settlement_status',
      slots: { default: 'settlementStatus' },
      title: '状态',
      width: 90,
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 15, 20, 30, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) =>
        fetchList(page.pageSize, page.currentPage, formValues),
    },
  },
  rowConfig: { isHover: true, keyField: 'session_id' },
};

const [Grid] = useVbenVxeGrid({ formOptions, gridOptions });
</script>

<template>
  <Page>
    <Grid>
      <template #sessionStart="{ row }">
        {{ formatLiveApiDateTime(row.session_start) }}
      </template>

      <template #sessionEnd="{ row }">
        {{ formatLiveApiDateTime(row.session_end) }}
      </template>

      <template #lvbGb="{ row }">
        {{ fmt(row.lvb_play_gb) }}
      </template>

      <template #vodGb="{ row }">
        {{ fmt(row.vod_play_gb) }}
      </template>

      <template #quotaGb="{ row }">
        {{ fmt(row.quota_gb_total) }}
      </template>

      <template #pushHeader>
        <ElTooltip
          content="推流为推流域名级参考，非单场精确值，不计入配额"
          placement="top"
        >
          <span>推流</span>
        </ElTooltip>
      </template>

      <template #pushFlux="{ row }">
        {{ fmt(row.push_flux_gb) }}
      </template>

      <template #settlementStatus="{ row }">
        <ElTooltip
          v-if="row.settlement_status === 'failed' && row.settle_error"
          :content="row.settle_error"
          placement="top"
        >
          <span>{{ formatSettlementStatus(row.settlement_status) }}</span>
        </ElTooltip>
        <span v-else>{{ formatSettlementStatus(row.settlement_status) }}</span>
      </template>
    </Grid>
  </Page>
</template>
