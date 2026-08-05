import { requestClient } from '#/api/request';

export interface MemberLevel {
  id: number;
  name: string;
  rank: number;
  rules: string;
  benefits: string;
  status: number;
  version: number;
  assigned_count: number;
}

export interface MemberLevelInput {
  name: string;
  rank: number;
  rules: string;
  benefits: string;
  status: number;
  version?: number;
}

export const listMemberLevels = () => requestClient.get<{ list: MemberLevel[] }>('/member-levels');
export const createMemberLevel = (data: MemberLevelInput) => requestClient.post<void>('/member-levels', data);
export const updateMemberLevel = (id: number, data: MemberLevelInput) => requestClient.put(`/member-levels/${id}`, data);
export const deleteMemberLevel = (id: number) => requestClient.delete(`/member-levels/${id}`);
