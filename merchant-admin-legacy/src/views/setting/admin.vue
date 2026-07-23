<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button v-if="canWrite" type="primary" @click="openCreate">新建子账号</a-button>
      <a-button @click="reload">刷新</a-button>
      <span class="hint">商户后台登录账号；可绑定共享模板或本店角色（店员核销见「员工账号」）</span>
    </div>
    <a-table
      row-key="merchant_admin_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'level'">
          {{ record.level === 0 ? '主账号' : '子账号' }}
        </template>
        <template v-else-if="column.key === 'status'">
          {{ record.status === 1 ? '启用' : '禁用' }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button v-if="canWrite" type="link" @click="openEdit(record)">编辑</a-button>
          <a-button
            v-if="canWrite && record.level !== 0"
            type="link"
            @click="toggle(record)"
          >
            {{ record.status === 1 ? '禁用' : '启用' }}
          </a-button>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑子账号' : '新建子账号'"
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
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { useAuthStore } from '@/stores/auth';
import {
  createAdmin,
  fetchAdmins,
  fetchRoles,
  updateAdmin,
  type MerchantAdmin,
} from '@/api/setting';

const auth = useAuthStore();
const canWrite = computed(() => auth.hasPerm('admins/write'));

const loading = ref(false);
const saving = ref(false);
const list = ref<MerchantAdmin[]>([]);
const roleOptions = ref<{ label: string; value: number }[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({
  account: '',
  password: '',
  real_name: '',
  phone: '',
  roleIds: [2] as number[],
});

const columns = [
  { title: 'ID', dataIndex: 'merchant_admin_id', width: 70 },
  { title: '账号', dataIndex: 'account' },
  { title: '姓名', dataIndex: 'real_name' },
  { title: '类型', key: 'level', width: 90 },
  { title: '角色', dataIndex: 'roles', width: 120 },
  { title: '状态', key: 'status', width: 80 },
  { title: '操作', key: 'action', width: 160 },
];

function rolesCSV() {
  return form.roleIds.filter((n) => n > 0).join(',') || '2';
}

function parseRoles(roles: string): number[] {
  const ids = (roles || '')
    .split(',')
    .map((s) => Number(s.trim()))
    .filter((n) => n > 0);
  return ids.length ? ids : [2];
}

async function loadRoles() {
  const { data } = await fetchRoles({ page: 1, limit: 100 });
  roleOptions.value = (data.list || []).map((r) => ({
    label: `${r.role_name}${r.mer_id === 0 ? '（模板）' : ''} (#${r.role_id})`,
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
  form.roleIds = roleOptions.value.length ? [roleOptions.value[0].value] : [2];
  modalOpen.value = true;
}

function openEdit(row: MerchantAdmin) {
  editingId.value = row.merchant_admin_id;
  form.account = row.account;
  form.password = '';
  form.real_name = row.real_name;
  form.phone = row.phone || '';
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

async function toggle(row: MerchantAdmin) {
  await updateAdmin(row.merchant_admin_id, { status: row.status === 1 ? 0 : 1 });
  message.success('已更新状态');
  void load();
}

onMounted(() => {
  void loadRoles();
  void load();
});
</script>

<style scoped>
.page-card {
  border-radius: 14px;
}
.toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.hint {
  color: #888;
  font-size: 13px;
}
</style>
