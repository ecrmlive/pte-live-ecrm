import { requestClient } from '#/api/request';

export interface PaginatedList<T> {
  data: T[];
  total: number;
}

export interface OrderListQuery {
  create_time?: '' | string[];
  dataType?: string;
  list_rows?: number;
  order_no?: string;
  order_source?: number | string;
  page?: number;
  search?: string;
  store_id?: number | string;
  style_id?: number | string;
}

export interface OrderListItem {
  advance?: { pay_price: string | number };
  assemble_status?: number;
  balance?: number;
  cancel_status?: number;
  create_time: string;
  delivery_status: { value: number };
  delivery_type: { text: string; value: number };
  express_id?: number;
  express_no?: string;
  express_price: string | number;
  is_comment: number;
  is_delete?: number;
  is_label?: number;
  is_single?: number;
  online_money?: number;
  order_id: number;
  order_no: string;
  order_source: number;
  order_source_text: string;
  order_status: { value: number };
  order_status_value?: number;
  pay_price: string | number;
  pay_source?: string;
  pay_status: { value: number };
  pay_type: { text: string; value: number };
  product: Array<{
    image: { file_path: string };
    is_gift?: number;
    product_attr?: string;
    product_name: string;
    product_price: string | number;
    refund?: { status: { text: string }; type: { text: string } };
    total_num: number;
  }>;
  receipt_status: { value: number };
  state_text: string;
  user?: { nickName: string; user_id: number };
  verify_status?: number;
  virtual_auto?: number;
  wx_delivery_status?: number;
}

export interface OrderTableRow extends Partial<OrderListItem> {
  is_top_row?: boolean;
}

export interface OrderListResult {
  ex_style: Array<{ name: string; value: number | string }>;
  is_send_wx: boolean | string;
  list: PaginatedList<OrderListItem>;
  order_count: { order_count: Record<string, number> };
  shop_list: Array<{ store_id: number | string; store_name: string }>;
  sourceList: Array<{ name: string; value: number | string }>;
}

export async function getOrderListApi(params: OrderListQuery) {
  return requestClient.post<OrderListResult>('/shop/order.order/index', params);
}

export async function getOrderExpressApi(payload: {
  express_id?: number;
  express_no?: string;
  order_id: number;
}) {
  return requestClient.post<{ express: { list: unknown } }>(
    '/shop/order.order/express',
    payload,
  );
}

export async function wxDeliveryOrderApi(orderId: number) {
  return requestClient.post('/shop/order.order/wxDelivery', { order_id: orderId });
}

export async function batchDeliveryOrderApi(file: File) {
  return requestClient.upload<{ msg?: string }>(
    '/shop/order.operate/batchDelivery',
    { iFile: file } as Record<string, File>,
  );
}

function saveExportBlob(blob: Blob) {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `orders-${Date.now()}.xlsx`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}

export async function exportOrderApi(params: OrderListQuery) {
  const blob = await requestClient.download<Blob>('/shop/order.operate/export', {
    params,
  });
  saveExportBlob(blob);
}

export async function cancelOrderApi(payload: {
  cancel_remark: string;
  order_no: string;
}) {
  return requestClient.upload<{ msg?: string }>(
    '/shop/order.order/orderCancel',
    payload as Record<string, string>,
  );
}

export interface OrderEnumField {
  text: string;
  value: number;
}

export interface OrderDetailProduct {
  delivery_num?: number;
  grade_product_price?: number | string;
  image?: { file_path?: string };
  is_gift?: number;
  is_new?: number;
  is_user_grade?: number;
  new_price?: number | string;
  order_product_id?: number;
  product_attr?: string;
  product_name: string;
  product_no?: string;
  product_price: number | string;
  product_type?: number;
  product_weight?: number | string;
  refund?: { status: { text: string }; type: { text: string } };
  tableData?: Array<{ name?: string; type?: string; value?: string }>;
  table_record_id?: number;
  total_num: number;
  total_pay_price?: number | string;
  total_price?: number | string;
  use_verify_num?: number;
  verify_num?: number;
}

export interface OrderDeliveryPackageProduct {
  delivery_num?: number;
  image_path?: string;
  order_product_id?: number;
  product_attr?: string;
  product_name?: string;
  total_num?: number;
}

export interface OrderDeliveryRecord {
  create_time?: string;
  express_id?: number;
  express_name?: string;
  express_no?: string;
  is_label?: number;
  label?: string;
  label_print_type?: number;
  order_delivery_id?: number;
  order_id?: number;
  product_list?: OrderDeliveryPackageProduct[];
  remark?: string;
}

export interface OrderLabelSetting {
  setting_id: number;
  setting_name: string;
}

export interface OrderLabelTemplate {
  template_id: number;
  template_name: string;
}

export interface OrderDeliveryPackageForm {
  address_id?: number;
  delivery_data: OrderDeliveryPackageProduct[];
  express_id?: number;
  express_no?: string;
  remark?: string;
  setting_id?: number;
  template_id?: number;
}

export interface OrderDeliveryFormPayload {
  address_id?: number;
  delivery_list?: OrderDeliveryPackageForm[];
  express_id?: number;
  express_no?: string;
  is_label?: number;
  is_single?: number;
  order_id: number;
  setting_id?: number;
  template_id?: number;
}

export interface OrderDetailData {
  advance?: {
    money_return?: number;
    pay_price?: number | string;
    reduce_money?: number | string;
  };
  address?: {
    city?: string;
    detail?: string;
    district?: string;
    name?: string;
    phone?: string;
    province?: string;
    region?: {
      city?: string;
      province?: string;
      region?: string;
    };
  };
  assemble_status?: number;
  balance?: number;
  buyer_remark?: string;
  cancel_status?: number;
  card_type?: number;
  cancel_remark?: string;
  coupon_money?: number;
  create_time?: string;
  custom_form?: Array<{ label?: string; title?: string; value?: unknown }>;
  delivery_status?: OrderEnumField;
  delivery_time?: string;
  delivery_type?: OrderEnumField;
  expire_text?: string;
  express?: { express_id?: number; express_name?: string };
  express_no?: string;
  express_price?: number | string;
  extract?: { linkman?: string; phone?: string };
  extractStore?: {
    address?: string;
    linkman?: string;
    phone?: string;
    region?: { city?: string; province?: string; region?: string };
    store_id?: number;
    store_name?: string;
  };
  fullreduce_money?: number;
  is_comment?: number;
  is_label?: number;
  is_single?: number;
  label?: string;
  label_print_type?: number;
  online_money?: number;
  orderDeliverList?: OrderDeliveryRecord[];
  order_id: number;
  order_no: string;
  order_price?: number | string;
  order_source?: number;
  order_status?: OrderEnumField;
  pay_price?: number | string;
  pay_status?: OrderEnumField;
  total_price?: number | string;
  pay_time?: string;
  pay_type?: OrderEnumField;
  points_money?: number;
  product?: OrderDetailProduct[];
  receipt_status?: OrderEnumField;
  receipt_time?: string;
  setting_id?: number;
  surplus_num?: number;
  template_id?: number;
  transaction_id?: string;
  user?: { mobile?: string; nickName?: string; user_id?: number };
  verify_status?: number;
  pay_source?: string;
  virtual_auto?: number;
  virtual_content?: unknown;
  wx_delivery_status?: number;
  extractClerk?: { real_name?: string };
}

export interface OrderDetailResult {
  address_list?: Array<{
    address_id: number;
    detail: string;
    name: string;
    phone: string;
  }>;
  detail: OrderDetailData;
  expressList?: Array<{
    express_code?: string;
    express_id: number;
    express_name: string;
  }>;
  extractList?: Array<{
    clerk?: { real_name?: string };
    create_time?: string;
    store?: { store_name?: string };
    verify_num?: number;
    verify_status?: number;
  }>;
  label_list?: OrderLabelSetting[];
  shopClerkList?: Array<{ clerk_id: number; real_name: string }>;
  template_list?: OrderLabelTemplate[];
}

export async function getOrderDetailApi(orderId: number) {
  return requestClient.post<OrderDetailResult>('/shop/order.order/detail', {
    order_id: orderId,
  });
}

export async function deliveryOrderApi(payload: OrderDeliveryFormPayload | Record<string, unknown>) {
  return requestClient.post('/shop/order.order/delivery', payload);
}

export async function labelCancelOrderApi(payload: {
  is_multi: number;
  order_id: number;
}) {
  return requestClient.post('/shop/order.order/labelCancel', payload);
}

export async function printRepeateOrderApi(payload: {
  is_multi: number;
  order_id: number;
}) {
  return requestClient.post('/shop/order.order/printRepeate', payload);
}

export async function updateOrderExpressApi(payload: {
  express_id: number;
  express_no: string;
  order_delivery_id?: number;
  order_id: number;
}) {
  return requestClient.post('/shop/order.order/updateExpress', payload);
}

export async function updateOrderPriceApi(payload: {
  order: { update_express_price: number | string; update_price: number | string };
  order_id: number;
}) {
  return requestClient.post('/shop/order.order/updatePrice', payload);
}

export async function updateOrderAddressApi(payload: {
  city_id: number;
  detail: string;
  name: string;
  order_id: number;
  phone: string;
  province_id: number;
  region_id: number;
}) {
  return requestClient.post('/shop/order.order/updateAddress', payload);
}

export async function virtualDeliveryOrderApi(payload: {
  order_id: number;
  virtual_content: string;
}) {
  return requestClient.post('/shop/order.order/virtual', payload);
}

export async function confirmCancelOrderApi(payload: {
  is_cancel: number;
  order_id: number;
}) {
  return requestClient.post('/shop/order.operate/confirmCancel', payload);
}

export async function extractOrderApi(payload: {
  extract_form: {
    order: {
      extract_clerk_id?: number;
      extract_status?: number;
      verify_num?: number | string;
    };
    order_id: number;
  };
}) {
  return requestClient.post('/shop/order.operate/extract', payload);
}

export async function getRegionTreeApi() {
  return requestClient.post<{ regionData: Record<string, unknown> }>(
    '/shop/data.region/lists',
    {},
  );
}

export interface RefundListQuery {
  create_time?: string[];
  list_rows?: number;
  order_no?: string;
  page?: number;
  status?: number;
  type?: number | string;
}

export interface RefundListItem {
  create_time: string;
  is_agree: OrderEnumField;
  is_receipt: number;
  is_user_send: number;
  order_id: number;
  order_no: string;
  order_refund_id: number;
  orderMaster: { create_time: string; order_id: number; order_source?: number };
  orderproduct: {
    image?: { file_path?: string };
    product_name: string;
    product_price: number | string;
    total_num: number;
    total_pay_price: number | string;
  };
  refund_money: string;
  status: OrderEnumField;
  type: OrderEnumField;
  user?: { nickName?: string; user_id?: number };
}

export interface RefundListResult {
  arr: Record<number, { status: OrderEnumField; total: number }>;
  list: PaginatedList<RefundListItem>;
}

export async function getRefundListApi(params: RefundListQuery) {
  return requestClient.post<RefundListResult>('/shop/order.refund/index', params);
}

export interface RefundDetailItem {
  address?: { detail?: string; name?: string; phone?: string };
  apply_desc?: string;
  deliver_time?: string;
  express?: { express_name?: string };
  express_no?: string;
  image?: Array<{ file_path: string }>;
  is_agree: OrderEnumField;
  is_plate_send?: number;
  is_receipt: number;
  is_user_send: number;
  order_id: number;
  order_master: { order_id: number; order_no: string };
  orderproduct: {
    image?: { file_path?: string };
    line_price?: number | string;
    max_refund_money?: number | string;
    product_name: string;
    product_weight?: number | string;
    total_num: number;
    total_pay_price: number | string;
  };
  order_refund_id: number;
  refuse_desc?: string;
  refund_money?: number | string;
  send_express_no?: string;
  send_time?: string;
  sendexpress?: { express_name?: string };
  status: OrderEnumField;
  type: OrderEnumField;
  user?: { nickName?: string; user_id?: number };
}

export interface RefundDetailResult {
  address: Array<{ address_id: number; detail: string }>;
  detail: RefundDetailItem;
  expressList: Array<{ express_id: number; express_name: string }>;
  order?: {
    advance?: { money_return?: number; pay_price?: number | string };
    order_source?: number;
  };
}

export async function getRefundDetailApi(orderRefundId: number) {
  return requestClient.post<RefundDetailResult>('/shop/order.refund/detail', {
    order_refund_id: orderRefundId,
  });
}

export async function auditRefundApi(payload: {
  address_id?: number | string;
  is_agree: number;
  order_refund_id: number;
  refuse_desc?: string;
  refund_money?: number | string;
}) {
  return requestClient.post('/shop/order.refund/audit', payload);
}

export async function receiptRefundApi(payload: {
  order_refund_id: number;
  refund_money?: number | string;
  send_express_id?: number | string;
  send_express_no?: string;
}) {
  return requestClient.post('/shop/order.refund/receipt', payload);
}

/** 将接口订单列表转为带表头行的表格数据 */
export function flattenOrderTableRows(list: OrderListItem[]): OrderTableRow[] {
  const rows: OrderTableRow[] = [];
  for (const item of list) {
    rows.push({
      create_time: item.create_time,
      is_top_row: true,
      order_no: item.order_no,
      order_source: item.order_source,
      order_source_text: item.order_source_text,
      order_status_value: item.order_status.value,
    });
    rows.push(item);
  }
  return rows;
}

export function orderSpanMethod({
  columnIndex,
  rowIndex,
}: {
  columnIndex: number;
  rowIndex: number;
}) {
  if (rowIndex % 2 === 0 && columnIndex === 0) {
    return [1, 8] as [number, number];
  }
  return [1, 1] as [number, number];
}
