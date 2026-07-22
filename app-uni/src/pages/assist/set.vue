<template>
  <view class="page">
    <view v-if="set" class="hero">
      <text class="title">{{ set.store_name || "助力单" }}</text>
      <text class="sub">{{ statusText(set.status) }} · ¥{{ set.assist_price }}</text>
      <text class="meta">
        发起人 {{ set.nickname || "用户" }} · 进度 {{ set.yet_assist_count }}/{{ set.assist_count }}
      </text>
    </view>

    <view class="section">
      <text class="section-title">已助力好友</text>
      <view v-if="!(set?.helpers || []).length" class="empty">还没有人助力</view>
      <view v-for="h in set?.helpers || []" :key="h.uid" class="helper">
        {{ h.nickname || "好友" }}
      </view>
    </view>

    <view class="footer">
      <button v-if="set && set.status === 1" class="cta" type="primary" @click="onHelp">
        帮 TA 助力
      </button>
      <button v-else-if="set && set.status === 10 && isOwner" class="cta" type="primary" @click="onOrder">
        助力价下单
      </button>
      <view v-else-if="set && set.status === 20" class="done">已支付完成</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { onLoad, onShow } from "@dcloudio/uni-app";
import {
  assistCheck,
  assistCreate,
  fetchAssistSet,
  helpAssist,
  type AssistSet,
} from "@/api/assist";
import { fetchAddresses } from "@/api/cart";
import { useUserStore } from "@/stores/user";

const user = useUserStore();
const id = ref(0);
const set = ref<AssistSet | null>(null);
const isOwner = computed(() => !!set.value && set.value.uid === user.profile?.uid);

onLoad((q) => {
  id.value = Number(q?.id || 0);
});

onShow(async () => {
  if (!id.value) return;
  try {
    set.value = await fetchAssistSet(id.value);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function statusText(s: number) {
  if (s === 1) return "助力中";
  if (s === 10) return "已满员可下单";
  if (s === 20) return "已支付";
  if (s === -1) return "已失败";
  return `状态${s}`;
}

async function onHelp() {
  if (!user.isLogin) {
    uni.navigateTo({ url: "/pages/login/index" });
    return;
  }
  try {
    set.value = await helpAssist(id.value);
    uni.showToast({ title: "助力成功", icon: "none" });
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "助力失败", icon: "none" });
  }
}

async function onOrder() {
  if (!user.isLogin) {
    uni.navigateTo({ url: "/pages/login/index" });
    return;
  }
  try {
    const addrs = await fetchAddresses();
    const arr = addrs.list || [];
    if (!arr.length) {
      uni.showToast({ title: "请先添加收货地址", icon: "none" });
      uni.navigateTo({ url: "/pages/address/list" });
      return;
    }
    const def = arr.find((a) => a.is_default === 1) || arr[0];
    await assistCheck({
      product_assist_set_id: id.value,
      cart_num: 1,
      address_id: def.address_id,
    });
    const g = await assistCreate({
      product_assist_set_id: id.value,
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
  margin-bottom: 24rpx;
}
.title {
  font-size: 34rpx;
  font-weight: 700;
}
.sub {
  display: block;
  margin-top: 12rpx;
  color: #e23030;
  font-size: 26rpx;
}
.meta {
  display: block;
  margin-top: 10rpx;
  color: #888;
  font-size: 24rpx;
}
.section-title {
  font-size: 28rpx;
  font-weight: 600;
  margin-bottom: 16rpx;
  display: block;
}
.empty {
  color: #999;
  font-size: 24rpx;
  padding: 24rpx 0;
}
.helper {
  background: #fff;
  border-radius: 12rpx;
  padding: 20rpx 24rpx;
  margin-bottom: 12rpx;
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
  border-radius: 44rpx;
}
.done {
  text-align: center;
  color: #888;
  padding: 16rpx 0;
}
</style>
