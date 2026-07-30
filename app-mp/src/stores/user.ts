import { defineStore } from "pinia";
import {
  loginByAccount,
  registerAccount,
  fetchMe,
  type AppUser,
  type TokenPair,
} from "@/api/auth";
import {
  clearToken,
  clearUser,
  getToken,
  getUserJSON,
  setTokenPair,
  setUserJSON,
} from "@/utils/storage";

function applyToken(pair: TokenPair) {
  setTokenPair(pair.access_token, pair.refresh_token);
}

export const useUserStore = defineStore("user", {
  state: () => ({
    token: getToken() as string,
    profile: getUserJSON<AppUser>(),
  }),
  getters: {
    isLogin: (s) => !!s.token,
    displayName: (s) => s.profile?.nickname || s.profile?.account || "未登录",
  },
  actions: {
    async login(account: string, password: string) {
      const data = await loginByAccount(account, password);
      this.acceptSession(data.token, data.user);
    },
    async register(account: string, password: string, nickname?: string) {
      const data = await registerAccount(account, password, nickname);
      this.acceptSession(data.token, data.user);
    },
    acceptSession(pair: TokenPair, user: AppUser) {
      applyToken(pair);
      this.token = pair.access_token;
      this.profile = user;
      setUserJSON(user);
    },
    async refreshMe() {
      if (!this.token) return;
      try {
        const me = await fetchMe();
        this.profile = me;
        setUserJSON(me);
      } catch {
        // token 失效时由 request 清 token
      }
    },
    logout() {
      this.token = "";
      this.profile = null;
      clearToken();
      clearUser();
    },
  },
});
