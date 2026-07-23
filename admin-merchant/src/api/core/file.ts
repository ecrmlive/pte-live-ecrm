import { requestClient } from '#/api/request';

export interface ShopFileItem {
  file_id: number;
  file_path: string;
  file_name?: string;
  real_name?: string;
  selected?: boolean;
}

export interface ShopFileGroupItem {
  group_id: null | number;
  group_name: string;
}

export interface ShopFileListResult {
  file_list: {
    data: ShopFileItem[];
    total: number;
  };
}

export interface ShopFileCategoryResult {
  group_list: ShopFileGroupItem[];
}

/** Platform system icon bank (`pte_live_image_bank`, COS `pte-live/system/*`). */
export interface ShopSystemImageCategoryItem {
  category_id: number;
  name: string;
}

export interface ShopSystemImageItem {
  category_id: number;
  image: string;
  name: string;
}

export interface ShopSystemImageListResult {
  list: {
    data: ShopSystemImageItem[];
    total: number;
  };
}

export interface ShopSystemImageCategoryResult {
  list: ShopSystemImageCategoryItem[];
}

export async function getShopFileCategoryApi(type: 'image' | 'video') {
  return requestClient.post<ShopFileCategoryResult>('/shop/file.file/category', {
    type,
  });
}

export async function getShopSystemImageCategoryApi() {
  return requestClient.post<ShopSystemImageCategoryResult>('/shop/file.image/index', {});
}

export async function getShopSystemImageListApi(params: {
  list_rows?: number;
  page?: number;
  parentId?: null | number;
}) {
  const payload: Record<string, unknown> = {
    fileType: 'image',
    list_rows: params.list_rows ?? 36,
    page: params.page ?? 1,
  };
  if (params.parentId != null && params.parentId > 0) {
    payload.parentId = params.parentId;
  }
  return requestClient.post<ShopSystemImageListResult>('/shop/file.image/list', payload);
}

export async function getShopFileListApi(params: {
  group_id?: null | number;
  list_rows?: number;
  page?: number;
  type?: string;
}) {
  const payload: Record<string, unknown> = {
    list_rows: params.list_rows ?? 36,
    page: params.page ?? 1,
    type: params.type ?? 'image',
  };
  if (params.group_id != null) {
    payload.group_id = params.group_id;
  }
  return requestClient.post<ShopFileListResult>('/shop/file.file/lists', payload);
}

function serializeFileIds(fileIds: number[]) {
  return fileIds.filter((id) => id > 0).join(',');
}

export async function deleteShopFilesApi(fileIds: number[]) {
  return requestClient.post<{ msg?: string }>('/shop/file.file/deleteFiles', {
    fileIds: serializeFileIds(fileIds),
  });
}

export async function addShopFileGroupApi(payload: {
  group_name: string;
  group_type: string;
}) {
  return requestClient.post<{ msg?: string }>('/shop/file.file/addGroup', payload);
}

export async function editShopFileGroupApi(payload: {
  group_id: number;
  group_name: string;
}) {
  return requestClient.post<{ msg?: string }>('/shop/file.file/editGroup', payload);
}

export async function deleteShopFileGroupApi(groupId: number) {
  return requestClient.post<{ msg?: string }>('/shop/file.file/deleteGroup', {
    group_id: groupId,
  });
}

export async function moveShopFilesApi(payload: {
  fileIds: number[];
  group_id: number | null;
}) {
  const body: Record<string, unknown> = {
    fileIds: serializeFileIds(payload.fileIds),
  };
  if (payload.group_id != null) {
    body.group_id = payload.group_id;
  }
  return requestClient.post<{ msg?: string }>(
    '/shop/file.upload/moveFiles',
    body,
  );
}

export async function uploadShopFilePresignApi(params: {
  file_name: string;
  file_size: number;
  file_type?: string;
  group_id?: number;
}) {
  return requestClient.post<{
    content_type?: string;
    key: string;
    method?: string;
    upload_url: string;
  }>('/shop/file.upload/presign', {
    file_type: params.file_type ?? 'image',
    group_id: params.group_id ?? 0,
    ...params,
  });
}

export async function uploadShopFileConfirmApi(params: {
  file_size: number;
  file_type?: string;
  group_id?: number;
  key: string;
  real_name: string;
}) {
  return requestClient.post<{ file_id: number; file_path: string }>(
    '/shop/file.upload/confirm',
    {
      file_type: params.file_type ?? 'image',
      group_id: params.group_id ?? 0,
      ...params,
    },
  );
}

export async function uploadShopFileDirect(
  file: File,
  options: { file_type?: string; group_id?: number } = {},
) {
  const fileType = options.file_type ?? 'image';
  const groupId = options.group_id ?? 0;
  const presign = await uploadShopFilePresignApi({
    file_name: file.name,
    file_size: file.size,
    file_type: fileType,
    group_id: groupId,
  });
  const method = (presign.method ?? 'PUT').toUpperCase();
  const putRes = await fetch(presign.upload_url, {
    body: file,
    headers: {
      'Content-Type': presign.content_type ?? file.type ?? 'application/octet-stream',
    },
    method,
  });
  if (!putRes.ok) {
    throw new Error(`云存储上传失败 HTTP ${putRes.status}`);
  }
  return uploadShopFileConfirmApi({
    file_size: file.size,
    file_type: fileType,
    group_id: groupId,
    key: presign.key,
    real_name: file.name,
  });
}

export async function uploadShopImageDirect(file: File, groupId = 0) {
  return uploadShopFileDirect(file, { file_type: 'image', group_id: groupId });
}
