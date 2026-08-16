import type { Ref } from 'vue';

import type { DiyEditorContext } from './use-diy-editor';

export function changeSeckillColumn(
  curItem: Record<string, unknown>,
  column: number,
) {
  const params = (curItem.params as Record<string, unknown>) || {};
  params.column = column;
  curItem.params = params;
  const config = {
    linePrice: 1,
    productName: 1,
    productPrice: 1,
    productSales: 1,
    product_btn: 1,
  };
  switch (column) {
    case 2: {
      config.productSales = 0;
      break;
    }
    case 3: {
      config.productSales = 0;
      config.product_btn = 0;
      break;
    }
    case 4: {
      config.productSales = 0;
      break;
    }
  }
  Object.assign(curItem.params as Record<string, unknown>, config);
}

export function changeProductColumn(
  curItem: Record<string, unknown>,
  column: number,
) {
  const params = (curItem.params ?? {}) as Record<string, unknown>;
  const style = (curItem.style ?? {}) as Record<string, unknown>;
  curItem.params = params;
  curItem.style = style;
  params.column = column;
  style.column = column;

  const defaults = {
    cartAction: 'detail',
    cartType: 2,
    comment: 1,
    filterVisible: 0,
    linePrice: 1,
    memberPrice: 1,
    productComment: 1,
    productLabel: 1,
    productName: 1,
    productPrice: 1,
    productSales: 1,
    productScore: 1,
    showCart: 1,
    storeDistance: 1,
  };

  for (const [key, value] of Object.entries(defaults)) {
    if (params[key] === undefined || params[key] === null) {
      params[key] = value;
    }
  }
}

export function resetStyleColors(
  editor: DiyEditorContext,
  curItem: Record<string, unknown>,
  name1: string,
  name2?: string,
) {
  editor.onEditorResetColor(curItem.style as Record<string, unknown>, name1, '');
  if (name2) {
    editor.onEditorResetColor(curItem.style as Record<string, unknown>, name2, '');
  }
}

export function marketingDisplayFlag(value: unknown): boolean {
  return value === 1 || value === '1' || value === true;
}

export function resolveMarketingColumn(item: Record<string, unknown>): number {
  const params = (item.params as Record<string, unknown>) || {};
  const style = (item.style as Record<string, unknown>) || {};
  const column = params.column ?? style.column ?? 1;
  return Number(column) || 1;
}

/** Editor preview: coupon card count from params.limit (1–30). */
export function resolveCouponPreviewCount(item: Record<string, unknown>): number {
  const params = (item.params as Record<string, unknown>) || {};
  const limit = Number(params.limit) || 1;
  return Math.max(1, Math.min(30, limit));
}

/** Editor preview: pad/slice placeholder rows to match params.showNum. */
export function resolveMarketingPreviewProducts(
  item: Record<string, unknown>,
): Array<Record<string, unknown>> {
  const params = (item.params as Record<string, unknown>) || {};
  const showNum = Math.max(1, Number(params.showNum) || 1);
  const data = Array.isArray(item.data)
    ? (item.data as Array<Record<string, unknown>>)
    : [];
  const defaultData = Array.isArray(item.defaultData)
    ? (item.defaultData as Array<Record<string, unknown>>)
    : [];
  const template = data[0] ||
    defaultData[0] || {
      image: '',
      product_image: '',
      product_name: '此处是商品',
    };
  return Array.from(
    { length: showNum },
    (_, index) => data[index] ?? defaultData[index] ?? { ...template },
  );
}

/** Editor preview: auto/choice source lists with params.auto.showNum. */
export function resolveAutoSourcePreviewItems(
  item: Record<string, unknown>,
  template: Record<string, unknown>,
): Array<Record<string, unknown>> {
  const params = (item.params as Record<string, unknown>) || {};
  if (params.source === 'choice') {
    const data = Array.isArray(item.data)
      ? (item.data as Array<Record<string, unknown>>)
      : [];
    return data.length > 0 ? data : [{ ...template }];
  }
  const auto = (params.auto as Record<string, unknown>) || {};
  const showNum = Math.max(1, Number(auto.showNum) || 1);
  const data = Array.isArray(item.data)
    ? (item.data as Array<Record<string, unknown>>)
    : [];
  const defaultData = Array.isArray(item.defaultData)
    ? (item.defaultData as Array<Record<string, unknown>>)
    : [];
  const tpl = data[0] || defaultData[0] || template;
  return Array.from({ length: showNum }, (_, index) =>
    data[index] ?? defaultData[index] ?? { ...tpl },
  );
}

/** Editor preview: product group with params.auto.showNum or manual choice. */
export function resolveProductPreviewProducts(
  item: Record<string, unknown>,
): Array<Record<string, unknown>> {
  return resolveAutoSourcePreviewItems(item, {
    image: '',
    line_price: 0,
    member_price: '89.00',
    product_comment: 12,
    product_label: '精选',
    product_name: '此处是商品',
    product_price: '99.00',
    product_score: '5.0',
    product_sales: 0,
    store_distance: '1.2km',
  });
}

/** Editor preview: live lists with params.showNum placeholder rows. */
export function resolveLivePreviewItems(
  item: Record<string, unknown>,
): Array<Record<string, unknown>> {
  const params = (item.params as Record<string, unknown>) || {};
  const showNum = Math.max(1, Number(params.showNum) || 1);
  const data = Array.isArray(item.data)
    ? (item.data as Array<Record<string, unknown>>)
    : [];
  const defaultData = Array.isArray(item.defaultData)
    ? (item.defaultData as Array<Record<string, unknown>>)
    : [];
  const template = data[0] ||
    defaultData[0] || {
      image: '',
      name: '直播预告',
    };
  return Array.from({ length: showNum }, (_, index) => data[index] ?? { ...template });
}

export function changeBargainLikeColumn(
  curItem: Record<string, unknown>,
  column: number,
) {
  const params = curItem.params as Record<string, unknown>;
  const style = curItem.style as Record<string, unknown>;
  params.column = column;
  style.column = column;

  const config = {
    linePrice: 1,
    productName: 1,
    productPrice: 1,
    productSales: 0,
    product_btn: 1,
  };
  switch (column) {
    case 2: {
      config.productSales = 0;
      break;
    }
    case 3:
    case 4: {
      config.productSales = 0;
      config.product_btn = 0;
      break;
    }
  }
  Object.assign(curItem.params as Record<string, unknown>, config);
}

export function applyTitleStylePreset(
  curItem: Record<string, unknown>,
  type: number | string,
) {
  const params = curItem as Record<string, unknown>;
  const style = params.style as Record<string, unknown>;
  const presets: Record<number, Record<string, unknown>> = {
    1: {
      background: '#FFFFFF',
      bgcolor: '#FFFFFF',
      bottomRadio: 0,
      lineColor: '#FF0000',
      paddingBottom: 0,
      paddingLeft: 0,
      paddingTop: 0,
      textColor: '#FF0000',
      textSize: 20,
      topRadio: 0,
      weight: 800,
    },
    2: {
      background: '#FFFFFF',
      bgcolor: '#FFFFFF',
      bottomRadio: 0,
      lineColor: '#DDDDDD',
      paddingBottom: 10,
      paddingLeft: 0,
      paddingTop: 10,
      textColor: '#FF0000',
      textSize: 20,
      topRadio: 0,
      weight: 800,
    },
    3: {
      background: '#FFFFFF',
      bgcolor: '#FFFFFF',
      bottomRadio: 0,
      lineColor: '',
      paddingBottom: 10,
      paddingLeft: 0,
      paddingTop: 10,
      textColor: '#FF0000',
      textSize: 18,
      topRadio: 0,
      weight: 800,
    },
    4: {
      background: '#FFFFFF',
      bgcolor: '#FFFFFF',
      bottomRadio: 0,
      lineColor: '',
      paddingBottom: 10,
      paddingLeft: 0,
      paddingTop: 10,
      subtextColor: '#DDDDDD',
      subtextSize: 12,
      textColor: '#FF0000',
      textSize: 18,
      topRadio: 0,
      weight: 800,
    },
    5: {
      background: '#FFFFFF',
      bgcolor: '#FFFFFF',
      bottomRadio: 0,
      lineColor: '',
      paddingBottom: 10,
      paddingLeft: 0,
      paddingTop: 10,
      subtextColor: '#999999',
      subtextSize: 12,
      textColor: '#FF0000',
      textSize: 20,
      topRadio: 0,
      weight: 800,
    },
    6: {
      background: '#FFFFFF',
      bgcolor: '#FFFFFF',
      bottomRadio: 0,
      lineColor: '',
      paddingBottom: 10,
      paddingLeft: 0,
      paddingTop: 10,
      subtextColor: '#eeeeee',
      subtextSize: 18,
      textColor: '#FF0000',
      textSize: 20,
      topRadio: 0,
      weight: 800,
    },
    7: {
      background: '#FFFFFF',
      bgcolor: '#FFFFFF',
      bottomRadio: 0,
      lineColor: '',
      paddingBottom: 10,
      paddingLeft: 0,
      paddingTop: 10,
      subbackground: '#FFCCCC',
      subtextColor: '#FF0000',
      subtextSize: 14,
      textColor: '#FF0000',
      textSize: 20,
      topRadio: 0,
      weight: 800,
    },
    8: {
      background: '#FFFFFF',
      bgcolor: '#FFFFFF',
      bottomRadio: 0,
      isLine: 1,
      isMore: 1,
      isSub: 1,
      lineColor: '#FF0000',
      moretextColor: '#FF0000',
      paddingBottom: 10,
      paddingLeft: 0,
      paddingTop: 10,
      subbackground: '#FFCCCC',
      subtextColor: '#FF0000',
      subtextSize: 14,
      textColor: '#FF0000',
      textSize: 18,
      topRadio: 0,
      weight: 800,
    },
  };
  Object.assign(style, presets[Number(type)] ?? {});
  style.type = Number(type);
}
