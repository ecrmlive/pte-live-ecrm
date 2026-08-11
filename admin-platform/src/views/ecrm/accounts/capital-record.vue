<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElButton, ElMessage } from 'element-plus';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  exportPlatformUserBillsApi,
  listPlatformUserBillsApi,
  listPlatformUserBillTypesApi,
  type PlatformUserBillQuery,
  type PlatformUserBillRow,
} from '#/api/core/platform-user-bill';
import {
  listUserSearchFormField,
  parseUserSearch,
} from '#/components/ecrm/user-search-field';
import { platformListPagerConfig } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const lastFormValues = ref<Record<string, unknown>>({});
const canRead = ref(false);
const canExport = ref(false);
const exporting = ref(false);
const typeOptions = ref<{ label: string; value: string }[]>([]);

function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || String(v).trim() === '') return '-';
  return String(v);
}

function formatMoney(v?: number) {
  return Number(v || 0).toFixed(2);
}

function buildFilterParams(
  formValues?: Record<string, unknown>,
): PlatformUserBillQuery {
  const range = Array.isArray(formValues?.date_range)
    ? formValues.date_range
    : [];
  const userSearch = parseUserSearch(formValues);
  return {
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    type: String(formValues?.type ?? '').trim() || undefined,
    user_type: userSearch.type || 'nickname',
    user_keyword: userSearch.keyword || undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults(
  [
    {
      ...LIST_DATE_RANGE_FIELD,
      label: '时间选择',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        filterable: true,
        options: typeOptions,
        placeholder: '请选择',
      },
      fieldName: 'type',
      label: '明细类型',
    },
    listUserSearchFormField({ typeWidth: '96px' }),
  ],
  {
    commonConfig: { componentProps: { class: 'w-full' } },
    submitButtonOptions: { content: '搜索' },
    handleSubmit: async (values) => {
      lastFormValues.value = { ...values };
      await gridApi.reload(values);
    },
    handleReset: async () => {
      await formApi.resetForm();
      const values = (await formApi.getValues()) ?? {};
      lastFormValues.value = { ...values };
      await gridApi.reload(values);
    },
  },
);

const [Form, formApi] = useVbenForm(formOptions);

const gridOptions: VxeGridProps<PlatformUserBillRow> = {
  columns: [
    { field: 'uid', title: '会员ID', width: 100 },
    {
      field: 'nickname',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 140,
      showOverflow: 'tooltip',
      title: '昵称',
    },
    {
      field: 'number',
      formatter: ({ cellValue }) => formatMoney(Number(cellValue)),
      minWidth: 120,
      title: '金额',
    },
    {
      field: 'title',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 160,
      showOverflow: 'tooltip',
      title: '明细类型',
    },
    {
      field: 'mark',
      className: 'col--remark',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 200,
      showOverflow: 'tooltip',
      title: '备注',
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '-',
      minWidth: 170,
      title: '创建时间',
    },
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const values =
          formValues && Object.keys(formValues).length > 0
            ? formValues
            : lastFormValues.value;
        const data = await listPlatformUserBillsApi({
          page: page.currentPage,
          limit: page.pageSize,
          ...buildFilterParams(values),
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'bill_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

async function exportRows() {
  if (!canExport.value || exporting.value) return;
  exporting.value = true;
  try {
    const result = await exportPlatformUserBillsApi(
      buildFilterParams(lastFormValues.value),
    );
    const blob = new Blob([result.content || ''], {
      type: 'text/csv;charset=utf-8',
    });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = result.file_name || '资金记录.csv';
    a.click();
    URL.revokeObjectURL(url);
    ElMessage.success(
      result.truncated
        ? `已导出前 ${result.row_count} 条（已截断）`
        : `已导出 ${result.row_count} 条`,
    );
  } finally {
    exporting.value = false;
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canRead.value = codes.includes('accounts.user_assets.read');
  canExport.value = codes.includes('accounts.user_assets.export');
  if (canRead.value) {
    try {
      const types = await listPlatformUserBillTypesApi();
      typeOptions.value = (types.list || []).map((item) => ({
        label: item.title,
        value: item.type,
      }));
      formApi.updateSchema([
        {
          fieldName: 'type',
          componentProps: {
            clearable: true,
            filterable: true,
            options: typeOptions.value,
            placeholder: '请选择',
          },
        },
      ]);
    } catch {
      typeOptions.value = [];
    }
    await gridApi.reload();
  }
});
</script>

<template>
  <Page auto-content-height>
    <div class="capital-filter">
      <Form />
    </div>

    <Grid>
      <template #toolbar-actions>
        <ElButton
          v-if="canExport"
          type="primary"
          :loading="exporting"
          @click="exportRows"
        >
          导出列表
        </ElButton>
      </template>
    </Grid>
  </Page>
</template>

<style scoped>
.capital-filter {
  padding: 12px 8px 4px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}
</style>
