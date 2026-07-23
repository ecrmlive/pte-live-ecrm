import { apiLiveAdminPost } from '#/utils/api-live-admin';
import type { AccessNode } from '#/views/access/types';

export interface MerchantTemplateRoleRow {
  create_time?: number;
  remark?: string;
  role_id: number;
  role_name: string;
  sort: number;
  update_time?: number;
}

export interface MerchantTemplateRoleForm {
  access_id: number[];
  remark?: string;
  role_id?: number;
  role_name: string;
  sort: number;
}

const MerchantRoleApi = {
  roleList(errorback?: boolean) {
    return apiLiveAdminPost<{ list: MerchantTemplateRoleRow[] }>(
      '/admin/platform/merchant-role/index',
      {},
      errorback,
    );
  },
  roleAddInfo(errorback?: boolean) {
    return apiLiveAdminPost<{ menu: AccessNode[] }>(
      '/admin/platform/merchant-role/add-info',
      {},
      errorback,
    );
  },
  roleAdd(params: string, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/merchant-role/add', { params }, errorback);
  },
  roleEditInfo(roleId: number, errorback?: boolean) {
    return apiLiveAdminPost<{
      menu: AccessNode[];
      model: MerchantTemplateRoleRow;
      select_menu: number[];
    }>('/admin/platform/merchant-role/edit-info', { role_id: roleId }, errorback);
  },
  roleEdit(roleId: number, params: string, errorback?: boolean) {
    return apiLiveAdminPost(
      '/admin/platform/merchant-role/edit',
      { role_id: roleId, params },
      errorback,
    );
  },
  roleDelete(roleId: number, errorback?: boolean) {
    return apiLiveAdminPost(
      '/admin/platform/merchant-role/delete',
      { role_id: roleId },
      errorback,
    );
  },
};

export default MerchantRoleApi;
