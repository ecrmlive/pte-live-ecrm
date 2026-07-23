import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface TableListItem {
  create_time: string;
  name: string;
  table_id: number;
  total_count: number;
}

export interface TableOptionItem {
  name: string;
  table_id: number;
}

export interface TableRecordItem {
  create_time: string;
  table_record_id: number;
  tableData?: Array<{ name: string; value: string }>;
  tableM?: { name: string };
  user?: { nickName: string; user_id: number };
}

export interface TableRecordQuery {
  list_rows?: number;
  page?: number;
  search?: string;
  table_id?: number;
}

export async function getTableListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<TableListItem> }>(
    '/shop/plus.table.table/index',
    params,
  );
}

export async function deleteTableApi(tableId: number) {
  return requestClient.post('/shop/plus.table.table/delete', { table_id: tableId });
}

export interface TableFieldRow {
  is_required: boolean;
  name: string;
  rule: string;
  select_value?: string;
  type: string;
  _uid?: number;
}

export interface TableFormModel {
  name: string;
  sort: number;
  table_id?: number;
  tableData: TableFieldRow[];
}

export interface TableEditMetaResult {
  model: TableFormModel;
}

export async function addTableApi(payload: TableFormModel) {
  return requestClient.post('/shop/plus.table.table/add', payload);
}

export async function getTableEditMetaApi(tableId: number) {
  return requestClient.get<TableEditMetaResult>('/shop/plus.table.table/edit', {
    params: { table_id: tableId },
  });
}

export async function editTableApi(payload: TableFormModel) {
  return requestClient.post('/shop/plus.table.table/edit', payload);
}

export async function getTableRecordListApi(params: TableRecordQuery) {
  return requestClient.post<{
    list: PaginatedList<TableRecordItem>;
    table_list: TableOptionItem[];
  }>('/shop/plus.table.record/index', params);
}

export async function deleteTableRecordApi(tableRecordId: number) {
  return requestClient.post('/shop/plus.table.record/delete', {
    table_record_id: tableRecordId,
  });
}

function saveExportBlob(blob: Blob) {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `table-records-${Date.now()}.xlsx`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}

export async function exportTableRecordApi(params: { search?: string; table_id: number }) {
  const blob = await requestClient.download<Blob>('/shop/plus.table.record/export', {
    params,
  });
  saveExportBlob(blob);
}
