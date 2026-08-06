import { requestClient } from '#/api/request';

export type UserAssetType = 'balance' | 'commission' | 'points';

export interface UserAssetLedgerRow {
  amount: number;
  asset_type: UserAssetType;
  created_at: string;
  id: number;
  reference_id: string;
  reference_type: string;
  user_id: number;
}

export interface UserAssetLedgerPage {
  limit: number;
  list: UserAssetLedgerRow[];
  page: number;
  total: number;
}

export interface UserAssetSummary {
  asset_type: UserAssetType;
  count: number;
  expense: number;
  income: number;
}

export function listPlatformUserAssetsApi(params: {
  asset_type?: UserAssetType;
  limit: number;
  page: number;
  user_id?: number;
  keyword?: string;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<UserAssetLedgerPage>('/finance/user-assets', { params });
}

export function getPlatformUserAssetSummaryApi(params?: {
  asset_type?: UserAssetType;
  user_id?: number;
  keyword?: string;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<{ list: UserAssetSummary[] }>(
    '/finance/user-assets/summary',
    { params },
  );
}
