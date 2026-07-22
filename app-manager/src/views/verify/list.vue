<template>
  <div>
    <a-card :bordered="false" class="card">
      <a-input-search
        v-model:value="code"
        enter-button="查码"
        placeholder="输入核销码"
        @search="onSearchCode"
      />
    </a-card>
    <a-card :bordered="false" class="card" title="待核销订单">
      <a-list :loading="loading" :data-source="list" item-layout="horizontal">
        <template #renderItem="{ item }">
          <a-list-item>
            <a-list-item-meta
              :title="item.order_sn"
              :description="`${item.real_name} ${item.user_phone} · ¥${Number(item.pay_price).toFixed(2)}`"
            />
            <template #actions>
              <a @click="$router.push(`/verify/${item.order_id}`)">核销</a>
            </template>
          </a-list-item>
        </template>
      </a-list>
      <a-empty v-if="!loading && !list.length" description="暂无待核销" />
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { message } from 'ant-design-vue';
import { fetchAwaitVerify, fetchOrderByCode, type StoreOrder } from '@/api/order';

const router = useRouter();
const loading = ref(false);
const list = ref<StoreOrder[]>([]);
const code = ref('');

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchAwaitVerify();
    list.value = data.list || [];
  } finally {
    loading.value = false;
  }
}

async function onSearchCode() {
  const c = code.value.trim();
  if (!c) return;
  try {
    const { data } = await fetchOrderByCode(c);
    router.push(`/verify/${data.order_id}`);
  } catch {
    message.error('未找到核销码对应订单');
  }
}

onMounted(load);
</script>

<style scoped>
.card {
  margin-bottom: 12px;
  border-radius: 12px;
}
</style>
