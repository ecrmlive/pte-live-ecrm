import { http } from "@/utils/request";

export interface CartItem {
  cart_id: number;
  product_id: number;
  product_attr_unique: string;
  cart_num: number;
  mer_id: number;
  mer_name?: string;
  store_name?: string;
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

export interface CartListResult {
  list: CartBucket[];
  total_num: number;
  total_price: number;
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

export function fetchCart() {
  return http.get<CartListResult>("/cart");
}

export function addCart(data: {
  product_id: number;
  product_attr_unique?: string;
  cart_num?: number;
  is_new?: number;
}) {
  return http.post<CartItem>("/cart", data);
}

export function setCartNum(id: number, cart_num: number) {
  return http.put<{ ok: boolean }>(`/cart/${id}`, { cart_num });
}

export function deleteCart(id: number) {
  return http.delete<{ ok: boolean }>(`/cart/${id}`);
}

/** @deprecated alias */
export function removeCart(id: number) {
  return deleteCart(id);
}

export function fetchAddresses() {
  return http.get<{ list: Address[] }>("/address");
}

export function createAddress(data: Partial<Address> & { real_name: string; phone: string; detail: string }) {
  return http.post<Address>("/address", data as Record<string, unknown>);
}

export function updateAddress(id: number, data: Partial<Address>) {
  return http.put<Address>(`/address/${id}`, data as Record<string, unknown>);
}

export function deleteAddress(id: number) {
  return http.delete<{ ok: boolean }>(`/address/${id}`);
}
