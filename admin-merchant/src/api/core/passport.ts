import { requestClient } from '#/api/request';

export interface LoginCodeData {
  codeImage?: string;
  codeKey?: string;
}

export interface ShopLoginBaseResult {
  codeData?: LoginCodeData;
  settings?: {
    shop_bg_img?: string;
    shop_logo_img?: string;
    shop_name?: string;
  };
}

export async function getShopLoginBaseApi() {
  return requestClient.post<ShopLoginBaseResult>('/shop/index/base', {});
}

export interface ShopEditPasswordForm {
  confirmPass: string;
  oldpass: string;
  password: string;
}

export async function editShopPasswordApi(data: ShopEditPasswordForm) {
  return requestClient.post<null>('/shop/passport/editPass', data);
}
