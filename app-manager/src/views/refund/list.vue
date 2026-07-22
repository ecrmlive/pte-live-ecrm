<template>
  <a-card :bordered="false" class="card" title="本店售后">
    <a-table
      row-key="refund_order_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="false"
      size="small"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'price'">¥{{ Number(record.refund_price).toFixed(2) }}</template>
        <template v-else-if="column.key === 'status'">{{ statusText(record.status) }}</template>
        <template v-else-if="column.key === 'action'">
          <a-button
            v-if="record.status === 0"
            type="link"
            @click="onApprove(record.refund_order_id)"
          >
            同意退款
          </a-button>
        </template>
      </template>
    </a-table>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { message } from 'ant-design-vue';
import { approveRefund, fetchRefunds } from '@/api/refund';

const loading = ref(false);
const list = ref<Array<Record<string, unknown>>>([]);
const columns = [
  { title: 'ID', dataIndex: 'refund_order_id', width: 70 },
  { title: '单号', dataIndex: 'refund_order_sn' },
  { title: '金额', key: 'price', width: 90 },
  { title: '状态', key: 'status', width: 90 },
  { title: '操作', key: 'action', width: 100 },
];

function statusText(s: number) {
  if (s === 0) return '待审';
  if (s === 3) return '已退款';
  if (s === -1) return '拒绝';
  return String(s);
}

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchRefunds({ page: 1, limit: 50 });
    list.value = (data as { list?: Array<Record<string, unknown>> }).list || [];
  } finally {
    loading.value = false;
  }
}

async function onApprove(id: number) {
  await approveRefund(id);
  message.success('已同意退款');
  load();
}

onMounted(load);
</script>

<style scoped>
.card {
  border-radius: 12px;
}
</style>
