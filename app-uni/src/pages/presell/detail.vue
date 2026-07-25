<template>
  <view class="page">
    <view v-if="active" class="hero">
      <text class="title">{{ active.store_name }}</text>
      <text class="sub">
        {{ active.presell_type === 2 ? "定金预售" : "全款预售" }} · ¥{{ active.price }} · 库存
        {{ active.stock }}
      </text>
      <text v-if="active.presell_type === 2" class="sub">
        定金 ¥{{ active.down_price }} · 尾款 ¥{{ active.final_price }}
      </text>
      <text v-if="active.delivery_day" class="sub">约 {{ active.delivery_day }} 天内发货</text>
    </view>
    <view class="footer">
      <button class="cta" @click="buy">
        {{ active?.presell_type === 2 ? "支付定金" : "立即预订" }}
      </button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { fetchAddresses } from "@/api/cart";
import { fetchPresell, presellCheck, presellCreate, type PresellActive } from "@/api/presell";
import { useUserStore } from "@/stores/user";

const user = useUserStore();
const id = ref(0);
const active = ref<PresellActive | null>(null);

onLoad(async (q) => {
  id.value = Number(q?.id || 0);
  if (!id.value) return;
  try {
    active.value = await fetchPresell(id.value);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

async function buy() {
  if (!user.isLogin) {
    uni.navigateTo({ url: "/pages/login/index" });
    return;
  }
  try {
    const data = await fetchAddresses();
    const arr = data.list || [];
    if (!arr.length) {
      uni.showToast({ title: "请先添加收货地址", icon: "none" });
      uni.navigateTo({ url: "/pages/address/list" });
      return;
    }
    const def = arr.find((a) => a.is_default === 1) || arr[0];
    await presellCheck({ product_presell_id: id.value, cart_num: 1, address_id: def.address_id });
    const g = await presellCreate({
      product_presell_id: id.value,
      cart_num: 1,
      address_id: def.address_id,
    });
    uni.navigateTo({ url: `/pages/order/pay?id=${g.group_order_id}` });
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "下单失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  padding: 24rpx 24rpx 160rpx;
  background: #f5f5f5;
}
.hero {
  background: #fff;
  border-radius: 12rpx;
  padding: 32rpx;
}
.title {
  font-size: 34rpx;
  font-weight: 700;
}
.sub {
  display: block;
  margin-top: 12rpx;
  color: #666;
  font-size: 26rpx;
}
.footer {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 20rpx 24rpx calc(20rpx + env(safe-area-inset-bottom));
  background: #fff;
}
.cta {
  background: #e23030;
  color: #fff;
}
</style>
