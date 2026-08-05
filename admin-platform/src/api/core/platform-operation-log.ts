import { requestClient } from '#/api/request';

export interface PlatformOperationLog {
  id: number;
  admin_user_id: number;
  role_code: string;
  action: string;
  resource_type: string;
  resource_id: string;
  request_id: string;
  created_at: string;
}

export interface PlatformOperationLogFilter {
  page: number;
  limit: number;
  admin_user_id?: number;
  role_code?: string;
  action?: string;
  resource_type?: string;
  start_date?: string;
  end_date?: string;
}

export function listPlatformOperationLogs(params: PlatformOperationLogFilter) {
  return requestClient.get<{ list: PlatformOperationLog[]; total: number }>('/operation-logs', { params });
}

export interface PlatformLoginLog { id:number; admin_user_id?:number; username:string; role_code:string; success:boolean; ip:string; user_agent:string; created_at:string; }
export function listPlatformLoginLogs(params:{page:number;limit:number;username?:string;success?:0|1;start_date?:string;end_date?:string}) { return requestClient.get<{list:PlatformLoginLog[];total:number}>('/login-logs',{params}); }
