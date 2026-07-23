import type { H5DomainListResult, H5DomainRow } from '#/api/core/h5-domain';

import { apiLiveAdminPost } from '#/utils/api-live-admin';

const DEPLOY_REQUEST_TIMEOUT_MS = 130_000;

export type { H5DomainRow };

const PlatformH5DomainApi = {
  index(
    data: { app_id: number; page?: number; list_rows?: number },
    errorback = false,
  ) {
    return apiLiveAdminPost<H5DomainListResult>(
      '/admin/platform/shop-h5-domain/index',
      data,
      errorback,
    );
  },
  deploy(appId: number, domainId: number, errorback = false) {
    return apiLiveAdminPost(
      '/admin/platform/shop-h5-domain/deploy',
      { app_id: appId, domain_id: domainId },
      errorback,
      DEPLOY_REQUEST_TIMEOUT_MS,
      { skipErrorMessage: true },
    );
  },
  disable(appId: number, domainId: number, errorback = false) {
    return apiLiveAdminPost(
      '/admin/platform/shop-h5-domain/disable',
      { app_id: appId, domain_id: domainId },
      errorback,
      DEPLOY_REQUEST_TIMEOUT_MS,
      { skipErrorMessage: true },
    );
  },
  delete(appId: number, domainId: number, errorback = false) {
    return apiLiveAdminPost(
      '/admin/platform/shop-h5-domain/delete',
      { app_id: appId, domain_id: domainId },
      errorback,
    );
  },
  add(
    data: {
      app_id: number;
      category: number;
      cert_key: string;
      cert_pem: string;
      cert_type: number;
      domain: string;
    },
    errorback = false,
  ) {
    return apiLiveAdminPost('/admin/platform/shop-h5-domain/add', data, errorback);
  },
};

export default PlatformH5DomainApi;
