<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import { ElImage } from 'element-plus';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  getPlatformVerifySummaryApi,
  listPlatformVerifyRecordsApi,
  type PlatformVerifyRecord,
  type PlatformVerifySummary,
} from '#/api/core/platform-trade';
import { platformListPagerConfig } from '#/constants/platform-list-grid';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const emptySummary = (): PlatformVerifySummary => ({
  paid_count: 0,
  pay_amount: 0,
  refund_amount: 0,
  wechat_amount: 0,
  balance_amount: 0,
  alipay_amount: 0,
});

const summary = ref<PlatformVerifySummary>(emptySummary());
const lastFormValues = ref<Record<string, unknown>>({});

const summaryCards = computed(() => [
  {
    key: 'paid_count',
    label: '已支付订单数量',
    value: String(summary.value.paid_count || 0),
    icon: 'ant-design:shopping-outlined',
    tone: 'blue',
  },
  {
    key: 'pay_amount',
    label: '实际支付金额',
    value: Number(summary.value.pay_amount || 0).toFixed(2),
    icon: 'ant-design:file-text-outlined',
    tone: 'orange',
  },
  {
    key: 'refund_amount',
    label: '已退款金额',
    value: Number(summary.value.refund_amount || 0).toFixed(2),
    icon: 'mdi:briefcase-outline',
    tone: 'green',
  },
  {
    key: 'wechat_amount',
    label: '微信支付金额',
    value: Number(summary.value.wechat_amount || 0).toFixed(2),
    icon: 'mdi:purse-outline',
    tone: 'pink',
  },
  {
    key: 'balance_amount',
    label: '余额支付金额',
    value: Number(summary.value.balance_amount || 0).toFixed(2),
    icon: 'ant-design:wallet-outlined',
    tone: 'purple',
  },
  {
    key: 'alipay_amount',
    label: '支付宝支付金额',
    value: Number(summary.value.alipay_amount || 0).toFixed(2),
    icon: 'mdi:bag-suitcase-outline',
    tone: 'sky',
  },
]);

function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || v === '') return '--';
  return String(v);
}

function moneyPlain(v?: number) {
  return Number(v || 0).toFixed(2);
}

function buildFilterParams(formValues?: Record<string, unknown>) {
  const range = Array.isArray(formValues?.date_range)
    ? formValues.date_range
    : [];
  return {
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    order_keyword: String(formValues?.order_keyword ?? '').trim() || undefined,
    pay_type:
      formValues?.pay_type === 0 ||
      formValues?.pay_type === 1 ||
      formValues?.pay_type === 2 ||
      formValues?.pay_type === 7
        ? Number(formValues.pay_type)
        : undefined,
    is_trader:
      formValues?.is_trader === 0 || formValues?.is_trader === 1
        ? Number(formValues.is_trader)
        : undefined,
    user_keyword: String(formValues?.user_keyword ?? '').trim() || undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults(
  [
    {
      ...LIST_DATE_RANGE_FIELD,
      label: '核销时间',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入订单号/收货人/联系方式',
      },
      fieldName: 'order_keyword',
      label: '订单号',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '余额支付', value: 0 },
          { label: '微信支付', value: 1 },
          { label: '支付宝支付', value: 2 },
          { label: '模拟支付', value: 7 },
        ],
        placeholder: '请选择',
      },
      fieldName: 'pay_type',
      label: '支付方式',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '自营', value: 1 },
          { label: '非自营', value: 0 },
        ],
        placeholder: '请选择',
      },
      fieldName: 'is_trader',
      label: '店铺类别',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入用户信息/联系电话',
      },
      fieldName: 'user_keyword',
      label: '用户信息',
    },
  ],
  {
    // 独立 Form（不经 useVbenVxeGrid）：对齐金标准/Grid 默认 chrome；勿设 actionButtonsReverse（重置 → 搜索）
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

const gridOptions: VxeGridProps<PlatformVerifyRecord> = {
  columns: [
    {
      type: 'expand',
      width: 46,
      align: 'center',
      slots: { content: 'expandContent' },
    },
    {
      field: 'order_sn',
      minWidth: 180,
      showOverflow: 'tooltip',
      title: '订单编号',
    },
    {
      field: 'order_type_label',
      minWidth: 100,
      title: '订单类型',
      formatter: ({ row }) => row.order_type_label || '核销订单',
    },
    {
      field: 'store_name',
      minWidth: 120,
      showOverflow: 'tooltip',
      title: '店铺名称',
      formatter: ({ row }) =>
        row.store_name || row.mer_name || `店铺 #${row.mer_id}`,
    },
    {
      field: 'store_category_name',
      minWidth: 88,
      title: '店铺类别',
      formatter: ({ row }) => dash(row.store_category_name),
    },
    {
      field: 'real_name',
      minWidth: 96,
      showOverflow: 'tooltip',
      title: '收货人',
      formatter: ({ row }) => dash(row.real_name),
    },
    {
      // 弹性主列：吸收容器剩余宽度，避免末列后空白 gutter
      field: 'product',
      minWidth: 220,
      align: 'left',
      headerAlign: 'center',
      showOverflow: false,
      slots: { default: 'productInfo' },
      title: '商品信息',
    },
    {
      field: 'pay_price',
      title: '实际支付',
      width: 96,
      formatter: ({ cellValue }) => moneyPlain(Number(cellValue)),
    },
    {
      field: 'pay_type_label',
      title: '支付方式',
      width: 100,
      formatter: ({ row }) => dash(row.pay_type_label),
    },
    {
      field: 'verifier_name',
      minWidth: 96,
      showOverflow: 'tooltip',
      title: '核销员',
      formatter: ({ row }) => dash(row.verifier_name),
    },
    {
      field: 'verify_status_label',
      slots: { default: 'verifyStatus' },
      title: '核销状态',
      width: 88,
    },
    {
      field: 'verify_time',
      width: 168,
      showOverflow: 'tooltip',
      title: '核销时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  expandConfig: { trigger: 'default' },
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
          listPlatformVerifyRecordsApi({
            page: page.currentPage,
            limit: page.pageSize,
            ...filters,
          }),
          getPlatformVerifySummaryApi(filters).catch(() => emptySummary()),
        ]);
        summary.value = stats || emptySummary();
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'order_id' },
  // 本页无 fixed 操作列：关闭平台默认 scrollX(gt:0)，否则列宽锁在 minWidth 之和、
  // 宽屏时末列（核销时间）右侧会留出大块空白 gutter。窄屏由 platform-list-page 外层 overflow 横滑。
  scrollX: { enabled: false },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
    enabled: false,
  },
};

// 与金标准相同：不传 height / class；筛选与统计在 Grid 外独立带，Grid 内仅表体+pager
const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });
</script>

<template>
  <Page auto-content-height>
    <!-- 筛选：独立白底区（不进 Grid formOptions，保证统计卡可夹在中间） -->
    <div class="cancellation-filter">
      <Form />
    </div>

    <!-- 统计：独立区域，禁止塞进 vxe toolbar -->
    <div class="cancellation-summary">
      <div class="verify-summary">
        <div
          v-for="card in summaryCards"
          :key="card.key"
          class="verify-summary__card"
          :class="`verify-summary__card--${card.tone}`"
        >
          <div class="verify-summary__icon">
            <IconifyIcon :icon="card.icon" />
          </div>
          <div class="verify-summary__body">
            <div class="verify-summary__value">{{ card.value }}</div>
            <div class="verify-summary__label">{{ card.label }}</div>
          </div>
        </div>
      </div>
    </div>

    <Grid>
      <template #productInfo="{ row }">
        <div v-if="row.product" class="order-product">
          <ElImage
            class="order-product__thumb"
            :src="resolveCosMediaUrl(row.product.product_image || '')"
            fit="cover"
          />
          <div class="order-product__text">
            <div class="order-product__name">
              {{ row.product.product_info || '—' }}
            </div>
            <div v-if="row.product.product_sku" class="order-product__sku">
              {{ row.product.product_sku }}
            </div>
            <div class="order-product__price">
              ¥ {{ moneyPlain(row.product.product_price) }} x
              {{ row.product.product_num || 0 }}
            </div>
          </div>
        </div>
        <span v-else>—</span>
      </template>

      <template #verifyStatus="{ row }">
        <span class="verify-status">
          {{ row.verify_status_label || '已核销' }}
        </span>
      </template>

      <template #expandContent="{ row }">
        <div class="order-expand">
          <div class="order-expand__item">
            <span class="label">商品总价：</span>
            <span class="value">{{ moneyPlain(row.total_price) }}</span>
          </div>
          <div class="order-expand__item">
            <span class="label">用户备注：</span>
            <span class="value">{{ dash(row.user_remark) }}</span>
          </div>
          <div class="order-expand__item">
            <span class="label">店铺备注：</span>
            <span class="value">{{ dash(row.merchant_remark) }}</span>
          </div>
        </div>
      </template>
    </Grid>
  </Page>
</template>

<style scoped>
.cancellation-filter {
  padding: 12px 8px 4px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.cancellation-summary {
  padding: 16px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.verify-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  width: 100%;
}

.verify-summary__card {
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

.verify-summary__icon {
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

.verify-summary__card--blue .verify-summary__icon {
  background: #409eff;
}

.verify-summary__card--orange .verify-summary__icon {
  background: #e6a23c;
}

.verify-summary__card--green .verify-summary__icon {
  background: #67c23a;
}

.verify-summary__card--pink .verify-summary__icon {
  background: #f56c6c;
}

.verify-summary__card--purple .verify-summary__icon {
  background: #9b59b6;
}

.verify-summary__card--sky .verify-summary__icon {
  background: #409eff;
}

.verify-summary__value {
  color: var(--el-text-color-primary);
  font-size: 24px;
  font-weight: 600;
  line-height: 1.2;
}

.verify-summary__label {
  margin-top: 6px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.2;
}

.order-product {
  display: flex;
  gap: 8px;
  align-items: flex-start;
  padding: 2px 0;
  text-align: left;
}

.order-product__thumb {
  flex: 0 0 48px;
  width: 48px;
  height: 48px;
  border-radius: 4px;
  background: var(--el-fill-color-light);
}

.order-product__text {
  min-width: 0;
  line-height: 1.45;
  word-break: break-all;
}

.order-product__name {
  color: var(--el-text-color-primary);
  font-size: 13px;
}

.order-product__sku,
.order-product__price {
  margin-top: 2px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.verify-status {
  color: var(--el-text-color-primary);
  font-size: 13px;
}

.order-expand {
  display: flex;
  flex-wrap: wrap;
  gap: 16px 48px;
  padding: 10px 16px 10px 46px;
  background: #fff;
}

.order-expand__item {
  display: flex;
  gap: 4px;
  align-items: baseline;
  min-width: 160px;
  font-size: 13px;
  line-height: 1.6;
}

.order-expand__item .label {
  flex-shrink: 0;
  color: var(--el-text-color-secondary);
}

.order-expand__item .value {
  color: var(--el-text-color-primary);
}

@media (min-width: 1600px) {
  .verify-summary {
    grid-template-columns: repeat(6, minmax(0, 1fr));
  }
}

@media (max-width: 1100px) {
  .verify-summary {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 720px) {
  .verify-summary {
    grid-template-columns: 1fr;
  }
}
</style>
