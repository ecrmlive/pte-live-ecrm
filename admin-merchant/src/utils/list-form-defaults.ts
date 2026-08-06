import type { VbenFormProps } from '#/adapter/form';

/**
 * 列表筛选表单默认布局（对齐平台店铺列表金标准）。
 * @see docs/acceptance/LAYOUT-FIDELITY-CHECKLIST.md
 */
export function listFormOptionsDefaults(
  schema: NonNullable<VbenFormProps['schema']>,
  overrides: Partial<VbenFormProps> = {},
): VbenFormProps {
  return {
    collapsed: false,
    showCollapseButton: false,
    wrapperClass: 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3',
    schema,
    ...overrides,
  };
}

/** 对齐店铺列表：Element Plus DatePicker + daterange（勿用未注册的 RangePicker） */
export const LIST_DATE_RANGE_FIELD = {
  component: 'DatePicker' as const,
  componentProps: {
    type: 'daterange',
    valueFormat: 'YYYY-MM-DD',
    startPlaceholder: '开始时间',
    endPlaceholder: '结束时间',
    class: 'w-full',
  },
  fieldName: 'date_range',
  label: '选择时间',
};
