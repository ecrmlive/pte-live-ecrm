import { requestClient } from '#/api/request';

export interface SvipInterest {
  id: number;
  name: string;
  display_name: string;
  description: string;
  icon_url: string;
  on_icon_url: string;
  link: string;
  status: number;
  sort: number;
  version: number;
}

export interface SvipInterestInput {
  name: string;
  display_name: string;
  description: string;
  icon_url: string;
  on_icon_url: string;
  link: string;
  status: number;
  sort: number;
  version?: number;
}

export const listSvipInterests = () =>
  requestClient.get<{ list: SvipInterest[] }>('/svip/interests');

export const createSvipInterest = (data: SvipInterestInput) =>
  requestClient.post<void>('/svip/interests', data);

export const updateSvipInterest = (id: number, data: SvipInterestInput) =>
  requestClient.put(`/svip/interests/${id}`, data);

export const updateSvipInterestStatus = (
  id: number,
  data: { status: number; version: number },
) => requestClient.put(`/svip/interests/${id}/status`, data);

export const deleteSvipInterest = (id: number) =>
  requestClient.delete(`/svip/interests/${id}`);
