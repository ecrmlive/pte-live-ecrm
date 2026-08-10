<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';
import { useRoute } from 'vue-router';

import { Page } from '@vben/common-ui';
import { ElAlert, ElTag } from 'element-plus';
import { Check, Close } from '@element-plus/icons-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listCouponReceiptRecords,
  type CouponReceiptRecord,
} from '#/api/core/platform-coupon-command';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const route = useRoute();
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
  const userIdRaw = String(formValues?.user_id ?? '').trim();
  const couponIdRaw = String(formValues?.coupon_id ?? '').trim();
  const statusRaw = formValues?.status;
  const sourceRaw = String(formValues?.source ?? '').trim();
  const scopeRaw = String(formValues?.coupon_scope ?? '').trim();
  const recipient = String(formValues?.recipient ?? '').trim();
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
    source: sourceRaw || undefined,
    coupon_scope:
      scopeRaw === 'platform' || scopeRaw === 'store' ? scopeRaw : undefined,
    recipient: recipient || undefined,
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: Object.entries(statusLabel).map(([value, label]) => ({
        label,
        value,
      })),
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '使用状态',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入领取人',
    },
    fieldName: 'recipient',
    label: '领取人',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '请输入优惠券ID',
    },
    fieldName: 'coupon_id',
    label: '优惠券',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '领取', value: 'receive' },
        { label: '新人', value: 'onboarding' },
        { label: '后台发放', value: 'platform_manual' },
      ],
      placeholder: '全部',
    },
    fieldName: 'source',
    label: '获取方式',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '平台通用券', value: 'platform' },
        { label: '店铺券', value: 'store' },
      ],
      placeholder: '全部',
    },
    fieldName: 'coupon_scope',
    label: '优惠类型',
  },
]);

const gridOptions: VxeGridProps<CouponReceiptRecord> = {
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    {
      field: 'coupon_name',
      minWidth: 160,
      showOverflow: false,
      title: '优惠券名称',
    },
    {
      field: 'recipient',
      minWidth: 140,
      showOverflow: false,
      title: '领取人',
      formatter: ({ row }) => row.recipient || `用户#${row.user_id}`,
    },
    {
      field: 'coupon_type_name',
      title: '优惠券类型',
      width: 120,
      formatter: ({ cellValue, row }) =>
        cellValue || (row.store_id ? '店铺券' : '平台通用券'),
    },
    {
      field: 'coupon_price',
      title: '面值',
      width: 90,
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(0),
    },
    {
      field: 'use_min_price',
      title: '最低消费额',
      width: 110,
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(0),
    },
    {
      field: 'use_start_time',
      minWidth: 170,
      title: '开始使用时间',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '不限时',
    },
    {
      field: 'use_end_time',
      minWidth: 170,
      title: '结束使用时间',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '不限时',
    },
    {
      field: 'source_name',
      title: '获取方式',
      width: 100,
      formatter: ({ cellValue, row }) => cellValue || row.source || '—',
    },
    {
      field: 'available',
      slots: { default: 'available' },
      title: '是否可用',
      width: 90,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 100,
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const data = await listCouponReceiptRecords(
          buildParams(page, formValues),
        );
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
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canRead.value =
    profile.roles.includes('platform') &&
    (codes.includes('marketing.coupon.record.read') ||
      codes.includes('marketing.store_coupon.user') ||
      codes.includes('marketing.store_coupon.list'));
  const couponId = String(route.query.coupon_id ?? '').trim();
  const status = String(route.query.status ?? '').trim();
  const scope = String(route.query.coupon_scope ?? '').trim();
  const allowed = Object.keys(statusLabel);
  if (
    couponId ||
    (status && allowed.includes(status)) ||
    scope === 'platform' ||
    scope === 'store'
  ) {
    await gridApi.formApi?.setValues({
      ...(couponId ? { coupon_id: couponId } : {}),
      ...(status && allowed.includes(status) ? { status } : {}),
      ...(scope === 'platform' || scope === 'store'
        ? { coupon_scope: scope }
        : {}),
    });
  }
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      title="当前账号没有优惠券领取记录查看权限"
      type="warning"
      :closable="false"
      class="mb-3"
    />
    <Grid v-else>
      <template #available="{ row }">
        <span
          class="avail-icon"
          :class="row.available ? 'is-yes' : 'is-no'"
        >
          <Check v-if="row.available" />
          <Close v-else />
        </span>
      </template>
      <template #status="{ row }">
        <ElTag :type="statusType[row.status]">
          {{ statusLabel[row.status] }}
        </ElTag>
      </template>
    </Grid>
  </Page>
</template>

<style scoped>
.avail-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  font-size: 14px;
}

.avail-icon.is-yes {
  color: hsl(var(--primary));
}

.avail-icon.is-no {
  color: hsl(var(--muted-foreground));
}
</style>
