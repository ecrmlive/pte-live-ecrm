import { apiLiveAdminPost } from '#/utils/api-live-admin';
import type { AccessNode } from '#/views/access/types';

export interface PlatformRoleRow {
  create_time?: number;
  remark?: string;
  role_id: number;
  role_name: string;
  sort: number;
  update_time?: number;
}

export interface PlatformRoleForm {
  access_id: number[];
  remark?: string;
  role_id?: number;
  role_name: string;
  sort: number;
}

const PlatformRoleApi = {
  roleList(errorback?: boolean) {
    return apiLiveAdminPost<{ list: PlatformRoleRow[] }>(
      '/admin/platform/role/index',
      {},
      errorback,
    );
  },
  roleAddInfo(errorback?: boolean) {
    return apiLiveAdminPost<{ menu: AccessNode[] }>(
      '/admin/platform/role/add-info',
      {},
      errorback,
    );
  },
  roleAdd(params: string, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/role/add', { params }, errorback);
  },
  roleEditInfo(roleId: number, errorback?: boolean) {
    return apiLiveAdminPost<{
      menu: AccessNode[];
      model: PlatformRoleRow;
      select_menu: number[];
    }>('/admin/platform/role/edit-info', { role_id: roleId }, errorback);
  },
  roleEdit(roleId: number, params: string, errorback?: boolean) {
    return apiLiveAdminPost(
      '/admin/platform/role/edit',
      { role_id: roleId, params },
      errorback,
    );
  },
  roleDelete(roleId: number, errorback?: boolean) {
    return apiLiveAdminPost(
      '/admin/platform/role/delete',
      { role_id: roleId },
      errorback,
    );
  },
};

export default PlatformRoleApi;
