import { defineStore } from 'pinia';
import { ref } from 'vue';
import {
  changePassword as apiChangePassword,
  fetchMe,
  fetchMenus,
  fetchPermissions,
  login as apiLogin,
  type MenuNode,
  type PlatformUser,
} from '@/api/auth';
import { tokenStore } from '@/api/http';

export const useAuthStore = defineStore('auth', () => {
  const user = ref<PlatformUser | null>(null);
  const menus = ref<MenuNode[]>([]);
  const permissions = ref<string[]>([]);
  const booted = ref(false);

  async function login(account: string, password: string) {
    const { data } = await apiLogin(account, password);
    tokenStore.set(data.token.access_token, data.token.refresh_token);
    user.value = data.user;
    await Promise.all([loadMenus(), loadPermissions()]);
  }

  async function bootstrap() {
    if (!tokenStore.getAccess()) {
      booted.value = true;
      return;
    }
    try {
      const { data } = await fetchMe();
      user.value = data;
      await Promise.all([loadMenus(), loadPermissions()]);
    } catch {
      tokenStore.clear();
      user.value = null;
      menus.value = [];
      permissions.value = [];
    } finally {
      booted.value = true;
    }
  }

  async function loadMenus() {
    const { data } = await fetchMenus();
    menus.value = data.menus || [];
  }

  async function loadPermissions() {
    const { data } = await fetchPermissions();
    permissions.value = data.permissions || [];
  }

  function hasPerm(code: string) {
    if (user.value?.level === 0) return true;
    return permissions.value.includes(code);
  }

  async function changePassword(oldPassword: string, newPassword: string) {
    await apiChangePassword(oldPassword, newPassword);
  }

  function logout() {
    tokenStore.clear();
    user.value = null;
    menus.value = [];
    permissions.value = [];
  }

  return {
    user,
    menus,
    permissions,
    booted,
    login,
    bootstrap,
    loadMenus,
    loadPermissions,
    hasPerm,
    changePassword,
    logout,
  };
});
