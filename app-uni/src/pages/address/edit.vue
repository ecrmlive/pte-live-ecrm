<template>
  <view class="page">
    <view class="card">
      <view class="field"><text class="label">收货人</text><input v-model="form.real_name" placeholder="姓名" /></view>
      <view class="field"><text class="label">手机号</text><input v-model="form.phone" type="number" placeholder="手机号" /></view>
      <view class="field"><text class="label">省</text><input v-model="form.province" placeholder="省" /></view>
      <view class="field"><text class="label">市</text><input v-model="form.city" placeholder="市" /></view>
      <view class="field"><text class="label">区</text><input v-model="form.district" placeholder="区" /></view>
      <view class="field"><text class="label">详细</text><input v-model="form.detail" placeholder="街道门牌" /></view>
      <view class="field switch">
        <text class="label">默认地址</text>
        <switch :checked="form.is_default === 1" @change="onDefault" />
      </view>
    </view>
    <view class="qx-btn qx-btn-primary save" @click="save">保存</view>
  </view>
</template>

<script setup lang="ts">
import { reactive, ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import { createAddress, fetchAddresses, updateAddress } from "@/api/cart";

const id = ref(0);
const form = reactive({
  real_name: "",
  phone: "",
  province: "广东省",
  city: "深圳市",
  district: "南山区",
  detail: "",
  is_default: 0 as number,
});

onLoad(async (q) => {
  id.value = Number(q?.id || 0);
  if (!id.value) return;
  try {
    const data = await fetchAddresses();
    const a = (data.list || []).find((x) => x.address_id === id.value);
    if (a) {
      form.real_name = a.real_name;
      form.phone = a.phone;
      form.province = a.province;
      form.city = a.city;
      form.district = a.district;
      form.detail = a.detail;
      form.is_default = a.is_default;
    }
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "加载失败", icon: "none" });
  }
});

function onDefault(e: { detail: { value: boolean } }) {
  form.is_default = e.detail.value ? 1 : 0;
}

async function save() {
  try {
    if (id.value) {
      await updateAddress(id.value, { ...form, is_default: form.is_default });
    } else {
      await createAddress({ ...form, is_default: form.is_default });
    }
    uni.showToast({ title: "已保存", icon: "success" });
    setTimeout(() => uni.navigateBack(), 400);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "保存失败", icon: "none" });
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
  padding: 24rpx;
}
.card {
  background: #fff;
  border-radius: 16rpx;
  padding: 8rpx 24rpx;
}
.field {
  display: flex;
  align-items: center;
  gap: 20rpx;
  min-height: 88rpx;
  border-bottom: 1px solid #f2f2f2;
}
.label {
  width: 140rpx;
  color: #666;
}
.switch {
  justify-content: space-between;
}
.save {
  margin-top: 40rpx;
}
</style>
