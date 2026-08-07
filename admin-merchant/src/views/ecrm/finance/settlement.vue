<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page, confirm, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElDescriptions,
  ElDescriptionsItem,
  ElMessage,
  ElTag,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  applyMerchantSettlementApi,
  getMerchantSettlementApi,
  listMerchantSettlementsApi,
  type MerchantSettlement,
  type MerchantSettlementStatus,
} from '#/api/core/merchant-finance';
import {
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListActionColumn,
} from '#/constants/merchant-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const canApply = ref(false);
const applying = ref(false);
const detail = ref<MerchantSettlement>();

const statusLabels: Record<MerchantSettlementStatus, string> = {
  approved: '已审核',
  bill_frozen: '账期已冻结',
  bill_pending: '账期待生成',
  paid: '已打款',
  rejected: '已拒绝',
  withdraw_applied: '待平台审核',
};

const statusTypes: Record<
  MerchantSettlementStatus,
  'danger' | 'info' | 'success' | 'warning'
> = {
  approved: 'success',
  bill_frozen: 'warning',
  bill_pending: 'info',
  paid: 'success',
  rejected: 'danger',
  withdraw_applied: 'warning',
};

function canApplyRow(row: MerchantSettlement) {
  return canApply.value && row.status === 'bill_frozen' && row.amount > 0;
}

function idempotencyKey(settlementId: number) {
  return `merchant-apply-${settlementId}-${crypto.randomUUID()}`;
}

const formOptions: VbenFormProps = listFormOptionsDefaults([
  LIST_DATE_RANGE_FIELD,
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: Object.entries(statusLabels).map(([value, label]) => ({
        label,
        value,
      })),
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '结算状态',
  },
]);

const gridOptions: VxeGridProps<MerchantSettlement> = {
  ...MERCHANT_LIST_GRID_LAYOUT,
  columns: [
    { field: 'settlement_id', title: '结算 ID', width: 100 },
    {
      field: 'period_start',
      minWidth: 280,
      showOverflow: false,
      slots: { default: 'period' },
      title: '结算周期',
    },
    {
      field: 'amount',
      title: '结算金额',
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
      field: 'application_no',
      minWidth: 140,
      showOverflow: false,
      title: '申请单号',
    },
    {
      field: 'updated_at',
      minWidth: 170,
      title: '更新时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    merchantListActionColumn({ width: 148 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const range = Array.isArray(formValues?.date_range)
          ? formValues.date_range
          : [];
        const status = String(formValues?.status ?? '').trim();
        const data = await listMerchantSettlementsApi({
          page: page.currentPage,
          limit: page.pageSize,
          status: (status || undefined) as MerchantSettlementStatus | undefined,
          date_from: range[0],
          date_to: range[1],
        });
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [DetailDrawer, detailDrawerApi] = useVbenDrawer({
  class: 'w-[560px] max-w-[96vw]',
  showConfirmButton: false,
  cancelText: '关闭',
  placement: 'right',
});

async function openDetail(row: MerchantSettlement) {
  detailDrawerApi.setState({ title: '结算详情', loading: true }).open();
  try {
    detail.value = await getMerchantSettlementApi(row.settlement_id);
  } finally {
    detailDrawerApi.setState({ loading: false });
  }
}

async function apply(row: MerchantSettlement) {
  try {
    await confirm({
      content: `确认提交结算申请？账期 ${formatShanghaiDateTime(row.period_start)} 至 ${formatShanghaiDateTime(row.period_end)}，金额 ¥${Number(row.amount).toFixed(2)}。提交后由平台审核打款，本页不录入收款账户。`,
      icon: 'warning',
      title: '提交结算申请',
    });
  } catch {
    return;
  }
  applying.value = true;
  try {
    await applyMerchantSettlementApi(row.settlement_id, {
      idempotency_key: idempotencyKey(row.settlement_id),
    });
    ElMessage.success('结算申请已提交，等待平台审核');
    gridApi.reload();
    if (detail.value?.settlement_id === row.settlement_id) {
      detail.value = await getMerchantSettlementApi(row.settlement_id);
    }
  } finally {
    applying.value = false;
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi().catch(() => [] as string[]);
  canApply.value = permissions.includes('finance.settlement.apply');
});
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #period="{ row }">
        {{ formatShanghaiDateTime(row.period_start) }}
        至
        {{ formatShanghaiDateTime(row.period_end) }}
      </template>
      <template #status="{ row }">
        <ElTag :type="statusTypes[row.status]" size="small">
          {{ statusLabels[row.status] }}
        </ElTag>
      </template>
      <template #action="{ row }">
        <ElButton link type="primary" @click="openDetail(row)">详情</ElButton>
        <ElButton
          v-if="canApplyRow(row)"
          :loading="applying"
          link
          type="success"
          @click="apply(row)"
        >
          申请结算
        </ElButton>
      </template>
    </Grid>

    <DetailDrawer>
      <template v-if="detail">
        <ElDescriptions :column="1" border>
          <ElDescriptionsItem label="结算 ID">
            {{ detail.settlement_id }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="结算周期">
            {{ formatShanghaiDateTime(detail.period_start) }}
            至
            {{ formatShanghaiDateTime(detail.period_end) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="结算金额">
            ¥{{ Number(detail.amount).toFixed(2) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="状态">
            <ElTag :type="statusTypes[detail.status]">
              {{ statusLabels[detail.status] }}
            </ElTag>
          </ElDescriptionsItem>
          <ElDescriptionsItem label="申请单号">
            {{ detail.application_no || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="审核备注">
            {{ detail.review_note || '—' }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="打款时间">
            {{ formatShanghaiDateTime(detail.paid_at) }}
          </ElDescriptionsItem>
          <ElDescriptionsItem label="更新时间">
            {{ formatShanghaiDateTime(detail.updated_at) }}
          </ElDescriptionsItem>
        </ElDescriptions>
        <div v-if="canApplyRow(detail)" class="mt-6">
          <ElButton :loading="applying" type="primary" @click="apply(detail)">
            提交结算申请
          </ElButton>
        </div>
      </template>
    </DetailDrawer>
  </Page>
</template>
