import { requestClient } from '#/api/request';

export interface DiyPageRow {
  doc?: DiyPageDoc;
  id: number;
  name: string;
  title: string;
  cover_image?: string;
  status: number;
  is_diy: number;
  is_default: number;
  add_time?: string;
  update_time?: string;
}

export interface DiyPageDoc {
  items: Array<Record<string, unknown>>;
  page: Record<string, unknown>;
}

export interface DiyEditorBootstrap {
  defaultData: Record<string, unknown>;
  defaultPage?: Record<string, unknown>;
  jsonData: DiyPageDoc;
  opts?: Record<string, unknown>;
}

export type CategoryDecorationLayout = 'card' | 'grid' | 'list';

export interface CategoryDecoration {
  layout: CategoryDecorationLayout;
}

export interface ProductDetailDecoration {
  config: Record<string, unknown>;
}

export type DiyLinkScope = 'merchant' | 'platform';

export interface DiyPageCategory {
	add_time?: string;
	children?: DiyPageCategory[];
	id: number;
	is_mer: number;
	level: number;
	name: string;
	pid: number;
	sort: number;
	status: number;
	type: string;
}

export interface DiyPageLink {
	add_time?: string;
	cate_id: number;
	category?: DiyPageCategory;
	example: string;
	id: number;
	is_mer: number;
	name: string;
	param: string;
	sort: number;
	status: number;
	url: string;
}

export type DiyEditorMode =
  | 'center-add'
  | 'center-edit'
  | 'custom-add'
  | 'custom-edit'
  | 'home-add'
  | 'home-edit';

export async function listDiyPagesApi(params: {
  page?: number;
  limit?: number;
  is_diy?: number;
  name?: string;
  status?: number;
}) {
  return requestClient.get<{ list: DiyPageRow[]; total: number; page: number; limit: number }>(
    '/diy/pages',
    { params },
  );
}

export async function listDiyDefaultsApi(params?: { page?: number; limit?: number }) {
  return requestClient.get<{ list: DiyPageRow[]; total: number }>('/diy/defaults', { params });
}

export async function getCategoryDecorationApi() {
  return requestClient.get<CategoryDecoration>('/diy/category-decoration');
}

export async function saveCategoryDecorationApi(layout: CategoryDecorationLayout) {
  return requestClient.put<CategoryDecoration>('/diy/category-decoration', { layout });
}

export async function getProductDetailDecorationApi() {
  return requestClient.get<ProductDetailDecoration>('/diy/product-detail-decoration');
}

export async function saveProductDetailDecorationApi(config: Record<string, unknown>) {
  return requestClient.put<ProductDetailDecoration>('/diy/product-detail-decoration', { config });
}

export async function applyDiyDefaultApi(id: number) {
  return requestClient.post<DiyPageRow>(`/diy/defaults/${id}/apply`, {});
}

export async function deleteDiyPageApi(id: number) {
  return requestClient.delete(`/diy/pages/${id}`);
}

export async function activeDiyPageApi(id: number) {
  return requestClient.post(`/diy/pages/${id}/active`, {});
}

export async function copyDiyPageApi(id: number) {
  return requestClient.post<DiyPageRow>(`/diy/pages/${id}/copy`, {});
}

export async function recoveryDiyPageApi(id: number) {
  return requestClient.post<DiyPageRow>(`/diy/pages/${id}/recovery`, {});
}

export async function listDiyPageCategoriesApi(scope: DiyLinkScope) {
	return requestClient.get<{ list: DiyPageCategory[] }>('/diy/page-categories', {
		params: { scope },
	});
}

export async function createDiyPageCategoryApi(
	scope: DiyLinkScope,
	body: Partial<DiyPageCategory>,
) {
	return requestClient.post<DiyPageCategory>('/diy/page-categories', body, { params: { scope } });
}

export async function updateDiyPageCategoryApi(
	id: number,
	scope: DiyLinkScope,
	body: Partial<DiyPageCategory>,
) {
	return requestClient.put<DiyPageCategory>(`/diy/page-categories/${id}`, body, { params: { scope } });
}

export async function deleteDiyPageCategoryApi(id: number, scope: DiyLinkScope) {
	return requestClient.delete(`/diy/page-categories/${id}`, { params: { scope } });
}

export async function listDiyPageLinksApi(
	scope: DiyLinkScope,
	params?: { limit?: number; name?: string; page?: number; status?: number },
) {
	return requestClient.get<{ limit: number; list: DiyPageLink[]; page: number; total: number }>(
		'/diy/page-links',
		{ params: { ...params, scope } },
	);
}

export async function createDiyPageLinkApi(scope: DiyLinkScope, body: Partial<DiyPageLink>) {
	return requestClient.post<DiyPageLink>('/diy/page-links', body, { params: { scope } });
}

export async function updateDiyPageLinkApi(
	id: number,
	scope: DiyLinkScope,
	body: Partial<DiyPageLink>,
) {
	return requestClient.put<DiyPageLink>(`/diy/page-links/${id}`, body, { params: { scope } });
}

export async function deleteDiyPageLinkApi(id: number, scope: DiyLinkScope) {
	return requestClient.delete(`/diy/page-links/${id}`, { params: { scope } });
}

export async function loadDiyEditorApi(mode: DiyEditorMode, pageId?: number) {
  const id = pageId && pageId > 0 ? pageId : 0;
  const path = id > 0 ? `/diy/editor/bootstrap/${id}` : '/diy/editor/bootstrap';
  return requestClient.get<DiyEditorBootstrap>(path);
}

export async function saveDiyEditorApi(
  mode: DiyEditorMode,
  paramsJson: string,
  pageId?: number,
) {
  const doc = JSON.parse(paramsJson) as DiyPageDoc;
  const pageParams = (doc.page?.params || {}) as Record<string, unknown>;
  const name = String(pageParams.name || pageParams.title || '装修页');
  const title = String(pageParams.title || name);
  const isDiy = mode.startsWith('custom') ? 0 : 1;
  const body = {
    name,
    title,
    template_name: mode.startsWith('home') ? 'home' : 'page',
    is_diy: isDiy,
    doc,
  };
  if (pageId && pageId > 0) {
    const row = await requestClient.put<DiyPageRow>(`/diy/pages/${pageId}`, body);
    return { page_id: row?.id ?? pageId };
  }
  const row = await requestClient.post<DiyPageRow>('/diy/pages', body);
  return { page_id: row?.id ?? 0 };
}

export function resolveDiyEditorMode(path: string): DiyEditorMode {
  if (path.includes('/diy/index') || path.includes('/diy/editor')) {
    // query types: 1 home, 0 micro — resolved in page via query; default home-edit/add
    return 'home-edit';
  }
  if (path === '/page/center/add') return 'center-add';
  if (path === '/page/center/edit') return 'center-edit';
  if (path === '/page/page/add') return 'custom-add';
  if (path === '/page/page/edit') return 'custom-edit';
  if (path === '/page/page/addPage') return 'home-add';
  if (path === '/page/page/editPage') return 'home-edit';
  return 'home-edit';
}

export function isCenterDiyMode(mode: DiyEditorMode) {
  return mode === 'center-add' || mode === 'center-edit';
}

export function isHomeDiyMode(mode: DiyEditorMode) {
  return mode === 'home-add' || mode === 'home-edit';
}

/** theme stub for model mixins */
export async function getPageThemeApi() {
  return { vars: { values: { theme: 'red' } } };
}
