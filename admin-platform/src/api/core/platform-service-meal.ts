import { requestClient } from '#/api/request';

export interface ServiceMealRow {
  meal_id: number;
  name: string;
  type: 1 | 2;
  price: number;
  num: number;
  sort: number;
  status: 0 | 1;
  create_time: string;
}

export interface ServiceMealInput {
  name: string;
  type: 1 | 2;
  price: number;
  num: number;
  sort: number;
  status: 0 | 1;
}

export function fetchServiceMeals(params: { page: number; limit: number }) {
  return requestClient.get<{
    limit: number;
    list: ServiceMealRow[];
    page: number;
    total: number;
  }>('/service-meals', { params });
}

export function createServiceMeal(data: ServiceMealInput) {
  return requestClient.post<{ ok: boolean }>('/service-meals', data);
}

export function updateServiceMeal(id: number, data: ServiceMealInput) {
  return requestClient.put<{ ok: boolean }>(`/service-meals/${id}`, data);
}

export function updateServiceMealStatus(id: number, status: 0 | 1) {
  return requestClient.put<{ ok: boolean }>(`/service-meals/${id}/status`, { status });
}

export function deleteServiceMeal(id: number) {
  return requestClient.delete<{ ok: boolean }>(`/service-meals/${id}`);
}
