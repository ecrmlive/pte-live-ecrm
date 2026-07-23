<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button type="primary" @click="openCreate">新建保障</a-button>
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="guarantee_id"
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
      :title="editingId ? '编辑保障' : '新建保障'"
      :confirm-loading="saving"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="名称" required>
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item label="内容">
          <a-textarea v-model:value="form.content" :rows="3" />
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
  createGuarantee,
  deleteGuarantee,
  fetchGuarantees,
  updateGuarantee,
  type Guarantee,
} from '@/api/productmeta';

const loading = ref(false);
const saving = ref(false);
const list = ref<Guarantee[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({ name: '', content: '', sort: 0, enabled: true });
const columns = [
  { title: 'ID', dataIndex: 'guarantee_id', width: 80 },
  { title: '名称', dataIndex: 'name' },
  { title: '内容', dataIndex: 'content' },
  { title: '状态', key: 'status', width: 90 },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchGuarantees({ page: pagination.current, limit: pagination.pageSize });
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
  form.content = '';
  form.sort = 0;
  form.enabled = true;
  modalOpen.value = true;
}
function openEdit(row: Guarantee) {
  editingId.value = row.guarantee_id;
  form.name = row.name;
  form.content = row.content;
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
      content: form.content.trim(),
      sort: form.sort,
      status: form.enabled ? 1 : 0,
    };
    if (editingId.value) await updateGuarantee(editingId.value, payload);
    else await createGuarantee(payload);
    message.success('已保存');
    modalOpen.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}
async function onDelete(row: Guarantee) {
  await deleteGuarantee(row.guarantee_id);
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
