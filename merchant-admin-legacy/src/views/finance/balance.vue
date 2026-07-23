<template>
  <a-card :bordered="false" class="page-card">
    <a-statistic title="店铺可提现余额" :precision="2" :value="balance" prefix="¥" />
    <a-button type="primary" style="margin-top: 24px" @click="$router.push('/finance/withdraw')">
      去提现
    </a-button>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { fetchBalance } from '@/api/finance';

const balance = ref(0);

onMounted(async () => {
  const { data } = await fetchBalance();
  balance.value = Number(data.mer_money || 0);
});
</script>

<style scoped>
.page-card {
  border-radius: 14px;
  max-width: 480px;
}
</style>
