<template>
  <view class="page">
    <view class="card">
      <text class="title">分销中心</text>
      <text class="line">我的 UID：{{ me?.uid ?? "-" }}</text>
      <text class="line">上级推广员：{{ me?.spread_uid || "未绑定" }}</text>
      <text class="line">是否推广员：{{ me?.is_promoter ? "是" : "否" }}</text>
      <text class="line">下级人数：{{ me?.spread_count ?? 0 }}</text>
    </view>
    <view v-if="me && !me.spread_uid" class="card">
      <text class="title">绑定上级</text>
      <input v-model="spreadUID" class="input" type="number" placeholder="推广员 UID" />
      <view class="qx-btn qx-btn-primary btn" @click="onBind">绑定</view>
    </view>
    <view class="card">
      <text class="title">佣金流水</text>
      <view v-for="b in bills" :key="b.bill_id" class="row">
        <text>{{ b.title || "佣金" }}</text>
        <text class="amt">+¥{{ Number(b.number).toFixed(2) }}</text>
      </view>
      <view v-if="!bills.length" class="empty">暂无流水</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import {
  bindSpread,
  fetchSpreadBills,
  fetchSpreadMe,
  type SpreadBill,
  type SpreadMe,
} from "@/api/spread";

const me = ref<SpreadMe | null>(null);
const bills = ref<SpreadBill[]>([]);
const spreadUID = ref("");

async function load() {
  try {
    me.value = await fetchSpreadMe();
    const res = await fetchSpreadBills();
    bills.value = res.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
}

onShow(() => void load());

async function onBind() {
  const uid = Number(spreadUID.value);
  if (!uid) {
    uni.showToast({ title: "请输入推广员 UID", icon: "none" });
    return;
  }
  try {
    await bindSpread(uid);
    uni.showToast({ title: "绑定成功", icon: "success" });
    await load();
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "绑定失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding: 20rpx 0;
}
.card {
  margin: 0 24rpx 20rpx;
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
  margin-top: 8rpx;
  font-size: 26rpx;
  color: #666;
}
.input {
  margin-top: 16rpx;
  background: #f7f7f7;
  border-radius: 12rpx;
  padding: 20rpx;
}
.btn {
  margin-top: 20rpx;
}
.row {
  display: flex;
  justify-content: space-between;
  padding: 16rpx 0;
  border-top: 1px solid #f0f0f0;
  font-size: 26rpx;
}
.amt {
  color: var(--qx-price, #e23030);
  font-weight: 600;
}
.empty {
  color: #999;
  padding: 24rpx 0;
  text-align: center;
}
</style>
