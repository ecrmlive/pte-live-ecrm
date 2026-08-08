import { requestClient } from '#/api/request';

export interface ProductMetaPage<T> {
  limit: number;
  list: T[];
  page: number;
  total: number;
}

export interface ProductLabel {
  create_time: string;
  info: string;
  label_id: number;
  name: string;
  sort: number;
  status: 0 | 1;
}

export interface ProductGuarantee {
  content: string;
  create_time: string;
  guarantee_id: number;
  name: string;
  sort: number;
  status: 0 | 1;
}

export interface ProductAttrTemplate {
  create_time: string;
  sort: number;
  template_id: number;
  template_name: string;
  template_value: string;
}

export interface ProductParameterItem {
  name: string;
  required: 0 | 1;
  sort: number;
  values: string[];
}

export interface ProductParameterTemplate {
  cate_id: number;
  create_time: string;
  is_required: 0 | 1;
  params: ProductParameterItem[];
  sort: number;
  template_id: number;
  template_name: string;
}

export type ProductLabelInput = Pick<ProductLabel, 'info' | 'name' | 'sort' | 'status'>;
export type ProductGuaranteeInput = Pick<ProductGuarantee, 'content' | 'name' | 'sort' | 'status'>;
export type ProductAttrTemplateInput = Pick<ProductAttrTemplate, 'sort' | 'template_name' | 'template_value'>;
export type ProductParameterTemplateInput = Pick<
  ProductParameterTemplate,
  'cate_id' | 'is_required' | 'params' | 'sort' | 'template_name'
>;

export interface ProductMetaListParams {
  date_from?: string;
  date_to?: string;
  keyword?: string;
  limit: number;
  page: number;
  status?: 0 | 1;
}

export function listProductLabelsApi(params: ProductMetaListParams) {
  return requestClient.get<ProductMetaPage<ProductLabel>>('/product/labels', { params });
}

export function createProductLabelApi(body: ProductLabelInput) {
  return requestClient.post<ProductLabel>('/product/labels', body);
}

export function updateProductLabelApi(id: number, body: ProductLabelInput) {
  return requestClient.put<ProductLabel>(`/product/labels/${id}`, body);
}

export function deleteProductLabelApi(id: number) {
  return requestClient.delete(`/product/labels/${id}`);
}

export function listProductGuaranteesApi(params: ProductMetaListParams) {
  return requestClient.get<ProductMetaPage<ProductGuarantee>>('/product/guarantees', { params });
}

export function createProductGuaranteeApi(body: ProductGuaranteeInput) {
  return requestClient.post<ProductGuarantee>('/product/guarantees', body);
}

export function updateProductGuaranteeApi(id: number, body: ProductGuaranteeInput) {
  return requestClient.put<ProductGuarantee>(`/product/guarantees/${id}`, body);
}

export function deleteProductGuaranteeApi(id: number) {
  return requestClient.delete(`/product/guarantees/${id}`);
}

export function listProductAttrTemplatesApi(params: ProductMetaListParams) {
  return requestClient.get<ProductMetaPage<ProductAttrTemplate>>('/product/attr-templates', { params });
}

export function createProductAttrTemplateApi(body: ProductAttrTemplateInput) {
  return requestClient.post<ProductAttrTemplate>('/product/attr-templates', body);
}

export function updateProductAttrTemplateApi(id: number, body: ProductAttrTemplateInput) {
  return requestClient.put<ProductAttrTemplate>(`/product/attr-templates/${id}`, body);
}

export function deleteProductAttrTemplateApi(id: number) {
  return requestClient.delete(`/product/attr-templates/${id}`);
}

export function listProductParameterTemplatesApi(params: ProductMetaListParams & { cate_id?: number }) {
  return requestClient.get<ProductMetaPage<ProductParameterTemplate>>('/product/parameter-templates', {
    params,
  });
}

export function getProductParameterTemplateApi(id: number) {
  return requestClient.get<ProductParameterTemplate>(`/product/parameter-templates/${id}`);
}

export function createProductParameterTemplateApi(body: ProductParameterTemplateInput) {
  return requestClient.post<ProductParameterTemplate>('/product/parameter-templates', body);
}

export function updateProductParameterTemplateApi(id: number, body: ProductParameterTemplateInput) {
  return requestClient.put<ProductParameterTemplate>(`/product/parameter-templates/${id}`, body);
}

export function deleteProductParameterTemplateApi(id: number) {
  return requestClient.delete(`/product/parameter-templates/${id}`);
}
