<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { listFinanceStatementsApi, type FinanceStatement } from '#/api/core/merchant-ledger';

const loading = ref(false);
const rows = ref<FinanceStatement[]>([]);
const total = ref(0);
const query = reactive({ page: 1, limit: 20 });

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
  <Page title="对账单" description="按账期查看本店结算对账单。">
    <el-card shadow="never">
      <el-table v-loading="loading" :data="rows">
        <el-table-column prop="statement_id" label="账单 ID" width="100" />
        <el-table-column label="账期" min-width="200"><template #default="{ row }">{{ row.period_start }} ~ {{ row.period_end }}</template></el-table-column>
        <el-table-column label="金额" width="120"><template #default="{ row }">¥{{ Number(row.amount).toFixed(2) }}</template></el-table-column>
        <el-table-column prop="status" label="状态" width="120" />
        <el-table-column prop="updated_at" label="更新时间" width="180" />
      </el-table>
      <div class="mt-4 flex justify-end"><el-pagination :current-page="query.page" :page-size="query.limit" :total="total" layout="total,prev,pager,next" @current-change="(page) => { query.page = page; load(); }" /></div>
    </el-card>
  </Page>
</template>
