/** Compatibility layer for pages still importing `@/api/trade`. Prefer `@/api/cart` + `@/api/order`. */
export {
  fetchCart,
  addCart,
  setCartNum as updateCartNum,
  removeCart,
  fetchAddresses,
  type CartBucket,
  type Address,
} from "./cart";

export {
  fetchOrders,
  fetchOrderDetail,
  payGroup as payOrder,
  type GroupOrder,
} from "./order";

import { v2Create } from "./order";

export function orderCreate(cart_ids: number[], address_id: number, mark = "") {
  return v2Create({ cart_ids, address_id, mark });
}
