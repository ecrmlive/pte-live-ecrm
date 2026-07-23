import { publicRequestClient } from '#/api/request';

export interface LoginCodeData {
  codeImage?: string;
  codeKey?: string;
}

export interface AdminLoginBaseResult {
  codeData?: LoginCodeData;
  settings?: {
    admin_bg_img?: string;
    admin_name?: string;
  };
}

export async function getAdminLoginBaseApi() {
  return publicRequestClient.post<AdminLoginBaseResult>('/admin/index/base', {});
}
