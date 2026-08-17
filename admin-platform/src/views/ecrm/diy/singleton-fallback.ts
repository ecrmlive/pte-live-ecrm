import type {
  DiyEditorBootstrap,
  DiyPageDoc,
  DiySingletonScope,
} from '#/api/core/diy';

type DiyItem = Record<string, unknown>;

const draftStorageKey = (scope: DiySingletonScope) =>
  `ecrm:diy:singleton-draft:${scope}`;

const clone = <T>(value: T): T => JSON.parse(JSON.stringify(value));

const baseStyle = {
  background: 'rgba(255, 255, 255, 0)',
  bottomBackground: 'rgba(255, 255, 255, 0)',
  marginTop: 0,
  paddingBottom: 0,
  paddingLeft: 0,
  paddingRight: 0,
  paddingTop: 0,
  radius: 0,
  shadow: false,
};

function createItem(
  name: string,
  type: string,
  group: string,
  icon: string,
  options: Partial<DiyItem> = {},
): DiyItem {
  return {
    name,
    type,
    group,
    icon,
    params: {},
    style: clone(baseStyle),
    ...options,
  };
}

function createDefaultData(): Record<string, DiyItem> {
  return {
    search: createItem('搜索框', 'search', 'tools', 'icon-sousuo', {
      params: {
        hotWords: ['雅诗兰黛', 'Only澳白瓶'],
        location: '定位中',
        placeholder: '搜索商品',
        searchType: 'search',
      },
    }),
    banner: createItem('图片轮播', 'banner', 'media', 'icon-tupian', {
      data: [{ imgUrl: '', linkUrl: '', title: '轮播图片' }],
      params: { autoplay: true, interval: 3 },
    }),
    imageSingle: createItem('单图组', 'imageSingle', 'media', 'icon-tupian', {
      data: [{ imgUrl: '', linkUrl: '', title: '图片' }],
    }),
    navBar: createItem('导航组', 'navBar', 'tools', 'icon-yingyong', {
      data: [
        { hide: false, icon: '店', iconBg: '#e8f1ff', iconColor: '#3777ff', text: '品牌好店' },
        { hide: false, icon: '驻', iconBg: '#e8f8ed', iconColor: '#35a863', text: '商户入驻' },
        { hide: false, icon: '分', iconBg: '#fff3d7', iconColor: '#e39a00', text: '积分商城' },
        { hide: false, icon: '券', iconBg: '#fff0f2', iconColor: '#e85d79', text: '领券中心' },
        { hide: false, icon: '会', iconBg: '#f2eaff', iconColor: '#8863e8', text: '会员中心' },
      ],
      style: { ...baseStyle, background: '#ffffff', textColor: '#333333' },
    }),
    option: createItem('选项卡', 'option', 'tools', 'icon-xuanxiang', {
      data: [{ text: '首页' }, { text: '果蔬生鲜' }, { text: '健康医疗' }],
    }),
    title: createItem('标题', 'title', 'tools', 'icon-biaoti', {
      params: { title: '购物车' },
      style: { ...baseStyle, background: '#ffffff', paddingLeft: 12, paddingRight: 12 },
    }),
    product: createItem('商品组', 'product', 'shop', 'icon-shangpin', {
      params: {
        cartText: '加入购物车',
        cartType: 'icon',
        column: 2,
        comment: false,
        linePrice: true,
        productName: true,
        productPrice: true,
        productSales: false,
        showCart: true,
      },
      style: {
        ...baseStyle,
        background: '#f5f5f5',
        bottomRadio: 10,
        cart_color1: '#ff4c01',
        cart_color2: '#ff7a45',
        cart_text_color: '#ffffff',
        line_price_color: '#999999',
        paddingBottom: 10,
        paddingLeft: 10,
        paddingRight: 10,
        paddingTop: 10,
        productBottomRadio: 8,
        productTopRadio: 8,
        product_name_color: '#333333',
        product_price_color: '#ff4c01',
        topRadio: 10,
      },
    }),
    store: createItem('店铺组', 'store', 'shop', 'icon-dianpu', {
      params: { title: '品牌好店' },
      style: { ...baseStyle, background: '#ffffff', paddingLeft: 10, paddingRight: 10 },
    }),
    coupon: createItem('优惠券', 'coupon', 'marketing', 'icon-youhuiquan'),
    service: createItem('在线客服', 'service', 'tools', 'icon-kefu'),
    richText: createItem('富文本', 'richText', 'tools', 'icon-wenzhang'),
    video: createItem('视频组', 'video', 'media', 'icon-shipin'),
    blank: createItem('辅助空白', 'blank', 'tools', 'icon-kongbai', {
      params: { height: 20 },
    }),
    bottomNav: createItem('底部导航', 'bottomNav', 'tools', 'icon-daohang', {
      data: [
        { hide: false, icon: '⌂', linkUrl: '/', text: '首页' },
        { hide: false, icon: '▦', linkUrl: '/pages/goods_cate/goods_cate', text: '分类' },
        { hide: false, icon: '🛒', linkUrl: '/pages/order_addcart/order_addcart', text: '购物车' },
        { hide: false, icon: '◉', linkUrl: '/pages/user/index', text: '我的' },
      ],
      style: {
        ...baseStyle,
        background: '#ffffff',
        positionType: 'fixed',
        textColor: '#333333',
      },
    }),
  };
}

function itemFromDefault(defaultData: Record<string, DiyItem>, key: string) {
  return clone(defaultData[key]);
}

function createDefaultPage(scope: DiySingletonScope, defaultData: Record<string, DiyItem>): DiyPageDoc {
  const keysByScope: Record<DiySingletonScope, string[]> = {
    cart: ['search', 'title', 'product', 'bottomNav'],
    home: ['search', 'banner', 'navBar', 'product', 'bottomNav'],
    store: ['search', 'banner', 'store', 'product', 'bottomNav'],
  };
  return {
    items: keysByScope[scope].map((key) => itemFromDefault(defaultData, key)),
    page: {
      background: '#f5f5f5',
      title: scope === 'cart' ? '购物车' : scope === 'store' ? '店铺' : '首页',
    },
  };
}

export function readSingletonDraft(scope: DiySingletonScope): DiyPageDoc | null {
  if (typeof window === 'undefined') return null;
  try {
    const value = window.localStorage.getItem(draftStorageKey(scope));
    if (!value) return null;
    const draft = JSON.parse(value) as DiyPageDoc;
    if (!Array.isArray(draft?.items) || !draft.page || typeof draft.page !== 'object') {
      return null;
    }
    return draft;
  } catch {
    return null;
  }
}

export function saveSingletonFallbackDraft(scope: DiySingletonScope, page: DiyPageDoc) {
  if (typeof window === 'undefined') return;
  try {
    window.localStorage.setItem(draftStorageKey(scope), JSON.stringify(page));
  } catch {
    // 浏览器禁用本地存储时保持当前编辑态，不阻断用户操作。
  }
}

export function createSingletonFallbackBootstrap(scope: DiySingletonScope): DiyEditorBootstrap {
  const defaultData = createDefaultData();
  return {
    defaultData,
    jsonData: readSingletonDraft(scope) ?? createDefaultPage(scope, defaultData),
    opts: {},
    pageId: 0,
    scope,
  };
}
