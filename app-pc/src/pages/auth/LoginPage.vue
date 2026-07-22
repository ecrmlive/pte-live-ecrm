<script setup lang="ts">
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useUserStore } from "@/stores/user";
import { ApiError } from "@/utils/request";

const user = useUserStore();
const router = useRouter();
const route = useRoute();

const mode = ref<"login" | "register">("login");
const account = ref("");
const password = ref("");
const nickname = ref("");
const loading = ref(false);
const message = ref("");

async function submit() {
  message.value = "";
  if (!account.value.trim() || !password.value) {
    message.value = "请输入账号和密码";
    return;
  }
  loading.value = true;
  try {
    if (mode.value === "login") {
      await user.login(account.value.trim(), password.value);
    } else {
      await user.register(account.value.trim(), password.value, nickname.value.trim() || undefined);
    }
    const redirect = typeof route.query.redirect === "string" ? route.query.redirect : "/user";
    await router.replace(redirect || "/user");
  } catch (e) {
    const err = e as ApiError;
    if (err.status === 404 || /ping|Not Found|404/i.test(err.message)) {
      message.value = "登录接口尚未就绪（api-app 阶段 1）。可先浏览首页骨架。";
    } else {
      message.value = err.message || "登录失败";
    }
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div class="pc-container">
    <div class="panel">
      <div class="copy">
        <h1>{{ mode === "login" ? "欢迎回来" : "注册账号" }}</h1>
        <p>手机号快捷登录 / 微信扫码将在后续阶段接入。当前走账号密码，对接 `/api/app/v1/auth/*`。</p>
      </div>
      <form class="form" @submit.prevent="submit">
        <label>
          <span>账号</span>
          <input v-model="account" autocomplete="username" placeholder="手机号或账号" />
        </label>
        <label v-if="mode === 'register'">
          <span>昵称（可选）</span>
          <input v-model="nickname" autocomplete="nickname" placeholder="展示名称" />
        </label>
        <label>
          <span>密码</span>
          <input
            v-model="password"
            type="password"
            autocomplete="current-password"
            placeholder="登录密码"
          />
        </label>
        <p v-if="message" class="msg">{{ message }}</p>
        <button class="pc-btn" type="submit" :disabled="loading">
          {{ loading ? "提交中…" : mode === "login" ? "登录" : "注册并登录" }}
        </button>
        <button
          class="switch"
          type="button"
          @click="mode = mode === 'login' ? 'register' : 'login'"
        >
          {{ mode === "login" ? "没有账号？去注册" : "已有账号？去登录" }}
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.panel {
  max-width: 880px;
  margin: 1rem auto 0;
  display: grid;
  grid-template-columns: 1.1fr 1fr;
  background: var(--pc-surface);
  border: 1px solid var(--pc-line);
  border-radius: calc(var(--pc-radius) + 4px);
  overflow: hidden;
  box-shadow: var(--pc-shadow);
}

.copy {
  padding: 2.4rem 2rem;
  background:
    linear-gradient(160deg, rgba(15, 107, 92, 0.92), rgba(20, 55, 70, 0.95)),
    #0f6b5c;
  color: #fff;
}

.copy h1 {
  margin: 0 0 0.8rem;
  font-size: 1.8rem;
}

.copy p {
  margin: 0;
  line-height: 1.7;
  opacity: 0.92;
}

.form {
  padding: 2rem 1.8rem;
  display: grid;
  gap: 0.9rem;
}

label {
  display: grid;
  gap: 0.35rem;
}

label span {
  color: var(--pc-muted);
  font-size: 0.9rem;
}

input {
  border: 1px solid var(--pc-line);
  border-radius: 8px;
  padding: 0.7rem 0.85rem;
  outline: none;
}

input:focus {
  border-color: rgba(15, 107, 92, 0.55);
  box-shadow: 0 0 0 3px rgba(15, 107, 92, 0.12);
}

.msg {
  margin: 0;
  color: var(--pc-danger);
  font-size: 0.92rem;
}

.switch {
  border: 0;
  background: transparent;
  color: var(--pc-brand);
  padding: 0.25rem 0;
  text-align: left;
}

@media (max-width: 780px) {
  .panel {
    grid-template-columns: 1fr;
  }
}
</style>
