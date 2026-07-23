<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="seckill_active_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'price'">
          ¥{{ record.seckill_price }}
          <span v-if="record.price" class="ot">原 ¥{{ record.price }}</span>
        </template>
        <template v-else-if="column.key === 'window'">
          {{ record.in_window ? '场次中' : '未开场' }}
        </template>
        <template v-else-if="column.key === 'status'">
          {{ record.status === 1 ? '开启' : '关闭' }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="toggle(record)">
            {{ record.status === 1 ? '关闭' : '开启' }}
          </a-button>
          <a-button type="link" danger @click="onDelete(record)">删除</a-button>
        </template>
      </template>
    </a-table>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  deleteSeckillActive,
  fetchSeckillActives,
  updateSeckillActive,
  type SeckillActive,
} from '@/api/seckill';

const loading = ref(false);
const list = ref<SeckillActive[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });

const columns = [
  { title: 'ID', dataIndex: 'seckill_active_id', width: 70 },
  { title: '活动', dataIndex: 'name' },
  { title: '商户', dataIndex: 'mer_name', width: 120 },
  { title: '商品', dataIndex: 'store_name' },
  { title: '秒杀价', key: 'price', width: 140 },
  { title: '日期', dataIndex: 'start_day', width: 110 },
  { title: '场次', key: 'window', width: 90 },
  { title: '状态', key: 'status', width: 80 },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchSeckillActives({ page: pagination.current, limit: pagination.pageSize });
    list.value = data.list || [];
    pagination.total = data.total || 0;
  } finally {
    loading.value = false;
  }
}

function reload() {
  pagination.current = 1;
  void load();
}

function onTableChange(p: { current?: number; pageSize?: number }) {
  pagination.current = p.current || 1;
  pagination.pageSize = p.pageSize || 20;
  void load();
}

async function toggle(row: SeckillActive) {
  const next = row.status === 1 ? 0 : 1;
  await updateSeckillActive(row.seckill_active_id, { status: next, name: row.name });
  message.success(next === 1 ? '已开启' : '已关闭');
  void load();
}

async function onDelete(row: SeckillActive) {
  await deleteSeckillActive(row.seckill_active_id);
  message.success('已删除');
  void load();
}

onMounted(() => void load());
</script>

<style scoped>
.toolbar {
  margin-bottom: 16px;
}
.ot {
  color: #999;
  margin-left: 6px;
  font-size: 12px;
}
</style>
