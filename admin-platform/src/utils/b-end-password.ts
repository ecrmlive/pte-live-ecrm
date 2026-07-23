export const B_END_PASSWORD_MIN = 6;
export const B_END_PASSWORD_MAX = 15;

/** 对齐 PHP PasswordService / api-platform ValidatePasswordStrength */
export function validateBEndPassword(value: unknown): string {
  const password = String(value ?? '');
  if (!password) return '请输入登录密码';
  if (/\s/.test(password)) return '密码不能包含空格';
  if (password.length < B_END_PASSWORD_MIN || password.length > B_END_PASSWORD_MAX) {
    return `密码长度需为${B_END_PASSWORD_MIN}-${B_END_PASSWORD_MAX}位`;
  }
  if (!/[A-Za-z]/.test(password) || !/\d/.test(password)) {
    return '密码需同时包含字母和数字';
  }
  return '';
}

export function isBEndPassword(value: unknown): boolean {
  return validateBEndPassword(value) === '';
}
