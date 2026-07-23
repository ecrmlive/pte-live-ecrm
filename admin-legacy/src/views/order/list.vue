<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-select v-model:value="paid" allow-clear placeholder="支付状态" style="width: 140px" @change="reload">
        <a-select-option :value="0">未支付</a-select-option>
        <a-select-option :value="1">已支付</a-select-option>
      </a-select>
      <a-button type="primary" @click="load">查询</a-button>
    </div>
    <a-table
      row-key="order_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'price'">¥{{ Number(record.pay_price).toFixed(2) }}</template>
        <template v-else-if="column.key === 'status'">
          <a-tag>{{ statusText(record) }}</a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="openDetail(record)">详情</a-button>
        </template>
      </template>
    </a-table>

    <a-drawer v-model:open="detailOpen" title="订单监管详情" width="480">
      <template v-if="current">
        <p>商户：{{ current.mer_name || current.mer_id }}</p>
        <p>单号：{{ current.order_sn }}</p>
        <p>主单：{{ current.group_order_id }}</p>
        <p>收货人：{{ current.real_name }} {{ current.user_phone }}</p>
        <p>地址：{{ current.user_address }}</p>
        <p>金额：¥{{ Number(current.pay_price).toFixed(2) }}</p>
        <a-divider />
        <div v-for="(p, idx) in current.products || []" :key="idx">
          {{ productName(p) }} ×{{ p.product_num }} · ¥{{ Number(p.total_price).toFixed(2) }}
        </div>
      </template>
    </a-drawer>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { fetchOrder, fetchOrders, type StoreOrder } from '@/api/order';

const loading = ref(false);
const list = ref<StoreOrder[]>([]);
const paid = ref<number | undefined>(1);
const detailOpen = ref(false);
const current = ref<StoreOrder | null>(null);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });

const columns = [
  { title: '订单ID', dataIndex: 'order_id', width: 90 },
  { title: '商户', dataIndex: 'mer_name', width: 140 },
  { title: '单号', dataIndex: 'order_sn' },
  { title: '收货人', dataIndex: 'real_name', width: 100 },
  { title: '金额', key: 'price', width: 100 },
  { title: '状态', key: 'status', width: 120 },
  { title: '操作', key: 'action', width: 100 },
];

function statusText(o: StoreOrder) {
  if (!o.paid) return '待支付';
  if (o.status === 0) return '待发货';
  if (o.status === 1) return '待收货';
  if (o.status === 3) return '已完成';
  return `状态${o.status}`;
}

function productName(p: { product_info?: string; product_id: number }) {
  try {
    return JSON.parse(p.product_info || '{}').store_name || `商品${p.product_id}`;
  } catch {
    return `商品${p.product_id}`;
  }
}

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchOrders({
      page: pagination.current,
      limit: pagination.pageSize,
      paid: paid.value,
    });
    list.value = data.list || [];
    pagination.total = data.total;
  } finally {
    loading.value = false;
  }
}

function reload() {
  pagination.current = 1;
  load();
}

function onTableChange(p: { current?: number; pageSize?: number }) {
  pagination.current = p.current || 1;
  pagination.pageSize = p.pageSize || 20;
  load();
}

async function openDetail(row: StoreOrder) {
  const { data } = await fetchOrder(row.order_id);
  current.value = data;
  detailOpen.value = true;
}

onMounted(load);
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}
</style>
