<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-input v-model:value="keyword" allow-clear placeholder="商户名/手机号" style="width: 220px" @press-enter="load" />
      <a-select v-model:value="status" allow-clear placeholder="状态" style="width: 140px" @change="load">
        <a-select-option :value="1">正常</a-select-option>
        <a-select-option :value="0">锁定</a-select-option>
      </a-select>
      <a-button type="primary" @click="load">查询</a-button>
    </div>
    <a-table
      row-key="mer_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="record.status === 1 ? 'green' : 'red'">
            {{ record.status === 1 ? '正常' : '锁定' }}
          </a-tag>
          <a-tag>{{ record.mer_state === 1 ? '营业' : '关闭' }}</a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="toggle(record)">
            {{ record.status === 1 ? '锁定' : '启用' }}
          </a-button>
        </template>
      </template>
    </a-table>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { fetchMerchants, setMerchantStatus, type Merchant } from '@/api/merchant';

const loading = ref(false);
const list = ref<Merchant[]>([]);
const keyword = ref('');
const status = ref<number | undefined>();
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });

const columns = [
  { title: 'ID', dataIndex: 'mer_id', width: 80 },
  { title: '商户名称', dataIndex: 'mer_name' },
  { title: '联系人', dataIndex: 'real_name', width: 120 },
  { title: '手机', dataIndex: 'mer_phone', width: 140 },
  { title: '状态', key: 'status', width: 160 },
  { title: '操作', key: 'action', width: 100 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchMerchants({
      page: pagination.current,
      limit: pagination.pageSize,
      keyword: keyword.value || undefined,
      status: status.value,
    });
    list.value = data.list || [];
    pagination.total = data.total;
  } finally {
    loading.value = false;
  }
}

function onTableChange(p: { current?: number; pageSize?: number }) {
  pagination.current = p.current || 1;
  pagination.pageSize = p.pageSize || 20;
  load();
}

async function toggle(row: Merchant) {
  await setMerchantStatus(row.mer_id, row.status !== 1);
  message.success('已更新');
  load();
}

onMounted(load);
</script>

<style scoped>
.page-card {
  border-radius: 14px;
}
.toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}
</style>
