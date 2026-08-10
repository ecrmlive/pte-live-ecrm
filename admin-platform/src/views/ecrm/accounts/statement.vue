<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import {
  ElButton,
  ElMessage,
  ElSkeleton,
  ElTabPane,
  ElTabs,
} from 'element-plus';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  exportPlatformStatementApi,
  getPlatformStatementDetailApi,
  getPlatformStatementTitleApi,
  listPlatformStatementsApi,
  type PlatformStatementDetail,
  type PlatformStatementRow,
  type PlatformStatementTitle,
  type PlatformStatementType,
} from '#/api/core/platform-statement';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const emptyTitle = (): PlatformStatementTitle => ({
  brokerage_expense: 0,
  coupon_amount: 0,
  merchant_count: 0,
  order_income: 0,
  platform_charge: 0,
  recharge_amount: 0,
  recharge_consume: 0,
  refund_expense: 0,
});

const title = ref<PlatformStatementTitle>(emptyTitle());
const lastFormValues = ref<Record<string, unknown>>({});
const billType = ref<PlatformStatementType>(1);
const canDownload = ref(false);
const detailLoading = ref(false);
const detail = ref<PlatformStatementDetail | null>(null);
const downloading = ref(false);

const summaryCards = computed(() => [
  {
    key: 'order_income',
    label: '订单收入总金额',
    value: formatMoney(title.value.order_income),
    icon: 'lucide:shopping-bag',
    tone: 'blue',
  },
  {
    key: 'refund_expense',
    label: '退款支出金额',
    value: formatMoney(title.value.refund_expense),
    icon: 'lucide:receipt',
    tone: 'orange',
  },
  {
    key: 'brokerage_expense',
    label: '佣金支出金额',
    value: formatMoney(title.value.brokerage_expense),
    icon: 'lucide:handshake',
    tone: 'green',
  },
  {
    key: 'platform_charge',
    label: '平台手续费',
    value: formatMoney(title.value.platform_charge),
    icon: 'lucide:badge-percent',
    tone: 'pink',
  },
  {
    key: 'recharge_amount',
    label: '充值金额',
    value: formatMoney(title.value.recharge_amount),
    icon: 'lucide:wallet',
    tone: 'purple',
  },
  {
    key: 'recharge_consume',
    label: '充值消费金额',
    value: formatMoney(title.value.recharge_consume),
    icon: 'lucide:credit-card',
    tone: 'sky',
  },
  {
    key: 'merchant_count',
    label: '产生交易的商户数',
    value: String(title.value.merchant_count || 0),
    icon: 'lucide:store',
    tone: 'blue',
  },
  {
    key: 'coupon_amount',
    label: '优惠券金额',
    value: formatMoney(title.value.coupon_amount),
    icon: 'lucide:ticket',
    tone: 'orange',
  },
]);

function formatMoney(v?: number) {
  return Number(v || 0).toFixed(2);
}

function buildFilterParams(formValues?: Record<string, unknown>) {
  const range = Array.isArray(formValues?.date_range)
    ? formValues.date_range
    : [];
  return {
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults(
  [
    {
      ...LIST_DATE_RANGE_FIELD,
      componentProps: {
        ...LIST_DATE_RANGE_FIELD.componentProps,
        startPlaceholder: '开始日期',
        endPlaceholder: '结束日期',
      },
      label: '选择时间',
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

const gridOptions: VxeGridProps<PlatformStatementRow> = {
  columns: [
    { type: 'seq', title: '序号', width: 70 },
    {
      field: 'date',
      minWidth: 120,
      title: '日期',
    },
    {
      field: 'income',
      formatter: ({ cellValue }) => formatMoney(Number(cellValue)),
      minWidth: 120,
      title: '账期内收入',
    },
    {
      field: 'expend',
      formatter: ({ cellValue }) => formatMoney(Number(cellValue)),
      minWidth: 120,
      title: '账期内支出',
    },
    {
      field: 'offline',
      formatter: ({ cellValue }) => formatMoney(Number(cellValue)),
      minWidth: 130,
      title: '店铺线下已收',
    },
    {
      field: 'charge',
      formatter: ({ cellValue }) => formatMoney(Number(cellValue)),
      minWidth: 140,
      title: '平台应入账金额',
    },
    platformListActionColumn({ minWidth: 160, width: 180 }),
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
          listPlatformStatementsApi({
            page: page.currentPage,
            limit: page.pageSize,
            type: billType.value,
            ...filters,
          }),
          getPlatformStatementTitleApi(filters).catch(() => emptyTitle()),
        ]);
        title.value = stats || emptyTitle();
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'date' },
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

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  showConfirmButton: false,
  title: '账单详情',
});

async function onBillTypeChange(name: string | number) {
  billType.value = Number(name) === 2 ? 2 : 1;
  await gridApi.reload(lastFormValues.value);
}

async function openDetail(row: PlatformStatementRow) {
  detail.value = null;
  detailLoading.value = true;
  detailDrawerApi.setState({ title: `账单详情 · ${row.date}` });
  detailDrawerApi.open();
  try {
    detail.value = await getPlatformStatementDetailApi({
      date: row.date,
      type: billType.value,
    });
  } catch {
    ElMessage.error('加载账单详情失败');
    detailDrawerApi.close();
  } finally {
    detailLoading.value = false;
  }
}

async function downloadBill(row: PlatformStatementRow) {
  if (!canDownload.value || downloading.value) return;
  downloading.value = true;
  try {
    const result = await exportPlatformStatementApi({
      date: row.date,
      type: billType.value,
    });
    const blob = new Blob([result.content], {
      type: 'text/csv;charset=utf-8',
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = result.file_name || `平台账单_${row.date}.csv`;
    link.click();
    URL.revokeObjectURL(url);
    ElMessage.success('账单已下载');
  } catch {
    ElMessage.error('下载失败，请稍后重试');
  } finally {
    downloading.value = false;
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canDownload.value = codes.includes('accounts.statement.download');
});
</script>

<template>
  <Page auto-content-height>
    <div class="statement-filter">
      <Form />
    </div>

    <div class="statement-summary">
      <div class="statement-summary__grid">
        <div
          v-for="card in summaryCards"
          :key="card.key"
          class="statement-summary__card"
          :class="`statement-summary__card--${card.tone}`"
        >
          <div class="statement-summary__icon">
            <IconifyIcon :icon="card.icon" />
          </div>
          <div class="statement-summary__body">
            <div class="statement-summary__value">{{ card.value }}</div>
            <div class="statement-summary__label">{{ card.label }}</div>
          </div>
        </div>
      </div>
    </div>

    <Grid>
      <template #toolbar-actions>
        <ElTabs
          :model-value="String(billType)"
          class="statement-tabs"
          @tab-change="onBillTypeChange"
        >
          <ElTabPane label="日账单" name="1" />
          <ElTabPane label="月账单" name="2" />
        </ElTabs>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton
          v-if="canDownload"
          link
          type="primary"
          :loading="downloading"
          @click="downloadBill(row)"
        >
          下载账单
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <ElSkeleton :loading="detailLoading" animated :rows="8">
        <div v-if="detail" class="statement-detail">
          <div class="statement-detail__period">账期：{{ detail.date }}</div>

          <section
            v-for="block in [detail.income, detail.bill, detail.expend, detail.charge]"
            :key="block.title"
            class="statement-detail__block"
          >
            <div class="statement-detail__head">
              <div class="statement-detail__title">{{ block.title }}</div>
              <div class="statement-detail__amount">
                ¥{{ formatMoney(block.number) }}
                <span v-if="block.count" class="statement-detail__count">
                  （{{ block.count }}）
                </span>
              </div>
            </div>
            <div v-if="block.data?.length" class="statement-detail__lines">
              <div
                v-for="(line, idx) in block.data"
                :key="`${block.title}-${idx}`"
                class="statement-detail__line"
              >
                <span>{{ line.label }}</span>
                <span>{{ line.amount }}</span>
                <span>{{ line.count }}</span>
              </div>
            </div>
          </section>
        </div>
      </ElSkeleton>
    </DetailDrawer>
  </Page>
</template>

<style scoped>
.statement-filter {
  padding: 12px 8px 4px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.statement-summary {
  padding: 16px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.statement-summary__grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  width: 100%;
}

.statement-summary__card {
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

.statement-summary__icon {
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

.statement-summary__card--blue .statement-summary__icon {
  background: #409eff;
}

.statement-summary__card--orange .statement-summary__icon {
  background: #e6a23c;
}

.statement-summary__card--green .statement-summary__icon {
  background: #67c23a;
}

.statement-summary__card--pink .statement-summary__icon {
  background: #f56c6c;
}

.statement-summary__card--purple .statement-summary__icon {
  background: #9b59b6;
}

.statement-summary__card--sky .statement-summary__icon {
  background: #36cfc9;
}

.statement-summary__value {
  color: var(--el-text-color-primary);
  font-size: 24px;
  font-weight: 600;
  line-height: 1.2;
}

.statement-summary__label {
  margin-top: 6px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.2;
}

.statement-tabs {
  margin-bottom: -8px;
}

.statement-tabs :deep(.el-tabs__header) {
  margin: 0;
}

.statement-tabs :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.statement-detail__period {
  margin-bottom: 16px;
  color: var(--el-text-color-regular);
  font-size: 14px;
}

.statement-detail__block {
  padding: 16px;
  margin-bottom: 12px;
  background: hsl(var(--background));
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
}

.statement-detail__head {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
  align-items: baseline;
  justify-content: space-between;
}

.statement-detail__title {
  color: var(--el-text-color-primary);
  font-size: 15px;
  font-weight: 600;
}

.statement-detail__amount {
  color: var(--el-color-primary);
  font-size: 18px;
  font-weight: 600;
}

.statement-detail__count {
  color: var(--el-text-color-secondary);
  font-size: 13px;
  font-weight: 400;
}

.statement-detail__lines {
  display: grid;
  gap: 8px;
  margin-top: 12px;
}

.statement-detail__line {
  display: grid;
  grid-template-columns: 1.2fr 1fr 0.6fr;
  gap: 8px;
  color: var(--el-text-color-regular);
  font-size: 13px;
}

@media (max-width: 1200px) {
  .statement-summary__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .statement-summary__grid {
    grid-template-columns: 1fr;
  }
}
</style>
