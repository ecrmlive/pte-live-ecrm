import { callbackPost, http } from "@/utils/request";

export interface CheckResult {
  merchants: Array<{
    mer_id: number;
    mer_name: string;
    postage: number;
    total_price: number;
    total_num: number;
    items: Array<{
      cart_id: number;
      product_id: number;
      store_name: string;
      image: string;
      price: number;
      cart_num: number;
      subtotal: number;
    }>;
  }>;
  total_price: number;
  total_postage: number;
  coupon_price: number;
  integral: number;
  integral_price: number;
  user_integral: number;
  pay_price: number;
  total_num: number;
  give_integral: number;
  used_svip?: boolean;
  svip_discount?: number;
}

export interface PointsCheckResult {
  product_id: number;
  store_name: string;
  cart_num: number;
  integral: number;
  user_integral: number;
  pay_price: number;
  mer_id: number;
  mer_name?: string;
  activity_type: number;
}

export interface GroupOrder {
  group_order_id: number;
  group_order_sn: string;
  uid?: number;
  pay_price: number;
  total_price: number;
  total_num: number;
  integral?: number;
  paid: number;
  pay_type: number;
  activity_type?: number;
  real_name: string;
  user_phone: string;
  user_address: string;
  create_time: string;
  orders?: Array<{
    order_id: number;
    order_sn: string;
    mer_id: number;
    mer_name?: string;
    pay_price: number;
    paid: number;
    status: number;
    delivery_name?: string;
    delivery_id?: string;
    products?: Array<{
      product_id: number;
      product_num: number;
      product_price: number;
      total_price: number;
      product_info?: string;
    }>;
  }>;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

/** FUNCTIONAL-TRUTH：POST …/v2/order/check|create */
export function v2Check(
  cart_ids: number[],
  coupon_user_ids: number[] = [],
  use_integral = 0,
) {
  return http.post<CheckResult>("/v2/order/check", {
    cart_ids,
    coupon_user_ids,
    use_integral,
  });
}

export function v2Create(data: {
  cart_ids: number[];
  address_id: number;
  mark?: string;
  coupon_user_ids?: number[];
  use_integral?: number;
}) {
  return http.post<GroupOrder>("/v2/order/create", data);
}

export type PayType = "mock" | "balance" | "wechat" | "alipay";

export interface PayIntent {
  status: "pending" | "paid";
  channel: string;
  group_order_id: number;
  out_trade_no: string;
  pay_price: number;
  sandbox?: boolean;
  notify_token?: string;
}

export function payGroup(id: number, pay_type: PayType = "mock") {
  return http.post<GroupOrder | PayIntent>(`/order/pay/${id}`, { pay_type });
}

export function notifyChannelPay(
  channel: "wechat" | "alipay",
  body: {
    group_order_id: number;
    uid: number;
    out_trade_no: string;
    pay_price: number;
    notify_token: string;
  },
) {
  return callbackPost<GroupOrder>(`/pay/${channel}`, body);
}

/** 积分商城：/order/v3/*（禁止 /v3/order/*） */
export function v3Check(data: {
  product_id: number;
  cart_num?: number;
  product_attr_unique?: string;
  address_id?: number;
}) {
  return http.post<PointsCheckResult>("/order/v3/check", data);
}

export function v3Create(data: {
  product_id: number;
  address_id: number;
  cart_num?: number;
  product_attr_unique?: string;
}) {
  return http.post<GroupOrder>("/order/v3/create", data);
}

export function pointsPay(id: number) {
  return http.post<GroupOrder>(`/order/points/pay/${id}`, {});
}

export function fetchIntegral() {
  return http.get<{ integral: number }>("/integral");
}

export function fetchOrders(page = 1, limit = 20) {
  return http.get<PageResult<GroupOrder>>(`/orders?page=${page}&limit=${limit}`);
}

export function fetchOrderDetail(id: number) {
  return http.get<GroupOrder>(`/orders/${id}`);
}
