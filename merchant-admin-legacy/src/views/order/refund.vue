<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-select v-model:value="status" allow-clear placeholder="状态" style="width: 140px" @change="reload">
        <a-select-option :value="0">待审核</a-select-option>
        <a-select-option :value="3">已退款</a-select-option>
        <a-select-option :value="-1">已拒绝</a-select-option>
      </a-select>
      <a-button type="primary" @click="reload">查询</a-button>
    </div>
    <a-table
      row-key="refund_order_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'price'">¥{{ Number(record.refund_price).toFixed(2) }}</template>
        <template v-else-if="column.key === 'status'">{{ statusText(record.status) }}</template>
        <template v-else-if="column.key === 'action'">
          <template v-if="record.status === 0">
            <a-button v-if="canApprove" type="link" @click="onApprove(record.refund_order_id)">同意退款</a-button>
            <a-button v-if="canReject" type="link" danger @click="openReject(record.refund_order_id)">拒绝</a-button>
          </template>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="rejectOpen" title="拒绝退款" :confirm-loading="saving" @ok="submitReject">
      <a-input v-model:value="failMessage" placeholder="拒绝原因" />
    </a-modal>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { approveRefund, fetchRefunds, rejectRefund, type RefundOrder } from '@/api/refund';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const canApprove = computed(() => auth.hasPerm('order/refund/approve'));
const canReject = computed(() => auth.hasPerm('order/refund/reject'));

const loading = ref(false);
const saving = ref(false);
const list = ref<RefundOrder[]>([]);
const status = ref<number | undefined>(0);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const rejectOpen = ref(false);
const failMessage = ref('不符合退款条件');
const currentId = ref(0);

const columns = [
  { title: 'ID', dataIndex: 'refund_order_id', width: 80 },
  { title: '退款单号', dataIndex: 'refund_order_sn' },
  { title: '子单', dataIndex: 'order_id', width: 90 },
  { title: '金额', key: 'price', width: 100 },
  { title: '原因', dataIndex: 'refund_message' },
  { title: '状态', key: 'status', width: 100 },
  { title: '操作', key: 'action', width: 180 },
];

function statusText(s: number) {
  if (s === 0) return '待审核';
  if (s === 3) return '已退款';
  if (s === -1) return '已拒绝';
  if (s === -2) return '已取消';
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
    const { data } = await fetchRefunds(params);
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
  saving.value = true;
  try {
    await approveRefund(id);
    message.success('已同意并退款');
    load();
  } finally {
    saving.value = false;
  }
}

function openReject(id: number) {
  currentId.value = id;
  failMessage.value = '不符合退款条件';
  rejectOpen.value = true;
}

async function submitReject() {
  saving.value = true;
  try {
    await rejectRefund(currentId.value, failMessage.value || '拒绝');
    message.success('已拒绝');
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
