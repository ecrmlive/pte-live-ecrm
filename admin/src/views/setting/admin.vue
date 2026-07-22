<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button type="primary" @click="openCreate">新建管理员</a-button>
      <a-button @click="reload">刷新</a-button>
    </div>
    <a-table
      row-key="admin_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          {{ record.status === 1 ? '启用' : '禁用' }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button type="link" @click="openEdit(record)">编辑</a-button>
          <a-button type="link" @click="toggle(record)">
            {{ record.status === 1 ? '禁用' : '启用' }}
          </a-button>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑管理员' : '新建管理员'"
      :confirm-loading="saving"
      @ok="submit"
    >
      <a-form layout="vertical">
        <a-form-item v-if="!editingId" label="账号" required>
          <a-input v-model:value="form.account" />
        </a-form-item>
        <a-form-item :label="editingId ? '新密码（可空）' : '密码'" :required="!editingId">
          <a-input-password v-model:value="form.password" />
        </a-form-item>
        <a-form-item label="姓名">
          <a-input v-model:value="form.real_name" />
        </a-form-item>
        <a-form-item label="手机">
          <a-input v-model:value="form.phone" />
        </a-form-item>
        <a-form-item label="角色">
          <a-select
            v-model:value="form.roleIds"
            mode="multiple"
            allow-clear
            placeholder="选择角色"
            :options="roleOptions"
            style="width: 100%"
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
  createAdmin,
  fetchAdmins,
  fetchRoles,
  updateAdmin,
  type PlatformAdmin,
} from '@/api/setting';

const loading = ref(false);
const saving = ref(false);
const list = ref<PlatformAdmin[]>([]);
const roleOptions = ref<{ label: string; value: number }[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({
  account: '',
  password: '',
  real_name: '',
  phone: '',
  roleIds: [1] as number[],
});

const columns = [
  { title: 'ID', dataIndex: 'admin_id', width: 70 },
  { title: '账号', dataIndex: 'account' },
  { title: '姓名', dataIndex: 'real_name' },
  { title: '角色', dataIndex: 'roles', width: 120 },
  { title: '状态', key: 'status', width: 80 },
  { title: '操作', key: 'action', width: 160 },
];

function rolesCSV() {
  return form.roleIds.filter((n) => n > 0).join(',') || '1';
}

function parseRoles(roles: string): number[] {
  const ids = (roles || '')
    .split(',')
    .map((s) => Number(s.trim()))
    .filter((n) => n > 0);
  return ids.length ? ids : [1];
}

async function loadRoles() {
  const { data } = await fetchRoles({ page: 1, limit: 100 });
  roleOptions.value = (data.list || []).map((r) => ({
    label: `${r.role_name} (#${r.role_id})`,
    value: r.role_id,
  }));
}

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchAdmins({ page: pagination.current, limit: pagination.pageSize });
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
  form.account = '';
  form.password = 'admin123';
  form.real_name = '';
  form.phone = '';
  form.roleIds = [1];
  modalOpen.value = true;
}

function openEdit(row: PlatformAdmin) {
  editingId.value = row.admin_id;
  form.account = row.account;
  form.password = '';
  form.real_name = row.real_name;
  form.phone = row.phone;
  form.roleIds = parseRoles(row.roles);
  modalOpen.value = true;
}

async function submit() {
  saving.value = true;
  try {
    const roles = rolesCSV();
    if (editingId.value) {
      const body: Record<string, unknown> = {
        real_name: form.real_name,
        phone: form.phone,
        roles,
      };
      if (form.password) body.password = form.password;
      await updateAdmin(editingId.value, body);
    } else {
      await createAdmin({
        account: form.account,
        password: form.password,
        real_name: form.real_name,
        phone: form.phone,
        roles,
      });
    }
    message.success('已保存');
    modalOpen.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

async function toggle(row: PlatformAdmin) {
  const next = row.status === 1 ? 0 : 1;
  await updateAdmin(row.admin_id, { status: next });
  message.success(next === 1 ? '已启用' : '已禁用');
  void load();
}

onMounted(() => {
  void loadRoles();
  void load();
});
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}
</style>
