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
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '管理员 ID' },
    fieldName: 'admin_user_id',
    label: '管理员 ID',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, maxlength: 32, placeholder: '角色代码' },
    fieldName: 'role_code',
    label: '角色',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      maxlength: 128,
      placeholder: '如 POST /coupons',
    },
    fieldName: 'action',
    label: '操作',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, maxlength: 64, placeholder: '资源类型' },
    fieldName: 'resource_type',
    label: '资源',
  },
]);

const gridOptions: VxeGridProps<PlatformOperationLog> = {
  columns: [
    { field: 'id', title: '日志 ID', width: 100 },
    { field: 'admin_user_id', title: '管理员 ID', width: 120 },
    { field: 'role_code', title: '角色', width: 130 },
    {
      field: 'action',
      minWidth: 300,
      showOverflow: false,
      title: '成功操作',
    },
    { field: 'resource_type', title: '资源类型', width: 140 },
    {
      field: 'resource_id',
      formatter: ({ cellValue }) => cellValue || '—',
      title: '资源 ID',
      width: 120,
    },
    {
      field: 'request_id',
      minWidth: 220,
      showOverflow: false,
      title: '请求号',
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '操作时间',
    },
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const adminIdRaw = String(formValues?.admin_user_id ?? '').trim();
        const result = await listPlatformOperationLogs({
          page: page.currentPage,
          limit: page.pageSize,
          admin_user_id: adminIdRaw ? Number(adminIdRaw) : undefined,
          role_code: String(formValues?.role_code ?? '').trim() || undefined,
          action: String(formValues?.action ?? '').trim() || undefined,
          resource_type:
            String(formValues?.resource_type ?? '').trim() || undefined,
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
