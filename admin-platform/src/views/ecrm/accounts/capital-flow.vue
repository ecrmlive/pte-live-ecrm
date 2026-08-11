<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElMessage, ElSkeleton } from 'element-plus';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  exportPlatformCapitalFlowsApi,
  getPlatformCapitalFlowApi,
  listPlatformCapitalFlowsApi,
  type PlatformCapitalFlowQuery,
  type PlatformCapitalFlowRow,
} from '#/api/core/platform-capital-flow';
import {
  listUserSearchFormField,
  parseUserSearch,
} from '#/components/ecrm/user-search-field';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';
import OrderDetailDrawer from '../order/components/OrderDetailDrawer.vue';

/** CRMEB：订单支付/分佣类流水详情打开订单详情 */
const ORDER_DETAIL_TYPES = new Set([
  'order',
  'brokerage_one',
  'brokerage_two',
  'refund_brokerage_one',
  'refund_brokerage_two',
  'refund_order',
  'order_platform_coupon',
  'order_svip_coupon',
]);

const lastFormValues = ref<Record<string, unknown>>({});
const canRead = ref(false);
const canExport = ref(false);
const exporting = ref(false);
const detailLoading = ref(false);
const current = ref<PlatformCapitalFlowRow | null>(null);
const orderDetailDrawerRef = ref<InstanceType<typeof OrderDetailDrawer>>();

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
): PlatformCapitalFlowQuery {
  const range = Array.isArray(formValues?.date_range)
    ? formValues.date_range
    : [];
  const userSearch = parseUserSearch(formValues);
  return {
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    order_sn: String(formValues?.order_sn ?? '').trim() || undefined,
    pay_type: String(formValues?.pay_type ?? '').trim() || undefined,
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
    listUserSearchFormField({
      options: [
        { label: '昵称', value: 'nickname' },
        { label: '用户ID', value: 'uid' },
        { label: '手机号', value: 'phone' },
      ],
      typeWidth: '96px',
    }),
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入订单号',
      },
      fieldName: 'order_sn',
      label: '订单号',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '余额', value: '0' },
          { label: '微信', value: '1' },
          { label: '支付宝', value: '2' },
          { label: '线下支付', value: '3' },
        ],
        placeholder: '请选择',
      },
      fieldName: 'pay_type',
      label: '支付方式',
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

const gridOptions: VxeGridProps<PlatformCapitalFlowRow> = {
  columns: [
    {
      field: 'order_sn',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 180,
      showOverflow: 'tooltip',
      title: '订单号',
    },
    {
      field: 'financial_record_sn',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 180,
      showOverflow: 'tooltip',
      title: '交易流水号',
    },
    {
      field: 'transaction_id',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 160,
      showOverflow: 'tooltip',
      title: '第三方交易单号',
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '-',
      minWidth: 170,
      title: '交易时间',
    },
    {
      field: 'user_info',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 140,
      showOverflow: 'tooltip',
      title: '对方信息',
    },
    {
      field: 'financial_type_name',
      formatter: ({ cellValue, row }) =>
        dash(cellValue || row.financial_type),
      minWidth: 120,
      title: '交易类型',
    },
    {
      field: 'pay_type_name',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 110,
      title: '支付方式',
    },
    {
      field: 'signed_number',
      formatter: ({ cellValue, row }) => {
        const n =
          cellValue === undefined || cellValue === null
            ? row.financial_pm === 1
              ? Number(row.number || 0)
              : -Number(row.number || 0)
            : Number(cellValue);
        return formatMoney(n);
      },
      minWidth: 130,
      title: '收支金额（元）',
    },
    platformListActionColumn({ width: 90 }),
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
        const data = await listPlatformCapitalFlowsApi({
          page: page.currentPage,
          limit: page.pageSize,
          ...buildFilterParams(values),
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'financial_record_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [FlowDrawer, flowDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  showConfirmButton: false,
  title: '流水详情',
});

function shouldOpenOrderDetail(row: PlatformCapitalFlowRow) {
  const type = String(row.financial_type || '');
  const orderId = Number(row.order_id || 0);
  return orderId > 0 && ORDER_DETAIL_TYPES.has(type);
}

async function openFlowDrawer(row: PlatformCapitalFlowRow) {
  current.value = row;
  detailLoading.value = true;
  flowDrawerApi.setState({
    title: `流水详情 · ${row.financial_record_sn || row.financial_record_id}`,
  });
  flowDrawerApi.open();
  try {
    current.value = await getPlatformCapitalFlowApi(row.financial_record_id);
  } finally {
    detailLoading.value = false;
  }
}

async function openDetail(row: PlatformCapitalFlowRow) {
  if (shouldOpenOrderDetail(row)) {
    const ok = await orderDetailDrawerRef.value?.open(Number(row.order_id));
    if (ok) return;
  }
  await openFlowDrawer(row);
}

async function exportRows() {
  if (!canExport.value || exporting.value) return;
  exporting.value = true;
  try {
    const result = await exportPlatformCapitalFlowsApi(
      buildFilterParams(lastFormValues.value),
    );
    const blob = new Blob([result.content || ''], {
      type: 'text/csv;charset=utf-8',
    });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = result.file_name || '资金流水.csv';
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
  canRead.value = codes.includes('accounts.capital_flow.read');
  canExport.value = codes.includes('accounts.capital_flow.export');
  if (canRead.value) {
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
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
      </template>
    </Grid>

    <OrderDetailDrawer ref="orderDetailDrawerRef" />

    <FlowDrawer>
      <ElSkeleton :loading="detailLoading" animated :rows="8">
        <div v-if="current" class="flow-detail">
          <section class="flow-detail__section">
            <div class="flow-detail__title">交易信息</div>
            <div class="flow-detail__grid">
              <div class="flow-kv">
                <span class="flow-kv__label">订单号：</span>
                <span>{{ dash(current.order_sn) }}</span>
              </div>
              <div class="flow-kv">
                <span class="flow-kv__label">交易流水号：</span>
                <span>{{ dash(current.financial_record_sn) }}</span>
              </div>
              <div class="flow-kv">
                <span class="flow-kv__label">第三方交易单号：</span>
                <span>{{ dash(current.transaction_id) }}</span>
              </div>
              <div class="flow-kv">
                <span class="flow-kv__label">交易时间：</span>
                <span>
                  {{
                    current.create_time
                      ? formatShanghaiDateTime(current.create_time)
                      : '-'
                  }}
                </span>
              </div>
              <div class="flow-kv">
                <span class="flow-kv__label">对方信息：</span>
                <span>{{ dash(current.user_info) }}</span>
              </div>
              <div class="flow-kv">
                <span class="flow-kv__label">用户ID：</span>
                <span>{{ dash(current.user_id) }}</span>
              </div>
              <div class="flow-kv">
                <span class="flow-kv__label">交易类型：</span>
                <span>
                  {{
                    dash(current.financial_type_name || current.financial_type)
                  }}
                </span>
              </div>
              <div class="flow-kv">
                <span class="flow-kv__label">支付方式：</span>
                <span>{{ dash(current.pay_type_name) }}</span>
              </div>
              <div class="flow-kv">
                <span class="flow-kv__label">收支金额：</span>
                <span>¥{{ formatMoney(current.signed_number) }}</span>
              </div>
              <div class="flow-kv">
                <span class="flow-kv__label">商户ID：</span>
                <span>{{ dash(current.mer_id) }}</span>
              </div>
            </div>
          </section>
        </div>
      </ElSkeleton>
    </FlowDrawer>
  </Page>
</template>

<style scoped>
.capital-filter {
  padding: 12px 8px 4px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.flow-detail__section {
  padding: 4px 0;
}

.flow-detail__title {
  position: relative;
  margin-bottom: 16px;
  padding-left: 10px;
  color: var(--el-text-color-primary);
  font-size: 15px;
  font-weight: 600;
}

.flow-detail__title::before {
  position: absolute;
  top: 2px;
  left: 0;
  width: 3px;
  height: 16px;
  background: var(--el-color-primary);
  border-radius: 2px;
  content: '';
}

.flow-detail__grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px 24px;
}

.flow-kv {
  display: flex;
  gap: 4px;
  align-items: flex-start;
  color: var(--el-text-color-primary);
  font-size: 14px;
  line-height: 1.5;
  word-break: break-all;
}

.flow-kv__label {
  flex-shrink: 0;
  color: var(--el-text-color-secondary);
}

@media (max-width: 768px) {
  .flow-detail__grid {
    grid-template-columns: 1fr;
  }
}
</style>
