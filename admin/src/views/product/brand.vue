<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button type="primary" @click="openCreate">新增品牌</a-button>
      <a-button @click="load">刷新</a-button>
    </div>
    <a-table row-key="brand_id" :loading="loading" :columns="columns" :data-source="list" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'is_show'">
          <a-tag :color="record.is_show ? 'green' : 'default'">{{ record.is_show ? '显示' : '隐藏' }}</a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="openEdit(record)">编辑</a-button>
          <a-popconfirm title="确认删除？" @confirm="onDelete(record.brand_id)">
            <a-button type="link" danger>删除</a-button>
          </a-popconfirm>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="open" :title="form.id ? '编辑品牌' : '新增品牌'" :confirm-loading="saving" @ok="submit">
      <a-form layout="vertical">
        <a-form-item label="品牌名称">
          <a-input v-model:value="form.brand_name" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="form.sort" style="width: 100%" />
        </a-form-item>
        <a-form-item label="显示">
          <a-switch v-model:checked="form.show" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { createBrand, deleteBrand, fetchBrands, updateBrand, type Brand } from '@/api/catalog';

const loading = ref(false);
const saving = ref(false);
const open = ref(false);
const list = ref<Brand[]>([]);
const form = reactive({ id: 0, brand_name: '', sort: 0, show: true });

const columns = [
  { title: 'ID', dataIndex: 'brand_id', width: 90 },
  { title: '品牌', dataIndex: 'brand_name' },
  { title: '排序', dataIndex: 'sort', width: 90 },
  { title: '显示', key: 'is_show', width: 90 },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchBrands();
    list.value = data.list || [];
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  form.id = 0;
  form.brand_name = '';
  form.sort = 0;
  form.show = true;
  open.value = true;
}

function openEdit(row: Brand) {
  form.id = row.brand_id;
  form.brand_name = row.brand_name;
  form.sort = row.sort;
  form.show = !!row.is_show;
  open.value = true;
}

async function submit() {
  if (!form.brand_name.trim()) {
    message.warning('请填写品牌名');
    return;
  }
  saving.value = true;
  try {
    const body = { brand_name: form.brand_name, sort: form.sort, is_show: form.show ? 1 : 0 };
    if (form.id) await updateBrand(form.id, body);
    else await createBrand(body);
    message.success('已保存');
    open.value = false;
    load();
  } finally {
    saving.value = false;
  }
}

async function onDelete(id: number) {
  await deleteBrand(id);
  message.success('已删除');
  load();
}

onMounted(load);
</script>

<style scoped>
.page-card {
  border-radius: 14px;
}
.toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}
</style>
