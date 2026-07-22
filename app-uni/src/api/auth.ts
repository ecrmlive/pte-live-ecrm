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

export function loginByAccount(account: string, password: string) {
  return http.post<LoginResult>("/auth/login", { account, password }, false);
}

export function registerAccount(account: string, password: string, nickname?: string) {
  return http.post<LoginResult>(
    "/auth/register",
    { account, password, nickname: nickname || account },
    false
  );
}

export function fetchMe() {
  return http.get<AppUser>("/auth/me");
}
