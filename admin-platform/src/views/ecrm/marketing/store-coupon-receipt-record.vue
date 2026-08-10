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
import { fetchPlatformMerchants } from '#/api/core/ecrm';
import {
  listCouponReceiptRecords,
  type CouponReceiptRecord,
} from '#/api/core/platform-coupon-command';
import { formatShanghaiDateTime } from '#/utils/date-time';
import { listFormOptionsDefaults } from '#/utils/list-form-defaults';

const route = useRoute();
const canRead = ref(false);
const merchantOptions = ref<{ label: string; value: number }[]>([]);

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

function formatCouponTime(value?: string | null) {
  if (!value) return '—';
  const raw = String(value);
  if (
    raw.startsWith('0000-00-00') ||
    raw.startsWith('0001-01-01') ||
    raw === 'Invalid Date'
  ) {
    return '—';
  }
  return formatShanghaiDateTime(raw);
}

function buildParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const statusRaw = formValues?.status;
  const merRaw = formValues?.mer_id;
  const couponIdRaw = String(formValues?.coupon_id ?? '').trim();
  const allowed = Object.keys(statusLabel);
  return {
    page: page.currentPage,
    limit: page.pageSize,
    coupon_scope: 'store' as const,
    status:
      typeof statusRaw === 'string' && allowed.includes(statusRaw)
        ? (statusRaw as CouponReceiptRecord['status'])
        : undefined,
    mer_id:
      merRaw === 0 || merRaw === undefined || merRaw === null || merRaw === ''
        ? undefined
        : Number(merRaw),
    recipient: String(formValues?.recipient ?? '').trim() || undefined,
    coupon_name: String(formValues?.coupon_name ?? '').trim() || undefined,
    coupon_id: couponIdRaw ? Number(couponIdRaw) : undefined,
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
      placeholder: '请输入优惠券名称',
    },
    fieldName: 'coupon_name',
    label: '优惠券',
  },
]);

const gridOptions: VxeGridProps<CouponReceiptRecord> = {
  columns: [
    { field: 'id', title: 'ID', width: 80 },
    {
      field: 'coupon_name',
      minWidth: 140,
      showOverflow: false,
      title: '优惠券名称',
    },
    {
      field: 'recipient',
      minWidth: 120,
      showOverflow: false,
      title: '领取人',
      formatter: ({ row }) => row.recipient || `用户#${row.user_id}`,
    },
    {
      field: 'store_name',
      minWidth: 120,
      showOverflow: false,
      title: '店铺名称',
      formatter: ({ cellValue, row }) =>
        cellValue || (row.store_id ? `店铺#${row.store_id}` : '—'),
    },
    {
      field: 'coupon_price',
      title: '面值',
      width: 80,
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(0),
    },
    {
      field: 'use_min_price',
      title: '最低消费额',
      width: 100,
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(0),
    },
    {
      field: 'use_start_time',
      minWidth: 170,
      title: '开始使用时间',
      formatter: ({ cellValue }) => formatCouponTime(cellValue),
    },
    {
      field: 'use_end_time',
      minWidth: 170,
      title: '结束使用时间',
      formatter: ({ cellValue }) => formatCouponTime(cellValue),
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

async function loadMerchants() {
  const data = await fetchPlatformMerchants({ page: 1, limit: 200 });
  merchantOptions.value = (data.list || []).map((row) => ({
    label: row.mer_name || `店铺#${row.mer_id}`,
    value: row.mer_id,
  }));
  await gridApi.formApi?.updateSchema?.([
    {
      fieldName: 'mer_id',
      componentProps: {
        clearable: true,
        filterable: true,
        options: merchantOptions.value,
        placeholder: '请选择',
      },
    },
  ]);
}

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canRead.value =
    profile.roles.includes('platform') &&
    (codes.includes('marketing.store_coupon.user') ||
      codes.includes('marketing.store_coupon.list') ||
      codes.includes('marketing.coupon.record.read'));
  await loadMerchants().catch(() => undefined);

  const couponId = String(route.query.coupon_id ?? '').trim();
  const status = String(route.query.status ?? '').trim();
  const merId = String(route.query.mer_id ?? '').trim();
  const allowed = Object.keys(statusLabel);
  const values: Record<string, unknown> = {};
  if (couponId) values.coupon_id = couponId;
  if (status && allowed.includes(status)) values.status = status;
  if (merId) values.mer_id = Number(merId);
  if (Object.keys(values).length) {
    await gridApi.formApi?.setValues(values);
  }
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      title="当前账号没有商户优惠券领取记录查看权限"
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
  color: hsl(142 70% 35%);
}

.avail-icon.is-no {
  color: hsl(var(--destructive));
}
</style>
