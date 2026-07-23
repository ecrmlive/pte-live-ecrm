import { requestClient } from '#/api/request';

export interface PaginatedList<T> {
  data: T[];
  total: number;
}

export interface AgentApplyStatus {
  text: string;
  value: number;
}

export interface AgentApplyType {
  text: string;
  value: number;
}

export interface AgentApplyUser {
  avatarUrl?: string;
  nickName?: string;
  user_id: number;
}

export interface AgentApplyReferee {
  nickName?: string;
  user_id: number;
}

export interface AgentApplyItem {
  apply_id: number;
  apply_status: AgentApplyStatus;
  apply_time: string;
  apply_type: AgentApplyType;
  mobile: string;
  nickName?: string;
  real_name: string;
  referee?: AgentApplyReferee | null;
  referee_id: number;
  reject_reason?: string;
  user?: AgentApplyUser;
  user_id: number;
}

export interface AgentApplyListResult {
  apply_list: PaginatedList<AgentApplyItem>;
}

export async function getAgentApplyListApi(params: {
  list_rows?: number;
  nick_name?: string;
  page?: number;
}) {
  return requestClient.post<AgentApplyListResult>(
    '/shop/plus.agent.apply/index',
    params,
  );
}

export async function editAgentApplyStatusApi(data: {
  apply_id: number;
  apply_status: number | string;
  mobile?: string;
  real_name?: string;
  referee_id?: number;
  reject_reason?: string;
  user_id?: number;
}) {
  return requestClient.post('/shop/plus.agent.apply/editApplyStatus', data);
}

export async function getAgentSettingApi() {
  return requestClient.post<AgentSettingResult>(
    '/shop/plus.agent.setting/index',
    {},
  );
}

export interface AgentPayTypeOption {
  name: string;
  value: number;
}

export interface AgentSettingSection<T> {
  values: T;
}

export interface AgentApplyProductRow {
  product_id: number;
  product_image?: string;
  product_name: string;
  product_price?: number | string;
  product_stock?: number;
}

export interface AgentBasicSetting {
  apply_condition?: string;
  apply_type?: string;
  become?: string;
  bind_time?: number | string;
  bind_type?: string;
  bind_user?: string;
  downline?: string;
  is_open?: string;
  level?: string;
  productList?: AgentApplyProductRow[];
  product_ids?: number[];
  self_buy?: string;
  store_status?: string;
}

export interface AgentCommissionSetting {
  first_money?: number | string;
  money_type?: string;
  second_money?: number | string;
  settle_days?: number | string;
}

export interface AgentSettlementSetting {
  alipay_type?: string;
  auto_audit?: string;
  cash_ratio?: number | string;
  max_money?: number | string;
  min_money?: number | string;
  pay_type?: number[];
  scene_id?: string;
  wechat_pay_auto?: string;
}

export interface AgentWordField {
  value: string;
}

export interface AgentWordsSetting {
  apply: {
    title: AgentWordField;
    words: Record<string, AgentWordField>;
  };
  cash_apply: {
    title: AgentWordField;
    words: Record<string, AgentWordField>;
  };
  cash_list: {
    title: AgentWordField;
    words: Record<string, AgentWordField>;
  };
  grade: { title: AgentWordField };
  index: {
    title: AgentWordField;
    words: Record<string, AgentWordField>;
  };
  order: {
    title: AgentWordField;
    words: Record<string, AgentWordField>;
  };
  product: { title: AgentWordField };
  qrcode: { title: AgentWordField };
  team: { title: AgentWordField };
}

export interface AgentLicenseSetting {
  license?: string;
}

export interface AgentBackgroundSetting {
  apply?: string;
  index?: string;
  product?: string;
  share?: string;
  store?: string;
}

export interface AgentSettingResult {
  data: {
    background: AgentSettingSection<AgentBackgroundSetting>;
    basic: AgentSettingSection<AgentBasicSetting>;
    commission: AgentSettingSection<AgentCommissionSetting>;
    license: AgentSettingSection<AgentLicenseSetting>;
    settlement: AgentSettingSection<AgentSettlementSetting>;
    words: AgentSettingSection<AgentWordsSetting>;
  };
  pay_type?: AgentPayTypeOption[];
}

export async function saveAgentBasicSettingApi(payload: Record<string, unknown>) {
  return requestClient.post<{ msg?: string }>('/shop/plus.agent.setting/basic', payload);
}

export async function saveAgentCommissionSettingApi(payload: Record<string, unknown>) {
  return requestClient.post<{ msg?: string }>('/shop/plus.agent.setting/commission', payload);
}

export async function saveAgentSettlementSettingApi(payload: { form: AgentSettlementSetting }) {
  return requestClient.post<{ msg?: string }>('/shop/plus.agent.setting/settlement', payload);
}

export async function saveAgentWordsSettingApi(payload: AgentWordsSetting) {
  return requestClient.post<{ msg?: string }>('/shop/plus.agent.setting/words', payload);
}

export async function saveAgentLicenseSettingApi(payload: AgentLicenseSetting) {
  return requestClient.post<{ msg?: string }>('/shop/plus.agent.setting/license', payload);
}

export async function saveAgentBackgroundSettingApi(payload: AgentBackgroundSetting) {
  return requestClient.post<{ msg?: string }>('/shop/plus.agent.setting/background', payload);
}

export interface AgentGradeOption {
  grade_id: number;
  name: string;
}

export interface AgentUserBaseData {
  buyMoney: number | string;
  buyOrder: number | string;
  cashNum: number | string;
  directMoney: number | string;
  directOrder: number | string;
  orderMoney: number | string;
  orderNum: number | string;
  referNum: number | string;
  totalMoney: number | string;
  totalNum: number | string;
}

export interface AgentUserReferee {
  nickName?: string;
  user_id: number;
}

export interface AgentUserGrade {
  grade_id?: number;
  name?: string;
}

export interface AgentUserItem {
  avatarUrl?: string;
  cashNum: number;
  create_time: string;
  grade?: AgentUserGrade | null;
  grade_id: number;
  mobile: string;
  money: number;
  nickName?: string;
  orderMoney: number | string;
  orderNum: number;
  real_name: string;
  referee?: AgentUserReferee | null;
  referee_id: number;
  referNum: number;
  totalMoney: string;
  total_money: number;
  user_id: number;
}

export interface AgentUserListResult {
  baseData: AgentUserBaseData;
  gradeList: AgentGradeOption[];
  list: PaginatedList<AgentUserItem>;
}

export async function getAgentUserListApi(params: {
  create_time?: string[];
  list_rows?: number;
  nick_name?: string;
  page?: number;
}) {
  return requestClient.post<AgentUserListResult>(
    '/shop/plus.agent.user/index',
    params,
  );
}

export async function deleteAgentUserApi(userId: number) {
  return requestClient.post('/shop/plus.agent.user/delete', {
    user_id: userId,
  });
}

export interface AgentGradeItem {
  condition_type?: string;
  create_time?: string;
  first_percent: number;
  font_color?: string;
  grade_id: number;
  image?: string;
  is_default: number;
  name: string;
  remark?: string;
  second_percent: number;
  third_percent?: number;
  weight: number;
}

export interface AgentGradeLogItem {
  agent?: { nickName?: string; real_name?: string };
  change_type: number;
  create_time: string;
  grade?: { name?: string };
  log_id: number;
  oldGrade?: { name?: string };
  user_id: number;
}

export async function getAgentGradeListApi(params: {
  list_rows?: number;
  page?: number;
}) {
  return requestClient.post<{ list: PaginatedList<AgentGradeItem> }>(
    '/shop/plus.agent.grade/index',
    params,
  );
}

export async function getAgentGradeLogListApi(params: {
  list_rows?: number;
  page?: number;
}) {
  return requestClient.post<{ list: PaginatedList<AgentGradeLogItem> }>(
    '/shop/plus.agent.grade/log',
    params,
  );
}

export async function deleteAgentGradeApi(gradeId: number) {
  return requestClient.post('/shop/plus.agent.grade/delete', {
    grade_id: gradeId,
  });
}

export interface AgentProductCategory {
  category_id: number;
  child?: AgentProductCategory[];
  name: string;
}

export interface AgentProductItem {
  image?: Array<{ file_path: string }>;
  is_agent: number;
  is_ind_agent?: number;
  product_id: number;
  product_name: string;
  product_price: string;
  product_stock: number;
  sales_actual: number;
}

export interface AgentProductListResult {
  category: AgentProductCategory[];
  list: PaginatedList<AgentProductItem>;
}

export async function getAgentProductListApi(params: {
  category_id?: number | string;
  is_agent?: number | string;
  list_rows?: number;
  page?: number;
  product_name?: string;
}) {
  return requestClient.get<AgentProductListResult>(
    '/shop/plus.agent.product/index',
    { params },
  );
}

export async function setAgentProductsApi(
  productIds: number[],
  isAgent: number,
) {
  return requestClient.post('/shop/plus.agent.product/setAgent', {
    is_agent: isAgent,
    productIds: productIds.join(','),
  });
}

export interface AgentOrderItem {
  agent_first?: { nickName?: string; user_id: number };
  agent_second?: { nickName?: string; user_id: number };
  create_time: string;
  first_money: number;
  first_user_id: number;
  id: number;
  is_invalid: number;
  is_settled: number;
  order_id: number;
  order_master?: {
    order_no: string;
    pay_price: string;
    user?: { mobile?: string; nickName?: string; user_id: number };
  };
  order_no: string;
  pay_price: number;
  second_money: number;
  second_user_id: number;
}

export async function getAgentOrderListApi(params: {
  is_settled?: number | string;
  list_rows?: number;
  order_no?: string;
  page?: number;
  search?: string;
  user_id?: number | string;
}) {
  return requestClient.post<{ list: PaginatedList<AgentOrderItem> }>(
    '/shop/plus.agent.order/index',
    params,
  );
}

function saveAgentExportBlob(blob: Blob, prefix: string) {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `${prefix}-${Date.now()}.csv`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}

export async function exportAgentOrderApi(params: {
  is_settled?: number | string;
  order_no?: string;
  search?: string;
  user_id?: number | string;
}) {
  const blob = await requestClient.download<Blob>(
    '/shop/plus.agent.order/export',
    { params },
  );
  saveAgentExportBlob(blob, 'agent-order');
}

export interface AgentCashOption {
  name: string;
  text?: string;
  value: number;
}

export interface AgentCashItem {
  alipay_account?: string;
  alipay_name?: string;
  apply_status: { name?: string; text?: string; value: number };
  audit_time?: string;
  avatarUrl?: string;
  bank_account?: string;
  bank_card?: string;
  bank_name?: string;
  cancelStatus?: number;
  cash_money?: number;
  create_time: string;
  id: number;
  mobile?: string;
  money: number;
  nickName?: string;
  pay_type: { name?: string; text?: string; value: number };
  real_money?: number;
  real_name?: string;
  reject_reason?: string;
  user_id: number;
}

export interface AgentCashListResult {
  alipay_type?: number;
  applyStatus: AgentCashOption[];
  list: PaginatedList<AgentCashItem>;
  payType: AgentCashOption[];
}

export async function getAgentCashListApi(params: {
  apply_status?: number | string;
  list_rows?: number;
  page?: number;
  pay_type?: number | string;
  search?: string;
  user_id?: number | string;
}) {
  return requestClient.post<AgentCashListResult>(
    '/shop/plus.agent.cash/index',
    params,
  );
}

export async function submitAgentCashApi(data: {
  apply_status: number | string;
  id: number;
  reject_reason?: string;
}) {
  return requestClient.post('/shop/plus.agent.cash/submit', data);
}

export async function confirmAgentCashMoneyApi(id: number) {
  return requestClient.post('/shop/plus.agent.cash/money', { id });
}

export async function agentCashWechatPayApi(id: number) {
  return requestClient.post('/shop/plus.agent.cash/wechat_pay', { id });
}

export async function cancelAgentCashPayApi(id: number) {
  return requestClient.post('/shop/plus.agent.cash/cancel', { id });
}

export async function exportAgentCashApi(params: {
  apply_status?: number | string;
  pay_type?: number | string;
  search?: string;
  user_id?: number | string;
}) {
  const blob = await requestClient.download<Blob>(
    '/shop/plus.agent.cash/export',
    { params },
  );
  saveAgentExportBlob(blob, 'agent-cash');
}

export interface AgentPosterItem {
  create_time: string;
  poster_id: number;
  poster_name: string;
  sort: number;
  status: number;
  update_time?: string;
}

export async function getAgentPosterListApi(params: {
  list_rows?: number;
  page?: number;
}) {
  return requestClient.post<{ list: PaginatedList<AgentPosterItem> }>(
    '/shop/plus.agent.poster/index',
    params,
  );
}

export async function deleteAgentPosterApi(posterId: number) {
  return requestClient.post('/shop/plus.agent.poster/delete', {
    poster_id: posterId,
  });
}

export async function setAgentPosterStateApi(posterId: number, status: number) {
  return requestClient.post('/shop/plus.agent.poster/state', {
    poster_id: posterId,
    status,
  });
}

export async function getAgentPosterEditMetaApi(posterId: number) {
  return requestClient.get<Record<string, unknown>>('/shop/plus.agent.poster/edit', {
    params: { poster_id: posterId },
  });
}

export interface AgentPosterElementBox {
  left: number;
  style?: string;
  top: number;
  width?: number;
}

export interface AgentPosterNickNameBox {
  color: string;
  fontSize: number;
  left: number;
  top: number;
}

export interface AgentPosterData {
  avatar: AgentPosterElementBox;
  backdrop: { src: string };
  nickName: AgentPosterNickNameBox;
  qrcode: AgentPosterElementBox;
}

export interface AgentPosterFormPayload {
  poster_data: AgentPosterData;
  poster_image: string;
  poster_name: string;
  poster_id?: number;
  sort: number | string;
}

export function createDefaultAgentPosterData(): AgentPosterData {
  return {
    avatar: { width: 60, style: 'circle', left: 10, top: 460 },
    backdrop: { src: '' },
    nickName: { fontSize: 16, color: '#363131', left: 80, top: 480 },
    qrcode: { width: 80, style: 'circle', left: 205, top: 445 },
  };
}

export async function getAgentPosterAddMetaApi() {
  return requestClient.get<{ data?: AgentPosterData }>('/shop/plus.agent.poster/add');
}

export async function addAgentPosterApi(payload: AgentPosterFormPayload) {
  return requestClient.post('/shop/plus.agent.poster/add', payload);
}

export async function editAgentPosterApi(
  payload: AgentPosterFormPayload & { poster_id: number },
) {
  return requestClient.post('/shop/plus.agent.poster/edit', payload);
}

export interface AgentProductEditMeta {
  agent_product: null | Record<string, unknown>;
  basicSetting: { level: number };
}

export async function getAgentProductEditMetaApi(productId: number) {
  return requestClient.get<AgentProductEditMeta>('/shop/plus.agent.product/edit', {
    params: { product_id: productId },
  });
}

export async function editAgentProductApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.agent.product/edit', payload);
}

export async function addAgentUserApi(payload: {
  mobile: string;
  real_name: string;
  referee_id?: number | string;
  user_id: number | string;
}) {
  return requestClient.post('/shop/plus.agent.user/add', payload);
}

export async function getAgentRefereeCheckApi(userId: number | string) {
  return requestClient.get<{ model: { real_name?: string } | null }>(
    '/shop/plus.agent.user/edit',
    { params: { user_id: userId } },
  );
}

export async function editAgentUserApi(payload: {
  grade_id?: number;
  mobile?: string;
  real_name?: string;
  referee_id?: number | string;
  user_id: number;
}) {
  return requestClient.post('/shop/plus.agent.user/edit', payload);
}

export async function getAgentUserFansApi(params: {
  level: number | string;
  list_rows?: number;
  page?: number;
  user_id: number;
}) {
  return requestClient.post<{ list: PaginatedList<Record<string, unknown>> }>(
    '/shop/plus.agent.user/fans',
    params,
  );
}

export interface AgentGradeAddMeta {
  basicSetting: { level: number };
}

export async function getAgentGradeAddMetaApi() {
  return requestClient.get<AgentGradeAddMeta>('/shop/plus.agent.grade/add');
}

export async function addAgentGradeApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.agent.grade/add', payload);
}

export async function getAgentGradeEditMetaApi() {
  return requestClient.get<AgentGradeAddMeta>('/shop/plus.agent.grade/edit');
}

export async function editAgentGradeApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.agent.grade/edit', payload);
}

export interface AgentTaskTypeOption {
  name: string;
  type_id: number;
}

export interface AgentTaskItem {
  name: string;
  number: number | string;
  sort: number | string;
  status: number;
  task_id: number;
  task_type: number;
  type_name?: string;
}

export async function getAgentTaskListApi(params: {
  grade_id: number;
  list_rows?: number;
  page?: number;
}) {
  return requestClient.post<{
    list: PaginatedList<AgentTaskItem>;
    typeList: AgentTaskTypeOption[];
  }>('/shop/plus.agent.task/index', params);
}

export async function addAgentTaskApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.agent.task/add', payload);
}

export async function editAgentTaskApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.agent.task/edit', payload);
}

export async function deleteAgentTaskApi(taskId: number) {
  return requestClient.post('/shop/plus.agent.task/delete', { task_id: taskId });
}

export async function setAgentTaskStateApi(payload: {
  status: number;
  task_id: number;
}) {
  return requestClient.post('/shop/plus.agent.task/state', payload);
}
