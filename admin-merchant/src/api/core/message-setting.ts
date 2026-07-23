import { requestClient } from '#/api/request';

export type MessageChannelType = 'mp' | 'mt' | 'sms' | 'wx';

export interface MessageListItem {
  message_id: number;
  message_name: string;
  message_settings_id: number | null;
  message_type: number | { text: string; value: number };
  mp_status: number;
  mp_template: string | null;
  mt_status: number;
  mt_template: string | null;
  remark: string;
  sms_status: number;
  sms_template: string | null;
  wx_status: number;
  wx_template: string | null;
}

export interface MessageFieldItem {
  field_ename: string;
  field_name: string;
  field_new_ename?: string;
  filed_new_value?: string;
  filed_value: string;
  is_var: number;
  message_field_id: number;
  message_id: number;
}

export interface MessageTemplateSettings {
  template_id: string;
  var_data: Record<
    string,
    {
      field_name: string;
      filed_value: string;
    }
  >;
}

export const MESSAGE_TYPE_LABELS: Record<number, string> = {
  10: '订单',
  20: '分销',
  30: '通知',
};

export function resolveMessageTypeLabel(value: MessageListItem['message_type']) {
  if (value && typeof value === 'object' && 'text' in value) {
    return value.text;
  }
  return MESSAGE_TYPE_LABELS[Number(value)] ?? String(value ?? '');
}

export async function getMessageListApi(messageTo: number) {
  return requestClient.post<{ list: MessageListItem[] }>('/shop/setting.message/index', {
    message_to: messageTo,
  });
}

export async function getMessageFieldListApi(messageId: number, messageType: MessageChannelType) {
  return requestClient.post<{ list: MessageFieldItem[]; settings: MessageTemplateSettings | null }>(
    '/shop/setting.message/field',
    {
      message_id: messageId,
      message_type: messageType,
    },
  );
}

export async function saveMessageSettingsApi(payload: {
  fieldList: MessageFieldItem[];
  message_id: number;
  message_type: MessageChannelType;
  template_id: string;
}) {
  return requestClient.post('/shop/setting.message/saveSettings', {
    fieldList: JSON.stringify(payload.fieldList),
    message_id: payload.message_id,
    message_type: payload.message_type,
    template_id: payload.template_id,
  });
}

export async function updateMessageSettingsStatusApi(
  messageSettingsId: number,
  messageType: MessageChannelType,
) {
  return requestClient.post('/shop/setting.message/updateSettingsStatus', {
    message_settings_id: messageSettingsId,
    message_type: messageType,
  });
}
