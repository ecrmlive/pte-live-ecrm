<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button type="primary" @click="openCreate">新建角色</a-button>
      <a-button @click="reload">刷新</a-button>
      <span class="hint">勾选菜单 + 按钮；半选父级也会写入 rules（侧栏需要）</span>
    </div>
    <a-table row-key="role_id" :loading="loading" :columns="columns" :data-source="list" :pagination="false">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          {{ record.status === 1 ? '启用' : '禁用' }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="openEdit(record)">编辑权限</a-button>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑角色' : '新建角色'"
      :confirm-loading="saving"
      width="560px"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item label="角色名" required>
          <a-input v-model:value="form.role_name" />
        </a-form-item>
        <a-form-item label="启用">
          <a-switch v-model:checked="form.on" />
        </a-form-item>
        <a-form-item label="菜单/按钮权限">
          <a-tree
            v-if="treeData.length"
            v-model:checkedKeys="checkedKeys"
            checkable
            default-expand-all
            :tree-data="treeData"
            :field-names="{ title: 'title', key: 'menu_id', children: 'children' }"
            @check="onCheck"
          />
          <span v-else class="hint">暂无菜单</span>
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  createRole,
  fetchMenuTree,
  fetchRoles,
  updateRole,
  type MenuNode,
  type SystemRole,
} from '@/api/setting';

const loading = ref(false);
const saving = ref(false);
const list = ref<SystemRole[]>([]);
const treeData = ref<MenuNode[]>([]);
const modalOpen = ref(false);
const editingId = ref(0);
const checkedKeys = ref<(string | number)[]>([]);
const halfCheckedKeys = ref<(string | number)[]>([]);
const form = reactive({ role_name: '', on: true });

const columns = [
  { title: 'ID', dataIndex: 'role_id', width: 70 },
  { title: '名称', dataIndex: 'role_name', width: 160 },
  { title: 'rules', dataIndex: 'rules', ellipsis: true },
  { title: '状态', key: 'status', width: 80 },
  { title: '操作', key: 'action', width: 120 },
];

function parseRules(rules: string): number[] {
  return (rules || '')
    .split(',')
    .map((s) => Number(s.trim()))
    .filter((n) => n > 0);
}

function collectMenuIDs(): number[] {
  const ids = [...checkedKeys.value, ...halfCheckedKeys.value]
    .map((k) => Number(k))
    .filter((n) => n > 0);
  return [...new Set(ids)];
}

function onCheck(
  _checked: unknown,
  info: { halfCheckedKeys?: (string | number)[] },
) {
  halfCheckedKeys.value = info.halfCheckedKeys || [];
}

function decorateTree(nodes: MenuNode[]): Array<MenuNode & { title: string }> {
  return (nodes || []).map((n) => ({
    ...n,
    title: n.is_menu === 2 ? `按钮 · ${n.menu_name}` : n.menu_name,
    children: n.children?.length ? decorateTree(n.children) : undefined,
  }));
}

async function load() {
  loading.value = true;
  try {
    const [{ data }, tree] = await Promise.all([fetchRoles({ page: 1, limit: 50 }), fetchMenuTree()]);
    list.value = data.list || [];
    treeData.value = decorateTree(tree.data.list || []) as unknown as MenuNode[];
  } finally {
    loading.value = false;
  }
}

function reload() {
  void load();
}

function openCreate() {
  editingId.value = 0;
  form.role_name = '';
  form.on = true;
  checkedKeys.value = [];
  halfCheckedKeys.value = [];
  modalOpen.value = true;
}

function openEdit(row: SystemRole) {
  editingId.value = row.role_id;
  form.role_name = row.role_name;
  form.on = row.status === 1;
  checkedKeys.value = parseRules(row.rules);
  halfCheckedKeys.value = [];
  modalOpen.value = true;
}

async function submit() {
  if (!form.role_name.trim()) {
    message.warning('请填写角色名');
    return;
  }
  saving.value = true;
  try {
    const body = {
      role_name: form.role_name,
      status: form.on ? 1 : 0,
      menu_ids: collectMenuIDs(),
    };
    if (editingId.value) {
      await updateRole(editingId.value, body);
    } else {
      await createRole(body);
    }
    message.success('已保存');
    modalOpen.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

onMounted(() => void load());
</script>

<style scoped>
.page-card {
  border-radius: 14px;
}
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
