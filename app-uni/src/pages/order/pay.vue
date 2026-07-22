<template>
  <view class="page">
    <view class="card">
      <text class="label">主单号</text>
      <text class="val">{{ group?.group_order_sn || "-" }}</text>
      <text class="price">¥{{ Number(group?.pay_price || 0).toFixed(2) }}</text>
      <text class="tip">子单数 {{ group?.orders?.length || 0 }}（多商户一单）</text>
    </view>
    <view class="qx-btn qx-btn-primary" @click="pay('mock')">Mock 支付</view>
    <view class="qx-btn qx-btn-ghost mt" @click="pay('balance')">余额支付</view>
    <view class="qx-btn qx-btn-ghost mt" @click="pay('wechat')">微信支付（沙箱）</view>
    <view class="qx-btn qx-btn-ghost mt" @click="pay('alipay')">支付宝（沙箱）</view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import {
  fetchOrderDetail,
  notifyChannelPay,
  payGroup,
  type GroupOrder,
  type PayIntent,
  type PayType,
} from "@/api/order";

const id = ref(0);
const group = ref<GroupOrder | null>(null);

onLoad(async (q) => {
  id.value = Number(q?.id || 0);
  if (!id.value) return;
  try {
    group.value = await fetchOrderDetail(id.value);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function isPayIntent(v: GroupOrder | PayIntent): v is PayIntent {
  return typeof (v as PayIntent).status === "string" && typeof (v as PayIntent).channel === "string";
}

async function pay(type: PayType) {
  try {
    const res = await payGroup(id.value, type);
    let g: GroupOrder;
    if (isPayIntent(res)) {
      if (res.status === "paid") {
        g = await fetchOrderDetail(id.value);
      } else {
        if (!res.notify_token) {
          uni.showToast({ title: "非沙箱环境请完成第三方支付", icon: "none" });
          return;
        }
        const uid = group.value?.uid || 0;
        if (!uid) {
          uni.showToast({ title: "缺少用户信息", icon: "none" });
          return;
        }
        g = await notifyChannelPay(type as "wechat" | "alipay", {
          group_order_id: res.group_order_id,
          uid,
          out_trade_no: res.out_trade_no,
          pay_price: res.pay_price,
          notify_token: res.notify_token,
        });
      }
    } else {
      g = res;
    }
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
