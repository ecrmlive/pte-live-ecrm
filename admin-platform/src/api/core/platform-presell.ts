import { requestClient } from '#/api/request';
export interface PlatformPresell { end_time:string; mer_id:number; mer_name?:string; price:number; product_presell_id:number; start_time:string; status:number; store_name:string; }
export interface PlatformPresellPage { limit:number; list:PlatformPresell[]; page:number; total:number; }
export function listPlatformPresellsApi(params:{limit:number;mer_id?:number;page:number}) { return requestClient.get<PlatformPresellPage>('/presell/actives',{params}); }
