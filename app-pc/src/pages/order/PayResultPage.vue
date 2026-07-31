<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  fetchOrderDetail,
	fetchOrderPaymentChannels,
  payOrder,
  type GroupOrder,
	type PaymentChannel,
  type PayType,
} from "@/api/trade";
import { ApiError } from "@/utils/request";

const route = useRoute();
const router = useRouter();
const order = ref<GroupOrder | null>(null);
const hint = ref("");
const paying = ref(false);
const paymentChannels = ref<PaymentChannel[]>([]);

const id = () => Number(route.params.id);

async function load() {
  try {
    const [detail, channels] = await Promise.all([
      fetchOrderDetail(id()),
      fetchOrderPaymentChannels(id()),
    ]);
    order.value = detail;
    paymentChannels.value = channels.list ?? [];
    hint.value = "";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "加载失败";
  }
}

onMounted(() => void load());

function channelEnabled(channel: "wechat" | "alipay") {
  return paymentChannels.value.some((item) => item.channel === channel && item.enabled);
}

async function pay(type: PayType) {
  paying.value = true;
  try {
    const res = await payOrder(id(), type);
    order.value = res as GroupOrder;
    hint.value = "支付成功";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "支付失败";
  } finally {
    paying.value = false;
  }
}
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <h1>支付订单</h1>
      <p v-if="hint" class="hint">{{ hint }}</p>
      <div v-if="order">
        <p>主单号 {{ order.group_order_sn }}</p>
        <p>应付 <strong>¥{{ Number(order.pay_price).toFixed(2) }}</strong> · {{ order.total_num }} 件</p>
        <p>状态：{{ order.paid === 1 ? "已支付" : "待支付" }}</p>
        <div v-if="order.paid !== 1" class="actions">
          <button v-if="channelEnabled('wechat')" class="pc-btn ghost" type="button" :disabled="paying" @click="pay('wechat')">微信支付</button>
          <button v-if="channelEnabled('alipay')" class="pc-btn ghost" type="button" :disabled="paying" @click="pay('alipay')">支付宝支付</button>
          <p v-if="!channelEnabled('wechat') && !channelEnabled('alipay')" class="hint">当前订单暂无可用的第三方支付方式，请联系店铺或稍后重试。</p>
        </div>
        <div v-else class="actions">
          <RouterLink class="pc-btn" to="/orders">查看订单</RouterLink>
          <button class="pc-btn ghost" type="button" @click="router.push('/goods')">继续购物</button>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.panel {
  background: var(--pc-surface);
  border: 1px solid var(--pc-line);
  border-radius: calc(var(--pc-radius) + 4px);
  padding: 1.5rem;
  box-shadow: var(--pc-shadow);
}
.hint { color: var(--pc-muted); }
.actions { display: flex; flex-wrap: wrap; gap: 0.8rem; margin-top: 1.2rem; }
.ghost { background: transparent; color: inherit; border: 1px solid var(--pc-line); }
</style>
