import { requestClient } from '#/api/request';

/** CRMEB 订单统计时间：lately7 / lately30 / month / year / `YYYY/MM/DD-YYYY/MM/DD` */
export type OrderStatDate = string;

export interface OrderTopCard {
  count: number;
  mom: number;
  statistic: number;
  title: string;
}

export interface OrderLinePoint {
  order_num: number;
  pay_price: number;
  refund_num: number;
  refund_price: number;
  xaxis: string;
}

export interface OrderPieSlice {
  name: string;
  value: number;
}

export function getOrderStatTopApi(date: OrderStatDate) {
  return requestClient.get<OrderTopCard[]>('/analytics/order/top', {
    params: { date },
  });
}

export function getOrderStatLineApi(date: OrderStatDate) {
  return requestClient.get<OrderLinePoint[]>('/analytics/order/line_chart', {
    params: { date },
  });
}

/** type=0 订单类型；type=1 发货方式（对齐 CRMEB pie_chart/:type） */
export function getOrderStatPieApi(type: 0 | 1, date: OrderStatDate) {
  return requestClient.get<OrderPieSlice[]>(`/analytics/order/pie_chart/${type}`, {
    params: { date },
  });
}
