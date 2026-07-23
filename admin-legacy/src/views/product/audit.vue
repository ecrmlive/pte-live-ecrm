<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-select v-model:value="status" allow-clear placeholder="审核状态" style="width: 160px" @change="load">
        <a-select-option :value="0">待审核</a-select-option>
        <a-select-option :value="1">已通过</a-select-option>
        <a-select-option :value="-1">已拒绝</a-select-option>
        <a-select-option :value="-2">强制下架</a-select-option>
      </a-select>
      <a-input v-model:value="keyword" allow-clear placeholder="商品名" style="width: 220px" @press-enter="load" />
      <a-button type="primary" @click="load">查询</a-button>
    </div>
    <a-table row-key="product_id" :loading="loading" :columns="columns" :data-source="list" :pagination="pagination" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="statusColor(record.status)">{{ statusText(record.status) }}</a-tag>
        </template>
        <template v-else-if="column.key === 'price'">¥{{ Number(record.price).toFixed(2) }}</template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="openAudit(record)">审核</a-button>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="open" title="商品审核" :confirm-loading="saving" @ok="submit">
      <p class="meta">{{ current?.store_name }} · {{ current?.mer_name || `mer#${current?.mer_id}` }}</p>
      <a-form layout="vertical">
        <a-form-item label="结果">
          <a-radio-group v-model:value="form.status">
            <a-radio :value="1">通过</a-radio>
            <a-radio :value="-1">拒绝</a-radio>
            <a-radio :value="-2">强制下架</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item v-if="form.status === -1" label="拒绝原因">
          <a-textarea v-model:value="form.refusal" :rows="3" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { auditProduct, fetchProducts, type Product } from '@/api/catalog';

const loading = ref(false);
const saving = ref(false);
const open = ref(false);
const list = ref<Product[]>([]);
const keyword = ref('');
const status = ref<number | undefined>(0);
const current = ref<Product | null>(null);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const form = reactive({ status: 1, refusal: '' });

const columns = [
  { title: 'ID', dataIndex: 'product_id', width: 80 },
  { title: '商品', dataIndex: 'store_name' },
  { title: '商户', dataIndex: 'mer_name', width: 140 },
  { title: '分类', dataIndex: 'cate_name', width: 120 },
  { title: '价格', key: 'price', width: 100 },
  { title: '库存', dataIndex: 'stock', width: 90 },
  { title: '状态', key: 'status', width: 110 },
  { title: '操作', key: 'action', width: 90 },
];

function statusText(s: number) {
  if (s === 1) return '已通过';
  if (s === -1) return '已拒绝';
  if (s === -2) return '强制下架';
  return '待审核';
}
function statusColor(s: number) {
  if (s === 1) return 'green';
  if (s === -1 || s === -2) return 'red';
  return 'orange';
}

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchProducts({
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

function openAudit(row: Product) {
  current.value = row;
  form.status = 1;
  form.refusal = '';
  open.value = true;
}

async function submit() {
  if (!current.value) return;
  if (form.status === -1 && !form.refusal.trim()) {
    message.warning('请填写拒绝原因');
    return;
  }
  saving.value = true;
  try {
    await auditProduct(current.value.product_id, { ...form });
    message.success('已提交');
    open.value = false;
    load();
  } finally {
    saving.value = false;
  }
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
.meta {
  margin: 0 0 12px;
  color: #516070;
}
</style>
