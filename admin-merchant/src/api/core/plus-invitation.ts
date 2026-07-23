import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface InvitationListItem {
  end_time?: { text: string };
  invitation_gift_id: number;
  inv_type: number;
  name: string;
  partake_num: number;
  start_time?: { text: string };
  status?: { text: string; value: number };
}

export interface InvitationPartakeItem {
  create_time: string;
  name: string;
  partake?: { nickName: string };
  user?: { nickName: string };
}

export interface InvitationReceiveCouponItem {
  coupon_num: number;
  name: string;
}

export interface InvitationReceiveItem {
  balance?: number | string;
  coupon?: InvitationReceiveCouponItem[];
  create_time: string;
  invitation_num?: number;
  invitation_receive_id: number;
  inviteName?: string;
  name: string;
  nickName?: string;
  point?: number | string;
}

export interface InvitationListQuery {
  create_time?: string[];
  list_rows?: number;
  page?: number;
  search?: string;
  status?: number | string;
}

export async function getInvitationListApi(params: InvitationListQuery) {
  return requestClient.post<{ list: PaginatedList<InvitationListItem> }>(
    '/shop/plus.invitation.active/index',
    params,
  );
}

export async function deleteInvitationApi(id: number) {
  return requestClient.post('/shop/plus.invitation.active/delete', { id });
}

export async function sendInvitationApi(id: number) {
  return requestClient.post('/shop/plus.invitation.active/send', { id });
}

export async function endInvitationApi(id: number) {
  return requestClient.post('/shop/plus.invitation.active/end', { id });
}

export async function getInvitationPartakeListApi(params: {
  id?: number;
  list_rows?: number;
  page?: number;
  search?: string;
}) {
  return requestClient.post<{ list: PaginatedList<InvitationPartakeItem> }>(
    '/shop/plus.invitation.active/partake',
    params,
  );
}

export async function getInvitationReceiveListApi(params: {
  create_time?: string[];
  id?: number;
  list_rows?: number;
  page?: number;
  search?: string;
  type: number;
}) {
  return requestClient.post<{ list: PaginatedList<InvitationReceiveItem> }>(
    '/shop/plus.invitation.active/receive',
    params,
  );
}

export async function getInvitationQrcodeApi(params: { id: number; source: string }) {
  return requestClient.post<{ image: string }>('/shop/plus.invitation.active/qrcode', params);
}

export interface InvitationFormPayload {
  [key: string]: unknown;
}

export async function addInvitationApi(payload: InvitationFormPayload) {
  return requestClient.post('/shop/plus.invitation.active/add', payload);
}

export async function editInvitationApi(payload: InvitationFormPayload) {
  return requestClient.post('/shop/plus.invitation.active/edit', payload);
}

function unwrapInvitationEditPayload(res: Record<string, unknown>) {
  const nested = res.data;
  if (nested && typeof nested === 'object' && !Array.isArray(nested)) {
    return { ...(nested as Record<string, unknown>) };
  }
  return { ...res };
}

export async function getInvitationEditMetaApi(invitationGiftId: number) {
  const res = await requestClient.get<Record<string, unknown>>(
    '/shop/plus.invitation.active/edit',
    { params: { invitation_gift_id: invitationGiftId } },
  );
  const data = unwrapInvitationEditPayload(res);
  if (data.reward) {
    data.reward_data = data.reward;
  }
  if (data.mreward) {
    data.member_reward = data.mreward;
  }
  if (data.creward) {
    data.consumption_reward = data.creward;
  }
  const rule = data.rule as { file_path?: string } | undefined;
  if (rule?.file_path) {
    data.rule_file_path = rule.file_path;
  } else {
    data.rule_file_path = '';
    data.rule_id = 0;
  }
  const share = data.share as { file_path?: string } | undefined;
  if (share?.file_path) {
    data.share_file_path = share.file_path;
  } else {
    data.share_file_path = '';
    data.share_image_id = 0;
  }
  const invImage = data.invImage as { file_path?: string } | undefined;
  if (invImage?.file_path) {
    data.inv_image_path = invImage.file_path;
  } else {
    data.inv_image_path = '';
    data.inv_image_id = 0;
  }
  return data;
}
