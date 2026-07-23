export type FormSelectOption = {
  label: string;
  value: number | string;
};

/** Match API value to an option's value type so ElSelect shows label, not raw id. */
export function matchFormSelectValue(
  value: unknown,
  options: FormSelectOption[],
  fallback?: number | string,
): number | string {
  if (value == null || value === '') {
    return fallback ?? options[0]?.value ?? '';
  }
  const matched = options.find(
    (opt) =>
      opt.value === value ||
      String(opt.value) === String(value) ||
      Number(opt.value) === Number(value),
  );
  if (matched) {
    return matched.value;
  }
  return fallback ?? (value as number | string);
}
