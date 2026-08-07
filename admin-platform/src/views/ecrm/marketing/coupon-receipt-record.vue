<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listCouponReceiptRecords,
  type CouponReceiptRecord,
} from '#/api/core/platform-coupon-command';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canRead = ref(false);

const statusLabel: Record<CouponReceiptRecord['status'], string> = {
  unused: '未使用',
  locked: '已锁定',
  used: '已使用',
  expired: '已过期',
};

const statusType: Record<
  CouponReceiptRecord['status'],
  'danger' | 'info' | 'success' | 'warning'
> = {
  unused: 'success',
  locked: 'warning',
  used: 'info',
  expired: 'danger',
};

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const userIdRaw = String(formValues?.user_id ?? '').trim();
  const couponIdRaw = String(formValues?.coupon_id ?? '').trim();
  const statusRaw = formValues?.status;
  const allowed = Object.keys(statusLabel);
  return {
    page: page.currentPage,
    limit: page.pageSize,
    user_id: userIdRaw ? Number(userIdRaw) : undefined,
    coupon_id: couponIdRaw ? Number(couponIdRaw) : undefined,
    status:
      typeof statusRaw === 'string' && allowed.includes(statusRaw)
        ? (statusRaw as CouponReceiptRecord['status'])
        : undefined,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '用户 ID' },
    fieldName: 'user_id',
    label: '用户 ID',
  },
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '优惠券 ID' },
    fieldName: 'coupon_id',
    label: '优惠券 ID',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: Object.entries(statusLabel).map(([value, label]) => ({ label, value })),
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '券状态',
  },
]);

const gridOptions: VxeGridProps<CouponReceiptRecord> = {
  columns: [
    { field: 'id', title: '用户券 ID', width: 110 },
    { field: 'user_id', title: '用户 ID', width: 100 },
    {
      field: 'coupon_name',
      minWidth: 220,
      showOverflow: false,
      title: '优惠券',
      formatter: ({ row }) => `${row.coupon_name}（#${row.coupon_id}）`,
    },
    {
      field: 'store_id',
      title: '归属',
      width: 110,
      formatter: ({ cellValue }) => (cellValue ? `店铺 #${cellValue}` : '平台券'),
    },
    { field: 'source', title: '来源', width: 150 },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 100,
    },
    {
      field: 'used_order_id',
      title: '关联订单 ID',
      width: 140,
      formatter: ({ cellValue }) => cellValue || '—',
    },
    {
      field: 'obtained_at',
      minWidth: 180,
      title: '领取时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const data = await listCouponReceiptRecords(buildParams(page, formValues));
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

onMounted(async () => {
  const [profile, codes] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canRead.value = profile.roles.includes('platform') && codes.includes('marketing.coupon.record.read');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page
    auto-content-height
    description="只读监管用户已领取优惠券及其真实状态；页面不返回账号、手机号或地址，不提供人工核销、解锁、删除或状态修复。"
    title="优惠券领取记录"
  >
    <ElAlert
      v-if="!canRead"
      title="当前账号没有优惠券领取记录查看权限"
      type="warning"
      :closable="false"
    />
    <template v-else>
      <Grid>
        <template #status="{ row }">
          <ElTag :type="statusType[row.status]">{{ statusLabel[row.status] }}</ElTag>
        </template>
      </Grid>
      <ElAlert
        class="mt-4"
        type="info"
        :closable="false"
        title="“已锁定、已使用、已过期”等状态由下单、支付、取消和券规则状态机写入；平台本页仅监管事实，不能绕过订单状态机修改。"
      />
    </template>
  </Page>
</template>
