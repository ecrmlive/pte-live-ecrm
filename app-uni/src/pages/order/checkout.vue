<template>
  <view class="page">
    <view class="card">
      <text class="title">收货地址</text>
      <view v-if="addr" class="addr">
        <text>{{ addr.real_name }} {{ addr.phone }}</text>
        <text class="line">{{ addr.province }}{{ addr.city }}{{ addr.district }} {{ addr.detail }}</text>
      </view>
      <text v-else class="line">暂无地址，请先新增</text>
      <text class="link" @click="goAddress">管理地址</text>
    </view>
    <view class="card">
      <text class="title">优惠券</text>
      <view v-for="c in usable" :key="c.coupon_user_id" class="coupon" @click="toggle(c.coupon_user_id)">
        <text class="check">{{ selected.includes(c.coupon_user_id) ? "✓" : "○" }}</text>
        <text>{{ c.coupon_title }} · ¥{{ Number(c.coupon_price).toFixed(2) }}</text>
      </view>
      <text v-if="!usable.length" class="line">暂无可用券</text>
      <text class="link" @click="goCenter">去领券</text>
    </view>
    <view class="card">
      <text class="title">积分抵扣</text>
      <text class="line">可用 {{ userIntegral }} · 最多抵应付 20%（100 积分=¥1）</text>
      <view class="row">
        <switch :checked="useIntegralOn" @change="onIntegralSwitch" color="#e23030" />
        <text class="line">使用积分抵扣</text>
      </view>
      <text v-if="useIntegralOn" class="line">
        预估抵 ¥{{ integralPrice }}（扣 {{ integralUsed }} 积分）
      </text>
    </view>
    <view class="card">
      <text>共 {{ cartIds.length }} 个购物车行</text>
      <text v-if="svipDiscount > 0" class="line">会员立减 ¥{{ svipDiscount.toFixed(2) }}</text>
      <text class="line">预估 ¥{{ payPreview }}（含优惠）</text>
    </view>
    <view class="foot qx-safe-bottom">
      <view class="qx-btn qx-btn-primary" @click="submit">提交并支付(Mock)</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad, onShow } from "@dcloudio/uni-app";
import { fetchAddresses, type Address } from "@/api/cart";
import { fetchUsableCoupons, type CouponUser } from "@/api/coupon";
import { payGroup, v2Check, v2Create } from "@/api/order";

const cartIds = ref<number[]>([]);
const addr = ref<Address | null>(null);
const payPreview = ref("0.00");
const usable = ref<CouponUser[]>([]);
const selected = ref<number[]>([]);
const useIntegralOn = ref(false);
const useIntegralWant = ref(0);
const userIntegral = ref(0);
const integralUsed = ref(0);
const integralPrice = ref("0.00");
const svipDiscount = ref(0);

async function refreshCheck() {
  if (!cartIds.value.length) return;
  try {
    const check = await v2Check(
      cartIds.value,
      selected.value,
      useIntegralOn.value ? useIntegralWant.value || userIntegral.value || 999999 : 0,
    );
    payPreview.value = Number(check.pay_price || 0).toFixed(2);
    userIntegral.value = check.user_integral || 0;
    integralUsed.value = check.integral || 0;
    integralPrice.value = Number(check.integral_price || 0).toFixed(2);
    svipDiscount.value = Number(check.svip_discount || 0);
    if (!useIntegralWant.value && userIntegral.value > 0) {
      useIntegralWant.value = userIntegral.value;
    }
  } catch {
    /* ignore */
  }
}

onLoad(async (q) => {
  cartIds.value = String(q?.cart_ids || "")
    .split(",")
    .map((x) => Number(x))
    .filter((x) => x > 0);
  if (cartIds.value.length) {
    try {
      const u = await fetchUsableCoupons(cartIds.value);
      usable.value = u.list || [];
    } catch {
      /* ignore */
    }
    await refreshCheck();
  }
});

onShow(async () => {
  try {
    const res = await fetchAddresses();
    addr.value = (res.list || []).find((a) => a.is_default === 1) || res.list?.[0] || null;
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "地址加载失败", icon: "none" });
  }
});

function toggle(id: number) {
  const i = selected.value.indexOf(id);
  if (i >= 0) selected.value = selected.value.filter((x) => x !== id);
  else selected.value = [...selected.value, id];
  void refreshCheck();
}

function onIntegralSwitch(e: { detail: { value: boolean } }) {
  useIntegralOn.value = !!e.detail.value;
  void refreshCheck();
}

function goAddress() {
  uni.navigateTo({ url: "/pages/address/list" });
}

function goCenter() {
  uni.navigateTo({ url: "/pages/coupon/center" });
}

async function submit() {
  if (!addr.value || !cartIds.value.length) {
    uni.showToast({ title: "缺少地址或商品", icon: "none" });
    return;
  }
  try {
    const g = await v2Create({
      cart_ids: cartIds.value,
      address_id: addr.value.address_id,
      coupon_user_ids: selected.value,
      use_integral: useIntegralOn.value ? useIntegralWant.value || userIntegral.value : 0,
    });
    await payGroup(g.group_order_id, "mock");
    uni.redirectTo({ url: `/pages/order/result?id=${g.group_order_id}` });
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "下单失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding-bottom: 140rpx;
}
.card {
  margin: 20rpx 24rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
}
.title {
  font-weight: 700;
  display: block;
  margin-bottom: 12rpx;
}
.line {
  display: block;
  color: var(--qx-text-secondary);
  margin-top: 8rpx;
  font-size: 24rpx;
}
.link {
  display: block;
  margin-top: 16rpx;
  color: var(--qx-primary, #e23030);
  font-size: 24rpx;
}
.coupon {
  display: flex;
  gap: 12rpx;
  align-items: center;
  padding: 12rpx 0;
  font-size: 26rpx;
}
.check {
  width: 36rpx;
  color: var(--qx-primary, #e23030);
}
.row {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-top: 12rpx;
}
.foot {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 20rpx 32rpx;
  background: #fff;
}
</style>
