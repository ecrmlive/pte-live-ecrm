import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface PageListItem {
  create_time: number | string;
  is_default?: number;
  page_id: number;
  page_name: string;
  page_type: number;
  update_time: number | string;
}

export interface DiySaveResult {
  page_id?: number;
}

export interface DiyPreviewData {
  items: Array<Record<string, unknown>>;
  page?: Record<string, unknown>;
}

export interface HomePageDefaultRow extends PageListItem {
  page_data?: DiyPreviewData;
}

export interface PageCategoryForm {
  category_style: number;
  share_title?: string;
  wind_style: number;
}

export interface PageThemeForm {
  theme: string;
}

export interface TabbarFormData {
  backgroundColor?: string;
  is_auto?: string;
  list?: Array<Record<string, unknown>>;
  textColor?: string;
  textHoverColor?: string;
  type?: string;
}

export async function getCustomPageListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<PageListItem> }>(
    '/shop/page.page/index',
    params,
  );
}

export async function getHomePageListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{
    default?: HomePageDefaultRow;
    list: PaginatedList<PageListItem>;
  }>('/shop/page.page/list', params);
}

export async function deletePageApi(pageId: number) {
  return requestClient.post('/shop/page.page/delete', { page_id: pageId });
}

export async function setCustomHomePageApi(pageId: number) {
  return requestClient.post('/shop/page.page/setHome', { page_id: pageId });
}

export async function setDefaultHomePageApi(pageId: number) {
  return requestClient.post('/shop/page.page/setPage', { page_id: pageId });
}

export async function getPageCategoryApi() {
  return requestClient.get<{ model: PageCategoryForm }>('/shop/page.page/category');
}

export async function savePageCategoryApi(payload: PageCategoryForm) {
  return requestClient.post('/shop/page.page/category', payload);
}

export async function getPageThemeApi() {
  return requestClient.get<{ vars: { values: PageThemeForm } }>('/shop/page.theme/index');
}

export async function savePageThemeApi(payload: PageThemeForm) {
  return requestClient.post('/shop/page.theme/index', payload);
}

export async function getTabbarApi() {
  return requestClient.get<{ vars: { data: TabbarFormData } }>('/shop/page.tabbar/index');
}

export async function saveTabbarApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/page.tabbar/edit', payload);
}

export async function getCenterPageListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{
    default?: HomePageDefaultRow;
    list: PaginatedList<PageListItem>;
  }>('/shop/page.center/index', params);
}

export async function deleteCenterPageApi(pageId: number) {
  return requestClient.post('/shop/page.center/delete', { page_id: pageId });
}

export async function setDefaultCenterPageApi(pageId: number) {
  return requestClient.post('/shop/page.center/set', { page_id: pageId });
}

export type DiyEditorMode =
  | 'center-add'
  | 'center-edit'
  | 'custom-add'
  | 'custom-edit'
  | 'home-add'
  | 'home-edit';

export async function loadDiyEditorApi(mode: DiyEditorMode, pageId?: number) {
  switch (mode) {
    case 'center-add':
      return requestClient.get<Record<string, unknown>>('/shop/page.center/add');
    case 'center-edit':
      return requestClient.get<Record<string, unknown>>('/shop/page.center/edit', {
        params: { page_id: pageId },
      });
    case 'custom-add':
      return requestClient.get<Record<string, unknown>>('/shop/page.page/add');
    case 'custom-edit':
      return requestClient.get<Record<string, unknown>>('/shop/page.page/edit', {
        params: { page_id: pageId },
      });
    case 'home-add':
      return requestClient.get<Record<string, unknown>>('/shop/page.page/addPage');
    case 'home-edit':
      return requestClient.get<Record<string, unknown>>('/shop/page.page/editPage', {
        params: { page_id: pageId },
      });
    default:
      throw new Error(`unknown diy mode: ${mode}`);
  }
}

export async function saveDiyEditorApi(
  mode: DiyEditorMode,
  paramsJson: string,
  pageId?: number,
) {
  const payload: Record<string, unknown> = { params: paramsJson };
  if (pageId) {
    payload.page_id = pageId;
  }
  switch (mode) {
    case 'center-add':
      return requestClient.post<DiySaveResult>('/shop/page.center/add', payload);
    case 'center-edit':
      return requestClient.post<DiySaveResult>('/shop/page.center/edit', payload);
    case 'custom-add':
      return requestClient.post<DiySaveResult>('/shop/page.page/add', payload);
    case 'custom-edit':
      return requestClient.post<DiySaveResult>('/shop/page.page/edit', payload);
    case 'home-add':
      return requestClient.post<DiySaveResult>('/shop/page.page/addPage', payload);
    case 'home-edit':
      return requestClient.post<DiySaveResult>('/shop/page.page/editPage', payload);
    default:
      throw new Error(`unknown diy mode: ${mode}`);
  }
}

export function resolveDiyEditorMode(path: string): DiyEditorMode {
  if (path === '/page/center/add') return 'center-add';
  if (path === '/page/center/edit') return 'center-edit';
  if (path === '/page/page/add') return 'custom-add';
  if (path === '/page/page/edit') return 'custom-edit';
  if (path === '/page/page/addPage') return 'home-add';
  if (path === '/page/page/editPage') return 'home-edit';
  throw new Error(`unsupported diy path: ${path}`);
}

export function isCenterDiyMode(mode: DiyEditorMode) {
  return mode === 'center-add' || mode === 'center-edit';
}

export function isHomeDiyMode(mode: DiyEditorMode) {
  return mode === 'home-add' || mode === 'home-edit';
}

