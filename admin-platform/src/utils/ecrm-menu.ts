import type { PlatformAccessMenuItem } from '#/utils/platform-menu';

import { normalizePlatformMenuIcon } from '#/constants/platform-lucide-icons';
import { resolveMergersComponent } from '#/views/ecrm/registry';

export interface MergersMenuNode {
  /** 旧菜单接口字段。 */
  menu_id?: number;
  pid?: number;
  path?: string;
  icon?: string;
  menu_name?: string;
  route?: string;
  /** qixi_crm_a_menu / api-platform 字段。 */
  id?: number;
  parent_id?: number;
  route_path?: string;
  title?: string;
  code?: string;
  kind?: 'button' | 'directory' | 'page' | string;
  sort?: number;
  children?: MergersMenuNode[];
}

function nodeID(node: MergersMenuNode) {
  return Number(node.id ?? node.menu_id ?? 0);
}

function nodeParentID(node: MergersMenuNode) {
  return Number(node.parent_id ?? node.pid ?? 0);
}

function nodePath(node: MergersMenuNode) {
  return String(node.route_path ?? node.path ?? '').trim() || `/menu-${nodeID(node)}`;
}

function nodeTitle(node: MergersMenuNode) {
  return String(node.title ?? node.menu_name ?? node.code ?? nodePath(node));
}

/**
 * 初始化数据已为所有一级模块配置图标；这里保留按业务代码的兜底，防止老库或
 * 运营新建菜单没有填图标时退化为 Vben 的 circle-dot 占位图标。
 */
const MENU_ICON_FALLBACK: Record<string, string> = {
  dashboard: 'layout-dashboard',
  merchant: 'store',
  region: 'map-pinned',
  service: 'messages-square',
  product: 'puzzle',
  order: 'receipt-text',
  marketing: 'activity',
  user: 'users',
  content: 'images',
  freight: 'map-plus',
  accounts: 'wallet',
  operations: 'pen-line',
  setting: 'settings',
};

function nodeIcon(node: MergersMenuNode, path: string) {
  const explicit = normalizePlatformMenuIcon(node.icon);
  if (explicit) {
    return explicit;
  }
  const code = String(node.code || '').split('.')[0];
  const segment = path.split('/').filter(Boolean)[0];
  return normalizePlatformMenuIcon(MENU_ICON_FALLBACK[code] || MENU_ICON_FALLBACK[segment]);
}

function mapNode(node: MergersMenuNode): PlatformAccessMenuItem {
  const path = nodePath(node);
  const children = (node.children || []).map(mapNode);
  const component = resolveMergersComponent(path);
  const isLeaf = children.length === 0;
  // qixi_crm_a_menu 的 directory 仅用于组织侧栏，不应注册成路由。
  // 否则父目录若碰巧与子页路径同名（例如 /region）会被 Vben 当作
  // 叶子页注册，造成子菜单脱离父级、落到侧栏底部。
  const isDirectory = node.kind === 'directory' || !isLeaf;
  return {
    access_id: nodeID(node),
    name: nodeTitle(node),
    path,
    icon: nodeIcon(node, path),
    is_menu: 1,
    is_route: !isDirectory && (component || isLeaf) ? 1 : 0,
    is_show: 1,
    component: !isDirectory
      ? component || (isLeaf ? 'ecrm/placeholder/index' : undefined)
      : undefined,
    children: children.length ? children : undefined,
  };
}

export function mapMergersMenusToAccess(
  menus: unknown[],
): PlatformAccessMenuItem[] {
  const flat = (menus as MergersMenuNode[]).filter(
    (item) => item.kind !== 'button',
  );
  // api-platform 返回平铺的 qixi_crm_a_menu；旧接口可能已经带 children。
  // 这里统一建树，避免把每个叶子错误注册到 /menu-0。
  const byID = new Map<number, MergersMenuNode>();
  for (const item of flat) {
    const id = nodeID(item);
    if (id > 0) {
      byID.set(id, { ...item, children: [] });
    }
  }
  const roots: MergersMenuNode[] = [];
  for (const item of byID.values()) {
    const parent = byID.get(nodeParentID(item));
    if (parent && parent !== item) {
      parent.children?.push(item);
    } else {
      roots.push(item);
    }
  }
  roots.sort((left, right) => Number(left.sort || 0) - Number(right.sort || 0));
  return roots.map(mapNode);
}
