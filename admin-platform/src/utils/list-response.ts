/** 解析 api-platform / PHP 列表：支持 list[] 或 list{ data, total } */
export function parseApiList<T>(payload: unknown): T[] {
  if (Array.isArray(payload)) {
    return payload as T[];
  }
  if (!payload || typeof payload !== 'object') {
    return [];
  }
  const root = payload as Record<string, unknown>;
  const list = root.list;
  if (Array.isArray(list)) {
    return list as T[];
  }
  if (list && typeof list === 'object') {
    const data = (list as { data?: unknown }).data;
    if (Array.isArray(data)) {
      return data as T[];
    }
  }
  return [];
}

export function parseApiListPage<T>(payload: unknown): {
  list: T[];
  total: number;
} {
  if (!payload || typeof payload !== 'object') {
    return { list: [], total: 0 };
  }
  const root = payload as Record<string, unknown>;
  const list = root.list;
  if (list && typeof list === 'object' && !Array.isArray(list)) {
    const page = list as { data?: T[]; total?: number };
    return {
      list: Array.isArray(page.data) ? page.data : [],
      total: Number(page.total ?? 0),
    };
  }
  const rows = parseApiList<T>(payload);
  return { list: rows, total: rows.length };
}
