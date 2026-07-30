<template>
  <view class="page">
    <view class="card">
      <text class="label">主单号</text>
      <text class="val">{{ group?.group_order_sn || "-" }}</text>
      <text class="price">¥{{ Number(group?.pay_price || 0).toFixed(2) }}</text>
      <text class="tip">子单数 {{ group?.orders?.length || 0 }}（多商户一单）</text>
    </view>
    <view class="qx-btn qx-btn-primary" @click="pay('mock')">Mock 支付</view>
    <view v-if="channelEnabled('wechat')" class="qx-btn qx-btn-ghost mt" @click="pay('wechat')">微信支付</view>
    <view v-if="channelEnabled('alipay')" class="qx-btn qx-btn-ghost mt" @click="pay('alipay')">支付宝支付</view>
    <text v-if="!channelEnabled('wechat') && !channelEnabled('alipay')" class="tip">当前订单暂无可用的第三方支付方式</text>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import {
  fetchOrderDetail,
	fetchPaymentChannels,
  payGroup,
  type GroupOrder,
	type PaymentChannel,
  type PayType,
} from "@/api/order";

const id = ref(0);
const group = ref<GroupOrder | null>(null);
const paymentChannels = ref<PaymentChannel[]>([]);

onLoad(async (q) => {
  id.value = Number(q?.id || 0);
  if (!id.value) return;
  try {
    const [detail, channels] = await Promise.all([fetchOrderDetail(id.value), fetchPaymentChannels(id.value)]);
    group.value = detail;
    paymentChannels.value = channels.list ?? [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function channelEnabled(channel: "wechat" | "alipay") {
  return paymentChannels.value.some((item) => item.channel === channel && item.enabled);
}

async function pay(type: PayType) {
  try {
    const res = await payGroup(id.value, type);
    const g = res as GroupOrder;
    const awaitFinal = (g.orders || []).find((o) => o.status === 10);
    if (awaitFinal) {
      const { fetchPresellFinals } = await import("@/api/presell");
      const finals = await fetchPresellFinals(true);
      const hit = (finals.list || []).find((f) => f.order_id === awaitFinal.order_id);
      if (hit) {
        uni.redirectTo({ url: `/pages/presell/final-pay?id=${hit.presell_order_id}` });
        return;
      }
    }
    uni.redirectTo({ url: `/pages/order/result?id=${g.group_order_id}` });
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "支付失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding: 40rpx 24rpx;
}
.card {
  background: #fff;
  border-radius: 16rpx;
  padding: 32rpx;
  margin-bottom: 40rpx;
}
.label {
  color: #999;
  font-size: 24rpx;
}
.val {
  display: block;
  margin-top: 8rpx;
}
.price {
  display: block;
  margin-top: 24rpx;
  font-size: 48rpx;
  font-weight: 700;
  color: var(--qx-price, #e23030);
}
.tip {
  display: block;
  margin-top: 12rpx;
  color: #999;
  font-size: 24rpx;
}
.mt {
  margin-top: 20rpx;
}
</style>
