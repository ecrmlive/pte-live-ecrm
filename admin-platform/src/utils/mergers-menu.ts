import type { PlatformAccessMenuItem } from '#/utils/platform-menu';

import { resolveMergersComponent } from '#/views/mergers/registry';

export interface MergersMenuNode {
  menu_id: number;
  pid: number;
  path: string;
  icon: string;
  menu_name: string;
  route: string;
  sort: number;
  children?: MergersMenuNode[];
}

function mapNode(node: MergersMenuNode): PlatformAccessMenuItem {
  const path = String(node.path || '').trim() || `/menu-${node.menu_id}`;
  const children = (node.children || []).map(mapNode);
  const component = resolveMergersComponent(path);
  const isLeaf = children.length === 0;
  return {
    access_id: node.menu_id,
    name: node.menu_name || path,
    path,
    icon: node.icon || undefined,
    is_menu: 1,
    is_route: component || isLeaf ? 1 : 0,
    is_show: 1,
    component: component || (isLeaf ? 'mergers/placeholder/index' : undefined),
    children: children.length ? children : undefined,
  };
}

export function mapMergersMenusToAccess(
  menus: unknown[],
): PlatformAccessMenuItem[] {
  return (menus as MergersMenuNode[]).map(mapNode);
}
