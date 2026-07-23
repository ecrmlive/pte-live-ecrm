<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button type="primary" @click="openCreate">新建参数模板</a-button>
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="template_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <a-button type="link" @click="openEdit(record)">编辑</a-button>
          <a-button type="link" danger @click="onDelete(record)">删除</a-button>
        </template>
      </template>
    </a-table>
    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑参数模板' : '新建参数模板'"
      :confirm-loading="saving"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="模板名" required>
          <a-input v-model:value="form.template_name" />
        </a-form-item>
        <a-form-item label="参数 JSON" required>
          <a-textarea v-model:value="form.template_value" :rows="5" placeholder='[{"name":"材质","value":"棉"}]' />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="form.sort" :min="0" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  createAttrTemplate,
  deleteAttrTemplate,
  fetchAttrTemplates,
  updateAttrTemplate,
  type AttrTemplate,
} from '@/api/productmeta';

const loading = ref(false);
const saving = ref(false);
const list = ref<AttrTemplate[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({ template_name: '', template_value: '[]', sort: 0 });
const columns = [
  { title: 'ID', dataIndex: 'template_id', width: 80 },
  { title: '名称', dataIndex: 'template_name' },
  { title: '参数', dataIndex: 'template_value', ellipsis: true },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchAttrTemplates({
      page: pagination.current,
      limit: pagination.pageSize,
    });
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
  form.template_name = '';
  form.template_value = '[{"name":"材质","value":"棉"}]';
  form.sort = 0;
  modalOpen.value = true;
}
function openEdit(row: AttrTemplate) {
  editingId.value = row.template_id;
  form.template_name = row.template_name;
  form.template_value = row.template_value;
  form.sort = row.sort;
  modalOpen.value = true;
}
async function submit() {
  if (!form.template_name.trim()) {
    message.warning('请填写模板名');
    return;
  }
  saving.value = true;
  try {
    const payload = {
      template_name: form.template_name.trim(),
      template_value: form.template_value.trim() || '[]',
      sort: form.sort,
    };
    if (editingId.value) await updateAttrTemplate(editingId.value, payload);
    else await createAttrTemplate(payload);
    message.success('已保存');
    modalOpen.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}
async function onDelete(row: AttrTemplate) {
  await deleteAttrTemplate(row.template_id);
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
