import type { SystemField } from './types';

/** 对齐 CRMEB apply-setting READONLY_LIST：系统固定字段，不可删改 */
export const SYSTEM_FIELDS: SystemField[] = [
  {
    title: '商户名称',
    placeholder: '请输入商户名称',
    kind: 'text',
    required: true,
  },
  {
    title: '联系人姓名',
    placeholder: '请输入联系人姓名',
    kind: 'text',
    required: true,
  },
  {
    title: '联系人电话',
    placeholder: '请输入联系人电话',
    kind: 'text',
    required: true,
  },
  {
    title: '店铺分类',
    placeholder: '请选择店铺分类',
    kind: 'select',
    required: true,
  },
  {
    title: '店铺类型',
    placeholder: '请选择店铺类型',
    kind: 'select',
    required: true,
  },
  {
    title: '平台分类',
    placeholder: '请选择平台分类',
    kind: 'select',
    required: true,
  },
  {
    title: '商品类型',
    placeholder: '请选择商品类型',
    kind: 'select',
    required: true,
  },
  {
    title: '申请资质',
    placeholder: '请上传申请资质',
    kind: 'image',
    required: true,
    imageSlots: 3,
  },
  {
    title: '申请说明',
    placeholder: '请输入申请说明',
    kind: 'textarea',
    required: true,
  },
];
