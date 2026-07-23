<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button v-if="canCreate" type="primary" @click="openCreate">新建助力</a-button>
      <a-button @click="reload">刷新</a-button>
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
          <a-button type="link" @click="openEdit(record)">编辑</a-button>
          <a-button v-if="canToggle" type="link" @click="toggleShow(record)">
            {{ record.is_show === 1 ? '下架' : '上架' }}
          </a-button>
          <a-button v-if="canDelete" type="link" danger @click="onDelete(record)">删除</a-button>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑助力' : '新建助力'"
      :confirm-loading="saving"
      width="520px"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="商品 ID" required>
          <a-input-number v-model:value="form.product_id" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="活动标题">
          <a-input v-model:value="form.store_name" />
        </a-form-item>
        <a-form-item label="助力价" required>
          <a-input-number v-model:value="form.assist_price" :min="0.01" :precision="2" style="width: 100%" />
        </a-form-item>
        <a-form-item label="需助力人数">
          <a-input-number v-model:value="form.assist_count" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="活动库存">
          <a-input-number v-model:value="form.stock" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="上架">
          <a-switch v-model:checked="form.on" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  createAssist,
  deleteAssist,
  fetchAssists,
  setAssistShow,
  updateAssist,
  type ProductAssist,
} from '@/api/assist';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canToggle = computed(() => auth.hasPerm('assist/toggle'));
const canCreate = computed(() => auth.hasPerm('assist/create'));
const canDelete = computed(() => auth.hasPerm('assist/delete'));

const loading = ref(false);
const saving = ref(false);
const list = ref<ProductAssist[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({
  product_id: 13,
  store_name: '',
  assist_price: 9.9,
  assist_count: 1,
  stock: 100,
  on: true,
});

const columns = [
  { title: 'ID', dataIndex: 'product_assist_id', width: 70 },
  { title: '活动', dataIndex: 'store_name' },
  { title: '商品', dataIndex: 'product_id', width: 80 },
  { title: '助力价', key: 'price', width: 100 },
  { title: '人数', dataIndex: 'assist_count', width: 80 },
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

function openCreate() {
  editingId.value = 0;
  form.product_id = 13;
  form.store_name = '';
  form.assist_price = 9.9;
  form.assist_count = 1;
  form.stock = 100;
  form.on = true;
  modalOpen.value = true;
}

function openEdit(row: ProductAssist) {
  editingId.value = row.product_assist_id;
  form.product_id = row.product_id;
  form.store_name = row.store_name;
  form.assist_price = Number(row.assist_price);
  form.assist_count = row.assist_count;
  form.stock = row.stock;
  form.on = row.is_show === 1;
  modalOpen.value = true;
}

async function submit() {
  saving.value = true;
  try {
    const body: Record<string, unknown> = {
      product_id: form.product_id,
      store_name: form.store_name,
      assist_price: form.assist_price,
      assist_count: form.assist_count,
      stock: form.stock,
    };
    if (!editingId.value || canToggle.value) {
      body.is_show = form.on ? 1 : 0;
    }
    if (editingId.value) {
      await updateAssist(editingId.value, body);
    } else {
      await createAssist(body);
    }
    message.success('已保存');
    modalOpen.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

async function toggleShow(row: ProductAssist) {
  const next = row.is_show === 1 ? 0 : 1;
  await setAssistShow(row.product_assist_id, next);
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
  gap: 12px;
  margin-bottom: 16px;
}
</style>
