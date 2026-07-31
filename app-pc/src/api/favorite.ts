import { http } from "@/utils/request";
import type { ProductItem } from "@/api/catalog";

export interface FavoriteStore {
  store_id: number;
  mer_id: number;
  store_name: string;
  merchant_app_id: string;
  follower_count: number;
}

export interface FavoriteState {
  id: number;
  followed: boolean;
}

export function fetchProductFavorites() {
  return http.get<{ type: "product"; list: ProductItem[] }>("/favorites?type=product");
}
export function fetchStoreFavorites() {
  return http.get<{ type: "store"; list: FavoriteStore[] }>("/favorites?type=store");
}

export function fetchProductFavoriteState(productID: number) {
  return http.get<FavoriteState>(`/favorites/products/${productID}`);
}
export function saveProductFavorite(productID: number) {
  return http.put<FavoriteState>(`/favorites/products/${productID}`);
}
export function removeProductFavorite(productID: number) {
  return http.delete<FavoriteState>(`/favorites/products/${productID}`);
}
export function fetchStoreFavoriteState(storeID: number) {
  return http.get<FavoriteState>(`/favorites/stores/${storeID}`);
}
export function saveStoreFavorite(storeID: number) {
  return http.put<FavoriteState>(`/favorites/stores/${storeID}`);
}
export function removeStoreFavorite(storeID: number) {
  return http.delete<FavoriteState>(`/favorites/stores/${storeID}`);
}
