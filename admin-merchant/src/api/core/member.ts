import { requestClient } from '#/api/request';

export interface PaginatedList<T> {
  data: T[];
  total: number;
}

export interface MemberGrade {
  grade_id: number | string;
  name: string;
}

export interface MemberTag {
  tag_id: number | string;
  tag_name: string;
}

export interface MemberListItem {
  avatarUrl?: string;
  balance?: number | string;
  create_time: string;
  grade?: { name: string };
  grade_id?: number;
  mobile?: string;
  nickName: string;
  pay_money?: number | string;
  points?: number | string;
  referee?: { nickName: string; user_id: number };
  reg_source?: string;
  user_id: number;
}

export interface MemberListQuery {
  grade_id?: number | string;
  list_rows?: number;
  nick_name?: string;
  page?: number;
  reg_date?: '' | string[];
  reg_source?: number | string;
  tag_id?: number | string;
}

export interface MemberListResult {
  allTag: MemberTag[];
  grade: MemberGrade[];
  list: PaginatedList<MemberListItem>;
}

export async function getMemberListApi(params: MemberListQuery) {
  return requestClient.post<MemberListResult>('/shop/user.user/index', params);
}

export async function deleteMemberApi(userId: number) {
  return requestClient.post('/shop/user.user/delete', { user_id: userId });
}

export interface EnumField {
  text: string;
  value: number;
}

export interface GradeListItem {
  create_time?: string;
  equity: number | string;
  font_color?: string;
  give_points?: number | string;
  grade_id: number;
  image?: string;
  is_default?: number;
  name: string;
  open_invite?: number;
  open_money?: number;
  open_points?: number;
  remark?: string;
  upgrade_invite?: number | string;
  upgrade_money?: number | string;
  upgrade_points?: number | string;
  weight: number | string;
}

export interface GradeLogItem {
  change_type: number;
  create_time: string;
  grade?: { name: string };
  log_id?: number;
  oldGrade?: { name: string };
  remark?: string;
  user?: { avatarUrl?: string; nickName?: string; user_id: number };
}

export interface GradeFormPayload extends Record<string, unknown> {
  equity: number | string;
  font_color?: string;
  give_points?: number | string;
  image: string;
  name: string;
  open_invite?: number | boolean;
  open_money?: number | boolean;
  open_points?: number | boolean;
  upgrade_invite?: number | string;
  upgrade_money?: number | string;
  upgrade_points?: number | string;
  weight: number | string;
}

export async function getGradeListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<GradeListItem> }>(
    '/shop/user.grade/index',
    params,
  );
}

export async function getGradeLogListApi(params: {
  list_rows?: number;
  new_grade_id?: number;
  old_grade_id?: number;
  page?: number;
  search?: string;
}) {
  return requestClient.post<{
    gradeList: MemberGrade[];
    list: PaginatedList<GradeLogItem>;
  }>('/shop/user.grade/log', params);
}

export async function addGradeApi(payload: GradeFormPayload) {
  return requestClient.post('/shop/user.grade/add', payload);
}

export async function editGradeApi(payload: GradeFormPayload & { grade_id: number }) {
  return requestClient.post('/shop/user.grade/edit', payload);
}

export async function deleteGradeApi(gradeId: number) {
  return requestClient.post('/shop/user.grade/delete', { grade_id: gradeId });
}

export interface TagListItem {
  create_time: string;
  tag_id: number;
  tag_name: string;
  user_count: number;
}

export async function getTagListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<TagListItem> }>(
    '/shop/user.tag/index',
    params,
  );
}

export async function addTagApi(tagName: string) {
  return requestClient.post('/shop/user.tag/add', { tag_name: tagName });
}

export async function editTagApi(payload: { tag_id: number; tag_name: string }) {
  return requestClient.post('/shop/user.tag/edit', payload);
}

export async function deleteTagApi(tagId: number) {
  return requestClient.post('/shop/user.tag/delete', { tag_id: tagId });
}

export interface EquityListItem {
  create_time: string;
  equity_id: number;
  image?: string;
  name: string;
  sort: number | string;
  status: number;
}

export async function getEquityListApi(params: {
  list_rows?: number;
  name?: string;
  page?: number;
  status?: number;
}) {
  return requestClient.post<{ list: PaginatedList<EquityListItem> }>(
    '/shop/user.equity/index',
    params,
  );
}

export async function getEquityDetailApi(equityId: number) {
  return requestClient.get<{ model: EquityListItem }>('/shop/user.equity/edit', {
    params: { equity_id: equityId },
  });
}

export async function addEquityApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/user.equity/add', payload);
}

export async function editEquityApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/user.equity/edit', payload);
}

export async function deleteEquityApi(equityId: number) {
  return requestClient.post('/shop/user.equity/delete', { equity_id: equityId });
}

export interface BalanceLogItem {
  create_time: string;
  describe?: string;
  log_id?: number;
  money: number | string;
  remark?: string;
  scene?: EnumField;
  user?: { avatarUrl?: string; nickName?: string; user_id: number };
  user_id: number;
}

export interface BalanceLogResult {
  attributes?: {
    scene?: Record<string, { name?: string; value?: number }>;
  };
  list: PaginatedList<BalanceLogItem>;
}

export async function getBalanceLogListApi(params: {
  create_time?: string[];
  list_rows?: number;
  page?: number;
  scene?: number;
  search?: string;
}) {
  return requestClient.get<BalanceLogResult>('/shop/user.balance/log', {
    params: {
      list_rows: params.list_rows,
      page: params.page,
      scene: params.scene,
      search: params.search,
      value1: params.create_time,
    },
  });
}

export async function getBalanceSettingApi() {
  return requestClient.get<{ values: Record<string, string> }>(
    '/shop/user.balance/setting',
  );
}

export async function saveBalanceSettingApi(payload: Record<string, string>) {
  return requestClient.post('/shop/user.balance/setting', payload);
}

export interface PointsLogItem {
  create_time: string;
  describe?: string;
  log_id?: number;
  remark?: string;
  user?: { avatarUrl?: string; nickName?: string; user_id: number };
  user_id: number;
  value: number | string;
}

export async function getPointsLogListApi(params: {
  create_time?: string[];
  list_rows?: number;
  page?: number;
  search?: string;
}) {
  return requestClient.get<{ list: PaginatedList<PointsLogItem> }>(
    '/shop/user.points/log',
    {
      params: {
        list_rows: params.list_rows,
        page: params.page,
        search: params.search,
        value1: params.create_time,
      },
    },
  );
}

export async function getPointsSettingApi() {
  return requestClient.get<{ values: Record<string, unknown> }>(
    '/shop/user.points/setting',
  );
}

export async function savePointsSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/user.points/setting', payload);
}

export interface MemberRechargeBalanceParams {
  mode: 'dec' | 'final' | 'inc';
  money?: number | string;
  remark?: string;
}

export interface MemberRechargePointsParams {
  mode: 'dec' | 'final' | 'inc';
  remark?: string;
  value?: number | string;
}

export async function rechargeMemberApi(payload: {
  params: MemberRechargeBalanceParams | MemberRechargePointsParams;
  source: 0 | 1;
  user_id: number;
}) {
  return requestClient.post('/shop/user.user/recharge', {
    params: JSON.stringify(payload.params),
    source: payload.source,
    user_id: payload.user_id,
  });
}

export async function editMemberGradeApi(payload: {
  grade_id: number | string;
  remark?: string;
  user_id: number;
}) {
  return requestClient.post('/shop/user.user/grade', payload);
}

export async function getMemberTagEditApi(userId: number) {
  return requestClient.get<{ allTag: MemberTag[]; userTag: Array<number | string> }>(
    '/shop/user.user/tag',
    { params: { user_id: userId } },
  );
}

export async function editMemberTagApi(payload: {
  checkedTag: Array<number | string>;
  user_id: number;
}) {
  return requestClient.post('/shop/user.user/tag', payload);
}

export interface MemberFormModel {
  avatarUrl?: string;
  gender?: number;
  mobile?: string;
  nickName?: string;
  password?: string;
  reg_source?: string;
  user_id?: number;
}

export interface MemberFormPayload {
  avatarUrl?: string;
  gender: number;
  mobile: string;
  nickName: string;
  password?: string;
  user_id?: number;
}

export async function getMemberEditMetaApi(userId: number) {
  return requestClient.get<{ model: MemberFormModel }>('/shop/user.user/edit', {
    params: { user_id: userId },
  });
}

export async function addMemberApi(payload: MemberFormPayload) {
  return requestClient.post('/shop/user.user/add', payload);
}

export async function editMemberApi(payload: MemberFormPayload & { user_id: number }) {
  return requestClient.post('/shop/user.user/edit', payload);
}
