<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-select v-model:value="status" allow-clear placeholder="状态" style="width: 140px" @change="reload">
        <a-select-option :value="0">待审核</a-select-option>
        <a-select-option :value="1">已通过</a-select-option>
        <a-select-option :value="-1">已拒绝</a-select-option>
      </a-select>
      <a-button type="primary" @click="reload">查询</a-button>
    </div>
    <a-table
      row-key="financial_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'money'">¥{{ Number(record.extract_money).toFixed(2) }}</template>
        <template v-else-if="column.key === 'status'">{{ statusText(record.status) }}</template>
        <template v-else-if="column.key === 'action'">
          <template v-if="record.status === 0">
            <a-button v-if="canApprove" type="link" @click="onApprove(record.financial_id)">通过</a-button>
            <a-button v-if="canReject" type="link" danger @click="openReject(record.financial_id)">拒绝退余额</a-button>
          </template>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="rejectOpen" title="拒绝提现" :confirm-loading="saving" @ok="submitReject">
      <a-input v-model:value="refusal" placeholder="拒绝原因（将退回商户余额）" />
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import {
  approveWithdraw,
  fetchWithdraws,
  rejectWithdraw,
  type Financial,
} from '@/api/finance';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canApprove = computed(() => auth.hasPerm('accounts/withdraw/approve'));
const canReject = computed(() => auth.hasPerm('accounts/withdraw/reject'));

const loading = ref(false);
const saving = ref(false);
const list = ref<Financial[]>([]);
const status = ref<number | undefined>(0);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const rejectOpen = ref(false);
const refusal = ref('资料不全');
const currentId = ref(0);

const columns = [
  { title: 'ID', dataIndex: 'financial_id', width: 80 },
  { title: '商户', dataIndex: 'mer_id', width: 90 },
  { title: '单号', dataIndex: 'financial_sn' },
  { title: '金额', key: 'money', width: 120 },
  { title: '账户', dataIndex: 'financial_account' },
  { title: '状态', key: 'status', width: 100 },
  { title: '操作', key: 'action', width: 200 },
];

function statusText(s: number) {
  if (s === 0) return '待审核';
  if (s === 1) return '已通过';
  if (s === -1) return '已拒绝';
  return String(s);
}

async function load() {
  loading.value = true;
  try {
    const params: Record<string, unknown> = {
      page: pagination.current,
      limit: pagination.pageSize,
    };
    if (status.value !== undefined) params.status = status.value;
    const { data } = await fetchWithdraws(params);
    list.value = data.list || [];
    pagination.total = data.total;
  } finally {
    loading.value = false;
  }
}

function reload() {
  pagination.current = 1;
  load();
}

function onTableChange(p: { current?: number; pageSize?: number }) {
  pagination.current = p.current || 1;
  pagination.pageSize = p.pageSize || 20;
  load();
}

async function onApprove(id: number) {
  await approveWithdraw(id);
  message.success('已通过');
  load();
}

function openReject(id: number) {
  currentId.value = id;
  refusal.value = '资料不全';
  rejectOpen.value = true;
}

async function submitReject() {
  saving.value = true;
  try {
    await rejectWithdraw(currentId.value, refusal.value || '拒绝');
    message.success('已拒绝并退回余额');
    rejectOpen.value = false;
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
</style>
