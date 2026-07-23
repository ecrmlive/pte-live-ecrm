/** PHP 接口常返回 1 / "1" / true；与 legacy `== 1` 一致 */
export function flagOn(value: unknown): boolean {
  return value == 1 || value === true;
}

/** 与 legacy `== 0` 一致 */
export function flagOff(value: unknown): boolean {
  return value == 0 || value === false;
}

/** 与 legacy `== expected` 一致（兼容字符串数字） */
export function numEq(value: unknown, expected: number): boolean {
  return value == expected;
}
