<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElButton, ElMessage, ElMessageBox } from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  deductMerchantDeposit,
  fetchMerchantDepositRefunds,
  fetchMerchantDeposits,
  markMerchantDepositRefundPaid,
  reviewMerchantDepositRefund,
  type MerchantDepositAccount,
  type MerchantDepositRefund,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import { platformListActionColumn, platformListPagerConfig } from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

const activeTab = ref<'accounts' | 'refunds'>('accounts');
const canManage = ref(false);

const depositStateLabels: Record<string, string> = {
  funded: '已缴',
  pending: '待缴',
  shortfall: '不足',
  refund_pending: '退款中',
  refunded: '已退',
  not_required: '无',
};

const refundStatusLabels: Record<string, string> = {
  applied: '待审核',
  approved: '已通过',
  rejected: '已拒绝',
  paid: '已打款',
};

function formatTime(value?: string) {
  return value ? formatShanghaiDateTime(value) : '—';
}

function splitFields(value: string, count: number) {
  const fields = value.split(',').map((item) => item.trim());
  return fields.length === count ? fields : [];
}

function validateDeduction(value: string) {
  const [amount, reason, key] = splitFields(value, 3);
  const numericAmount = Number(amount);
  const isCentAmount =
    Number.isFinite(numericAmount) &&
    Math.abs(numericAmount * 100 - Math.round(numericAmount * 100)) < 0.000001;
  if (!amount || !reason || !key || !isCentAmount || numericAmount <= 0) {
    return '请填写精确到分的正数金额、扣减原因和幂等键，并以英文逗号分隔。';
  }
  return true;
}

function validatePayout(value: string) {
  const [key, reference] = splitFields(value, 2);
  if (
    !key ||
    !reference ||
    [...key].length < 8 ||
    [...key].length > 128 ||
    [...reference].length > 128
  ) {
    return '请填写 8–128 位幂等键和不超过 128 字的内部打款凭证号。';
  }
  return true;
}

function validateNote(value: string) {
  const note = value.trim();
  return note && [...note].length <= 500
    ? true
    : '审核说明不能为空，且不能超过 500 个字符。';
}

function isPromptDismissed(error: unknown) {
  return error === 'cancel' || error === 'close' || error === 'escape';
}

function buildAccountParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const merIdRaw = String(formValues?.mer_id ?? '').trim();
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    merchant_id: merIdRaw ? Number(merIdRaw) : undefined,
    status: String(formValues?.status ?? '').trim() || undefined,
  };
}

function buildRefundParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const merIdRaw = String(formValues?.mer_id ?? '').trim();
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    merchant_id: merIdRaw ? Number(merIdRaw) : undefined,
    status: String(formValues?.status ?? '').trim() || undefined,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
  };
}

const accountStatusOptions = [
  { label: '已缴', value: 'funded' },
  { label: '待缴', value: 'pending' },
  { label: '不足', value: 'shortfall' },
  { label: '退款中', value: 'refund_pending' },
  { label: '已退', value: 'refunded' },
  { label: '无', value: 'not_required' },
];

const refundStatusOptions = [
  { label: '待审核', value: 'applied' },
  { label: '已通过', value: 'approved' },
  { label: '已拒绝', value: 'rejected' },
  { label: '已打款', value: 'paid' },
];

const formOptions: VbenFormProps = listFormOptionsDefaults([
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '商户 ID',
    },
    fieldName: 'mer_id',
    label: '商户搜索',
  },
  {
    component: 'Input',
    componentProps: {
      clearable: true,
      placeholder: '状态关键词（可选）',
    },
    fieldName: 'keyword',
    label: '关键词',
  },
  {
    component: 'Select',
    componentProps: {
      clearable: true,
      options: accountStatusOptions,
      placeholder: '全部状态',
    },
    fieldName: 'status',
    label: '保证金状态',
  },
]);

const accountGridOptions: VxeGridProps<MerchantDepositAccount> = {
  columns: [
    { field: 'merchant_id', title: '商户 ID', width: 100 },
    {
      field: 'required_amount',
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
      title: '应缴',
      width: 120,
    },
    {
      field: 'available_amount',
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
      title: '可用',
      width: 120,
    },
    {
      field: 'state',
      formatter: ({ cellValue }) => depositStateLabels[cellValue] || cellValue,
      title: '状态',
      width: 100,
    },
    platformListActionColumn({ width: 88 }),
  ],
  pagerConfig: platformListPagerConfig({ pageSize: 20, pageSizes: [10, 20, 50, 100] }),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await fetchMerchantDeposits(buildAccountParams(page, formValues));
        const list = data.list || [];
        return { items: list, total: Number(data.total ?? 0) };
      },
    },
  },
  rowConfig: { isHover: true, keyField: 'merchant_id' },
  toolbarConfig: {
    custom: false,
    export: false,
    refresh: false,
    search: false,
    zoom: false,
  },
};

const refundGridOptions: VxeGridProps<MerchantDepositRefund> = {
  columns: [
    { field: 'id', title: '申请 ID', width: 100 },
    { field: 'merchant_id', title: '商户 ID', width: 100 },
    {
      field: 'amount',
      formatter: ({ cellValue }) => `¥${Number(cellValue || 0).toFixed(2)}`,
      title: '退款金额',
      width: 120,
    },
    {
      field: 'status',
      formatter: ({ cellValue }) => refundStatusLabels[cellValue] || cellValue,
      title: '状态',
      width: 100,
    },
    {
      field: 'reason',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 160,
      showOverflow: false,
      title: '申请原因',
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatTime(cellValue),
      minWidth: 170,
      title: '申请时间',
    },
    platformListActionColumn({ width: 200 }),
  ],
  pagerConfig: platformListPagerConfig({ pageSize: 20, pageSizes: [10, 20, 50, 100] }),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await fetchMerchantDepositRefunds(buildRefundParams(page, formValues));
        const list = data.list || [];
        return { items: list, total: Number(data.total ?? 0) };
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

const [AccountGrid, accountGridApi] = useVbenVxeGrid({
  formOptions,
  gridOptions: accountGridOptions,
});

const [RefundGrid, refundGridApi] = useVbenVxeGrid({
  formOptions: listFormOptionsDefaults([
    LIST_DATE_RANGE_FIELD,
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '商户 ID' },
      fieldName: 'mer_id',
      label: '商户 ID',
    },
    {
      component: 'Input',
      componentProps: { clearable: true, placeholder: '原因 / 备注关键词' },
      fieldName: 'keyword',
      label: '关键词',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: refundStatusOptions,
        placeholder: '全部状态',
      },
      fieldName: 'status',
      label: '退款状态',
    },
  ]),
  gridOptions: refundGridOptions,
});

function setTab(tab: 'accounts' | 'refunds') {
  if (activeTab.value === tab) return;
  activeTab.value = tab;
  if (tab === 'accounts') accountGridApi.reload();
  else refundGridApi.reload();
}

async function deduct(row: MerchantDepositAccount) {
  try {
    const { value } = await ElMessageBox.prompt(
      '填写“金额,扣减原因,幂等键”，例如 10,虚构违规扣减,deposit-demo-001。',
      '扣除保证金',
      { inputValidator: validateDeduction },
    );
    const [amount, reason, idempotency_key] = splitFields(value, 3);
    await deductMerchantDeposit(row.merchant_id, {
      amount: Number(amount),
      reason,
      idempotency_key,
    });
    ElMessage.success('保证金扣减已登记');
    accountGridApi.reload();
  } catch (error) {
    if (!isPromptDismissed(error)) throw error;
  }
}

async function review(row: MerchantDepositRefund, approved: boolean) {
  try {
    const { value } = await ElMessageBox.prompt(
      '填写审核说明。',
      approved ? '同意退款' : '拒绝退款',
      { inputValidator: validateNote },
    );
    await reviewMerchantDepositRefund(row.id, approved, value.trim());
    ElMessage.success('退款审核已保存');
    refundGridApi.reload();
  } catch (error) {
    if (!isPromptDismissed(error)) throw error;
  }
}

async function paid(row: MerchantDepositRefund) {
  try {
    const { value } = await ElMessageBox.prompt(
      '填写“幂等键,内部打款凭证号”；不得填写账户信息。',
      '登记保证金退款打款',
      { inputValidator: validatePayout },
    );
    const [idempotency_key, payout_reference] = splitFields(value, 2);
    await markMerchantDepositRefundPaid(row.id, { idempotency_key, payout_reference });
    ElMessage.success('打款登记已保存');
    refundGridApi.reload();
  } catch (error) {
    if (!isPromptDismissed(error)) throw error;
  }
}

onMounted(async () => {
  const codes = await getAccessCodesApi();
  canManage.value = codes.includes('merchant.deposit.review');
});
</script>

<template>
  <Page auto-content-height>
    <AccountGrid v-if="activeTab === 'accounts'">
      <template #toolbar-actions>
        <div class="deposit-tabs" role="tablist">
          <button
            type="button"
            role="tab"
            class="deposit-tabs__item is-active"
            aria-selected="true"
            @click="setTab('accounts')"
          >
            保证金账户
          </button>
          <button
            type="button"
            role="tab"
            class="deposit-tabs__item"
            aria-selected="false"
            @click="setTab('refunds')"
          >
            退款申请
          </button>
        </div>
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="canManage"
          link
          type="danger"
          @click="deduct(row)"
        >
          扣减
        </ElButton>
        <span v-else>—</span>
      </template>
    </AccountGrid>

    <RefundGrid v-if="activeTab === 'refunds'">
      <template #toolbar-actions>
        <div class="deposit-tabs" role="tablist">
          <button
            type="button"
            role="tab"
            class="deposit-tabs__item"
            aria-selected="false"
            @click="setTab('accounts')"
          >
            保证金账户
          </button>
          <button
            type="button"
            role="tab"
            class="deposit-tabs__item is-active"
            aria-selected="true"
            @click="setTab('refunds')"
          >
            退款申请
          </button>
        </div>
      </template>
      <template #action="{ row }">
        <template v-if="canManage">
          <template v-if="row.status === 'applied'">
            <ElButton link type="success" @click="review(row, true)">同意</ElButton>
            <ElButton link type="danger" @click="review(row, false)">拒绝</ElButton>
          </template>
          <ElButton
            v-else-if="row.status === 'approved'"
            link
            type="primary"
            @click="paid(row)"
          >
            登记打款
          </ElButton>
          <span v-else>—</span>
        </template>
        <span v-else>—</span>
      </template>
    </RefundGrid>
  </Page>
</template>

<style scoped>
.deposit-tabs {
  display: flex;
  gap: 28px;
  border-bottom: 1px solid hsl(var(--border));
  width: 100%;
}

.deposit-tabs__item {
  margin-bottom: -1px;
  padding: 10px 2px 12px;
  border: 0;
  border-bottom: 2px solid transparent;
  background: transparent;
  color: hsl(var(--foreground) / 70%);
  font-size: 14px;
  line-height: 22px;
  cursor: pointer;
}

.deposit-tabs__item:hover {
  color: hsl(var(--primary));
}

.deposit-tabs__item.is-active {
  border-bottom-color: hsl(var(--primary));
  color: hsl(var(--primary));
  font-weight: 600;
}
</style>
