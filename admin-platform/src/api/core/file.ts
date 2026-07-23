import { adminPost, adminUpload } from '#/utils/admin-api';

const FileApi = {
  SystemPictureList(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/file.image/list', data, errorback);
  },
  PictureIndex(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/file.image/index', data, errorback);
  },
  deleteFiles(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/file.image/deleteFiles', data, errorback);
  },
  addCategory(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/file.image/addCategory', data, errorback);
  },
  editCategory(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/file.image/edit', data, errorback);
  },
  deleteCategory(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/file.image/delete', data, errorback);
  },
  uploadPresign(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/file.upload/presign', data, errorback);
  },
  uploadConfirm(data: Record<string, unknown>, errorback?: boolean) {
    return adminPost('/admin/file.upload/confirm', data, errorback);
  },
  async uploadCosDirect(
    file: File,
    options: { category_id?: number; file_type?: string } = {},
    errorback?: boolean,
  ) {
    const fileType = options.file_type || 'image';
    const categoryId = options.category_id ?? 0;
    const presignRes = await this.uploadPresign(
      {
        file_name: file.name,
        file_type: fileType,
        file_size: file.size,
        category_id: categoryId,
      },
      errorback,
    );
    const info = (presignRes as { data?: Record<string, unknown> }).data || {};
    const uploadUrl = info.upload_url as string;
    const key = info.key as string;
    const contentType =
      (info.content_type as string) || file.type || 'application/octet-stream';
    const method = String(info.method || 'PUT').toUpperCase();
    if (!uploadUrl || !key) {
      throw new Error('获取云存储上传地址失败');
    }
    const putRes = await fetch(uploadUrl, {
      method,
      body: file,
      headers: { 'Content-Type': contentType },
    });
    if (!putRes.ok) {
      throw new Error(`云存储上传失败 HTTP ${putRes.status}`);
    }
    return this.uploadConfirm(
      {
        key,
        file_type: fileType,
        real_name: file.name,
        file_size: file.size,
        category_id: categoryId,
      },
      errorback,
    );
  },
  moveFile(formData: FormData, errorback?: boolean) {
    return adminUpload('/admin/file.image/moveFiles', formData, errorback);
  },
};

export default FileApi;
