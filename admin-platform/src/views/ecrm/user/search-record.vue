<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert, ElButton, ElMessage, ElMessageBox } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  clearUserSearchRecords,
  exportUserSearchRecords,
  listUserSearchRecords,
  type UserSearchRecord,
} from '#/api/core/platform-user-search';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canRead = ref(false);
const canClear = ref(false);
const canExport = ref(false);
const lastFilters = ref<Record<string, unknown>>({});

const sourceLabel: Record<string, string> = {
  pc: 'PC',
  h5: 'H5',
  mini: '小程序',
};

function buildFilters(formValues?: Record<string, unknown>) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const userIdRaw = String(formValues?.user_id ?? '').trim();
  return {
    user_id: userIdRaw ? Number(userIdRaw) : undefined,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    source: (String(formValues?.source ?? '').trim() ||
      undefined) as UserSearchRecord['source'] | undefined,
    start_date: range[0],
    end_date: range[1],
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '用户 ID' },
    fieldName: 'user_id',
    label: '用户 ID',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '搜索关键词' },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: 'PC', value: 'pc' },
        { label: 'H5', value: 'h5' },
        { label: '小程序', value: 'mini' },
      ],
      placeholder: '全部来源',
    },
    fieldName: 'source',
    label: '来源',
  },
]);

const gridOptions: VxeGridProps<UserSearchRecord> = {
  columns: [
    { field: 'id', title: '记录 ID', width: 100 },
    { field: 'user_id', title: '用户 ID', width: 100 },
    {
      field: 'keyword',
      minWidth: 220,
      showOverflow: false,
      title: '搜索关键词',
    },
    {
      field: 'source',
      formatter: ({ cellValue }) =>
        sourceLabel[String(cellValue)] || String(cellValue),
      title: '来源',
      width: 100,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 180,
      title: '搜索时间',
    },
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const filters = buildFilters(formValues);
        lastFilters.value = filters;
        const result = await listUserSearchRecords({
          page: page.currentPage,
          limit: page.pageSize,
          ...filters,
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

async function clearRecords() {
  try {
    const user = await ElMessageBox.prompt(
      '请输入需清理的用户 ID。本操作仅逻辑删除该用户当前可见搜索记录，不影响订单、收藏或浏览记录。',
      '按用户清理搜索记录',
      {
        inputPattern: /^[1-9]\d*$/,
        inputErrorMessage: '请输入正整数用户 ID',
      },
    );
    const reason = await ElMessageBox.prompt(
      '请填写清理原因（2 至 500 个字符）。',
      '清理原因',
      { inputPattern: /.{2,}/, inputErrorMessage: '清理原因至少 2 个字符' },
    );
    const result = await clearUserSearchRecords({
      user_id: Number(user.value),
      reason: reason.value.trim(),
      idempotency_key: `search-clear-${user.value}-${crypto.randomUUID()}`,
    });
    ElMessage.success(`已逻辑清理 ${result.cleared_count} 条搜索记录`);
    gridApi.reload();
  } catch {
    /* 用户取消 */
  }
}

async function exportRows() {
  try {
    const prompt = await ElMessageBox.prompt(
      '请填写导出原因。导出含用户 ID 与搜索词，最多 5000 条，请按最小必要原则使用。',
      '导出搜索记录',
      { inputPattern: /.{2,}/, inputErrorMessage: '导出原因至少 2 个字符' },
    );
    const filters = lastFilters.value;
    const result = await exportUserSearchRecords({
      user_id: filters.user_id as number | undefined,
      keyword: filters.keyword as string | undefined,
      source: filters.source as UserSearchRecord['source'] | undefined,
      start_date: filters.start_date as string | undefined,
      end_date: filters.end_date as string | undefined,
      reason: prompt.value.trim(),
    });
    const blob = new Blob([result.content], {
      type: 'text/csv;charset=utf-8',
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name;
    link.click();
    URL.revokeObjectURL(url);
    ElMessage.success(
      `已导出 ${result.row_count} 条搜索记录${result.truncated ? '（已按 5000 条上限截断）' : ''}`,
    );
  } catch {
    /* 用户取消 */
  }
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  const platform = profile.roles.includes('platform');
  canRead.value = platform && codes.includes('user.search_record.read');
  canClear.value = platform && codes.includes('user.search_record.clear');
  canExport.value = platform && codes.includes('user.search_record.export');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有搜索记录查看权限"
      type="warning"
      :closable="false"
    />
    <template v-else>
      <Grid>
        <template #toolbar-actions>
          <ElButton v-if="canClear" plain type="warning" @click="clearRecords">
            按用户清理
          </ElButton>
          <ElButton v-if="canExport" plain type="success" @click="exportRows">
            导出 CSV
          </ElButton>
        </template>
      </Grid>
      <ElAlert
        class="mt-4"
        title="导出会写入审计。CSV 使用 UTF-8 BOM 并对以 =、+、-、@ 开头的关键词转义，避免表格公式注入。"
        type="info"
        :closable="false"
      />
    </template>
  </Page>
</template>
