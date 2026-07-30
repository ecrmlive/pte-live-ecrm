import { requestClient } from '#/api/request';

export type AttachmentKind = 'image' | 'video';

export interface AttachmentCategory {
  attachment_category_id: number;
  attachment_category_name: string;
  attachment_category_enname: string;
  pid: number;
  sort: number;
}

export interface AttachmentItem {
  attachment_id: number;
  attachment_category_id: number;
  attachment_name: string;
  attachment_src: string;
  attachment_type: 0 | 1;
  create_time: string;
}

export interface AttachmentPage {
  limit: number;
  list: AttachmentItem[];
  page: number;
  total: number;
}

export function listAttachmentCategoriesApi() {
  return requestClient.get<{ list: AttachmentCategory[] }>('/attachments/categories');
}

export function createAttachmentCategoryApi(data: {
  attachment_category_enname?: string;
  attachment_category_name: string;
  pid?: number;
  sort?: number;
}) {
  return requestClient.post<AttachmentCategory>('/attachments/categories', data);
}

export function updateAttachmentCategoryApi(
  id: number,
  data: {
    attachment_category_enname?: string;
    attachment_category_name: string;
    pid?: number;
    sort?: number;
  },
) {
  return requestClient.put<AttachmentCategory>(`/attachments/categories/${id}`, data);
}

export function deleteAttachmentCategoryApi(id: number) {
  return requestClient.delete(`/attachments/categories/${id}`);
}

export function listAttachmentsApi(params: {
  category_id?: number;
  limit?: number;
  page?: number;
  type?: AttachmentKind;
}) {
  return requestClient.get<AttachmentPage>('/attachments', { params });
}

export function uploadAttachmentApi(file: File, categoryID = 0) {
  const form = new FormData();
  form.append('file', file);
  form.append('category_id', String(categoryID));
  return requestClient.post<AttachmentItem>('/attachments/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
}

export function deleteAttachmentApi(id: number) {
  return requestClient.delete(`/attachments/${id}`);
}
