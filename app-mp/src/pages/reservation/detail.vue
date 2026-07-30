<template>
  <view class="page" v-if="product">
    <view class="hero card">
      <text class="title">{{ product.store_name }}</text>
      <text class="sub">¥{{ product.price }} · {{ product.mer_name }}</text>
      <text class="hint">到店服务，支付后生成核销码</text>
    </view>

    <view class="card">
      <text class="label">预约日期</text>
      <scroll-view scroll-x class="days">
        <view
          v-for="d in days"
          :key="d"
          class="day"
          :class="{ on: d === date }"
          @click="pickDate(d)"
        >
          {{ d.slice(5) }}
        </view>
      </scroll-view>
    </view>

    <view class="card">
      <text class="label">可选时段</text>
      <view v-if="!slots.length" class="empty">当日暂无可约</view>
      <view
        v-for="s in slots"
        :key="s.attr_reservation_id"
        class="slot"
        :class="{ on: s.attr_reservation_id === slotId, disabled: s.remain <= 0 }"
        @click="pickSlot(s)"
      >
        <text>{{ s.label || `${s.start_time}-${s.end_time}` }}</text>
        <text class="remain">余 {{ s.remain }}</text>
      </view>
    </view>

    <view class="foot qx-safe-bottom">
      <view class="qx-btn qx-btn-primary" @click="submit">确认预约</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { fetchAddresses } from "@/api/cart";
import {
  fetchDaySlots,
  fetchReserveProduct,
  reservationCheck,
  reservationCreate,
  type ReserveProduct,
  type SlotDay,
} from "@/api/reservation";
import { useUserStore } from "@/stores/user";

const user = useUserStore();
const productId = ref(0);
const product = ref<ReserveProduct | null>(null);
const days = ref<string[]>([]);
const date = ref("");
const slots = ref<SlotDay[]>([]);
const slotId = ref(0);

function buildDays(n: number) {
  const out: string[] = [];
  const base = new Date();
  for (let i = 0; i < n; i++) {
    const d = new Date(base);
    d.setDate(base.getDate() + i);
    const y = d.getFullYear();
    const m = String(d.getMonth() + 1).padStart(2, "0");
    const day = String(d.getDate()).padStart(2, "0");
    out.push(`${y}-${m}-${day}`);
  }
  return out;
}

async function loadSlots() {
  if (!productId.value || !date.value) return;
  const data = await fetchDaySlots(productId.value, date.value);
  slots.value = data.list || [];
  slotId.value = 0;
}

async function pickDate(d: string) {
  date.value = d;
  await loadSlots();
}

function pickSlot(s: SlotDay) {
  if (s.remain <= 0) return;
  slotId.value = s.attr_reservation_id;
}

onLoad(async (q) => {
  productId.value = Number(q?.id || 0);
  if (!productId.value) return;
  try {
    product.value = await fetchReserveProduct(productId.value);
    const n = Math.max(1, product.value.show_reservation_days || 7);
    days.value = buildDays(n);
    date.value = days.value[0] || "";
    await loadSlots();
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

async function submit() {
  if (!user.isLogin) {
    uni.navigateTo({ url: "/pages/login/index" });
    return;
  }
  if (!slotId.value || !date.value) {
    uni.showToast({ title: "请选择日期与时段", icon: "none" });
    return;
  }
  try {
    const addrRes = await fetchAddresses();
    const arr = addrRes.list || [];
    if (!arr.length) {
      uni.showToast({ title: "请先添加地址", icon: "none" });
      uni.navigateTo({ url: "/pages/address/list" });
      return;
    }
    const addr = arr.find((a) => a.is_default === 1) || arr[0];
    await reservationCheck({
      product_id: productId.value,
      slot_id: slotId.value,
      date: date.value,
      address_id: addr.address_id,
    });
    const g = await reservationCreate({
      product_id: productId.value,
      slot_id: slotId.value,
      date: date.value,
      address_id: addr.address_id,
    });
    uni.navigateTo({ url: `/pages/order/pay?id=${g.group_order_id}` });
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "预约失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  padding: 24rpx 24rpx 160rpx;
  background: #f5f5f5;
}
.card {
  background: #fff;
  border-radius: 12rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
}
.title {
  font-size: 34rpx;
  font-weight: 700;
}
.sub {
  display: block;
  margin-top: 10rpx;
  color: #e23030;
}
.hint {
  display: block;
  margin-top: 8rpx;
  font-size: 22rpx;
  color: #999;
}
.label {
  font-size: 28rpx;
  font-weight: 600;
  display: block;
  margin-bottom: 16rpx;
}
.days {
  white-space: nowrap;
}
.day {
  display: inline-block;
  padding: 16rpx 24rpx;
  margin-right: 12rpx;
  border-radius: 8rpx;
  background: #f5f5f5;
  font-size: 26rpx;
}
.day.on {
  background: #fff1f0;
  color: #e23030;
  font-weight: 600;
}
.slot {
  display: flex;
  justify-content: space-between;
  padding: 22rpx 16rpx;
  border-radius: 8rpx;
  margin-bottom: 12rpx;
  background: #fafafa;
}
.slot.on {
  background: #fff1f0;
  color: #e23030;
}
.slot.disabled {
  opacity: 0.4;
}
.remain {
  color: #999;
  font-size: 24rpx;
}
.empty {
  color: #999;
  padding: 24rpx 0;
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
