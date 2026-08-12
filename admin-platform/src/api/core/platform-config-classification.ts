import { requestClient } from '#/api/request';

export interface ConfigClassification {
  classify_key: string;
  created_at: string;
  description: string;
  icon: string;
  id: number;
  name: string;
  parent_id: number;
  sort: number;
  status: number;
  updated_at: string;
}

export interface ConfigClassificationItem {
  backend_type: number;
  classification_id: number;
  config_key: string;
  content: string;
  created_at: string;
  description: string;
  field_type: ConfigFieldType;
  id: number;
  name: string;
  sort: number;
  status: number;
  updated_at: string;
}

export type ConfigFieldType =
  | 'file'
  | 'image'
  | 'input'
  | 'number'
  | 'radio'
  | 'switch'
  | 'textarea';

export interface ConfigClassificationPage<T> {
  limit: number;
  list: T[];
  page: number;
  total: number;
}

export interface ConfigClassificationInput {
  classify_key: string;
  description?: string;
  icon?: string;
  name: string;
  sort?: number;
  status?: number;
}

export interface ConfigClassificationItemInput {
  backend_type?: number;
  classification_id?: number;
  config_key: string;
  content: string;
  description?: string;
  field_type?: ConfigFieldType;
  name: string;
  sort?: number;
  status?: number;
}

const BASE_PATH = '/maintain/config-classifications';

export function listConfigClassificationsApi(params: {
  limit: number;
  name?: string;
  page: number;
  status?: number;
}) {
  return requestClient.get<ConfigClassificationPage<ConfigClassification>>(BASE_PATH, { params });
}

export function createConfigClassificationApi(body: ConfigClassificationInput) {
  return requestClient.post<ConfigClassification>(BASE_PATH, body);
}

export function updateConfigClassificationApi(id: number, body: ConfigClassificationInput) {
  return requestClient.put<ConfigClassification>(`${BASE_PATH}/${id}`, body);
}

export function setConfigClassificationStatusApi(id: number, status: 0 | 1) {
  return requestClient.put<{ id: number; status: number }>(`${BASE_PATH}/${id}/status`, { status });
}

export function deleteConfigClassificationApi(id: number) {
  return requestClient.delete<{ ok: boolean }>(`${BASE_PATH}/${id}`);
}

export function listConfigClassificationItemsApi(
  classificationID: number,
  params: { limit: number; page: number },
) {
  return requestClient.get<
    ConfigClassificationPage<ConfigClassificationItem> & { classification: ConfigClassification }
  >(`${BASE_PATH}/${classificationID}/items`, { params });
}

export function createConfigClassificationItemApi(
  classificationID: number,
  body: ConfigClassificationItemInput,
) {
  return requestClient.post<ConfigClassificationItem>(`${BASE_PATH}/${classificationID}/items`, body);
}

export function updateConfigClassificationItemApi(
  classificationID: number,
  itemID: number,
  body: ConfigClassificationItemInput,
) {
  return requestClient.put<ConfigClassificationItem>(
    `${BASE_PATH}/${classificationID}/items/${itemID}`,
    body,
  );
}

export function setConfigClassificationItemStatusApi(
  classificationID: number,
  itemID: number,
  status: 0 | 1,
) {
  return requestClient.put<{ id: number; status: number }>(
    `${BASE_PATH}/${classificationID}/items/${itemID}/status`,
    { status },
  );
}

export function deleteConfigClassificationItemApi(classificationID: number, itemID: number) {
  return requestClient.delete<{ ok: boolean }>(`${BASE_PATH}/${classificationID}/items/${itemID}`);
}
