export interface MessageItem {
  message_id: number;
  message_name: string;
  message_ename: string;
  message_to: { text: string; value: number } | number;
  message_type: { text: string; value: number } | number;
  remark?: string;
  sort?: number;
  create_time?: string;
  children?: MessageItem[];
}

export interface MessageFormModel {
  message_id?: number;
  message_name: string;
  message_ename: string;
  message_to: number;
  message_type: number;
  sort?: number | string;
  remark?: string;
  status?: number;
}

export interface MessageFieldRow {
  message_field_id: number;
  message_id: number;
  field_name: string;
  field_ename: string;
  filed_value: string;
  is_var?: number;
  sort?: number | string;
}
