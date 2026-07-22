<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="product_group_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'price'">¥{{ record.price }}</template>
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
import { deleteGroup, fetchGroups, updateGroup, type ProductGroup } from '@/api/combination';

const loading = ref(false);
const list = ref<ProductGroup[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });

const columns = [
  { title: 'ID', dataIndex: 'product_group_id', width: 70 },
  { title: '商品', dataIndex: 'store_name' },
  { title: '商户', dataIndex: 'mer_name', width: 120 },
  { title: '拼团价', key: 'price', width: 100 },
  { title: '成团人数', dataIndex: 'buying_count_num', width: 100 },
  { title: '状态', key: 'show', width: 80 },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchGroups({ page: pagination.current, limit: pagination.pageSize });
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

async function toggle(row: ProductGroup) {
  const next = row.is_show === 1 ? 0 : 1;
  await updateGroup(row.product_group_id, { is_show: next, price: row.price });
  message.success(next === 1 ? '已上架' : '已下架');
  void load();
}

async function onDelete(row: ProductGroup) {
  await deleteGroup(row.product_group_id);
  message.success('已删除');
  void load();
}

onMounted(() => void load());
</script>

<style scoped>
.toolbar {
  margin-bottom: 16px;
}
</style>
