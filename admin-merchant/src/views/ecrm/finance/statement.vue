<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';

import { listFinanceStatementsApi, type FinanceStatement } from '#/api/core/merchant-ledger';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const rows = ref<FinanceStatement[]>([]);
const total = ref(0);
const query = reactive({ page: 1, limit: 20 });

const statusLabels: Record<string, string> = {
  approved: '已审核',
  bill_frozen: '账期已冻结',
  bill_pending: '账期待生成',
  paid: '已打款',
  rejected: '已拒绝',
  withdraw_applied: '待平台审核',
};

function statusLabel(status: string) {
  return statusLabels[status] || status;
}

async function load() {
  loading.value = true;
  try {
    const data = await listFinanceStatementsApi({ ...query });
    rows.value = data.list ?? [];
    total.value = data.total ?? 0;
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <EcrmListPage title="对账单" description="按账期查看本店结算对账单；数据来自 qixi_crm_m_settlement_bill。">
    <template #filters>
      <div class="flex flex-wrap items-center gap-2">
        <el-button @click="load">刷新</el-button>
      </div>
    </template>

    <el-table v-loading="loading" :data="rows" row-key="statement_id">
      <el-table-column prop="statement_id" label="账单 ID" width="100" />
      <el-table-column label="账期" min-width="200">
        <template #default="{ row }">{{ row.period_start }} ~ {{ row.period_end }}</template>
      </el-table-column>
      <el-table-column label="金额" width="120">
        <template #default="{ row }">¥{{ Number(row.amount).toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag>{{ statusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="updated_at" label="更新时间" width="180" />
    </el-table>
    <el-empty v-if="!loading && !rows.length" description="暂无对账单" />

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
