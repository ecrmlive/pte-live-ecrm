<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button @click="reload">刷新</a-button>
      <span class="hint">助力活动监管 · 可上下架/删除</span>
    </div>
    <a-table
      row-key="product_assist_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'price'">¥{{ record.assist_price }}</template>
        <template v-else-if="column.key === 'show'">
          {{ record.is_show === 1 ? '上架' : '下架' }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="toggle(record)">
            {{ record.is_show === 1 ? '下架' : '上架' }}
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
import { deleteAssist, fetchAssists, updateAssist, type ProductAssist } from '@/api/assist';

const loading = ref(false);
const list = ref<ProductAssist[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });

const columns = [
  { title: 'ID', dataIndex: 'product_assist_id', width: 70 },
  { title: '活动', dataIndex: 'store_name' },
  { title: '商户', dataIndex: 'mer_name', width: 120 },
  { title: '助力价', key: 'price', width: 100 },
  { title: '需人数', dataIndex: 'assist_count', width: 80 },
  { title: '库存', dataIndex: 'stock', width: 80 },
  { title: '状态', key: 'show', width: 80 },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchAssists({ page: pagination.current, limit: pagination.pageSize });
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

async function toggle(row: ProductAssist) {
  const next = row.is_show === 1 ? 0 : 1;
  await updateAssist(row.product_assist_id, { is_show: next, assist_price: row.assist_price });
  message.success(next === 1 ? '已上架' : '已下架');
  void load();
}

async function onDelete(row: ProductAssist) {
  await deleteAssist(row.product_assist_id);
  message.success('已删除');
  void load();
}

onMounted(() => void load());
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}
.hint {
  color: #888;
  font-size: 13px;
}
</style>
