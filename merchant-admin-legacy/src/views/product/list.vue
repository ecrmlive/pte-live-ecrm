<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-input v-model:value="keyword" allow-clear placeholder="商品名" style="width: 200px" @press-enter="load" />
      <a-select v-model:value="status" allow-clear placeholder="审核状态" style="width: 140px" @change="load">
        <a-select-option :value="0">待审核</a-select-option>
        <a-select-option :value="1">已通过</a-select-option>
        <a-select-option :value="-1">已拒绝</a-select-option>
      </a-select>
      <a-button type="primary" @click="load">查询</a-button>
      <a-button v-if="canCreate" type="primary" ghost @click="$router.push('/product/edit')">发布商品</a-button>
    </div>
    <a-table
      row-key="product_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'price'">¥{{ Number(record.price).toFixed(2) }}</template>
        <template v-else-if="column.key === 'status'">
          <a-tag :color="statusColor(record.status)">{{ statusText(record.status) }}</a-tag>
          <a-tag>{{ record.is_show ? '上架' : '下架' }}</a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="$router.push(`/product/edit?id=${record.product_id}`)">编辑</a-button>
          <a-button v-if="canShow" type="link" @click="toggleShow(record)">{{ record.is_show ? '下架' : '上架' }}</a-button>
          <a-button v-if="canStock" type="link" @click="openStock(record)">库存</a-button>
          <a-popconfirm v-if="canDelete" title="确认删除？" @confirm="onDelete(record.product_id)">
            <a-button type="link" danger>删除</a-button>
          </a-popconfirm>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="stockOpen" title="修改库存" :confirm-loading="saving" @ok="submitStock">
      <a-input-number v-model:value="stockValue" :min="0" style="width: 100%" />
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  deleteProduct,
  fetchProducts,
  setProductShow,
  setProductStock,
  type Product,
} from '@/api/catalog';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canShow = computed(() => auth.hasPerm('product/show'));
const canStock = computed(() => auth.hasPerm('product/stock'));
const canCreate = computed(() => auth.hasPerm('product/create'));
const canDelete = computed(() => auth.hasPerm('product/delete'));

const loading = ref(false);
const saving = ref(false);
const list = ref<Product[]>([]);
const keyword = ref('');
const status = ref<number | undefined>();
const stockOpen = ref(false);
const stockValue = ref(0);
const stockId = ref(0);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });

const columns = [
  { title: 'ID', dataIndex: 'product_id', width: 80 },
  { title: '商品', dataIndex: 'store_name' },
  { title: '分类', dataIndex: 'cate_name', width: 120 },
  { title: '价格', key: 'price', width: 100 },
  { title: '库存', dataIndex: 'stock', width: 90 },
  { title: '状态', key: 'status', width: 160 },
  { title: '操作', key: 'action', width: 260 },
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

async function toggleShow(row: Product) {
  await setProductShow(row.product_id, !row.is_show);
  message.success('已更新');
  load();
}

function openStock(row: Product) {
  stockId.value = row.product_id;
  stockValue.value = row.stock;
  stockOpen.value = true;
}

async function submitStock() {
  saving.value = true;
  try {
    await setProductStock(stockId.value, Number(stockValue.value) || 0);
    message.success('已更新库存');
    stockOpen.value = false;
    load();
  } finally {
    saving.value = false;
  }
}

async function onDelete(id: number) {
  await deleteProduct(id);
  message.success('已删除');
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
  flex-wrap: wrap;
}
</style>
