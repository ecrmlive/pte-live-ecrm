import { apiLiveAdminPost } from '#/utils/api-live-admin';

const PlatformAccessApi = {
  accessList(data: Record<string, unknown>, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/access/index', data, errorback);
  },
  addAccess(data: Record<string, unknown>, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/access/add', data, errorback);
  },
  editAccess(data: Record<string, unknown>, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/access/edit', data, errorback);
  },
  delAccess(data: Record<string, unknown>, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/access/delete', data, errorback);
  },
  status(data: Record<string, unknown>, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/access/status', data, errorback);
  },
};

export default PlatformAccessApi;
