import { adminGet, adminPost } from '#/utils/admin-api';

const PlugsApi = {
  plugslist(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/plus.plus/index', data, errorback);
  },
  getplugs(data: Record<string, unknown>, errorback?: boolean) {
    return adminGet('/admin/plus.plus/add', data, errorback);
  },
  addplugs(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/plus.plus/add', data, errorback);
  },
  editplugs(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/plus.plus/edit', data, errorback);
  },
  deleteplugs(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/plus.plus/delete', data, errorback);
  },
  updatePlugsStatus(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/plus.plus/updateStatus', data, errorback);
  },
  updatePlugsRecom(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/plus.plus/updateRecom', data, errorback);
  },
};

export default PlugsApi;
