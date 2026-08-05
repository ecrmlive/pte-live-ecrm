<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { Page } from '@vben/common-ui';
import { getAccessCodesApi } from '#/api/core/auth';
import { getPlatformUserAssetSummaryApi } from '#/api/core/platform-maintain';

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
  <Page title="资金账单" description="基于用户资产流水汇总的只读监管视图；完整对账单待 finance 域 API 接入。">
    <el-alert class="mb-4" title="不含用户姓名、手机号、收款账户或外部支付凭据。" type="warning" :closable="false" />
    <el-alert v-if="!canRead" class="mb-4" title="当前账号无查看资产汇总权限。" type="info" :closable="false" />
    <el-row v-else v-loading="loading" :gutter="16">
      <el-col v-for="item in summary" :key="item.asset_type" :md="8" :xs="24" class="mb-4">
        <el-card shadow="never">
          <div class="text-sm text-gray-500">{{ labels[item.asset_type] || item.asset_type }}</div>
          <div class="mt-2 text-sm">入账 {{ item.income.toFixed(2) }} · 支出 {{ item.expense.toFixed(2) }}</div>
          <div class="mt-1 text-lg font-semibold">{{ item.count }} 笔</div>
        </el-card>
      </el-col>
      <el-col v-if="!loading && !summary.length" :span="24"><el-empty description="暂无资产汇总数据" /></el-col>
    </el-row>
  </Page>
</template>
