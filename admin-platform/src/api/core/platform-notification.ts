import { requestClient } from '#/api/request';

export type NotificationAudience = 'member' | 'store';
export type NotificationChannel = 'mini_program' | 'wechat';

export interface NotificationConfig {
  audience: NotificationAudience;
  mini_program_enabled: 0 | 1;
  mini_program_text: string;
  notice_type: string;
  notification_id: number;
  scene: string;
  sms_enabled: 0 | 1;
  sms_text: string;
  updated_at: string;
  wechat_enabled: 0 | 1;
  wechat_text: string;
}

export interface NotificationConfigSaveInput {
  mini_program_enabled: 0 | 1;
  mini_program_text: string;
  sms_enabled: 0 | 1;
  sms_text: string;
  wechat_enabled: 0 | 1;
  wechat_text: string;
}

export function listNotificationConfigsApi(
  audience: NotificationAudience,
  params: { limit: number; page: number },
) {
  return requestClient.get<{
    limit: number;
    list: NotificationConfig[];
    page: number;
    total: number;
  }>('/notification-configs', { params: { ...params, audience } });
}

export function getNotificationConfigApi(id: number) {
  return requestClient.get<NotificationConfig>(`/notification-configs/${id}`);
}

export function saveNotificationConfigApi(
  id: number,
  data: NotificationConfigSaveInput,
) {
  return requestClient.put<NotificationConfig>(`/notification-configs/${id}`, data);
}

export function syncNotificationTemplatesApi(
  audience: NotificationAudience,
  channel: NotificationChannel,
) {
  return requestClient.post<{ ok: boolean }>('/notification-configs/sync', {
    audience,
    channel,
  });
}
