<template>
  <view class="page">
    <view v-if="group" class="hero">
      <text class="title">{{ group.store_name }}</text>
      <text class="sub">{{ group.buying_count_num }}人团 · 拼团价 ¥{{ group.price }}</text>
    </view>

    <view class="section">
      <text class="section-title">进行中的团</text>
      <view v-if="!buyings.length" class="empty">暂无进行中的团，可开新团</view>
      <view v-for="b in buyings" :key="b.group_buying_id" class="buying">
        <text>还差 {{ b.remain }} 人 · 已有 {{ b.yet_buying_num }}/{{ b.buying_count_num }}</text>
        <text class="join" @click="join(b.group_buying_id)">参团</text>
      </view>
    </view>

    <view class="footer">
      <button class="cta" @click="open()">我要开团</button>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad, onShow } from "@dcloudio/uni-app";
import {
  fetchBuyings,
  fetchGroup,
  groupCheck,
  groupCreate,
  type Buying,
  type ProductGroup,
} from "@/api/combination";
import { fetchAddresses } from "@/api/cart";

const id = ref(0);
const group = ref<ProductGroup | null>(null);
const buyings = ref<Buying[]>([]);

onLoad((q) => {
  id.value = Number(q?.id || 0);
});

onShow(async () => {
  if (!id.value) return;
  try {
    const [g, bs] = await Promise.all([fetchGroup(id.value), fetchBuyings(id.value)]);
    group.value = g;
    buyings.value = bs.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

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

async function submit(buyingID: number) {
  const addressID = await pickAddressID();
  if (!addressID) return;
  try {
    await groupCheck({
      product_group_id: id.value,
      group_buying_id: buyingID,
      cart_num: 1,
      address_id: addressID,
    });
    const g = await groupCreate({
      product_group_id: id.value,
      group_buying_id: buyingID,
      cart_num: 1,
      address_id: addressID,
    });
    uni.navigateTo({ url: `/pages/order/pay?id=${g.group_order_id}` });
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "下单失败", icon: "none" });
  }
}

function open() {
  void submit(0);
}
function join(buyingID: number) {
  void submit(buyingID);
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
