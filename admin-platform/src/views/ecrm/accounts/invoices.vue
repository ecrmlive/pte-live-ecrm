<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElSkeleton } from 'element-plus';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  fetchPlatformMerchants,
  type PlatformMerchantRow,
} from '#/api/core/ecrm';
import {
  getPlatformInvoiceApi,
  listPlatformInvoices,
  type PlatformInvoice,
  type PlatformInvoiceQuery,
} from '#/api/core/platform-invoice';
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

const ORDER_TYPE_TABS = [
  { key: '', name: '全部' },
  { key: '1', name: '待付款' },
  { key: '2', name: '待发货' },
  { key: '3', name: '待收货' },
  { key: '4', name: '待评价' },
  { key: '5', name: '交易完成' },
  { key: '6', name: '已退款' },
  { key: '7', name: '已删除' },
] as const;

const orderType = ref('');
const lastFormValues = ref<Record<string, unknown>>({});
const canRead = ref(false);
const storeOptions = ref<{ label: string; value: number }[]>([]);
const detailLoading = ref(false);
const current = ref<PlatformInvoice | null>(null);

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
): Omit<PlatformInvoiceQuery, 'page' | 'limit'> {
  const range = Array.isArray(formValues?.date_range)
    ? formValues.date_range
    : [];
  const merId = Number(formValues?.mer_id || 0);
  const userSearch = parseUserSearch(formValues);
  return {
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    mer_id: merId > 0 ? merId : undefined,
    status: String(formValues?.status ?? '').trim() || undefined,
    order_type: orderType.value || undefined,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
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
        options: [],
        placeholder: '请选择',
      },
      fieldName: 'mer_id',
      label: '店铺名称',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [
          { label: '已开票', value: '1' },
          { label: '未开票', value: '0' },
        ],
        placeholder: '请选择',
      },
      fieldName: 'status',
      label: '开票状态',
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
        placeholder: '请输入发票联系人',
      },
      fieldName: 'keyword',
      label: '发票搜索',
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

const gridOptions: VxeGridProps<PlatformInvoice> = {
  columns: [
    {
      field: 'order_sn',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 180,
      showOverflow: 'tooltip',
      title: '订单号',
    },
    {
      field: 'mer_name',
      formatter: ({ cellValue, row }) =>
        dash(cellValue || row.merchant_name || row.store_name),
      minWidth: 140,
      showOverflow: 'tooltip',
      title: '店铺名称',
    },
    {
      field: 'nickname',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 120,
      showOverflow: 'tooltip',
      title: '用户昵称',
    },
    {
      field: 'pay_price',
      formatter: ({ cellValue }) => formatMoney(Number(cellValue)),
      minWidth: 110,
      title: '订单金额',
    },
    {
      field: 'order_status_label',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 100,
      title: '订单状态',
    },
    {
      field: 'invoice_amount',
      formatter: ({ cellValue }) => formatMoney(Number(cellValue)),
      minWidth: 110,
      title: '发票金额',
    },
    {
      field: 'receipt_sn',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 160,
      showOverflow: 'tooltip',
      title: '发票单号',
    },
    {
      field: 'invoice_type_label',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 100,
      title: '发票类型',
    },
    {
      field: 'title_type_label',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 90,
      title: '抬头类型',
    },
    {
      field: 'contact_name',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 100,
      showOverflow: 'tooltip',
      title: '发票联系人',
    },
    {
      field: 'contact_info',
      formatter: ({ cellValue }) => dash(cellValue),
      minWidth: 160,
      showOverflow: 'tooltip',
      title: '发票联系信息',
    },
    {
      field: 'create_time',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '-',
      minWidth: 170,
      title: '下单时间',
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
        const data = await listPlatformInvoices({
          page: page.currentPage,
          limit: page.pageSize,
          ...buildFilterParams(values),
        });
        return { items: data.list || [], total: data.total || 0 };
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

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  showConfirmButton: false,
  title: '发票详情',
});

async function setOrderType(key: string) {
  orderType.value = key;
  await gridApi.reload(lastFormValues.value);
}

function detailTitleOf(row?: PlatformInvoice | null) {
  if (!row) return '发票详情';
  if (row.detail_title) return row.detail_title;
  if (Number(row.invoice_type) === 2) return '企业专用纸质发票';
  return row.profile_type === 'enterprise'
    ? '企业电子普通发票'
    : '个人电子普通发票';
}

function detailText(v?: string | number | null) {
  if (v === undefined || v === null) return '';
  return String(v);
}

function detailAmount(v?: number) {
  const n = Number(v || 0);
  return n > 0 ? n.toFixed(2) : '';
}

async function openDetail(row: PlatformInvoice) {
  current.value = row;
  detailLoading.value = true;
  detailDrawerApi.setState({ title: detailTitleOf(row) });
  detailDrawerApi.open();
  try {
    current.value = await getPlatformInvoiceApi(row.id);
    detailDrawerApi.setState({ title: detailTitleOf(current.value) });
  } finally {
    detailLoading.value = false;
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi().catch(() => [] as string[]);
  canRead.value = codes.includes('accounts.invoice.read');
  const merchants = await fetchPlatformMerchants({
    page: 1,
    limit: 200,
    status: 1,
  }).catch(() => ({ list: [] as PlatformMerchantRow[] }));
  storeOptions.value = (merchants.list || []).map((m) => ({
    label: m.merchant_name || `店铺${m.merchant_id}`,
    value: Number(m.merchant_id),
  }));
  formApi.updateSchema([
    {
      fieldName: 'mer_id',
      componentProps: {
        clearable: true,
        filterable: true,
        options: storeOptions.value,
        placeholder: '请选择',
      },
    },
  ]);
  if (canRead.value) {
    await gridApi.reload();
  }
});
</script>

<template>
  <Page auto-content-height>
    <div class="invoice-filter">
      <div class="invoice-status-tabs" role="tablist">
        <span class="invoice-status-tabs__label">订单状态：</span>
        <button
          v-for="tab in ORDER_TYPE_TABS"
          :key="tab.key || 'all'"
          type="button"
          role="tab"
          class="invoice-status-tabs__item"
          :aria-selected="orderType === tab.key"
          :class="{ 'is-active': orderType === tab.key }"
          @click="setOrderType(tab.key)"
        >
          {{ tab.name }}
        </button>
      </div>
      <Form />
    </div>

    <Grid>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <ElSkeleton :loading="detailLoading" animated :rows="8">
        <div v-if="current" class="invoice-detail">
          <section class="invoice-detail__section">
            <div class="invoice-detail__title">发票详情</div>
            <div class="invoice-detail__stack">
              <div class="invoice-kv">
                <span class="invoice-kv__label">发票申请单号：</span>
                <span>{{ detailText(current.receipt_sn) }}</span>
              </div>
            </div>
          </section>

          <section class="invoice-detail__section">
            <div class="invoice-detail__title">发票信息</div>
            <div class="invoice-detail__grid">
              <div class="invoice-kv">
                <span class="invoice-kv__label">发票抬头：</span>
                <span>{{ detailText(current.title) }}</span>
              </div>
              <div class="invoice-kv">
                <span class="invoice-kv__label">发票类型：</span>
                <span>{{ detailText(current.invoice_type_label) }}</span>
              </div>
              <div class="invoice-kv">
                <span class="invoice-kv__label">发票抬头类型：</span>
                <span>{{ detailText(current.title_type_label) }}</span>
              </div>
              <div class="invoice-kv">
                <span class="invoice-kv__label">发票金额：</span>
                <span>{{ detailAmount(current.invoice_amount) }}</span>
              </div>
              <div
                v-if="current.profile_type === 'enterprise'"
                class="invoice-kv"
              >
                <span class="invoice-kv__label">企业税号：</span>
                <span>{{ detailText(current.tax_no) }}</span>
              </div>
            </div>
          </section>

          <section class="invoice-detail__section">
            <div class="invoice-detail__title">联系信息</div>
            <div class="invoice-detail__stack">
              <div class="invoice-kv">
                <span class="invoice-kv__label">联系邮箱：</span>
                <span>{{ detailText(current.email) }}</span>
              </div>
              <div class="invoice-kv">
                <span class="invoice-kv__label">开票状态：</span>
                <span>{{ detailText(current.status_label) }}</span>
              </div>
              <div class="invoice-kv">
                <span class="invoice-kv__label">发票备注：</span>
                <span>{{ detailText(current.mark) }}</span>
              </div>
            </div>
          </section>
        </div>
      </ElSkeleton>
    </DetailDrawer>
  </Page>
</template>

<style scoped>
.invoice-filter {
  padding: 12px 8px 4px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.invoice-status-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
  padding: 0 4px 4px;
}

.invoice-status-tabs__label {
  color: var(--el-text-color-regular);
  font-size: 14px;
}

.invoice-status-tabs__item {
  padding: 4px 12px;
  color: var(--el-text-color-regular);
  font-size: 13px;
  line-height: 1.4;
  background: var(--el-fill-color-light);
  border: 1px solid transparent;
  border-radius: 4px;
  cursor: pointer;
}

.invoice-status-tabs__item.is-active {
  color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
  border-color: var(--el-color-primary-light-5);
}

.invoice-detail__section + .invoice-detail__section {
  margin-top: 28px;
}

.invoice-detail__title {
  position: relative;
  margin-bottom: 18px;
  padding-left: 10px;
  color: var(--el-text-color-primary);
  font-size: 15px;
  font-weight: 600;
}

.invoice-detail__title::before {
  position: absolute;
  top: 2px;
  left: 0;
  width: 3px;
  height: 16px;
  background: var(--el-color-primary);
  border-radius: 2px;
  content: '';
}

.invoice-detail__grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px 48px;
}

.invoice-detail__stack {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.invoice-kv {
  display: flex;
  gap: 4px;
  align-items: flex-start;
  min-height: 22px;
  color: var(--el-text-color-primary);
  font-size: 14px;
  line-height: 1.6;
  word-break: break-all;
}

.invoice-kv__label {
  flex-shrink: 0;
  color: var(--el-text-color-secondary);
}

@media (max-width: 768px) {
  .invoice-detail__grid {
    grid-template-columns: 1fr;
  }
}
</style>
