<template>
  <a-card :bordered="false" class="page-card">
    <div class="toolbar">
      <a-select v-model:value="paidFilter" allow-clear placeholder="支付" style="width: 120px" @change="reload">
        <a-select-option :value="1">已支付</a-select-option>
        <a-select-option :value="0">未支付</a-select-option>
      </a-select>
      <a-select
        v-model:value="statusFilter"
        :disabled="pendingOnly"
        allow-clear
        placeholder="状态"
        style="width: 140px"
        @change="reload"
      >
        <a-select-option :value="0">待发货</a-select-option>
        <a-select-option :value="1">待收货</a-select-option>
        <a-select-option :value="2">待评价</a-select-option>
        <a-select-option :value="3">已完成</a-select-option>
      </a-select>
      <a-button type="primary" @click="reload">查询</a-button>
    </div>

    <a-table
      row-key="order_id"
      :loading="loading"
      :columns="columns"
      :data-source="list"
      :pagination="pagination"
      @change="onTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'paid'">
          <a-tag :color="record.paid ? 'green' : 'orange'">{{ record.paid ? '已支付' : '未支付' }}</a-tag>
        </template>
        <template v-else-if="column.key === 'status'">{{ statusText(record.status, record.paid) }}</template>
        <template v-else-if="column.key === 'price'">¥{{ Number(record.pay_price).toFixed(2) }}</template>
        <template v-else-if="column.key === 'action'">
          <a-button
            v-if="canDeliver && record.paid === 1 && record.status === 0"
            type="link"
            @click="openDeliver(record)"
          >
            发货
          </a-button>
          <a-button
            v-if="canVerify && record.paid === 1 && record.status !== 3 && record.status !== -1"
            type="link"
            @click="doVerify(record.order_id)"
          >
            核销
          </a-button>
          <a-button type="link" @click="openDetail(record.order_id)">详情</a-button>
        </template>
      </template>
    </a-table>

    <a-modal v-model:open="deliverOpen" title="发货" :confirm-loading="saving" @ok="submitDeliver">
      <a-form layout="vertical">
        <a-form-item label="快递公司">
          <a-input v-model:value="deliverForm.delivery_name" placeholder="如：顺丰速运" />
        </a-form-item>
        <a-form-item label="快递单号">
          <a-input v-model:value="deliverForm.delivery_id" />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-drawer v-model:open="detailOpen" title="订单详情" width="420">
      <template v-if="detail">
        <p>子单号：{{ detail.order_sn }}</p>
        <p>金额：¥{{ Number(detail.pay_price).toFixed(2) }}</p>
        <p>收货：{{ detail.real_name }} {{ detail.user_phone }}</p>
        <p>{{ detail.user_address }}</p>
        <p v-if="detail.delivery_id">物流：{{ detail.delivery_name }} {{ detail.delivery_id }}</p>
      </template>
    </a-drawer>
  </a-card>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { message } from 'ant-design-vue';
import { deliverOrder, fetchOrder, fetchOrders, verifyOrder, type StoreOrder } from '@/api/order';
import { useAuthStore } from '@/stores/auth';

const route = useRoute();
const auth = useAuthStore();
const canDeliver = computed(() => auth.hasPerm('order/deliver'));
const canVerify = computed(() => auth.hasPerm('order/verify'));
const pendingOnly = computed(() => route.path.includes('/order/delivery'));

const loading = ref(false);
const saving = ref(false);
const list = ref<StoreOrder[]>([]);
const paidFilter = ref<number | undefined>(1);
const statusFilter = ref<number | undefined>(pendingOnly.value ? 0 : undefined);
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showSizeChanger: true });
const deliverOpen = ref(false);
const detailOpen = ref(false);
const currentId = ref(0);
const detail = ref<StoreOrder | null>(null);
const deliverForm = reactive({ delivery_name: '演示快递', delivery_id: '' });

const columns = [
  { title: 'ID', dataIndex: 'order_id', width: 80 },
  { title: '子单号', dataIndex: 'order_sn', width: 180 },
  { title: '收货人', dataIndex: 'real_name', width: 100 },
  { title: '金额', key: 'price', width: 100 },
  { title: '支付', key: 'paid', width: 90 },
  { title: '状态', key: 'status', width: 100 },
  { title: '操作', key: 'action', width: 140 },
];

function statusText(status: number, paid: number) {
  if (!paid) return '待支付';
  if (status === 0) return '待发货';
  if (status === 1) return '待收货';
  if (status === 2) return '待评价';
  if (status === 3) return '已完成';
  if (status === -1) return '已退款';
  return String(status);
}

async function load() {
  loading.value = true;
  try {
    const params: Record<string, unknown> = {
      page: pagination.current,
      limit: pagination.pageSize,
    };
    if (paidFilter.value !== undefined) params.paid = paidFilter.value;
    if (pendingOnly.value) {
      params.paid = 1;
      params.status = 0;
    } else if (statusFilter.value !== undefined) {
      params.status = statusFilter.value;
    }
    const { data } = await fetchOrders(params);
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

function openDeliver(record: StoreOrder) {
  currentId.value = record.order_id;
  deliverForm.delivery_name = '演示快递';
  deliverForm.delivery_id = `SF${Date.now().toString().slice(-10)}`;
  deliverOpen.value = true;
}

async function submitDeliver() {
  if (!deliverForm.delivery_name || !deliverForm.delivery_id) {
    message.warning('请填写物流信息');
    return;
  }
  saving.value = true;
  try {
    await deliverOrder(currentId.value, {
      delivery_name: deliverForm.delivery_name,
      delivery_id: deliverForm.delivery_id,
      delivery_type: '1',
    });
    message.success('已发货');
    deliverOpen.value = false;
    load();
  } finally {
    saving.value = false;
  }
}

async function openDetail(id: number) {
  const { data } = await fetchOrder(id);
  detail.value = data;
  detailOpen.value = true;
}

async function doVerify(id: number) {
  saving.value = true;
  try {
    await verifyOrder(id);
    message.success('已核销完成');
    load();
  } finally {
    saving.value = false;
  }
}

watch(
  () => route.path,
  () => {
    statusFilter.value = pendingOnly.value ? 0 : undefined;
    paidFilter.value = 1;
    reload();
  },
);

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
  flex-wrap: wrap;
}
</style>
