<template>
  <view class="page">
    <view class="hero">
      <text class="brand">栖息商城</text>
      <text class="sub">多商户商城 · 用户登录</text>
    </view>

    <view class="form">
      <view class="field">
        <text class="label">账号</text>
        <input
          v-model="account"
          class="input"
          type="text"
          maxlength="32"
          placeholder="请输入账号"
          confirm-type="next"
        />
      </view>
      <view class="field">
        <text class="label">密码</text>
        <input
          v-model="password"
          class="input"
          password
          maxlength="64"
          placeholder="至少 6 位"
          confirm-type="done"
          @confirm="submit"
        />
      </view>
      <view v-if="mode === 'register'" class="field">
        <text class="label">昵称</text>
        <input
          v-model="nickname"
          class="input"
          type="text"
          maxlength="16"
          placeholder="可选，默认与账号相同"
        />
      </view>

      <view class="qx-btn qx-btn-primary btn" :class="{ disabled: loading }" @click="submit">
        {{ mode === "login" ? "登录" : "注册并登录" }}
      </view>
      <view class="switch" @click="toggleMode">
        <text>{{ mode === "login" ? "没有账号？去注册" : "已有账号？去登录" }}</text>
      </view>
      <text class="hint">演示账号 demo / admin123</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useUserStore } from "@/stores/user";

const user = useUserStore();
const mode = ref<"login" | "register">("login");
const account = ref("demo");
const password = ref("admin123");
const nickname = ref("");
const loading = ref(false);

function toggleMode() {
  mode.value = mode.value === "login" ? "register" : "login";
}

async function submit() {
  if (loading.value) return;
  const a = account.value.trim();
  const p = password.value;
  if (!a || p.length < 6) {
    uni.showToast({ title: "请填写账号且密码至少 6 位", icon: "none" });
    return;
  }
  loading.value = true;
  try {
    if (mode.value === "login") {
      await user.login(a, p);
    } else {
      await user.register(a, p, nickname.value.trim() || undefined);
    }
    uni.showToast({ title: "登录成功", icon: "success" });
    setTimeout(() => {
      uni.switchTab({ url: "/pages/user/index" });
    }, 300);
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "失败", icon: "none" });
  } finally {
    loading.value = false;
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  padding: 80rpx 48rpx 60rpx;
  background: linear-gradient(180deg, #fff5f5 0%, var(--qx-bg) 42%);
  box-sizing: border-box;
}

.hero {
  margin-bottom: 64rpx;
}

.brand {
  display: block;
  font-size: 56rpx;
  font-weight: 700;
  color: var(--qx-brand);
}

.sub {
  display: block;
  margin-top: 12rpx;
  color: var(--qx-text-secondary);
  font-size: 26rpx;
}

.form {
  background: var(--qx-card);
  border-radius: 20rpx;
  padding: 40rpx 36rpx 48rpx;
}

.field {
  margin-bottom: 28rpx;
}

.label {
  display: block;
  margin-bottom: 12rpx;
  font-size: 26rpx;
  color: var(--qx-text-secondary);
}

.input {
  height: 88rpx;
  padding: 0 24rpx;
  background: #f7f7f7;
  border-radius: 12rpx;
  font-size: 30rpx;
}

.btn {
  margin-top: 20rpx;
}

.btn.disabled {
  opacity: 0.6;
}

.switch {
  margin-top: 28rpx;
  text-align: center;
  color: var(--qx-brand);
  font-size: 26rpx;
  height: 44px;
  line-height: 44px;
}

.hint {
  display: block;
  margin-top: 12rpx;
  text-align: center;
  color: #bbb;
  font-size: 22rpx;
}
</style>
