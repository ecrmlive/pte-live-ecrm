<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { Page } from '@vben/common-ui';
import { getMerchantBalanceApi } from '#/api/core/merchant-finance';
import { listMerchantWithdrawsApi, type MerchantWithdraw } from '#/api/core/merchant-finance';

const router = useRouter();
const balance = ref(0);
const rows = ref<MerchantWithdraw[]>([]);
const loading = ref(false);

async function load() {
  loading.value = true;
  try {
    const [bal, withdraws] = await Promise.all([
      getMerchantBalanceApi(),
      listMerchantWithdrawsApi({ page: 1, limit: 10 }),
    ]);
    balance.value = bal.mer_money;
    rows.value = withdraws.list ?? [];
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="转账管理" description="查看余额与近期提现记录；完整提现流程请前往提现管理。">
    <el-card v-loading="loading" class="mb-4 max-w-xl" shadow="never">
      <div class="text-sm text-muted-foreground">可用余额（元）</div>
      <div class="my-3 text-3xl font-semibold">¥{{ Number(balance).toFixed(2) }}</div>
      <el-button type="primary" @click="router.push('/finance/withdraw')">申请提现</el-button>
    </el-card>
    <el-card shadow="never">
      <template #header>近期提现记录</template>
      <el-table :data="rows">
        <el-table-column prop="financial_sn" label="单号" min-width="160" />
        <el-table-column label="金额" width="120"><template #default="{ row }">¥{{ Number(row.extract_money).toFixed(2) }}</template></el-table-column>
        <el-table-column prop="status" label="状态" width="100" />
        <el-table-column prop="create_time" label="申请时间" width="180" />
      </el-table>
    </el-card>
  </Page>
</template>
