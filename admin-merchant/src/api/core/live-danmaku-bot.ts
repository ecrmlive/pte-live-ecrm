import { requestClient } from '#/api/request';

export interface DanmakuScriptTemplateItem {
  content?: string;
  line_count?: number;
  name: string;
  script_template_id: number;
}

export interface DanmakuRobotTemplateItem {
  bot_count: number;
  loop_enabled?: number;
  name: string;
  robot_template_id: number;
  script_template_id: number;
  script_template_name?: string;
}

export interface DanmakuTemplateListResult<T> {
  list: T[];
  total: number;
}

export async function getDanmakuScriptTemplateListApi(params: {
  list_rows?: number;
  page?: number;
}) {
  return requestClient.post<DanmakuTemplateListResult<DanmakuScriptTemplateItem>>(
    '/api/v1/shop/live/danmaku-bot/script-template/list',
    params,
  );
}

export async function saveDanmakuScriptTemplateApi(payload: {
  content: string;
  name: string;
  script_template_id?: number;
}) {
  return requestClient.post<{ msg?: string; script_template_id?: number }>(
    '/api/v1/shop/live/danmaku-bot/script-template/save',
    payload,
  );
}

export async function deleteDanmakuScriptTemplateApi(payload: {
  script_template_id: number;
}) {
  return requestClient.post<{ msg?: string }>(
    '/api/v1/shop/live/danmaku-bot/script-template/delete',
    payload,
  );
}

export async function getDanmakuRobotTemplateListApi(params: {
  list_rows?: number;
  page?: number;
}) {
  return requestClient.post<DanmakuTemplateListResult<DanmakuRobotTemplateItem>>(
    '/api/v1/shop/live/danmaku-bot/robot-template/list',
    params,
  );
}

export async function saveDanmakuRobotTemplateApi(payload: {
  bot_count: number;
  loop_enabled?: number;
  name: string;
  robot_template_id?: number;
  script_template_id: number;
}) {
  return requestClient.post<{ msg?: string; robot_template_id?: number }>(
    '/api/v1/shop/live/danmaku-bot/robot-template/save',
    payload,
  );
}

export async function deleteDanmakuRobotTemplateApi(payload: {
  robot_template_id: number;
}) {
  return requestClient.post<{ msg?: string }>(
    '/api/v1/shop/live/danmaku-bot/robot-template/delete',
    payload,
  );
}
