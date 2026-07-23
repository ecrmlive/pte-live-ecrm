import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface LotteryPrizeItem {
  award_id?: number | string;
  balance?: number | string;
  draw_num?: number;
  image?: string;
  is_default?: number;
  is_play?: number;
  name?: string;
  points?: number | string;
  prize_id?: number;
  prompt?: string;
  stock?: number | string;
  type?: number;
  weight?: number | string;
}

export interface LotterySettingForm {
  content: string;
  coupon_num?: number;
  end_time: string;
  file_path: string;
  grades: number[];
  image_id: number | string;
  is_times?: number;
  lottery_id?: number;
  name?: string;
  points: number | string;
  prize?: LotteryPrizeItem[];
  prize_ids?: number[] | string;
  start_time: string;
  status: number;
  times: number | string;
  total_num: number | string;
  user_type: number;
}

export interface LotteryRecordItem {
  avatarUrl?: string;
  create_time?: string;
  deliveryType?: number;
  express_id?: number;
  express_no?: string;
  lottery_type_text?: string;
  mobile?: string;
  nickName?: string;
  record_id: number;
  record_name?: string;
  remark?: string;
  status?: number;
  user_id?: number;
}

export interface LotteryRecordQuery {
  list_rows?: number;
  page?: number;
  record_name?: string;
  reg_date?: string[];
  search?: string;
  status?: number;
  type?: number;
}

export async function getLotteryRecordListApi(params: LotteryRecordQuery) {
  return requestClient.post<{
    list: PaginatedList<LotteryRecordItem>;
    lotteryType: string[];
  }>('/shop/plus.lottery/record', params);
}

export async function saveLotteryRemarkApi(payload: {
  record_id: number;
  remark: string;
}) {
  return requestClient.post('/shop/plus.lottery/remark', payload);
}

export async function getLotteryExpressApi(payload: {
  express_id?: number;
  express_no?: string;
  record_id: number;
}) {
  return requestClient.post<{
    express: { list: Array<{ context?: string; time?: string }> };
  }>('/shop/plus.lottery/express', payload);
}

function saveExportBlob(blob: Blob) {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `lottery-records-${Date.now()}.xlsx`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}

export async function exportLotteryRecordApi(params: Omit<LotteryRecordQuery, 'list_rows' | 'page'>) {
  const blob = await requestClient.download<Blob>('/shop/plus.lottery/export', {
    params,
  });
  saveExportBlob(blob);
}

export async function getLotterySettingApi() {
  return requestClient.get<{
    data: LotterySettingForm;
    gradeList: Array<{ grade_id: number; name: string }>;
  }>('/shop/plus.lottery/setting');
}

export async function saveLotterySettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.lottery/setting', payload);
}
