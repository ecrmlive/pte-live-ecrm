export function isAllSpace(val: string) {
  return /^[ ]*$/.test(val);
}

export function replaceSpace(val?: string) {
  if (val != undefined) {
    return val.replace(/\s*/g, '');
  }
  return '';
}

export function hasSpace(val?: string) {
  if (val != undefined) {
    return /\s/g.test(val);
  }
  return false;
}
