<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getUserInfoApi } from '#/api/core/auth';
import {
  approvePlatformWithdrawApi,
  getPlatformWithdrawApi,
  listPlatformWithdrawsApi,
  rejectPlatformWithdrawApi,
  type PlatformWithdraw,
} from '#/api/core/platform-finance';

const loading = ref(false);
const rows = ref<PlatformWithdraw[]>([]);
const total = ref(0);
const detailOpen = ref(false);
const rejectOpen = ref(false);
const current = ref<PlatformWithdraw>();
const rejecting = ref(false);
const isPlatformOperator = ref(false);
const rejectForm = reactive({ refusal: '' });
const query = reactive({ limit: 20, page: 1, status: undefined as number | undefined });

const auditStatus: Record<number, { label: string; type: 'danger' | 'info' | 'success' | 'warning' }> = {
  [-1]: { label: '审核拒绝', type: 'danger' },
  0: { label: '待平台审核', type: 'warning' },
  1: { label: '审核通过', type: 'success' },
};

function auditInfo(status: number) {
  return auditStatus[status] || { label: '未知状态', type: 'info' as const };
}

function accountType(type: number) {
  return ({ 1: '银行卡', 2: '微信', 3: '支付宝' } as Record<number, string>)[type] || '未知';
}

function transferInfo(row: PlatformWithdraw) {
  return row.financial_status === 1
    ? { label: '已打款', type: 'success' as const }
    : { label: '未打款', type: 'info' as const };
}

function canAudit(row: PlatformWithdraw) {
  return isPlatformOperator.value && row.status === 0;
}

async function load() {
  loading.value = true;
  try {
    const result = await listPlatformWithdrawsApi(query);
    rows.value = result.list || [];
    total.value = result.total || 0;
  } finally {
    loading.value = false;
  }
}

function search() {
  query.page = 1;
  void load();
}

function reset() {
  query.status = undefined;
  query.page = 1;
  void load();
}

async function openDetail(row: PlatformWithdraw) {
  current.value = await getPlatformWithdrawApi(row.financial_id);
  detailOpen.value = true;
}

async function approve(row: PlatformWithdraw) {
  try {
    await ElMessageBox.confirm(
      '确认审核通过该提现申请？该操作只更新审核状态，不会执行实际打款；打款凭证功能尚未接入。',
      '审核通过确认',
      { confirmButtonText: '确认通过', cancelButtonText: '取消', type: 'warning' },
    );
    await approvePlatformWithdrawApi(row.financial_id);
    ElMessage.success('提现申请已审核通过');
    await load();
  } catch {
    // 用户取消或接口已返回错误时，requestClient 统一处理提示。
  }
}

function openReject(row: PlatformWithdraw) {
  current.value = row;
  rejectForm.refusal = '';
  rejectOpen.value = true;
}

async function reject() {
  const refusal = rejectForm.refusal.trim();
  if (!refusal) {
    ElMessage.warning('请填写拒绝原因');
    return;
  }
  if (!current.value) return;
  rejecting.value = true;
  try {
    await rejectPlatformWithdrawApi(current.value.financial_id, refusal);
    rejectOpen.value = false;
    ElMessage.success('提现申请已拒绝，金额已退回商户余额');
    await load();
  } finally {
    rejecting.value = false;
  }
}

onMounted(async () => {
  const [profile] = await Promise.all([getUserInfoApi(), load()]);
  isPlatformOperator.value = profile.is_agent !== 1;
});
</script>

<template>
  <Page title="提现审核" description="平台账号可审核商户提现；区域账号仅查看所属区域商户记录。审核通过不等于实际打款，当前接口未提供打款凭证上传。">
    <el-card shadow="never">
      <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
        <el-form-item label="审核状态">
          <el-select v-model="query.status" clearable class="w-44" placeholder="全部状态">
            <el-option label="待平台审核" :value="0" />
            <el-option label="审核通过" :value="1" />
            <el-option label="审核拒绝" :value="-1" />
          </el-select>
        </el-form-item>
        <el-form-item><el-button type="primary" @click="search">查询</el-button><el-button @click="reset">重置</el-button></el-form-item>
      </el-form>
    </el-card>

    <el-card class="mt-4" shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="financial_id">
        <el-table-column label="申请单号" min-width="180" prop="financial_sn" />
        <el-table-column label="商户 ID" width="100" prop="mer_id" />
        <el-table-column label="提现金额" width="116"><template #default="{ row }">¥{{ Number(row.extract_money).toFixed(2) }}</template></el-table-column>
        <el-table-column label="收款方式" width="104"><template #default="{ row }">{{ accountType(row.financial_type) }}</template></el-table-column>
        <el-table-column label="收款账户" min-width="150" prop="financial_account" show-overflow-tooltip />
        <el-table-column label="审核状态" width="116"><template #default="{ row }"><el-tag :type="auditInfo(row.status).type">{{ auditInfo(row.status).label }}</el-tag></template></el-table-column>
        <el-table-column label="打款状态" width="100"><template #default="{ row }"><el-tag :type="transferInfo(row).type">{{ transferInfo(row).label }}</el-tag></template></el-table-column>
        <el-table-column label="申请时间" min-width="170" prop="create_time" />
        <el-table-column fixed="right" label="操作" width="172"><template #default="{ row }"><el-button link type="primary" @click="openDetail(row)">详情</el-button><template v-if="canAudit(row)"><el-button link type="success" @click="approve(row)">通过</el-button><el-button link type="danger" @click="openReject(row)">拒绝</el-button></template></template></el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :page-sizes="[10, 20, 50, 100]" :total="total" background layout="total, sizes, prev, pager, next" @current-change="(page: number) => { query.page = page; load(); }" @size-change="(limit: number) => { query.limit = limit; query.page = 1; load(); }" /></div>
    </el-card>

    <el-drawer v-model="detailOpen" :with-header="false" size="560px"><template v-if="current"><div class="mb-5 text-lg font-medium">提现详情</div><el-descriptions :column="1" border><el-descriptions-item label="申请单号">{{ current.financial_sn }}</el-descriptions-item><el-descriptions-item label="商户 ID">{{ current.mer_id }}</el-descriptions-item><el-descriptions-item label="提现金额">¥{{ Number(current.extract_money).toFixed(2) }}</el-descriptions-item><el-descriptions-item label="收款方式">{{ accountType(current.financial_type) }}</el-descriptions-item><el-descriptions-item label="收款账户">{{ current.financial_account }}</el-descriptions-item><el-descriptions-item label="审核状态"><el-tag :type="auditInfo(current.status).type">{{ auditInfo(current.status).label }}</el-tag></el-descriptions-item><el-descriptions-item label="打款状态"><el-tag :type="transferInfo(current).type">{{ transferInfo(current).label }}</el-tag></el-descriptions-item><el-descriptions-item v-if="current.refusal" label="拒绝原因">{{ current.refusal }}</el-descriptions-item><el-descriptions-item label="申请备注">{{ current.mark || '—' }}</el-descriptions-item><el-descriptions-item label="申请时间">{{ current.create_time || '—' }}</el-descriptions-item></el-descriptions></template></el-drawer>

    <el-dialog v-model="rejectOpen" title="拒绝提现" width="480px" destroy-on-close><el-form label-width="84px"><el-form-item label="拒绝原因" required><el-input v-model="rejectForm.refusal" :rows="4" maxlength="200" placeholder="请向商户说明拒绝原因" show-word-limit type="textarea" /></el-form-item></el-form><template #footer><el-button @click="rejectOpen = false">取消</el-button><el-button :loading="rejecting" type="danger" @click="reject">确认拒绝</el-button></template></el-dialog>
  </Page>
</template>
