import type { VxeTableGridOptions } from '@vben/plugins/vxe-table';

import type { ComponentPropsMap, ComponentType } from './component';

import { h } from 'vue';

import {
  setupVbenVxeTable,
  useVbenVxeGrid as useGrid,
} from '@vben/plugins/vxe-table';

import { ElButton, ElImage } from 'element-plus';

import {
  MERCHANT_LIST_GRID_CLASS,
  MERCHANT_LIST_GRID_LAYOUT,
  merchantListPagerConfig,
} from '#/constants/merchant-list-grid';

import { useVbenForm } from './form';

setupVbenVxeTable({
  configVxeTable: (vxeUI) => {
    vxeUI.setConfig({
      grid: {
        align: 'center',
        border: false,
        columnConfig: {
          resizable: true,
        },
        minHeight: 180,
        height: 'auto',
        formConfig: {
          // 全局禁用vxe-table的表单配置，使用formOptions
          enabled: false,
        },
        proxyConfig: {
          autoLoad: true,
          response: {
            result: 'items',
            total: 'total',
            list: 'items',
          },
          showActiveMsg: true,
          showResponseMsg: false,
        },
        round: true,
        // 默认 true（单行省略）；主列表页由 useVbenVxeGrid 合并 MERCHANT_LIST_GRID_LAYOUT 覆盖
        showOverflow: true,
        size: 'small',
      } as VxeTableGridOptions,
    });

    // 表格配置项可以用 cellRender: { name: 'CellImage' },
    vxeUI.renderer.add('CellImage', {
      renderTableDefault(renderOpts, params) {
        const { props } = renderOpts;
        const { column, row } = params;
        const src = row[column.field];
        return h(ElImage, { src, previewSrcList: [src], ...props });
      },
    });

    // 表格配置项可以用 cellRender: { name: 'CellLink' },
    vxeUI.renderer.add('CellLink', {
      renderTableDefault(renderOpts) {
        const { props } = renderOpts;
        return h(
          ElButton,
          { size: 'small', link: true },
          { default: () => props?.text },
        );
      },
    });

    // 这里可以自行扩展 vxe-table 的全局配置，比如自定义格式化
    // vxeUI.formats.add
  },
  useVbenForm,
});

type GridHookOptions<T extends Record<string, any> = Record<string, any>> =
  Parameters<typeof useGrid<T, ComponentType, ComponentPropsMap>>[0];

/** Disable virtual Y when height is auto — avoids vxe #7141 huge body-inner-wrapper. */
function withAutoHeightVirtualYGuard<T extends Record<string, any>>(
  options?: GridHookOptions<T>,
): GridHookOptions<T> | undefined {
  if (!options?.gridOptions) {
    return options;
  }

  const { gridOptions } = options;
  const height = gridOptions.height;
  if (height !== 'auto' && height !== undefined && height !== null && height !== '') {
    return options;
  }

  return {
    ...options,
    gridOptions: {
      ...gridOptions,
      virtualYConfig: {
        ...gridOptions.virtualYConfig,
        enabled: false,
      },
    },
  };
}

/** 主列表默认：merchant-vxe-grid + showOverflow false + 单元格顶对齐 + 真实分页 */
function withMerchantListGridDefaults<T extends Record<string, any>>(
  options?: GridHookOptions<T>,
): GridHookOptions<T> {
  const base = options ?? {};
  const gridOptions = base.gridOptions ?? {};

  return {
    ...base,
    gridClass: base.gridClass ?? MERCHANT_LIST_GRID_CLASS,
    gridOptions: {
      ...MERCHANT_LIST_GRID_LAYOUT,
      ...gridOptions,
      cellConfig: {
        ...MERCHANT_LIST_GRID_LAYOUT.cellConfig,
        ...gridOptions.cellConfig,
      },
      pagerConfig: {
        ...merchantListPagerConfig(),
        ...(gridOptions.pagerConfig ?? {}),
      },
    },
  };
}

export const useVbenVxeGrid = <T extends Record<string, any>>(
  options?: GridHookOptions<T>,
) => {
  const normalized = withMerchantListGridDefaults(options);
  return useGrid<T, ComponentType, ComponentPropsMap>(
    withAutoHeightVirtualYGuard(normalized) ?? normalized,
  );
};

export type * from '@vben/plugins/vxe-table';
