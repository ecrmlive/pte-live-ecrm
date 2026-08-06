<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page } from '@vben/common-ui';
import { ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  listMerchantExpressApi,
  type MerchantExpressRow,
} from '#/api/core/merchant-logistics';
import { MERCHANT_LIST_GRID_LAYOUT } from '#/constants/merchant-list-grid';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '名称 / 编码' },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '展示', value: 1 },
        { label: '隐藏', value: 0 },
      ],
      placeholder: '全部',
    },
    fieldName: 'is_show',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<MerchantExpressRow> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'express_id', title: 'ID', width: 90 },
    { field: 'name', minWidth: 160, showOverflow: false, title: '名称' },
    { field: 'code', minWidth: 140, showOverflow: false, title: '编码' },
    { field: 'sort', title: '排序', width: 90 },
    {
      field: 'is_show',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
  ],
  pagerConfig: { enabled: true, pageSize: 50, pageSizes: [20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const isShow = formValues?.is_show;
        const data = await listMerchantExpressApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          is_show: isShow === 0 || isShow === 1 ? Number(isShow) : undefined,
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'express_id' },
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
        <ElTag :type="row.is_show === 1 ? 'success' : 'info'">
          {{ row.is_show === 1 ? '展示' : '隐藏' }}
        </ElTag>
      </template>
    </Grid>
  </Page>
</template>
