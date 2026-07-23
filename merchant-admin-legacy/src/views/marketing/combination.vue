<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button v-if="canCreate" type="primary" @click="openCreate">新建拼团</a-button>
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
      :title="editingId ? '编辑拼团' : '新建拼团'"
      :confirm-loading="saving"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="商品 ID" required>
          <a-input-number v-model:value="form.product_id" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="拼团价" required>
          <a-input-number v-model:value="form.price" :min="0.01" :step="0.1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="成团人数" required>
          <a-input-number v-model:value="form.buying_count_num" :min="2" style="width: 100%" />
        </a-form-item>
        <a-form-item label="开团时长(小时)">
          <a-input-number v-model:value="form.time" :min="1" style="width: 100%" />
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
  createGroup,
  deleteGroup,
  fetchGroups,
  setGroupShow,
  updateGroup,
  type ProductGroup,
} from '@/api/combination';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canToggle = computed(() => auth.hasPerm('combination/toggle'));
const canCreate = computed(() => auth.hasPerm('combination/create'));
const canDelete = computed(() => auth.hasPerm('combination/delete'));

const loading = ref(false);
const saving = ref(false);
const list = ref<ProductGroup[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({
  product_id: 2,
  price: 19.9,
  buying_count_num: 2,
  time: 24,
  on: true,
});

const columns = [
  { title: 'ID', dataIndex: 'product_group_id', width: 70 },
  { title: '商品', dataIndex: 'store_name' },
  { title: '商品ID', dataIndex: 'product_id', width: 90 },
  { title: '拼团价', key: 'price', width: 100 },
  { title: '人数', dataIndex: 'buying_count_num', width: 80 },
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

function openCreate() {
  editingId.value = 0;
  form.product_id = 2;
  form.price = 19.9;
  form.buying_count_num = 2;
  form.time = 24;
  form.on = true;
  modalOpen.value = true;
}

function openEdit(row: ProductGroup) {
  editingId.value = row.product_group_id;
  form.product_id = row.product_id;
  form.price = row.price;
  form.buying_count_num = row.buying_count_num;
  form.time = row.time || 24;
  form.on = row.is_show === 1;
  modalOpen.value = true;
}

async function submit() {
  if (!form.product_id || form.price <= 0 || form.buying_count_num < 2) {
    message.warning('请填写完整');
    return;
  }
  saving.value = true;
  const payload: Record<string, unknown> = {
    product_id: form.product_id,
    price: form.price,
    buying_count_num: form.buying_count_num,
    time: form.time,
  };
  if (!editingId.value || canToggle.value) {
    payload.is_show = form.on ? 1 : 0;
  }
  try {
    if (editingId.value) {
      await updateGroup(editingId.value, payload);
    } else {
      await createGroup(payload);
    }
    message.success('已保存');
    modalOpen.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

async function toggleShow(row: ProductGroup) {
  const next = row.is_show === 1 ? 0 : 1;
  await setGroupShow(row.product_group_id, next);
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
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}
</style>
