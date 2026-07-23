<script lang="ts" setup>
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import LiveTrafficApi from '#/api/core/live-traffic';
import { Page } from '@vben/common-ui';
import { formatLiveApiDateTime } from '#/utils/live-api-time';
import { parseApiListPage } from '#/utils/list-response';

interface LiveSummaryRow {
  app_id: number;
  app_name: string;
  live_id: number;
  room_name: string;
  session_count: number;
  lvb_play_gb: number;
  vod_play_gb: number;
  quota_gb_total: number;
  last_session_end: string;
}

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
  if (appId) params.app_id = Number(appId);
  const res = await LiveTrafficApi.liveList(params);
  const pageData = parseApiListPage<LiveSummaryRow>(res.data);
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
    { field: 'room_name', minWidth: 140, title: '房间名' },
    { field: 'session_count', title: '场次数', width: 80 },
    {
      field: 'lvb_play_gb',
      slots: { default: 'lvbGb' },
      title: 'LVB合计',
      width: 90,
    },
    {
      field: 'vod_play_gb',
      slots: { default: 'vodGb' },
      title: 'VOD合计',
      width: 90,
    },
    {
      field: 'quota_gb_total',
      slots: { default: 'quotaGb' },
      title: '配额合计',
      width: 90,
    },
    {
      field: 'last_session_end',
      slots: { default: 'lastSessionEnd' },
      title: '最近场次结束',
      width: 170,
    },
  ],
  pagerConfig: { enabled: true, pageSize: 15, pageSizes: [10, 15, 20, 30, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) =>
        fetchList(page.pageSize, page.currentPage, formValues),
    },
  },
  rowConfig: { isHover: true, keyField: 'live_id' },
};

const [Grid] = useVbenVxeGrid({ formOptions, gridOptions });
</script>

<template>
  <Page>
    <Grid>
      <template #lvbGb="{ row }">
        {{ fmt(row.lvb_play_gb) }}
      </template>

      <template #vodGb="{ row }">
        {{ fmt(row.vod_play_gb) }}
      </template>

      <template #quotaGb="{ row }">
        {{ fmt(row.quota_gb_total) }}
      </template>

      <template #lastSessionEnd="{ row }">
        {{ formatLiveApiDateTime(row.last_session_end) }}
      </template>
    </Grid>
  </Page>
</template>
