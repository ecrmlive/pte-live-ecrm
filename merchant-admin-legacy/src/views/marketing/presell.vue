<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button v-if="canCreate" type="primary" @click="openCreate">新建预售</a-button>
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="product_presell_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'type'">
          {{ record.presell_type === 2 ? '定金' : '全款' }}
        </template>
        <template v-else-if="column.key === 'price'">
          ¥{{ record.price }}
          <span v-if="record.presell_type === 2" class="muted">
            （定金¥{{ record.down_price }}）
          </span>
        </template>
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
      :title="editingId ? '编辑预售' : '新建预售'"
      :confirm-loading="saving"
      width="560px"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="商品 ID" required>
          <a-input-number v-model:value="form.product_id" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="活动标题">
          <a-input v-model:value="form.store_name" />
        </a-form-item>
        <a-form-item label="预售类型">
          <a-radio-group v-model:value="form.presell_type">
            <a-radio :value="1">全款</a-radio>
            <a-radio :value="2">定金</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="预售总价" required>
          <a-input-number v-model:value="form.price" :min="0.01" :precision="2" style="width: 100%" />
        </a-form-item>
        <a-form-item v-if="form.presell_type === 2" label="定金" required>
          <a-input-number v-model:value="form.down_price" :min="0.01" :precision="2" style="width: 100%" />
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
  createPresell,
  deletePresell,
  fetchPresells,
  setPresellShow,
  updatePresell,
  type PresellActive,
} from '@/api/presell';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canToggle = computed(() => auth.hasPerm('presell/toggle'));
const canCreate = computed(() => auth.hasPerm('presell/create'));
const canDelete = computed(() => auth.hasPerm('presell/delete'));

const loading = ref(false);
const saving = ref(false);
const list = ref<PresellActive[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({
  product_id: 3,
  store_name: '',
  price: 29.9,
  down_price: 10,
  stock: 100,
  presell_type: 1,
  on: true,
});

const columns = [
  { title: 'ID', dataIndex: 'product_presell_id', width: 70 },
  { title: '活动', dataIndex: 'store_name' },
  { title: '类型', key: 'type', width: 80 },
  { title: '商品', dataIndex: 'product_id', width: 80 },
  { title: '价格', key: 'price', width: 160 },
  { title: '库存', dataIndex: 'stock', width: 80 },
  { title: '状态', key: 'show', width: 80 },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchPresells({ page: pagination.current, limit: pagination.pageSize });
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
  form.product_id = 3;
  form.store_name = '';
  form.price = 29.9;
  form.down_price = 10;
  form.stock = 100;
  form.presell_type = 1;
  form.on = true;
  modalOpen.value = true;
}

function openEdit(row: PresellActive) {
  editingId.value = row.product_presell_id;
  form.product_id = row.product_id;
  form.store_name = row.store_name;
  form.price = Number(row.price);
  form.down_price = Number(row.down_price || 0);
  form.stock = row.stock;
  form.presell_type = row.presell_type === 2 ? 2 : 1;
  form.on = row.is_show === 1;
  modalOpen.value = true;
}

async function submit() {
  saving.value = true;
  try {
    const body: Record<string, unknown> = {
      product_id: form.product_id,
      store_name: form.store_name,
      price: form.price,
      stock: form.stock,
      presell_type: form.presell_type,
    };
    if (!editingId.value || canToggle.value) {
      body.is_show = form.on ? 1 : 0;
    }
    if (form.presell_type === 2) {
      body.down_price = form.down_price;
      body.final_price = Number((form.price - form.down_price).toFixed(2));
    }
    if (editingId.value) {
      await updatePresell(editingId.value, body);
    } else {
      await createPresell(body);
    }
    message.success('已保存');
    modalOpen.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

async function toggleShow(row: PresellActive) {
  const next = row.is_show === 1 ? 0 : 1;
  await setPresellShow(row.product_presell_id, next);
  message.success(next === 1 ? '已上架' : '已下架');
  void load();
}

async function onDelete(row: PresellActive) {
  await deletePresell(row.product_presell_id);
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
.muted {
  color: #999;
  font-size: 12px;
}
</style>
