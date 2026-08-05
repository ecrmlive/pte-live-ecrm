<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';
import { deductMerchantDeposit, fetchMerchantDepositRefunds, fetchMerchantDeposits, markMerchantDepositRefundPaid, reviewMerchantDepositRefund, type MerchantDepositAccount, type MerchantDepositRefund } from '#/api/core/ecrm';
import { getAccessCodesApi } from '#/api/core/auth';
const accounts = ref<MerchantDepositAccount[]>([]);
const refunds = ref<MerchantDepositRefund[]>([]);
const canManage = ref(false);
const loading = ref(false);

function splitFields(value: string, count: number) {
  const fields = value.split(',').map((item) => item.trim());
  return fields.length === count ? fields : [];
}

function validateDeduction(value: string) {
  const [amount, reason, key] = splitFields(value, 3);
  const numericAmount = Number(amount);
  const isCentAmount = Number.isFinite(numericAmount) && Math.abs(numericAmount * 100 - Math.round(numericAmount * 100)) < 0.000001;
  if (!amount || !reason || !key || !isCentAmount || numericAmount <= 0) return '请填写精确到分的正数金额、扣减原因和幂等键，并以英文逗号分隔。';
  return true;
}

function validatePayout(value: string) {
  const [key, reference] = splitFields(value, 2);
  if (!key || !reference || [...key].length < 8 || [...key].length > 128 || [...reference].length > 128) return '请填写 8–128 位幂等键和不超过 128 字的内部打款凭证号。';
  return true;
}

function validateNote(value: string) {
  const note = value.trim();
  return note && [...note].length <= 500 ? true : '审核说明不能为空，且不能超过 500 个字符。';
}

function isPromptDismissed(error: unknown) {
  return error === 'cancel' || error === 'close' || error === 'escape';
}

async function load() {
  loading.value = true;
  try {
    const [accountResult, refundResult] = await Promise.all([fetchMerchantDeposits(), fetchMerchantDepositRefunds()]);
    accounts.value = accountResult.list || [];
    refunds.value = refundResult.list || [];
  } finally {
    loading.value = false;
  }
}

async function deduct(row: MerchantDepositAccount) {
  try {
    const { value } = await ElMessageBox.prompt('填写“金额,扣减原因,幂等键”，例如 10,虚构违规扣减,deposit-demo-001。', '扣除保证金', { inputValidator: validateDeduction });
    const [amount, reason, idempotency_key] = splitFields(value, 3);
    await deductMerchantDeposit(row.merchant_id, { amount: Number(amount), reason, idempotency_key });
    ElMessage.success('保证金扣减已登记');
    await load();
  } catch (error) {
    if (!isPromptDismissed(error)) throw error;
  }
}

async function review(row: MerchantDepositRefund, approved: boolean) {
  try {
    const { value } = await ElMessageBox.prompt('填写审核说明。', approved ? '同意退款' : '拒绝退款', { inputValidator: validateNote });
    await reviewMerchantDepositRefund(row.id, approved, value.trim());
    ElMessage.success('退款审核已保存');
    await load();
  } catch (error) {
    if (!isPromptDismissed(error)) throw error;
  }
}

async function paid(row: MerchantDepositRefund) {
  try {
    const { value } = await ElMessageBox.prompt('填写“幂等键,内部打款凭证号”；不得填写账户信息。', '登记保证金退款打款', { inputValidator: validatePayout });
    const [idempotency_key, payout_reference] = splitFields(value, 2);
    await markMerchantDepositRefundPaid(row.id, { idempotency_key, payout_reference });
    ElMessage.success('打款登记已保存');
    await load();
  } catch (error) {
    if (!isPromptDismissed(error)) throw error;
  }
}

onMounted(async () => {
  const [codes] = await Promise.all([getAccessCodesApi(), load()]);
  canManage.value = codes.includes('merchant.deposit.review');
});
</script>
<template><Page title="店铺保证金" description="平台监管保证金余额、不可变流水和退款申请；不展示或保存收款账户资料。"><el-tabs><el-tab-pane label="保证金账户"><el-table v-loading="loading" :data="accounts"><el-table-column prop="merchant_id" label="商户 ID"/><el-table-column label="应缴"><template #default="{row}">¥{{Number(row.required_amount).toFixed(2)}}</template></el-table-column><el-table-column label="可用"><template #default="{row}">¥{{Number(row.available_amount).toFixed(2)}}</template></el-table-column><el-table-column prop="state" label="状态"/><el-table-column v-if="canManage" label="操作"><template #default="{row}"><el-button link type="danger" @click="deduct(row)">扣减</el-button></template></el-table-column></el-table></el-tab-pane><el-tab-pane label="退款申请"><el-table :data="refunds"><el-table-column prop="id" label="申请 ID"/><el-table-column prop="merchant_id" label="商户 ID"/><el-table-column label="退款金额"><template #default="{row}">¥{{Number(row.amount).toFixed(2)}}</template></el-table-column><el-table-column prop="status" label="状态"/><el-table-column prop="reason" label="申请原因"/><el-table-column v-if="canManage" label="操作"><template #default="{row}"><template v-if="row.status==='applied'"><el-button link type="success" @click="review(row,true)">同意</el-button><el-button link type="danger" @click="review(row,false)">拒绝</el-button></template><el-button v-else-if="row.status==='approved'" link type="primary" @click="paid(row)">登记打款</el-button></template></el-table-column></el-table></el-tab-pane></el-tabs></Page></template>
