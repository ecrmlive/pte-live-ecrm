<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { listFinanceLedgerApi, type FinanceLedgerEntry } from '#/api/core/merchant-ledger';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const rows = ref<FinanceLedgerEntry[]>([]);
const total = ref(0);
const query = reactive({ page: 1, limit: 20 });

const entryLabels: Record<string, string> = {
  test_seed: '夹具入账',
  settlement_accrual: '结算应计',
  settlement_payout: '结算打款',
  withdraw: '提现扣减',
  refund_reversal: '退款冲销',
};

function entryLabel(type: string) {
  return entryLabels[type] || type;
}

async function load() {
  loading.value = true;
  try {
    const data = await listFinanceLedgerApi({ ...query });
    rows.value = data.list ?? [];
    total.value = data.total ?? 0;
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <EcrmListPage title="资金流水" description="本店资金账本流水明细；数据来自 qixi_crm_m_finance_ledger，按 store_id 隔离。">
    <template #filters>
      <div class="flex flex-wrap items-center gap-2">
        <el-button @click="load">刷新</el-button>
      </div>
    </template>

    <el-table v-loading="loading" :data="rows" row-key="id">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column label="类型" width="140">
        <template #default="{ row }">{{ entryLabel(row.entry_type) }}</template>
      </el-table-column>
      <el-table-column label="金额" width="120">
        <template #default="{ row }">¥{{ Number(row.amount).toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="reference_type" label="关联类型" width="120" />
      <el-table-column prop="reference_id" label="关联 ID" min-width="140" />
      <el-table-column prop="created_at" label="时间" width="180" />
    </el-table>
    <el-empty v-if="!loading && !rows.length" description="暂无资金流水" />

    <template #pager>
      <el-pagination
        :current-page="query.page"
        :page-size="query.limit"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="(page: number) => { query.page = page; load(); }"
      />
    </template>
  </EcrmListPage>
</template>
