<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElMessage,
  ElMessageBox,
  ElTag,
} from 'element-plus';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  againPlatformOrderProfitsharingApi,
  exportPlatformOrderProfitsharingsApi,
  getPlatformOrderProfitsharingApi,
  listPlatformOrderProfitsharingsApi,
  type PlatformOrderProfitsharingQuery,
  type PlatformOrderProfitsharingRow,
} from '#/api/core/platform-order-profitsharing';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const lastFormValues = ref<Record<string, unknown>>({});
const canManage = ref(false);
const canExport = ref(false);
const exporting = ref(false);
const detailLoading = ref(false);
const current = ref<PlatformOrderProfitsharingRow | null>(null);

const statusOptions = [
  { label: '待分账', value: 0 },
  { label: '已分账', value: 1 },
  { label: '分账中', value: 2 },
  { label: '已退款', value: -1 },
  { label: '分账失败', value: -2 },
];

const typeOptions = [
  { label: '订单支付', value: 'order' },
  { label: '尾款支付', value: 'presell' },
];

const storeOptions = [
  { label: 'CRM Live服饰旗舰店', value: 1 },
  { label: 'CRM Live居家优选店', value: 2 },
  { label: 'CRM Live数码生活店', value: 3 },
];

function formatMoney(v?: number) {
  return Number(v || 0).toFixed(2);
}

function statusInfo(status: number) {
  switch (status) {
    case 2:
      return { label: '分账中', type: 'warning' as const };
    case 1:
      return { label: '已分账', type: 'success' as const };
    case 0:
      return { label: '待分账', type: 'info' as const };
    case -1:
      return { label: '已退款', type: 'danger' as const };
    case -2:
      return { label: '分账失败', type: 'danger' as const };
    default:
      return { label: '未知', type: 'info' as const };
  }
}

function typeLabel(type: string) {
  return type === 'presell' ? '尾款支付' : '订单支付';
}

function buildFilterParams(
  formValues?: Record<string, unknown>,
): PlatformOrderProfitsharingQuery {
  const createRange = Array.isArray(formValues?.date_range)
    ? formValues.date_range
    : [];
  const profitRange = Array.isArray(formValues?.profit_date_range)
    ? formValues.profit_date_range
    : [];
  const statusRaw = formValues?.status;
  const typeRaw = String(formValues?.type ?? '').trim();
  const merRaw = formValues?.mer_id;
  const allowedStatus = [0, 1, 2, -1, -2];
  return {
    date_from: createRange[0] as string | undefined,
    date_to: createRange[1] as string | undefined,
    profit_date_from: profitRange[0] as string | undefined,
    profit_date_to: profitRange[1] as string | undefined,
    status: allowedStatus.includes(Number(statusRaw))
      ? Number(statusRaw)
      : undefined,
    type: typeRaw || undefined,
    mer_id:
      merRaw === 0 || merRaw
        ? Number(merRaw) || undefined
        : undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults(
  [
    {
      ...LIST_DATE_RANGE_FIELD,
      fieldName: 'date_range',
      label: '创建时间',
    },
    {
      ...LIST_DATE_RANGE_FIELD,
      fieldName: 'profit_date_range',
      label: '分账时间',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: statusOptions,
        placeholder: '请选择',
      },
      fieldName: 'status',
      label: '状态',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: typeOptions,
        placeholder: '请选择',
      },
      fieldName: 'type',
      label: '账单类型',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: storeOptions,
        placeholder: '请选择',
      },
      fieldName: 'mer_id',
      label: '店铺名称',
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

const gridOptions: VxeGridProps<PlatformOrderProfitsharingRow> = {
  columns: [
    {
      field: 'profitsharing_id',
      minWidth: 100,
      title: '分账ID',
    },
    {
      field: 'order_sn',
      minWidth: 180,
      showOverflow: 'tooltip',
      title: '订单编号',
    },
    {
      field: 'mer_name',
      minWidth: 140,
      showOverflow: 'tooltip',
      title: '店铺名称',
    },
    {
      field: 'profitsharing_price',
      formatter: ({ cellValue }) => formatMoney(Number(cellValue)),
      minWidth: 110,
      title: '订单金额',
    },
    {
      field: 'type',
      formatter: ({ cellValue }) => typeLabel(String(cellValue || '')),
      minWidth: 110,
      title: '账单类型',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      minWidth: 110,
      title: '状态',
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 170,
      title: '创建时间',
    },
    platformListActionColumn({ minWidth: 140, width: 180 }),
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
        const data = await listPlatformOrderProfitsharingsApi({
          page: page.currentPage,
          limit: page.pageSize,
          ...filters,
        });
        return {
          items: data.list || [],
          total: data.total || 0,
        };
      },
    },
  },
  toolbarConfig: {
    custom: true,
    export: false,
    refresh: true,
    search: false,
    zoom: true,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({
  gridOptions,
});

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  title: '分账详情',
  footer: false,
});

async function reloadGrid() {
  await gridApi.reload(lastFormValues.value);
}

async function openDetail(row: PlatformOrderProfitsharingRow) {
  detailLoading.value = true;
  current.value = row;
  detailDrawerApi.open();
  try {
    current.value = await getPlatformOrderProfitsharingApi(row.profitsharing_id);
  } catch {
    // requestClient 已提示
  } finally {
    detailLoading.value = false;
  }
}

async function retryAgain(row: PlatformOrderProfitsharingRow) {
  if (!canManage.value || row.status !== -2) return;
  try {
    await ElMessageBox.confirm(
      `确认对分账 #${row.profitsharing_id} 重新发起分账？`,
      '重新分账',
      { type: 'warning' },
    );
  } catch {
    return;
  }
  try {
    await againPlatformOrderProfitsharingApi(row.profitsharing_id);
    ElMessage.success('分账成功');
    await reloadGrid();
  } catch {
    // requestClient 已提示
  }
}

async function exportRows() {
  if (!canExport.value || exporting.value) return;
  exporting.value = true;
  try {
    const filters = buildFilterParams(lastFormValues.value);
    const result = await exportPlatformOrderProfitsharingsApi(filters);
    const blob = new Blob([result.content], {
      type: 'text/csv;charset=utf-8',
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name || '分账管理.csv';
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

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canManage.value = codes.includes('accounts.profitsharing.manage');
  canExport.value = codes.includes('accounts.profitsharing.export');
});
</script>

<template>
  <Page auto-content-height>
    <div class="profitsharing-filter">
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
      <template #status="{ row }">
        <div>
          <ElTag :type="statusInfo(row.status).type" size="small">
            {{ statusInfo(row.status).label }}
          </ElTag>
          <div
            v-if="row.status === -2 && row.error_msg"
            class="profitsharing-error"
          >
            {{ row.error_msg }}
          </div>
        </div>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton
          v-if="canManage && row.status === -2"
          link
          type="primary"
          @click="retryAgain(row)"
        >
          重新分账
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <div v-loading="detailLoading" class="profitsharing-detail">
        <ElDescriptions v-if="current" :column="2" border>
          <ElDescriptionsItem label="分账ID">
            {{ current.profitsharing_id }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="分账单号">
            {{ current.profitsharing_sn || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="订单编号">
            {{ current.order_sn || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="店铺名称">
            {{ current.mer_name }}（{{ current.mer_id }}）
          </ElDescriptionsItem>
          <ElDescriptionsItem label="订单金额">
            ¥{{ formatMoney(current.profitsharing_price) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="平台手续费">
            ¥{{ formatMoney(current.profitsharing_mer_price) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="已退款金额">
            ¥{{ formatMoney(current.profitsharing_refund) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="账单类型">
            {{ typeLabel(current.type) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="状态">
            <ElTag :type="statusInfo(current.status).type" size="small">
              {{ statusInfo(current.status).label }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="支付交易号">
            {{ current.transaction_id || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="创建时间">
            {{ formatShanghaiDateTime(current.create_time) || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="分账时间">
            {{ formatShanghaiDateTime(current.profitsharing_time) || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem v-if="current.error_msg" :span="2" label="失败原因">
            {{ current.error_msg }}
          </ElDescriptionsItem>
        </ElDescriptions>
        <div
          v-if="current && canManage && current.status === -2"
          class="profitsharing-detail__actions"
        >
          <ElButton type="primary" @click="retryAgain(current)">
            重新分账
          </ElButton>
        </div>
      </div>
    </DetailDrawer>
  </Page>
</template>

<style scoped>
.profitsharing-filter {
  margin-bottom: 12px;
}

.profitsharing-error {
  color: var(--el-color-danger);
  font-size: 12px;
  line-height: 1.4;
  margin-top: 4px;
  max-width: 160px;
}

.profitsharing-detail {
  min-height: 120px;
  padding: 4px 0 16px;
}

.profitsharing-detail__actions {
  margin-top: 16px;
}
</style>
