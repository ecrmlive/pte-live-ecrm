<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElAlert, ElButton, ElCard, ElCol, ElMessage, ElMessageBox, ElRow, ElTag } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getAccessCodesApi, getUserInfoApi } from '#/api/core/auth';
import {
  approvePlatformMerchantSettlementApi,
  getPlatformMerchantSettlementSummaryApi,
  listPlatformMerchantSettlementsApi,
  markPlatformMerchantSettlementPaidApi,
  rejectPlatformMerchantSettlementApi,
  type MerchantSettlementRow,
  type MerchantSettlementStatus,
  type MerchantSettlementSummary,
} from '#/api/core/platform-merchant-settlement';
import { platformListActionColumn } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const statusLabels: Record<MerchantSettlementStatus, string> = {
  approved: '已审核',
  bill_frozen: '账期已冻结',
  bill_pending: '账期待生成',
  cancelled: '已撤销（历史）',
  paid: '已打款',
  rejected: '已拒绝',
  withdraw_applied: '待平台审核',
};

const SUMMARY_STATUSES: MerchantSettlementStatus[] = [
  'bill_frozen',
  'withdraw_applied',
  'paid',
];

const summary = ref<MerchantSettlementSummary[]>([]);
const canRead = ref(false);
const canReview = ref(false);

const summaryByStatus = computed(
  () => new Map(summary.value.map((item) => [item.status, item])),
);

function summaryText(status: MerchantSettlementStatus) {
  const item = summaryByStatus.value.get(status);
  return item ? `${item.count} 笔 · ${item.amount.toFixed(2)}` : '暂无记录';
}

function idempotencyKey(action: string, settlementId: number) {
  return `${action}-${settlementId}-${crypto.randomUUID()}`;
}

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
      undefined) as MerchantSettlementStatus | undefined,
    date_from: range[0],
    date_to: range[1],
  };
}

async function loadSummary() {
  if (!canRead.value) return;
  const summaryData = await getPlatformMerchantSettlementSummaryApi();
  summary.value = summaryData.list || [];
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
      placeholder: '全部',
    },
    fieldName: 'status',
    label: '结算状态',
  },
]);

const gridOptions: VxeGridProps<MerchantSettlementRow> = {
  columns: [
    { field: 'settlement_id', title: '结算 ID', width: 110 },
    {
      field: 'merchant_name',
      minWidth: 160,
      showOverflow: false,
      title: '商户',
      formatter: ({ row }) => `${row.merchant_name}（${row.merchant_id}）`,
    },
    { field: 'store_id', title: '店铺 ID', width: 100 },
    {
      field: 'period_start',
      minWidth: 300,
      showOverflow: false,
      title: '结算周期',
      formatter: ({ row }) =>
        `${formatShanghaiDateTime(row.period_start)} 至 ${formatShanghaiDateTime(row.period_end)}`,
    },
    {
      field: 'amount',
      title: '结算金额',
      width: 130,
      formatter: ({ cellValue }) => Number(cellValue || 0).toFixed(2),
    },
    {
      field: 'status',
      slots: { default: 'status' },
      title: '状态',
      width: 130,
    },
    {
      field: 'updated_at',
      minWidth: 180,
      title: '投影更新时间',
      formatter: ({ cellValue }) => formatShanghaiDateTime(cellValue),
    },
    platformListActionColumn({ width: 170 }),
  ],
  pagerConfig: { enabled: true, pageSize: 10, pageSizes: [10, 20, 50, 100] },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        if (!canRead.value) {
          return { items: [], total: 0 };
        }
        const data = await listPlatformMerchantSettlementsApi(
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

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

async function reloadGrid() {
  await gridApi.reload();
  window.setTimeout(() => void loadSummary(), 2300);
}

async function approve(row: MerchantSettlementRow) {
  try {
    await ElMessageBox.confirm(
      '确认审核通过该店铺结算申请？此操作不会直接访问商户库，状态将通过受控命令与事件投影更新。',
      '结算审核',
      { type: 'warning' },
    );
    await approvePlatformMerchantSettlementApi(row.settlement_id, {
      idempotency_key: idempotencyKey('approve', row.settlement_id),
    });
    ElMessage.success('审核命令已完成，监管投影正在刷新');
    await reloadGrid();
  } catch {
    /* 用户取消或接口已返回错误时，requestClient 统一提示。 */
  }
}

async function reject(row: MerchantSettlementRow) {
  try {
    const { value } = await ElMessageBox.prompt('请填写拒绝原因', '拒绝结算申请', {
      inputPattern: /\S+/,
      inputErrorMessage: '拒绝原因必填',
    });
    await rejectPlatformMerchantSettlementApi(row.settlement_id, {
      idempotency_key: idempotencyKey('reject', row.settlement_id),
      review_note: value.trim(),
    });
    ElMessage.success('拒绝命令已完成，监管投影正在刷新');
    await reloadGrid();
  } catch {
    /* 用户取消或接口已返回错误时，requestClient 统一提示。 */
  }
}

async function markPaid(row: MerchantSettlementRow) {
  try {
    const { value } = await ElMessageBox.prompt(
      '请输入内部打款凭证编号（不录入银行卡、账号或密钥）',
      '登记打款凭证',
      { inputPattern: /\S{3,}/, inputErrorMessage: '凭证编号至少 3 个字符' },
    );
    await markPlatformMerchantSettlementPaidApi(row.settlement_id, {
      idempotency_key: idempotencyKey('paid', row.settlement_id),
      payout_reference: value.trim(),
    });
    ElMessage.success('打款登记命令已完成，监管投影正在刷新');
    await reloadGrid();
  } catch {
    /* 用户取消或接口已返回错误时，requestClient 统一提示。 */
  }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([
    getUserInfoApi(),
    getAccessCodesApi(),
  ]);
  canRead.value =
    profile.roles.some((role) => role === 'platform' || role === 'region') &&
    permissions.includes('accounts.merchant_settlement.read');
  canReview.value =
    profile.roles.includes('platform') &&
    permissions.includes('accounts.merchant_settlement.review');
  if (canRead.value) {
    await loadSummary();
  }
});
</script>

<template>
  <Page auto-content-height>
    <ElAlert
      v-if="!canRead"
      class="mb-4"
      title="当前账号没有查看店铺结算监管投影的权限"
      type="warning"
      :closable="false"
    />
    <template v-else>
      <ElRow :gutter="16" class="mb-4">
        <ElCol
          v-for="status in SUMMARY_STATUSES"
          :key="status"
          :md="8"
          :xs="24"
        >
          <ElCard shadow="never">
            <div class="text-sm text-gray-500">{{ statusLabels[status] }}</div>
            <div class="mt-2 text-sm">{{ summaryText(status) }}</div>
          </ElCard>
        </ElCol>
      </ElRow>

      <Grid>
        <template #status="{ row }">
          <ElTag>{{ statusLabels[row.status] }}</ElTag>
        </template>
        <template v-if="canReview" #action="{ row }">
          <template v-if="row.status === 'withdraw_applied'">
            <ElButton link type="success" @click="approve(row)">通过</ElButton>
            <ElButton link type="danger" @click="reject(row)">拒绝</ElButton>
          </template>
          <ElButton
            v-else-if="row.status === 'approved'"
            link
            type="primary"
            @click="markPaid(row)"
          >
            登记打款
          </ElButton>
        </template>
      </Grid>
    </template>
  </Page>
</template>
