<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import { ElAlert, ElTag } from 'element-plus';

import { useVbenForm } from '#/adapter/form';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  getSvipOrderSummary,
  listSvipOrders,
  type SvipOrder,
  type SvipOrderSummary,
} from '#/api/core/platform-svip-plan';
import { platformListPagerConfig } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canRead = ref(false);
const emptySummary = (): SvipOrderSummary => ({
  paid_member_count: 0,
  paid_amount: 0,
  expired_member_count: 0,
});
const summary = ref<SvipOrderSummary>(emptySummary());
const lastFormValues = ref<Record<string, unknown>>({});

const summaryCards = computed(() => [
  {
    key: 'paid_member_count',
    label: '累计付费会员人数',
    value: String(summary.value.paid_member_count || 0),
    icon: 'lucide:users',
    tone: 'blue',
  },
  {
    key: 'paid_amount',
    label: '累计支付会员费',
    value: Number(summary.value.paid_amount || 0).toFixed(2),
    icon: 'lucide:circle-dollar-sign',
    tone: 'orange',
  },
  {
    key: 'expired_member_count',
    label: '累计已过期人数',
    value: String(summary.value.expired_member_count || 0),
    icon: 'lucide:history',
    tone: 'pink',
  },
]);

const PAY_TYPE_OPTIONS = [
  { label: '微信', value: 'weixin' },
  { label: '支付宝', value: 'alipay' },
  { label: '小程序', value: 'routine' },
  { label: '平台赠送', value: 'sys' },
  { label: '免费', value: 'free' },
];

function payTypeText(payType?: string) {
  const key = String(payType || '').trim();
  if (!key) return '—';
  if (key === 'wechat' || key === 'weixin') return '微信';
  return (
    {
      alipay: '支付宝',
      routine: '小程序',
      sys: '平台赠送',
      free: '免费',
    }[key] || key
  );
}

function payStatusText(row: SvipOrder) {
  const payType = String(row.pay_type || '').trim();
  if (payType === 'sys' || payType === 'free') return '无需支付';
  return (
    { pending: '未支付', paid: '已支付', closed: '已关闭' }[row.status] ||
    row.status
  );
}

function payStatusType(row: SvipOrder) {
  const payType = String(row.pay_type || '').trim();
  if (payType === 'sys' || payType === 'free') return 'info' as const;
  return (
    { pending: 'warning', paid: 'success', closed: 'info' }[row.status] || 'info'
  ) as 'info' | 'success' | 'warning';
}

function endTimeText(row: SvipOrder) {
  if (row.plan_type === 'lifetime') return '永久';
  return row.end_time ? formatShanghaiDateTime(row.end_time) : '—';
}

function buildFilterParams(formValues?: Record<string, unknown>) {
  const range = Array.isArray(formValues?.date_range)
    ? formValues.date_range
    : [];
  return {
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
    pay_type: String(formValues?.pay_type ?? '').trim() || undefined,
    title: String(formValues?.title ?? '').trim() || undefined,
    nickname: String(formValues?.nickname ?? '').trim() || undefined,
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
        options: PAY_TYPE_OPTIONS,
        placeholder: '请选择',
      },
      fieldName: 'pay_type',
      label: '支付方式',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入会员卡名称',
      },
      fieldName: 'title',
      label: '会员卡名称',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入用户名称搜索',
      },
      fieldName: 'nickname',
      label: '用户名',
    },
  ],
  {
    commonConfig: { componentProps: { class: 'w-full' } },
    wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-4',
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

const gridOptions: VxeGridProps<SvipOrder> = {
  columns: [
    {
      field: 'order_no',
      minWidth: 180,
      showOverflow: 'tooltip',
      title: '订单号',
    },
    {
      field: 'nickname',
      minWidth: 120,
      showOverflow: 'tooltip',
      title: '用户名',
      formatter: ({ row }) => row.nickname || `用户 #${row.user_id}`,
    },
    {
      field: 'phone',
      minWidth: 120,
      showOverflow: 'tooltip',
      title: '手机号码',
      formatter: ({ cellValue }) =>
        cellValue === undefined || cellValue === null || cellValue === ''
          ? '—'
          : String(cellValue),
    },
    {
      field: 'plan_name',
      minWidth: 120,
      showOverflow: 'tooltip',
      title: '会员卡名称',
    },
    {
      field: 'amount',
      title: '支付金额(元)',
      width: 110,
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
    },
    {
      field: 'pay_type',
      title: '支付方式',
      width: 100,
      formatter: ({ cellValue }) => payTypeText(String(cellValue || '')),
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '支付状态',
      width: 100,
    },
    {
      field: 'created_at',
      minWidth: 170,
      showOverflow: 'tooltip',
      title: '购买时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    {
      field: 'end_time',
      minWidth: 170,
      showOverflow: 'tooltip',
      title: '到期时间',
      formatter: ({ row }) => endTimeText(row),
    },
  ],
  emptyText: '暂无数据',
  pagerConfig: platformListPagerConfig(),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const values =
          formValues && Object.keys(formValues).length > 0
            ? formValues
            : lastFormValues.value;
        const filters = buildFilterParams(values);
        const [listResult, stats] = await Promise.all([
          listSvipOrders({
            page: page.currentPage,
            limit: page.pageSize,
            ...filters,
          }),
          getSvipOrderSummary(filters).catch(() => emptySummary()),
        ]);
        summary.value = stats || emptySummary();
        return {
          items: listResult.list || [],
          total: listResult.total || 0,
        };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'order_no' },
  scrollX: { enabled: true },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ gridOptions });

onMounted(async () => {
  const [profile, codes] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canRead.value =
    profile.roles.some((role) => role === 'platform' || role === 'operations') &&
    codes.includes('user.svip.record.read');
  if (canRead.value) gridApi.reload();
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有会员记录查看权限"
      type="warning"
      :closable="false"
    />
    <template v-else>
      <div class="svip-record-filter">
        <Form />
      </div>

      <div class="svip-record-summary">
        <div class="svip-summary">
          <div
            v-for="card in summaryCards"
            :key="card.key"
            class="svip-summary__card"
            :class="`svip-summary__card--${card.tone}`"
          >
            <div class="svip-summary__icon">
              <IconifyIcon :icon="card.icon" />
            </div>
            <div class="svip-summary__body">
              <div class="svip-summary__value">{{ card.value }}</div>
              <div class="svip-summary__label">{{ card.label }}</div>
            </div>
          </div>
        </div>
      </div>

      <Grid>
        <template #status="{ row }">
          <ElTag :type="payStatusType(row)">
            {{ payStatusText(row) }}
          </ElTag>
        </template>
      </Grid>
    </template>
  </Page>
</template>

<style scoped>
.svip-record-filter {
  padding: 12px 8px 4px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.svip-record-summary {
  padding: 16px;
  margin-bottom: 12px;
  background: hsl(var(--card));
  border-radius: 0.375rem;
}

.svip-summary {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  width: 100%;
}

.svip-summary__card {
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

.svip-summary__icon {
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

.svip-summary__card--blue .svip-summary__icon {
  background: #409eff;
}

.svip-summary__card--orange .svip-summary__icon {
  background: #e6a23c;
}

.svip-summary__card--pink .svip-summary__icon {
  background: #f56c6c;
}

.svip-summary__value {
  color: var(--el-text-color-primary);
  font-size: 24px;
  font-weight: 600;
  line-height: 1.2;
}

.svip-summary__label {
  margin-top: 6px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  line-height: 1.2;
}

@media (max-width: 1100px) {
  .svip-summary {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .svip-summary {
    grid-template-columns: 1fr;
  }
}
</style>
