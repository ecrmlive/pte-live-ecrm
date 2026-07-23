import type { RouteRecordStringComponent } from '@vben/types';

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
  const path = String(node.path || '').trim() || `/menu-${node.menu_id}`;
  const children = (node.children || [])
    .map(convertNode)
    .filter(Boolean) as RouteRecordStringComponent[];
  const componentKey = resolveMergersComponent(path);
  const isLeaf = children.length === 0;
  const record: RouteRecordStringComponent = {
    name: routeNameFromPath(path),
    path,
    meta: {
      title: node.menu_name || path,
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
  return (menus as MergersMenuNode[])
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
