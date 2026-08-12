import { requestClient } from '#/api/request';
import { resolveApiBaseUrl } from '#/utils/pte-live-api';
import { getDecryptedToken } from '#/utils/pte-live-token';

export type DatabaseScope = 'admin' | 'business';

export interface DatabaseTableRow {
  database_scope: DatabaseScope;
  engine: string;
  row_count: number;
  size_bytes: number;
  table_comment: string;
  table_name: string;
  updated_at?: string;
}

export interface DatabaseTableColumn {
  column_comment: string;
  column_default?: string | null;
  column_name: string;
  column_type: string;
  extra: string;
  is_nullable: string;
}

export interface DatabaseBackupRecord {
  created_at: string;
  created_by: number;
  database_scope: DatabaseScope;
  file_name: string;
  id: number;
  size_bytes: number;
  status: 'deleted' | 'failed' | 'ready';
  table_count: number;
  table_names: string[];
}

export function listDatabaseTablesApi(params: { limit: number; page: number }) {
  return requestClient.get<{ list: DatabaseTableRow[]; total: number }>(
    '/maintain/database/tables',
    { params },
  );
}

export function getDatabaseTableDetailApi(scope: DatabaseScope, tableName: string) {
  return requestClient.get<{
    columns: DatabaseTableColumn[];
    database_scope: DatabaseScope;
    table_name: string;
  }>(`/maintain/database/tables/${scope}/${encodeURIComponent(tableName)}`);
}

export function createDatabaseBackupApi(data: {
  scope: DatabaseScope;
  tables: string[];
}) {
  return requestClient.post<{ file_name: string; size_bytes: number; table_count: number }>(
    '/maintain/database/backups',
    data,
  );
}

export function listDatabaseBackupsApi(params: { limit: number; page: number }) {
  return requestClient.get<{ list: DatabaseBackupRecord[]; total: number }>(
    '/maintain/database/backups',
    { params },
  );
}

export function deleteDatabaseBackupApi(id: number) {
  return requestClient.delete(`/maintain/database/backups/${id}`);
}

export function maintainDatabaseTablesApi(
  action: 'optimize' | 'repair',
  data: { scope: DatabaseScope; tables: string[] },
) {
  return requestClient.post<{ ok: boolean; tables: string[] }>(
    `/maintain/database/tables/${action}`,
    data,
  );
}

export async function downloadDatabaseBackupApi(id: number) {
  const token = getDecryptedToken();
  const response = await fetch(
    `${resolveApiBaseUrl()}/api/platform/v1/maintain/database/backups/${id}/download`,
    {
      headers: token ? { 'Authori-zation': `Bearer ${token}` } : {},
    },
  );
  if (!response.ok) {
    let message = '下载备份失败';
    try {
      const body = (await response.json()) as { message?: string };
      message = body.message || message;
    } catch {
      // Keep the user-facing fallback when a proxy returns a non-JSON error.
    }
    throw new Error(message);
  }
  const disposition = response.headers.get('content-disposition') || '';
  const encoded = disposition.match(/filename\*=UTF-8''([^;]+)/i)?.[1];
  const plain = disposition.match(/filename="?([^";]+)"?/i)?.[1];
  return {
    blob: await response.blob(),
    fileName: encoded ? decodeURIComponent(encoded) : (plain || `database-backup-${id}.sql`),
  };
}
