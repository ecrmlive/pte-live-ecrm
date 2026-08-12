<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { ElAlert } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listPlatformOperationLogs,
  type PlatformOperationLog,
} from '#/api/core/platform-operation-log';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canRead = ref(false);

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      maxlength: 128,
      placeholder: '请输入管理员ID/名称',
    },
    fieldName: 'admin_keyword',
    label: '管理员搜索',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: 'GET', value: 'GET' },
        { label: 'POST', value: 'POST' },
        { label: 'PUT', value: 'PUT' },
        { label: 'PATCH', value: 'PATCH' },
        { label: 'DELETE', value: 'DELETE' },
      ],
      placeholder: '请选择',
    },
    fieldName: 'request_method',
    label: '请求方式',
  },
  {
    ...LIST_DATE_RANGE_FIELD,
    label: '操作时间',
  },
]);

const gridOptions: VxeGridProps<PlatformOperationLog> = {
  columns: [
    { field: 'admin_user_id', title: '管理员 ID', width: 120 },
    { field: 'admin_name', title: '管理员姓名', width: 150 },
    {
      field: 'permission_name',
      title: '权限名称',
      width: 170,
    },
    {
      field: 'request',
      minWidth: 260,
      showOverflow: 'tooltip',
      title: '请求',
    },
    { field: 'request_method', title: '请求方式', width: 120 },
    {
      field: 'link',
      minWidth: 300,
      showOverflow: 'tooltip',
      title: '链接',
    },
    { field: 'ip', title: 'IP', width: 150 },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '操作时间',
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const result = await listPlatformOperationLogs({
          page: page.currentPage,
          limit: page.pageSize,
          admin_keyword:
            String(formValues?.admin_keyword ?? '').trim() || undefined,
          request_method:
            String(formValues?.request_method ?? '').trim() || undefined,
          start_date: range[0],
          end_date: range[1],
        });
        return { items: result.list || [], total: result.total || 0 };
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canRead.value =
    profile.roles.includes('platform') &&
    codes.includes('setting.operation_log.read');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有操作日志查看权限"
      type="warning"
      :closable="false"
    />
    <Grid v-else />
  </Page>
</template>
