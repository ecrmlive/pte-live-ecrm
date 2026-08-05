import { requestClient } from '#/api/request';

export interface MerchantProductComment {
  content: string;
  created_at: string;
  id: number;
  product_id: number;
  product_title: string;
  reply_content: string;
  score: number;
  source: string;
  status: string;
  user_id: number;
}

export interface MerchantProductCommentPage {
  limit: number;
  list: MerchantProductComment[];
  page: number;
  total: number;
}

export function listMerchantProductCommentsApi(params: {
  limit: number;
  page: number;
  product_id?: number;
  status?: string;
}) {
  return requestClient.get<MerchantProductCommentPage>('/product/comments', { params });
}

export function replyMerchantProductCommentApi(id: number, body: { reply_content: string }) {
  return requestClient.post<{ ok: boolean }>(`/product/comments/${id}/reply`, body);
}
