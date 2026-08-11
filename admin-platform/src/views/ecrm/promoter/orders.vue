<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElButton, ElImage, ElTag } from 'element-plus';

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
import OrderDetailDrawer from '../order/components/OrderDetailDrawer.vue';

const STATUS_TABS: Array<{ key: keyof PlatformOrderTabCounts; name: string }> =
  [
    { key: 'all', name: '全部' },
    { key: 'unpaid', name: '待付款' },
    { key: 'unshipped', name: '待发货' },
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
  const activityType = values.activity_type;
  const productType = values.product_type;
  const userSearch = parseUserSearch(values);
  return {
    page: page.currentPage,
    limit: page.pageSize,
    tab_status: tabStatus.value,
    is_spread: 1 as const,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    activity_type:
      activityType === 0 || activityType ? Number(activityType) : undefined,
    mer_id: merId > 0 ? merId : undefined,
    mer_category_id: merCategoryId > 0 ? merCategoryId : undefined,
    product_name: String(values.product_name ?? '').trim() || undefined,
    delivery_type: String(values.delivery_type ?? '').trim() || undefined,
    product_type:
      productType === 0 || productType ? Number(productType) : undefined,
    keyword: String(values.keyword ?? '').trim() || undefined,
    spread_keyword: String(values.spread_keyword ?? '').trim() || undefined,
    top_spread_keyword:
      String(values.top_spread_keyword ?? '').trim() || undefined,
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
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入订单号/收货人/联系方式',
    },
    fieldName: 'keyword',
    label: '关键字',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入推广人' },
    fieldName: 'spread_keyword',
    label: '推广人',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '请输入上级推广' },
    fieldName: 'top_spread_keyword',
    label: '上级推广',
  },
  listUserSearchFormField({ label: '用户信息' }),
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
      formatter: ({ row }) => row.order_type_label || '普通订单',
      minWidth: 100,
      title: '订单类型',
    },
    {
      field: 'real_name',
      formatter: ({ row }) => dash(row.real_name),
      minWidth: 110,
      showOverflow: false,
      title: '收货人/订购人',
    },
    {
      field: 'store_name',
      formatter: ({ row }) =>
        row.store_name || row.mer_name || `店铺 #${row.mer_id}`,
      minWidth: 130,
      showOverflow: false,
      title: '店铺名称',
    },
    {
      field: 'product',
      minWidth: 220,
      showOverflow: false,
      slots: { default: 'productInfo' },
      title: '商品信息',
    },
    {
      field: 'unit_pay',
      slots: { default: 'unitPay' },
      title: '实际支付',
      width: 120,
    },
    {
      field: 'pay_price',
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
      title: '支付金额',
      width: 100,
    },
    {
      field: 'commission_total',
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
      title: '订单佣金',
      width: 100,
    },
    {
      field: 'pay_type_label',
      formatter: ({ row }) => dash(row.pay_type_label),
      title: '支付类型',
      width: 96,
    },
    {
      field: 'pay_status_label',
      formatter: ({ row }) => dash(row.pay_status_label),
      title: '支付状态',
      width: 96,
    },
    {
      field: 'status_label',
      slots: { default: 'status' },
      title: '订单',
      width: 100,
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

      <template #unitPay="{ row }">
        <span v-if="row.product">
          {{ money(row.product.product_price) }} x
          {{ row.product.product_num || 1 }}
        </span>
        <span v-else>{{ money(row.pay_price) }}</span>
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
            <span class="label">推广人</span>
            <span class="value">{{ dash(row.spread_name) }}</span>
          </div>
          <div class="order-expand__item">
            <span class="label">上级推广</span>
            <span class="value">{{ dash(row.top_spread_name) }}</span>
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
