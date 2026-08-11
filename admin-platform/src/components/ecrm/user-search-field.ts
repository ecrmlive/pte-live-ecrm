import type { VbenFormProps } from '#/adapter/form';

import { markRaw } from 'vue';

import PrefixedKeywordField, {
  type PrefixedKeywordOption,
  type PrefixedKeywordValue,
} from './prefixed-keyword-field.vue';

export type { PrefixedKeywordOption, PrefixedKeywordValue };
export { default as PrefixedKeywordField } from './prefixed-keyword-field.vue';

/** 平台列表「用户搜索」标准选项（昵称 / 用户ID / 手机号） */
export const USER_SEARCH_OPTIONS: PrefixedKeywordOption[] = [
  { label: '昵称', value: 'nickname' },
  { label: '用户ID', value: 'uid' },
  { label: '手机号', value: 'phone' },
];

export const DEFAULT_USER_SEARCH_VALUE: PrefixedKeywordValue = {
  type: 'nickname',
  keyword: '',
};

type FormSchemaItem = NonNullable<VbenFormProps['schema']>[number];

export type ListUserSearchFormFieldOptions = {
  /** 默认 nickname */
  defaultType?: string;
  fieldName?: string;
  label?: string;
  options?: PrefixedKeywordOption[];
  placeholder?: string;
  typeWidth?: string;
};

/**
 * Vben 列表筛选项：统一「用户搜索」左右布局（类型 Select + 关键词 Input）。
 *
 * @example
 * listFormOptionsDefaults([
 *   listUserSearchFormField(),
 * ])
 */
export function listUserSearchFormField(
  overrides: ListUserSearchFormFieldOptions = {},
): FormSchemaItem {
  const options = overrides.options?.length
    ? overrides.options
    : USER_SEARCH_OPTIONS;
  const defaultType =
    overrides.defaultType || options[0]?.value || 'nickname';
  return {
    component: markRaw(PrefixedKeywordField),
    componentProps: {
      options,
      placeholder: overrides.placeholder ?? '请输入内容',
      typeWidth: overrides.typeWidth ?? '96px',
    },
    defaultValue: {
      type: defaultType,
      keyword: '',
    } satisfies PrefixedKeywordValue,
    fieldName: overrides.fieldName ?? 'user_search',
    label: overrides.label ?? '用户搜索',
  };
}

/**
 * 通用「类型 + 关键词」筛选项（订单搜索等可复用同一 UI）。
 */
export function listPrefixedKeywordFormField(
  overrides: ListUserSearchFormFieldOptions & {
    options: PrefixedKeywordOption[];
  },
): FormSchemaItem {
  return listUserSearchFormField(overrides);
}

/** 从列表表单值解析用户搜索 */
export function parseUserSearch(
  formValues?: Record<string, unknown>,
  fieldName = 'user_search',
): PrefixedKeywordValue {
  const raw = formValues?.[fieldName];
  if (raw && typeof raw === 'object') {
    const obj = raw as Record<string, unknown>;
    return {
      type: String(obj.type ?? obj.field ?? 'nickname').trim() || 'nickname',
      keyword: String(obj.keyword ?? '').trim(),
    };
  }
  return { type: 'nickname', keyword: '' };
}
