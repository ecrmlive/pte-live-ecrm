import type { VbenFormProps } from '#/adapter/form';

/**
 * 列表筛选表单默认布局（对齐店铺列表金标准）。
 * @see admin-platform/src/views/ecrm/merchant/list.vue
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

/** 从 formOptions 值提取通用列表查询参数（日期、关键词、状态、商户 ID） */
export function buildStandardListParams(
  page: { currentPage: number; pageSize: number },
  formValues?: Record<string, unknown>,
  options: { merField?: string; statusValues?: number[] } = {},
) {
  const range = Array.isArray(formValues?.date_range) ? formValues.date_range : [];
  const merField = options.merField ?? 'mer_id';
  const merIdRaw = String(formValues?.[merField] ?? '').trim();
  const statusRaw = formValues?.status;
  const allowedStatus = options.statusValues ?? [0, 1];
  return {
    page: page.currentPage,
    limit: page.pageSize,
    keyword: String(formValues?.keyword ?? '').trim() || undefined,
    [merField]: merIdRaw ? Number(merIdRaw) : undefined,
    status: allowedStatus.includes(Number(statusRaw)) ? Number(statusRaw) : undefined,
    date_from: range[0] as string | undefined,
    date_to: range[1] as string | undefined,
  };
}

export const LIST_MER_ID_FIELD = {
  component: 'Input' as const,
  componentProps: { clearable: true, placeholder: '商户 ID' },
  fieldName: 'mer_id',
  label: '商户 ID',
};

export const LIST_KEYWORD_FIELD = (placeholder = '活动名称 / 商品') => ({
  component: 'Input' as const,
  componentProps: { clearable: true, placeholder },
  fieldName: 'keyword',
  label: '关键词',
});

export const LIST_ENABLE_STATUS_FIELD = (label = '活动状态') => ({
  component: 'Select' as const,
  componentProps: {
    clearable: true,
    options: [
      { label: '启用', value: 1 },
      { label: '停用', value: 0 },
    ],
    placeholder: '全部状态',
  },
  fieldName: 'status',
  label,
});
