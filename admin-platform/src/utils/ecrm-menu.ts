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
  home: 'ant-design:home-outlined',
  dashboard: 'ant-design:dashboard-outlined',
  store: 'ant-design:shop-outlined',
  merchant: 'ant-design:shop-outlined',
  region: 'ant-design:cluster-outlined',
  service: 'ant-design:customer-service-outlined',
  product: 'ant-design:shopping-outlined',
  order: 'ant-design:file-text-outlined',
  promoter: 'ant-design:send-outlined',
  marketing: 'ant-design:flag-outlined',
  user: 'ant-design:user-outlined',
  content: 'ant-design:read-outlined',
  freight: 'lucide:map-plus',
  accounts: 'ant-design:bar-chart-outlined',
  app: 'ant-design:appstore-outlined',
  operations: 'ant-design:format-painter-outlined',
  setting: 'ant-design:setting-outlined',
  maintain: 'ant-design:tool-outlined',
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

function mapNode(node: MergersMenuNode): null | PlatformAccessMenuItem {
  const path = nodePath(node);
  const children = (node.children || [])
    .map(mapNode)
    .filter((item): item is PlatformAccessMenuItem => Boolean(item));
  const component = resolveMergersComponent(path);
  const isLeaf = children.length === 0;
  // qixi_crm_a_menu 的 directory 仅用于组织侧栏，不应注册成路由。
  // 否则父目录若碰巧与子页路径同名（例如 /region、/banlace、/brokerage）
  // 会被 Vben 当作叶子页注册，造成子菜单脱离父级或同级缩进错乱。
  const isDirectory = node.kind === 'directory' || !isLeaf;
  // 服务端菜单不是前端代码生成器。未注册真实组件的叶子不能降级到
  // placeholder，否则会把“菜单可见”误报成“功能已实现”。目录在所有
  // 子项都被过滤后同样隐藏，避免产生空的侧栏分组。
  if ((isLeaf && !component) || (isDirectory && children.length === 0)) {
    return null;
  }
  return {
    access_id: nodeID(node),
    name: nodeTitle(node),
    path,
    icon: nodeIcon(node, path),
    is_menu: 1,
    is_route: !isDirectory && component ? 1 : 0,
    is_show: 1,
    component: !isDirectory ? component : undefined,
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
  return roots
    .map(mapNode)
    .filter((item): item is PlatformAccessMenuItem => Boolean(item));
}
