import { requestClient } from '#/api/request';

export interface DataGroupField {
  field: string;
  name: string;
  type: string;
}

export interface DataGroup {
  created_at: string;
  description: string;
  fields: DataGroupField[];
  group_key: string;
  id: number;
  name: string;
  sort: number;
  updated_at: string;
}

export interface DataGroupItem {
  created_at: string;
  data: Record<string, unknown>;
  group_id: number;
  id: number;
  sort: number;
  status: number;
  updated_at: string;
}

export interface DataGroupPage<T> {
  limit: number;
  list: T[];
  page: number;
  total: number;
}

export interface DataGroupInput {
  description?: string;
  fields?: DataGroupField[];
  group_key: string;
  name: string;
  sort?: number;
}

export interface DataGroupItemInput {
  data: Record<string, unknown>;
  sort?: number;
  status?: number;
}

const BASE_PATH = '/maintain/group-data';

export function listDataGroupsApi(params: {
  keyword?: string;
  limit: number;
  page: number;
}) {
  return requestClient.get<DataGroupPage<DataGroup>>(BASE_PATH, { params });
}

export function createDataGroupApi(body: DataGroupInput) {
  return requestClient.post<DataGroup>(BASE_PATH, body);
}

export function updateDataGroupApi(id: number, body: DataGroupInput) {
  return requestClient.put<DataGroup>(`${BASE_PATH}/${id}`, body);
}

export function deleteDataGroupApi(id: number) {
  return requestClient.delete<{ ok: boolean }>(`${BASE_PATH}/${id}`);
}

export function listDataGroupItemsApi(
  groupID: number,
  params: { limit: number; page: number },
) {
  return requestClient.get<DataGroupPage<DataGroupItem> & { group: DataGroup }>(
    `${BASE_PATH}/${groupID}/items`,
    { params },
  );
}

export function createDataGroupItemApi(groupID: number, body: DataGroupItemInput) {
  return requestClient.post<DataGroupItem>(`${BASE_PATH}/${groupID}/items`, body);
}

export function updateDataGroupItemApi(
  groupID: number,
  itemID: number,
  body: DataGroupItemInput,
) {
  return requestClient.put<DataGroupItem>(`${BASE_PATH}/${groupID}/items/${itemID}`, body);
}

export function setDataGroupItemStatusApi(groupID: number, itemID: number, status: 0 | 1) {
  return requestClient.put<{ id: number; status: number }>(
    `${BASE_PATH}/${groupID}/items/${itemID}/status`,
    { status },
  );
}

export function deleteDataGroupItemApi(groupID: number, itemID: number) {
  return requestClient.delete<{ ok: boolean }>(`${BASE_PATH}/${groupID}/items/${itemID}`);
}
