import { requestClient } from '#/api/request';

/** CRMEB 用户统计时间：lately7 / lately30 / month / year / `YYYY/MM/DD-YYYY/MM/DD` */
export type UserStatDate = string;

export interface UserTopCard {
  count: number;
  mom: number;
  statistic: number;
  title: string;
}

export interface UserLinePoint {
  count: number;
  xaxis: string;
}

export interface UserDealPoint {
  new: number;
  old: number;
  xaxis: string;
}

export function getUserStatTopApi(date: UserStatDate) {
  return requestClient.get<UserTopCard[]>('/analytics/user/top', {
    params: { date },
  });
}

/** type=0 新增用户；1 活跃用户；2 新增付费会员 */
export function getUserStatLineApi(date: UserStatDate, type: 0 | 1 | 2) {
  return requestClient.get<UserLinePoint[]>('/analytics/user/line_chart', {
    params: { date, type },
  });
}

/** 成交用户：老用户 / 新用户分组柱图 */
export function getUserStatDealApi(date: UserStatDate) {
  return requestClient.get<UserDealPoint[]>('/analytics/user/pie_chart', {
    params: { date },
  });
}
