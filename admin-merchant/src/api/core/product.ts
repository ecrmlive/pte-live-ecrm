import { requestClient } from '#/api/request';

export interface PaginatedList<T> {
  current_page?: number;
  data: T[];
  last_page?: number;
  per_page?: number;
  total: number;
}

export interface ProductCategoryOption {
  category_id: number | string;
  child?: ProductCategoryOption[];
  name: string;
}

export interface ProductListItem {
  category_list?: Array<{ category_id: number; name: string }>;
  category_name?: string;
  create_time: string;
  image: Array<{ file_path: string }>;
  product_id: number;
  product_name: string;
  product_price: string | number;
  line_price?: string | number;
  product_sort: number;
  product_status?: number;
  product_stock: number;
  product_type: number;
  product_type_text: string;
  sales_actual: number;
  view_times: number;
}

export function resolveProductTypeText(
  row: Pick<ProductListItem, 'product_type' | 'product_type_text'>,
) {
  if (row.product_type_text) {
    return row.product_type_text;
  }
  switch (Number(row.product_type)) {
    case 2:
      return '虚拟商品';
    case 3:
      return '卡密商品';
    case 4:
      return '次卡商品';
    default:
      return '普通商品';
  }
}

export interface ProductListQuery {
  category_id?: number | string;
  list_rows?: number;
  page?: number;
  product_name?: string;
  product_type?: number;
  type?: string;
}

export interface ProductListResult {
  category: ProductCategoryOption[];
  list: PaginatedList<ProductListItem>;
  productType: Array<{ name: string; value: number }>;
  product_count: Record<string, number>;
}

export interface ProductStatePayload {
  product_id: number;
  state: number;
}

export async function getProductListApi(params: ProductListQuery) {
  return requestClient.post<ProductListResult>('/shop/product.product/index', params);
}

export interface ProductChooseListQuery {
  category_id?: number;
  is_virtual?: number;
  list_rows?: number;
  page?: number;
  product_type?: number;
  search?: string;
  spec_type?: number;
}

export interface ProductChooseItem extends ProductListItem {
  noChoose?: boolean;
  product_image?: string;
}

/** 商品选择弹窗 confirm 统一返回数组 */
export function normalizeProductPickerResult(params: unknown): ProductChooseItem[] {
  if (!params) return [];
  if (Array.isArray(params)) {
    return params as ProductChooseItem[];
  }
  return [params as ProductChooseItem];
}

export async function getProductChooseListApi(params: ProductChooseListQuery) {
  return requestClient.post<{
    category: ProductCategoryOption[];
    list: PaginatedList<ProductChooseItem>;
  }>('/shop/data.product/lists', params);
}

export interface ProductSpecChooseItem {
  noChoose?: boolean;
  product_sku_id: number;
  product_stock: number;
  spec_form: {
    product_price: number | string;
  };
  spec_name: string;
}

export async function chooseProductSpecApi(productId: number) {
  return requestClient.post<{ specList: ProductSpecChooseItem[] }>(
    '/shop/data.product/spec',
    { product_id: productId },
  );
}

export interface ProductIndexChooseQuery extends ProductChooseListQuery {
  is_form?: number;
  is_virtual?: number;
}

export interface ProductIndexChooseItem extends ProductChooseItem {
  category?: { name: string };
  product_stock?: number;
  sku?: Array<{
    product_attr?: string;
    product_price?: number | string;
    product_sku_id?: number;
    spec_sku_id?: number | string;
    spec_type?: number;
    stock_num?: number;
  }>;
  spec_type?: number;
}

export async function getProductIndexListApi(params: ProductIndexChooseQuery) {
  return requestClient.post<{
    category: ProductCategoryOption[];
    list: PaginatedList<ProductIndexChooseItem>;
  }>('/shop/data.product/index', params);
}

export async function deleteProductApi(productId: number) {
  return requestClient.post('/shop/product.product/delete', {
    product_id: productId,
  });
}

export async function changeProductStateApi(payload: ProductStatePayload) {
  return requestClient.post('/shop/product.product/state', payload);
}

export interface ProductAddBaseData {
  basicSetting: Record<string, unknown>;
  category: ProductCategoryOption[];
  delivery: Array<{ delivery_id: number; name: string }>;
  gradeList: Array<{ grade_id: number; name: string; product_equity?: number }>;
  isSpecLocked: boolean;
  logistics: Array<{ name: string; value: number }>;
  productType: Array<{ describe: string; name: string; value: number }>;
  specData: false | unknown;
  tableList: Array<{ name: string; table_id: number }>;
}

export async function getProductAddBaseDataApi() {
  return requestClient.get<ProductAddBaseData>('/shop/product.product/add');
}

export async function addProductApi(params: Record<string, unknown>) {
  return requestClient.post('/shop/product.product/add', {
    params: JSON.stringify(params),
  });
}

export interface ProductEditBaseData extends ProductAddBaseData {
  model: Record<string, unknown>;
}

export async function getProductEditDataApi(
  productId: number,
  scene: 'copy' | 'edit' = 'edit',
) {
  return requestClient.get<ProductEditBaseData>('/shop/product.product/edit', {
    params: { product_id: productId, scene },
  });
}

export async function editProductApi(payload: {
  product_id: number;
  scene: string;
  params: string;
}) {
  return requestClient.post<{ msg?: string }>(
    '/shop/product.product/edit',
    payload,
  );
}

export async function addProductSpecApi(payload: {
  spec_name: string;
  spec_value: string;
}) {
  return requestClient.post<{ spec_id: number; spec_value_id: number }>(
    '/shop/product.spec/addSpec',
    payload,
  );
}

export async function addProductSpecValueApi(payload: {
  spec_id: number;
  spec_value: string;
}) {
  return requestClient.post<{ spec_value_id: number }>(
    '/shop/product.spec/addSpecValue',
    payload,
  );
}

export async function importProductVirtualApi(file: File) {
  return requestClient.upload<{
    list: Array<{ card_no: string; card_pwd: string }>;
  }>('/shop/product.product/importVirtual', { file, iFile: file });
}

export interface ProductCategoryItem {
  category_id: number;
  child?: ProductCategoryItem[];
  create_time?: string;
  image_id?: number | string;
  images?: { file_path: string };
  name: string;
  parent_id?: number | string;
  sort: number;
  status?: number;
}

export interface ProductCategoryListResult {
  list: ProductCategoryItem[];
}

export interface ProductCategoryFormPayload {
  category_id?: number;
  image_id: number | string;
  name: string;
  parent_id: number | string;
  sort: number;
}

export async function getProductCategoryListApi() {
  return requestClient.post<ProductCategoryListResult>(
    '/shop/product.category/index',
    {},
  );
}

export async function addProductCategoryApi(payload: ProductCategoryFormPayload) {
  return requestClient.post('/shop/product.category/add', payload);
}

export async function editProductCategoryApi(payload: ProductCategoryFormPayload) {
  return requestClient.post('/shop/product.category/edit', payload);
}

export async function deleteProductCategoryApi(categoryId: number) {
  return requestClient.post('/shop/product.category/delete', {
    category_id: categoryId,
  });
}

export async function setProductCategoryStatusApi(
  categoryId: number,
  status: number,
) {
  return requestClient.post<{ msg?: string }>('/shop/product.category/set', {
    category_id: categoryId,
    status,
  });
}

export interface ProductCommentItem {
  comment_id: number;
  content: string;
  create_time: string;
  is_picture: number;
  order_id?: number;
  orderM?: { order_no: string };
  product: {
    image: Array<{ file_path: string }>;
    product_name: string;
    product_price: string | number;
  };
  score: number;
  sort: number;
  status: number;
  user?: { nickName: string };
}

export interface ProductCommentListQuery {
  list_rows?: number;
  name?: string;
  page?: number;
  score?: number;
  status?: number | string;
}

export interface ProductCommentListResult {
  list: PaginatedList<ProductCommentItem>;
  num: number;
}

export async function getProductCommentListApi(params: ProductCommentListQuery) {
  return requestClient.post<ProductCommentListResult>(
    '/shop/product.comment/index',
    params,
  );
}

export async function deleteProductCommentApi(commentId: number) {
  return requestClient.post('/shop/product.comment/delete', {
    comment_id: commentId,
  });
}

export interface ProductCommentDetail {
  comment_id: number;
  content: string;
  create_time: string;
  image?: Array<{ file_path: string }>;
  product: {
    image: Array<{ file_path: string }>;
    product_name: string;
  };
  score: number;
  sort: number;
  status: number;
  user: { nickName: string };
}

export async function getProductCommentDetailApi(commentId: number) {
  return requestClient.post<{ data: ProductCommentDetail }>(
    '/shop/product.comment/detail',
    { comment_id: commentId },
  );
}

export async function editProductCommentApi(payload: {
  comment_id: number;
  content?: string;
  score?: number;
  sort: number;
  status: number;
}) {
  return requestClient.post<{ msg?: string }>('/shop/product.comment/edit', payload);
}

export interface ProductVirtualDetail {
  product_name: string;
  product_sku_id: number;
  product_stock: number;
  spec_type: number;
  specList: Array<{
    product_sku_id: number;
    spec_form: {
      card_info?: string;
      card_type?: number;
      stock_num?: number;
      virtualInfo?: unknown;
    };
    spec_name: string;
    spec_sku_id?: string | number;
  }>;
}

export interface ProductVirtualItem {
  card_no: string;
  card_pwd: string;
  card_type?: number;
  order_id?: number;
  product_id: number;
  product_sku_id: number;
  spec_name?: string;
  spec_sku_id?: string | number;
  use_status: number;
  use_time: string;
  virtual_id: number;
}

export interface ProductVirtualListResult {
  detail: ProductVirtualDetail;
  list: PaginatedList<ProductVirtualItem>;
}

export async function getProductVirtualListApi(params: {
  list_rows?: number;
  page?: number;
  product_id: number;
  use_status?: number;
}) {
  return requestClient.get<ProductVirtualListResult>(
    '/shop/product.virtual/index',
    { params },
  );
}

export async function addProductVirtualApi(payload: {
  product_id: number;
  product_sku_id: number;
  spec_sku_id?: number | string;
  virtualInfo: Array<{ card_no: string; card_pwd: string }>;
}) {
  return requestClient.post('/shop/product.virtual/add', payload);
}

export async function editProductVirtualCardApi(payload: {
  card_no: string;
  card_pwd: string;
  card_type?: number;
  virtual_id: number;
}) {
  return requestClient.post('/shop/product.virtual/edit', payload);
}

export async function editProductVirtualFixedApi(payload: {
  card_info: string;
  card_type?: number;
  product_sku_id: number;
  stock_num: number | string;
}) {
  return requestClient.post('/shop/product.virtual/edit', payload);
}

export async function deleteProductVirtualApi(virtualId: number | string) {
  return requestClient.post('/shop/product.virtual/delete', {
    virtual_id: virtualId,
  });
}
