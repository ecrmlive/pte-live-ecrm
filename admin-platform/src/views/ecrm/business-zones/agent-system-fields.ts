import type { SystemField } from '../merchant/apply-form-diy/types';

/** 对齐 CRMEB businessZonesSettings READONLY_LIST */
export const AGENT_SYSTEM_FIELDS: SystemField[] = [
  {
    title: '代理名称',
    placeholder: '请输入代理名称',
    kind: 'text',
    required: true,
  },
  {
    title: '联系电话',
    placeholder: '请输入联系电话',
    kind: 'text',
    required: true,
  },
  {
    title: '申请资质',
    placeholder: '请上传申请资质',
    kind: 'image',
    required: true,
    imageSlots: 3,
  },
];
