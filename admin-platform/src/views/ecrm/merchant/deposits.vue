<script setup lang="ts">
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, onMounted, reactive, ref } from 'vue';

import { Page, useVbenDrawer } from '@vben/common-ui';
import {
  ElButton,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElMessageBox,
  ElPagination,
  ElRadio,
  ElRadioGroup,
  ElTable,
  ElTableColumn,
} from 'element-plus';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  deductMerchantDeposit,
  fetchMerchantCategories,
  fetchMerchantDepositLedgers,
  fetchMerchantDepositRefunds,
  fetchMerchantDeposits,
  fetchMerchantTypes,
  fundMerchantDepositOffline,
  markMerchantDepositRefundNote,
  markMerchantDepositRefundPaid,
  reviewMerchantDepositRefund,
  type MerchantCategoryRow,
  type MerchantDepositAccount,
  type MerchantDepositLedger,
  type MerchantDepositRefund,
  type MerchantTypeRow,
} from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
import {
  platformListActionColumn,
  platformListPagerConfig,
} from '#/constants/platform-list-grid';
import { formatShanghaiDateTime } from '#/utils/date-time';
import {
  LIST_DATE_RANGE_FIELD,
  listFormOptionsDefaults,
} from '#/utils/list-form-defaults';

type DepositTab = 'pending' | 'funded' | 'refunds';

/** CRMEB：缴存保证金默认选中 */
const activeTab = ref<DepositTab>('funded');
const canManage = ref(false);
const categories = ref<MerchantCategoryRow[]>([]);
const types = ref<MerchantTypeRow[]>([]);

const ledgerTitleLabels: Record<string, string> = {
  fund: '线下补缴保证金',
  deduct: '保证金扣除',
  refund_approved: '退款审核通过',
  refund_rejected: '退款审核拒绝',
  refund_paid: '保证金退还',
};

const refundStatusLabels: Record<string, string> = {
  applied: '待审核',
  approved: '已通过',
  rejected: '未通过',
  paid: '已退回',
};

const ledgerMerchantId = ref(0);
const ledgerLoading = ref(false);
const ledgerRows = ref<MerchantDepositLedger[]>([]);
const ledgerPage = ref(1);
const ledgerLimit = ref(10);
const ledgerTotal = ref(0);

/** CRMEB 操作记录弹窗 width:800px */
const [LedgerDrawer, ledgerDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  footer: false,
  title: '操作记录',
});

/** CRMEB localMarginForm：线下缴纳保证金（展示额度/待缴 + 备注 + 状态） */
const offlineTarget = ref<MerchantDepositAccount | null>(null);
const offlineForm = reactive({
  mark: '',
  /** 0 未缴纳仅关闭；1 已缴纳走 fund-offline（金额取待缴，与 CRMEB hidden number 一致） */
  status: 1 as 0 | 1,
});

const [OfflineDrawer, offlineDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  confirmText: '保存',
  cancelText: '取消',
  onConfirm: async () => submitOffline(),
  title: '线下缴纳保证金',
});

/** CRMEB setMarginForm：扣除保证金（金额 + 原因），幂等键前端静默生成 */
const deductTarget = ref<MerchantDepositAccount | null>(null);
const deductForm = reactive({
  amount: undefined as number | undefined,
  reason: '',
});

const [DeductDrawer, deductDrawerApi] = useVbenDrawer({
  class: 'w-[1000px] max-w-[96vw]',
  placement: 'right',
  confirmText: '保存',
  cancelText: '取消',
  onConfirm: async () => submitDeduct(),
  title: '扣除保证金',
});

function formatMoney(value?: number | string | null) {
  return Number(value || 0).toFixed(2);
}

function formatTime(value?: string | null) {
  return value ? formatShanghaiDateTime(value) : '—';
}

function isCentAmount(value: number) {
  return (
    Number.isFinite(value) &&
    Math.abs(value * 100 - Math.round(value * 100)) < 0.000001
  );
}

function splitFields(value: string, count: number) {
  const fields = value.split(',').map((item) => item.trim());
  return fields.length === count ? fields : [];
}

function isPromptDismissed(error: unknown) {
  return error === 'cancel' || error === 'close' || error === 'escape';
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
  return [...note].length <= 500 ? true : '备注不能超过 500 个字符。';
}

function validateReviewNote(value: string) {
  const note = value.trim();
  return note && [...note].length <= 500
    ? true
    : '审核说明不能为空，且不能超过 500 个字符。';
}

function newIdempotencyKey(prefix: string, merchantId: number) {
  const uuid =
    typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function'
      ? crypto.randomUUID()
      : `${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
  return `${prefix}-${merchantId}-${uuid}`;
}

function payableOf(row: MerchantDepositAccount | null | undefined) {
  if (!row) return 0;
  return Number(
    row.payable_amount ??
      Math.max(Number(row.required_amount || 0) - Number(row.available_amount || 0), 0),
  );
}

const traderOptions = [
  { label: '自营', value: 1 },
  { label: '非自营', value: 0 },
];

const refundStatusOptions = [
  { label: '待审核', value: 'applied' },
  { label: '已通过', value: 'approved' },
  { label: '未通过', value: 'rejected' },
  { label: '已退回', value: 'paid' },
];

const categoryOptions = computed(() =>
  categories.value.map((item) => ({
    label: item.category_name,
    value: item.merchant_category_id,
  })),
);

const typeOptions = computed(() =>
  types.value.map((item) => ({
    label: item.name,
    value: item.id,
  })),
);

function buildAccountParams(
  page: { currentPage: number; pageSize: number },
  formValues: Record<string, unknown> | undefined,
  tab: 'pending' | 'funded',
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const typeRaw = formValues?.type_id;
  const categoryRaw = formValues?.category_id;
  const traderRaw = formValues?.is_trader;
  return {
    tab,
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    type_id: typeRaw === 0 || typeRaw ? Number(typeRaw) : undefined,
    category_id: categoryRaw === 0 || categoryRaw ? Number(categoryRaw) : undefined,
    is_trader:
      traderRaw === 0 || traderRaw === 1 || traderRaw === '0' || traderRaw === '1'
        ? Number(traderRaw)
        : undefined,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
  };
}

function buildRefundParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const typeRaw = formValues?.type_id;
  const categoryRaw = formValues?.category_id;
  const traderRaw = formValues?.is_trader;
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    status: String(formValues?.status ?? '').trim() || undefined,
    type_id: typeRaw === 0 || typeRaw ? Number(typeRaw) : undefined,
    category_id: categoryRaw === 0 || categoryRaw ? Number(categoryRaw) : undefined,
    is_trader:
      traderRaw === 0 || traderRaw === 1 || traderRaw === '0' || traderRaw === '1'
        ? Number(traderRaw)
        : undefined,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
  };
}

function accountFilterSchema(traderLabel: string): NonNullable<VbenFormProps['schema']> {
  return [
    LIST_DATE_RANGE_FIELD,
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: traderOptions,
        placeholder: '全部',
      },
      fieldName: 'is_trader',
      label: traderLabel,
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [],
        placeholder: '全部',
      },
      fieldName: 'category_id',
      label: '店铺分类',
    },
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: [],
        placeholder: '全部',
      },
      fieldName: 'type_id',
      label: '店铺类型',
    },
    {
      component: 'Input',
      componentProps: {
        clearable: true,
        placeholder: '请输入店铺关键字/店铺名/联系电话',
      },
      fieldName: 'keyword',
      label: '关键字',
    },
  ];
}

const pendingGridOptions: VxeGridProps<MerchantDepositAccount> = {
  columns: [
    { field: 'merchant_id', title: 'ID', width: 90 },
    {
      field: 'merchant_name',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 140,
      title: '店铺名称',
    },
    {
      field: 'type_name',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 110,
      title: '店铺类型',
    },
    {
      field: 'owner_name',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 110,
      title: '店铺姓名',
    },
    {
      field: 'payable_amount',
      formatter: ({ cellValue, row }) =>
        formatMoney(
          cellValue ??
            Math.max(Number(row.required_amount || 0) - Number(row.available_amount || 0), 0),
        ),
      title: '待缴金额',
      width: 120,
    },
    {
      field: 'required_amount',
      formatter: ({ cellValue }) => formatMoney(cellValue),
      title: '保证金额度',
      width: 120,
    },
    {
      field: 'mark',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 120,
      showOverflow: false,
      title: '备注',
    },
    platformListActionColumn({ width: 100 }),
  ],
  pagerConfig: platformListPagerConfig({ pageSize: 10, pageSizes: [10, 20, 50, 100] }),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await fetchMerchantDeposits(
          buildAccountParams(page, formValues, 'pending'),
        );
        return { items: data.list || [], total: Number(data.total ?? 0) };
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

const fundedGridOptions: VxeGridProps<MerchantDepositAccount> = {
  columns: [
    { field: 'merchant_id', title: 'ID', width: 90 },
    {
      field: 'merchant_name',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 140,
      title: '店铺名称',
    },
    {
      field: 'type_name',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 110,
      title: '店铺类型',
    },
    {
      field: 'owner_name',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 110,
      title: '店铺姓名',
    },
    {
      field: 'required_amount',
      formatter: ({ cellValue }) => formatMoney(cellValue),
      title: '剩余应缴',
      width: 120,
    },
    {
      field: 'paid_at',
      formatter: ({ cellValue }) => formatTime(cellValue),
      minWidth: 170,
      title: '支付时间',
    },
    {
      field: 'mark',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 120,
      showOverflow: false,
      title: '备注',
    },
    platformListActionColumn({ width: 180 }),
  ],
  pagerConfig: platformListPagerConfig({ pageSize: 10, pageSizes: [10, 20, 50, 100] }),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await fetchMerchantDeposits(
          buildAccountParams(page, formValues, 'funded'),
        );
        return { items: data.list || [], total: Number(data.total ?? 0) };
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
    { field: 'id', title: 'ID', width: 90 },
    {
      field: 'merchant_name',
      formatter: ({ cellValue, row }) => cellValue || `商户#${row.merchant_id}`,
      minWidth: 140,
      title: '店铺名称',
    },
    {
      field: 'owner_name',
      formatter: ({ cellValue }) => cellValue || '—',
      minWidth: 110,
      title: '店铺姓名',
    },
    {
      field: 'required_amount',
      formatter: ({ cellValue }) => formatMoney(cellValue),
      title: '剩余应缴',
      width: 110,
    },
    {
      field: 'status',
      formatter: ({ cellValue }) => refundStatusLabels[cellValue] || cellValue,
      title: '申请状态',
      width: 100,
    },
    {
      field: 'available_amount',
      formatter: ({ cellValue }) => formatMoney(cellValue),
      title: '结余保证金',
      width: 120,
    },
    {
      field: 'amount',
      formatter: ({ cellValue }) => formatMoney(cellValue),
      title: '退款金额',
      width: 110,
    },
    {
      field: 'refund_method',
      formatter: ({ cellValue }) => cellValue || '线下',
      title: '退回方式',
      width: 100,
    },
    {
      field: 'created_at',
      formatter: ({ cellValue }) => formatTime(cellValue),
      minWidth: 170,
      title: '申请时间',
    },
    {
      field: 'review_note',
      formatter: ({ cellValue, row }) => cellValue || row.reason || '—',
      minWidth: 120,
      showOverflow: false,
      title: '备注',
    },
    platformListActionColumn({ width: 220 }),
  ],
  pagerConfig: platformListPagerConfig({ pageSize: 10, pageSizes: [10, 20, 50, 100] }),
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        const data = await fetchMerchantDepositRefunds(buildRefundParams(page, formValues));
        return { items: data.list || [], total: Number(data.total ?? 0) };
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

const [PendingGrid, pendingGridApi] = useVbenVxeGrid({
  formOptions: listFormOptionsDefaults(accountFilterSchema('店铺级别')),
  gridOptions: pendingGridOptions,
});

const [FundedGrid, fundedGridApi] = useVbenVxeGrid({
  formOptions: listFormOptionsDefaults(accountFilterSchema('店铺类别')),
  gridOptions: fundedGridOptions,
});

const [RefundGrid, refundGridApi] = useVbenVxeGrid({
  formOptions: listFormOptionsDefaults([
    ...accountFilterSchema('店铺类别'),
    {
      component: 'Select',
      componentProps: {
        clearable: true,
        options: refundStatusOptions,
        placeholder: '全部状态',
      },
      fieldName: 'status',
      label: '状态',
    },
  ]),
  gridOptions: refundGridOptions,
});

function patchFilterSelectOptions() {
  const catOpts = categoryOptions.value;
  const typeOpts = typeOptions.value;
  const patch = [
    {
      fieldName: 'category_id',
      componentProps: { clearable: true, options: catOpts, placeholder: '全部' },
    },
    {
      fieldName: 'type_id',
      componentProps: { clearable: true, options: typeOpts, placeholder: '全部' },
    },
  ];
  pendingGridApi.formApi?.updateSchema?.(patch);
  fundedGridApi.formApi?.updateSchema?.(patch);
  refundGridApi.formApi?.updateSchema?.(patch);
}

function setTab(tab: DepositTab) {
  if (activeTab.value === tab) return;
  activeTab.value = tab;
  if (tab === 'pending') pendingGridApi.reload();
  else if (tab === 'funded') fundedGridApi.reload();
  else refundGridApi.reload();
}

function ledgerTitle(entryType: string) {
  return ledgerTitleLabels[entryType] || entryType || '—';
}

function ledgerSignedAmount(row: MerchantDepositLedger) {
  const raw = Number(row.amount || 0);
  const negative = row.entry_type === 'deduct' || row.entry_type === 'refund_paid';
  const value = Math.abs(raw);
  if (negative) {
    return { className: 'is-deduct', text: `-${value.toFixed(2)}` };
  }
  return { className: 'is-fund', text: `+${value.toFixed(2)}` };
}

function ledgerRemark(row: MerchantDepositLedger) {
  const reason = String(row.reason || '').trim();
  const opId = Number(row.operator_admin_id || 0);
  const opName = String(row.operator_name || '').trim() || '平台管理员';
  const opText = opId > 0 ? `【操作者：${opId}/${opName}】` : '';
  if (reason && opText) return `${reason} ${opText}`;
  if (reason) return reason;
  if (opText) return `操作者：${opId}/${opName}`;
  return '—';
}

async function loadLedgers() {
  const merchantId = Number(ledgerMerchantId.value || 0);
  if (!merchantId) return;
  ledgerLoading.value = true;
  try {
    const data = await fetchMerchantDepositLedgers(merchantId, {
      page: ledgerPage.value,
      limit: ledgerLimit.value,
    });
    ledgerRows.value = data?.list || [];
    ledgerTotal.value = Number(data?.total ?? 0);
  } catch {
    ledgerRows.value = [];
    ledgerTotal.value = 0;
  } finally {
    ledgerLoading.value = false;
  }
}

async function openLedgers(merchantId: number) {
  const id = Number(merchantId || 0);
  if (!id) {
    ElMessage.warning('商户 ID 无效');
    return;
  }
  ledgerMerchantId.value = id;
  ledgerPage.value = 1;
  ledgerRows.value = [];
  ledgerTotal.value = 0;
  ledgerDrawerApi.open();
  await loadLedgers();
}

function onLedgerPageChange(page: number) {
  ledgerPage.value = page;
  void loadLedgers();
}

function ledgerIndex(index: number) {
  return (ledgerPage.value - 1) * ledgerLimit.value + index + 1;
}

function openOffline(row: MerchantDepositAccount) {
  if (!canManage.value) return;
  const payable = payableOf(row);
  if (payable <= 0) {
    ElMessage.warning('当前无待缴保证金');
    return;
  }
  offlineTarget.value = row;
  offlineForm.mark = String(row.mark || '');
  // 按钮语义是「线下付款」，默认选已缴纳；仍保留 CRMEB 未缴纳选项
  offlineForm.status = 1;
  offlineDrawerApi.open();
}

async function submitOffline() {
  const row = offlineTarget.value;
  if (!row) return;
  const mark = offlineForm.mark.trim();
  if ([...mark].length > 500) {
    ElMessage.warning('备注不能超过 500 个字符');
    return;
  }
  // CRMEB：status=0 仅改备注不入账；本仓无单独改备注接口，关闭即可
  if (offlineForm.status !== 1) {
    ElMessage.info('未选择已缴纳，未登记线下付款');
    offlineDrawerApi.close();
    return;
  }
  const payable = payableOf(row);
  if (!isCentAmount(payable) || payable <= 0) {
    ElMessage.warning('待缴金额无效，无法登记线下付款');
    return;
  }
  offlineDrawerApi.lock();
  try {
    await fundMerchantDepositOffline(row.merchant_id, {
      // 与 CRMEB hidden number 一致：按待缴全额入账（API amount=0 也会回落到 payable）
      amount: payable,
      mark: mark || undefined,
      idempotency_key: newIdempotencyKey('deposit-offline', row.merchant_id),
    });
    ElMessage.success('线下付款已登记');
    offlineDrawerApi.close();
    pendingGridApi.reload();
    fundedGridApi.reload();
  } finally {
    offlineDrawerApi.unlock();
  }
}

function openDeduct(row: MerchantDepositAccount) {
  if (!canManage.value) return;
  const available = Number(row.available_amount || 0);
  if (available <= 0) {
    ElMessage.warning('商户剩余保证金不足，无法扣费');
    return;
  }
  deductTarget.value = row;
  deductForm.amount = undefined;
  deductForm.reason = '';
  deductDrawerApi.open();
}

async function submitDeduct() {
  const row = deductTarget.value;
  if (!row) return;
  const amount = Number(deductForm.amount);
  const reason = deductForm.reason.trim();
  const available = Number(row.available_amount || 0);
  if (!isCentAmount(amount) || amount <= 0) {
    ElMessage.warning('请输入精确到分的正数扣除金额');
    return;
  }
  if (amount > available + 0.000001) {
    ElMessage.warning(`扣除金额不能超过剩余保证金 ${formatMoney(available)}`);
    return;
  }
  if (!reason) {
    ElMessage.warning('请输入保证金扣除原因');
    return;
  }
  if ([...reason].length > 500) {
    ElMessage.warning('扣除原因不能超过 500 个字符');
    return;
  }
  deductDrawerApi.lock();
  try {
    await deductMerchantDeposit(row.merchant_id, {
      amount,
      reason,
      idempotency_key: newIdempotencyKey('deposit-deduct', row.merchant_id),
    });
    ElMessage.success('扣除保证金成功');
    deductDrawerApi.close();
    fundedGridApi.reload();
  } finally {
    deductDrawerApi.unlock();
  }
}

async function remarkRefund(row: MerchantDepositRefund) {
  if (!canManage.value) return;
  try {
    const { value } = await ElMessageBox.prompt('填写备注。', '备注', {
      inputValue: row.review_note || '',
      inputValidator: validateNote,
    });
    await markMerchantDepositRefundNote(row.id, value.trim());
    ElMessage.success('备注已保存');
    refundGridApi.reload();
  } catch (error) {
    if (!isPromptDismissed(error)) throw error;
  }
}

async function review(row: MerchantDepositRefund) {
  if (!canManage.value) return;
  try {
    await ElMessageBox.confirm('请选择审核结果。', '审核退回保证金', {
      distinguishCancelAndClose: true,
      confirmButtonText: '通过',
      cancelButtonText: '驳回',
      type: 'warning',
    })
      .then(async () => {
        const { value } = await ElMessageBox.prompt('填写审核说明。', '同意退款', {
          inputValidator: validateReviewNote,
        });
        await reviewMerchantDepositRefund(row.id, true, value.trim());
        ElMessage.success('退款审核已通过');
        refundGridApi.reload();
      })
      .catch(async (error) => {
        if (error === 'cancel') {
          const { value } = await ElMessageBox.prompt('填写审核说明。', '拒绝退款', {
            inputValidator: validateReviewNote,
          });
          await reviewMerchantDepositRefund(row.id, false, value.trim());
          ElMessage.success('退款已驳回');
          refundGridApi.reload();
          return;
        }
        if (!isPromptDismissed(error)) throw error;
      });
  } catch (error) {
    if (!isPromptDismissed(error)) throw error;
  }
}

async function paid(row: MerchantDepositRefund) {
  if (!canManage.value) return;
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
  const [codes, categoryData, typeData] = await Promise.all([
    getAccessCodesApi(),
    fetchMerchantCategories().catch(() => ({ list: [] as MerchantCategoryRow[] })),
    fetchMerchantTypes().catch(() => ({ list: [] as MerchantTypeRow[] })),
  ]);
  canManage.value = codes.includes('merchant.deposit.review');
  categories.value = categoryData.list || [];
  types.value = typeData.list || [];
  patchFilterSelectOptions();
});
</script>

<template>
  <Page auto-content-height>
    <FundedGrid v-if="activeTab === 'funded'">
      <template #toolbar-actions>
        <div class="deposit-tabs" role="tablist">
          <button
            type="button"
            role="tab"
            class="deposit-tabs__item"
            aria-selected="false"
            @click="setTab('pending')"
          >
            待缴保证金
          </button>
          <button
            type="button"
            role="tab"
            class="deposit-tabs__item is-active"
            aria-selected="true"
            @click="setTab('funded')"
          >
            缴存保证金
          </button>
          <button
            type="button"
            role="tab"
            class="deposit-tabs__item"
            aria-selected="false"
            @click="setTab('refunds')"
          >
            退回保证金
          </button>
        </div>
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="canManage && Number(row.available_amount || 0) > 0"
          link
          type="primary"
          @click="openDeduct(row)"
        >
          保证金扣费
        </ElButton>
        <ElButton link type="primary" @click="openLedgers(row.merchant_id)">
          操作记录
        </ElButton>
      </template>
    </FundedGrid>

    <PendingGrid v-else-if="activeTab === 'pending'">
      <template #toolbar-actions>
        <div class="deposit-tabs" role="tablist">
          <button
            type="button"
            role="tab"
            class="deposit-tabs__item is-active"
            aria-selected="true"
            @click="setTab('pending')"
          >
            待缴保证金
          </button>
          <button
            type="button"
            role="tab"
            class="deposit-tabs__item"
            aria-selected="false"
            @click="setTab('funded')"
          >
            缴存保证金
          </button>
          <button
            type="button"
            role="tab"
            class="deposit-tabs__item"
            aria-selected="false"
            @click="setTab('refunds')"
          >
            退回保证金
          </button>
        </div>
      </template>
      <template #action="{ row }">
        <ElButton
          v-if="canManage"
          link
          type="primary"
          @click="openOffline(row)"
        >
          线下付款
        </ElButton>
        <span v-else>—</span>
      </template>
    </PendingGrid>

    <RefundGrid v-else>
      <template #toolbar-actions>
        <div class="deposit-tabs" role="tablist">
          <button
            type="button"
            role="tab"
            class="deposit-tabs__item"
            aria-selected="false"
            @click="setTab('pending')"
          >
            待缴保证金
          </button>
          <button
            type="button"
            role="tab"
            class="deposit-tabs__item"
            aria-selected="false"
            @click="setTab('funded')"
          >
            缴存保证金
          </button>
          <button
            type="button"
            role="tab"
            class="deposit-tabs__item is-active"
            aria-selected="true"
            @click="setTab('refunds')"
          >
            退回保证金
          </button>
        </div>
      </template>
      <template #action="{ row }">
        <template v-if="canManage">
          <ElButton link type="primary" @click="remarkRefund(row)">备注</ElButton>
          <ElButton
            v-if="row.status === 'applied'"
            link
            type="primary"
            @click="review(row)"
          >
            审核
          </ElButton>
          <ElButton
            v-else-if="row.status === 'approved'"
            link
            type="primary"
            @click="paid(row)"
          >
            登记打款
          </ElButton>
          <ElButton link type="primary" @click="openLedgers(row.merchant_id)">
            操作记录
          </ElButton>
        </template>
        <template v-else>
          <ElButton link type="primary" @click="openLedgers(row.merchant_id)">
            操作记录
          </ElButton>
        </template>
      </template>
    </RefundGrid>

    <OfflineDrawer>
      <ElForm label-width="140px" class="deposit-action-form">
        <ElFormItem label="商户名称：">
          <span>{{ offlineTarget?.merchant_name || '—' }}</span>
        </ElFormItem>
        <ElFormItem label="商户ID：">
          <span>{{ offlineTarget?.merchant_id || '—' }}</span>
        </ElFormItem>
        <ElFormItem label="商户保证金额度：">
          <span>{{ formatMoney(offlineTarget?.required_amount) }}</span>
        </ElFormItem>
        <ElFormItem label="商户剩余保证金：">
          <span>{{ formatMoney(offlineTarget?.available_amount) }}</span>
        </ElFormItem>
        <ElFormItem label="待缴保证金金额：">
          <span>{{ formatMoney(payableOf(offlineTarget)) }}</span>
        </ElFormItem>
        <ElFormItem label="备注：">
          <ElInput
            v-model="offlineForm.mark"
            type="textarea"
            :rows="3"
            maxlength="500"
            show-word-limit
            placeholder="请输入备注"
          />
        </ElFormItem>
        <ElFormItem label="状态：" required>
          <ElRadioGroup v-model="offlineForm.status">
            <ElRadio :value="0">未缴纳</ElRadio>
            <ElRadio :value="1">已缴纳</ElRadio>
          </ElRadioGroup>
        </ElFormItem>
      </ElForm>
    </OfflineDrawer>

    <DeductDrawer>
      <ElForm label-width="140px" class="deposit-action-form">
        <ElFormItem label="商户名称：">
          <span>{{ deductTarget?.merchant_name || '—' }}</span>
        </ElFormItem>
        <ElFormItem label="商户ID：">
          <span>{{ deductTarget?.merchant_id || '—' }}</span>
        </ElFormItem>
        <ElFormItem label="商户保证金额度：">
          <span>{{ formatMoney(deductTarget?.required_amount) }}</span>
        </ElFormItem>
        <ElFormItem label="商户剩余保证金：">
          <span>{{ formatMoney(deductTarget?.available_amount) }}</span>
        </ElFormItem>
        <ElFormItem label="保证金扣除金额：" required>
          <ElInputNumber
            v-model="deductForm.amount"
            :min="0.01"
            :max="Number(deductTarget?.available_amount || 0)"
            :precision="2"
            :step="0.01"
            :controls="false"
            class="deposit-deduct-amount"
            placeholder="请输入保证金扣除金额"
          />
        </ElFormItem>
        <ElFormItem label="保证金扣除原因：" required>
          <ElInput
            v-model="deductForm.reason"
            maxlength="500"
            show-word-limit
            placeholder="请输入保证金扣除原因"
          />
        </ElFormItem>
      </ElForm>
    </DeductDrawer>

    <LedgerDrawer>
      <div class="deposit-ledger-body">
        <ElTable
          v-loading="ledgerLoading"
          :data="ledgerRows"
          border
          size="small"
          empty-text="暂无数据"
          class="deposit-ledger-table"
        >
          <ElTableColumn
            type="index"
            label="序号"
            min-width="60"
            :index="ledgerIndex"
          />
          <ElTableColumn label="标题" min-width="120">
            <template #default="{ row }">
              {{ ledgerTitle(row.entry_type) }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="金额" min-width="100">
            <template #default="{ row }">
              <span class="deposit-ledger-amount" :class="ledgerSignedAmount(row).className">
                {{ ledgerSignedAmount(row).text }}
              </span>
            </template>
          </ElTableColumn>
          <ElTableColumn label="保证金结余" min-width="110">
            <template #default="{ row }">
              {{ formatMoney(row.balance_after) }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="备注" min-width="180" show-overflow-tooltip>
            <template #default="{ row }">
              {{ ledgerRemark(row) }}
            </template>
          </ElTableColumn>
          <ElTableColumn label="操作时间" min-width="160">
            <template #default="{ row }">
              {{ formatTime(row.created_at) }}
            </template>
          </ElTableColumn>
        </ElTable>
        <div class="deposit-ledger-pagination">
          <ElPagination
            background
            layout="total, prev, pager, next, jumper"
            :current-page="ledgerPage"
            :page-size="ledgerLimit"
            :total="ledgerTotal"
            @current-change="onLedgerPageChange"
          />
        </div>
      </div>
    </LedgerDrawer>
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

.deposit-action-form :deep(.el-form-item) {
  margin-bottom: 16px;
}

.deposit-deduct-amount {
  width: 100%;
}

.deposit-ledger-body {
  min-height: 360px;
  display: flex;
  flex-direction: column;
}

.deposit-ledger-table {
  flex: 1;
  width: 100%;
}

.deposit-ledger-amount.is-fund {
  color: #19be6b;
}

.deposit-ledger-amount.is-deduct {
  color: #ed4014;
}

.deposit-ledger-pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
