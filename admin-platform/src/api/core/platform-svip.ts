import { requestClient } from '#/api/request';

export interface PlatformSvipUser {
  integral: number;
  is_svip: number;
  is_svip_active: boolean;
  nickname: string;
  now_money: number;
  phone_masked?: string;
  svip_endtime?: string | null;
  uid: number;
}

export interface PlatformSvipUserPage {
  limit: number;
  list: PlatformSvipUser[];
  page: number;
  total: number;
}

export function listPlatformSvipUsersApi(params: { limit: number; page: number }) {
  return requestClient.get<PlatformSvipUserPage>('/users', { params });
}

export function setPlatformUserSvipApi(
  id: number,
  body: { is_svip: number; svip_endtime?: string },
) {
  return requestClient.put(`/users/${id}/svip`, body);
}
