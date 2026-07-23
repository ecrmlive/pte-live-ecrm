import { adminPost } from '#/utils/admin-api';

const ShopApi = {
  shopList(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/shop/index', data, errorback);
  },
  addShop(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/shop/add', data, errorback);
  },
  editShop(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/shop/edit', data, errorback);
  },
  updateStatus(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/shop/updateStatus', data, errorback);
  },
  storeEnter(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/shop/enter', data, errorback);
  },
  deleteShop(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/shop/delete', data, errorback);
  },
  updateWxStatus(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/shop/updateWxStatus', data, errorback);
  },
};

export default ShopApi;
