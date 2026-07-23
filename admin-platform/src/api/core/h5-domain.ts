export interface H5DomainRow {
  domain_id: number;
  app_id: number;
  domain: string;
  cert_type: number;
  cert_type_text: string;
  category?: number;
  category_text?: string;
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
  can_delete?: boolean;
  can_edit_cert?: boolean;
}

export interface H5DomainListResult {
  list: {
    data: H5DomainRow[];
    total: number;
    per_page?: number;
    current_page?: number;
  };
  gateway_cname?: string;
  max_count?: number;
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

export function enableStatusTagType(status?: number) {
  return status === 1 ? 'success' : 'info';
}
