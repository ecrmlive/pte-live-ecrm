<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

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

const loading = ref(false);
const rows = ref<MerchantSettlementRow[]>([]);
const summary = ref<MerchantSettlementSummary[]>([]);
const total = ref(0);
const canRead = ref(false);
const canReview = ref(false);
const query = reactive({ limit: 20, merchant_id: undefined as number | undefined, page: 1, status: undefined as MerchantSettlementStatus | undefined });
const summaryByStatus = computed(() => new Map(summary.value.map((item) => [item.status, item])));

const statusLabels: Record<MerchantSettlementStatus, string> = {
  approved: '已审核',
  bill_frozen: '账期已冻结',
  bill_pending: '账期待生成',
  cancelled: '已撤销（历史）',
  paid: '已打款',
  rejected: '已拒绝',
  withdraw_applied: '待平台审核',
};

function summaryText(status: MerchantSettlementStatus) {
  const item = summaryByStatus.value.get(status);
  return item ? `${item.count} 笔 · ${item.amount.toFixed(2)}` : '暂无记录';
}

async function load() {
  if (!canRead.value) return;
  loading.value = true;
  try {
    const [pageData, summaryData] = await Promise.all([
      listPlatformMerchantSettlementsApi(query),
      getPlatformMerchantSettlementSummaryApi(),
    ]);
    rows.value = pageData.list || [];
    total.value = pageData.total || 0;
    summary.value = summaryData.list || [];
  } finally {
    loading.value = false;
  }
}

function search() { query.page = 1; void load(); }
function reset() { query.merchant_id = undefined; query.status = undefined; query.page = 1; void load(); }

function idempotencyKey(action: string, settlementId: number) {
  return `${action}-${settlementId}-${crypto.randomUUID()}`;
}

async function approve(row: MerchantSettlementRow) {
  try {
    await ElMessageBox.confirm('确认审核通过该店铺结算申请？此操作不会直接访问商户库，状态将通过受控命令与事件投影更新。', '结算审核', { type: 'warning' });
    await approvePlatformMerchantSettlementApi(row.settlement_id, { idempotency_key: idempotencyKey('approve', row.settlement_id) });
    ElMessage.success('审核命令已完成，监管投影正在刷新');
    window.setTimeout(() => void load(), 2300);
  } catch { /* 用户取消或接口已返回错误时，requestClient 统一提示。 */ }
}

async function reject(row: MerchantSettlementRow) {
  try {
    const { value } = await ElMessageBox.prompt('请填写拒绝原因', '拒绝结算申请', { inputPattern: /\S+/, inputErrorMessage: '拒绝原因必填' });
    await rejectPlatformMerchantSettlementApi(row.settlement_id, { idempotency_key: idempotencyKey('reject', row.settlement_id), review_note: value.trim() });
    ElMessage.success('拒绝命令已完成，监管投影正在刷新');
    window.setTimeout(() => void load(), 2300);
  } catch { /* 用户取消或接口已返回错误时，requestClient 统一提示。 */ }
}

async function markPaid(row: MerchantSettlementRow) {
  try {
    const { value } = await ElMessageBox.prompt('请输入内部打款凭证编号（不录入银行卡、账号或密钥）', '登记打款凭证', { inputPattern: /\S{3,}/, inputErrorMessage: '凭证编号至少 3 个字符' });
    await markPlatformMerchantSettlementPaidApi(row.settlement_id, { idempotency_key: idempotencyKey('paid', row.settlement_id), payout_reference: value.trim() });
    ElMessage.success('打款登记命令已完成，监管投影正在刷新');
    window.setTimeout(() => void load(), 2300);
  } catch { /* 用户取消或接口已返回错误时，requestClient 统一提示。 */ }
}

onMounted(async () => {
  const [profile, permissions] = await Promise.all([getUserInfoApi(), getAccessCodesApi()]);
  canRead.value = profile.roles.some((role) => role === 'platform' || role === 'region') && permissions.includes('accounts.merchant_settlement.read');
  canReview.value = profile.roles.includes('platform') && permissions.includes('accounts.merchant_settlement.review');
  await load();
});
</script>

<template>
  <Page title="店铺结算监管" description="平台可读取全量投影，区域角色仅可读取授权区域。平台审核与内部凭证登记仅发送受控命令；不直连商户库、不发起转账，也不录入收款账户或密钥。">
    <el-alert v-if="!canRead" class="mb-4" title="当前账号没有查看店铺结算监管投影的权限" type="warning" :closable="false" />
    <template v-else>
      <el-row :gutter="16" class="mb-4">
        <el-col v-for="status in (['bill_frozen', 'withdraw_applied', 'paid'] as MerchantSettlementStatus[])" :key="status" :md="8" :xs="24">
          <el-card shadow="never"><div class="text-sm text-gray-500">{{ statusLabels[status] }}</div><div class="mt-2 text-sm">{{ summaryText(status) }}</div></el-card>
        </el-col>
      </el-row>
      <el-card shadow="never">
        <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
          <el-form-item label="商户 ID"><el-input-number v-model="query.merchant_id" :min="1" controls-position="right" /></el-form-item>
          <el-form-item label="结算状态"><el-select v-model="query.status" clearable class="w-32" placeholder="全部"><el-option v-for="(label, status) in statusLabels" :key="status" :label="label" :value="status" /></el-select></el-form-item>
          <el-form-item><el-button type="primary" @click="search">查询</el-button><el-button @click="reset">重置</el-button></el-form-item>
        </el-form>
        <el-table v-loading="loading" :data="rows" class="mt-3" row-key="settlement_id"><el-table-column label="结算 ID" prop="settlement_id" width="110" /><el-table-column label="商户" min-width="160"><template #default="{ row }">{{ row.merchant_name }}（{{ row.merchant_id }}）</template></el-table-column><el-table-column label="店铺 ID" prop="store_id" width="100" /><el-table-column label="结算周期" min-width="300"><template #default="{ row }">{{ row.period_start }} 至 {{ row.period_end }}</template></el-table-column><el-table-column label="结算金额" width="130"><template #default="{ row }">{{ row.amount.toFixed(2) }}</template></el-table-column><el-table-column label="状态" width="130"><template #default="{ row }"><el-tag>{{ statusLabels[row.status] }}</el-tag></template></el-table-column><el-table-column label="投影更新时间" prop="updated_at" min-width="180" /><el-table-column v-if="canReview" fixed="right" label="操作" width="170"><template #default="{ row }"><template v-if="row.status === 'withdraw_applied'"><el-button link type="success" @click="approve(row)">通过</el-button><el-button link type="danger" @click="reject(row)">拒绝</el-button></template><el-button v-else-if="row.status === 'approved'" link type="primary" @click="markPaid(row)">登记打款</el-button></template></el-table-column></el-table>
        <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page: number) => { query.page = page; load(); }" @size-change="(limit: number) => { query.limit = limit; query.page = 1; load(); }" /></div>
      </el-card>
    </template>
  </Page>
</template>
