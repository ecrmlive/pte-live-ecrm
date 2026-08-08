import {
  ArrowDown,
  Calendar,
  Checked,
  Clock,
  Document,
  EditPen,
  Location,
  Picture,
  Select as SelectIcon,
} from '@element-plus/icons-vue';

import type {
  ApplyCityLevel,
  ApplyContentType,
  ApplyFieldType,
  ApplyFormField,
  ComponentMeta,
} from './types';

/** 左侧组件库：类型 / 显示名 / 默认 schema 元数据 / 图标 */
export const COMPONENT_LIBRARY: ComponentMeta[] = [
  { type: 'checkbox', label: '多选框', defaultTitle: '多选框', icon: Checked },
  { type: 'city', label: '城市', defaultTitle: '城市', icon: Location },
  { type: 'date', label: '日期', defaultTitle: '日期', icon: Calendar },
  {
    type: 'daterange',
    label: '日期范围',
    defaultTitle: '日期范围',
    icon: Calendar,
  },
  { type: 'radio', label: '单选框', defaultTitle: '单选框', icon: SelectIcon },
  { type: 'select', label: '下拉框', defaultTitle: '下拉框', icon: ArrowDown },
  { type: 'text', label: '文本框', defaultTitle: '文本框', icon: EditPen },
  {
    type: 'textarea',
    label: '多行文本框',
    defaultTitle: '多行文本框',
    icon: Document,
  },
  { type: 'time', label: '时间', defaultTitle: '时间', icon: Clock },
  {
    type: 'timerange',
    label: '时间范围',
    defaultTitle: '时间范围',
    icon: Clock,
  },
  { type: 'image', label: '图片', defaultTitle: '图片', icon: Picture },
];

export const CONTENT_TYPE_OPTIONS: Array<{
  value: ApplyContentType;
  label: string;
}> = [
  { value: 'text', label: '文本' },
  { value: 'mobile', label: '手机号' },
  { value: 'idcard', label: '身份证号' },
  { value: 'email', label: '邮箱' },
  { value: 'number', label: '数字' },
];

export const CITY_LEVEL_OPTIONS: Array<{
  value: ApplyCityLevel;
  label: string;
}> = [
  { value: 'province_city', label: '省市' },
  { value: 'province_city_district', label: '省市区' },
  {
    value: 'province_city_district_street',
    label: '省市区街道',
  },
];

const DATE_LIKE: ApplyFieldType[] = ['date', 'daterange', 'time', 'timerange'];

export function isDateLikeType(type: ApplyFieldType) {
  return DATE_LIKE.includes(type);
}

export function isOptionsType(type: ApplyFieldType) {
  return type === 'radio' || type === 'checkbox' || type === 'select';
}

export function isSelectLikeType(type: ApplyFieldType) {
  return (
    type === 'select' ||
    type === 'city' ||
    type === 'date' ||
    type === 'daterange' ||
    type === 'time' ||
    type === 'timerange'
  );
}

export function fieldTypeLabel(type: ApplyFieldType) {
  return COMPONENT_LIBRARY.find((item) => item.type === type)?.label || type;
}

function newFieldId(type: ApplyFieldType) {
  return `${type}_${Date.now().toString(36)}_${Math.random().toString(36).slice(2, 6)}`;
}

/** 按类型生成默认 schema（对齐 CRMEB defaultConfig） */
export function createField(type: ApplyFieldType): ApplyFormField {
  const meta = COMPONENT_LIBRARY.find((item) => item.type === type);
  const base: ApplyFormField = {
    id: newFieldId(type),
    type,
    title: meta?.defaultTitle || '字段',
    content_type: 'text',
    default_value: '',
    placeholder: '请填写',
    required: false,
  };

  switch (type) {
    case 'text':
      return {
        ...base,
        placeholder: '请填写',
      };
    case 'textarea':
      return {
        ...base,
        placeholder: '请填写',
      };
    case 'city':
      return {
        ...base,
        placeholder: '请选择',
        city_level: 'province_city_district',
      };
    case 'image':
      return {
        ...base,
        placeholder: '点击上传图片',
        max_upload: 8,
      };
    case 'radio':
    case 'checkbox':
    case 'select':
      return {
        ...base,
        placeholder: type === 'select' ? '请选择' : '',
        options: ['选项一', '选项二'],
      };
    case 'date':
    case 'daterange':
      return {
        ...base,
        placeholder: '请选择',
        default_visible: 'show',
        default_mode: 'current',
        specify_value: '',
      };
    case 'time':
    case 'timerange':
      return {
        ...base,
        placeholder: '请选择',
        default_visible: 'show',
        default_mode: 'current',
        specify_value: '',
      };
    default:
      return base;
  }
}

export function cloneField(source: ApplyFormField): ApplyFormField {
  return {
    ...source,
    id: newFieldId(source.type),
    options: source.options ? [...source.options] : undefined,
  };
}

/** 预览区展示文案 */
export function previewDisplayValue(field: ApplyFormField): string {
  if (isDateLikeType(field.type)) {
    if (field.default_visible === 'hide') {
      return field.placeholder || '请选择';
    }
    if (field.default_mode === 'specify' && field.specify_value) {
      return field.specify_value;
    }
    if (field.default_mode === 'current') {
      if (field.type === 'date' || field.type === 'daterange') return '当前日期';
      return '当前时间';
    }
  }
  if (field.default_value) return field.default_value;
  if (field.placeholder) return field.placeholder;
  switch (field.type) {
    case 'city':
      return '请选择城市';
    case 'date':
      return '请选择日期';
    case 'daterange':
      return '请选择日期范围';
    case 'time':
      return '请选择时间';
    case 'timerange':
      return '请选择时间范围';
    case 'select':
      return '请选择';
    case 'image':
      return '点击上传图片';
    default:
      return '请填写';
  }
}
