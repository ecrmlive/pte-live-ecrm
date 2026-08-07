<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page } from '@vben/common-ui';
import { ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listStoreCustomersApi,
  type StoreCustomer,
} from '#/api/core/merchant-store-customer';
import { MERCHANT_LIST_GRID_LAYOUT } from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '用户 ID 或昵称',
    },
    fieldName: 'keyword',
    label: '关键词',
  },
]);

const gridOptions: VxeGridProps<StoreCustomer> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'user_id', title: '用户 ID', width: 90 },
    {
      field: 'nickname',
      minWidth: 140,
      showOverflow: false,
      title: '昵称',
    },
    { field: 'mobile', title: '手机号', width: 140 },
    { field: 'order_count', title: '订单数', width: 90 },
    {
      field: 'total_pay',
      title: '累计消费',
      width: 120,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'last_order_at',
      minWidth: 170,
      title: '最近下单',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listStoreCustomersApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
        });
        return { items: data.list ?? [], total: data.total ?? 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'user_id' },
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
      <template #status="{ row }">
        <ElTag :type="row.status ? 'success' : 'info'">
          {{ row.status ? '正常' : '停用' }}
        </ElTag>
      </template>
    </Grid>
  </Page>
</template>
