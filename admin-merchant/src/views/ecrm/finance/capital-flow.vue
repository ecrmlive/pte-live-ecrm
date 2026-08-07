<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page } from '@vben/common-ui';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listFinanceLedgerApi,
  type FinanceLedgerEntry,
} from '#/api/core/merchant-ledger';
import { MERCHANT_LIST_GRID_LAYOUT } from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const entryLabels: Record<string, string> = {
  test_seed: '夹具入账',
  settlement_accrual: '结算应计',
  settlement_payout: '结算打款',
  withdraw: '提现扣减',
  refund_reversal: '退款冲销',
};

function entryLabel(type: string) {
  return entryLabels[type] || type;
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: Object.entries(entryLabels).map(([value, label]) => ({
        label,
        value,
      })),
      placeholder: '全部类型',
    },
    fieldName: 'entry_type',
    label: '流水类型',
  },
]);

const gridOptions: VxeGridProps<FinanceLedgerEntry> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    {
      field: 'entry_type',
      title: '类型',
      width: 140,
      formatter: ({ cellValue }) => entryLabel(String(cellValue || '')),
    },
    {
      field: 'amount',
      title: '金额',
      width: 120,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    { field: 'reference_type', title: '关联类型', width: 120 },
    {
      field: 'reference_id',
      minWidth: 140,
      showOverflow: false,
      title: '关联 ID',
    },
    {
      field: 'created_at',
      title: '时间',
      width: 180,
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const data = await listFinanceLedgerApi({
          page: page.currentPage,
          limit: page.pageSize,
          entry_type: String(formValues?.entry_type ?? '').trim() || undefined,
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
