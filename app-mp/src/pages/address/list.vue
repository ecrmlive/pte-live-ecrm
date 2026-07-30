<template>
  <view class="page">
    <view v-if="!list.length" class="empty">暂无地址</view>
    <view v-for="a in list" :key="a.address_id" class="card" @click="onSelect(a)">
      <view class="top">
        <text class="name">{{ a.real_name }} {{ a.phone }}</text>
        <text v-if="a.is_default" class="tag">默认</text>
      </view>
      <text class="line">{{ a.province }}{{ a.city }}{{ a.district }}{{ a.detail }}</text>
      <view class="ops" @click.stop>
        <text @click="edit(a)">编辑</text>
        <text class="danger" @click="remove(a.address_id)">删除</text>
      </view>
    </view>
    <view class="foot qx-safe-bottom">
      <view class="qx-btn qx-btn-primary" @click="edit()">新增地址</view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { onLoad, onShow } from "@dcloudio/uni-app";
import { deleteAddress, fetchAddresses, type Address } from "@/api/cart";

const list = ref<Address[]>([]);
const selectMode = ref(false);

onLoad((q) => {
  selectMode.value = q?.select === "1";
});

async function load() {
  try {
    const data = await fetchAddresses();
    list.value = data.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
}

onShow(load);

function onSelect(a: Address) {
  if (!selectMode.value) return;
  const pages = getCurrentPages();
  const prev = pages[pages.length - 2] as { $vm?: { address?: Address } };
  if (prev?.$vm) {
    prev.$vm.address = a;
  }
  uni.$emit("address-selected", a);
  uni.navigateBack();
}

function edit(a?: Address) {
  const q = a ? `?id=${a.address_id}` : "";
  uni.navigateTo({ url: `/pages/address/edit${q}` });
}

async function remove(id: number) {
  const ok = await new Promise<boolean>((resolve) => {
    uni.showModal({
      title: "删除地址",
      content: "确认删除？",
      success: (r) => resolve(!!r.confirm),
    });
  });
  if (!ok) return;
  try {
    await deleteAddress(id);
    await load();
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "删除失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding: 20rpx 24rpx 160rpx;
}
.empty {
  text-align: center;
  padding: 120rpx 0;
  color: #999;
}
.card {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
}
.top {
  display: flex;
  align-items: center;
  gap: 12rpx;
}
.name {
  font-weight: 700;
}
.tag {
  font-size: 20rpx;
  color: var(--qx-primary, #e23030);
  border: 1px solid currentColor;
  padding: 0 8rpx;
  border-radius: 6rpx;
}
.line {
  display: block;
  margin-top: 12rpx;
  color: var(--qx-text-secondary);
  font-size: 24rpx;
}
.ops {
  margin-top: 16rpx;
  display: flex;
  gap: 32rpx;
  color: #666;
  font-size: 24rpx;
}
.danger {
  color: #e23030;
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
