<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page } from '@vben/common-ui';
import { ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listFinanceStatementsApi,
  type FinanceStatement,
  type FinanceStatementStatus,
} from '#/api/core/merchant-ledger';
import {
  MERCHANT_LIST_GRID_LAYOUT,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const statusLabels: Record<FinanceStatementStatus, string> = {
  approved: '已审核',
  bill_frozen: '账期已冻结',
  bill_pending: '账期待生成',
  paid: '已打款',
  rejected: '已拒绝',
  withdraw_applied: '待平台审核',
};

const statusTypes: Record<
  FinanceStatementStatus,
  'danger' | 'info' | 'success' | 'warning'
> = {
  approved: 'success',
  bill_frozen: 'warning',
  bill_pending: 'info',
  paid: 'success',
  rejected: 'danger',
  withdraw_applied: 'warning',
};

function statusLabel(status: string) {
  return statusLabels[status as FinanceStatementStatus] || status;
}

function statusType(status: string) {
  return statusTypes[status as FinanceStatementStatus] || 'info';
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: Object.entries(statusLabels).map(([value, label]) => ({
        label,
        value,
      })),
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '账单状态',
  },
]);

const gridOptions: VxeGridProps<FinanceStatement> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'statement_id', title: '账单 ID', width: 100 },
    {
      field: 'period_start',
      minWidth: 220,
      showOverflow: false,
      slots: { default: 'period' },
      title: '账期',
    },
    {
      field: 'amount',
      title: '金额',
      width: 120,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 120,
    },
    {
      field: 'updated_at',
      minWidth: 170,
      title: '更新时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = String(formValues?.status ?? '').trim();
        const data = await listFinanceStatementsApi({
          page: page.currentPage,
          limit: page.pageSize,
          status: (status || undefined) as FinanceStatementStatus | undefined,
          date_from: range[0],
          date_to: range[1],
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'statement_id' },
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
    <Grid>
      <template #period="{ row }">
        {{ row.period_start }} ~ {{ row.period_end }}
      </template>
      <template #status="{ row }">
        <ElTag :type="statusType(row.status)">{{ statusLabel(row.status) }}</ElTag>
      </template>
    </Grid>
  </Page>
</template>
