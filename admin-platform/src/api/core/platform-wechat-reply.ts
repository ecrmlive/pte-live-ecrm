import { requestClient } from '#/api/request';

export interface WechatReplyRow {
  wechat_reply_id: number;
  key: string;
  type: string;
  content: string;
  status: number;
  sort: number;
  create_time?: string;
  update_time?: string;
}

export function listWechatRepliesApi(params?: {
  page?: number;
  limit?: number;
  keyword?: string;
  kind?: 'keyword' | 'special' | 'all';
}) {
  return requestClient.get<{ list: WechatReplyRow[]; count: number }>(
    '/wechat/replies',
    { params },
  );
}

export function getWechatReplyApi(id: number) {
  return requestClient.get<WechatReplyRow>(`/wechat/replies/${id}`);
}

export function getWechatReplySpecialApi(key: 'subscribe' | 'default') {
  return requestClient.get<WechatReplyRow>(`/wechat/replies/special/${key}`);
}

export function saveWechatReplySpecialApi(
  key: 'subscribe' | 'default',
  data: { content: string; status?: number; type?: string },
) {
  return requestClient.put<WechatReplyRow>(
    `/wechat/replies/special/${key}`,
    data,
  );
}

export function createWechatReplyApi(data: {
  key: string;
  type?: string;
  content: string;
  status?: number;
  sort?: number;
}) {
  return requestClient.post<WechatReplyRow>('/wechat/replies', data);
}

export function updateWechatReplyApi(
  id: number,
  data: {
    key: string;
    type?: string;
    content: string;
    status?: number;
    sort?: number;
  },
) {
  return requestClient.put(`/wechat/replies/${id}`, data);
}

export function setWechatReplyStatusApi(id: number, status: number) {
  return requestClient.put(`/wechat/replies/${id}/status`, { status });
}

export function deleteWechatReplyApi(id: number) {
  return requestClient.delete(`/wechat/replies/${id}`);
}

export function matchWechatReplyApi(key: string) {
  return requestClient.get<{ matched: boolean; reply: WechatReplyRow | null }>(
    '/wechat/replies/match',
    { params: { key } },
  );
}
