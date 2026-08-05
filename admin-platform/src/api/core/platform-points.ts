import { requestClient } from '#/api/request';

export interface PlatformPointsProduct { merchant_id:number; merchant_name:string; original_price:number; points_required:number; product_id:number; sale_status:number; stock:number; store_id:number; store_name:string; title:string; version:number; }
export interface PlatformPointsOrder { created_at:string; id:number; order_no:string; pay_status:string; points_amount:number; total_quantity:number; user_id:number; }
export function listPlatformPointsProductsApi(params:{keyword?:string;limit:number;merchant_id?:number;page:number;sale_status?:number}) { return requestClient.get<{list:PlatformPointsProduct[];total:number;page:number;limit:number}>('/points/products',{params}); }
export function getPlatformPointsSummaryApi() { return requestClient.get<{total:number;on_sale:number;stock:number}>('/points/products/summary'); }
export function updatePlatformPointsProductApi(id:number,payload:{points_required?:number;sale_status?:number;stock?:number;version:number}) { return requestClient.put<PlatformPointsProduct>(`/points/products/${id}`,payload); }
export function listPlatformPointsOrdersApi(params:{limit:number;page:number;pay_status?:string}) { return requestClient.get<{list:PlatformPointsOrder[];total:number;page:number;limit:number}>('/points/orders',{params}); }
