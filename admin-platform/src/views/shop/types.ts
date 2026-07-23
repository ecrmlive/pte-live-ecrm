export interface ShopRow {
  app_id: number;
  app_name: string;
  user_name: string;
  real_name?: string;
  platform_operator_name?: string;
  is_recycle: boolean;
  weixin_service: boolean;
  expire_time?: number;
  expire_time_text?: string;
  no_expire?: boolean;
  remain_gb?: null | number | string;
  create_time?: number | string;
  password?: string;
  password_confirm?: string;
}
