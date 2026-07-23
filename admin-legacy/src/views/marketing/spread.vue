<template>
  <a-card :bordered="false" class="page-card">
    <a-tabs v-model:activeKey="tab">
      <a-tab-pane key="logs" tab="绑定日志">
        <a-table
          row-key="user_spread_log_id"
          :loading="loading"
          :columns="logColumns"
          :data-source="logs"
          :pagination="pagination"
          @change="onTableChange"
        />
      </a-tab-pane>
      <a-tab-pane key="bills" tab="佣金流水">
        <a-table
          row-key="bill_id"
          :loading="loading"
          :columns="billColumns"
          :data-source="bills"
          :pagination="pagination"
          @change="onTableChange"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'number'">¥{{ Number(record.number).toFixed(2) }}</template>
          </template>
        </a-table>
      </a-tab-pane>
    </a-tabs>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue';
import { fetchSpreadBills, fetchSpreadLogs, type SpreadLog } from '@/api/coupon';

const tab = ref('logs');
const loading = ref(false);
const logs = ref<SpreadLog[]>([]);
const bills = ref<Record<string, unknown>[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });

const logColumns = [
  { title: 'ID', dataIndex: 'user_spread_log_id', width: 90 },
  { title: '用户', dataIndex: 'uid', width: 90 },
  { title: '原推广人', dataIndex: 'old_spread_uid', width: 110 },
  { title: '新推广人', dataIndex: 'spread_uid', width: 110 },
  { title: '时间', dataIndex: 'create_time' },
];

const billColumns = [
  { title: 'ID', dataIndex: 'bill_id', width: 90 },
  { title: '用户', dataIndex: 'uid', width: 90 },
  { title: '标题', dataIndex: 'title' },
  { title: '金额', key: 'number', width: 120 },
  { title: '关联单', dataIndex: 'link_id' },
  { title: '时间', dataIndex: 'create_time' },
];

async function load() {
  loading.value = true;
  try {
    const params = { page: pagination.current, limit: pagination.pageSize };
    if (tab.value === 'logs') {
      const { data } = await fetchSpreadLogs(params);
      logs.value = data.list || [];
      pagination.total = data.total;
    } else {
      const { data } = await fetchSpreadBills(params);
      bills.value = (data.list || []) as Record<string, unknown>[];
      pagination.total = data.total;
    }
  } finally {
    loading.value = false;
  }
}

function onTableChange(p: { current?: number; pageSize?: number }) {
  pagination.current = p.current || 1;
  pagination.pageSize = p.pageSize || 20;
  load();
}

watch(tab, () => {
  pagination.current = 1;
  load();
});

onMounted(load);
</script>
