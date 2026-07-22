import { defineStore } from 'pinia';
import { ref } from 'vue';
import { fetchMe, login as apiLogin, type ManagerUser } from '@/api/auth';
import { tokenStore } from '@/api/http';

export const useAuthStore = defineStore('auth', () => {
  const user = ref<ManagerUser | null>(null);
  const booted = ref(false);

  async function login(account: string, password: string) {
    const { data } = await apiLogin(account, password);
    tokenStore.set(data.token.access_token, data.token.refresh_token);
    user.value = data.user;
  }

  async function bootstrap() {
    if (!tokenStore.getAccess()) {
      booted.value = true;
      return;
    }
    try {
      const { data } = await fetchMe();
      user.value = data;
    } catch {
      tokenStore.clear();
      user.value = null;
    } finally {
      booted.value = true;
    }
  }

  function logout() {
    tokenStore.clear();
    user.value = null;
  }

  return { user, booted, login, bootstrap, logout };
});
