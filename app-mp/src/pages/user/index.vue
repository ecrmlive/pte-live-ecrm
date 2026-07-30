<template>
  <view class="page">
    <view class="nav" :style="{ paddingTop: statusBar + 'px' }">
      <view class="nav-inner">
        <text class="nav-title">我的</text>
      </view>
    </view>

    <view class="body" :style="{ paddingTop: navPad + 'px' }">
      <view class="card profile" @click="onProfileTap">
        <view class="avatar">
          <text class="avatar-text">{{ avatarText }}</text>
        </view>
        <view class="meta">
          <text class="name">{{ user.displayName }}</text>
          <text class="sub">{{ user.isLogin ? tipLogged : "点击登录 / 注册" }}</text>
        </view>
      </view>

      <view class="card menu">
        <view class="row" @click="goPath('/pages/order/list')">
          <text>我的订单</text>
          <text class="arrow">›</text>
        </view>
        <view class="row" @click="goPath('/pages/presell/finals')">
          <text>待付尾款</text>
          <text class="arrow">›</text>
        </view>
        <view class="row" @click="goPath('/pages/assist/list')">
          <text>好友助力</text>
          <text class="arrow">›</text>
        </view>
        <view class="row" @click="goPath('/pages/refund/list')">
          <text>我的售后</text>
          <text class="arrow">›</text>
        </view>
        <view class="row" @click="goPath('/pages/address/list')">
          <text>收货地址</text>
          <text class="arrow">›</text>
        </view>
        <view class="row" @click="goPath('/pages/coupon/mine')">
          <text>我的优惠券</text>
          <text class="arrow">›</text>
        </view>
        <view class="row" @click="goPath('/pages/coupon/center')">
          <text>领券中心</text>
          <text class="arrow">›</text>
        </view>
        <view class="row" @click="goPath('/pages/spread/index')">
          <text>分销中心</text>
          <text class="arrow">›</text>
        </view>
        <view class="row" @click="goPath('/pages/points/list')">
          <text>积分商城</text>
          <text class="arrow">›</text>
        </view>
        <view class="row" @click="goCommunity">
          <text>社区种草</text>
          <text class="arrow">›</text>
        </view>
        <view class="row" @click="goCreatePost">
          <text>发帖</text>
          <text class="arrow">›</text>
        </view>
        <view class="row" @click="goPath('/pages/notice/list')">
          <text>平台公告</text>
          <text class="arrow">›</text>
        </view>
        <view class="row" @click="goPath('/pages/agreement/detail?key=sys_user_agree')">
          <text>用户协议</text>
          <text class="arrow">›</text>
        </view>
        <view class="row" @click="goPath('/pages/agreement/detail?key=sys_userr_privacy')">
          <text>隐私政策</text>
          <text class="arrow">›</text>
        </view>
      </view>

      <view v-if="user.isLogin" class="qx-btn qx-btn-ghost logout" @click="logout">退出登录</view>
      <text class="phase">阶段 7 · 公告 / 协议</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { useUserStore } from "@/stores/user";

const user = useUserStore();
const statusBar = ref(20);
const navPad = computed(() => statusBar.value + 44);

const avatarText = computed(() => {
  const n = user.displayName;
  return n ? n.slice(0, 1) : "游";
});

const tipLogged = computed(() => {
  const phone = user.profile?.phone;
  const pts = user.profile?.integral ?? 0;
  const base = phone ? `手机 ${phone}` : `账号 ${user.profile?.account || ""}`;
  const vip = user.profile?.is_svip_active ? " · SVIP" : "";
  return `${base} · 积分 ${pts}${vip}`;
});

onShow(() => {
  const sys = uni.getSystemInfoSync();
  statusBar.value = sys.statusBarHeight || 20;
  if (user.isLogin) {
    user.refreshMe();
  }
});

function onProfileTap() {
  if (!user.isLogin) {
    uni.navigateTo({ url: "/pages/login/index" });
  }
}

function goPath(url: string) {
  if (!user.isLogin) {
    uni.navigateTo({ url: "/pages/login/index" });
    return;
  }
  uni.navigateTo({ url });
}

function goCommunity() {
  uni.navigateTo({ url: "/pages/community/list" });
}

function goCreatePost() {
  if (!user.isLogin) {
    uni.navigateTo({ url: "/pages/login/index" });
    return;
  }
  uni.navigateTo({ url: "/pages/community/create" });
}

function logout() {
  uni.showModal({
    title: "退出登录",
    content: "确定退出当前账号？",
    success: (res) => {
      if (res.confirm) {
        user.logout();
        uni.showToast({ title: "已退出", icon: "none" });
      }
    },
  });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
}

.nav {
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  z-index: 20;
  background: #fff;
}

.nav-inner {
  height: 44px;
  display: flex;
  align-items: center;
  padding: 0 24rpx;
}

.nav-title {
  font-size: 34rpx;
  font-weight: 700;
}

.body {
  padding: 24rpx 24rpx 40rpx;
}

.card {
  background: var(--qx-card);
  border-radius: 16rpx;
  margin-bottom: 24rpx;
}

.profile {
  display: flex;
  align-items: center;
  gap: 24rpx;
  padding: 36rpx 28rpx;
}

.avatar {
  width: 112rpx;
  height: 112rpx;
  border-radius: 56rpx;
  background: linear-gradient(135deg, #e23030, #f06a3d);
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-text {
  color: #fff;
  font-size: 44rpx;
  font-weight: 700;
}

.meta {
  flex: 1;
}

.name {
  display: block;
  font-size: 36rpx;
  font-weight: 700;
}

.sub {
  display: block;
  margin-top: 8rpx;
  color: var(--qx-text-secondary);
  font-size: 24rpx;
}

.menu .row {
  height: 96rpx;
  padding: 0 28rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--qx-line);
  font-size: 30rpx;
}

.menu .row:last-child {
  border-bottom: none;
}

.arrow {
  color: #ccc;
  font-size: 40rpx;
}

.logout {
  margin-top: 12rpx;
}

.phase {
  display: block;
  margin-top: 36rpx;
  text-align: center;
  color: #bbb;
  font-size: 22rpx;
}
</style>
