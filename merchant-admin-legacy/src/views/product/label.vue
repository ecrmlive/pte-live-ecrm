<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button type="primary" @click="openCreate">新建标签</a-button>
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="label_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">{{ record.status === 1 ? '启用' : '停用' }}</template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="openEdit(record)">编辑</a-button>
          <a-button type="link" danger @click="onDelete(record)">删除</a-button>
        </template>
      </template>
    </a-table>
    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑标签' : '新建标签'"
      :confirm-loading="saving"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="名称" required>
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item label="说明">
          <a-input v-model:value="form.info" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="form.sort" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model:checked="form.enabled" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  createLabel,
  deleteLabel,
  fetchLabels,
  updateLabel,
  type ProductLabel,
} from '@/api/productmeta';

const loading = ref(false);
const saving = ref(false);
const list = ref<ProductLabel[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({ name: '', info: '', sort: 0, enabled: true });
const columns = [
  { title: 'ID', dataIndex: 'label_id', width: 80 },
  { title: '名称', dataIndex: 'name' },
  { title: '说明', dataIndex: 'info' },
  { title: '状态', key: 'status', width: 90 },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchLabels({ page: pagination.current, limit: pagination.pageSize });
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
  form.name = '';
  form.info = '';
  form.sort = 0;
  form.enabled = true;
  modalOpen.value = true;
}
function openEdit(row: ProductLabel) {
  editingId.value = row.label_id;
  form.name = row.name;
  form.info = row.info;
  form.sort = row.sort;
  form.enabled = row.status === 1;
  modalOpen.value = true;
}
async function submit() {
  if (!form.name.trim()) {
    message.warning('请填写名称');
    return;
  }
  saving.value = true;
  try {
    const payload = {
      name: form.name.trim(),
      info: form.info.trim(),
      sort: form.sort,
      status: form.enabled ? 1 : 0,
    };
    if (editingId.value) await updateLabel(editingId.value, payload);
    else await createLabel(payload);
    message.success('已保存');
    modalOpen.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}
async function onDelete(row: ProductLabel) {
  await deleteLabel(row.label_id);
  message.success('已删除');
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
