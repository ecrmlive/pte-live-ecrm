import { requestClient } from '#/api/request';
export interface SpreadLog { create_time: string; old_spread_uid: number; spread_uid: number; uid: number; user_spread_log_id: number; }
export interface SpreadLogPage { limit: number; list: SpreadLog[]; page: number; total: number; }
export function listPlatformSpreadLogsApi(params: { limit: number; page: number }) { return requestClient.get<SpreadLogPage>('/spread/logs', { params }); }
