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
  circle_id: number;
}

export interface MerchantCategoryRow {
  merchant_category_id: number;
  category_name: string;
  commission_rate: number;
}

export function fetchMerchantCategories() {
  return requestClient.get<{ list: MerchantCategoryRow[] }>('/merchant-categories');
}

export interface MerchantTypeRow {
  id: number;
  name: string;
  type_info: string;
  is_margin: number;
  margin: number;
  description: string;
  remark: string;
  status: number;
  menu_codes: string[];
}

export interface MerchantTypeSaveInput {
  name: string;
  type_info: string;
  is_margin: boolean;
  margin: number;
  description: string;
  remark: string;
  status: boolean;
  menu_codes: string[];
}

export function fetchMerchantTypes() {
  return requestClient.get<{ list: MerchantTypeRow[] }>('/merchant-types');
}

export interface MerchantDepositAccount { merchant_id: number; required_amount: number; available_amount: number; state: string; }
export interface ProductLabelRow { id: number; name: string; description: string; color: string; sort: number; status: number; created_at: string; }
export interface ProductGuaranteeRow { id: number; name: string; content: string; icon_url: string; sort: number; status: number; created_at: string; }
export interface ProductParameterTemplateRow { id: number; name: string; values_json: string; sort: number; status: number; created_at: string; }
export function fetchProductLabels() { return requestClient.get<{ list: ProductLabelRow[] }>('/product/labels'); }
export function saveProductLabel(id: number | undefined, value: Omit<ProductLabelRow, 'id' | 'created_at'>) { return id ? requestClient.put(`/product/labels/${id}`, value) : requestClient.post('/product/labels', value); }
export function deleteProductLabel(id: number) { return requestClient.delete(`/product/labels/${id}`); }
export function fetchProductGuarantees() { return requestClient.get<{ list: ProductGuaranteeRow[] }>('/product/guarantees'); }
export function saveProductGuarantee(id: number | undefined, value: Omit<ProductGuaranteeRow, 'id' | 'created_at'>) { return id ? requestClient.put(`/product/guarantees/${id}`, value) : requestClient.post('/product/guarantees', value); }
export function deleteProductGuarantee(id: number) { return requestClient.delete(`/product/guarantees/${id}`); }
export function fetchProductParameterTemplates() { return requestClient.get<{ list: ProductParameterTemplateRow[] }>('/product/parameter-templates'); }
export function saveProductParameterTemplate(id: number | undefined, value: { name: string; values: string[]; sort: number; status: number }) { return id ? requestClient.put(`/product/parameter-templates/${id}`, value) : requestClient.post('/product/parameter-templates', value); }
export function deleteProductParameterTemplate(id: number) { return requestClient.delete(`/product/parameter-templates/${id}`); }
export type ProductCommentStatus = 'hidden' | 'pending' | 'published';
export interface ProductCommentRow { id: number; product_id: number; store_id: number; score: number; content: string; media: string; reply_content: string; source: 'user' | 'virtual'; virtual_author_name: string; sort: number; status: ProductCommentStatus; created_at: string; product_title: string; }
export function fetchProductComments(params: { page: number; limit: number; product_id?: number; status?: ProductCommentStatus }) { return requestClient.get<PageResult<ProductCommentRow>>('/product/comments', { params }); }
export function fetchProductComment(id: number) { return requestClient.get<ProductCommentRow>(`/product/comments/${id}`); }
export function moderateProductComment(id: number, input: { action: 'hide' | 'publish'; idempotency_key: string; note?: string }) { return requestClient.post<{ comment_id: number; status: ProductCommentStatus }>(`/product/comments/${id}/moderate`, input); }
export interface VirtualProductCommentInput { product_id?: number; score: number; content: string; virtual_author_name: string; sort: number; attachment_ids?: number[]; idempotency_key: string; }
export function createVirtualProductComment(input: Required<VirtualProductCommentInput>) { return requestClient.post<{ comment_id: number; status: ProductCommentStatus }>('/product/comments/virtual', input); }
export function updateVirtualProductComment(id: number, input: Omit<VirtualProductCommentInput, 'product_id'>) { return requestClient.put<{ comment_id: number; status: ProductCommentStatus }>(`/product/comments/${id}/virtual`, input); }
export function sortVirtualProductComment(id: number, input: { sort: number; idempotency_key: string }) { return requestClient.put<{ comment_id: number; status: ProductCommentStatus }>(`/product/comments/${id}/sort`, input); }
export function deleteVirtualProductComment(id: number, input: { idempotency_key: string; note?: string }) { return requestClient.delete<{ comment_id: number; status: string }>(`/product/comments/${id}`, { data: input }); }
export interface UserFeedbackRow { id:number; user_id:number; type:string; content:string; status:'pending'|'replied'|'closed'; reply:string; created_at:string; updated_at:string; }
export function fetchUserFeedback(params:{page:number;limit:number;status?:string}) { return requestClient.get<PageResult<UserFeedbackRow>>('/user-feedback',{params}); }
export function replyUserFeedback(id:number,input:{reply:string;idempotency_key:string}) { return requestClient.post(`/user-feedback/${id}/reply`,input); }
export function closeUserFeedback(id:number,input:{reply?:string;idempotency_key:string}) { return requestClient.post(`/user-feedback/${id}/close`,input); }
export function deleteUserFeedback(id:number,input:{idempotency_key:string}) { return requestClient.delete(`/user-feedback/${id}`,{data:input}); }
export interface UserFeedbackCategory { id:number; name:string; sort:number; status:0|1; created_at:string; updated_at:string; }
export function fetchUserFeedbackCategories() { return requestClient.get<{list:UserFeedbackCategory[]}>('/user-feedback/categories'); }
export function createUserFeedbackCategory(input:{name:string;sort:number;status:0|1;idempotency_key:string}) { return requestClient.post('/user-feedback/categories',input); }
export function updateUserFeedbackCategory(id:number,input:{name:string;sort:number;status:0|1;idempotency_key:string}) { return requestClient.put(`/user-feedback/categories/${id}`,input); }
export function setUserFeedbackCategoryStatus(id:number,input:{status:0|1;idempotency_key:string}) { return requestClient.put(`/user-feedback/categories/${id}/status`,input); }
export function deleteUserFeedbackCategory(id:number,input:{idempotency_key:string}) { return requestClient.delete(`/user-feedback/categories/${id}`,{data:input}); }
export interface PlatformUserRow { id:number; nickname:string; mobile:string; status:0|1; balance:number; points:number; level_name:string; created_at:string; }
export function fetchPlatformUsers(params:{page:number;limit:number;id?:number;keyword?:string;status?:0|1}) { return requestClient.get<PageResult<PlatformUserRow>>('/user-list',{params}); }
export interface PlatformUserExport { file_name:string; content:string; row_count:number; truncated:boolean; }
export function exportPlatformUsers(input:{id?:number;keyword?:string;status?:0|1;reason:string}) { return requestClient.post<PlatformUserExport>('/user-list/export',input); }
export interface PlatformUserDetail { profile: PlatformUserRow & { commission:number }; assets: Array<{id:number;asset_type:'balance'|'points'|'commission';amount:number;reference_type:string;reference_id:string;created_at:string}>; membership_logs: Array<{id:number;level_name:string;previous_level_name:string;change_type:string;note:string;created_at:string}>; signs: Array<{id:number;sign_date:string;points:number;continuous_days:number;created_at:string}>; browse_history: Array<{id:number;product_id:number;store_id:number;store_name:string;title:string;viewed_at:string}>; orders: Array<{id:number;order_no:string;store_name:string;pay_amount:number;total_quantity:number;status:string;created_at:string}>; coupons:Array<{id:number;coupon_id:number;store_id:number;name:string;discount_type:'amount'|'rate';discount_value:number;min_amount:number;status:'unused'|'locked'|'used'|'expired';obtained_at:string;ends_at?:string}>; distribution:{parent_user_id:number;parent_nickname:string;direct_user_count:number;promoter_status:-1|0|1}; }
export interface PlatformUserGroupOption { group_id:number; group_name:string; sort:number; }
export function fetchPlatformUserGroupOptions() { return requestClient.get<{list:PlatformUserGroupOption[]}>('/user-list/groups'); }
export function assignPlatformUserGroups(input:{user_ids:number[];group_id:number;reason:string;idempotency_key:string}) { return requestClient.post('/user-list/groups/assign', input); }
export function assignPlatformUserGroup(id:number,input:{group_id:number;reason:string;idempotency_key:string}) { return requestClient.post(`/user-list/${id}/group`,input); }
export interface PlatformUserLabelOption { label_id:number; label_name:string; sort:number; }
export function fetchPlatformUserLabelOptions() { return requestClient.get<{list:PlatformUserLabelOption[]}>('/user-list/labels'); }
export function assignPlatformUserLabels(input:{user_ids:number[];label_ids:number[];reason:string;idempotency_key:string}) { return requestClient.post('/user-list/labels/assign', input); }
export function assignPlatformUserLabel(id:number,input:{label_ids:number[];reason:string;idempotency_key:string}) { return requestClient.post(`/user-list/${id}/labels`,input); }
export function changePlatformUserStatus(id:number,input:{status:0|1;reason:string;idempotency_key:string}) { return requestClient.post(`/user-list/${id}/status`,input); }
export function createPlatformUser(input:{account:string;password:string;nickname:string;reason:string;idempotency_key:string}) { return requestClient.post<{user_id:number}>('/user-list',input); }
export function updatePlatformUserProfile(id:number,input:{nickname:string;avatar_url:string;gender:0|1|2;bio:string;reason:string;idempotency_key:string}) { return requestClient.put(`/user-list/${id}/profile`,input); }
export function resetPlatformUserPassword(id:number,input:{password:string;reason:string;idempotency_key:string}) { return requestClient.post(`/user-list/${id}/password`,input); }
export function assignPlatformUserPromoters(input:{user_ids:number[];status:0|1;reason:string;idempotency_key:string}) { return requestClient.post('/user-list/promoters/assign',input); }
export function sendPlatformUserInAppNotification(id:number,input:{title:string;body:string;cover_url:string;reason:string;idempotency_key:string}) { return requestClient.post(`/user-list/${id}/notifications`,input); }
export function fetchPlatformUserDetail(id:number) { return requestClient.get<PlatformUserDetail>(`/user-list/${id}/detail`); }
export function adjustPlatformUserAsset(id:number,input:{asset_type:'balance'|'points';amount:number;reason:string;idempotency_key:string}) { return requestClient.post(`/user-list/${id}/assets/adjust`,input); }
export interface PlatformMemberLevel { id:number; name:string; rank:number; status:number; }
export function fetchPlatformMemberLevels() { return requestClient.get<{list:PlatformMemberLevel[]}>('/user-list/member-levels'); }
export function adjustPlatformUserMemberLevel(id:number,input:{level_id:number;reason:string;idempotency_key:string}) { return requestClient.post(`/user-list/${id}/member-level`,input); }
export interface PlatformCouponTemplate { coupon_id:number; name:string; store_id:number; discount_type:'amount'|'rate'; discount_value:number; min_amount:number; }
export function fetchPlatformCouponTemplates() { return requestClient.get<{list:PlatformCouponTemplate[]}>('/user-list/coupon-templates'); }
export function issuePlatformUserCoupon(userId:number,couponId:number,input:{reason:string;idempotency_key:string}) { return requestClient.post(`/user-list/${userId}/coupons/${couponId}/issue`,input); }
export function revokePlatformUserCoupon(userId:number,couponId:number,input:{reason:string;idempotency_key:string}) { return requestClient.post(`/user-list/${userId}/coupons/${couponId}/revoke`,input); }
export function changePlatformUserReferrer(userId:number,input:{parent_user_id:number;reason:string;idempotency_key:string}) { return requestClient.post(`/user-list/${userId}/referrer`,input); }
export interface MerchantDepositRefund { id: number; merchant_id: number; amount: number; status: string; reason: string; review_note: string; payout_reference?: string; created_at: string; }
export function fetchMerchantDeposits() { return requestClient.get<{ list: MerchantDepositAccount[] }>('/merchant-deposits'); }
export interface ProfitsharingApplication { id:number; merchant_id:number; application_no:string; status:string; description:string; review_note:string; created_at:string; }
export function fetchProfitsharingApplications(status?:string){return requestClient.get<{list:ProfitsharingApplication[]}>('/merchant-profitsharing-applications',{params:{status}})}
export function reviewProfitsharingApplication(id:number,approved:boolean,note:string){return requestClient.post(`/merchant-profitsharing-applications/${id}/review`,{approved,note})}
export function saveProfitsharingApplicationNote(id:number,note:string){return requestClient.put(`/merchant-profitsharing-applications/${id}/note`,{note})}
export function fetchMerchantDepositRefunds(status?: string) { return requestClient.get<{ list: MerchantDepositRefund[] }>('/merchant-deposit-refunds', { params: { status } }); }
export function deductMerchantDeposit(merchantId: number, input: { amount: number; reason: string; idempotency_key: string }) { return requestClient.post(`/merchant-deposits/${merchantId}/deduct`, input); }
export function reviewMerchantDepositRefund(id: number, approved: boolean, note: string) { return requestClient.post(`/merchant-deposit-refunds/${id}/${approved ? 'approve' : 'reject'}`, { note }); }
export function markMerchantDepositRefundPaid(id: number, input: { idempotency_key: string; payout_reference: string }) { return requestClient.post(`/merchant-deposit-refunds/${id}/mark-paid`, input); }

export function fetchMerchantType(id: number) {
  return requestClient.get<MerchantTypeRow>(`/merchant-types/${id}`);
}

export function saveMerchantType(id: number | undefined, input: MerchantTypeSaveInput) {
  return id ? requestClient.put<MerchantTypeRow>(`/merchant-types/${id}`, input) : requestClient.post<MerchantTypeRow>('/merchant-types', input);
}

export function setMerchantTypeRemark(id: number, remark: string) {
  return requestClient.put<{ ok: boolean }>(`/merchant-types/${id}/remark`, { remark });
}

export function setMerchantTypeStatus(id: number, enabled: boolean) {
  return requestClient.put<{ ok: boolean }>(`/merchant-types/${id}/status`, { enabled });
}

export function deleteMerchantType(id: number) {
  return requestClient.delete<{ ok: boolean }>(`/merchant-types/${id}`);
}

export function createMerchantCategory(input: { category_name: string; commission_rate: number }) {
  return requestClient.post<MerchantCategoryRow>('/merchant-categories', input);
}

export function updateMerchantCategory(id: number, input: { category_name: string; commission_rate: number }) {
  return requestClient.put<{ ok: boolean }>(`/merchant-categories/${id}`, input);
}

export function deleteMerchantCategory(id: number) {
  return requestClient.delete<{ ok: boolean }>(`/merchant-categories/${id}`);
}

export interface StoreGroupRow {
  id: number;
  parent_id: number;
  path: string;
  level: number;
  name: string;
  sort: number;
  status: number;
  diy_page_id: number;
  positioning_status: number;
  longitude?: number | null;
  latitude?: number | null;
  address: string;
  merchant_count: number;
  merchant_ids: number[];
  children?: StoreGroupRow[];
}

export interface StoreGroupMerchantRow {
  merchant_id: number;
  merchant_name: string;
  region_id: number;
  status: number;
}

export interface StoreGroupSaveInput {
  parent_id: number;
  name: string;
  sort: number;
  status: boolean;
  diy_page_id: number;
  positioning_status: boolean;
  longitude?: number;
  latitude?: number;
  address: string;
  merchant_ids: number[];
}

export interface DiyPageOption {
  id: number;
  name: string;
  status: number;
}

export function fetchStoreGroups(keyword?: string) {
  return requestClient.get<{ list: StoreGroupRow[] }>('/store-groups', { params: { keyword } });
}

export function fetchStoreGroup(id: number) {
  return requestClient.get<StoreGroupRow>(`/store-groups/${id}`);
}

export function saveStoreGroup(id: number | undefined, input: StoreGroupSaveInput) {
  return id
    ? requestClient.put<StoreGroupRow>(`/store-groups/${id}`, input)
    : requestClient.post<StoreGroupRow>('/store-groups', input);
}

export function deleteStoreGroup(id: number) {
  return requestClient.delete<{ ok: boolean }>(`/store-groups/${id}`);
}

export function setStoreGroupStatus(id: number, enabled: boolean) {
  return requestClient.post<{ ok: boolean }>(`/store-groups/${id}/status`, { enabled });
}

export function setStoreGroupTemplate(id: number, diy_page_id: number) {
  return requestClient.post<{ ok: boolean }>(`/store-groups/${id}/template`, { diy_page_id });
}

export function fetchStoreGroupMerchants(id: number) {
  return requestClient.get<{ list: StoreGroupMerchantRow[] }>(`/store-groups/${id}/merchants`);
}

export function fetchPlatformDiyPages() {
  return requestClient.get<PageResult<DiyPageOption>>('/diy/pages', { params: { page: 1, limit: 100 } });
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

export function assignMerchantIntentionRegion(id: number, region_id: number) {
  return requestClient.post<MerchantIntentionRow>(`/merchant-intentions/${id}/assign-region`, { region_id });
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
  /** 结算账号资料只写不回传；仅表示后台是否已配置。 */
  payment_configured: boolean;
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
export interface BusinessZoneAgentOption { circle_agent_id:number; name:string; type:0|1; }
export interface BusinessZoneAgentMerchant { merchant_id:number; merchant_name:string; region_id:number; status:number; }
export function fetchBusinessZoneAgentOptions() { return requestClient.get<{list:BusinessZoneAgentOption[]}>('/business-zone-agents/options'); }
export function fetchBusinessZoneAgentMerchants(id:number) { return requestClient.get<{list:BusinessZoneAgentMerchant[]}>(`/business-zone-agents/${id}/merchants`); }
export interface BusinessZoneAgentSettings { status_counts:{pending:number;approved:number;rejected:number;revoked:number}; review:{platform_review_required:boolean;rejection_reason_required:boolean}; security:{payment_credentials_write_only:boolean;admin_binding_required:boolean;password_min_length:number;password_max_length:number}; revocation:{hard_delete:boolean;blocked_when:string[]}; }
export function fetchBusinessZoneAgentSettings() { return requestClient.get<BusinessZoneAgentSettings>('/business-zone-agents/settings'); }
export function resetBusinessZoneAgentPassword(id:number,input:{password:string;reason:string;idempotency_key:string}) { return requestClient.post<{circle_agent_id:number;replayed:boolean}>(`/business-zone-agents/${id}/password`, input); }

export function createBusinessZoneAgent(data: Partial<BusinessZoneAgentRow>) {
  return requestClient.post<BusinessZoneAgentRow>('/business-zone-agents', data);
}

export function updateBusinessZoneAgent(id: number, data: Partial<BusinessZoneAgentRow>) {
  return requestClient.put<BusinessZoneAgentRow>(`/business-zone-agents/${id}`, data);
}

export function auditBusinessZoneAgent(id: number, status: -1 | 1, audit_reason = '') {
  return requestClient.post(`/business-zone-agents/${id}/audit`, { status, audit_reason });
}
export function revokeBusinessZoneAgent(id: number, input: { reason: string; idempotency_key: string }) {
  return requestClient.delete(`/business-zone-agents/${id}`, { data: input });
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
  merchant_ids: string;
  region_ids: string;
  service_store_ids: string;
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

export function deletePlatformAdmin(id: number) {
  return requestClient.delete<{ ok: true }>(`/setting/admins/${id}`);
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
