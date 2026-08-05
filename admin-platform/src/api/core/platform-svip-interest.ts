import { requestClient } from '#/api/request';

export interface SvipInterest {
  id: number;
  name: string;
  description: string;
  icon_url: string;
  status: number;
  sort: number;
  version: number;
}

export interface SvipInterestInput {
  name: string;
  description: string;
  icon_url: string;
  status: number;
  sort: number;
  version?: number;
}

export const listSvipInterests = () => requestClient.get<{ list: SvipInterest[] }>('/svip/interests');
export const createSvipInterest = (data: SvipInterestInput) => requestClient.post<void>('/svip/interests', data);
export const updateSvipInterest = (id: number, data: SvipInterestInput) => requestClient.put(`/svip/interests/${id}`, data);
export const deleteSvipInterest = (id: number) => requestClient.delete(`/svip/interests/${id}`);
