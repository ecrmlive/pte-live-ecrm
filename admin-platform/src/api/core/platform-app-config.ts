import { requestClient } from '#/api/request';

export interface AppStubConfig {
  name: string;
  enabled: boolean;
  remark: string;
}

export type AppConfigKey = 'routine' | 'wechat-reply' | 'wechat-menus' | 'wechat-template' | 'wechat-news';

const PATHS: Record<AppConfigKey, string> = {
  routine: '/setting/routine-app',
  'wechat-reply': '/setting/wechat-reply',
  'wechat-menus': '/setting/wechat-menus',
  'wechat-template': '/setting/wechat-template',
  'wechat-news': '/setting/wechat-news',
};

function parseConfig(raw: string): AppStubConfig {
  return JSON.parse(raw) as AppStubConfig;
}

function stringifyConfig(config: AppStubConfig): string {
  return JSON.stringify(config);
}

export function getAppStubConfigApi(key: AppConfigKey) {
  return requestClient
    .get<{ config: string; note: string }>(PATHS[key])
    .then((data) => ({ note: data.note, config: parseConfig(data.config) }));
}

export function saveAppStubConfigApi(key: AppConfigKey, config: AppStubConfig) {
  return requestClient
    .put<{ config: string }>(PATHS[key], { config: stringifyConfig(config) })
    .then((data) => parseConfig(data.config));
}

export const APP_MANAGE_CODES: Record<AppConfigKey, string> = {
  routine: 'app.routine.manage',
  'wechat-reply': 'app.wechat_reply.manage',
  'wechat-menus': 'app.wechat_menus.manage',
  'wechat-template': 'app.wechat_template.manage',
  'wechat-news': 'app.wechat_news.manage',
};
