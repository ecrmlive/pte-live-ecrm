import { requestClient } from '#/api/request';

export type AttachmentKind = 'image' | 'video';

export interface AttachmentCategory {
  attachment_category_id: number;
  attachment_category_name: string;
  attachment_category_enname: string;
  pid: number;
  sort: number;
  is_system?: number;
}

export interface AttachmentItem {
  attachment_id: number;
  attachment_category_id: number;
  attachment_name: string;
  attachment_src: string;
  attachment_type: 0 | 1;
  is_system?: number;
  create_time: string;
}

export interface AttachmentPage {
  limit: number;
  list: AttachmentItem[];
  page: number;
  total: number;
}

export function listAttachmentCategoriesApi(params?: { type?: AttachmentKind }) {
  return requestClient.get<{ list: AttachmentCategory[] }>('/attachments/categories', {
    params,
  });
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
  /** 1 = 仅行级系统预置素材（侧栏「系统素材」） */
  is_system?: 0 | 1;
  keyword?: string;
  limit?: number;
  page?: number;
  type?: AttachmentKind;
}) {
  return requestClient.get<AttachmentPage>('/attachments', { params });
}

export function uploadAttachmentApi(
  file: File,
  categoryID = 0,
  options?: { isSystem?: boolean },
) {
  const form = new FormData();
  form.append('file', file);
  form.append('category_id', String(categoryID));
  if (options?.isSystem) {
    form.append('is_system', '1');
  }
  return requestClient.post<AttachmentItem>('/attachments/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
}

export function deleteAttachmentApi(id: number) {
  return requestClient.delete(`/attachments/${id}`);
}

export function moveAttachmentsApi(attachmentIDs: number[], categoryID: number) {
  return requestClient.patch('/attachments/move', {
    attachment_ids: attachmentIDs,
    category_id: categoryID,
  });
}
