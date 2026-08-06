<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert, ElCard, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  getSvipOrderSummary,
  listSvipOrders,
  type SvipOrder,
  type SvipOrderSummary,
} from '#/api/core/platform-svip-plan';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canRead = ref(false);
const summary = ref<SvipOrderSummary>();

function orderStatusText(status: SvipOrder['status']) {
  return (
    { pending: '待支付', paid: '已支付', closed: '已关闭' }[status] || status
  );
}

function orderStatusType(status: SvipOrder['status']) {
  return (
    { pending: 'warning', paid: 'success', closed: 'info' }[status] || 'info'
  ) as 'info' | 'success' | 'warning';
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '订单号 / 用户 ID',
    },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待支付', value: 'pending' },
        { label: '已支付', value: 'paid' },
        { label: '已关闭', value: 'closed' },
      ],
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '状态',
  },
]);

const gridOptions: VxeGridProps<SvipOrder> = {
  columns: [
    { field: 'order_no', minWidth: 180, showOverflow: false, title: '会员订单号' },
    { field: 'user_id', title: '用户 ID', width: 100 },
    { field: 'plan_name', minWidth: 130, showOverflow: false, title: '会员类型' },
    {
      field: 'duration_days',
      formatter: ({ row }) =>
        row.plan_type === 'lifetime' ? '永久' : `${row.duration_days} 天`,
      title: '有效期',
      width: 100,
    },
    {
      field: 'amount',
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
      title: '金额',
      width: 100,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 90,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
      minWidth: 170,
      title: '创建时间',
    },
    {
      field: 'paid_at',
      formatter: ({ cellValue }) =>
        cellValue ? formatShanghaiDateTime(cellValue) : '—',
      minWidth: 170,
      title: '支付时间',
    },
  ],
  pagerConfig: { enabled: true, pageSize: 20, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) return { items: [], total: 0 };
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = String(formValues?.status ?? '').trim() || undefined;
        const [listResult, stats] = await Promise.all([
          listSvipOrders({
            page: page.currentPage,
            limit: page.pageSize,
            status,
            keyword: String(formValues?.keyword ?? '').trim() || undefined,
            date_from: range[0] as string | undefined,
            date_to: range[1] as string | undefined,
          }),
          getSvipOrderSummary(),
        ]);
        summary.value = stats;
        return {
          items: listResult.list || [],
          total: listResult.total || 0,
        };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'order_no' },
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
      <div class="mb-4 grid grid-cols-2 gap-4 md:grid-cols-5">
        <ElCard shadow="never">总订单：{{ summary?.total || 0 }}</ElCard>
        <ElCard shadow="never">待支付：{{ summary?.pending || 0 }}</ElCard>
        <ElCard shadow="never">已支付：{{ summary?.paid || 0 }}</ElCard>
        <ElCard shadow="never">已关闭：{{ summary?.closed || 0 }}</ElCard>
        <ElCard shadow="never">
          已支付金额：¥{{ Number(summary?.paid_amount || 0).toFixed(2) }}
        </ElCard>
      </div>
      <Grid>
        <template #status="{ row }">
          <ElTag :type="orderStatusType(row.status)">
            {{ orderStatusText(row.status) }}
          </ElTag>
        </template>
      </Grid>
    </template>
  </Page>
</template>
