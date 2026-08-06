import { requestClient } from '#/api/request';

/** CRMEB 商品统计时间：lately7 / lately30 / month / year / `YYYY/MM/DD-YYYY/MM/DD` */
export type ProductStatDate = string;

export interface ProductTopCard {
  count: number;
  mom: number;
  statistic: number;
  title: string;
}

export interface ProductLinePoint {
  paid_num: number;
  relation: number;
  total_num: number;
  visit: number;
  xaxis: string;
}

export interface ProductPieSlice {
  name: string;
  value: number;
}

export function getProductStatTopApi(date: ProductStatDate) {
  return requestClient.get<ProductTopCard[]>('/analytics/product/top', {
    params: { date },
  });
}

export function getProductStatLineApi(date: ProductStatDate) {
  return requestClient.get<ProductLinePoint[]>('/analytics/product/line_chart', {
    params: { date },
  });
}

/** type=1 商品分类；type=0 商品类型（对齐 CRMEB pie_chart/:type） */
export function getProductStatPieApi(type: 0 | 1) {
  return requestClient.get<ProductPieSlice[]>(`/analytics/product/pie_chart/${type}`);
}
