import { requestClient } from '#/api/request';

export interface MarketingCacheListItem {
  id: string;
  name: string;
  enabled: boolean;
  remark: string;
}

function listApi(path: string) {
  return requestClient.get<{ list: MarketingCacheListItem[] }>(path);
}

function saveApi(path: string, list: MarketingCacheListItem[]) {
  return requestClient.put<{ list: MarketingCacheListItem[] }>(path, { list });
}

export const fetchMarketingDiscounts = () => listApi('/marketing/discounts');
export const saveMarketingDiscounts = (list: MarketingCacheListItem[]) => saveApi('/marketing/discounts', list);
export const fetchMarketingApplications = () => listApi('/marketing/applications');
export const saveMarketingApplications = (list: MarketingCacheListItem[]) => saveApi('/marketing/applications', list);
export const fetchMarketingAtmosphere = () => listApi('/marketing/atmosphere');
export const saveMarketingAtmosphere = (list: MarketingCacheListItem[]) => saveApi('/marketing/atmosphere', list);
export const fetchMarketingBorder = () => listApi('/marketing/border');
export const saveMarketingBorder = (list: MarketingCacheListItem[]) => saveApi('/marketing/border', list);
export const fetchMarketingTopic = () => listApi('/marketing/topic');
export const saveMarketingTopic = (list: MarketingCacheListItem[]) => saveApi('/marketing/topic', list);
