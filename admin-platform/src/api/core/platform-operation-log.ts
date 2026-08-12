import { requestClient } from '#/api/request';

export interface PlatformOperationLog {
  admin_name: string;
  admin_user_id: number;
  created_at: string;
  id: number;
  ip: string;
  link: string;
  permission_name: string;
  request: string;
  request_method: string;
}

export interface PlatformOperationLogFilter {
  admin_keyword?: string;
  end_date?: string;
  limit: number;
  page: number;
  request_method?: string;
  start_date?: string;
}

export function listPlatformOperationLogs(params: PlatformOperationLogFilter) {
  return requestClient.get<{ list: PlatformOperationLog[]; total: number }>('/operation-logs', { params });
}

export interface PlatformLoginLog { admin_user_id?:number; created_at:string; id:number; ip:string; role_code:string; success:boolean; user_agent:string; username:string; }
export function listPlatformLoginLogs(params:{end_date?:string;limit:number;page:number;start_date?:string;success?:0|1;username?:string;}) { return requestClient.get<{list:PlatformLoginLog[];total:number}>('/login-logs',{params}); }
