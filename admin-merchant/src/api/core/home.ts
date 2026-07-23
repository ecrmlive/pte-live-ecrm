import { requestClient } from '#/api/request';

export interface HomeQueryParams {
  product_type?: number;
  product_time?: number;
  sale_time?: number;
  top_time?: number;
  user_time?: number;
}

export interface HomeTopData {
  comment_rate: number;
  comment_today: number;
  comment_total: string;
  comment_yesterday: number;
  order_rate: number;
  order_today: number;
  order_total: string;
  order_yesterday: number;
  product_rate: number;
  product_today: number;
  product_total: string;
  product_yesterday: number;
  sale_rate: number;
  sale_today: string;
  sale_total: string;
  sale_yesterday: string;
  user_rate: number;
  user_today: number;
  user_total: string;
  user_yesterday: number;
}

export interface HomeWaitData {
  agent: {
    apply: number;
    cash_apply: number;
    cash_money: number;
  };
  order: {
    card_count: number;
    disposal: number;
    refund: number;
  };
  review: {
    balance_apply: number;
    balance_money: number;
    comment: number;
  };
  stock: {
    product: number;
  };
}

export interface HomeChartSeries {
  data: Array<Record<string, number | string>>;
  days: string[];
  saleMoney?: string;
}

export interface HomeProductRankItem {
  product_name: string;
  total_num: number;
  total_price: number;
}

export interface HomeDashboardPayload {
  productRank: HomeProductRankItem[];
  saleData: HomeChartSeries;
  top_data: HomeTopData;
  update_time: string;
  userData: HomeChartSeries;
  wait_data: HomeWaitData;
}

interface HomeDashboardResponse {
  data: HomeDashboardPayload;
}

/** 商户首页统计（原 legacy POST /shop/Index/index） */
export async function getHomeDashboardApi(params: HomeQueryParams) {
  const result = await requestClient.post<HomeDashboardResponse>(
    '/shop/Index/index',
    params,
  );
  return result.data;
}
