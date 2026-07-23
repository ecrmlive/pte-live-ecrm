import { requestClient } from '#/api/request';

export interface DayCompareValue {
  today: number | string;
  yesterday: number | string;
}

export interface SalesStatisticsIndex {
  order: {
    order_per_price: DayCompareValue;
    order_refund_money: DayCompareValue;
    order_refund_total: DayCompareValue;
    order_total: DayCompareValue;
    order_total_price: DayCompareValue;
    order_user_total: DayCompareValue;
  };
  product: {
    no_pay: DayCompareValue;
    pay: DayCompareValue;
    sale: DayCompareValue;
  };
  productRefundRanking: ProductRankItem[];
  productSaleRanking: ProductRankItem[];
  productViewRanking: ProductRankItem[];
}

export interface ProductRankItem {
  image?: Array<{ file_path?: string }> | { file?: { file_path?: string } };
  product_name: string;
  refund_count?: number;
  total_sales_num?: number;
  view_times?: number;
}

export interface UserStatisticsIndex {
  inviteRanking: UserRankItem[];
  payRanking: UserRankItem[];
  pointsRanking: UserRankItem[];
  user: {
    user_add: DayCompareValue;
    user_pay: DayCompareValue;
    user_total: DayCompareValue;
  };
}

export interface UserRankItem {
  avatarUrl?: string;
  expend_money?: number | string;
  nickName: string;
  total_invite?: number;
  total_points?: number | string;
  user_id?: number;
}

export interface StatisticsSeriesData {
  data: Array<{ day?: string; total_money?: number | string; total_num?: number | string }>;
  days: string[];
}

export async function getSalesStatisticsIndexApi() {
  return requestClient.post<SalesStatisticsIndex>('/shop/statistics.sales/index', {});
}

export async function getSalesOrderByDateApi(params: {
  search_time: string[];
  type: 'order' | 'refund';
}) {
  return requestClient.post<StatisticsSeriesData>(
    '/shop/statistics.sales/order',
    params,
  );
}

export async function getSalesProductByDateApi(params: { search_time: string[] }) {
  return requestClient.post<StatisticsSeriesData>(
    '/shop/statistics.sales/product',
    params,
  );
}

export async function getUserStatisticsIndexApi() {
  return requestClient.post<UserStatisticsIndex>('/shop/statistics.user/index', {});
}

export async function getUserNewByDateApi(params: { search_time: string[] }) {
  return requestClient.post<StatisticsSeriesData>(
    '/shop/statistics.user/new_user',
    params,
  );
}

export async function getUserPayByDateApi(params: { search_time: string[] }) {
  return requestClient.post<StatisticsSeriesData>(
    '/shop/statistics.user/pay_user',
    params,
  );
}

export function productRankImage(item: ProductRankItem): string {
  const img = item.image;
  if (Array.isArray(img)) {
    return img[0]?.file_path ?? '';
  }
  if (img && typeof img === 'object' && 'file' in img) {
    return (img as { file?: { file_path?: string } }).file?.file_path ?? '';
  }
  return '';
}

export function formatDateYmd(date: Date): string {
  const y = date.getFullYear();
  const m = String(date.getMonth() + 1).padStart(2, '0');
  const d = String(date.getDate()).padStart(2, '0');
  return `${y}-${m}-${d}`;
}

export function defaultStatisticsDateRange(days = 7): [string, string] {
  const end = new Date();
  const start = new Date();
  start.setDate(start.getDate() - days);
  return [formatDateYmd(start), formatDateYmd(end)];
}

export const STATISTICS_DATE_SHORTCUTS = [
  {
    text: '最近一周',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setDate(start.getDate() - 7);
      return [start, end];
    },
  },
  {
    text: '最近一个月',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setDate(start.getDate() - 30);
      return [start, end];
    },
  },
  {
    text: '最近三个月',
    value: () => {
      const end = new Date();
      const start = new Date();
      start.setDate(start.getDate() - 90);
      return [start, end];
    },
  },
];
