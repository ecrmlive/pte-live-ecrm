<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import { ElButton, ElMessage } from 'element-plus';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  exportPlatformIntegralLogsApi,
  getPlatformIntegralLogTitleApi,
  listPlatformIntegralLogsApi,
  type PlatformIntegralLogRow,
  type PlatformIntegralLogTitle,
} from '#/api/core/platform-integral-log';
import { platformListPagerConfig } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const emptyTitle = (): PlatformIntegralLogTitle => ({
  total_integral: 0,
  sign_count: 0,
  sign_integral: 0,
  used_integral: 0,
  order_integral: 0,
  freeze_integral: 0,
});

const title = ref<PlatformIntegralLogTitle>(emptyTitle());
const lastFormValues = ref<Record<string, unknown>>({});
const exporting = ref(false);

// 图标须落在离线 Iconify bundle（public/iconify/lucide.json / ant-design.json）
const summaryCards = computed(() => [
  {
    key: 'total_integral',
    label: '总积分',
    value: formatStat(title.value.total_integral),
    icon: 'lucide:award',
    tone: 'blue',
  },
  {
    key: 'sign_count',
    label: '客户签到次数',
    value: String(title.value.sign_count || 0),
    icon: 'lucide:calendar-clock',
    tone: 'orange',
  },
  {
    key: 'sign_integral',
    label: '签到送出积分',
    value: formatStat(title.value.sign_integral),
    icon: 'lucide:gift',
    tone: 'green',
  },
  {
    key: 'used_integral',
    label: '使用积分',
    value: formatStat(title.value.used_integral),
    icon: 'lucide:shopping-bag',
    tone: 'pink',
  },
  {
    key: 'order_integral',
    label: '下单赠送积分',
    value: formatStat(title.value.order_integral),
    icon: 'lucide:package',
    tone: 'purple',
  },
  {
    key: 'freeze_integral',
    label: '冻结积分',
    value: formatStat(title.value.freeze_integral),
    icon: 'lucide:cloud',
    tone: 'sky',
  },
]);

function formatStat(v?: number) {
  const n = Number(v || 0);
  if (Number.isInteger(n)) return String(n);
  return n.toFixed(2);
}

function formatChange(row: PlatformIntegralLogRow) {
  const abs = formatStat(row.number);
  return row.pm === 1 ? `+${abs}` : `-${abs}`;
}

function changeClass(row: PlatformIntegralLogRow) {
  // CRMEB：正红负绿
  return row.pm === 1 ? 'change--plus' : 'change--minus';
}

function buildFilterParams(formValues?: Record<string, unknown>) {
  const range = Array.isArray(formValues?.date_range)
    ? formValues.date_range
    : [];
  return {
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults(
  [
    {
      ...LIST_DATE_RANGE_FIELD,
      label: '选择时间',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入用户UID、用户昵称、标题',
      },
      fieldName: 'keyword',
      label: '搜索',
    },
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

const gridOptions: VxeGridProps<PlatformIntegralLogRow> = {
  columns: [
    { field: 'bill_id', title: 'ID', width: 90 },
    {
      field: 'nickname',
      minWidth: 120,
      showOverflow: 'tooltip',
      title: '用户昵称',
      formatter: ({ row }) => row.nickname || `用户 #${row.uid}`,
    },
    {
      field: 'title',
      minWidth: 140,
      showOverflow: 'tooltip',
      title: '积分标题',
    },
    {
      field: 'number',
      slots: { default: 'numberChange' },
      title: '积分变动',
      width: 110,
    },
    {
      field: 'balance',
      title: '当前积分额度',
      width: 120,
      formatter: ({ cellValue }) => formatStat(Number(cellValue)),
    },
    {
      field: 'mark',
      className: 'col--remark',
      minWidth: 180,
      width: 260,
      showOverflow: 'tooltip',
      title: '备注',
      formatter: ({ cellValue }) =>
        cellValue === undefined || cellValue === null || cellValue === ''
          ? '—'
          : String(cellValue),
    },
    {
      field: 'create_time',
      minWidth: 170,
      showOverflow: 'tooltip',
      title: '添加时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const values =
          formValues && Object.keys(formValues).length > 0
            ? formValues
            : lastFormValues.value;
        const filters = buildFilterParams(values);
        const [data, stats] = await Promise.all([
          listPlatformIntegralLogsApi({
            page: page.currentPage,
            limit: page.pageSize,
            ...filters,
          }),
          getPlatformIntegralLogTitleApi(filters).catch(() => emptyTitle()),
        ]);
        title.value = stats || emptyTitle();
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'bill_id' },
  scrollX: { enabled: false },
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
  if (exporting.value) return;
  exporting.value = true;
  try {
    const filters = buildFilterParams(lastFormValues.value);
    const result = await exportPlatformIntegralLogsApi(filters);
    const blob = new Blob([result.content], {
      type: 'text/csv;charset=utf-8',
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name || '积分日志.csv';
    link.click();
    URL.revokeObjectURL(url);
    ElMessage.success(
      `已导出 ${result.row_count} 条${result.truncated ? '（已截断）' : ''}`,
    );
  } catch {
    ElMessage.error('导出失败，请稍后重试');
  } finally {
    exporting.value = false;
  }
}
</script>

<template>
  <Page auto-content-height>
    <div class="integral-log-filter">
      <Form />
    </div>

    <div class="integral-log-summary">
      <div class="integral-summary">
        <div
          v-for="card in summaryCards"
          :key="card.key"
          class="integral-summary__card"
          :class="`integral-summary__card--${card.tone}`"
        >
          <div class="integral-summary__icon">
            <IconifyIcon :icon="card.icon" />
          </div>
          <div class="integral-summary__body">
            <div class="integral-summary__value">{{ card.value }}</div>
            <div class="integral-summary__label">{{ card.label }}</div>
          </div>
        </div>
      </div>
    </div>

    <Grid>
      <template #toolbar-actions>
        <ElButton type="primary" :loading="exporting" @click="exportRows">
          导出
        </ElButton>
      </template>
      <template #numberChange="{ row }">
        <span :class="changeClass(row)">{{ formatChange(row) }}</span>
      </template>
    </Grid>
  </Page>
</template>

<style scoped>
.integral-log-filter {
  padding: 12px 8px 4px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.integral-log-summary {
  padding: 16px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.integral-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  width: 100%;
}

.integral-summary__card {
  display: flex;
  gap: 16px;
  align-items: center;
  min-height: 88px;
  padding: 20px 22px;
  background: hsl(var(--background));
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
  box-shadow: 0 1px 2px rgb(0 0 0 / 3%);
}

.integral-summary__icon {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  color: #fff;
  font-size: 24px;
  border-radius: 50%;
}

.integral-summary__card--blue .integral-summary__icon {
  background: #409eff;
}

.integral-summary__card--orange .integral-summary__icon {
  background: #e6a23c;
}

.integral-summary__card--green .integral-summary__icon {
  background: #67c23a;
}

.integral-summary__card--pink .integral-summary__icon {
  background: #f56c6c;
}

.integral-summary__card--purple .integral-summary__icon {
  background: #9b59b6;
}

.integral-summary__card--sky .integral-summary__icon {
  background: #36cfc9;
}

.integral-summary__value {
  color: var(--el-text-color-primary);
  font-size: 24px;
  font-weight: 600;
  line-height: 1.2;
}

.integral-summary__label {
  margin-top: 6px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.2;
}

.change--plus {
  color: #f56c6c;
  font-weight: 600;
}

.change--minus {
  color: #67c23a;
  font-weight: 600;
}

@media (min-width: 1600px) {
  .integral-summary {
    grid-template-columns: repeat(6, minmax(0, 1fr));
  }
}

@media (max-width: 1100px) {
  .integral-summary {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .integral-summary {
    grid-template-columns: 1fr;
  }
}
</style>
