<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-select v-model:value="status" allow-clear placeholder="审核状态" style="width: 160px" @change="load">
        <a-select-option :value="0">待审核</a-select-option>
        <a-select-option :value="1">已通过</a-select-option>
        <a-select-option :value="2">已拒绝</a-select-option>
      </a-select>
      <a-input v-model:value="keyword" allow-clear placeholder="商户/联系人/手机" style="width: 220px" @press-enter="load" />
      <a-button type="primary" @click="load">查询</a-button>
    </div>
    <a-table row-key="mer_intention_id" :loading="loading" :columns="columns" :data-source="list" :pagination="pagination" @change="onTableChange">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'status'">
          <a-tag :color="statusColor(record.status)">{{ statusText(record.status) }}</a-tag>
        </template>
        <template v-else-if="column.key === 'action'">
          <a-button v-if="record.status === 0" type="link" @click="openAudit(record)">审核</a-button>
          <span v-else class="muted">—</span>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="open" title="入驻审核" :confirm-loading="saving" @ok="submit">
      <a-form layout="vertical">
        <a-form-item label="审核结果">
          <a-radio-group v-model:value="form.status">
            <a-radio :value="1">通过并创建商户</a-radio>
            <a-radio :value="2">拒绝</a-radio>
          </a-radio-group>
        </a-form-item>
        <template v-if="form.status === 1">
          <a-form-item label="商户账号（可空，默认 mer{id}）">
            <a-input v-model:value="form.account" />
          </a-form-item>
          <a-form-item label="初始密码（可空，默认 admin123）">
            <a-input-password v-model:value="form.password" />
          </a-form-item>
        </template>
        <a-form-item v-else label="拒绝原因">
          <a-textarea v-model:value="form.fail_msg" :rows="3" />
        </a-form-item>
        <a-form-item label="备注">
          <a-input v-model:value="form.mark" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { auditIntention, fetchIntentions, type Intention } from '@/api/merchant';

const loading = ref(false);
const saving = ref(false);
const open = ref(false);
const list = ref<Intention[]>([]);
const keyword = ref('');
const status = ref<number | undefined>(0);
const currentId = ref(0);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const form = reactive({ status: 1, fail_msg: '', mark: '', account: '', password: '' });

const columns = [
  { title: 'ID', dataIndex: 'mer_intention_id', width: 80 },
  { title: '商户名称', dataIndex: 'mer_name' },
  { title: '联系人', dataIndex: 'name', width: 100 },
  { title: '手机', dataIndex: 'phone', width: 130 },
  { title: '状态', key: 'status', width: 100 },
  { title: '关联 mer_id', dataIndex: 'mer_id', width: 110 },
  { title: '操作', key: 'action', width: 100 },
];

function statusText(s: number) {
  return s === 1 ? '已通过' : s === 2 ? '已拒绝' : '待审核';
}
function statusColor(s: number) {
  return s === 1 ? 'green' : s === 2 ? 'red' : 'orange';
}

async function load() {
  loading.value = true;
  try {
    const { data } = await fetchIntentions({
      page: pagination.current,
      limit: pagination.pageSize,
      keyword: keyword.value || undefined,
      status: status.value,
    });
    list.value = data.list || [];
    pagination.total = data.total;
  } finally {
    loading.value = false;
  }
}

function onTableChange(p: { current?: number; pageSize?: number }) {
  pagination.current = p.current || 1;
  pagination.pageSize = p.pageSize || 20;
  load();
}

function openAudit(row: Intention) {
  currentId.value = row.mer_intention_id;
  form.status = 1;
  form.fail_msg = '';
  form.mark = '';
  form.account = '';
  form.password = '';
  open.value = true;
}

async function submit() {
  if (form.status === 2 && !form.fail_msg.trim()) {
    message.warning('请填写拒绝原因');
    return;
  }
  saving.value = true;
  try {
    const { data } = await auditIntention(currentId.value, { ...form });
    if (form.status === 1) {
      message.success(`已通过，mer_id=${data.mer_id}，账号=${data.account}`);
    } else {
      message.success('已拒绝');
    }
    open.value = false;
    load();
  } finally {
    saving.value = false;
  }
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
.muted {
  color: #98a2b3;
}
</style>
