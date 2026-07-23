export interface ProductSpecItem {
  item_id: number;
  spec_value: string;
}

export interface ProductSpecAttrGroup {
  group_id: number;
  group_name: string;
  loading?: boolean;
  spec_items: ProductSpecItem[];
  tempValue?: string;
}

export interface ProductSpecFormData {
  card_info?: string;
  card_type?: number;
  image_id?: number;
  image_path?: string;
  line_price?: number | string;
  low_price?: number | string;
  product_no?: string;
  product_price?: number | string;
  product_weight?: number | string;
  stock_num?: number | string;
  virtualInfo?: Array<{ card_no: string; card_pwd: string }> | string;
}

export interface ProductSpecListRow {
  product_sku_id: number;
  rows: ProductSpecItem[];
  spec_form: ProductSpecFormData;
  spec_sku_id: string;
}

export interface ProductSpecMany {
  spec_attr: ProductSpecAttrGroup[];
  spec_list: ProductSpecListRow[];
}

export function calcDescartes<T>(array: T[][]): T[] | T[][] {
  if (array.length < 2) {
    return array[0] || [];
  }
  return array.reduce<T[][]>((col, set) => {
    const res: T[][] = [];
    col.forEach((c) => {
      set.forEach((s) => {
        const t = ([] as T[]).concat(Array.isArray(c) ? c : [c]);
        t.push(s);
        res.push(t);
      });
    });
    return res;
  });
}

export function buildProductSpecList(
  specAttr: ProductSpecAttrGroup[],
  existingList: ProductSpecListRow[],
): ProductSpecListRow[] {
  const specArr = specAttr.map((group) => group.spec_items);
  const specListTemp = calcDescartes(specArr);
  const specList: ProductSpecListRow[] = [];

  for (let i = 0; i < specListTemp.length; i++) {
    const item = specListTemp[i];
    const rows = Array.isArray(item) ? item : [item as ProductSpecItem];
    specList.push({
      product_sku_id: 0,
      spec_sku_id: rows.map((row) => row.item_id).join('_'),
      rows,
      spec_form: {},
    });
  }

  if (existingList.length > 0 && specList.length > 0) {
    for (let i = 0; i < specList.length; i++) {
      const overlap = existingList.filter(
        (row) => row.spec_sku_id === specList[i]!.spec_sku_id,
      );
      if (overlap.length > 0) {
        specList[i]!.spec_form = overlap[0]!.spec_form;
        specList[i]!.product_sku_id = overlap[0]!.product_sku_id;
      }
    }
  }

  return specList;
}

export function productSpecSpanMethod(
  specAttr: ProductSpecAttrGroup[],
  rowIndex: number,
  columnIndex: number,
) {
  if (columnIndex >= specAttr.length - 1) {
    return undefined;
  }
  let totalRow = 1;
  for (let i = columnIndex + 1; i < specAttr.length; i++) {
    totalRow *= specAttr[i]!.spec_items.length;
  }
  if (rowIndex % totalRow === 0) {
    return { rowspan: totalRow, colspan: 1 };
  }
  return { rowspan: 0, colspan: 0 };
}
