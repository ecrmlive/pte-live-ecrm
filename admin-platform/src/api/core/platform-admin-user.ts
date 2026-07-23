import { apiLiveAdminPost } from '#/utils/api-live-admin';

export interface PlatformAdminUserRow {
  admin_user_id: number;
  create_time?: number;
  is_super?: number;
  role_ids?: number[];
  role_names?: string[];
  update_time?: number;
  user_name: string;
}

const PlatformAdminUserApi = {
  userList(errorback?: boolean) {
    return apiLiveAdminPost<{
      list: PlatformAdminUserRow[];
      roleList: { role_id: number; role_name: string }[];
    }>('/admin/platform/admin-user/index', {}, errorback);
  },
  userAdd(
    data: { password: string; role_ids: number[]; user_name: string },
    errorback?: boolean,
  ) {
    return apiLiveAdminPost('/admin/platform/admin-user/add', data, errorback);
  },
  userEditInfo(adminUserId: number, errorback?: boolean) {
    return apiLiveAdminPost<{
      roleList: { role_id: number; role_name: string }[];
      role_ids: number[];
      user: PlatformAdminUserRow;
    }>('/admin/platform/admin-user/edit-info', { admin_user_id: adminUserId }, errorback);
  },
  userEdit(
    adminUserId: number,
    data: { is_super?: number; role_ids: number[] },
    errorback?: boolean,
  ) {
    return apiLiveAdminPost('/admin/platform/admin-user/edit', {
      admin_user_id: adminUserId,
      ...data,
    }, errorback);
  },
};

export default PlatformAdminUserApi;
