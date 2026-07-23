import { fetchPlatformSessionApi } from '#/api/core/platform-session';
import {
  convertPlatformMenusToVben,
  type PlatformAccessMenuItem,
} from '#/utils/platform-menu';

export async function getAllMenusApi() {
  const session = await fetchPlatformSessionApi();
  return convertPlatformMenusToVben(session.menus || []);
}

export type { PlatformAccessMenuItem };
