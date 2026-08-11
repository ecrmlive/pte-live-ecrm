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
  image?: string;
  synopsis: string;
  content: string;
  sort: number;
  status: number;
  visit?: number;
  create_time?: string;
}

export interface UserLabelRow {
  label_id: number;
  label_name: string;
  sort: number;
  create_time?: string;
}

/** 平台店铺（CRMEB 商户）管理。 */
export interface PlatformMerchantRow {
  mer_id: number;
  category_id: number;
  type_id: number;
  business_id: number;
  mer_name: string;
  owner_name: string;
  real_name: string;
  mer_phone: string;
  mer_address: string;
  mer_info: string;
  mer_keyword: string;
  mark: string;
  status: number;
  mer_state: number;
  is_audit: number;
  is_best: number;
  offline_pay: number;
  is_trader: number;
  is_bro_room: number;
  is_bro_goods: number;
  commission_switch: number;
  commission_rate: number;
  mer_account: string;
  sub_mchid: string;
  applyment_id: string;
  care_count: number;
  care_ficti: number;
  sort: number;
  region_id: number;
  category_name: string;
  type_name: string;
  region_name: string;
  deposit_state: string;
  deposit_required: number;
  deposit_available: number;
  type_margin: number;
  type_is_margin: number;
  store_group_ids?: number[];
  goods_type?: string;
  goods_types?: number[];
  platform_category_ids?: string;
  platform_category_id_list?: number[];
  mer_star?: number;
  mer_avatar?: string;
  create_time: string;
}

export interface PlatformMerchantSaveInput {
  mer_name: string;
  owner_name?: string;
  real_name?: string;
  mer_phone?: string;
  mer_address?: string;
  mer_info?: string;
  mer_keyword?: string;
  mark?: string;
  category_id?: number;
  type_id?: number;
  business_id?: number;
  region_id?: number;
  is_best?: boolean;
  offline_pay?: boolean;
  is_trader?: boolean;
  is_audit?: boolean;
  is_bro_room?: boolean;
  is_bro_goods?: boolean;
  commission_switch?: boolean;
  commission_rate?: number;
  mer_account?: string;
  mer_password?: string;
  sub_mchid?: string;
  applyment_id?: string;
  care_count?: number;
  care_ficti?: number;
  sort?: number;
  status?: boolean;
  store_group_ids?: number[];
  goods_types?: number[];
  platform_category_ids?: number[];
  mer_star?: number;
  mer_avatar?: string;
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
  images?: string;
  category_name?: string;
  type_name?: string;
  merchant_category_id: number;
  mer_type_id: number;
  circle_id: number;
}

export interface MerchantCategoryRow {
  merchant_category_id: number;
  category_name: string;
  /** 平台佣金比例，业务语义等同手续费比例 */
  commission_rate: number;
  create_time?: string;
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
  store_count: number;
  created_at?: string;
  updated_at?: string;
  menu_codes: string[];
}

export interface MerchantTypeSaveInput {
  name: string;
  type_info: string;
  is_margin: boolean;
  margin: number;
  description: string;
  remark: string;
  status?: boolean;
  menu_codes: string[];
}

export function fetchMerchantTypes(params?: { keyword?: string }) {
  return requestClient.get<{ list: MerchantTypeRow[] }>('/merchant-types', {
    params,
  });
}

export interface MerchantDepositAccount {
  merchant_id: number;
  merchant_name?: string;
  owner_name?: string;
  type_name?: string;
  category_name?: string;
  required_amount: number;
  available_amount: number;
  payable_amount?: number;
  state: string;
  mark?: string;
  is_trader?: number;
  type_id?: number;
  category_id?: number;
  paid_at?: string | null;
}
export interface MerchantDepositLedger {
  id: number;
  merchant_id: number;
  entry_type: string;
  amount: number;
  balance_after: number;
  reason: string;
  operator_admin_id: number;
  operator_name?: string;
  created_at: string;
}
export function fetchMerchantDepositLedgers(
  merchantId: number,
  params?: { page?: number; limit?: number },
) {
  return requestClient.get<{ list: MerchantDepositLedger[]; total?: number }>(
    `/merchant-deposits/${merchantId}/ledgers`,
    { params },
  );
}
export interface ProductLabelRow { id: number; name: string; description: string; color: string; sort: number; status: number; created_at: string; }
export interface ProductGuaranteeRow {
  id: number;
  name: string;
  content: string;
  icon_url: string;
  sort: number;
  status: number;
  mer_count: number;
  product_count: number;
  created_at: string;
}
export type ProductGuaranteeSaveInput = Pick<
  ProductGuaranteeRow,
  'name' | 'content' | 'icon_url' | 'sort' | 'status'
>;
export interface ProductParameterItem {
  name: string;
  values: string[];
  required: number;
  sort: number;
}
export interface ProductParameterTemplateRow {
  id: number;
  name: string;
  cate_ids: number[];
  cate_names: string[];
  cate_names_text: string;
  params: ProductParameterItem[];
  is_required: number;
  sort: number;
  status: number;
  created_at: string;
}
export type ProductParameterTemplateSaveInput = {
  name: string;
  cate_ids: number[];
  params: ProductParameterItem[];
  sort: number;
  status: number;
};
export interface StoreParameterItem {
  name: string;
  values: string[];
  required: number;
  sort: number;
}
export interface StoreParameterTemplateRow {
  id: number;
  store_id: number;
  mer_id: number;
  mer_name: string;
  template_name: string;
  sort: number;
  created_at: string;
}
export interface StoreParameterTemplateDetail extends StoreParameterTemplateRow {
  cate_id?: number;
  is_required?: number;
  params: StoreParameterItem[];
}
export function fetchProductLabels() { return requestClient.get<{ list: ProductLabelRow[] }>('/product/labels'); }
export function saveProductLabel(id: number | undefined, value: Omit<ProductLabelRow, 'id' | 'created_at'>) { return id ? requestClient.put(`/product/labels/${id}`, value) : requestClient.post('/product/labels', value); }
export function updateProductLabelStatus(id: number, status: number) {
  return requestClient.put(`/product/labels/${id}/status`, { status });
}
export function deleteProductLabel(id: number) { return requestClient.delete(`/product/labels/${id}`); }
export function fetchProductGuarantees(params?: {
  page?: number;
  limit?: number;
  keyword?: string;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<{ list: ProductGuaranteeRow[]; total?: number }>(
    '/product/guarantees',
    { params },
  );
}
export function saveProductGuarantee(
  id: number | undefined,
  value: ProductGuaranteeSaveInput,
) {
  return id
    ? requestClient.put(`/product/guarantees/${id}`, value)
    : requestClient.post('/product/guarantees', value);
}
export function updateProductGuaranteeStatus(id: number, status: number) {
  return requestClient.put(`/product/guarantees/${id}/status`, { status });
}
export function deleteProductGuarantee(id: number) { return requestClient.delete(`/product/guarantees/${id}`); }
export function fetchProductParameterTemplates(params?: {
  page?: number;
  limit?: number;
  name?: string;
  template_name?: string;
  cate_id?: number;
}) {
  return requestClient.get<{
    list: ProductParameterTemplateRow[];
    total?: number;
    page?: number;
    limit?: number;
  }>('/product/parameter-templates', { params });
}
export function fetchProductParameterTemplate(id: number) {
  return requestClient.get<ProductParameterTemplateRow>(
    `/product/parameter-templates/${id}`,
  );
}
export function saveProductParameterTemplate(
  id: number | undefined,
  value: ProductParameterTemplateSaveInput,
) {
  return id
    ? requestClient.put(`/product/parameter-templates/${id}`, value)
    : requestClient.post('/product/parameter-templates', value);
}
export function deleteProductParameterTemplate(id: number) {
  return requestClient.delete(`/product/parameter-templates/${id}`);
}
export interface ProductPriceRuleRow {
  id: number;
  name: string;
  cate_ids: number[];
  cate_names: string[];
  cate_names_text: string;
  is_default: number;
  content: string;
  sort: number;
  status: number;
  created_at: string;
  updated_at: string;
}
export type ProductPriceRuleSaveInput = {
  name: string;
  cate_ids: number[];
  content: string;
  sort: number;
  status: number;
};
export function fetchProductPriceRules(params?: {
  page?: number;
  limit?: number;
  keyword?: string;
  cate_id?: number;
  status?: number;
}) {
  return requestClient.get<{
    list: ProductPriceRuleRow[];
    total?: number;
    page?: number;
    limit?: number;
  }>('/product/price-rules', { params });
}
export function fetchProductPriceRule(id: number) {
  return requestClient.get<ProductPriceRuleRow>(`/product/price-rules/${id}`);
}
export function saveProductPriceRule(
  id: number | undefined,
  value: ProductPriceRuleSaveInput,
) {
  return id
    ? requestClient.put(`/product/price-rules/${id}`, value)
    : requestClient.post('/product/price-rules', value);
}
export function updateProductPriceRuleStatus(id: number, status: number) {
  return requestClient.put(`/product/price-rules/${id}/status`, { status });
}
export function deleteProductPriceRule(id: number) {
  return requestClient.delete(`/product/price-rules/${id}`);
}
export function fetchStoreParameterTemplates(params?: {
  page?: number;
  limit?: number;
  mer_id?: number;
  store_id?: number;
  template_name?: string;
}) {
  return requestClient.get<{
    list: StoreParameterTemplateRow[];
    total: number;
    page: number;
    limit: number;
  }>('/product/store-parameter-templates', { params });
}
export function fetchStoreParameterTemplate(id: number) {
  return requestClient.get<StoreParameterTemplateDetail>(
    `/product/store-parameter-templates/${id}`,
  );
}
export function createStoreParameterTemplate(value: {
  mer_id?: number;
  store_id?: number;
  template_name: string;
  cate_id: number;
  is_required?: number;
  params: StoreParameterItem[];
  sort?: number;
}) {
  return requestClient.post<StoreParameterTemplateDetail>(
    '/product/store-parameter-templates',
    value,
  );
}
export function copyStoreParameterTemplate(
  id: number,
  value: {
    template_name: string;
    cate_ids: number[];
    params: StoreParameterItem[];
    sort: number;
    status?: number;
  },
) {
  return requestClient.post<{
    ok: boolean;
    platform_template_id: number;
    name: string;
  }>(`/product/store-parameter-templates/${id}/copy`, value);
}
export type ProductCommentStatus = 'hidden' | 'pending' | 'published';
export interface ProductCommentRow {
  id: number;
  product_id: number;
  store_id: number;
  score: number;
  content: string;
  media: string;
  reply_content: string;
  source: 'user' | 'virtual';
  virtual_author_name: string;
  virtual_author_avatar?: string;
  sort: number;
  status: ProductCommentStatus;
  created_at: string;
  product_title: string;
  product_cover?: string;
  user_name?: string;
}
export function fetchProductComments(params: {
  page: number;
  limit: number;
  product_id?: number;
  status?: ProductCommentStatus;
  keyword?: string;
  user_name?: string;
  date_from?: string;
  date_to?: string;
  sort_field?: 'score';
  sort_order?: 'asc' | 'desc';
}) {
  return requestClient.get<PageResult<ProductCommentRow>>('/product/comments', { params });
}
export function fetchProductComment(id: number) { return requestClient.get<ProductCommentRow>(`/product/comments/${id}`); }
export function moderateProductComment(id: number, input: { action: 'hide' | 'publish'; idempotency_key: string; note?: string }) { return requestClient.post<{ comment_id: number; status: ProductCommentStatus }>(`/product/comments/${id}/moderate`, input); }
export interface VirtualProductCommentInput {
  product_id?: number;
  score: number;
  content: string;
  virtual_author_name: string;
  virtual_author_avatar?: string;
  sort: number;
  attachment_ids?: number[];
  idempotency_key: string;
}
export function createVirtualProductComment(input: Required<VirtualProductCommentInput>) { return requestClient.post<{ comment_id: number; status: ProductCommentStatus }>('/product/comments/virtual', input); }
export function updateVirtualProductComment(id: number, input: Omit<VirtualProductCommentInput, 'product_id'>) { return requestClient.put<{ comment_id: number; status: ProductCommentStatus }>(`/product/comments/${id}/virtual`, input); }
export function sortVirtualProductComment(id: number, input: { sort: number; idempotency_key: string }) { return requestClient.put<{ comment_id: number; status: ProductCommentStatus }>(`/product/comments/${id}/sort`, input); }
export function deleteVirtualProductComment(id: number, input: { idempotency_key: string; note?: string }) { return requestClient.delete<{ comment_id: number; status: string }>(`/product/comments/${id}`, { data: input }); }
export interface UserFeedbackRow { id:number; user_id:number; type:string; content:string; status:'pending'|'replied'|'closed'; reply:string; created_at:string; updated_at:string; }
export function fetchUserFeedback(params:{page:number;limit:number;status?:string;keyword?:string;date_from?:string;date_to?:string}) { return requestClient.get<PageResult<UserFeedbackRow>>('/user-feedback',{params}); }
export function replyUserFeedback(id:number,input:{reply:string;idempotency_key:string}) { return requestClient.post(`/user-feedback/${id}/reply`,input); }
export function closeUserFeedback(id:number,input:{reply?:string;idempotency_key:string}) { return requestClient.post(`/user-feedback/${id}/close`,input); }
export function deleteUserFeedback(id:number,input:{idempotency_key:string}) { return requestClient.delete(`/user-feedback/${id}`,{data:input}); }
export interface UserFeedbackCategory {
  id: number;
  pid: number;
  name: string;
  sort: number;
  status: 0 | 1;
  created_at: string;
  updated_at: string;
  children?: UserFeedbackCategory[];
}
export function fetchUserFeedbackCategories() {
  return requestClient.get<{ list: UserFeedbackCategory[] }>('/user-feedback/categories');
}
export function createUserFeedbackCategory(input: {
  name: string;
  pid: number;
  sort: number;
  status: 0 | 1;
  idempotency_key: string;
}) {
  return requestClient.post('/user-feedback/categories', input);
}
export function updateUserFeedbackCategory(
  id: number,
  input: {
    name: string;
    pid: number;
    sort: number;
    status: 0 | 1;
    idempotency_key: string;
  },
) {
  return requestClient.put(`/user-feedback/categories/${id}`, input);
}
export function setUserFeedbackCategoryStatus(
  id: number,
  input: { status: 0 | 1; idempotency_key: string },
) {
  return requestClient.put(`/user-feedback/categories/${id}/status`, input);
}
export function deleteUserFeedbackCategory(
  id: number,
  input: { idempotency_key: string },
) {
  return requestClient.delete(`/user-feedback/categories/${id}`, { data: input });
}
export interface PlatformUserRow {
  id: number;
  nickname: string;
  avatar_url?: string;
  mobile: string;
  source_channel?: string;
  status: 0 | 1;
  balance: number;
  points: number;
  level_name: string;
  group_id?: number;
  group_name?: string;
  parent_user_id?: number;
  parent_nickname?: string;
  svip_status?: string;
  svip_expires_at?: string | null;
  is_svip?: number;
  svip_label?: string;
  created_at: string;
}
export function fetchPlatformUsers(params: {
  page: number;
  limit: number;
  id?: number;
  keyword?: string;
  nickname?: string;
  phone?: string;
  status?: 0 | 1;
  label_id?: number;
  fields_type?: string;
  fields_value?: string;
  source_channel?: string;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<PageResult<PlatformUserRow>>('/user-list', { params });
}
export function setPlatformUserListSvip(
  id: number,
  body: { is_svip: number; svip_endtime?: string },
) {
  return requestClient.put(`/user-list/${id}/svip`, body);
}
export interface PlatformUserExport { file_name:string; content:string; row_count:number; truncated:boolean; }
export function exportPlatformUsers(input:{id?:number;keyword?:string;nickname?:string;phone?:string;status?:0|1;date_from?:string;date_to?:string;reason:string}) { return requestClient.post<PlatformUserExport>('/user-list/export',input); }
export interface PlatformUserDetail {
  profile: PlatformUserRow & {
    commission: number;
    gender?: number;
    bio?: string;
  };
  assets: Array<{
    id: number;
    asset_type: 'balance' | 'points' | 'commission';
    amount: number;
    reference_type: string;
    reference_id: string;
    created_at: string;
  }>;
  membership_logs: Array<{
    id: number;
    level_name: string;
    previous_level_name: string;
    change_type: string;
    note: string;
    created_at: string;
  }>;
  signs: Array<{
    id: number;
    sign_date: string;
    points: number;
    continuous_days: number;
    created_at: string;
  }>;
  browse_history: Array<{
    id: number;
    product_id: number;
    store_id: number;
    store_name: string;
    title: string;
    viewed_at: string;
  }>;
  orders: Array<{
    id: number;
    order_no: string;
    store_name: string;
    pay_amount: number;
    total_quantity: number;
    status: string;
    created_at: string;
  }>;
  coupons: Array<{
    id: number;
    coupon_id: number;
    store_id: number;
    name: string;
    discount_type: 'amount' | 'rate';
    discount_value: number;
    min_amount: number;
    status: 'unused' | 'locked' | 'used' | 'expired';
    obtained_at: string;
    ends_at?: string;
  }>;
  distribution: {
    parent_user_id: number;
    parent_nickname: string;
    direct_user_count: number;
    promoter_status: -1 | 0 | 1;
  };
}
export interface PlatformUserGroupOption { group_id:number; group_name:string; sort:number; }
export function fetchPlatformUserGroupOptions() { return requestClient.get<{list:PlatformUserGroupOption[]}>('/user-list/groups'); }
export function assignPlatformUserGroups(input:{user_ids:number[];group_id:number;reason:string;idempotency_key:string}) { return requestClient.post('/user-list/groups/assign', input); }
export function assignPlatformUserGroup(id:number,input:{group_id:number;reason:string;idempotency_key:string}) { return requestClient.post(`/user-list/${id}/group`,input); }
export interface PlatformUserLabelOption { label_id:number; label_name:string; sort:number; }
export function fetchPlatformUserLabelOptions() { return requestClient.get<{list:PlatformUserLabelOption[]}>('/user-list/labels'); }
export function assignPlatformUserLabels(input:{user_ids:number[];label_ids:number[];reason:string;idempotency_key:string}) { return requestClient.post('/user-list/labels/assign', input); }
export function assignPlatformUserLabel(id:number,input:{label_ids:number[];reason:string;idempotency_key:string}) { return requestClient.post(`/user-list/${id}/labels`,input); }
export function changePlatformUserStatus(id:number,input:{status:0|1;reason:string;idempotency_key:string}) { return requestClient.post(`/user-list/${id}/status`,input); }
export function createPlatformUser(input: {
  account: string;
  password: string;
  nickname?: string;
  avatar_url?: string;
  real_name?: string;
  phone?: string;
  id_card?: string;
  gender?: 0 | 1 | 2;
  status?: 0 | 1;
  is_promoter?: 0 | 1;
  reason?: string;
  idempotency_key: string;
}) {
  return requestClient.post<{ user_id: number }>('/user-list', input);
}
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
export interface MerchantDepositRefund {
  id: number;
  merchant_id: number;
  merchant_name?: string;
  owner_name?: string;
  required_amount?: number;
  available_amount?: number;
  amount: number;
  status: string;
  reason: string;
  review_note: string;
  refund_method?: string;
  payout_reference?: string;
  created_at: string;
}
export function fetchMerchantDeposits(params?: {
  tab?: 'pending' | 'funded' | string;
  merchant_id?: number;
  keyword?: string;
  status?: string;
  type_id?: number;
  category_id?: number;
  is_trader?: number;
  date_from?: string;
  date_to?: string;
  page?: number;
  limit?: number;
}) {
  return requestClient.get<{ list: MerchantDepositAccount[]; total?: number }>(
    '/merchant-deposits',
    { params },
  );
}
export interface ProfitsharingApplication {
  id: number;
  merchant_id: number;
  merchant_name: string;
  application_no: string;
  applyment_id: string;
  status: string;
  description: string;
  message: string;
  review_note: string;
  created_at: string;
}
export function fetchProfitsharingApplications(params?: {
  status?: string;
  keyword?: string;
  merchant_id?: number;
  page?: number;
  limit?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<{ list: ProfitsharingApplication[]; total?: number }>(
    '/merchant-profitsharing-applications',
    { params },
  );
}
export function fetchProfitsharingApplication(id: number) {
  return requestClient.get<ProfitsharingApplication>(
    `/merchant-profitsharing-applications/${id}`,
  );
}
export function reviewProfitsharingApplication(
  id: number,
  data: { status: string; note: string },
) {
  return requestClient.post(`/merchant-profitsharing-applications/${id}/review`, data);
}
export function saveProfitsharingApplicationNote(id: number, note: string) {
  return requestClient.put(`/merchant-profitsharing-applications/${id}/note`, {
    note,
  });
}
export function fetchMerchantDepositRefunds(params?: {
  status?: string;
  merchant_id?: number;
  keyword?: string;
  type_id?: number;
  category_id?: number;
  is_trader?: number;
  date_from?: string;
  date_to?: string;
  page?: number;
  limit?: number;
}) {
  return requestClient.get<{ list: MerchantDepositRefund[]; total?: number }>(
    '/merchant-deposit-refunds',
    { params },
  );
}
export function deductMerchantDeposit(merchantId: number, input: { amount: number; reason: string; idempotency_key: string }) { return requestClient.post(`/merchant-deposits/${merchantId}/deduct`, input); }
export function fundMerchantDepositOffline(
  merchantId: number,
  input: { amount?: number; mark?: string; idempotency_key: string },
) {
  return requestClient.post(`/merchant-deposits/${merchantId}/fund-offline`, input);
}
export function reviewMerchantDepositRefund(id: number, approved: boolean, note: string) { return requestClient.post(`/merchant-deposit-refunds/${id}/${approved ? 'approve' : 'reject'}`, { note }); }
export function markMerchantDepositRefundNote(id: number, note: string) {
  return requestClient.post(`/merchant-deposit-refunds/${id}/mark`, { note });
}
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
  /** 列表/详情由后端按 diy_page_id 关联装修页名称 */
  diy_page_name?: string;
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
  contact_name?: string;
  contact_mobile?: string;
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
  add_time?: string;
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

export function fetchPlatformDiyPages(params?: {
  page?: number;
  limit?: number;
  name?: string;
}) {
  return requestClient.get<PageResult<DiyPageOption>>('/diy/pages', {
    params: {
      page: params?.page ?? 1,
      limit: params?.limit ?? 10,
      name: params?.name || undefined,
    },
  });
}

export function fetchPlatformMerchants(params: {
  keyword?: string;
  limit: number;
  page: number;
  status?: number;
  category_id?: number;
  type_id?: number;
  region_id?: number;
  is_best?: number;
  offline_pay?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<PageResult<PlatformMerchantRow>>('/merchants', {
    params,
  });
}

export function fetchPlatformMerchant(id: number) {
  return requestClient.get<PlatformMerchantRow>(`/merchants/${id}`);
}

export interface MerchantOperateLogRow {
  id: number;
  action: string;
  action_label: string;
  terminal: string;
  role_code: string;
  role_label: string;
  operator_id: number;
  operator_name: string;
  created_at: string;
}

export function fetchMerchantOperateLogs(
  id: number,
  params: {
    page: number;
    limit: number;
    terminal?: string;
    start_date?: string;
    end_date?: string;
  },
) {
  return requestClient.get<{ list: MerchantOperateLogRow[]; total: number }>(
    `/merchants/${id}/operate-logs`,
    { params },
  );
}

export function createPlatformMerchant(input: PlatformMerchantSaveInput) {
  return requestClient.post<PlatformMerchantRow>('/merchants', input);
}

export function updatePlatformMerchant(id: number, input: PlatformMerchantSaveInput) {
  return requestClient.put<PlatformMerchantRow>(`/merchants/${id}`, input);
}

export function updatePlatformMerchantStatus(id: number, enabled: boolean) {
  return requestClient.put<{ ok: boolean }>(`/merchants/${id}/status`, {
    enabled,
  });
}

export function updatePlatformMerchantRecommend(id: number, enabled: boolean) {
  return requestClient.put<{ ok: boolean }>(`/merchants/${id}/recommend`, {
    enabled,
  });
}

/** 平台一键登录店铺：签发 store_console JWT，前端新开店铺管理系统。 */
export function loginPlatformMerchant(id: number) {
  return requestClient.post<{
    token: {
      access_token: string;
      refresh_token: string;
      expires_in: number;
    };
    mer_id: number;
    store_id: number;
    account: string;
    mer_name: string;
    store_name: string;
    store_app_id: string;
    path: string;
  }>(`/merchants/${id}/login`);
}

export function fetchMerchantIntentions(params: {
  keyword?: string;
  limit: number;
  page: number;
  status?: number;
  category_id?: number;
  type_id?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<PageResult<MerchantIntentionRow>>(
    '/merchant-intentions',
    { params },
  );
}

export function fetchMerchantIntention(id: number) {
  return requestClient.get<MerchantIntentionRow>(`/merchant-intentions/${id}`);
}

export function createMerchantIntention(data: {
  mer_name: string;
  name: string;
  phone: string;
  merchant_category_id?: number;
  mer_type_id?: number;
  circle_id?: number;
  images?: string;
  category_name?: string;
  type_name?: string;
}) {
  return requestClient.post<MerchantIntentionRow>('/merchant-intentions', data);
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

export function deleteMerchantIntention(id: number) {
  return requestClient.delete<{ ok: boolean }>(`/merchant-intentions/${id}`);
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

export function fetchArticles(params: {
  page: number;
  limit: number;
  title?: string;
  cid?: number;
}) {
  return requestClient.get<PageResult<ArticleRow>>('/articles', { params });
}

export function getArticle(id: number) {
  return requestClient.get<ArticleRow>(`/articles/${id}`);
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

/** 区域代理 / 商圈管理（type=0 区域；type=1 商户主体）。 */
export interface BusinessZoneAgentBrief {
  circle_agent_id: number;
  name: string;
  phone?: string;
}

export interface BusinessZoneMerchantBrief {
  mer_id: number;
  mer_name: string;
  real_name?: string;
  mer_phone?: string;
}

export interface BusinessZoneRow {
  circle_id: number;
  pid: number;
  path: string;
  name: string;
  circle_agent_id: number;
  circle_agent?: BusinessZoneAgentBrief;
  commission_type: number;
  commission_rate: number;
  level: number;
  remark: string;
  sort: number;
  status: number;
  type: number;
  role_id: number;
  business_store_category?: number;
  business_store_type?: number;
  goods_type?: number[];
  platform_category_ids?: number[];
  merchant_count?: number;
  merchant_ids?: number[];
  merchant?: BusinessZoneMerchantBrief[];
  /** 是否存在下级（列表懒加载树用） */
  has_child?: boolean;
  create_time: string;
}

export interface BusinessZoneSaveInput {
  pid?: number;
  name: string;
  circle_agent_id?: number;
  commission_type?: number;
  commission_rate?: number;
  remark?: string;
  sort?: number;
  status?: number;
  type?: number;
  role_id?: number;
  business_store_category?: number;
  business_store_type?: number;
  goods_type?: number[];
  platform_category_ids?: number[];
  merchant_ids?: number[];
}

export interface BusinessZoneAgentCircleBrief {
  circle_id: number;
  name: string;
  type: number;
  status: number;
  commission_type?: number;
  commission_rate?: number;
}

export interface BusinessZoneAgentRow {
  circle_agent_id: number;
  uid: number;
  name: string;
  phone: string;
  /** 身份资质：JSON 图片 URL 数组字符串，或历史纯文本。 */
  qualification: string;
  remark: string;
  /** CRMEB extend JSON；含 avatar。 */
  extend?: string;
  /** 区域代理标识图（来自 extend.avatar）。 */
  avatar?: string;
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
  /** 关联 C 端用户昵称（列表/详情 enrichment）。 */
  nickname?: string;
  /** C 端用户头像（非区域代理标识图）。 */
  avatar_url?: string;
  /** 统一后台登录账号（只读回显）。 */
  account?: string;
  circle_ids?: number[];
  circles?: BusinessZoneAgentCircleBrief[];
  admin?: { admin_id: number; account: string };
}

export function fetchBusinessZones(params: {
  keyword?: string;
  name?: string;
  status?: number;
  type?: number;
  circle_agent_id?: number;
  /** 父级区域 ID；传 0 仅一级区域，传具体 ID 拉下级 */
  pid?: number;
  page: number;
  limit: number;
}) {
  return requestClient.get<PageResult<BusinessZoneRow>>('/business-zones', { params });
}

export function fetchBusinessZone(id: number) {
  return requestClient.get<BusinessZoneRow>(`/business-zones/${id}`);
}

export interface BusinessZoneOptionNode {
  value: number;
  label: string;
  children?: BusinessZoneOptionNode[];
}

/** type=0 区域，type=1 所属商户（CRMEB business_id）。 */
export function fetchBusinessZoneOptions(type?: 0 | 1) {
  return requestClient.get<{ list: BusinessZoneOptionNode[] }>('/business-zones/options', {
    params: type === undefined ? undefined : { type },
  });
}

export function createBusinessZone(data: BusinessZoneSaveInput) {
  return requestClient.post<BusinessZoneRow>('/business-zones', data);
}

export function updateBusinessZone(id: number, data: BusinessZoneSaveInput) {
  return requestClient.put<BusinessZoneRow>(`/business-zones/${id}`, data);
}

export function updateBusinessZoneStatus(id: number, status: number) {
  return requestClient.put<BusinessZoneRow>(`/business-zones/${id}/status`, { status });
}

export function deleteBusinessZone(id: number) {
  return requestClient.delete(`/business-zones/${id}`);
}

/** H5 邀请入驻：返回带 region_id 的入驻页 URL，供前端生成二维码。 */
export interface BusinessZoneInvitePayload {
  circle_id: number;
  name: string;
  site_url: string;
  h5_url: string;
  label: string;
}

export function fetchBusinessZoneInvite(id: number) {
  return requestClient.get<BusinessZoneInvitePayload>(`/business-zones/${id}/invite`);
}

export function fetchBusinessZoneAgents(params: {
  keyword?: string;
  name?: string;
  phone?: string;
  status?: number;
  type?: number;
  uid?: number;
  nickname?: string;
  user_phone?: string;
  date_from?: string;
  date_to?: string;
  page: number;
  limit: number;
}) {
  return requestClient.get<PageResult<BusinessZoneAgentRow>>('/business-zone-agents', { params });
}

export function fetchBusinessZoneAgent(id: number) {
  return requestClient.get<BusinessZoneAgentRow>(`/business-zone-agents/${id}`);
}
export interface BusinessZoneAgentOption { circle_agent_id:number; name:string; phone?:string; type:0|1; }
export interface BusinessZoneAgentMerchant { merchant_id:number; merchant_name:string; region_id:number; status:number; }
export function fetchBusinessZoneAgentOptions(type?: 0 | 1) {
  return requestClient.get<{ list: BusinessZoneAgentOption[] }>('/business-zone-agents/options', {
    params: type === undefined ? undefined : { type },
  });
}
export function fetchBusinessZoneAgentMerchants(id:number) { return requestClient.get<{list:BusinessZoneAgentMerchant[]}>(`/business-zone-agents/${id}/merchants`); }
export interface BusinessZoneAgentSettings { status_counts:{pending:number;approved:number;rejected:number;revoked:number}; review:{platform_review_required:boolean;rejection_reason_required:boolean}; security:{payment_credentials_write_only:boolean;admin_binding_required:boolean;password_min_length:number;password_max_length:number}; revocation:{hard_delete:boolean;blocked_when:string[]}; }
export function fetchBusinessZoneAgentSettings() { return requestClient.get<BusinessZoneAgentSettings>('/business-zone-agents/settings'); }
export function resetBusinessZoneAgentPassword(id:number,input:{password:string;reason:string;idempotency_key:string}) { return requestClient.post<{circle_agent_id:number;replayed:boolean}>(`/business-zone-agents/${id}/password`, input); }

export type BusinessZoneAgentSaveInput = Partial<BusinessZoneAgentRow> & {
  account?: string;
  password?: string;
  circle_ids?: number[];
  /** 区域代理标识图 → 后端写入 extend.avatar */
  avatar?: string;
  /** 区域表单内新增代理人时立即通过审核 */
  auto_approve?: boolean;
};

export function createBusinessZoneAgent(data: BusinessZoneAgentSaveInput) {
  return requestClient.post<BusinessZoneAgentRow>('/business-zone-agents', data);
}

/** 区域列表表单内快速新增代理人（zone.manage 权限，立即通过）。 */
export function createBusinessZoneRegionAgent(data: BusinessZoneAgentSaveInput) {
  return requestClient.post<BusinessZoneAgentRow>('/business-zones/agents', {
    ...data,
    type: 0,
    auto_approve: true,
  });
}

export function updateBusinessZoneAgent(id: number, data: BusinessZoneAgentSaveInput) {
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
  linked_user_id?: number;
  avatar_url?: string;
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

export interface MerchantStoreMenuRow extends PlatformMenuRow {
  code: string;
  component: string;
  is_route: number;
  created_at: string;
}

export type MerchantStoreMenuInput = {
  parent_id?: number;
  code: string;
  name: string;
  path: string;
  component?: string;
  icon?: string;
  is_menu?: number;
  is_route?: number;
  sort?: number;
  status?: number;
};

export async function fetchPlatformMenus() {
  const result = await requestClient.get<{ list: PlatformMenuRow[] }>('/setting/menus');
  return result.list || [];
}

export function updatePlatformMenu(id: number, data: { menu_name?: string; sort?: number; is_show?: number }) {
  return requestClient.put<PlatformMenuRow>(`/setting/menus/${id}`, data);
}

export async function fetchMerchantStoreMenus() {
  const result = await requestClient.get<{ list: MerchantStoreMenuRow[] }>('/merchant-menus');
  return result.list || [];
}

export function createMerchantStoreMenu(data: MerchantStoreMenuInput) {
  return requestClient.post<MerchantStoreMenuRow>('/merchant-menus', data);
}

export function updateMerchantStoreMenu(id: number, data: MerchantStoreMenuInput) {
  return requestClient.put<MerchantStoreMenuRow>(`/merchant-menus/${id}`, data);
}

export function deleteMerchantStoreMenu(id: number) {
  return requestClient.delete<{ menu_id: number }>(`/merchant-menus/${id}`);
}
