import { requestClient } from '#/api/request';

export interface LiveSensitiveWordItem {
  word: string;
  word_id: number;
}

export async function getLiveSensitiveWordListApi() {
  return requestClient.post<{ list: LiveSensitiveWordItem[] }>(
    '/api/v1/shop/live/sensitive-word/list',
    {},
  );
}

export async function addLiveSensitiveWordApi(payload: { word: string }) {
  return requestClient.post<LiveSensitiveWordItem>(
    '/api/v1/shop/live/sensitive-word/add',
    payload,
  );
}

export async function deleteLiveSensitiveWordApi(payload: { word_id: number }) {
  return requestClient.post<{ msg?: string }>(
    '/api/v1/shop/live/sensitive-word/delete',
    payload,
  );
}
