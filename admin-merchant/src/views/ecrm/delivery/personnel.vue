<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { Page } from '@vben/common-ui';
import { ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { listMerchantStaffApi, type MerchantStaff } from '#/api/core/staff';
import { MERCHANT_LIST_GRID_LAYOUT } from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '姓名 / 手机号 / 账号',
    },
    fieldName: 'keyword',
    label: '配送员搜索',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '启用', value: 1 },
        { label: '停用', value: 0 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<MerchantStaff> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'service_id', title: 'ID', width: 80 },
    { field: 'nickname', minWidth: 140, showOverflow: false, title: '姓名' },
    { field: 'phone', title: '手机号', width: 140 },
    { field: 'role_code', title: '角色', width: 100 },
    {
      field: 'is_goods',
      title: '可发货',
      width: 90,
      formatter: ({ cellValue }) => (cellValue ? '是' : '否'),
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '创建时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const status = formValues?.status;
        const data = await listMerchantStaffApi({
          page: page.currentPage,
          limit: page.pageSize,
          keyword: String(formValues?.keyword ?? '').trim() || undefined,
          status: status === 0 || status === 1 ? Number(status) : undefined,
          staff_scope: 'delivery',
        });
        return { items: data.list ?? [], total: data.total ?? 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'service_id' },
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
          {{ row.status ? '启用' : '停用' }}
        </ElTag>
      </template>
    </Grid>
  </Page>
</template>
