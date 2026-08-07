<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  listPlatformTransferRecordsApi,
  type MerchantSettlementRow,
  type TransferRecordStatus,
} from '#/api/core/platform-merchant-settlement';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const statusLabels: Record<TransferRecordStatus, string> = {
  approved: '待登记打款',
  paid: '已打款',
  rejected: '已拒绝',
};

const canRead = ref(false);

function buildListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const merIdRaw = String(formValues?.mer_id ?? '').trim();
  return {
    page: page.currentPage,
    limit: page.pageSize,
    merchant_id: merIdRaw ? Number(merIdRaw) : undefined,
    status: (String(formValues?.status ?? '').trim() ||
      undefined) as TransferRecordStatus | undefined,
    date_from: range[0],
    date_to: range[1],
  };
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Input',
    componentProps: { clearable: true, placeholder: '商户 ID' },
    fieldName: 'mer_id',
    label: '商户 ID',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: Object.entries(statusLabels).map(([value, label]) => ({
        label,
        value,
      })),
      placeholder: '打款链路',
    },
    fieldName: 'status',
    label: '打款状态',
  },
]);

const gridOptions: VxeGridProps<MerchantSettlementRow> = {
  columns: [
    { field: 'settlement_id', title: '结算 ID', width: 110 },
    {
      field: 'merchant_name',
      minWidth: 180,
      showOverflow: false,
      title: '商户',
      formatter: ({ row }) => `${row.merchant_name}（${row.merchant_id}）`,
    },
    { field: 'store_id', title: '店铺 ID', width: 100 },
    {
      field: 'period_start',
      minWidth: 280,
      showOverflow: false,
      title: '结算周期',
      formatter: ({ row }) =>
        `${formatShanghaiDateTime(row.period_start)} 至 ${formatShanghaiDateTime(row.period_end)}`,
    },
    {
      field: 'amount',
      title: '金额',
      width: 120,
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 120,
    },
    {
      field: 'updated_at',
      minWidth: 180,
      title: '投影更新时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) {
          return { items: [], total: 0 };
        }
        const data = await listPlatformTransferRecordsApi(
          buildListParams(page, formValues),
        );
        return { items: data.list || [], total: data.total || 0 };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'settlement_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const [Grid] = useVbenVxeGrid({ formOptions, gridOptions });

onMounted(async () => {
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canRead.value =
    profile.roles.some((role) => role === 'platform' || role === 'region') &&
    permissions.includes('accounts.merchant_settlement.read');
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      class="mb-4"
      title="本页不发起转账、不保存打款密钥；登记内部凭证请在「店铺结算监管」完成。"
      type="warning"
      :closable="false"
    />
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有转账记录（结算投影）查看权限"
      type="warning"
      :closable="false"
    />
    <Grid v-else>
      <template #status="{ row }">
        <ElTag>
          {{ statusLabels[row.status as TransferRecordStatus] || row.status }}
        </ElTag>
      </template>
    </Grid>
  </Page>
</template>
