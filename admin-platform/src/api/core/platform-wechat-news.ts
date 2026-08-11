import { requestClient } from '#/api/request';

export interface WechatNewsItem {
  title: string;
  author: string;
  synopsis: string;
  image: string;
  content: string;
}

export interface WechatNewsRow {
  wechat_news_id: number;
  status: number;
  article: WechatNewsItem[];
  create_time?: string;
  update_time?: string;
}

export function listWechatNewsApi(params?: {
  page?: number;
  limit?: number;
  cate_name?: string;
}) {
  return requestClient.get<{ list: WechatNewsRow[]; count: number }>(
    '/wechat/news',
    { params },
  );
}

export function getWechatNewsApi(id: number) {
  return requestClient.get<WechatNewsRow>(`/wechat/news/${id}`);
}

export function createWechatNewsApi(data: {
  status?: number;
  data: WechatNewsItem[];
}) {
  return requestClient.post<{ ok: boolean; wechat_news_id?: number }>(
    '/wechat/news',
    data,
  );
}

export function updateWechatNewsApi(
  id: number,
  data: { status?: number; data: WechatNewsItem[] },
) {
  return requestClient.put<{ ok: boolean }>(`/wechat/news/${id}`, data);
}

export function deleteWechatNewsApi(id: number) {
  return requestClient.delete(`/wechat/news/${id}`);
}
