<template>
  <view class="page">
    <view v-if="active" class="hero">
      <text class="title">{{ active.store_name }}</text>
      <text class="sub">需 {{ active.assist_count }} 人助力 · 助力价 ¥{{ active.assist_price }}</text>
    </view>

    <view class="section">
      <text class="section-title">进行中的助力</text>
      <view v-if="!sets.length" class="empty">暂无进行中，可发起新助力</view>
      <view
        v-for="s in sets"
        :key="s.product_assist_set_id"
        class="buying"
        @click="goSet(s.product_assist_set_id)"
      >
        <view>
          <text>已助力 {{ s.yet_assist_count }}/{{ s.assist_count }}</text>
          <text class="status"> · {{ statusText(s.status) }}</text>
        </view>
        <view class="actions">
          <text v-if="s.status === 1" class="join" @click.stop="goSet(s.product_assist_set_id)">帮TA</text>
          <text v-if="s.status === 10" class="join" @click.stop="order(s.product_assist_set_id)">去下单</text>
        </view>
      </view>
    </view>

    <view class="footer">
      <button class="cta" @click="start()">发起助力</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad, onShow } from "@dcloudio/uni-app";
import {
  assistCheck,
  assistCreate,
  fetchAssist,
  fetchSets,
  startAssist,
  type AssistSet,
  type ProductAssist,
} from "@/api/assist";
import { fetchAddresses } from "@/api/cart";

const id = ref(0);
const active = ref<ProductAssist | null>(null);
const sets = ref<AssistSet[]>([]);

onLoad((q) => {
  id.value = Number(q?.id || 0);
});

onShow(async () => {
  if (!id.value) return;
  try {
    const [a, ss] = await Promise.all([fetchAssist(id.value), fetchSets(id.value)]);
    active.value = a;
    sets.value = ss.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function statusText(s: number) {
  if (s === 1) return "进行中";
  if (s === 10) return "可下单";
  if (s === 20) return "已支付";
  return "已结束";
}

function goSet(setID: number) {
  uni.navigateTo({ url: `/pages/assist/set?id=${setID}` });
}

async function start() {
  try {
    const row = await startAssist(id.value);
    if (row.status === 10) {
      await order(row.product_assist_set_id);
      return;
    }
    uni.navigateTo({ url: `/pages/assist/set?id=${row.product_assist_set_id}` });
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "发起失败", icon: "none" });
  }
}

async function pickAddressID(): Promise<number> {
  const data = await fetchAddresses();
  const arr = data.list || [];
  if (!arr.length) {
    uni.showToast({ title: "请先添加收货地址", icon: "none" });
    uni.navigateTo({ url: "/pages/address/list" });
    return 0;
  }
  const def = arr.find((a) => a.is_default === 1) || arr[0];
  return def.address_id;
}

async function order(setID: number) {
  const addressID = await pickAddressID();
  if (!addressID) return;
  try {
    await assistCheck({
      product_assist_set_id: setID,
      cart_num: 1,
      address_id: addressID,
    });
    const g = await assistCreate({
      product_assist_set_id: setID,
      cart_num: 1,
      address_id: addressID,
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
.buying {
  background: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 26rpx;
}
.status {
  color: #888;
  font-size: 22rpx;
}
.actions {
  display: flex;
  gap: 16rpx;
}
.join {
  color: #e23030;
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
</style>
