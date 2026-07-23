<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button v-if="canCreate" type="primary" @click="openCreate">新建秒杀</a-button>
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="seckill_active_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'price'">¥{{ record.seckill_price }}</template>
        <template v-else-if="column.key === 'status'">
          {{ record.status === 1 ? '开启' : '关闭' }}
          · {{ record.in_window ? '场次中' : '未开场' }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="openEdit(record)">编辑</a-button>
          <a-button v-if="canToggle" type="link" @click="toggleStatus(record)">
            {{ record.status === 1 ? '关闭' : '开启' }}
          </a-button>
          <a-button v-if="canDelete" type="link" danger @click="onDelete(record)">删除</a-button>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑秒杀' : '新建秒杀'"
      :confirm-loading="saving"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="活动名" required>
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item label="商品 ID" required>
          <a-input-number v-model:value="form.product_id" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="秒杀价" required>
          <a-input-number v-model:value="form.seckill_price" :min="0.01" :step="0.1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="开始日">
          <a-input v-model:value="form.start_day" placeholder="YYYY-MM-DD" />
        </a-form-item>
        <a-form-item label="结束日">
          <a-input v-model:value="form.end_day" placeholder="YYYY-MM-DD" />
        </a-form-item>
        <a-form-item label="场次 IDs">
          <a-input v-model:value="form.seckill_time_ids" placeholder="1" />
        </a-form-item>
        <a-form-item label="开启">
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
  createSeckillActive,
  deleteSeckillActive,
  fetchSeckillActives,
  setSeckillStatus,
  updateSeckillActive,
  type SeckillActive,
} from '@/api/seckill';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canToggle = computed(() => auth.hasPerm('seckill/toggle'));
const canCreate = computed(() => auth.hasPerm('seckill/create'));
const canDelete = computed(() => auth.hasPerm('seckill/delete'));

const loading = ref(false);
const saving = ref(false);
const list = ref<SeckillActive[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({
  name: '',
  product_id: 1,
  seckill_price: 9.9,
  start_day: '',
  end_day: '',
  seckill_time_ids: '1',
  on: true,
});

const columns = [
  { title: 'ID', dataIndex: 'seckill_active_id', width: 70 },
  { title: '活动', dataIndex: 'name' },
  { title: '商品', dataIndex: 'store_name' },
  { title: '商品ID', dataIndex: 'product_id', width: 90 },
  { title: '秒杀价', key: 'price', width: 100 },
  { title: '状态', key: 'status', width: 140 },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchSeckillActives({ page: pagination.current, limit: pagination.pageSize });
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
  form.name = '本店秒杀';
  form.product_id = 1;
  form.seckill_price = 9.9;
  form.start_day = '';
  form.end_day = '';
  form.seckill_time_ids = '1';
  form.on = true;
  modalOpen.value = true;
}

function openEdit(row: SeckillActive) {
  editingId.value = row.seckill_active_id;
  form.name = row.name;
  form.product_id = row.product_id;
  form.seckill_price = row.seckill_price;
  form.start_day = row.start_day;
  form.end_day = row.end_day;
  form.seckill_time_ids = row.seckill_time_ids || '1';
  form.on = row.status === 1;
  modalOpen.value = true;
}

async function submit() {
  if (!form.name.trim() || !form.product_id || form.seckill_price <= 0) {
    message.warning('请填写完整');
    return;
  }
  saving.value = true;
  const payload: Record<string, unknown> = {
    name: form.name.trim(),
    product_id: form.product_id,
    seckill_price: form.seckill_price,
    start_day: form.start_day,
    end_day: form.end_day,
    seckill_time_ids: form.seckill_time_ids || '1',
  };
  if (!editingId.value || canToggle.value) {
    payload.status = form.on ? 1 : 0;
  }
  try {
    if (editingId.value) {
      await updateSeckillActive(editingId.value, payload);
    } else {
      await createSeckillActive(payload);
    }
    message.success('已保存');
    modalOpen.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

async function toggleStatus(row: SeckillActive) {
  const next = row.status === 1 ? 0 : 1;
  await setSeckillStatus(row.seckill_active_id, next);
  message.success(next === 1 ? '已开启' : '已关闭');
  void load();
}

async function onDelete(row: SeckillActive) {
  await deleteSeckillActive(row.seckill_active_id);
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
