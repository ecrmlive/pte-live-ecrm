<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page } from '@vben/common-ui';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listUserSearchRecordsApi,
  type UserSearchRecord,
} from '#/api/core/merchant-search-record';
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
    componentProps: { clearable: true, placeholder: '商品名称' },
    fieldName: 'keyword',
    label: '商品搜索',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '用户 ID' },
    fieldName: 'user_id',
    label: '用户 ID',
  },
]);

const gridOptions: VxeGridProps<UserSearchRecord> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    { field: 'user_id', title: '用户 ID', width: 90 },
    {
      field: 'product_title',
      minWidth: 200,
      showOverflow: false,
      slots: { default: 'product' },
      title: '商品',
    },
    {
      field: 'viewed_at',
      title: '浏览时间',
      width: 180,
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const userID = Number(formValues?.user_id);
        const data = await listUserSearchRecordsApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          user_id:
            Number.isFinite(userID) && userID > 0 ? userID : undefined,
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
    <Grid>
      <template #product="{ row }">
        {{ row.product_title || `#${row.product_id}` }}
      </template>
    </Grid>
  </Page>
</template>
