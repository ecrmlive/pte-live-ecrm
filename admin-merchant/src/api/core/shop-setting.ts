import { requestClient } from '#/api/request';

import type { PaginatedList } from './member';

export interface EnumField {
  text: string;
  value: number;
}

export interface ReturnAddressItem {
  address_id: number;
  create_time: string;
  detail: string;
  name: string;
  phone: string;
  sort: number;
}

export interface ReturnAddressListResult {
  list: PaginatedList<ReturnAddressItem> & {
    current_page?: number;
    last_page?: number;
    per_page?: number;
  };
}

export interface ExpressListItem {
  create_time: string;
  express_code: string;
  express_id: number;
  express_name: string;
  sort: number;
  wx_code?: string;
}

export interface DeliveryListItem {
  create_time: string;
  delivery_id: number;
  method: EnumField;
  name: string;
  sort: number;
}

export interface DeliveryAreaNode {
  checked?: boolean;
  children?: DeliveryAreaNode[];
  disabled?: boolean;
  index?: number[] | null;
  indeterminate?: boolean;
  label: string;
  value: number | string;
}

export interface DeliveryRuleRow {
  additional?: number | string;
  additional_fee?: number | string;
  citys: Array<number | string>;
  first?: number | string;
  first_fee?: number | string;
  showData?: Record<
    number | string,
    {
      citys: Array<{ id: number | string; name: string }>;
      id: number | string;
      isAllCitys: boolean;
      name: string;
    }
  >;
}

export interface DeliveryFormPayload {
  delivery_id?: number;
  method?: EnumField | string;
  name: string;
  radio: string;
  rule: DeliveryRuleRow[];
  sort: number | string;
}

export interface DeliveryFormDataBlock {
  additional?: number | string;
  additional_fee?: number | string;
  citys: Array<number | string>;
  first?: number | string;
  first_fee?: number | string;
  province: number[];
}

export interface PrinterListItem {
  create_time: string;
  printer_id: number;
  printer_name: string;
  printer_type: EnumField;
  sort: number;
}

export interface SettingVarsResponse<T = Record<string, unknown>> {
  vars: {
    values: T;
    printerList?: Array<{ printer_id: number | string; printer_name: string }>;
  };
}

export async function getReturnAddressListApi(params: {
  list_rows?: number;
  page?: number;
}) {
  return requestClient.post<ReturnAddressListResult>(
    '/shop/setting.address/index',
    params,
  );
}

export async function addReturnAddressApi(data: {
  detail: string;
  name: string;
  phone: string;
  sort?: number;
}) {
  return requestClient.post('/shop/setting.address/add', data);
}

export async function getReturnAddressDetailApi(addressId: number) {
  return requestClient.get<{ detail: ReturnAddressItem }>(
    '/shop/setting.address/edit',
    { params: { address_id: addressId } },
  );
}

export async function editReturnAddressApi(data: {
  address_id: number;
  detail: string;
  name: string;
  phone: string;
  sort?: number;
}) {
  return requestClient.post('/shop/setting.address/edit', data);
}

export async function deleteReturnAddressApi(addressId: number) {
  return requestClient.post('/shop/setting.address/delete', {
    address_id: addressId,
  });
}

export interface DeliveryTypeOption {
  name: string;
  value: number | string;
}

export function normalizeDeliveryTypeOptions(
  raw: DeliveryTypeOption[] | Record<string, DeliveryTypeOption> | undefined,
): DeliveryTypeOption[] {
  if (Array.isArray(raw)) {
    return raw;
  }
  if (raw && typeof raw === 'object') {
    return Object.values(raw);
  }
  return [];
}

export async function getStoreSettingApi() {
  return requestClient.get<{
    all_type: DeliveryTypeOption[] | Record<string, DeliveryTypeOption>;
    vars: { values: Record<string, unknown> };
  }>('/shop/setting.store/index');
}

/** POST body aligned with api-platform StoreSave / legacy PHP store index. */
export function buildStoreSettingPayload(
  values: Record<string, unknown>,
  wxOpen: boolean,
) {
  const checkedCities = ((values.checkedCities ?? []) as Array<number | string>).map(
    Number,
  );
  return {
    avatarUrl: String(values.avatarUrl ?? ''),
    checkedCities,
    customer: String(values.customer ?? ''),
    isMergeAddress: values.isMergeAddress !== false,
    isMergeBalance: values.isMergeBalance !== false,
    isMergeCoupon: values.isMergeCoupon !== false,
    isMergeOrder: values.isMergeOrder !== false,
    isMergePoints: values.isMergePoints !== false,
    isMergeVerifySms: values.isMergeVerifySms !== false,
    isPassportMobile: !!values.isPassportMobile,
    is_get_log: !!values.is_get_log,
    is_send_wx: !!values.is_send_wx,
    key: String(values.key ?? ''),
    login_desc: String(values.login_desc ?? ''),
    login_logo: String(values.login_logo ?? ''),
    logoUrl: String(values.logoUrl ?? ''),
    mp_open: !!values.mp_open,
    mp_phone: !!values.mp_phone,
    name: String(values.name ?? ''),
    secret: String(values.secret ?? ''),
    sms_open: !!values.sms_open,
    tx_key: String(values.tx_key ?? ''),
    user_name: String(values.user_name ?? ''),
    wx_force_profile: !!values.wx_force_profile,
    wx_open: wxOpen,
    wx_phone: !!values.wx_phone,
    wx_qrcode_version: String(values.wx_qrcode_version ?? 'release'),
  };
}

export async function saveStoreSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/setting.store/index', payload);
}

export async function getTradeSettingApi() {
  return requestClient.get<SettingVarsResponse<Record<string, unknown>>>(
    '/shop/setting.trade/index',
  );
}

export async function saveTradeSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/setting.trade/index', payload);
}

export async function getStorageSettingApi() {
  return requestClient.get<SettingVarsResponse<Record<string, unknown>>>(
    '/shop/setting.storage/index',
  );
}

export async function saveStorageSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/setting.storage/index', payload);
}

export async function getSmsSettingApi() {
  return requestClient.get<SettingVarsResponse<Record<string, unknown>>>(
    '/shop/setting.sms/index',
  );
}

export async function saveSmsSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/setting.sms/index', payload);
}

export async function testSmsSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/setting.sms/smsTest', payload);
}

export async function getPrintingSettingApi() {
  return requestClient.get<
    SettingVarsResponse<Record<string, unknown>> & {
      vars: SettingVarsResponse['vars'] & {
        printerList?: Array<{ printer_id: number | string; printer_name: string }>;
      };
    }
  >('/shop/setting.printing/index');
}

export async function savePrintingSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/setting.printing/index', payload);
}

export type ClearCacheItem = {
  key?: string;
  name: string;
  type?: 'cache' | 'file';
};

export async function getClearSettingApi() {
  return requestClient.get<{ cacheList: Record<string, ClearCacheItem> }>(
    '/shop/setting.clear/index',
  );
}

export async function saveClearSettingApi(payload: { keys: string[] }) {
  return requestClient.post('/shop/setting.clear/index', payload);
}

export async function getMpServiceSettingApi() {
  return requestClient.get<SettingVarsResponse<{ mp_image?: string; qq?: string; wechat?: string }>>(
    '/shop/setting.MpService/index',
  );
}

export async function saveMpServiceSettingApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/setting.MpService/index', payload);
}

export async function getProtocolSettingApi() {
  return requestClient.get<{
    vars: {
      privacy?: { privacy?: string };
      service?: { service?: string };
    };
  }>('/shop/setting.protocol/index');
}

export async function saveProtocolSettingApi(payload: { type: 'privacy' | 'service'; value: string }) {
  return requestClient.post('/shop/setting.protocol/index', payload);
}

export async function getExpressListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<ExpressListItem> }>(
    '/shop/setting.express/index',
    params,
  );
}

export interface ExpressFormPayload {
  express_code: string;
  express_id?: number;
  express_name: string;
  sort: number;
  wx_code?: string;
}

export async function addExpressApi(payload: ExpressFormPayload) {
  return requestClient.post('/shop/setting.express/add', payload);
}

export async function getExpressDetailApi(expressId: number) {
  return requestClient.get<{ detail: ExpressListItem }>('/shop/setting.express/edit', {
    params: { express_id: expressId },
  });
}

export async function editExpressApi(payload: ExpressFormPayload) {
  return requestClient.post('/shop/setting.express/edit', payload);
}

export async function deleteExpressApi(expressId: number) {
  return requestClient.post('/shop/setting.express/delete', { express_id: expressId });
}

export interface ExpressCompanyItem {
  company_code: string;
  company_name: string;
}

export async function getExpressCompanyListApi() {
  return requestClient.post<{ data: ExpressCompanyItem[] }>(
    '/shop/setting.express/companyList',
    {},
  );
}

export async function getDeliveryListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<DeliveryListItem> }>(
    '/shop/setting.delivery/index',
    params,
  );
}

export async function getDeliveryAddMetaApi() {
  return requestClient.get<{ arr: DeliveryAreaNode[]; cityCount: number }>(
    '/shop/setting.delivery/add',
  );
}

export async function addDeliveryApi(payload: DeliveryFormPayload) {
  const cleanedRule = payload.rule.map(({ showData: _showData, ...row }) => row);
  const { rule: _rule, ...rest } = payload;
  return requestClient.post('/shop/setting.delivery/add', {
    ...rest,
    params: JSON.stringify({ rule: cleanedRule }),
  });
}

export async function getDeliveryEditMetaApi(deliveryId: number) {
  return requestClient.get<{
    arr: DeliveryAreaNode[];
    cityCount: number;
    detail: DeliveryFormPayload;
    formData: DeliveryFormDataBlock[];
  }>('/shop/setting.delivery/edit', { params: { delivery_id: deliveryId } });
}

export async function editDeliveryApi(payload: DeliveryFormPayload) {
  const cleanedRule = payload.rule.map(({ showData: _showData, ...row }) => row);
  const { rule: _rule, ...rest } = payload;
  return requestClient.post('/shop/setting.delivery/edit', {
    ...rest,
    params: JSON.stringify({ rule: cleanedRule }),
  });
}

export async function deleteDeliveryApi(deliveryId: number) {
  return requestClient.post('/shop/setting.delivery/delete', { delivery_id: deliveryId });
}

export async function getPrinterListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<PrinterListItem> }>(
    '/shop/setting.printer/index',
    params,
  );
}

export async function getPrinterAddMetaApi() {
  return requestClient.get<{ printerType: Record<string, string> }>(
    '/shop/setting.printer/add',
  );
}

export async function addPrinterApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/setting.printer/add', payload);
}

export async function getPrinterEditMetaApi(printerId: number) {
  return requestClient.get<{ detail: Record<string, unknown>; printerType: Record<string, string> }>(
    '/shop/setting.printer/edit',
    { params: { printer_id: printerId } },
  );
}

export async function editPrinterApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/setting.printer/edit', payload);
}

export async function deletePrinterApi(printerId: number) {
  return requestClient.post('/shop/setting.printer/delete', { printer_id: printerId });
}
