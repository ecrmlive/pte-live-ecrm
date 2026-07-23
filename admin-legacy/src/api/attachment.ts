import { http } from '@/api/http';

export interface AttachmentCategory {
  attachment_category_id: number;
  attachment_category_name: string;
  attachment_category_enname?: string;
  mer_id: number;
  sort: number;
}

export interface Attachment {
  attachment_id: number;
  attachment_category_id: number;
  attachment_name: string;
  attachment_src: string;
  upload_type: number;
  user_type: number;
}

export function fetchCategories() {
  return http.get<{ list: AttachmentCategory[] }>('/attachments/categories');
}

export function createCategory(data: Record<string, unknown>) {
  return http.post<AttachmentCategory>('/attachments/categories', data);
}

export function deleteCategory(id: number) {
  return http.delete(`/attachments/categories/${id}`);
}

export function fetchAttachments(params: Record<string, unknown>) {
  return http.get<{ list: Attachment[]; total: number }>('/attachments', { params });
}

export function deleteAttachment(id: number) {
  return http.delete(`/attachments/${id}`);
}

export function uploadAttachment(file: File, categoryId?: number) {
  const fd = new FormData();
  fd.append('file', file);
  if (categoryId) fd.append('category_id', String(categoryId));
  return http.post<Attachment>('/attachments/upload', fd, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
}
