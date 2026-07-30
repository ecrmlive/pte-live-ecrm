import { requestClient } from '#/api/request';

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export interface ExpressRow {
  express_id: number;
  name: string;
  code: string;
  sort: number;
  is_show: number;
}

export interface ArticleRow {
  article_id: number;
  cid: number;
  title: string;
  author: string;
  synopsis: string;
  content: string;
  sort: number;
  status: number;
  visit?: number;
}

export interface UserLabelRow {
  label_id: number;
  label_name: string;
  sort: number;
}

/** 平台店铺（CRMEB 商户）管理。 */
export interface PlatformMerchantRow {
  mer_id: number;
  category_id: number;
  mer_name: string;
  real_name: string;
  mer_phone: string;
  mer_address: string;
  mer_info: string;
  mark: string;
  status: number;
  mer_state: number;
  is_audit: number;
  create_time: string;
}

export interface MerchantIntentionRow {
  mer_intention_id: number;
  uid: number;
  phone: string;
  mer_name: string;
  name: string;
  create_time?: string | null;
  status: number;
  fail_msg: string;
  mark: string;
  mer_id: number;
  merchant_category_id: number;
  mer_type_id: number;
}

export function fetchPlatformMerchants(params: {
  keyword?: string;
  limit: number;
  page: number;
  status?: number;
}) {
  return requestClient.get<PageResult<PlatformMerchantRow>>('/merchants', {
    params,
  });
}

export function fetchPlatformMerchant(id: number) {
  return requestClient.get<PlatformMerchantRow>(`/merchants/${id}`);
}

export function updatePlatformMerchantStatus(id: number, enabled: boolean) {
  return requestClient.put<{ ok: boolean }>(`/merchants/${id}/status`, {
    enabled,
  });
}

export function fetchMerchantIntentions(params: {
  keyword?: string;
  limit: number;
  page: number;
  status?: number;
}) {
  return requestClient.get<PageResult<MerchantIntentionRow>>(
    '/merchant-intentions',
    { params },
  );
}

export function auditMerchantIntention(
  id: number,
  data: {
    account?: string;
    fail_msg?: string;
    mark?: string;
    password?: string;
    region_id?: number;
    status: number;
  },
) {
  return requestClient.post(`/merchant-intentions/${id}/audit`, data);
}

export function fetchExpressList(params: { page: number; limit: number }) {
  return requestClient.get<PageResult<ExpressRow>>('/express', { params });
}

export function createExpress(data: {
  name: string;
  code?: string;
  sort?: number;
  is_show?: number;
}) {
  return requestClient.post<ExpressRow>('/express', data);
}

export function updateExpress(
  id: number,
  data: { name: string; code?: string; sort?: number; is_show?: number },
) {
  return requestClient.put<ExpressRow>(`/express/${id}`, data);
}

export function deleteExpress(id: number) {
  return requestClient.delete(`/express/${id}`);
}

export function fetchArticles(params: { page: number; limit: number }) {
  return requestClient.get<PageResult<ArticleRow>>('/articles', { params });
}

export function createArticle(data: Partial<ArticleRow>) {
  return requestClient.post<ArticleRow>('/articles', data);
}

export function updateArticle(id: number, data: Partial<ArticleRow>) {
  return requestClient.put<ArticleRow>(`/articles/${id}`, data);
}

export function deleteArticle(id: number) {
  return requestClient.delete(`/articles/${id}`);
}

export function fetchUserLabels(params: { page: number; limit: number }) {
  return requestClient.get<PageResult<UserLabelRow>>('/user/labels', {
    params,
  });
}

export function createUserLabel(data: { label_name: string; sort?: number }) {
  return requestClient.post<UserLabelRow>('/user/labels', data);
}

export function updateUserLabel(
  id: number,
  data: { label_name: string; sort?: number },
) {
  return requestClient.put<UserLabelRow>(`/user/labels/${id}`, data);
}

export function deleteUserLabel(id: number) {
  return requestClient.delete(`/user/labels/${id}`);
}

export function markUserLabels(uid: number, label_ids: number[]) {
  return requestClient.put(`/user/${uid}/labels`, { label_ids });
}

/** 区域代理 / 商圈管理。 */
export interface BusinessZoneRow {
  circle_id: number;
  pid: number;
  path: string;
  name: string;
  circle_agent_id: number;
  commission_type: number;
  commission_rate: number;
  level: number;
  remark: string;
  sort: number;
  status: number;
  type: number;
  role_id: number;
  create_time: string;
}

export interface BusinessZoneAgentRow {
  circle_agent_id: number;
  uid: number;
  name: string;
  phone: string;
  qualification: string;
  remark: string;
  audit_admin_id: number;
  audit_reason: string;
  audit_time?: string | null;
  status: number;
  payment_method: number;
  payment_name: string;
  payment_account: string;
  payment_bank: string;
  payment_qr_img: string;
  balance: number;
  type: number;
  business_name: string;
  create_time: string;
}

export function fetchBusinessZones(params: { keyword?: string; status?: number; page: number; limit: number }) {
  return requestClient.get<PageResult<BusinessZoneRow>>('/business-zones', { params });
}

export function createBusinessZone(data: Partial<BusinessZoneRow>) {
  return requestClient.post<BusinessZoneRow>('/business-zones', data);
}

export function updateBusinessZone(id: number, data: Partial<BusinessZoneRow>) {
  return requestClient.put<BusinessZoneRow>(`/business-zones/${id}`, data);
}

export function deleteBusinessZone(id: number) {
  return requestClient.delete(`/business-zones/${id}`);
}

export function fetchBusinessZoneAgents(params: { keyword?: string; status?: number; page: number; limit: number }) {
  return requestClient.get<PageResult<BusinessZoneAgentRow>>('/business-zone-agents', { params });
}

export function createBusinessZoneAgent(data: Partial<BusinessZoneAgentRow>) {
  return requestClient.post<BusinessZoneAgentRow>('/business-zone-agents', data);
}

export function updateBusinessZoneAgent(id: number, data: Partial<BusinessZoneAgentRow>) {
  return requestClient.put<BusinessZoneAgentRow>(`/business-zone-agents/${id}`, data);
}

export function auditBusinessZoneAgent(id: number, status: -1 | 1, audit_reason = '') {
  return requestClient.post(`/business-zone-agents/${id}/audit`, { status, audit_reason });
}

/** 平台管理员（含区域管理员绑定）。 */
export interface PlatformAdminRow {
  admin_id: number;
  account: string;
  real_name: string;
  phone: string;
  roles: string;
  status: number;
  level: number;
  region_ids: string;
  is_agent: number;
  circle_agent_id: number;
  create_time?: string;
}

export function fetchPlatformAdmins(params: { page: number; limit: number }) {
  return requestClient.get<PageResult<PlatformAdminRow>>('/setting/admins', { params });
}

export function createPlatformAdmin(data: Partial<PlatformAdminRow> & { password: string }) {
  return requestClient.post<PlatformAdminRow>('/setting/admins', data);
}

export function updatePlatformAdmin(id: number, data: Partial<PlatformAdminRow> & { password?: string }) {
  return requestClient.put<PlatformAdminRow>(`/setting/admins/${id}`, data);
}

export interface PlatformRoleRow {
  role_id: number;
  code: string;
  role_name: string;
  rules: string;
  status: number;
  mer_id: number;
  is_agent: number;
  circle_id: number;
}

export interface PlatformMenuNode {
  menu_id: number;
  menu_name: string;
  children?: PlatformMenuNode[];
}

export function fetchPlatformRoles(params: { page: number; limit: number }) {
  return requestClient.get<PageResult<PlatformRoleRow>>('/setting/roles', { params });
}

export function createPlatformRole(data: { code: string; role_name: string; menu_ids: number[]; status: number; is_agent?: number; circle_id?: number }) {
  return requestClient.post<PlatformRoleRow>('/setting/roles', data);
}

export function updatePlatformRole(id: number, data: { role_name: string; menu_ids: number[]; status: number }) {
  return requestClient.put<PlatformRoleRow>(`/setting/roles/${id}`, data);
}

export async function fetchPlatformMenuTree() {
  const result = await requestClient.get<{ list: PlatformMenuNode[] }>('/setting/menus/tree');
  return result.list || [];
}

export interface PlatformMenuRow {
  menu_id: number;
  pid: number;
  path: string;
  icon: string;
  menu_name: string;
  route: string;
  sort: number;
  is_show: number;
  is_menu: number;
  is_agent: number;
}

export async function fetchPlatformMenus() {
  const result = await requestClient.get<{ list: PlatformMenuRow[] }>('/setting/menus');
  return result.list || [];
}

export function updatePlatformMenu(id: number, data: { menu_name?: string; sort?: number; is_show?: number }) {
  return requestClient.put<PlatformMenuRow>(`/setting/menus/${id}`, data);
}
