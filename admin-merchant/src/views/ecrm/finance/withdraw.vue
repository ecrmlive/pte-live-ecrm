<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage } from 'element-plus';

import {
  applyMerchantWithdrawApi,
  getMerchantBalanceApi,
  getMerchantWithdrawApi,
  listMerchantWithdrawsApi,
  type MerchantWithdraw,
} from '#/api/core/merchant-finance';

const loading = ref(false);
const saving = ref(false);
const rows = ref<MerchantWithdraw[]>([]);
const total = ref(0);
const balance = ref(0);
const applyOpen = ref(false);
const detailOpen = ref(false);
const detail = ref<MerchantWithdraw>();
const query = reactive({ limit: 20, page: 1 });
const form = reactive({ extract_money: 0, financial_account: '', financial_type: 1, mark: '' });

const auditStatus: Record<number, { label: string; type: 'danger' | 'info' | 'success' | 'warning' }> = {
  [-1]: { label: '审核拒绝', type: 'danger' },
  0: { label: '待平台审核', type: 'warning' },
  1: { label: '审核通过', type: 'success' },
};

function auditInfo(status: number) {
  return auditStatus[status] || { label: '未知状态', type: 'info' as const };
}

function accountType(type: number) {
  return ({ 1: '银行卡', 2: '微信', 3: '支付宝' }[type] || '未知');
}

function transferInfo(row: MerchantWithdraw) {
  return row.financial_status === 1 ? { label: '已打款', type: 'success' as const } : { label: '未打款', type: 'info' as const };
}

async function load() {
  loading.value = true;
  try {
    const [page, currentBalance] = await Promise.all([listMerchantWithdrawsApi(query), getMerchantBalanceApi()]);
    rows.value = page.list;
    total.value = page.total;
    balance.value = currentBalance.mer_money;
  } finally {
    loading.value = false;
  }
}

function openApply() {
  Object.assign(form, { extract_money: 0, financial_account: '', financial_type: 1, mark: '' });
  applyOpen.value = true;
}

async function apply() {
  if (form.extract_money <= 0) {
    ElMessage.warning('提现金额必须大于 0');
    return;
  }
  if (form.extract_money > balance.value) {
    ElMessage.warning('提现金额不能超过可用余额');
    return;
  }
  if (!form.financial_account.trim()) {
    ElMessage.warning('请填写收款账户');
    return;
  }
  saving.value = true;
  try {
    await applyMerchantWithdrawApi({ ...form, financial_account: form.financial_account.trim(), mark: form.mark.trim() });
    applyOpen.value = false;
    ElMessage.success('提现申请已提交，等待平台审核');
    await load();
  } finally {
    saving.value = false;
  }
}

async function openDetail(row: MerchantWithdraw) {
  detail.value = await getMerchantWithdrawApi(row.financial_id);
  detailOpen.value = true;
}

onMounted(() => void load());
</script>

<template>
  <Page title="提现管理" description="仅可提交提现申请和查看本商户记录；审核、拒绝、打款及凭证均由平台后台处理。">
    <el-card shadow="never"><div class="flex items-center justify-between"><div><span class="text-sm text-muted-foreground">当前可用余额</span><span class="ml-3 text-xl font-semibold text-primary">¥{{ Number(balance).toFixed(2) }}</span></div><el-button type="primary" @click="openApply">申请提现</el-button></div></el-card>
    <el-card class="mt-4" shadow="never"><el-table v-loading="loading" :data="rows" row-key="financial_id"><el-table-column label="申请单号" min-width="180" prop="financial_sn" /><el-table-column label="提现金额" width="116"><template #default="{ row }">¥{{ Number(row.extract_money).toFixed(2) }}</template></el-table-column><el-table-column label="收款方式" width="104"><template #default="{ row }">{{ accountType(row.financial_type) }}</template></el-table-column><el-table-column label="收款账户" min-width="150" prop="financial_account" show-overflow-tooltip /><el-table-column label="审核状态" width="116"><template #default="{ row }"><el-tag :type="auditInfo(row.status).type">{{ auditInfo(row.status).label }}</el-tag></template></el-table-column><el-table-column label="打款状态" width="100"><template #default="{ row }"><el-tag :type="transferInfo(row).type">{{ transferInfo(row).label }}</el-tag></template></el-table-column><el-table-column label="申请时间" min-width="170" prop="create_time" /><el-table-column fixed="right" label="操作" width="76"><template #default="{ row }"><el-button link type="primary" @click="openDetail(row)">详情</el-button></template></el-table-column></el-table><div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page) => { query.page = page; load(); }" @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }" /></div></el-card>
    <el-dialog v-model="applyOpen" title="申请提现" width="500px" destroy-on-close><el-form label-width="96px"><el-form-item label="可用余额"><span>¥{{ Number(balance).toFixed(2) }}</span></el-form-item><el-form-item label="提现金额" required><el-input-number v-model="form.extract_money" :min="0.01" :precision="2" class="w-full" /></el-form-item><el-form-item label="收款方式" required><el-select v-model="form.financial_type" class="w-full"><el-option label="银行卡" :value="1" /><el-option label="微信" :value="2" /><el-option label="支付宝" :value="3" /></el-select></el-form-item><el-form-item label="收款账户" required><el-input v-model="form.financial_account" placeholder="银行卡号、微信号或支付宝账号" /></el-form-item><el-form-item label="申请备注"><el-input v-model="form.mark" :rows="3" maxlength="200" show-word-limit type="textarea" /></el-form-item></el-form><template #footer><el-button @click="applyOpen = false">取消</el-button><el-button :loading="saving" type="primary" @click="apply">提交申请</el-button></template></el-dialog>
    <el-drawer v-model="detailOpen" :with-header="false" size="560px"><template v-if="detail"><div class="mb-5 text-lg font-medium">提现详情</div><el-descriptions :column="1" border><el-descriptions-item label="申请单号">{{ detail.financial_sn }}</el-descriptions-item><el-descriptions-item label="提现金额">¥{{ Number(detail.extract_money).toFixed(2) }}</el-descriptions-item><el-descriptions-item label="收款方式">{{ accountType(detail.financial_type) }}</el-descriptions-item><el-descriptions-item label="收款账户">{{ detail.financial_account }}</el-descriptions-item><el-descriptions-item label="审核状态"><el-tag :type="auditInfo(detail.status).type">{{ auditInfo(detail.status).label }}</el-tag></el-descriptions-item><el-descriptions-item label="打款状态"><el-tag :type="transferInfo(detail).type">{{ transferInfo(detail).label }}</el-tag></el-descriptions-item><el-descriptions-item v-if="detail.refusal" label="拒绝原因">{{ detail.refusal }}</el-descriptions-item><el-descriptions-item label="申请备注">{{ detail.mark || '—' }}</el-descriptions-item><el-descriptions-item label="申请时间">{{ detail.create_time || '—' }}</el-descriptions-item></el-descriptions></template></el-drawer>
  </Page>
</template>
