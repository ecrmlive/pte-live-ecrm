<script lang="ts" setup>
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import {
  ElButton,
  ElMessage,
  ElMessageBox,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  batchDeleteJwtSessionsApi,
  deleteJwtSessionApi,
  listJwtSessionsApi,
  type JwtSessionRow,
} from '#/api/core/auth-session';
import { Page } from '@vben/common-ui';
import {
  formatUnixTime,
  LOGIN_PLATFORM_OPTIONS,
  USER_ROLE_OPTIONS,
} from '#/utils/live-token-refresh';
import {
  PLATFORM_SEARCH_SELECT_PROPS
} from '#/utils/platform-list-search-form';

const selectedIds = ref<number[]>([]);

const formOptions: VbenFormProps = {
  showCollapseButton: false,
  schema: [
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: LOGIN_PLATFORM_OPTIONS,
        placeholder: '全部',
        ...PLATFORM_SEARCH_SELECT_PROPS,
      },
      fieldName: 'login_platform',
      formItemClass: 'pb-0',
      label: '登录平台',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '用户ID',
      },
      fieldName: 'uid',
      formItemClass: 'pb-0',
      label: '用户ID',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '昵称',
      },
      fieldName: 'nick_name',
      formItemClass: 'pb-0',
      label: '用户昵称',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: USER_ROLE_OPTIONS,
        placeholder: '全部',
        ...PLATFORM_SEARCH_SELECT_PROPS,
      },
      fieldName: 'user_role',
      formItemClass: 'pb-0',
      label: '用户角色',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '商户名称',
      },
      fieldName: 'shop_name',
      formItemClass: 'pb-0',
      label: '商户名称',
    },
  ],
};

async function fetchSessionPage(
  pageSize: number,
  currentPage: number,
  formValues?: Record<string, unknown>,
) {
  const res = await listJwtSessionsApi({
    page: currentPage,
    list_rows: pageSize,
    login_platform: String(formValues?.login_platform ?? '') || undefined,
    uid: String(formValues?.uid ?? '').trim() || undefined,
    nick_name: String(formValues?.nick_name ?? '').trim() || undefined,
    user_role: String(formValues?.user_role ?? '') || undefined,
    shop_name: String(formValues?.shop_name ?? '').trim() || undefined,
  });
  selectedIds.value = [];
  gridApi.grid?.clearCheckboxRow?.();
  return {
    items: res.data?.list ?? [],
    total: res.data?.total ?? 0,
  };
}

const gridOptions = {
  border: true,
  checkboxConfig: { highlight: true, reserve: true },
  columns: [
    { type: 'checkbox', width: 48 },
    { field: 'uid', title: '用户ID', width: 90 },
    {
      field: 'nick_name',
      minWidth: 120,
      showOverflow: true,
      title: '用户昵称',
    },
    {
      field: 'login_platform_label',
      title: '登录平台',
      width: 120,
    },
    {
      field: 'shop_name',
      minWidth: 120,
      showOverflow: true,
      title: '商户名称',
    },
    { field: 'user_role', title: '用户角色', width: 90 },
    {
      field: 'jwt_token',
      minWidth: 200,
      showOverflow: true,
      slots: { default: 'jwtToken' },
      title: 'JWT Token',
    },
    {
      field: 'expires_at',
      slots: { default: 'expiresAt' },
      title: '有效时间',
      width: 170,
    },
    {
      field: 'last_login_time',
      slots: { default: 'lastLoginTime' },
      title: '最后登录时间',
      width: 170,
    },
    {
      fixed: 'right',
      slots: { default: 'action' },
      title: '操作',
      width: 90,
    },
  ],
  pagerConfig: { enabled: true, pageSize: 15, pageSizes: [10, 15, 20, 30, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) =>
        fetchSessionPage(page.pageSize, page.currentPage, formValues),
    },
  },
  rowConfig: { isHover: true, keyField: 'session_id' },
};

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions,
  gridEvents: {
    checkboxAll: syncSelection,
    checkboxChange: syncSelection,
  },
  gridOptions,
});

function syncSelection() {
  const rows = (gridApi.grid?.getCheckboxRecords?.() ?? []) as JwtSessionRow[];
  selectedIds.value = rows.map((row) => row.session_id);
}

function reload() {
  gridApi.reload();
}

async function handleDelete(row: JwtSessionRow) {
  try {
    await ElMessageBox.confirm(
      `确定删除用户 ${row.nick_name || row.uid} 的登录会话？`,
      '删除确认',
      { type: 'warning' },
    );
    await deleteJwtSessionApi(row.session_id);
    ElMessage.success('删除成功');
    reload();
  } catch {
    // cancelled or failed
  }
}

async function handleBatchDelete() {
  if (!selectedIds.value.length) {
    ElMessage.warning('请先选择要删除的会话');
    return;
  }
  try {
    await ElMessageBox.confirm(
      `确定批量删除选中的 ${selectedIds.value.length} 条登录会话？`,
      '批量删除',
      { type: 'warning' },
    );
    await batchDeleteJwtSessionsApi(selectedIds.value);
    ElMessage.success('删除成功');
    reload();
  } catch {
    // cancelled or failed
  }
}

function shortToken(token: string) {
  if (!token) return '—';
  if (token.length <= 24) return token;
  return `${token.slice(0, 12)}…${token.slice(-8)}`;
}
</script>

<template>
  <Page>
    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-access:code="'platform:authSession:batchDelete'"
          :disabled="!selectedIds.length"
          type="danger"
          @click="handleBatchDelete"
        >
          批量删除
        </ElButton>
      </template>

      <template #jwtToken="{ row }">
        <span :title="row.jwt_token">{{ shortToken(row.jwt_token) }}</span>
      </template>

      <template #expiresAt="{ row }">
        {{ formatUnixTime(row.expires_at) }}
      </template>

      <template #lastLoginTime="{ row }">
        {{ formatUnixTime(row.last_login_time) }}
      </template>

      <template #action="{ row }">
        <ElButton
          v-access:code="'platform:authSession:delete'"
          link
          type="danger"
          @click="handleDelete(row)"
        >
          删除
        </ElButton>
      </template>
    </Grid>
  </Page>
</template>
