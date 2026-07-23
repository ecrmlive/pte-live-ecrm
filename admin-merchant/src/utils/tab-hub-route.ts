/** Tab 容器页：根据当前 path 解析应激活的子 Tab key */
export function resolveTabKeyFromPath<T extends string>(
  path: string,
  sources: ReadonlyArray<{ key: T; path: string }>,
): null | T {
  const hit = sources.find((item) => item.path === path);
  return hit?.key ?? null;
}

/** path 优先，其次 ?type= query，最后 fallback */
export function syncHubTabFromRoute<T extends string>(
  path: string,
  queryType: unknown,
  pathSources: ReadonlyArray<{ key: T; path: string }>,
  queryKeys: readonly T[],
  fallback: T,
): T {
  const fromPath = resolveTabKeyFromPath(path, pathSources);
  if (fromPath) return fromPath;
  const type = String(queryType ?? '');
  if (queryKeys.includes(type as T)) return type as T;
  return fallback;
}
