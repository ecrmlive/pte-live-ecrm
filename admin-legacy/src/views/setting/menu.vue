<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button @click="reload">刷新</a-button>
      <span class="hint">可改名称 / 排序 / 显示；隐藏后侧栏不再下发</span>
    </div>
    <a-table row-key="menu_id" :loading="loading" :columns="columns" :data-source="list" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'show'">
          {{ record.is_show === 1 ? '显示' : '隐藏' }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="toggle(record)">
            {{ record.is_show === 1 ? '隐藏' : '显示' }}
          </a-button>
          <a-button type="link" @click="openEdit(record)">编辑</a-button>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="modalOpen" title="编辑菜单" :confirm-loading="saving" @ok="submit">
      <a-form layout="vertical">
        <a-form-item label="名称">
          <a-input v-model:value="form.menu_name" />
        </a-form-item>
        <a-form-item label="排序">
          <a-input-number v-model:value="form.sort" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { fetchMenus, updateMenu, type SystemMenu } from '@/api/setting';

const loading = ref(false);
const saving = ref(false);
const list = ref<SystemMenu[]>([]);
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({ menu_name: '', sort: 0 });

const columns = [
  { title: 'ID', dataIndex: 'menu_id', width: 70 },
  { title: '父级', dataIndex: 'pid', width: 70 },
  { title: '名称', dataIndex: 'menu_name', width: 140 },
  { title: 'path', dataIndex: 'path' },
  { title: '排序', dataIndex: 'sort', width: 70 },
  { title: '显示', key: 'show', width: 80 },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchMenus();
    list.value = data.list || [];
  } finally {
    loading.value = false;
  }
}

function reload() {
  void load();
}

function openEdit(row: SystemMenu) {
  editingId.value = row.menu_id;
  form.menu_name = row.menu_name;
  form.sort = row.sort;
  modalOpen.value = true;
}

async function submit() {
  saving.value = true;
  try {
    await updateMenu(editingId.value, { menu_name: form.menu_name, sort: form.sort });
    message.success('已保存');
    modalOpen.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

async function toggle(row: SystemMenu) {
  const next = row.is_show === 1 ? 0 : 1;
  await updateMenu(row.menu_id, { is_show: next });
  message.success(next === 1 ? '已显示' : '已隐藏');
  void load();
}

onMounted(() => void load());
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}
.hint {
  color: #888;
  font-size: 13px;
}
</style>
