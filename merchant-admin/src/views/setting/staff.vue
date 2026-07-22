<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-button v-if="canWrite" type="primary" @click="openCreate">新建店员</a-button>
      <a-button @click="reload">刷新</a-button>
      <span class="hint">manager 端登录核销；商户后台子账号见「子账号」</span>
    </div>
    <a-table
      row-key="service_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          {{ record.status === 1 ? '启用' : '禁用' }}
          · 核销{{ record.is_verify === 1 ? '开' : '关' }}
          · 发货{{ record.is_goods === 1 ? '开' : '关' }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button v-if="canWrite" type="link" @click="openEdit(record)">编辑</a-button>
          <a-button v-if="canWrite" type="link" @click="toggle(record)">
            {{ record.status === 1 ? '禁用' : '启用' }}
          </a-button>
        </template>
      </template>
    </a-table>

    <a-modal
      v-model:open="modalOpen"
      :title="editingId ? '编辑店员' : '新建店员'"
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
        <a-form-item label="昵称">
          <a-input v-model:value="form.nickname" />
        </a-form-item>
        <a-form-item label="手机">
          <a-input v-model:value="form.phone" />
        </a-form-item>
        <a-form-item label="可核销">
          <a-switch v-model:checked="form.verify" />
        </a-form-item>
        <a-form-item label="可发货">
          <a-switch v-model:checked="form.goods" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { createStaff, fetchStaff, updateStaff, type Staff } from '@/api/setting';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canWrite = computed(() => auth.hasPerm('staff/write'));

const loading = ref(false);
const saving = ref(false);
const list = ref<Staff[]>([]);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const modalOpen = ref(false);
const editingId = ref(0);
const form = reactive({
  account: '',
  password: '',
  nickname: '',
  phone: '',
  verify: true,
  goods: false,
});

const columns = [
  { title: 'ID', dataIndex: 'service_id', width: 70 },
  { title: '账号', dataIndex: 'account' },
  { title: '昵称', dataIndex: 'nickname' },
  { title: '状态', key: 'status', width: 180 },
  { title: '操作', key: 'action', width: 160 },
];

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchStaff({ page: pagination.current, limit: pagination.pageSize });
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
  form.nickname = '';
  form.phone = '';
  form.verify = true;
  form.goods = false;
  modalOpen.value = true;
}

function openEdit(row: Staff) {
  editingId.value = row.service_id;
  form.account = row.account;
  form.password = '';
  form.nickname = row.nickname;
  form.phone = row.phone;
  form.verify = row.is_verify === 1;
  form.goods = row.is_goods === 1;
  modalOpen.value = true;
}

async function submit() {
  saving.value = true;
  try {
    const body: Record<string, unknown> = {
      nickname: form.nickname,
      phone: form.phone,
      is_verify: form.verify ? 1 : 0,
      is_goods: form.goods ? 1 : 0,
    };
    if (form.password) body.password = form.password;
    if (editingId.value) {
      await updateStaff(editingId.value, body);
    } else {
      await createStaff({ account: form.account, ...body, password: form.password });
    }
    message.success('已保存');
    modalOpen.value = false;
    void load();
  } finally {
    saving.value = false;
  }
}

async function toggle(row: Staff) {
  const next = row.status === 1 ? 0 : 1;
  await updateStaff(row.service_id, { status: next });
  message.success(next === 1 ? '已启用' : '已禁用');
  void load();
}

onMounted(() => void load());
</script>

<style scoped>
.toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}
.hint {
  color: #888;
  font-size: 13px;
}
</style>
