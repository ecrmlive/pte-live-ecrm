<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button type="primary" @click="openCreate">新建快递</a-button>
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="express_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'show'">{{ record.is_show === 1 ? '展示' : '隐藏' }}</template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="openEdit(record)">编辑</a-button>
          <a-button type="link" danger @click="onDelete(record)">删除</a-button>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑快递' : '新建快递'"
      :confirm-loading="saving"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="名称" required>
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item label="编码">
          <a-input v-model:value="form.code" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="form.sort" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="展示">
          <a-switch v-model:checked="form.show" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  createExpress,
  deleteExpress,
  fetchExpress,
  updateExpress,
  type Express,
} from '@/api/logistics';

const loading = ref(false);
const saving = ref(false);
const list = ref<Express[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({ name: '', code: '', sort: 0, show: true });

const columns = [
  { title: 'ID', dataIndex: 'express_id', width: 80 },
  { title: '名称', dataIndex: 'name' },
  { title: '编码', dataIndex: 'code', width: 140 },
  { title: '排序', dataIndex: 'sort', width: 80 },
  { title: '状态', key: 'show', width: 90 },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchExpress({ page: pagination.current, limit: pagination.pageSize });
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
  form.code = '';
  form.sort = 0;
  form.show = true;
  modalOpen.value = true;
}

function openEdit(row: Express) {
  editingId.value = row.express_id;
  form.name = row.name;
  form.code = row.code;
  form.sort = row.sort;
  form.show = row.is_show === 1;
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
      code: form.code.trim(),
      sort: form.sort,
      is_show: form.show ? 1 : 0,
    };
    if (editingId.value) {
      await updateExpress(editingId.value, payload);
    } else {
      await createExpress(payload);
    }
    message.success('已保存');
    modalOpen.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}

async function onDelete(row: Express) {
  await deleteExpress(row.express_id);
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
