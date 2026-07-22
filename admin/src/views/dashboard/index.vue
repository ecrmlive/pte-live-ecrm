<template>
  <div class="dash">
    <a-card :bordered="false" class="card">
      <h2>平台工作台</h2>
      <p>阶段 2：商户审核 / 类目品牌 / 商品审核已接通（见左侧菜单）。</p>
      <a-space style="margin-top: 16px" wrap>
        <a-button type="primary" @click="$router.push('/merchant/audit')">入驻审核</a-button>
        <a-button @click="$router.push('/merchant/list')">商户列表</a-button>
        <a-button @click="$router.push('/product/audit')">商品审核</a-button>
        <a-button @click="$router.push('/product/category')">平台分类</a-button>
        <a-button @click="$router.push('/product/brand')">品牌管理</a-button>
      </a-space>
      <a-descriptions bordered size="small" :column="1" style="max-width: 520px; margin-top: 20px">
        <a-descriptions-item label="账号">{{ auth.user?.account }}</a-descriptions-item>
        <a-descriptions-item label="姓名">{{ auth.user?.real_name }}</a-descriptions-item>
        <a-descriptions-item label="角色">{{ auth.user?.roles }}</a-descriptions-item>
        <a-descriptions-item label="菜单数">{{ menuCount }}</a-descriptions-item>
      </a-descriptions>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useAuthStore } from '@/stores/auth';
import type { MenuNode } from '@/api/auth';

const auth = useAuthStore();

function countMenus(nodes: MenuNode[]): number {
  return nodes.reduce((sum, n) => sum + 1 + countMenus(n.children || []), 0);
}

const menuCount = computed(() => countMenus(auth.menus));
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
