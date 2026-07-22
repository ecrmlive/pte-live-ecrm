import { callbackPost, http } from "@/utils/request";

export interface CartItem {
  cart_id: number;
  mer_id: number;
  product_id: number;
  product_attr_unique: string;
  cart_num: number;
  store_name?: string;
  mer_name?: string;
  image?: string;
  price?: number;
  stock?: number;
  is_fail?: number;
}

export interface CartBucket {
  mer_id: number;
  mer_name: string;
  subtotal: number;
  items: CartItem[];
}

export interface Address {
  address_id: number;
  real_name: string;
  phone: string;
  province: string;
  city: string;
  district: string;
  detail: string;
  post_code: number;
  is_default: number;
}

export interface GroupOrder {
  group_order_id: number;
  group_order_sn: string;
  uid?: number;
  pay_price: number;
  total_num: number;
  paid: number;
  pay_type: number;
  orders?: StoreOrder[];
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

export interface StoreOrder {
  order_id: number;
  order_sn: string;
  mer_id: number;
  mer_name?: string;
  pay_price: number;
  paid: number;
  status: number;
  delivery_name?: string;
  delivery_id?: string;
}

export function fetchCart() {
  return http.get<{ list: CartBucket[]; total_num?: number; total_price?: number }>("/cart");
}

export function addCart(data: {
  product_id: number;
  cart_num: number;
  product_attr_unique?: string;
}) {
  return http.post<CartItem>("/cart", data);
}

export function updateCartNum(cartId: number, cart_num: number) {
  return http.put(`/cart/${cartId}`, { cart_num });
}

export function removeCart(cartId: number) {
  return http.delete(`/cart/${cartId}`);
}

export function fetchAddresses() {
  return http.get<{ list: Address[] }>("/address");
}

export function createAddress(data: Partial<Address> & { real_name: string; phone: string; detail: string }) {
  return http.post<Address>("/address", data as Record<string, unknown>);
}

export function orderCheck(cart_ids: number[], address_id: number, coupon_user_ids: number[] = []) {
  return http.post<{
    pay_price: number;
    total_num: number;
    coupon_price?: number;
    merchants: unknown[];
  }>("/v2/order/check", {
    cart_ids,
    address_id,
    coupon_user_ids,
  });
}

export function orderCreate(
  cart_ids: number[],
  address_id: number,
  mark = "",
  coupon_user_ids: number[] = [],
) {
  return http.post<GroupOrder>("/v2/order/create", {
    cart_ids,
    address_id,
    mark,
    coupon_user_ids,
  });
}

export function payOrder(groupOrderId: number, pay_type: PayType = "mock") {
  return http.post<GroupOrder | PayIntent>(`/order/pay/${groupOrderId}`, { pay_type });
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

export function fetchOrders(page = 1) {
  return http.get<{ list: GroupOrder[]; total: number }>(`/orders?page=${page}`);
}

export function fetchOrderDetail(id: number) {
  return http.get<GroupOrder>(`/orders/${id}`);
}
