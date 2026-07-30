import type { RouteRecordStringComponent } from '@vben/types';

import { resolveMergersComponent } from '#/views/mergers/registry';

export interface MergersMenuNode {
  /** 旧身份服务菜单字段。 */
  menu_id?: number;
  pid?: number;
  path?: string;
  icon?: string;
  menu_name?: string;
  route?: string;
  sort?: number;
  /** qixi_crm_m_menu / api-merchant 字段。 */
  id?: number;
  parent_id?: number;
  name?: string;
  component?: string;
  is_menu?: number;
  is_route?: number;
  children?: MergersMenuNode[];
}

function nodeID(node: MergersMenuNode) {
  return Number(node.id ?? node.menu_id ?? 0);
}

function nodeParentID(node: MergersMenuNode) {
  return Number(node.parent_id ?? node.pid ?? 0);
}

function nodePath(node: MergersMenuNode) {
  return String(node.path || '').trim() || `/menu-${nodeID(node)}`;
}

function nodeTitle(node: MergersMenuNode) {
  return String(node.name ?? node.menu_name ?? nodePath(node));
}

function flattenMenuTree(nodes: MergersMenuNode[]) {
  const flat: MergersMenuNode[] = [];
  const visit = (node: MergersMenuNode) => {
    flat.push(node);
    for (const child of node.children || []) {
      visit(child);
    }
  };
  for (const node of nodes) {
    visit(node);
  }
  return flat;
}

/**
 * 只接受仓库内 views/mergers 的声明式组件路径，避免把数据库字符串当成任意动态导入路径。
 * 已注册路径优先由 registry 决定；此处用于新种子菜单无需再维护第二份映射。
 */
function declaredComponent(node: MergersMenuNode) {
  const component = String(node.component || '')
    .trim()
    .replace(/^views\//, '')
    .replace(/\.vue$/, '');
  return /^mergers\/[a-zA-Z0-9_/-]+$/.test(component)
    ? component
    : undefined;
}

function routeNameFromPath(path: string) {
  return (
    String(path || '')
      .split('/')
      .filter(Boolean)
      .map((segment) => segment.replace(/[^a-zA-Z0-9]/g, ''))
      .join('') || 'Root'
  );
}

function convertNode(node: MergersMenuNode): RouteRecordStringComponent | null {
  const path = nodePath(node);
  const children = (node.children || [])
    .map(convertNode)
    .filter(Boolean) as RouteRecordStringComponent[];
  const componentKey = resolveMergersComponent(path) || declaredComponent(node);
  const isLeaf = children.length === 0;
  const record: RouteRecordStringComponent = {
    name: routeNameFromPath(path),
    path,
    meta: {
      title: nodeTitle(node),
      icon: node.icon || undefined,
      hideInMenu: false,
    },
  };
  if (componentKey || isLeaf) {
    record.component = `../views/${componentKey || 'mergers/placeholder/index'}.vue`;
  } else if (children[0]?.path) {
    record.redirect = children[0].path;
  }
  if (children.length) {
    record.children = children;
  }
  return record;
}

export function convertMergersMenusToVben(
  menus: unknown[],
): RouteRecordStringComponent[] {
  const input = menus as MergersMenuNode[];
  // api-merchant 与旧接口都可能返回 children；也允许管理接口传入平铺行。
  // 统一建树，避免新菜单被错误解析为 /menu-0 或丢失父子关系。
  const byID = new Map<number, MergersMenuNode>();
  for (const item of flattenMenuTree(input)) {
    const id = nodeID(item);
    if (id > 0) {
      byID.set(id, { ...item, children: [] });
    }
  }
  if (!byID.size) {
    return input
      .map(convertNode)
      .filter(Boolean) as RouteRecordStringComponent[];
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
    .map(convertNode)
    .filter(Boolean) as RouteRecordStringComponent[];
}

export function extractRouterRoutes(
  routes: RouteRecordStringComponent[],
): RouteRecordStringComponent[] {
  const pages: RouteRecordStringComponent[] = [];
  function walk(route: RouteRecordStringComponent) {
    if (route.component) {
      pages.push({ ...route, children: undefined });
      return;
    }
    for (const child of route.children ?? []) {
      walk(child);
    }
  }
  for (const route of routes) {
    walk(route);
  }
  return pages;
}
