import type { VbenFormSchema } from '#/adapter/form';

export type FormRecord = Record<string, unknown>;

export function getByPath(obj: FormRecord, path: string): unknown {
  return path.split('.').reduce<unknown>((acc, key) => {
    if (acc && typeof acc === 'object') {
      return (acc as FormRecord)[key];
    }
    return undefined;
  }, obj);
}

export function setByPath(obj: FormRecord, path: string, value: unknown): void {
  const keys = path.split('.');
  let current: FormRecord = obj;
  for (let i = 0; i < keys.length - 1; i++) {
    const key = keys[i]!;
    if (!(key in current) || typeof current[key] !== 'object' || current[key] === null) {
      current[key] = {};
    }
    current = current[key] as FormRecord;
  }
  current[keys[keys.length - 1]!] = value;
}

export function pickPaths(obj: FormRecord, paths: string[]): FormRecord {
  const result: FormRecord = {};
  for (const path of paths) {
    result[path] = getByPath(obj, path);
  }
  return result;
}

/** Nested object for Vben form setValues (field names like `style.rowsNum`). */
export function pickPathsNested(obj: FormRecord, paths: string[]): FormRecord {
  const result: FormRecord = {};
  for (const path of paths) {
    setByPath(result, path, getByPath(obj, path));
  }
  return result;
}

export function extractFieldNames(schema: VbenFormSchema[]): string[] {
  return schema
    .map((item) => item.fieldName)
    .filter((name): name is string => Boolean(name && !name.startsWith('_')));
}

export function parseIntFields(obj: FormRecord, paths: string[]): void {
  for (const path of paths) {
    const value = getByPath(obj, path);
    if (value !== undefined && value !== null && value !== '') {
      setByPath(obj, path, Number.parseInt(String(value), 10));
    }
  }
}
