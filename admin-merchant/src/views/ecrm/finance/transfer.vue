<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';
import { ElButton, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  getMerchantBalanceApi,
  listMerchantWithdrawsApi,
  type MerchantWithdraw,
} from '#/api/core/merchant-finance';
import {
  MERCHANT_LIST_GRID_LAYOUT,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const router = useRouter();
const balance = ref(0);

const auditStatus: Record<
  number,
  { label: string; type: 'danger' | 'info' | 'success' | 'warning' }
> = {
  [-1]: { label: '审核拒绝', type: 'danger' },
  0: { label: '待平台审核', type: 'warning' },
  1: { label: '审核通过', type: 'success' },
};

function auditInfo(status: number) {
  return auditStatus[status] || { label: '未知状态', type: 'info' as const };
}

async function refreshBalance() {
  const currentBalance = await getMerchantBalanceApi();
  balance.value = currentBalance.mer_money;
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '申请单号' },
    fieldName: 'financial_sn',
    label: '申请单号',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: [
        { label: '待平台审核', value: 0 },
        { label: '审核通过', value: 1 },
        { label: '审核拒绝', value: -1 },
      ],
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '审核状态',
  },
]);

const gridOptions: VxeGridProps<MerchantWithdraw> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    {
      field: 'financial_sn',
      minWidth: 180,
      showOverflow: false,
      title: '申请单号',
    },
    {
      field: 'extract_money',
      title: '提现金额',
      width: 116,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'status',
      slots: { default: 'audit' },
      title: '审核状态',
      width: 116,
    },
    {
      field: 'create_time',
      minWidth: 170,
      title: '申请时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = formValues?.status;
        await refreshBalance();
        const data = await listMerchantWithdrawsApi({
          page: page.currentPage,
          limit: page.pageSize,
          financial_sn:
            String(formValues?.financial_sn ?? '').trim() || undefined,
          status:
            status === -1 || status === 0 || status === 1
              ? Number(status)
              : undefined,
          date_from: range[0],
          date_to: range[1],
        });
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'financial_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid] = useVbenVxeGrid({ formOptions, gridOptions });

onMounted(() => void refreshBalance());
</script>

<template>
  <Page auto-content-height>
    <template #extra>
      <ElButton type="primary" @click="router.push('/finance/withdraw')">
        申请提现
      </ElButton>
    </template>

    <div class="mb-4 flex flex-wrap items-center gap-3 rounded-lg border bg-card px-4 py-3">
      <div>
        <span class="text-sm text-muted-foreground">可用余额（元）</span>
        <span class="ml-3 text-2xl font-semibold">
          ¥{{ Number(balance).toFixed(2) }}
        </span>
      </div>
    </div>

    <Grid>
      <template #audit="{ row }">
        <ElTag :type="auditInfo(row.status).type">
          {{ auditInfo(row.status).label }}
        </ElTag>
      </template>
    </Grid>
  </Page>
</template>
