<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button type="primary" @click="openCreate(0)">新增一级分类</a-button>
      <a-button @click="load">刷新</a-button>
    </div>
    <a-table
      row-key="store_category_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="false"
      default-expand-all-rows
      :children-column-name="'children'"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'is_show'">
          <a-tag :color="record.is_show ? 'green' : 'default'">{{ record.is_show ? '显示' : '隐藏' }}</a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="openCreate(record.store_category_id)">加子类</a-button>
          <a-button type="link" @click="openEdit(record)">编辑</a-button>
          <a-popconfirm title="确认删除？" @confirm="onDelete(record.store_category_id)">
            <a-button type="link" danger>删除</a-button>
          </a-popconfirm>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="open" :title="form.id ? '编辑分类' : '新增分类'" :confirm-loading="saving" @ok="submit">
      <a-form layout="vertical">
        <a-form-item label="名称">
          <a-input v-model:value="form.cate_name" />
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
import {
  createCategory,
  deleteCategory,
  fetchCategoryTree,
  updateCategory,
  type CategoryNode,
} from '@/api/catalog';

const loading = ref(false);
const saving = ref(false);
const open = ref(false);
const list = ref<CategoryNode[]>([]);
const form = reactive({ id: 0, pid: 0, cate_name: '', sort: 0, show: true });

const columns = [
  { title: '名称', dataIndex: 'cate_name' },
  { title: 'ID', dataIndex: 'store_category_id', width: 90 },
  { title: '排序', dataIndex: 'sort', width: 90 },
  { title: '显示', key: 'is_show', width: 90 },
  { title: '操作', key: 'action', width: 220 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchCategoryTree();
    list.value = data.list || [];
  } finally {
    loading.value = false;
  }
}

function openCreate(pid: number) {
  form.id = 0;
  form.pid = pid;
  form.cate_name = '';
  form.sort = 0;
  form.show = true;
  open.value = true;
}

function openEdit(row: CategoryNode) {
  form.id = row.store_category_id;
  form.pid = row.pid;
  form.cate_name = row.cate_name;
  form.sort = row.sort;
  form.show = !!row.is_show;
  open.value = true;
}

async function submit() {
  if (!form.cate_name.trim()) {
    message.warning('请填写名称');
    return;
  }
  saving.value = true;
  try {
    if (form.id) {
      await updateCategory(form.id, {
        cate_name: form.cate_name,
        sort: form.sort,
        is_show: form.show ? 1 : 0,
      });
    } else {
      await createCategory({
        pid: form.pid,
        cate_name: form.cate_name,
        sort: form.sort,
        is_show: form.show ? 1 : 0,
      });
    }
    message.success('已保存');
    open.value = false;
    load();
  } finally {
    saving.value = false;
  }
}

async function onDelete(id: number) {
  await deleteCategory(id);
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
