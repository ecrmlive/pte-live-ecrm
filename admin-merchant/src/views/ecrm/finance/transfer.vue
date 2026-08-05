<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { getMerchantBalanceApi, listMerchantWithdrawsApi, type MerchantWithdraw } from '#/api/core/merchant-finance';
import { EcrmListPage } from '#/components/ecrm';

const router = useRouter();
const balance = ref(0);
const rows = ref<MerchantWithdraw[]>([]);
const loading = ref(false);

const auditStatus: Record<number, string> = {
  [-1]: '审核拒绝',
  0: '待平台审核',
  1: '审核通过',
};

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
  <EcrmListPage
    title="转账管理"
    description="查看余额与近期提现记录；完整提现流程请前往提现管理。本页不录入收款密钥。"
  >
    <template #filters>
      <div class="flex flex-wrap items-center justify-between gap-3 w-full">
        <div>
          <span class="text-sm text-muted-foreground">可用余额（元）</span>
          <span class="ml-3 text-2xl font-semibold">¥{{ Number(balance).toFixed(2) }}</span>
        </div>
        <div class="flex gap-2">
          <el-button @click="load">刷新</el-button>
          <el-button type="primary" @click="router.push('/finance/withdraw')">申请提现</el-button>
        </div>
      </div>
    </template>

    <div class="mb-3 text-sm font-medium">近期提现记录</div>
    <el-table v-loading="loading" :data="rows" row-key="financial_id">
      <el-table-column prop="financial_sn" label="单号" min-width="160" />
      <el-table-column label="金额" width="120">
        <template #default="{ row }">¥{{ Number(row.extract_money).toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="120">
        <template #default="{ row }">{{ auditStatus[row.status] || row.status }}</template>
      </el-table-column>
      <el-table-column prop="create_time" label="申请时间" width="180" />
    </el-table>
    <el-empty v-if="!loading && !rows.length" description="暂无提现记录" />
  </EcrmListPage>
</template>
