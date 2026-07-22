<template>
  <div class="dash">
    <a-card :bordered="false" class="card">
      <h2>店铺工作台</h2>
      <p>阶段 2：商品 CRUD / 上下架 / 库存已接通；平台审品后 C 端可见。</p>
      <a-space style="margin-top: 16px" wrap>
        <a-button type="primary" @click="$router.push('/product/list')">商品列表</a-button>
        <a-button @click="$router.push('/product/edit')">发布商品</a-button>
      </a-space>
      <a-descriptions bordered size="small" :column="1" style="max-width: 520px; margin-top: 20px">
        <a-descriptions-item label="商户">{{ auth.user?.mer_name }}</a-descriptions-item>
        <a-descriptions-item label="mer_id">{{ auth.user?.mer_id }}</a-descriptions-item>
        <a-descriptions-item label="账号">{{ auth.user?.account }}</a-descriptions-item>
      </a-descriptions>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
</script>

<style scoped>
.card {
  border-radius: 14px;
  background: var(--qx-panel);
}
.card h2 {
  margin: 0 0 8px;
}
.card p {
  margin: 0;
  color: #516070;
}
</style>
