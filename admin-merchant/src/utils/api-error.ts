const TECHNICAL_ERROR_PATTERNS = [
  /SQLSTATE/i,
  /Connection refused/i,
  /PDOException/i,
  /think\\db/i,
  /HY000/i,
  /Access denied for user/i,
  /Fatal error/i,
];

type ApiErrorLike = {
  message?: string;
  status?: number;
  code?: string;
  data?: Record<string, unknown>;
  response?: {
    status?: number;
    data?: Record<string, unknown>;
  };
};

export function isTechnicalApiErrorMessage(message: string): boolean {
  const text = String(message || '').trim();
  if (!text) return false;
  return TECHNICAL_ERROR_PATTERNS.some((pattern) => pattern.test(text));
}

/** 将接口/网络错误转为用户可读文案，过滤 SQL、PDO 等技术信息 */
export function formatUserFacingApiError(
  error: unknown,
  fallback = '请求失败，请稍后再试',
): string {
  const err = error as ApiErrorLike;
  const responseData = err?.response?.data ?? err?.data ?? {};
  const raw =
    (responseData?.msg as string) ||
    (responseData?.message as string) ||
    err?.message ||
    '';
  const status = err?.response?.status ?? err?.status;

  if (
    err?.message?.includes?.('Network Error') ||
    raw.includes('Network Error')
  ) {
    return '网络异常，请检查您的网络连接后重试。';
  }

  if (
    err?.message?.includes?.('timeout') ||
    err?.code === 'ECONNABORTED'
  ) {
    return '网络异常，请检查您的网络连接后重试。';
  }

  if ((status && status >= 500) || isTechnicalApiErrorMessage(raw)) {
    if (status === 502 || status === 503 || status === 504) {
      return '网络异常，请检查您的网络连接后重试。';
    }
    return fallback;
  }

  if (raw && !isTechnicalApiErrorMessage(raw)) {
    return raw;
  }

  return fallback;
}
