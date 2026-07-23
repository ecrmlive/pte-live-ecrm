import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface SurfaceTemplateItem {
  create_time: string;
  sort: number;
  template_id: number;
  template_name: string;
  template_num: string;
}

export interface SurfaceSettingItem {
  create_time: string;
  express?: { express_name: string };
  express_id: number;
  express_name?: string;
  partner_id: string;
  setting_id: number;
  setting_name: string;
  sort: number;
}

export interface ExpressOption {
  express_id: number;
  express_name: string;
}

export interface SurfaceTemplateFormValues {
  sort: number | string;
  template_id?: number;
  template_name: string;
  template_num: string;
  wx_code?: string;
}

export interface SurfaceSettingFormValues {
  check_man?: string;
  code?: string;
  express_id: number | string;
  net?: string;
  partner_id: string;
  partner_key?: string;
  partner_name?: string;
  partner_secret?: string;
  setting_id?: number;
  setting_name: string;
  sort: number | string;
}

export async function getSurfaceTemplateListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<SurfaceTemplateItem> }>(
    '/shop/plus.surface.template/index',
    params,
  );
}

export async function addSurfaceTemplateApi(payload: SurfaceTemplateFormValues) {
  return requestClient.post('/shop/plus.surface.template/add', payload);
}

export async function getSurfaceTemplateEditMetaApi(templateId: number) {
  return requestClient.get<{ model: SurfaceTemplateFormValues }>(
    '/shop/plus.surface.template/edit',
    { params: { template_id: templateId } },
  );
}

export async function editSurfaceTemplateApi(payload: SurfaceTemplateFormValues) {
  return requestClient.post('/shop/plus.surface.template/edit', payload);
}

export async function deleteSurfaceTemplateApi(templateId: number) {
  return requestClient.post('/shop/plus.surface.template/delete', { template_id: templateId });
}

export async function getSurfaceSettingListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<SurfaceSettingItem> }>(
    '/shop/plus.surface.setting/index',
    params,
  );
}

export async function getSurfaceSettingAddMetaApi() {
  return requestClient.get<{ expressList: ExpressOption[] }>('/shop/plus.surface.setting/add');
}

export async function addSurfaceSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.surface.setting/add', payload);
}

export async function getSurfaceSettingEditMetaApi(settingId: number) {
  return requestClient.get<{ expressList: ExpressOption[]; model: SurfaceSettingFormValues }>(
    '/shop/plus.surface.setting/edit',
    { params: { setting_id: settingId } },
  );
}

export async function editSurfaceSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.surface.setting/edit', payload);
}

export async function deleteSurfaceSettingApi(settingId: number) {
  return requestClient.post('/shop/plus.surface.setting/delete', { setting_id: settingId });
}
