import { liveAdminPost } from '#/api/live-request';

export interface JwtSessionRow {
  session_id: number;
  uid: number;
  nick_name: string;
  login_platform: string;
  login_platform_label: string;
  shop_name: string;
  user_role: string;
  jwt_token: string;
  expires_at: number;
  last_login_time: number;
  app_id: number;
  token_type: string;
}

export interface JwtSessionListResult {
  list: JwtSessionRow[];
  total: number;
  page: number;
  list_rows: number;
}

export async function listJwtSessionsApi(params: Record<string, unknown>) {
  return liveAdminPost<JwtSessionListResult>(
    '/api/v1/admin/auth/session/list',
    params,
  );
}

export async function deleteJwtSessionApi(sessionId: number) {
  return liveAdminPost('/api/v1/admin/auth/session/delete', {
    session_id: sessionId,
  });
}

export async function batchDeleteJwtSessionsApi(sessionIds: number[]) {
  return liveAdminPost('/api/v1/admin/auth/session/batch-delete', {
    session_ids: sessionIds.join(','),
  });
}
