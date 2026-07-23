export function mergeTable<T extends { product_id?: number; rowSpan?: number | null }>(list: T[]) {
  let curItem = {
    index: 0,
    isFirst: false,
    product_id: undefined as number | undefined,
    rowSpan: 1,
  };
  for (let i = 0; i < list.length; i++) {
    const item = list[i];
    item.rowSpan = null;
    if (!curItem.isFirst) {
      curItem.isFirst = true;
      curItem.index = i;
      curItem.product_id = item.product_id;
    } else if (curItem.product_id !== item.product_id) {
      list[curItem.index]!.rowSpan = curItem.rowSpan;
      curItem.rowSpan = 1;
      curItem.index = i;
      curItem.product_id = item.product_id;
    } else {
      curItem.rowSpan++;
    }
  }
  if (list.length > 0) {
    list[curItem.index]!.rowSpan = curItem.rowSpan;
  }
  return list;
}
