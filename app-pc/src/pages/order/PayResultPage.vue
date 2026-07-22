<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  fetchOrderDetail,
  notifyChannelPay,
  payOrder,
  type GroupOrder,
  type PayIntent,
  type PayType,
} from "@/api/trade";
import { ApiError } from "@/utils/request";

const route = useRoute();
const router = useRouter();
const order = ref<GroupOrder | null>(null);
const hint = ref("");
const paying = ref(false);

const id = () => Number(route.params.id);

async function load() {
  try {
    order.value = await fetchOrderDetail(id());
    hint.value = "";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "加载失败";
  }
}

onMounted(() => void load());

function isPayIntent(v: GroupOrder | PayIntent): v is PayIntent {
  return typeof (v as PayIntent).status === "string" && typeof (v as PayIntent).channel === "string";
}

async function pay(type: PayType) {
  paying.value = true;
  try {
    const res = await payOrder(id(), type);
    if (isPayIntent(res)) {
      if (res.status === "paid") {
        order.value = await fetchOrderDetail(id());
      } else {
        if (!res.notify_token) {
          hint.value = "非沙箱环境请完成第三方支付";
          return;
        }
        const uid = order.value?.uid || 0;
        if (!uid) {
          hint.value = "缺少用户信息";
          return;
        }
        order.value = await notifyChannelPay(type as "wechat" | "alipay", {
          group_order_id: res.group_order_id,
          uid,
          out_trade_no: res.out_trade_no,
          pay_price: res.pay_price,
          notify_token: res.notify_token,
        });
      }
    } else {
      order.value = res;
    }
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
          <button class="pc-btn" type="button" :disabled="paying" @click="pay('mock')">Mock 支付</button>
          <button class="pc-btn ghost" type="button" :disabled="paying" @click="pay('balance')">余额支付</button>
          <button class="pc-btn ghost" type="button" :disabled="paying" @click="pay('wechat')">微信支付（沙箱）</button>
          <button class="pc-btn ghost" type="button" :disabled="paying" @click="pay('alipay')">支付宝（沙箱）</button>
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
