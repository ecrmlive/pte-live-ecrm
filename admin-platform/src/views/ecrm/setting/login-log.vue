<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listPlatformLoginLogs,
  type PlatformLoginLog,
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
    componentProps: { clearable: true, maxlength: 64, placeholder: '登录账号' },
    fieldName: 'username',
    label: '账号',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '成功', value: 1 },
        { label: '失败', value: 0 },
      ],
      placeholder: '全部结果',
    },
    fieldName: 'success',
    label: '结果',
  },
]);

const gridOptions: VxeGridProps<PlatformLoginLog> = {
  columns: [
    { field: 'id', title: '日志 ID', width: 100 },
    { field: 'admin_user_id', title: '管理员 ID', width: 120 },
    { field: 'username', minWidth: 160, title: '账号' },
    { field: 'role_code', title: '角色', width: 130 },
    {
      field: 'success',
      slots: { default: 'success' },
      title: '结果',
      width: 100,
    },
    { field: 'ip', title: 'IP', width: 150 },
    {
      field: 'user_agent',
      minWidth: 220,
      showOverflow: false,
      title: '客户端标识',
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '时间',
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
        const successRaw = formValues?.success;
        const result = await listPlatformLoginLogs({
          page: page.currentPage,
          limit: page.pageSize,
          username: String(formValues?.username ?? '').trim() || undefined,
          success:
            successRaw === 0 || successRaw === 1
              ? (Number(successRaw) as 0 | 1)
              : undefined,
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
      title="当前账号没有登录日志查看权限"
      type="warning"
      :closable="false"
    />
    <Grid v-else>
      <template #success="{ row }">
        <ElTag :type="row.success ? 'success' : 'danger'">
          {{ row.success ? '成功' : '失败' }}
        </ElTag>
      </template>
    </Grid>
  </Page>
</template>
