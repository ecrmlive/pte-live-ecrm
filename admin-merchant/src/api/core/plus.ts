import { requestClient } from '#/api/request';

import type { PlusCategory } from '#/views/plus/types';

export async function getPlusCenterApi() {
  return requestClient.post<{ list: PlusCategory[] }>(
    '/shop/plus.plus/index',
    {},
  );
}
