import { requestClient } from '#/api/request';

export interface H5DomainRow {
  domain_id: number;
  app_id: number;
  domain: string;
  cert_type: number;
  cert_type_text: string;
  category: number;
  category_text: string;
  biz_status: number;
  platform_status: number;
  deploy_status: string;
  deploy_error?: string;
  cert_not_after?: number;
  cert_revision?: number;
  deployed_cert_revision?: number;
  last_deploy_at?: number;
  create_time?: string;
  enable_status?: number;
  enable_status_text?: string;
  cert_sync_status?: string;
  cert_sync_status_text?: string;
  deploy_display_status?: string;
  deploy_display_text?: string;
  merchant_status_text?: string;
  biz_disabled_pending?: boolean;
  can_edit_cert?: boolean;
  can_delete?: boolean;
}

export interface H5DomainListResult {
  list: {
    data: H5DomainRow[];
    total: number;
    per_page: number;
    current_page: number;
  };
  gateway_cname?: string;
  max_count?: number;
}

export function getH5DomainListApi(params: { page?: number; list_rows?: number }) {
  return requestClient.post<H5DomainListResult>('/shop/setting.h5domain/index', params);
}

export function addH5DomainApi(data: {
  domain: string;
  cert_type: number;
  category: number;
  cert_key: string;
  cert_pem: string;
}) {
  return requestClient.post('/shop/setting.h5domain/add', data);
}

export function editH5DomainApi(data: {
  domain_id: number;
  cert_type?: number;
  cert_key: string;
  cert_pem: string;
}) {
  return requestClient.post('/shop/setting.h5domain/edit', data);
}

export function setH5DomainBizStatusApi(data: { domain_id: number; biz_status: 0 | 1 }) {
  return requestClient.post('/shop/setting.h5domain/status', data);
}

export function deleteH5DomainApi(data: { domain_id: number }) {
  return requestClient.post('/shop/setting.h5domain/delete', data);
}

export interface EffectiveH5DomainItem {
  domain_id: number;
  domain: string;
  cert_type: number;
  origin: string;
  category: number;
}

export function getEffectiveH5DomainsApi() {
  return requestClient.post<{ list: EffectiveH5DomainItem[] }>(
    '/shop/setting.h5domain/effective-list',
    {},
  );
}

export function merchantStatusTagType(text: string) {
  if (text.includes('已启用')) return 'success';
  if (text === '待部署' || text.includes('待更新')) return 'warning';
  if (text === '部署失败') return 'danger';
  if (text === '已停用' || text === '未启用') return 'info';
  return 'info';
}

export function deployDisplayTagType(status?: string) {
  switch (status) {
    case 'deployed':
      return 'success';
    case 'updated':
      return 'warning';
    case 'failed':
      return 'danger';
    case 'offline':
      return 'info';
    case 'deploying':
      return 'warning';
    default:
      return 'info';
  }
}

export function certSyncTagType(status?: string) {
  if (status === 'outdated') return 'warning';
  if (status === 'latest') return 'success';
  return 'info';
}
