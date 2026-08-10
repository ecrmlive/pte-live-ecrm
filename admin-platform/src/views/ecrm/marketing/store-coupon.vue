<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import { ElButton, ElTag } from 'element-plus';
import { CloseBold } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  getStoreCouponDetailApi,
  listStoreCouponsApi,
  type StoreCouponDetail,
  type StoreCouponListItem,
} from '#/api/core/platform-promotion';
import CouponUserIssueModal, {
  type CouponIssueMode,
} from '#/components/marketing/coupon-user-issue-modal.vue';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const detail = ref<StoreCouponDetail | null>(null);
const detailLoading = ref(false);
const issueOpen = ref(false);
const issueCouponId = ref(0);
const issueMode = ref<CouponIssueMode>('used');

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const statusRaw = formValues?.status;
  const traderRaw = formValues?.is_trader;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    status:
      statusRaw === 0 || statusRaw === 1 ? Number(statusRaw) : undefined,
    is_trader:
      traderRaw === 0 || traderRaw === 1 ? Number(traderRaw) : undefined,
  };
}

function formatTime(value?: string | null) {
  return value ? formatShanghaiDateTime(value) : '—';
}

function formatMinPrice(value?: number) {
  const n = Number(value || 0);
  return n > 0 ? `最低消费${n}` : '无门槛';
}

function formatPublish(row: StoreCouponListItem) {
  if (row.is_limited === 1) {
    return `发布: ${row.total_count} / 剩余: ${row.remain_count}`;
  }
  return '不限量';
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '自营', value: 1 },
        { label: '非自营', value: 0 },
      ],
      placeholder: '全部',
    },
    fieldName: 'is_trader',
    label: '店铺类别',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '开启', value: 1 },
        { label: '未开启', value: 0 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '状态',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入优惠券名称',
    },
    fieldName: 'keyword',
    label: '关键字',
  },
]);

const gridOptions: VxeGridProps<StoreCouponListItem> = {
  columns: [
    { field: 'coupon_id', title: 'ID', width: 80 },
    { field: 'title', minWidth: 140, showOverflow: false, title: '优惠券名称' },
    {
      field: 'coupon_type_name',
      title: '优惠券类型',
      width: 100,
    },
    {
      field: 'mer_name',
      minWidth: 120,
      showOverflow: false,
      title: '店铺名称',
    },
    {
      field: 'trader_name',
      title: '店铺类别',
      width: 90,
    },
    {
      field: 'claim_text',
      minWidth: 120,
      showOverflow: false,
      title: '领取日期',
    },
    {
      field: 'validity_text',
      title: '使用时间',
      width: 100,
    },
    {
      field: 'total_count',
      minWidth: 140,
      showOverflow: false,
      title: '发布数量',
      formatter: ({ row }) => formatPublish(row),
    },
    {
      field: 'received_total',
      minWidth: 140,
      slots: { default: 'usageCount' },
      title: '使用数量',
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    platformListActionColumn({ width: 200 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await listStoreCouponsApi(buildParams(page, formValues));
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'coupon_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid] = useVbenVxeGrid({ formOptions, gridOptions });

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  footer: false,
  placement: 'right',
});

async function openDetail(row: StoreCouponListItem) {
  detail.value = null;
  detailDrawerApi.setState({ title: '优惠券详情', loading: true }).open();
  detailLoading.value = true;
  try {
    detail.value = await getStoreCouponDetailApi(row.coupon_id);
  } finally {
    detailLoading.value = false;
    detailDrawerApi.setState({ loading: false });
  }
}

function openIssue(row: StoreCouponListItem, mode: CouponIssueMode) {
  issueCouponId.value = row.coupon_id;
  issueMode.value = mode;
  issueOpen.value = true;
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #usageCount="{ row }">
        <div class="usage-count">
          <div>已领取总数: {{ row.received_total ?? 0 }}</div>
          <div>已使用总数: {{ row.used_total ?? 0 }}</div>
        </div>
      </template>
      <template #status="{ row }">
        <span :class="row.status === 1 ? 'status-on' : 'status-off'">
          {{ row.status === 1 ? '开启' : '未开启' }}
        </span>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton link type="primary" @click="openIssue(row, 'receive')">
          领取记录
        </ElButton>
        <ElButton link type="primary" @click="openIssue(row, 'used')">
          使用记录
        </ElButton>
      </template>
    </Grid>

    <CouponUserIssueModal
      v-model:open="issueOpen"
      :coupon-id="issueCouponId"
      coupon-scope="store"
      :mode="issueMode"
    />

    <DetailDrawer>
      <div v-loading="detailLoading" class="store-coupon-detail">
        <template v-if="detail">
          <div class="store-coupon-detail__grid">
            <div class="store-coupon-detail__item">
              <span class="label">优惠券名称</span>
              <span class="value">{{ detail.title || '—' }}</span>
            </div>
            <div class="store-coupon-detail__item">
              <span class="label">优惠券类型</span>
              <span class="value">{{
                detail.coupon_type_name || '店铺券'
              }}</span>
            </div>
            <div class="store-coupon-detail__item">
              <span class="label">店铺名称</span>
              <span class="value">{{ detail.mer_name || '—' }}</span>
            </div>
            <div class="store-coupon-detail__item">
              <span class="label">店铺类别</span>
              <span class="value">{{ detail.trader_name || '—' }}</span>
            </div>
            <div class="store-coupon-detail__item">
              <span class="label">优惠券面值</span>
              <span class="value">{{
                Number(detail.coupon_price || 0).toFixed(2)
              }}</span>
            </div>
            <div class="store-coupon-detail__item">
              <span class="label">使用门槛</span>
              <span class="value">{{
                formatMinPrice(detail.use_min_price)
              }}</span>
            </div>
            <div class="store-coupon-detail__item">
              <span class="label">领取日期</span>
              <span class="value">{{ detail.claim_text || '不限时' }}</span>
            </div>
            <div class="store-coupon-detail__item">
              <span class="label">使用时间</span>
              <span class="value">{{ detail.validity_text || '—' }}</span>
            </div>
            <div class="store-coupon-detail__item">
              <span class="label">是否限量</span>
              <span class="value">
                <ElTag
                  v-if="detail.is_limited === 1"
                  type="success"
                  size="small"
                >
                  是
                </ElTag>
                <span v-else class="limit-no">
                  <CloseBold class="limit-no__icon" />
                </span>
              </span>
            </div>
            <div class="store-coupon-detail__item">
              <span class="label">已领取总数</span>
              <span class="value">{{ detail.received_total ?? 0 }}</span>
            </div>
            <div class="store-coupon-detail__item">
              <span class="label">已使用总数</span>
              <span class="value">{{ detail.used_total ?? 0 }}</span>
            </div>
            <div class="store-coupon-detail__item">
              <span class="label">创建时间</span>
              <span class="value">{{ formatTime(detail.create_time) }}</span>
            </div>
            <div class="store-coupon-detail__item">
              <span class="label">状态</span>
              <span class="value">{{
                detail.status === 1 ? '开启' : '未开启'
              }}</span>
            </div>
          </div>
        </template>
      </div>
    </DetailDrawer>
  </Page>
</template>

<style scoped>
.usage-count {
  font-size: 13px;
  line-height: 1.6;
}

.status-on {
  color: hsl(142 70% 35%);
}

.status-off {
  color: hsl(var(--destructive));
}

.store-coupon-detail {
  min-height: 200px;
  padding: 8px 4px 16px;
}

.store-coupon-detail__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px 32px;
}

.store-coupon-detail__item {
  display: grid;
  grid-template-columns: 120px minmax(0, 1fr);
  gap: 8px;
  align-items: start;
  font-size: 13px;
  line-height: 1.6;
}

.store-coupon-detail__item .label {
  color: hsl(var(--muted-foreground));
}

.store-coupon-detail__item .value {
  color: hsl(var(--foreground));
  word-break: break-all;
}

.limit-no {
  display: inline-flex;
  color: hsl(var(--primary));
}

.limit-no__icon {
  width: 14px;
  height: 14px;
}
</style>
