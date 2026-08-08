import type { Component } from 'vue';

/**
 * 入驻表单 DIY schema（对齐 CRMEB systemFormDesign / business_application_form）。
 *
 * C 端渲染约定（与后台预览共用同一 JSON）：
 * - `form_fields[]`：自定义字段，按数组顺序渲染在系统固定字段之后
 * - 系统固定字段不落库，由端侧常量渲染
 * - 每字段按 `type` 选择渲染器；未知 type 应跳过
 *
 * 持久化键：`merchant_apply_setting`（含 background_image + form_fields）
 */

export type ApplyFieldType =
  | 'checkbox'
  | 'city'
  | 'date'
  | 'daterange'
  | 'radio'
  | 'select'
  | 'text'
  | 'textarea'
  | 'time'
  | 'timerange'
  | 'image';

/** 文本框「内容」类型，对齐 CRMEB valConfig.tabList */
export type ApplyContentType =
  | 'text'
  | 'mobile'
  | 'idcard'
  | 'email'
  | 'number';

/** 城市默认值级别，对齐 CRMEB citys.valConfig */
export type ApplyCityLevel =
  | 'province_city'
  | 'province_city_district'
  | 'province_city_district_street';

/** 日期/时间类默认值：显示 | 隐藏 */
export type ApplyDefaultVisible = 'show' | 'hide';

/** 日期/时间类默认值内容：当前 | 指定 */
export type ApplyDefaultMode = 'current' | 'specify';

/**
 * 单条自定义字段。
 * 基础字段始终存在；类型专属属性为可选扩展（存 JSON，向后兼容）。
 */
export interface ApplyFormField {
  id: string;
  type: ApplyFieldType;
  title: string;
  /** 文本框内容类型；其它类型可忽略，默认 text */
  content_type: ApplyContentType;
  /** 简单默认值（文本/选项类）；日期时间类优先看 default_* */
  default_value: string;
  placeholder: string;
  required: boolean;
  /** radio / checkbox / select */
  options?: string[];
  /** image：最多上传，对齐 CRMEB numConfig */
  max_upload?: number;
  /** city：省市级别 */
  city_level?: ApplyCityLevel;
  /** date / daterange / time / timerange */
  default_visible?: ApplyDefaultVisible;
  default_mode?: ApplyDefaultMode;
  /** 指定日期/时间时的值 */
  specify_value?: string;
}

export interface ApplyFormSchema {
  form_fields: ApplyFormField[];
}

export type SystemFieldKind = 'text' | 'textarea' | 'select' | 'image';

export interface SystemField {
  title: string;
  placeholder: string;
  kind: SystemFieldKind;
  /** 系统固定字段是否必填（预览展示 *） */
  required: boolean;
  imageSlots?: number;
}

export interface ComponentMeta {
  type: ApplyFieldType;
  label: string;
  defaultTitle: string;
  icon: Component;
}
