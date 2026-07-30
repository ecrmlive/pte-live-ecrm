import { http } from "@/utils/request";
import {
  clearMerchantAppID,
  getMerchantAppID,
  getRefreshToken,
  getToken,
  setMerchantAppID,
  setTokenPair,
} from "@/utils/storage";

export interface AppUser {
  uid: number;
  account: string;
  nickname: string;
  avatar?: string;
  phone?: string;
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
  return http.post<LoginResult>("/auth/login", { account, password, channel: "pc" }, false);
}

export function registerAccount(account: string, password: string, nickname?: string) {
  return http.post<LoginResult>(
    "/auth/register",
    { account, password, nickname: nickname || account, channel: "pc" },
    false,
  );
}

export function fetchMe() {
  return http.get<AppUser>("/auth/me");
}

/**
 * 由 api-business 解析 X-AppId 读模型并签发店铺上下文 JWT。
 * 客户端只保存店铺 AppId；商户 ID 与 IM SDK AppId 只来自服务端令牌。
 */
export async function activateMerchantContext(appID: string) {
  const previousAppID = getMerchantAppID();
  const previousAccessToken = getToken();
  const previousRefreshToken = getRefreshToken();
  setMerchantAppID(appID);
  try {
    const data = await http.post<{ token: TokenPair }>("/auth/store-context", {});
    setTokenPair(data.token.access_token, data.token.refresh_token);
    return data;
  } catch (error) {
    if (previousAppID) setMerchantAppID(previousAppID);
    else clearMerchantAppID();
    if (previousAccessToken) setTokenPair(previousAccessToken, previousRefreshToken);
    throw error;
  }
}
