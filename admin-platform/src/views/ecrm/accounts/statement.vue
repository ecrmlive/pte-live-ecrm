<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';

import { getAccessCodesApi } from '#/api/core/auth';
import { getPlatformUserAssetSummaryApi } from '#/api/core/platform-maintain';
import { EcrmListPage } from '#/components/ecrm';

const loading = ref(false);
const canRead = ref(false);
const summary = ref<Array<{ asset_type: string; count: number; expense: number; income: number }>>([]);

const labels: Record<string, string> = { balance: '余额', commission: '佣金', points: '积分' };

async function load() {
  if (!canRead.value) return;
  loading.value = true;
  try {
    summary.value = (await getPlatformUserAssetSummaryApi()).list || [];
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  const permissions = await getAccessCodesApi();
  canRead.value = permissions.includes('accounts.user_assets.read');
  await load();
});
</script>

<template>
  <Page
    title="资金账单"
    description="基于用户资产流水汇总的只读监管视图；明细请查看「用户资产流水」。"
  >
    <el-alert
      class="mb-4"
      title="不含用户姓名、手机号、收款账户或外部支付凭据。"
      type="warning"
      :closable="false"
    />
    <el-alert
      v-if="!canRead"
      class="mb-4"
      title="当前账号无查看资产汇总权限。"
      type="info"
      :closable="false"
    />
    <EcrmListPage v-else title="资产汇总">
      <template #filters>
        <el-button @click="load">刷新</el-button>
      </template>
      <el-table v-loading="loading" :data="summary" row-key="asset_type">
        <el-table-column label="资产类型" min-width="120">
          <template #default="{ row }">{{ labels[row.asset_type] || row.asset_type }}</template>
        </el-table-column>
        <el-table-column label="入账" width="140">
          <template #default="{ row }">{{ row.income.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="支出" width="140">
          <template #default="{ row }">{{ row.expense.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="笔数" width="100" prop="count" />
      </el-table>
      <el-empty v-if="!loading && !summary.length" description="暂无资产汇总数据" />
    </EcrmListPage>
  </Page>
</template>
