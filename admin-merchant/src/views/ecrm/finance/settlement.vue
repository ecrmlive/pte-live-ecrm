<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { ElMessage, ElMessageBox } from 'element-plus';

import { getAccessCodesApi } from '#/api/core/auth';
import {
  applyMerchantSettlementApi,
  getMerchantSettlementApi,
  listMerchantSettlementsApi,
  type MerchantSettlement,
  type MerchantSettlementStatus,
} from '#/api/core/merchant-finance';

const loading = ref(false);
const applying = ref(false);
const rows = ref<MerchantSettlement[]>([]);
const total = ref(0);
const detailOpen = ref(false);
const detail = ref<MerchantSettlement>();
const canApply = ref(false);
const query = reactive({
  limit: 20,
  page: 1,
  status: undefined as MerchantSettlementStatus | undefined,
});

const statusLabels: Record<MerchantSettlementStatus, string> = {
  approved: '已审核',
  bill_frozen: '账期已冻结',
  bill_pending: '账期待生成',
  paid: '已打款',
  rejected: '已拒绝',
  withdraw_applied: '待平台审核',
};

const statusTypes: Record<MerchantSettlementStatus, 'danger' | 'info' | 'success' | 'warning'> = {
  approved: 'success',
  bill_frozen: 'warning',
  bill_pending: 'info',
  paid: 'success',
  rejected: 'danger',
  withdraw_applied: 'warning',
};

function formatTime(value?: string | null) {
  if (!value) return '—';
  return String(value).replace('T', ' ').slice(0, 19);
}

function canApplyRow(row: MerchantSettlement) {
  return canApply.value && row.status === 'bill_frozen' && row.amount > 0;
}

function idempotencyKey(settlementId: number) {
  return `merchant-apply-${settlementId}-${crypto.randomUUID()}`;
}

async function load() {
  loading.value = true;
  try {
    const page = await listMerchantSettlementsApi(query);
    rows.value = page.list;
    total.value = page.total;
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

async function openDetail(row: MerchantSettlement) {
  detail.value = await getMerchantSettlementApi(row.settlement_id);
  detailOpen.value = true;
}

async function apply(row: MerchantSettlement) {
  try {
    await ElMessageBox.confirm(
      `确认提交结算申请？账期 ${formatTime(row.period_start)} 至 ${formatTime(row.period_end)}，金额 ¥${Number(row.amount).toFixed(2)}。提交后由平台审核打款，本页不录入收款账户。`,
      '提交结算申请',
      { confirmButtonText: '确认提交', cancelButtonText: '取消', type: 'warning' },
    );
  } catch {
    return;
  }
  applying.value = true;
  try {
    await applyMerchantSettlementApi(row.settlement_id, {
      idempotency_key: idempotencyKey(row.settlement_id),
    });
    ElMessage.success('结算申请已提交，等待平台审核');
    await load();
  } finally {
    applying.value = false;
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canApply.value = permissions.includes('finance.settlement.apply');
  await load();
});
</script>

<template>
  <Page title="结算管理" description="查看本店结算账期与申请进度；仅账期已冻结且金额大于 0 时可提交结算申请，审核与打款由平台处理。">
    <el-card shadow="never">
      <el-form class="flex flex-wrap gap-x-4" label-width="72px" @submit.prevent="search">
        <el-form-item label="结算状态">
          <el-select v-model="query.status" clearable class="w-40" placeholder="全部">
            <el-option v-for="(label, status) in statusLabels" :key="status" :label="label" :value="status" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">查询</el-button>
          <el-button @click="reset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="mt-4" shadow="never">
      <el-table v-loading="loading" :data="rows" row-key="settlement_id">
        <el-table-column label="结算 ID" prop="settlement_id" width="100" />
        <el-table-column label="结算周期" min-width="280">
          <template #default="{ row }">{{ formatTime(row.period_start) }} 至 {{ formatTime(row.period_end) }}</template>
        </el-table-column>
        <el-table-column label="结算金额" width="120">
          <template #default="{ row }">¥{{ Number(row.amount).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="statusTypes[row.status]">{{ statusLabels[row.status] }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="申请单号" min-width="140" prop="application_no" show-overflow-tooltip />
        <el-table-column label="更新时间" min-width="170">
          <template #default="{ row }">{{ formatTime(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="148">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">详情</el-button>
            <el-button v-if="canApplyRow(row)" :loading="applying" link type="success" @click="apply(row)">申请结算</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="mt-4 flex justify-end">
        <el-pagination
          :current-page="query.page"
          :page-size="query.limit"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          background
          layout="total, sizes, prev, pager, next"
          @current-change="(page) => { query.page = page; load(); }"
          @size-change="(limit) => { query.limit = limit; query.page = 1; load(); }"
        />
      </div>
    </el-card>

    <el-drawer v-model="detailOpen" :with-header="false" size="560px">
      <template v-if="detail">
        <div class="mb-5 text-lg font-medium">结算详情</div>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="结算 ID">{{ detail.settlement_id }}</el-descriptions-item>
          <el-descriptions-item label="结算周期">{{ formatTime(detail.period_start) }} 至 {{ formatTime(detail.period_end) }}</el-descriptions-item>
          <el-descriptions-item label="结算金额">¥{{ Number(detail.amount).toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="statusTypes[detail.status]">{{ statusLabels[detail.status] }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="申请单号">{{ detail.application_no || '—' }}</el-descriptions-item>
          <el-descriptions-item label="审核备注">{{ detail.review_note || '—' }}</el-descriptions-item>
          <el-descriptions-item label="打款时间">{{ formatTime(detail.paid_at) }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ formatTime(detail.updated_at) }}</el-descriptions-item>
        </el-descriptions>
        <div v-if="canApplyRow(detail)" class="mt-6">
          <el-button :loading="applying" type="primary" @click="apply(detail)">提交结算申请</el-button>
        </div>
      </template>
    </el-drawer>
  </Page>
</template>
