<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { Page } from '@vben/common-ui';

import { getMerchantBalanceApi } from '#/api/core/merchant-finance';

const router = useRouter();
const loading = ref(false);
const balance = ref(0);

async function load() {
  loading.value = true;
  try {
    balance.value = (await getMerchantBalanceApi()).mer_money;
  } finally {
    loading.value = false;
  }
}

onMounted(() => void load());
</script>

<template>
  <Page title="店铺余额" description="余额由订单结算、退款与提现申请实时变更；提现审核和打款仅由平台侧处理。">
    <el-card v-loading="loading" class="max-w-2xl" shadow="never">
      <div class="text-sm text-muted-foreground">可用余额（元）</div>
      <div class="my-4 text-4xl font-semibold text-primary">¥{{ Number(balance).toFixed(2) }}</div>
      <div class="text-sm text-muted-foreground">提交提现申请后会立即冻结相应余额；平台拒绝时自动退回。</div>
      <div class="mt-6"><el-button type="primary" @click="router.push('/finance/withdraw')">申请提现</el-button><el-button @click="load">刷新余额</el-button></div>
    </el-card>
  </Page>
</template>
