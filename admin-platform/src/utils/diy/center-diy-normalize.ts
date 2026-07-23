/** Legacy center pages may store product layout flags under style instead of params. */
export function normalizeCenterDiyItems(
  diyData: { items?: Array<Record<string, unknown>> },
) {
  for (const item of diyData.items ?? []) {
    if (item.type !== 'product') continue;

    const params = (item.params ??= {}) as Record<string, unknown>;
    const style = (item.style ?? {}) as Record<string, unknown>;
    const show = (style.show ?? {}) as Record<string, unknown>;

    if (params.column == null && style.column != null) {
      params.column = style.column;
    } else if (params.column != null && style.column == null) {
      style.column = params.column;
    }
    if (params.display == null && style.display != null) {
      params.display = style.display;
    }

    for (const key of [
      'productName',
      'productPrice',
      'linePrice',
      'productSales',
      'comment',
      'showCart',
      'cartType',
      'cartText',
      'sellingPoint',
    ]) {
      if (params[key] == null && show[key] != null) {
        params[key] = show[key];
      }
    }

    const colorDefaults: Record<string, string> = {
      product_name_color: '#333333',
      product_price_color: '#FF4C01',
      line_price_color: '#999999',
      product_sales_color: '#999999',
      product_comment_color: '#999999',
    };
    for (const [key, value] of Object.entries(colorDefaults)) {
      if (style[key] == null || style[key] === '') {
        style[key] = value;
      }
    }
  }
}
