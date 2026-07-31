import { defineStore } from "pinia";
import {
  fetchMe,
  loginByAccount,
  registerAccount,
  type AppUser,
  type TokenPair,
} from "@/api/auth";
import {
  clearToken,
  clearMerchantAppID,
  clearUser,
  getToken,
  getUserJSON,
  setTokenPair,
  setUserJSON,
} from "@/utils/storage";
import { fetchCart } from "@/api/trade";

export const useUserStore = defineStore("user", {
  state: () => ({
    token: getToken() as string,
    profile: getUserJSON<AppUser>(),
    cartCount: 0,
  }),
  getters: {
    isLogin: (s) => !!s.token,
    displayName: (s) => s.profile?.nickname || s.profile?.account || "未登录",
  },
  actions: {
    acceptSession(pair: TokenPair, user: AppUser) {
      setTokenPair(pair.access_token, pair.refresh_token);
      this.token = pair.access_token;
      this.profile = user;
      setUserJSON(user);
    },
    replaceTokenPair(pair: TokenPair) {
      setTokenPair(pair.access_token, pair.refresh_token);
      this.token = pair.access_token;
    },
    async login(account: string, password: string, captchaToken: string) {
      const data = await loginByAccount(account, password, captchaToken);
      this.acceptSession(data.token, data.user);
    },
    async register(account: string, password: string, nickname?: string) {
      const data = await registerAccount(account, password, nickname);
      this.acceptSession(data.token, data.user);
    },
    async refreshMe() {
      if (!this.token) return;
      try {
        const me = await fetchMe();
        this.profile = me;
        setUserJSON(me);
      } catch {
        // token 失效由 request 清 token
      }
    },
    async refreshCartCount() {
      if (!this.token) {
        this.cartCount = 0;
        return;
      }
      try {
        const cart = await fetchCart();
        this.cartCount = Number(cart.total_num || 0);
      } catch {
        this.cartCount = 0;
      }
    },
    logout() {
      this.token = "";
      this.profile = null;
      this.cartCount = 0;
      clearToken();
      clearMerchantAppID();
      clearUser();
    },
  },
});
