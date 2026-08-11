<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import {
  ElButton,
  ElImage,
  ElMessage,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  fetchMerchantCategories,
  fetchPlatformMerchants,
  type MerchantCategoryRow,
} from '#/api/core/ecrm';
import {
  getPlatformOrderTabCountsApi,
  listPlatformOrdersApi,
  type PlatformOrder,
  type PlatformOrderTabCounts,
} from '#/api/core/platform-trade';
import {
  listPrefixedKeywordFormField,
  listUserSearchFormField,
  parseUserSearch,
} from '#/components/ecrm/user-search-field';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

import OrderDetailDrawer from './components/OrderDetailDrawer.vue';

const STATUS_TABS: Array<{ key: keyof PlatformOrderTabCounts; name: string }> =
  [
    { key: 'all', name: '全部' },
    { key: 'unpaid', name: '待付款' },
    { key: 'unshipped', name: '待发货/核销' },
    { key: 'unreceived', name: '待收货' },
    { key: 'unevaluated', name: '待评价' },
    { key: 'completed', name: '交易完成' },
    { key: 'refunded', name: '已退款' },
    { key: 'deleted', name: '已删除' },
  ];

const tabStatus = ref<keyof PlatformOrderTabCounts>('all');
const tabCounts = ref<PlatformOrderTabCounts>({
  all: 0,
  unpaid: 0,
  unshipped: 0,
  unreceived: 0,
  unevaluated: 0,
  completed: 0,
  refunded: 0,
  deleted: 0,
});
const lastFormValues = ref<Record<string, unknown>>({});
const storeOptions = ref<{ label: string; value: number }[]>([]);
const categoryOptions = ref<{ label: string; value: number }[]>([]);
const orderDetailDrawerRef = ref<InstanceType<typeof OrderDetailDrawer>>();

const ORDER_SEARCH_OPTIONS = [
  { label: '订单号', value: 'order_sn' },
  { label: '收货人', value: 'real_name' },
  { label: '收货电话', value: 'phone' },
];

function money(v?: number) {
  return `¥${Number(v || 0).toFixed(2)}`;
}

function dash(v?: string | number | null) {
  if (v === 0) return '0';
  if (v === undefined || v === null || v === '') return '-';
  return String(v);
}

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const values = formValues || {};
  lastFormValues.value = values;
  const range = Array.isArray(values.date_range) ? values.date_range : [];
  const merId = Number(values.mer_id || 0);
  const merCategoryId = Number(values.mer_category_id || 0);
  const merTypeId = Number(values.mer_type_id || 0);
  const activityType = values.activity_type;
  const payType = values.pay_type;
  const productType = values.product_type;
  const orderRaw = values.order_search;
  const orderType =
    orderRaw && typeof orderRaw === 'object'
      ? String((orderRaw as { type?: string }).type ?? 'order_sn')
      : 'order_sn';
  const orderKeyword =
    orderRaw && typeof orderRaw === 'object'
      ? String((orderRaw as { keyword?: string }).keyword ?? '').trim()
      : '';
  const userSearch = parseUserSearch(values);
  return {
    page: page.currentPage,
    limit: page.pageSize,
    tab_status: tabStatus.value,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    activity_type:
      activityType === 0 || activityType ? Number(activityType) : undefined,
    mer_id: merId > 0 ? merId : undefined,
    mer_category_id: merCategoryId > 0 ? merCategoryId : undefined,
    mer_type_id: merTypeId > 0 ? merTypeId : undefined,
    product_name: String(values.product_name ?? '').trim() || undefined,
    pay_type:
      payType === 0 || payType === 1 || payType === 2 || payType === 7
        ? Number(payType)
        : undefined,
    delivery_type: String(values.delivery_type ?? '').trim() || undefined,
    product_type:
      productType === 0 || productType ? Number(productType) : undefined,
    order_search_type: orderType.trim() || 'order_sn',
    order_search_keyword: orderKeyword || undefined,
    user_search_type: userSearch.type || 'nickname',
    user_search_keyword: userSearch.keyword || undefined,
  };
}

async function loadTabCounts(formValues?: Record<string, unknown>) {
  const params = buildListParams(
    { currentPage: 1, pageSize: 10 },
    formValues ?? lastFormValues.value,
  );
  const {
    page: _p,
    limit: _l,
    tab_status: _t,
    ...filters
  } = params;
  try {
    tabCounts.value = await getPlatformOrderTabCountsApi(filters);
  } catch {
    tabCounts.value = {
      all: 0,
      unpaid: 0,
      unshipped: 0,
      unreceived: 0,
      unevaluated: 0,
      completed: 0,
      refunded: 0,
      deleted: 0,
    };
  }
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  { ...LIST_DATE_RANGE_FIELD, label: '时间选择' },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '普通', value: 0 },
        { label: '秒杀', value: 1 },
        { label: '预售', value: 2 },
        { label: '拼团', value: 3 },
        { label: '助力', value: 4 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'activity_type',
    label: '活动类型',
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
      filterable: true,
      options: [],
      placeholder: '请选择',
    },
    fieldName: 'mer_category_id',
    label: '店铺类别',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入商品名称' },
    fieldName: 'product_name',
    label: '商品名称',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '余额', value: 0 },
        { label: '微信', value: 1 },
        { label: '支付宝', value: 2 },
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
        { label: '快递发货', value: 'express' },
        { label: '到店自提', value: 'pickup' },
        { label: '同城配送', value: 'city' },
        { label: '虚拟发货', value: 'service' },
      ],
      placeholder: '请选择',
    },
    fieldName: 'delivery_type',
    label: '发货方式',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '普通商品', value: 0 },
        { label: '虚拟', value: 1 },
        { label: '云盘', value: 2 },
        { label: '卡密', value: 3 },
        { label: '预约', value: 4 },
        { label: '年/次卡', value: 5 },
      ],
      placeholder: '请选择',
    },
    fieldName: 'product_type',
    label: '商品类型',
  },
  listPrefixedKeywordFormField({
    fieldName: 'order_search',
    label: '订单搜索',
    defaultType: 'order_sn',
    options: ORDER_SEARCH_OPTIONS,
    typeWidth: '110px',
  }),
  listUserSearchFormField({ typeWidth: '96px' }),
]);

const gridOptions: VxeGridProps<PlatformOrder> = {
  columns: [
    {
      type: 'expand',
      width: 48,
      slots: { content: 'expandContent' },
    },
    {
      field: 'order_sn',
      minWidth: 190,
      showOverflow: false,
      slots: { default: 'orderSn' },
      title: '订单编号',
    },
    {
      field: 'nickname',
      minWidth: 140,
      showOverflow: false,
      slots: { default: 'userInfo' },
      title: '用户信息',
    },
    {
      field: 'order_type_label',
      minWidth: 100,
      title: '订单类型',
      formatter: ({ row }) => row.order_type_label || '普通订单',
    },
    {
      field: 'real_name',
      minWidth: 110,
      showOverflow: false,
      title: '收货人/订购人',
      formatter: ({ row }) => dash(row.real_name),
    },
    {
      field: 'store_name',
      minWidth: 130,
      showOverflow: false,
      title: '店铺名称',
      formatter: ({ row }) =>
        row.store_name || row.mer_name || `店铺 #${row.mer_id}`,
    },
    {
      field: 'product',
      minWidth: 220,
      showOverflow: false,
      slots: { default: 'productInfo' },
      title: '商品信息',
    },
    {
      field: 'pay_price',
      title: '实际支付',
      width: 100,
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
    },
    {
      field: 'pay_type_label',
      title: '支付方式',
      width: 96,
      formatter: ({ row }) => dash(row.pay_type_label),
    },
    {
      field: 'pay_status_label',
      title: '支付状态',
      width: 96,
      formatter: ({ row }) => dash(row.pay_status_label),
    },
    {
      field: 'status_label',
      slots: { default: 'status' },
      title: '订单状态',
      width: 110,
    },
    platformListActionColumn({ width: 88 }),
  ],
  expandConfig: { trigger: 'default' },
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const params = buildListParams(page, formValues);
        const [data] = await Promise.all([
          listPlatformOrdersApi(params),
          loadTabCounts(formValues),
        ]);
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'order_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

function setStatusTab(key: keyof PlatformOrderTabCounts) {
  if (tabStatus.value === key) return;
  tabStatus.value = key;
  gridApi.reload();
}

function exportStub() {
  ElMessage.info('导出列表接口尚未接入，已保留按钮占位');
}

function openDetail(id: number) {
  orderDetailDrawerRef.value?.open(id);
}

onMounted(async () => {
  const [stores, categories] = await Promise.all([
    fetchPlatformMerchants({ page: 1, limit: 200, status: 1 }).catch(() => ({
      list: [],
    })),
    fetchMerchantCategories().catch(() => ({
      list: [] as MerchantCategoryRow[],
    })),
  ]);
  storeOptions.value = (stores.list || []).map((m: any) => ({
    label: m.mer_name || m.merchant_name || `店铺 #${m.mer_id || m.id}`,
    value: Number(m.mer_id || m.id),
  }));
  categoryOptions.value = (categories.list || []).map((c) => ({
    label: c.category_name,
    value: c.merchant_category_id,
  }));
  gridApi.formApi?.updateSchema([
    {
      fieldName: 'mer_id',
      componentProps: {
        clearable: true,
        filterable: true,
        options: storeOptions.value,
        placeholder: '请选择',
      },
    },
    {
      fieldName: 'mer_category_id',
      componentProps: {
        clearable: true,
        filterable: true,
        options: categoryOptions.value,
        placeholder: '请选择',
      },
    },
  ]);
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-actions>
        <div class="order-toolbar">
          <div class="order-status-tabs" role="tablist">
            <button
              v-for="tab in STATUS_TABS"
              :key="tab.key"
              type="button"
              role="tab"
              class="order-status-tabs__item"
              :aria-selected="tabStatus === tab.key"
              :class="{ 'is-active': tabStatus === tab.key }"
              @click="setStatusTab(tab.key)"
            >
              {{ tab.name }}({{ tabCounts[tab.key] || 0 }})
            </button>
          </div>
          <div class="order-toolbar__actions">
            <ElButton type="primary" @click="exportStub">导出列表</ElButton>
          </div>
        </div>
      </template>

      <template #orderSn="{ row }">
        <div class="order-sn">
          <div>{{ row.order_sn }}</div>
          <div v-if="row.user_deleted" class="order-sn__tip">用户已删除</div>
        </div>
      </template>

      <template #userInfo="{ row }">
        <span class="order-user-link">
          {{ dash(row.nickname) }}/{{ row.uid || '—' }}
        </span>
      </template>

      <template #productInfo="{ row }">
        <div v-if="row.product" class="order-product">
          <ElImage
            class="order-product__thumb"
            :src="resolveCosMediaUrl(row.product.product_image || '')"
            fit="cover"
          />
          <div class="order-product__text">
            {{ row.product.product_info || '—' }}
          </div>
        </div>
        <span v-else>—</span>
      </template>

      <template #status="{ row }">
        <ElTag type="info" effect="plain">
          {{ row.status_label || '未知' }}
        </ElTag>
      </template>

      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row.order_id)">
          详情
        </ElButton>
      </template>

      <template #expandContent="{ row }">
        <div class="order-expand">
          <div class="order-expand__item">
            <span class="label">商品总价</span>
            <span class="value">{{ money(row.total_price) }}</span>
          </div>
          <div class="order-expand__item">
            <span class="label">下单时间</span>
            <span class="value">{{
              formatShanghaiDateTime(row.create_time)
            }}</span>
          </div>
          <div class="order-expand__item">
            <span class="label">用户备注</span>
            <span class="value">{{ dash(row.user_remark) }}</span>
          </div>
          <div class="order-expand__item">
            <span class="label">店铺备注</span>
            <span class="value">{{ dash(row.merchant_remark) }}</span>
          </div>
          <div class="order-expand__item">
            <span class="label">总单号</span>
            <span class="value">{{ dash(row.group_order_sn) }}</span>
          </div>
        </div>
      </template>
    </Grid>

    <OrderDetailDrawer ref="orderDetailDrawerRef" />
  </Page>
</template>

<style scoped>
.order-toolbar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
}

.order-status-tabs {
  display: flex;
  flex-wrap: wrap;
  gap: 4px 18px;
}

.order-status-tabs__item {
  padding: 6px 0;
  color: var(--el-text-color-regular);
  cursor: pointer;
  background: transparent;
  border: 0;
  border-bottom: 2px solid transparent;
}

.order-status-tabs__item:hover {
  color: var(--el-color-primary);
}

.order-status-tabs__item.is-active {
  color: var(--el-color-primary);
  border-bottom-color: var(--el-color-primary);
}

.order-toolbar__actions {
  display: flex;
  gap: 8px;
}

.order-sn__tip {
  margin-top: 2px;
  font-size: 12px;
  color: var(--el-color-danger);
}

.order-user-link {
  color: var(--el-color-primary);
}

.order-product {
  display: flex;
  gap: 8px;
  align-items: flex-start;
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

.order-expand {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 28px;
  padding: 8px 12px;
}

.order-expand__item {
  display: flex;
  gap: 8px;
  min-width: 180px;
}

.order-expand__item .label {
  color: var(--el-text-color-secondary);
}
</style>
