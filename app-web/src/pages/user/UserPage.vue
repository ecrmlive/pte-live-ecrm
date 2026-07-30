<script setup lang="ts">
import { onMounted } from "vue";
import { useRouter } from "vue-router";
import { useUserStore } from "@/stores/user";

const user = useUserStore();
const router = useRouter();

onMounted(() => {
  if (user.isLogin) user.refreshMe();
});

function goLogin() {
  router.push({ name: "login", query: { redirect: "/user" } });
}
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <header>
        <h1>个人中心</h1>
        <p>资产、订单、优惠券、地址、发票与设置入口（功能表 4 · 用户中心）。</p>
      </header>

      <div v-if="user.isLogin" class="profile">
        <div class="avatar">{{ user.displayName.slice(0, 1) }}</div>
        <div>
          <h2>{{ user.displayName }}</h2>
          <p class="muted">账号 {{ user.profile?.account || "-" }} · UID {{ user.profile?.uid || "-" }}</p>
        </div>
      </div>
      <div v-else class="guest">
        <p>登录后可查看订单、优惠券与地址。</p>
        <button class="pc-btn" type="button" @click="goLogin">去登录</button>
      </div>

      <div class="tiles">
        <article class="clickable" @click="router.push('/orders')">
          <h3>我的订单</h3>
          <p>待付款 / 已支付 / 多商户子单</p>
          <span class="tag">进入</span>
        </article>
        <article class="clickable" @click="router.push('/presell')">
          <h3>预售 / 尾款</h3>
          <p>全款预售 · 定金待付尾款</p>
          <span class="tag">进入</span>
        </article>
        <article class="clickable" @click="router.push('/assist')">
          <h3>好友助力</h3>
          <p>发起助力 · 满员下单</p>
          <span class="tag">进入</span>
        </article>
        <article class="clickable" @click="router.push('/community')">
          <h3>社区种草</h3>
          <p>浏览种草 · 发帖审核</p>
          <span class="tag">进入</span>
        </article>
        <article class="clickable" @click="router.push('/coupons')">
          <h3>优惠券</h3>
          <p>领券中心 / 我的未用券</p>
          <span class="tag">进入</span>
        </article>
        <article class="clickable" @click="router.push('/points')">
          <h3>积分商城</h3>
          <p>积分兑换 · 独立 v3 入口</p>
          <span class="tag">进入</span>
        </article>
        <article class="clickable" @click="router.push('/notices')">
          <h3>平台公告</h3>
          <p>运营通知与说明</p>
          <span class="tag">进入</span>
        </article>
        <article class="clickable" @click="router.push('/agreements/sys_user_agree')">
          <h3>用户协议</h3>
          <p>注册与使用条款</p>
          <span class="tag">进入</span>
        </article>
        <article class="clickable" @click="router.push('/agreements/sys_userr_privacy')">
          <h3>隐私政策</h3>
          <p>个人信息处理说明</p>
          <span class="tag">进入</span>
        </article>
        <article>
          <h3>收货地址</h3>
          <p>地址增删改查</p>
          <span class="tag">阶段 3</span>
        </article>
        <article>
          <h3>账户安全</h3>
          <p>改密 / 改手机 / 发票</p>
          <span class="tag">阶段 4+</span>
        </article>
      </div>
    </section>
  </div>
</template>

<style scoped>
.panel {
  background: var(--pc-surface);
  border: 1px solid var(--pc-line);
  border-radius: calc(var(--pc-radius) + 4px);
  padding: 1.6rem 1.8rem 1.8rem;
  box-shadow: var(--pc-shadow);
}

header h1 {
  margin: 0;
  font-size: 1.5rem;
}

header p {
  margin: 0.45rem 0 0;
  color: var(--pc-muted);
}

.profile,
.guest {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin: 1.4rem 0 1.6rem;
  padding: 1rem 1.1rem;
  border-radius: var(--pc-radius);
  background: var(--pc-brand-soft);
}

.avatar {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  background: var(--pc-brand);
  color: #fff;
  font-size: 1.3rem;
  font-weight: 700;
}

.profile h2 {
  margin: 0;
  font-size: 1.15rem;
}

.muted {
  margin: 0.25rem 0 0;
  color: var(--pc-muted);
}

.guest p {
  margin: 0;
  flex: 1;
}

.tiles {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.85rem;
}

article {
  border: 1px solid var(--pc-line);
  border-radius: 8px;
  padding: 1rem;
}

article h3 {
  margin: 0 0 0.4rem;
  font-size: 1rem;
}

article p {
  margin: 0;
  color: var(--pc-muted);
  font-size: 0.9rem;
  min-height: 2.6em;
}

.tag {
  display: inline-block;
  margin-top: 0.7rem;
  font-size: 0.75rem;
  color: var(--pc-brand);
  background: var(--pc-brand-soft);
  padding: 0.15rem 0.45rem;
  border-radius: 999px;
}

@media (max-width: 900px) {
  .tiles {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
