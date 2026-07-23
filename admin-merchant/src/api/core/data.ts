import { requestClient } from '#/api/request';

export type RegionTree = Record<
  number,
  {
    city?: Record<
      number,
      {
        id: number;
        name: string;
        region?: Record<number, { id: number; name: string }>;
      }
    >;
    id: number;
    name: string;
  }
>;

export async function getRegionDataApi() {
  return requestClient.post<{ regionData: RegionTree }>('/shop/data.region/lists', {});
}
