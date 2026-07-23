export function setSessionStorage(name: string, val: unknown) {
  sessionStorage.setItem(name, JSON.stringify(val));
}

export function getSessionStorage(name: string) {
  if (sessionStorage.hasOwnProperty(name)) {
    return JSON.parse(sessionStorage.getItem(name) || 'null');
  }
  return false;
}

export function addSessionStorage(name: string, val: Record<string, unknown>) {
  if (sessionStorage.hasOwnProperty(name)) {
    const old = JSON.parse(sessionStorage.getItem(name) || '{}');
    sessionStorage.setItem(name, JSON.stringify({ ...old, ...val }));
  } else {
    sessionStorage.setItem(name, JSON.stringify(val));
  }
}

export function deleteSessionStorage(name: string | null = null) {
  if (name != null) {
    sessionStorage.removeItem(name);
  } else {
    sessionStorage.clear();
  }
}

export function deepClone<T>(obj: T): T {
  if (obj === null || typeof obj !== 'object') {
    return obj;
  }
  const objClone = (Array.isArray(obj) ? [] : {}) as T;
  for (const key in obj as Record<string, unknown>) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      const value = (obj as Record<string, unknown>)[key];
      (objClone as Record<string, unknown>)[key] =
        value && typeof value === 'object'
          ? deepClone(value)
          : value;
    }
  }
  return objClone;
}

export function deepMerger(
  obj1: Record<string, unknown>,
  obj2: Record<string, unknown>,
) {
  for (const key in obj2) {
    if (Object.prototype.hasOwnProperty.call(obj2, key)) {
      if (obj2[key] && typeof obj2[key] === 'object') {
        obj1[key] = deepMerger(
          (obj1[key] as Record<string, unknown>) || {},
          obj2[key] as Record<string, unknown>,
        );
      } else {
        obj1[key] = obj2[key];
      }
    }
  }
  return obj1;
}

export function formatModel(
  thisObj: Record<string, unknown>,
  sourceObj?: Record<string, unknown>,
) {
  for (const key in thisObj) {
    if (sourceObj && typeof sourceObj[key] !== 'undefined') {
      if (
        thisObj[key] &&
        Object.prototype.toString.call(thisObj[key]) === '[object Object]'
      ) {
        formatModel(
          thisObj[key] as Record<string, unknown>,
          sourceObj[key] as Record<string, unknown>,
        );
      } else {
        thisObj[key] = sourceObj[key];
      }
    }
  }
  return thisObj;
}
