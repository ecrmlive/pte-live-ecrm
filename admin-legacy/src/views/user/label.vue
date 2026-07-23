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
        <template v-if="column.key === 'action'">
          <a-button type="link" @click="openEdit(record)">编辑</a-button>
          <a-button type="link" @click="openMark(record)">打标</a-button>
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
        <a-form-item label="标签名" required>
          <a-input v-model:value="form.label_name" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="form.sort" :min="0" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:open="markOpen" title="给用户打标" :confirm-loading="markSaving" @ok="submitMark">
      <a-form layout="vertical">
        <a-form-item label="用户 UID" required>
          <a-input-number v-model:value="markForm.uid" :min="1" style="width: 100%" />
        </a-form-item>
        <a-form-item label="标签（可多选）">
          <a-select
            v-model:value="markForm.label_ids"
            mode="multiple"
            style="width: 100%"
            :options="list.map((x) => ({ label: x.label_name, value: x.label_id }))"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  createUserLabel,
  deleteUserLabel,
  fetchUserLabels,
  markUserLabels,
  updateUserLabel,
  type UserLabel,
} from '@/api/usertag';

const loading = ref(false);
const saving = ref(false);
const list = ref<UserLabel[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({ label_name: '', sort: 0 });
const markOpen = ref(false);
const markSaving = ref(false);
const markForm = reactive<{ uid: number; label_ids: number[] }>({ uid: 1, label_ids: [] });

const columns = [
  { title: 'ID', dataIndex: 'label_id', width: 80 },
  { title: '标签名', dataIndex: 'label_name' },
  { title: '排序', dataIndex: 'sort', width: 80 },
  { title: '操作', key: 'action', width: 220 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchUserLabels({ page: pagination.current, limit: pagination.pageSize });
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
  form.label_name = '';
  form.sort = 0;
  modalOpen.value = true;
}

function openEdit(row: UserLabel) {
  editingId.value = row.label_id;
  form.label_name = row.label_name;
  form.sort = row.sort;
  modalOpen.value = true;
}

function openMark(row: UserLabel) {
  markForm.uid = 1;
  markForm.label_ids = [row.label_id];
  markOpen.value = true;
}

async function submit() {
  if (!form.label_name.trim()) {
    message.warning('请填写标签名');
    return;
  }
  saving.value = true;
  try {
    const payload = { label_name: form.label_name.trim(), sort: form.sort };
    if (editingId.value) {
      await updateUserLabel(editingId.value, payload);
    } else {
      await createUserLabel(payload);
    }
    message.success('已保存');
    modalOpen.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}

async function submitMark() {
  if (!markForm.uid) {
    message.warning('请填写用户 UID');
    return;
  }
  markSaving.value = true;
  try {
    await markUserLabels(markForm.uid, markForm.label_ids);
    message.success('已打标');
    markOpen.value = false;
  } finally {
    markSaving.value = false;
  }
}

async function onDelete(row: UserLabel) {
  await deleteUserLabel(row.label_id);
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
