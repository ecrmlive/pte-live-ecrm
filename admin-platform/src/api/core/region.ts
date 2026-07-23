import { adminGet, adminPost } from '#/utils/admin-api';

const RegionApi = {
  regionList(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/region/index', data, errorback);
  },
  toAddRegion(data: Record<string, unknown>, errorback?: boolean) {
    return adminGet('/admin/region/add', data, errorback);
  },
  addRegion(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/region/add', data, errorback);
  },
  regionDetail(data: Record<string, unknown>, errorback?: boolean) {
    return adminGet('/admin/region/edit', data, errorback);
  },
  editRegion(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/region/edit', data, errorback);
  },
  deleteRegion(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/region/delete', data, errorback);
  },
};

export default RegionApi;
