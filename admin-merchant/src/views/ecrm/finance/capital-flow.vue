<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { listFinanceLedgerApi, type FinanceLedgerEntry } from '#/api/core/merchant-ledger';

const loading = ref(false);
const rows = ref<FinanceLedgerEntry[]>([]);
const total = ref(0);
const query = reactive({ page: 1, limit: 20 });

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
  <Page title="资金流水" description="本店资金账本流水明细。">
    <el-card shadow="never">
      <el-table v-loading="loading" :data="rows">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="entry_type" label="类型" width="140" />
        <el-table-column label="金额" width="120"><template #default="{ row }">¥{{ Number(row.amount).toFixed(2) }}</template></el-table-column>
        <el-table-column prop="reference_type" label="关联类型" width="120" />
        <el-table-column prop="reference_id" label="关联 ID" min-width="140" />
        <el-table-column prop="created_at" label="时间" width="180" />
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :total="total" layout="total,prev,pager,next" @current-change="(page) => { query.page = page; load(); }" /></div>
    </el-card>
  </Page>
</template>
