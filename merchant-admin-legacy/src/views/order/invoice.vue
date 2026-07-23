<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="invoice_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">{{ statusText(record.status) }}</template>
        <template v-else-if="column.key === 'action'">
          <a-button
            v-if="record.status === 0"
            type="link"
            @click="onAudit(record, 1)"
          >
            开票
          </a-button>
          <a-button
            v-if="record.status === 0"
            type="link"
            danger
            @click="onAudit(record, -1)"
          >
            驳回
          </a-button>
        </template>
      </template>
    </a-table>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { auditInvoice, fetchInvoices, type Invoice } from '@/api/invoice';

const loading = ref(false);
const list = ref<Invoice[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });

const columns = [
  { title: 'ID', dataIndex: 'invoice_id', width: 80 },
  { title: '订单', dataIndex: 'order_id', width: 100 },
  { title: '用户', dataIndex: 'uid', width: 90 },
  { title: '抬头', dataIndex: 'header' },
  { title: '税号', dataIndex: 'tax_no', width: 160 },
  { title: '邮箱', dataIndex: 'email', width: 160 },
  { title: '状态', key: 'status', width: 90 },
  { title: '操作', key: 'action', width: 140 },
];

function statusText(s: number) {
  if (s === 1) return '已开';
  if (s === -1) return '驳回';
  return '待审';
}

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchInvoices({ page: pagination.current, limit: pagination.pageSize });
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

async function onAudit(row: Invoice, status: number) {
  await auditInvoice(row.invoice_id, { status, mark: status === 1 ? '已开具' : '驳回' });
  message.success(status === 1 ? '已开票' : '已驳回');
  await load();
}

onMounted(() => void load());
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}
</style>
