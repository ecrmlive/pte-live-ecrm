import { requestClient } from '#/api/request';

export interface PlatformUserGroup {
  create_time: string;
  group_id: number;
  group_name: string;
  sort: number;
}

export interface PlatformUserGroupPage {
  limit: number;
  list: PlatformUserGroup[];
  page: number;
  total: number;
}

export function listPlatformUserGroupsApi(params: { limit: number; page: number }) {
  return requestClient.get<PlatformUserGroupPage>('/user/groups', { params });
}

export function createPlatformUserGroupApi(body: { group_name: string; sort: number }) {
  return requestClient.post<PlatformUserGroup>('/user/groups', body);
}

export function updatePlatformUserGroupApi(id: number, body: { group_name: string; sort: number }) {
  return requestClient.put<PlatformUserGroup>(`/user/groups/${id}`, body);
}

export function deletePlatformUserGroupApi(id: number) {
  return requestClient.delete(`/user/groups/${id}`);
}
