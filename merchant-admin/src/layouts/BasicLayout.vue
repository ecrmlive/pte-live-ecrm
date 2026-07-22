<template>
  <a-layout class="shell">
    <a-layout-sider v-model:collapsed="collapsed" collapsible theme="light" width="232" class="sider">
      <div class="brand">
        <span class="mark">店</span>
        <div v-if="!collapsed" class="brand-text">
          <strong>{{ auth.user?.mer_name || '商户后台' }}</strong>
          <small>mer_id={{ auth.user?.mer_id ?? '-' }}</small>
        </div>
      </div>
      <a-menu
        v-model:selectedKeys="selectedKeys"
        v-model:openKeys="openKeys"
        mode="inline"
        :items="menuItems"
        @click="onMenuClick"
      />
    </a-layout-sider>
    <a-layout>
      <a-layout-header class="header">
        <div class="header-title">{{ pageTitle }}</div>
        <div class="header-actions">
          <span class="user">{{ auth.user?.real_name || auth.user?.account }}</span>
          <a-button type="link" @click="showPwd = true">改密</a-button>
          <a-button type="link" danger @click="onLogout">退出</a-button>
        </div>
      </a-layout-header>
      <a-layout-content class="content">
        <router-view />
      </a-layout-content>
    </a-layout>

    <a-modal v-model:open="showPwd" title="修改密码" :confirm-loading="pwdLoading" @ok="submitPwd">
      <a-form layout="vertical">
        <a-form-item label="原密码">
          <a-input-password v-model:value="pwdForm.oldPassword" />
        </a-form-item>
        <a-form-item label="新密码">
          <a-input-password v-model:value="pwdForm.newPassword" />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-layout>
</template>

<script setup lang="ts">
import { computed, h, reactive, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { message, type MenuProps } from 'ant-design-vue';
import {
  AccountBookOutlined,
  DashboardOutlined,
  GiftOutlined,
  ProfileOutlined,
  SettingOutlined,
  ShoppingOutlined,
} from '@ant-design/icons-vue';
import { useAuthStore } from '@/stores/auth';
import type { MenuNode } from '@/api/auth';

const auth = useAuthStore();
const route = useRoute();
const router = useRouter();
const collapsed = ref(false);
const selectedKeys = ref<string[]>([]);
const openKeys = ref<string[]>([]);
const showPwd = ref(false);
const pwdLoading = ref(false);
const pwdForm = reactive({ oldPassword: '', newPassword: '' });

const iconMap: Record<string, unknown> = {
  DashboardOutlined,
  ShoppingOutlined,
  ProfileOutlined,
  SettingOutlined,
  AccountBookOutlined,
  GiftOutlined,
};

const pageTitle = computed(() => (route.meta.title as string) || '商户后台');

function toItems(nodes: MenuNode[]): MenuProps['items'] {
  return nodes.map((n) => {
    const IconComp = n.icon ? iconMap[n.icon] : undefined;
    const item: NonNullable<MenuProps['items']>[number] = {
      key: n.path,
      label: n.menu_name,
      icon: IconComp ? () => h(IconComp as object) : undefined,
    };
    if (n.children?.length) {
      (item as { children?: MenuProps['items'] }).children = toItems(n.children);
    }
    return item;
  });
}

const menuItems = computed(() => toItems(auth.menus));

function findMenuPath(nodes: MenuNode[], target: string, trail: string[] = []): string[] | null {
  for (const n of nodes) {
    const next = [...trail, n.path];
    if (n.path === target) return next;
    if (n.children?.length) {
      const hit = findMenuPath(n.children, target, next);
      if (hit) return hit;
    }
  }
  return null;
}

function findMenuNode(nodes: MenuNode[], target: string): MenuNode | null {
  for (const n of nodes) {
    if (n.path === target) return n;
    if (n.children?.length) {
      const hit = findMenuNode(n.children, target);
      if (hit) return hit;
    }
  }
  return null;
}

watch(
  () => [route.path, auth.menus] as const,
  ([path]) => {
    selectedKeys.value = [path];
    const chain = findMenuPath(auth.menus, path);
    if (chain && chain.length > 1) {
      openKeys.value = chain.slice(0, -1);
      return;
    }
    const parts = path.split('/').filter(Boolean);
    if (parts.length > 1) {
      openKeys.value = [`/${parts[0]}`];
    }
  },
  { immediate: true },
);

function onMenuClick(info: { key: string | number }) {
  const path = String(info.key);
  const node = findMenuNode(auth.menus, path);
  if (node?.children?.length) {
    return;
  }
  if (path && path !== route.path) {
    router.push(path);
  }
}

function onLogout() {
  auth.logout();
  router.replace('/login');
}

async function submitPwd() {
  if (!pwdForm.oldPassword || !pwdForm.newPassword) {
    message.warning('请填写完整');
    return;
  }
  pwdLoading.value = true;
  try {
    await auth.changePassword(pwdForm.oldPassword, pwdForm.newPassword);
    message.success('密码已更新');
    showPwd.value = false;
    pwdForm.oldPassword = '';
    pwdForm.newPassword = '';
  } finally {
    pwdLoading.value = false;
  }
}
</script>

<style scoped>
.shell {
  min-height: 100vh;
  background: var(--qx-bg);
}
.sider {
  border-right: 1px solid var(--qx-line);
  background: var(--qx-panel);
}
.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 64px;
  padding: 0 16px;
  border-bottom: 1px solid var(--qx-line);
}
.mark {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: grid;
  place-items: center;
  background: linear-gradient(145deg, #1f6f8b, #124556);
  color: #fff;
  font-weight: 700;
}
.brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}
.brand-text strong {
  font-size: 15px;
}
.brand-text small {
  color: #6b7785;
  margin-top: 2px;
}
.header {
  background: var(--qx-panel);
  border-bottom: 1px solid var(--qx-line);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 64px;
}
.header-title {
  font-size: 16px;
  font-weight: 600;
}
.header-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}
.user {
  margin-right: 8px;
  color: #516070;
}
.content {
  margin: 16px;
  min-height: calc(100vh - 96px);
}
</style>
