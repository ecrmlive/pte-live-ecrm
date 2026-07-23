import { apiLiveAdminPost } from '#/utils/api-live-admin';

const MerchantAccessApi = {
  accessList(data: Record<string, unknown>, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/merchant-access/index', data, errorback);
  },
  addAccess(data: Record<string, unknown>, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/merchant-access/add', data, errorback);
  },
  editAccess(data: Record<string, unknown>, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/merchant-access/edit', data, errorback);
  },
  delAccess(data: Record<string, unknown>, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/merchant-access/delete', data, errorback);
  },
  status(data: Record<string, unknown>, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/merchant-access/status', data, errorback);
  },
};

export default MerchantAccessApi;
