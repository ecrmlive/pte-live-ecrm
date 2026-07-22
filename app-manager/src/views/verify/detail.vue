<template>
  <a-card :bordered="false" class="card" :loading="loading">
    <template v-if="order">
      <p><strong>{{ order.order_sn }}</strong></p>
      <p>金额 ¥{{ Number(order.pay_price).toFixed(2) }} · 状态 {{ statusText(order.status, order.paid) }}</p>
      <p>{{ order.real_name }} {{ order.user_phone }}</p>
      <p class="muted">{{ order.user_address }}</p>
      <p v-if="order.verify_code">核销码：{{ order.verify_code }}</p>
      <a-form layout="vertical" style="margin-top: 16px">
        <a-form-item label="核销码（若订单有码需一致）">
          <a-input v-model:value="verifyCode" placeholder="可扫码或手输" />
        </a-form-item>
        <a-space>
          <a-button type="primary" :loading="saving" @click="onVerify">确认核销</a-button>
          <a-button @click="onRefund">代退（仅退款）</a-button>
        </a-space>
      </a-form>
    </template>
  </a-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { message, Modal } from 'ant-design-vue';
import { fetchOrder, verifyOrder, type StoreOrder } from '@/api/order';
import { createRefund } from '@/api/refund';

const route = useRoute();
const router = useRouter();
const loading = ref(false);
const saving = ref(false);
const order = ref<StoreOrder | null>(null);
const verifyCode = ref('');

function statusText(status: number, paid: number) {
  if (!paid) return '未支付';
  if (status === 0) return '待发货';
  if (status === 1) return '待收货';
  if (status === 2) return '待评价';
  if (status === 3) return '已完成';
  return String(status);
}

async function load() {
  const id = Number(route.params.id);
  loading.value = true;
  try {
    const { data } = await fetchOrder(id);
    order.value = data;
    verifyCode.value = data.verify_code || '';
  } finally {
    loading.value = false;
  }
}

async function onVerify() {
  if (!order.value) return;
  saving.value = true;
  try {
    await verifyOrder(order.value.order_id, verifyCode.value || undefined);
    message.success('核销成功');
    router.replace('/verify');
  } finally {
    saving.value = false;
  }
}

function onRefund() {
  if (!order.value) return;
  Modal.confirm({
    title: '代用户申请仅退款？',
    onOk: async () => {
      await createRefund(order.value!.order_id, '店员代退');
      message.success('已提交退款申请');
      router.push('/refund');
    },
  });
}

onMounted(load);
</script>

<style scoped>
.card {
  border-radius: 12px;
}
.muted {
  color: #888;
}
</style>
