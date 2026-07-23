import { adminGet, adminPost } from '#/utils/admin-api';

const SettingApi = {
  serviceDetail(data: Record<string, unknown>, errorback?: boolean) {
    return adminGet('/admin/setting/index', data, errorback);
  },
  editService(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/setting/index', data, errorback);
  },
};

export default SettingApi;
