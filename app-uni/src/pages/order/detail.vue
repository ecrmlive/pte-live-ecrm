<template>
  <view class="page" v-if="order">
    <view class="card">
      <text class="st">{{ order.paid ? "已支付" : "待支付" }}</text>
      <text class="line">主单号 {{ order.group_order_sn }}</text>
      <text class="line">应付 ¥{{ Number(order.pay_price).toFixed(2) }}</text>
      <text class="line">{{ order.real_name }} {{ order.user_phone }}</text>
      <text class="line muted">{{ order.user_address }}</text>
    </view>
    <view v-for="o in order.orders || []" :key="o.order_id" class="card">
      <text class="mer">{{ o.mer_name || `商户${o.mer_id}` }} · {{ statusText(o.status, o.paid) }}</text>
      <view v-for="p in o.products || []" :key="p.product_id + '-' + p.product_num" class="item">
        <text>{{ productName(p) }} ×{{ p.product_num }}</text>
        <text>¥{{ Number(p.total_price).toFixed(2) }}</text>
      </view>
      <text v-if="o.delivery_id" class="line muted">物流 {{ o.delivery_name }} {{ o.delivery_id }}</text>
      <view
        v-if="o.paid === 1 && o.status !== -1"
        class="qx-btn qx-btn-ghost refund-btn"
        @click="goRefund(o.order_id)"
      >
        申请仅退款
      </view>
    </view>
    <view v-if="!order.paid" class="foot">
      <view class="qx-btn qx-btn-primary" @click="pay">立即支付(Mock)</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { fetchOrderDetail, payGroup, type GroupOrder } from "@/api/order";

const order = ref<GroupOrder | null>(null);
const id = ref(0);

onLoad(async (q) => {
  id.value = Number(q?.id || 0);
  await load();
});

async function load() {
  if (!id.value) return;
  try {
    order.value = await fetchOrderDetail(id.value);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
}

function statusText(status: number, paid: number) {
  if (!paid) return "待支付";
  if (status === 0) return "待发货";
  if (status === 1) return "待收货";
  if (status === 3) return "已完成";
  if (status === 9) return "拼团中";
  if (status === 10) return "待付尾款";
  if (status === 11) return "尾款超时";
  return `状态${status}`;
}

function productName(p: { product_info?: string; product_id: number }) {
  try {
    const info = JSON.parse(p.product_info || "{}");
    return info.store_name || `商品${p.product_id}`;
  } catch {
    return `商品${p.product_id}`;
  }
}

function goRefund(orderId: number) {
  uni.navigateTo({ url: `/pages/refund/apply?order_id=${orderId}` });
}

async function pay() {
  try {
    order.value = await payGroup(id.value, "mock");
    uni.showToast({ title: "支付成功", icon: "success" });
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "支付失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding: 20rpx 24rpx 140rpx;
}
.card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
}
.st {
  display: block;
  font-size: 34rpx;
  font-weight: 700;
  margin-bottom: 12rpx;
}
.mer {
  font-weight: 700;
}
.line {
  display: block;
  margin-top: 8rpx;
  font-size: 26rpx;
}
.muted {
  color: #999;
}
.item {
  display: flex;
  justify-content: space-between;
  margin-top: 12rpx;
  font-size: 26rpx;
}
.refund-btn {
  margin-top: 20rpx;
  height: 72rpx;
  line-height: 72rpx;
  font-size: 26rpx;
}
.foot {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 16rpx 24rpx;
  background: #fff;
}
</style>
