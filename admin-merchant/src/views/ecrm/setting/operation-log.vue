<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page } from '@vben/common-ui';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listMerchantOperationLogsApi,
  type MerchantOperationLog,
} from '#/api/core/merchant-operation-log';
import { MERCHANT_LIST_GRID_LAYOUT } from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '动作 / 资源类型 / 请求 ID' },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '动作标识' },
    fieldName: 'action',
    label: '动作',
  },
]);

const gridOptions: VxeGridProps<MerchantOperationLog> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    { field: 'account_id', title: '操作人', width: 90 },
    { field: 'action', minWidth: 140, showOverflow: false, title: '动作' },
    { field: 'resource_type', title: '资源类型', width: 120 },
    { field: 'resource_id', title: '资源 ID', width: 120 },
    { field: 'request_id', minWidth: 160, showOverflow: false, title: '请求 ID' },
    {
      field: 'created_at',
      title: '时间',
      width: 180,
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const data = await listMerchantOperationLogsApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          action: String(formValues?.action ?? '').trim() || undefined,
          date_from: range[0],
          date_to: range[1],
        });
        return { items: data.list ?? [], total: data.total ?? 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid] = useVbenVxeGrid({ formOptions, gridOptions });
</script>

<template>
  <Page auto-content-height>
    <Grid />
  </Page>
</template>
