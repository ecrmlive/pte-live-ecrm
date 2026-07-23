/** 列表页权限路由（add/edit/delete）→ 守卫静默回 listPath；edit 带 id 时 list-edit intent */

export type ListModalRouteRule = {
  listKey: string;
  listPath: string;
  idKey: string;
  /** 编辑页 query 备用 key（如 live_id / id） */
  altIdKeys?: string[];
  addPaths?: string[];
  editPaths?: string[];
  deletePaths?: string[];
  listQuery?: Record<string, string>;
};

function listBase(listPath: string) {
  return listPath.replace(/\/(index|list)$/, '');
}

export function pathsForRule(rule: ListModalRouteRule) {
  const base = listBase(rule.listPath);
  return {
    add: rule.addPaths ?? [`${base}/add`],
    delete: rule.deletePaths ?? [`${base}/delete`],
    edit: rule.editPaths ?? [`${base}/edit`],
  };
}

/** 权限 → 列表弹窗路由表（add/delete 静默；edit 打开编辑弹窗） */
export const LIST_MODAL_ROUTE_RULES: ListModalRouteRule[] = [
  // 门店
  { listKey: 'store-store', listPath: '/store/store/index', idKey: 'store_id' },
  { listKey: 'store-clerk', listPath: '/store/clerk/index', idKey: 'clerk_id' },
  // 权限
  { listKey: 'auth-user', listPath: '/auth/user/index', idKey: 'shop_user_id' },
  { listKey: 'auth-role', listPath: '/auth/role/index', idKey: 'role_id' },
  // 会员
  { listKey: 'user-member', listPath: '/user/user/index', idKey: 'user_id' },
  { listKey: 'user-tag', listPath: '/user/tag/index', idKey: 'tag_id' },
  { listKey: 'user-equity', listPath: '/user/equity/index', idKey: 'equity_id' },
  { listKey: 'user-grade', listPath: '/user/grade/index', idKey: 'grade_id' },
  // 设置
  { listKey: 'setting-address', listPath: '/setting/address/index', idKey: 'address_id' },
  { listKey: 'setting-delivery', listPath: '/setting/delivery/index', idKey: 'delivery_id' },
  { listKey: 'setting-express', listPath: '/setting/express/index', idKey: 'express_id' },
  { listKey: 'setting-printer', listPath: '/setting/printer/index', idKey: 'printer_id' },
  {
    listKey: 'live-h5domain',
    listPath: '/live/h5domain/index',
    idKey: 'domain_id',
    editPaths: ['/live/h5domain/edit'],
    deletePaths: ['/live/h5domain/delete', '/live/h5domain/status'],
  },
  // 插件 · 独立列表
  { listKey: 'plus-fullreduce', listPath: '/plus/fullreduce/index', idKey: 'fullreduce_id' },
  { listKey: 'plus-buyactivity', listPath: '/plus/buyactivity/index', idKey: 'buy_id' },
  { listKey: 'plus-package', listPath: '/plus/package/index', idKey: 'gift_package_id' },
  {
    listKey: 'plus-invitation',
    listPath: '/plus/invitation/active/index',
    idKey: 'invitation_gift_id',
  },
  {
    listKey: 'plus-coupon',
    listPath: '/plus/coupon/index',
    idKey: 'coupon_id',
    addPaths: ['/plus/coupon/coupon/add', '/plus/coupon/add'],
    editPaths: ['/plus/coupon/coupon/edit', '/plus/coupon/edit'],
    listQuery: { type: 'list' },
  },
  {
    listKey: 'plus-seckill-active',
    listPath: '/plus/seckill/index',
    idKey: 'seckill_activity_id',
    addPaths: ['/plus/seckill/active/add'],
    editPaths: ['/plus/seckill/active/edit'],
    listQuery: { type: 'first' },
  },
  {
    listKey: 'plus-seckill-time',
    listPath: '/plus/seckill/index',
    idKey: 'id',
    addPaths: ['/plus/seckill/time/add'],
    editPaths: ['/plus/seckill/time/edit'],
    listQuery: { type: 'third' },
  },
  {
    listKey: 'plus-article',
    listPath: '/plus/article/index',
    idKey: 'article_id',
    addPaths: ['/plus/article/article/add'],
    editPaths: ['/plus/article/article/edit'],
    listQuery: { type: 'article' },
  },
  {
    listKey: 'plus-surface-template',
    listPath: '/plus/surface/index',
    idKey: 'template_id',
    addPaths: ['/plus/surface/template/add'],
    editPaths: ['/plus/surface/template/edit'],
    listQuery: { type: 'template' },
  },
  {
    listKey: 'plus-surface-setting',
    listPath: '/plus/surface/index',
    idKey: 'setting_id',
    addPaths: ['/plus/surface/setting/add'],
    editPaths: ['/plus/surface/setting/edit'],
    listQuery: { type: 'setting' },
  },
  {
    listKey: 'plus-table-table',
    listPath: '/plus/table/event',
    idKey: 'table_id',
    addPaths: ['/plus/table/table/add'],
    editPaths: ['/plus/table/table/edit'],
    listQuery: { type: 'table' },
  },
  {
    listKey: 'plus-agent-grade',
    listPath: '/plus/agent/index',
    idKey: 'grade_id',
    addPaths: ['/plus/agent/grade/add'],
    editPaths: ['/plus/agent/grade/edit'],
    deletePaths: ['/plus/agent/grade/delete'],
    listQuery: { type: 'grade' },
  },
  {
    listKey: 'plus-agent-user',
    listPath: '/plus/agent/index',
    idKey: 'user_id',
    addPaths: ['/plus/agent/user/add'],
    editPaths: ['/plus/agent/user/edit'],
    deletePaths: ['/plus/agent/user/delete'],
    listQuery: { type: 'user' },
  },
  {
    listKey: 'plus-agent-poster',
    listPath: '/plus/agent/index',
    idKey: 'poster_id',
    addPaths: ['/plus/agent/poster/add'],
    editPaths: ['/plus/agent/poster/edit'],
    listQuery: { type: 'poster' },
  },
  {
    listKey: 'plus-live-room',
    listPath: '/plus/live/wx/index',
    idKey: 'live_id',
    altIdKeys: ['id'],
    addPaths: ['/plus/live/wx/add'],
    editPaths: ['/plus/live/wx/edit'],
    listQuery: { type: 'room' },
  },
  {
    listKey: 'plus-live-anchor',
    listPath: '/plus/live/wx/index',
    idKey: 'anchor_id',
    addPaths: ['/plus/live/anchor/add'],
    editPaths: ['/plus/live/anchor/edit', '/live/anchor/edit'],
    listQuery: { type: 'anchor' },
  },
  {
    listKey: 'plus-card-code',
    listPath: '/plus/card/event',
    idKey: 'code_id',
    editPaths: ['/plus/card/code/edit'],
    listQuery: { type: 'code' },
  },
];
