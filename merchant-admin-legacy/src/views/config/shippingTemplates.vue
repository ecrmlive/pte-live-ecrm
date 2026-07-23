<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button type="primary" @click="openCreate">新建模板</a-button>
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
        <template v-if="column.key === 'type'">{{ typeText(record.type) }}</template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="openEdit(record)">编辑</a-button>
          <a-button type="link" danger @click="onDelete(record)">删除</a-button>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑运费模板' : '新建运费模板'"
      :confirm-loading="saving"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="名称" required>
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item label="计费方式">
          <a-select v-model:value="form.type" :options="typeOptions" />
        </a-form-item>
        <a-form-item label="首件/首重">
          <a-input-number v-model:value="form.first" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="首费">
          <a-input-number v-model:value="form.first_price" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="续件/续重">
          <a-input-number v-model:value="form.continue" :min="0" style="width: 100%" />
        </a-form-item>
        <a-form-item label="续费">
          <a-input-number v-model:value="form.continue_price" :min="0" style="width: 100%" />
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
  createShippingTemplate,
  deleteShippingTemplate,
  fetchShippingTemplates,
  getShippingTemplate,
  updateShippingTemplate,
  type ShippingTemplate,
} from '@/api/logistics';

const loading = ref(false);
const saving = ref(false);
const list = ref<ShippingTemplate[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({
  name: '',
  type: 1,
  sort: 0,
  first: 1,
  first_price: 0,
  continue: 1,
  continue_price: 0,
});

const typeOptions = [
  { label: '按件数', value: 1 },
  { label: '按重量', value: 2 },
  { label: '包邮', value: 3 },
];

const columns = [
  { title: 'ID', dataIndex: 'template_id', width: 80 },
  { title: '名称', dataIndex: 'name' },
  { title: '计费', key: 'type', width: 100 },
  { title: '排序', dataIndex: 'sort', width: 80 },
  { title: '操作', key: 'action', width: 160 },
];

function typeText(t: number) {
  return typeOptions.find((x) => x.value === t)?.label || String(t);
}

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchShippingTemplates({
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
  form.name = '';
  form.type = 3;
  form.sort = 0;
  form.first = 1;
  form.first_price = 0;
  form.continue = 1;
  form.continue_price = 0;
  modalOpen.value = true;
}

async function openEdit(row: ShippingTemplate) {
  editingId.value = row.template_id;
  const { data } = await getShippingTemplate(row.template_id);
  form.name = data.name;
  form.type = data.type;
  form.sort = data.sort;
  const r = data.regions?.[0];
  form.first = r?.first ?? 1;
  form.first_price = r?.first_price ?? 0;
  form.continue = r?.continue ?? 1;
  form.continue_price = r?.continue_price ?? 0;
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
      type: form.type,
      sort: form.sort,
      regions: [
        {
          city_ids: '',
          first: form.first,
          first_price: form.first_price,
          continue: form.continue,
          continue_price: form.continue_price,
        },
      ],
    };
    if (editingId.value) {
      await updateShippingTemplate(editingId.value, payload);
    } else {
      await createShippingTemplate(payload);
    }
    message.success('已保存');
    modalOpen.value = false;
    await load();
  } finally {
    saving.value = false;
  }
}

async function onDelete(row: ShippingTemplate) {
  await deleteShippingTemplate(row.template_id);
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
