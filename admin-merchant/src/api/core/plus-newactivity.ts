import { requestClient } from '#/api/request';

export interface NewActivitySkuItem {
  activity_price: number | string;
  activity_product_sku_id?: number;
  product_attr?: string;
  product_price: number | string;
  product_sku_id: number;
  showSkuId?: string;
  stock_num?: number;
}

export interface NewActivityProductItem {
  activity_product_id?: number;
  activity_price?: number | string;
  product_attr: string;
  product_id: number;
  product_name?: string;
  product_num?: number;
  product_price: number | string;
  product_sku_id: number;
  showSkuId: string | number;
  sort: number;
  spec_sku_id?: number;
  spec_type?: number;
  stock_num: number;
  sku?: NewActivitySkuItem[];
}

export interface NewActivityFormValues {
  buy_money: number | string;
  content: string;
  delivery_id: number | string;
  image: string;
  is_open: boolean | number | string;
  limit_num: number | string;
  notice: string;
  product_list?: NewActivityProductItem[];
}

export interface DeliveryOption {
  delivery_id: number;
  name: string;
}

export async function getNewActivitySettingApi() {
  return requestClient.get<{
    delivery: DeliveryOption[];
    vars: { values: NewActivityFormValues };
  }>('/shop/plus.newactivity/index');
}

export async function saveNewActivitySettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.newactivity/index', payload);
}
