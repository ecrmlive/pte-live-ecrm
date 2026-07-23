import type { RouteLocationNormalized, RouteLocationRaw } from 'vue-router';

import {
  LIST_MODAL_ROUTE_RULES,
  pathsForRule,
  type ListModalRouteRule,
} from './list-modal-route-rules';

export const LIST_MODAL_SESSION_KEY = 'qixiLiveMerchantListModal';

export type PlusProductPlugin = 'advance' | 'assemble' | 'bargain' | 'points';

export type RefundDetailModalMode = 'audit' | 'detail' | 'receipt';

export type ListModalIntent =
  | { action: 'live-add' }
  | { action: 'live-edit'; liveId: number }
  | { action: 'order-detail'; orderId: number }
  | { action: 'order-extract'; orderId: number }
  | { action: 'product-add' }
  | { action: 'product-edit'; productId: number; scene?: 'copy' | 'edit' }
  | { action: 'plus-product-add'; plugin: PlusProductPlugin; productId?: number }
  | { action: 'plus-product-edit'; id: number; plugin: PlusProductPlugin }
  | { action: 'refund-open'; mode: RefundDetailModalMode; orderRefundId: number }
  | { action: 'package-orders'; giftPackageId: number }
  | { action: 'invitation-partake'; invitationGiftId: number }
  | { action: 'list-edit'; listKey: string; id: number }
  | { action: 'list-action'; listKey: string; mode: string; id: number };

const PLUS_PRODUCT_FORM_ROUTES: Array<{
  idKey: string;
  listPath: string;
  plugin: PlusProductPlugin;
}> = [
  { idKey: 'assemble_product_id', listPath: '/plus/assemble/index', plugin: 'assemble' },
  { idKey: 'bargain_product_id', listPath: '/plus/bargain/index', plugin: 'bargain' },
  { idKey: 'advance_product_id', listPath: '/plus/advance/index', plugin: 'advance' },
  { idKey: 'point_product_id', listPath: '/plus/points/index', plugin: 'points' },
];

/** Persist modal intent across access-guard redirects that may drop query params. */
export function markListModalIntent(intent: ListModalIntent) {
  if (typeof sessionStorage === 'undefined') return;
  sessionStorage.setItem(LIST_MODAL_SESSION_KEY, JSON.stringify(intent));
}

export function peekListModalIntent(): ListModalIntent | null {
  if (typeof sessionStorage === 'undefined') return null;
  const raw = sessionStorage.getItem(LIST_MODAL_SESSION_KEY);
  if (!raw) return null;
  try {
    return JSON.parse(raw) as ListModalIntent;
  } catch {
    return null;
  }
}

export function clearListModalIntent() {
  if (typeof sessionStorage === 'undefined') return;
  sessionStorage.removeItem(LIST_MODAL_SESSION_KEY);
}

export function consumeListModalIntent(): ListModalIntent | null {
  const intent = peekListModalIntent();
  if (intent) clearListModalIntent();
  return intent;
}

function resolveIdFromQuery(
  query: RouteLocationNormalized['query'],
  idKey: string,
  altIdKeys: string[] = [],
) {
  const keys = [idKey, ...altIdKeys];
  for (const key of keys) {
    const id = Number(query[key] ?? 0);
    if (id > 0) return id;
  }
  return 0;
}

function redirectForListRule(
  rule: ListModalRouteRule,
  to: RouteLocationNormalized,
  markEdit: boolean,
) {
  const id = markEdit ? resolveIdFromQuery(to.query, rule.idKey, rule.altIdKeys) : 0;
  if (markEdit && id > 0) {
    markListModalIntent({ action: 'list-edit', listKey: rule.listKey, id });
  }
  return {
    path: rule.listPath,
    ...(rule.listQuery ? { query: { ...rule.listQuery } } : {}),
    replace: true as const,
  };
}

/** 列表页 onMounted：消费 edit 权限路由 intent */
export function consumeListEditIntent(
  listKey: string,
  onEdit: (id: number) => void | Promise<void>,
) {
  const intent = consumeListModalIntent();
  if (intent?.action === 'list-edit' && intent.listKey === listKey && intent.id > 0) {
    void onEdit(intent.id);
  }
}

/** 列表页 onMounted：消费 edit / 扩展 action（如 qrcode） */
export function consumeListRouteIntent(
  listKey: string,
  handlers: {
    edit?: (id: number) => void | Promise<void>;
    [mode: string]: ((id: number) => void | Promise<void>) | undefined;
  },
) {
  const intent = consumeListModalIntent();
  if (!intent) return;
  if (intent.action === 'list-edit' && intent.listKey === listKey && intent.id > 0) {
    void handlers.edit?.(intent.id);
    return;
  }
  if (intent.action === 'list-action' && intent.listKey === listKey && intent.id > 0) {
    const handler = handlers[intent.mode];
    if (handler) void handler(intent.id);
  }
}

function resolveConfiguredListModalRedirect(
  to: RouteLocationNormalized,
): RouteLocationRaw | null {
  const path = to.path;

  for (const rule of LIST_MODAL_ROUTE_RULES) {
    const paths = pathsForRule(rule);
    if (paths.add.some((p) => path === p) || paths.delete.some((p) => path === p)) {
      return redirectForListRule(rule, to, false);
    }
    if (paths.edit.some((p) => path === p)) {
      return redirectForListRule(rule, to, true);
    }
  }

  // 插件直播 · 房间二维码
  if (path.includes('/plus/live/room/qrcode')) {
    const liveId = resolveIdFromQuery(to.query, 'live_id', ['id']);
    if (liveId > 0) {
      markListModalIntent({
        action: 'list-action',
        listKey: 'plus-live-room',
        mode: 'qrcode',
        id: liveId,
      });
    }
    return {
      path: '/plus/live/wx/index',
      query: { type: 'room' },
      replace: true,
    };
  }

  return null;
}

/**
 * Permission-only routes (add/edit buttons) map to list pages + modal flags.
 * Redirect in router guard so tab bar URL stays on the list path.
 */
export function resolveListModalRedirect(
  to: RouteLocationNormalized,
): RouteLocationRaw | null {
  const path = to.path;

  if (path === '/live/room/add' || path.endsWith('/live/room/add')) {
    markListModalIntent({ action: 'live-add' });
    return {
      path: '/live/index',
      replace: true,
    };
  }

  if (path.includes('/live/room/edit')) {
    const liveId = Number(to.query.live_id ?? to.query.id ?? 0);
    if (liveId > 0) {
      markListModalIntent({ action: 'live-edit', liveId });
    }
    return {
      path: '/live/index',
      replace: true,
    };
  }

  if (path === '/order/order/detail' || path.endsWith('/order/order/detail')) {
    const orderId = Number(to.query.order_id ?? 0);
    if (orderId > 0) {
      markListModalIntent({ action: 'order-detail', orderId });
    }
    return {
      path: '/order/order/index',
      replace: true,
    };
  }

  if (path === '/order/operate/extract' || path.endsWith('/order/operate/extract')) {
    const orderId = Number(to.query.order_id ?? 0);
    if (orderId > 0) {
      markListModalIntent({ action: 'order-extract', orderId });
    }
    return {
      path: '/order/order/index',
      replace: true,
    };
  }

  if (path === '/product/product/add') {
    markListModalIntent({ action: 'product-add' });
    return {
      path: '/product/product/index',
      replace: true,
    };
  }

  if (path === '/product/product/edit' || path.endsWith('/product/product/edit')) {
    const productId = Number(to.query.product_id ?? 0);
    const scene = String(to.query.scene || 'edit') === 'copy' ? 'copy' : 'edit';
    if (productId > 0) {
      markListModalIntent({ action: 'product-edit', productId, scene });
    }
    return {
      path: '/product/product/index',
      replace: true,
    };
  }

  if (
    path === '/order/refund/detail' ||
    path === '/order/refund/Detail' ||
    path === '/order/refund/audit' ||
    path === '/order/refund/receipt'
  ) {
    const orderRefundId = Number(to.query.order_refund_id ?? 0);
    const mode: RefundDetailModalMode = path.includes('/audit')
      ? 'audit'
      : path.includes('/receipt')
        ? 'receipt'
        : 'detail';
    if (orderRefundId > 0) {
      markListModalIntent({ action: 'refund-open', mode, orderRefundId });
    }
    return {
      path: '/order/refund/index',
      replace: true,
    };
  }

  if (path === '/plus/package/orderlist' || path.endsWith('/plus/package/orderlist')) {
    const giftPackageId = Number(to.query.gift_package_id ?? 0);
    if (giftPackageId > 0) {
      markListModalIntent({ action: 'package-orders', giftPackageId });
    }
    return {
      path: '/plus/package/index',
      replace: true,
    };
  }

  if (
    path === '/plus/invitation/active/partake' ||
    path.endsWith('/plus/invitation/active/partake')
  ) {
    const invitationGiftId = Number(to.query.invitation_gift_id ?? 0);
    if (invitationGiftId > 0) {
      markListModalIntent({ action: 'invitation-partake', invitationGiftId });
    }
    return {
      path: '/plus/invitation/active/index',
      replace: true,
    };
  }

  const configured = resolveConfiguredListModalRedirect(to);
  if (configured) {
    return configured;
  }

  for (const { idKey, listPath, plugin } of PLUS_PRODUCT_FORM_ROUTES) {
    const productBase = `/plus/${plugin}/product`;
    if (path.includes(`${productBase}/add`)) {
      const productId = Number(to.query.product_id ?? 0);
      markListModalIntent({
        action: 'plus-product-add',
        plugin,
        ...(productId > 0 ? { productId } : {}),
      });
      return {
        path: listPath,
        query: { type: 'product' },
        replace: true,
      };
    }
    if (path.includes(`${productBase}/edit`)) {
      const id = Number(to.query[idKey] ?? 0);
      if (id > 0) {
        markListModalIntent({ action: 'plus-product-edit', id, plugin });
      }
      return {
        path: listPath,
        query: { type: 'product' },
        replace: true,
      };
    }
  }

  return null;
}

/** Open add/edit modal after guard redirected a permission-only plus product route. */
export function consumePlusProductModalIntent(
  plugin: PlusProductPlugin,
  handlers: {
    onAdd: (productId?: number) => void;
    onEdit: (id: number) => void;
  },
) {
  const intent = consumeListModalIntent();
  if (!intent) return;

  if (intent.action === 'plus-product-add' && intent.plugin === plugin) {
    handlers.onAdd(intent.productId);
    return;
  }

  if (intent.action === 'plus-product-edit' && intent.plugin === plugin && intent.id > 0) {
    handlers.onEdit(intent.id);
  }
}

/** Strip modal opener query keys after the list page consumes them. */
export function stripListModalQuery(query: Record<string, unknown>) {
  const next = { ...query };
  delete next.openAdd;
  delete next.openEdit;
  return next;
}
