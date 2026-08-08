import type { SubMenuProvider } from '../types';

import { computed, getCurrentInstance } from 'vue';

import { findComponentUpward } from '../utils';

function useMenu() {
  const instance = getCurrentInstance();
  if (!instance) {
    throw new Error('instance is required');
  }

  /**
   * @zh_CN 获取所有父级菜单链路
   */
  const parentPaths = computed(() => {
    let parent = instance.parent;
    const paths: string[] = [instance.props.path as string];
    while (parent?.type.name !== 'MenuUI') {
      if (parent?.props.path) {
        paths.unshift(parent.props.path as string);
      }
      parent = parent?.parent ?? null;
    }

    return paths;
  });

  const parentMenu = computed(() => {
    return findComponentUpward(instance, ['MenuUI', 'SubMenu']);
  });

  return {
    parentMenu,
    parentPaths,
  };
}

function useMenuStyle(menu?: SubMenuProvider) {
  const subMenuStyle = computed(() => {
    // Parent SubMenu.level：0 表示当前 ul 内是 L2；1→L3；2→L4。
    // 注意：`?? 0 + 1` 会因运算符优先级变成 `?? 1`，导致层级错乱。
    return {
      '--menu-level': menu ? (menu.level ?? 0) : 0,
    };
  });
  return subMenuStyle;
}

export { useMenu, useMenuStyle };
