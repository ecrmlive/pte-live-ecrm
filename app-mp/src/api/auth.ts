import { http } from "@/utils/request";

export interface AppUser {
  uid: number;
  account: string;
  nickname: string;
  avatar?: string;
  phone?: string;
  integral?: number;
  now_money?: number;
  is_svip?: number;
  svip_endtime?: string | null;
  is_svip_active?: boolean;
}

export interface TokenPair {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface LoginResult {
  token: TokenPair;
  user: AppUser;
}

export type CUserChannel = "wechat" | "mini_program" | "h5" | "pc" | "ios" | "android" | "harmony";

export function currentClientChannel(): CUserChannel {
  // #ifdef MP-WEIXIN
  return "mini_program";
  // #endif
  return "h5";
}

export function loginByAccount(account: string, password: string, channel = currentClientChannel()) {
  return http.post<LoginResult>("/auth/login", { account, password, channel }, false);
}

export function registerAccount(account: string, password: string, nickname?: string, channel = currentClientChannel()) {
  return http.post<LoginResult>(
    "/auth/register",
    { account, password, nickname: nickname || account, channel },
    false
  );
}

export function fetchMe() {
  return http.get<AppUser>("/auth/me");
}
