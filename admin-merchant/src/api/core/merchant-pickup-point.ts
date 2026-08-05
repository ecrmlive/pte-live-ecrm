import { requestClient } from '#/api/request';

export interface StorePickupPoint {
  contact_name: string;
  detail: string;
  id: number;
  is_default: number;
  mobile: string;
  region_code: string;
}

export function listStorePickupPointsApi() {
  return requestClient.get<{ list: StorePickupPoint[] }>('/store-pickup-points');
}

export function createStorePickupPointApi(body: Omit<StorePickupPoint, 'id'>) {
  return requestClient.post<StorePickupPoint>('/store-pickup-points', body);
}

export function updateStorePickupPointApi(id: number, body: Omit<StorePickupPoint, 'id'>) {
  return requestClient.put<StorePickupPoint>(`/store-pickup-points/${id}`, body);
}

export function deleteStorePickupPointApi(id: number) {
  return requestClient.delete<{ ok: boolean }>(`/store-pickup-points/${id}`);
}
