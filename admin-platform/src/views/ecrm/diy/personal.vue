<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';

import { Icon as IconifyIcon } from '@iconify/vue';
import { Page } from '@vben/common-ui';
import {
  ElButton,
  ElColorPicker,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElRadio,
  ElRadioGroup,
  ElSlider,
  ElSwitch,
} from 'element-plus';

import { getPersonalDecorationApi, savePersonalDecorationApi } from '#/api/core/diy';
import ImagePickerDialog from '#/components/shop/image-picker-dialog.vue';

type SectionKey = 'header' | 'orders' | 'points' | 'services' | 'bottom';
type ActiveSectionKey = SectionKey | 'generic';
type ComponentStyle = {
  cardBackground: string;
  cardBackgroundEnd: string;
  cardBorderColor: string;
  cardBorderWidth: number;
  cardFloat: number;
  cardMarginTop: number;
  cardPaddingBottom: number;
  cardPaddingHorizontal: number;
  cardPaddingTop: number;
  cardRadius: number;
  cardShadow: 'off' | 'on';
  cardShadowBlur: number;
  cardShadowColor: string;
  cardShadowSpread: number;
  cardShadowX: number;
  cardShadowY: number;
  contentBackground: string;
  contentBorderColor: string;
  contentBorderWidth: number;
  contentFontSize: number;
  contentPaddingBottom: number;
  contentPaddingHorizontal: number;
  contentPaddingTop: number;
  contentRadius: number;
  contentShadow: 'off' | 'on';
  contentShadowBlur: number;
  contentShadowColor: string;
  contentShadowSpread: number;
  contentShadowX: number;
  contentShadowY: number;
  contentTextColor: string;
  headerBackground: string;
  headerBorderColor: string;
  headerBorderWidth: number;
  headerFontSize: number;
  headerPaddingBottom: number;
  headerPaddingTop: number;
  headerTextColor: string;
};
type NavigationItem = {
  enabled: boolean;
  icon: string;
  image?: string;
  label: string;
  link: string;
  value: number;
};
type HeaderStatistic = {
  enabled: boolean;
  label: string;
  value: number;
};
type OrderItem = {
  enabled: boolean;
  icon: string;
  label: string;
};
type LibraryItem = {
  icon: string;
  label: string;
  section?: SectionKey;
  type: string;
};
type GenericComponent = Pick<LibraryItem, 'icon' | 'label' | 'type'> & {
  bannerImages?: Array<{ image: string; link: string }>;
  description: string;
  enabled: boolean;
  id: string;
  link: string;
  productColumns?: 1 | 2 | 3;
  productCount?: number;
  productLayout?: 'big' | 'list' | 'waterfall';
  productSource?: 'all' | 'hot' | 'new';
  style: ComponentStyle;
  title: string;
};
type AdSlide = { actionText: string; image: string; link: string; subtitle: string; title: string };

function createComponentStyle(): ComponentStyle {
  return {
    // 默认不覆盖每个模块本身的初始视觉；用户修改后才应用对应层级的样式。
    headerBackground: 'transparent', headerTextColor: '#303133', headerFontSize: 16, headerPaddingTop: 0, headerPaddingBottom: 0, headerBorderColor: '#e7e9ee', headerBorderWidth: 0,
    contentBackground: 'transparent', contentTextColor: '#333333', contentFontSize: 14, contentPaddingTop: 0, contentPaddingBottom: 0, contentPaddingHorizontal: 0, contentBorderColor: '#eef0f4', contentBorderWidth: 0, contentRadius: 0, contentShadow: 'off', contentShadowColor: '#888888', contentShadowX: 0, contentShadowY: 0, contentShadowBlur: 0, contentShadowSpread: 0,
    cardBackground: 'transparent', cardBackgroundEnd: 'transparent', cardBorderColor: 'transparent', cardBorderWidth: 0, cardFloat: 0, cardMarginTop: 0, cardPaddingTop: 0, cardPaddingBottom: 0, cardPaddingHorizontal: 0, cardRadius: 0, cardShadow: 'off', cardShadowColor: '#888888', cardShadowX: 0, cardShadowY: 0, cardShadowBlur: 0, cardShadowSpread: 0,
  };
}

const saving = ref(false);
const activeSection = ref<ActiveSectionKey>('header');
const activeConfigTab = ref<'content' | 'style'>('content');
const activeGenericId = ref<string>();
const configScrollRef = ref<HTMLElement>();
const draggingService = ref<number>();
const draggingBottom = ref<number>();
const draggingOrder = ref<number>();
const imagePickerOpen = ref(false);
const imagePickerTarget = ref<'headerAvatar' | 'pointsSlide' | 'serviceIcon' | 'genericBanner'>('headerAvatar');
const imagePickerServiceIndex = ref<number>();
const imagePickerGenericId = ref<string>();
const imagePickerBannerIndex = ref<number>();
const imagePickerPointsIndex = ref<number>();
const activePointsSlide = ref(0);
const personalDraftStorageKey = 'ecrm:diy:personal-decoration:draft';
const genericComponents = ref<GenericComponent[]>([]);

function openImagePicker(target: 'headerAvatar' | 'pointsSlide' | 'serviceIcon' | 'genericBanner', serviceIndex?: number, genericId?: string, bannerIndex?: number) {
  imagePickerTarget.value = target;
  imagePickerServiceIndex.value = serviceIndex;
  imagePickerGenericId.value = genericId;
  imagePickerBannerIndex.value = bannerIndex;
  imagePickerPointsIndex.value = target === 'pointsSlide' ? serviceIndex : undefined;
  imagePickerOpen.value = true;
}

function selectImageAsset(items: Array<{ file_path: string }>) {
  const filePath = items[0]?.file_path;
  if (!filePath) return;
  if (imagePickerTarget.value === 'headerAvatar') {
    config.headerAvatar = filePath;
  } else if (imagePickerTarget.value === 'pointsSlide') {
    const index = imagePickerPointsIndex.value;
    if (index !== undefined && config.pointsSlides[index]) config.pointsSlides[index].image = filePath;
  } else if (imagePickerTarget.value === 'serviceIcon') {
    const index = imagePickerServiceIndex.value;
    if (index !== undefined && config.serviceList[index]) config.serviceList[index].image = filePath;
  } else {
    const component = genericComponents.value.find((item) => item.id === imagePickerGenericId.value);
    const index = imagePickerBannerIndex.value;
    if (component?.bannerImages && index !== undefined && component.bannerImages[index]) component.bannerImages[index].image = filePath;
  }
}

const services: NavigationItem[] = [
  { value: 0, label: '我的余额', icon: 'lucide:wallet-cards', link: '/pages/user/balance', enabled: true },
  { value: 1, label: '我的等级', icon: 'lucide:badge-check', link: '/pages/user/level', enabled: true },
  { value: 2, label: '积分中心', icon: 'lucide:circle-dollar-sign', link: '/pages/points/index', enabled: true },
  { value: 3, label: '签到', icon: 'lucide:calendar-check-2', link: '/pages/sign/index', enabled: true },
  { value: 4, label: '会员中心', icon: 'lucide:crown', link: '/pages/user/vip', enabled: true },
  { value: 5, label: '地址管理', icon: 'lucide:map-pin', link: '/pages/user/address', enabled: true },
  { value: 6, label: '发票管理', icon: 'lucide:receipt-text', link: '/pages/user/invoice', enabled: true },
  { value: 7, label: '我的收藏', icon: 'lucide:heart', link: '/pages/user/collect', enabled: true },
  { value: 8, label: '助力记录', icon: 'lucide:hand-heart', link: '/pages/activity/assist', enabled: true },
  { value: 9, label: '活动', icon: 'lucide:party-popper', link: '/pages/activity/index', enabled: true },
  { value: 10, label: '报名活动', icon: 'lucide:badge-plus', link: '/pages/activity/register', enabled: true },
];
const orderItems: OrderItem[] = [
  { label: '待付款', icon: 'lucide:credit-card', enabled: true },
  { label: '待发货/核销', icon: 'lucide:package-check', enabled: true },
  { label: '待收货', icon: 'lucide:truck', enabled: true },
  { label: '待评价', icon: 'lucide:message-circle-more', enabled: true },
  { label: '售后/退款', icon: 'lucide:circle-dollar-sign', enabled: true },
];

const config = reactive({
  pageName: '我的页面',
  shareTitle: '我的商城',
  showHeader: 1,
  showVip: 1,
  headerAvatar: '',
  headerLoginTitle: '请点击登录',
  headerChatEnabled: 1,
  headerStats: [
    { value: 0, label: '我的收藏', enabled: true },
    { value: 0, label: '关注店铺', enabled: true },
    { value: 0, label: '浏览记录', enabled: true },
    { value: 0, label: '优惠券', enabled: true },
  ] as HeaderStatistic[],
  showOrders: 1,
  orderTitle: '我的订单',
  orderMoreText: '全部订单',
  orderItems: [...orderItems],
  showPointsBanner: 1,
  pointsTitle: '来积分商城',
  pointsSubtitle: '积分还能兑换礼品',
  pointsActionText: '登录查看更多精彩',
  pointsLink: '/pages/points_mall/index',
  pointsImage: '',
  pointsSlides: [{ title: '来积分商城', subtitle: '积分还能兑换礼品', actionText: '登录查看更多精彩', link: '/pages/points_mall/index', image: '' }] as AdSlide[],
  pointsBackground: '#ff5b29',
  showServices: 1,
  serviceCount: 11,
  serviceList: [...services],
  visibleServices: services.map((item) => item.value),
  serviceNavigationType: 'icon-text',
  serviceRows: 4,
  serviceDisplayMode: 'fixed',
  serviceIconRadiusMode: 'all',
  serviceIconRadius: 8,
  serviceIconTopLeftRadius: 8,
  serviceIconTopRightRadius: 8,
  serviceIconBottomRightRadius: 8,
  serviceIconBottomLeftRadius: 8,
  serviceIconShadow: 'off',
  serviceFloat: 0,
  serviceBackground: 'rgba(255, 255, 255, 0)',
  serviceBackgroundEnd: 'rgba(255, 255, 255, 0)',
  serviceBottomBackground: 'rgba(255, 0, 0, 0)',
  serviceTextColor: '#333333',
  servicePaddingTop: 0,
  servicePaddingBottom: 10,
  servicePaddingLeft: 0,
  serviceMarginTop: 0,
  serviceRadiusMode: 'all',
  serviceRadius: 0,
  serviceTopLeftRadius: 0,
  serviceTopRightRadius: 0,
  serviceBottomRightRadius: 0,
  serviceBottomLeftRadius: 0,
  serviceShadow: 'off',
  bottomList: [
    { value: 0, label: '首页', icon: 'lucide:house', link: '/pages/index/index', enabled: true },
    { value: 1, label: '分类', icon: 'lucide:layout-grid', link: '/pages/category/index', enabled: true },
    { value: 2, label: '购物车', icon: 'lucide:shopping-cart', link: '/pages/cart/index', enabled: true },
    { value: 3, label: '我的', icon: 'lucide:circle-user-round', link: '/pages/user/index', enabled: true },
  ],
  visibleBottom: [0, 1, 2, 3],
  sectionStyles: {
    header: { ...createComponentStyle(), headerTextColor: '#ffffff', contentTextColor: '#ffffff' },
    orders: createComponentStyle(),
    points: createComponentStyle(),
    services: createComponentStyle(),
    bottom: createComponentStyle(),
  } as Record<SectionKey, ComponentStyle>,
});

const sections: Array<{ key: SectionKey; label: string; icon: string }> = [
  { key: 'header', label: '会员头部', icon: 'lucide:circle-user-round' },
  { key: 'orders', label: '我的订单', icon: 'lucide:package-check' },
  { key: 'points', label: '广告组', icon: 'lucide:gift' },
  { key: 'services', label: '我的服务', icon: 'lucide:grip' },
  { key: 'bottom', label: '底部导航', icon: 'lucide:panel-bottom' },
];

// 左侧始终展示完整的通用 DIY 组件库；个人页固有模块只是其中的一部分。
const libraryGroups: Array<{ label: string; items: LibraryItem[] }> = [
  {
    label: '页面模块',
    items: [
      { type: 'memberHeader', label: '会员头部', icon: 'lucide:circle-user-round', section: 'header' },
      { type: 'myOrders', label: '我的订单', icon: 'lucide:package-check', section: 'orders' },
      { type: 'pointsBanner', label: '广告组', icon: 'lucide:gift', section: 'points' },
      { type: 'myServices', label: '我的服务', icon: 'lucide:layout-grid', section: 'services' },
      { type: 'bottomNav', label: '底部导航', icon: 'lucide:panel-bottom', section: 'bottom' },
    ],
  },
  {
    label: '基础组件',
    items: [
      { type: 'search', label: '搜索框', icon: 'lucide:search' }, { type: 'topMerge', label: '轮播搜索', icon: 'lucide:images' },
      { type: 'option', label: '选项卡', icon: 'lucide:layout-grid' }, { type: 'banner', label: '图片轮播', icon: 'lucide:images' },
      { type: 'imageSingle', label: '单图组', icon: 'lucide:image' }, { type: 'window', label: '图片橱窗', icon: 'lucide:images' },
      { type: 'navBar', label: '导航组', icon: 'lucide:layout-grid' }, { type: 'product', label: '商品组', icon: 'lucide:package' },
      { type: 'article', label: '文章组', icon: 'lucide:book-open' }, { type: 'special', label: '头条快报', icon: 'lucide:newspaper' },
      { type: 'notice', label: '公告组', icon: 'lucide:megaphone' }, { type: 'hotspot', label: '热区', icon: 'lucide:mouse-pointer' },
      { type: 'video', label: '视频组', icon: 'lucide:video' }, { type: 'title', label: '标题', icon: 'ant-design:font-size-outlined' },
      { type: 'store', label: '店铺组', icon: 'lucide:store' },
    ],
  },
  {
    label: '营销组件',
    items: [
      { type: 'coupon', label: '优惠券', icon: 'lucide:ticket' }, { type: 'assembleProduct', label: '拼团商品', icon: 'lucide:users' },
      { type: 'bargainProduct', label: '砍价商品', icon: 'lucide:percent' }, { type: 'seckillProduct', label: '秒杀商品', icon: 'lucide:timer' },
      { type: 'previewProduct', label: '预售商品', icon: 'lucide:calendar-clock' }, { type: 'newActivity', label: '活动专区', icon: 'lucide:sparkles' },
      { type: 'qixiLive', label: '直播', icon: 'lucide:radio' }, { type: 'videoLive', label: '视频直播', icon: 'lucide:video' },
    ],
  },
  {
    label: '工具组件',
    items: [
      { type: 'guide', label: '辅助线', icon: 'lucide:scan-line' }, { type: 'blank', label: '辅助空白', icon: 'lucide:file' },
      { type: 'richText', label: '富文本', icon: 'lucide:align-left' }, { type: 'service', label: '在线客服', icon: 'lucide:headphones' },
      { type: 'surface', label: '悬浮按钮', icon: 'lucide:circle-dot' },
    ],
  },
];

const visibleServices = computed(() => config.serviceList
  .filter((item) => item.enabled)
  .slice(0, config.serviceCount));
const visibleBottom = computed(() => config.bottomList.filter((item) => item.enabled && config.visibleBottom.includes(item.value)));
const visibleHeaderStats = computed(() => config.headerStats.filter((item) => item.enabled));
const visibleOrderItems = computed(() => config.orderItems.filter((item) => item.enabled));
const currentPointsSlide = computed(() => config.pointsSlides[activePointsSlide.value] ?? config.pointsSlides[0]);
const activeGenericComponent = computed(() => genericComponents.value.find((item) => item.id === activeGenericId.value));
const activeStyle = computed<ComponentStyle | undefined>(() => activeSection.value === 'generic'
  ? activeGenericComponent.value?.style
  : config.sectionStyles[activeSection.value]);
const activeTitle = computed(() => activeSection.value === 'generic'
  ? activeGenericComponent.value?.label ?? '通用组件'
  : sections.find((item) => item.key === activeSection.value)?.label ?? '页面设置');

// 服务显示状态以服务卡片本身为唯一来源；保留 visibleServices 仅用于兼容旧的已保存配置。
watch(
  () => config.serviceList.map((item) => ({ enabled: item.enabled, value: item.value })),
  () => {
    config.visibleServices = config.serviceList
      .filter((item) => item.enabled)
      .map((item) => item.value);
  },
  { deep: true },
);
const servicePreviewStyle = computed(() => {
  const radius = config.serviceRadiusMode === 'individual'
    ? `${config.serviceTopLeftRadius}px ${config.serviceTopRightRadius}px ${config.serviceBottomRightRadius}px ${config.serviceBottomLeftRadius}px`
    : `${config.serviceRadius}px`;
  return {
    background: config.serviceBackgroundEnd === 'rgba(255, 255, 255, 0)'
      ? config.serviceBackground
      : `linear-gradient(180deg, ${config.serviceBackground}, ${config.serviceBackgroundEnd})`,
    borderRadius: radius,
    boxShadow: config.serviceShadow === 'on' ? '0 6px 16px rgb(25 39 70 / 14%)' : 'none',
    color: config.serviceTextColor,
    marginTop: `${config.serviceMarginTop - config.serviceFloat}px`,
    padding: `${config.servicePaddingTop}px ${config.servicePaddingLeft}px ${config.servicePaddingBottom}px`,
  };
});
const serviceIconStyle = computed(() => {
  const radius = config.serviceIconRadiusMode === 'individual'
    ? `${config.serviceIconTopLeftRadius}px ${config.serviceIconTopRightRadius}px ${config.serviceIconBottomRightRadius}px ${config.serviceIconBottomLeftRadius}px`
    : `${config.serviceIconRadius}px`;
  return { borderRadius: radius, boxShadow: config.serviceIconShadow === 'on' ? '0 4px 10px rgb(25 39 70 / 16%)' : 'none' };
});
const pointsPreviewStyle = computed(() => ({
  background: currentPointsSlide.value?.image
    ? `center / cover no-repeat url("${currentPointsSlide.value.image}")`
    : config.pointsBackground,
}));

function moduleCardStyle(style: ComponentStyle) {
  const background = style.cardBackgroundEnd === style.cardBackground
    ? style.cardBackground
    : `linear-gradient(180deg, ${style.cardBackground}, ${style.cardBackgroundEnd})`;
  return {
    ...(style.cardBackground !== 'transparent' || style.cardBackgroundEnd !== 'transparent' ? { background } : {}),
    border: `${style.cardBorderWidth}px solid ${style.cardBorderColor}`,
    ...(style.cardRadius > 0 ? { borderRadius: `${style.cardRadius}px` } : {}),
    boxShadow: style.cardShadow === 'on'
      ? `${style.cardShadowX}px ${style.cardShadowY}px ${style.cardShadowBlur}px ${style.cardShadowSpread}px ${style.cardShadowColor}`
      : 'none',
    marginTop: `${style.cardMarginTop - style.cardFloat}px`,
    ...(style.cardPaddingTop || style.cardPaddingHorizontal || style.cardPaddingBottom
      ? { padding: `${style.cardPaddingTop}px ${style.cardPaddingHorizontal}px ${style.cardPaddingBottom}px` }
      : {}),
  };
}

function moduleHeaderStyle(style: ComponentStyle) {
  return {
    ...(style.headerBackground !== 'transparent' ? { background: style.headerBackground } : {}),
    borderBottom: `${style.headerBorderWidth}px solid ${style.headerBorderColor}`,
    color: style.headerTextColor,
    fontSize: `${style.headerFontSize}px`,
    ...(style.headerPaddingTop || style.headerPaddingBottom
      ? { paddingTop: `${style.headerPaddingTop}px`, paddingBottom: `${style.headerPaddingBottom}px` }
      : {}),
  };
}

function moduleContentStyle(style: ComponentStyle) {
  return {
    ...(style.contentBackground !== 'transparent' ? { background: style.contentBackground } : {}),
    border: `${style.contentBorderWidth}px solid ${style.contentBorderColor}`,
    ...(style.contentRadius > 0 ? { borderRadius: `${style.contentRadius}px` } : {}),
    boxShadow: style.contentShadow === 'on'
      ? `${style.contentShadowX}px ${style.contentShadowY}px ${style.contentShadowBlur}px ${style.contentShadowSpread}px ${style.contentShadowColor}`
      : 'none',
    color: style.contentTextColor,
    fontSize: `${style.contentFontSize}px`,
    ...(style.contentPaddingTop || style.contentPaddingHorizontal || style.contentPaddingBottom
      ? { padding: `${style.contentPaddingTop}px ${style.contentPaddingHorizontal}px ${style.contentPaddingBottom}px` }
      : {}),
  };
}

function resetActiveStyle(key: keyof ComponentStyle) {
  const style = activeStyle.value;
  if (style) style[key] = createComponentStyle()[key] as never;
}

function copyObject(target: Record<string, any>, source: Record<string, any>) {
  for (const [key, value] of Object.entries(source)) {
    if (key in target) target[key] = value;
  }
}

function normalizeNavigationItems(items: NavigationItem[], fallback: NavigationItem[]) {
  if (!Array.isArray(items) || !items.length) return fallback.map((item) => ({ ...item }));
  return items.map((item, index) => ({
    ...fallback[index % fallback.length],
    ...item,
    enabled: item.enabled ?? true,
    link: item.link ?? '',
  }));
}

function applySavedConfig(saved: Record<string, any>) {
  if (!saved || !Object.keys(saved).length) return;
  copyObject(config, saved);
  const legacyPointsSlide: AdSlide = {
    title: config.pointsTitle || '来积分商城',
    subtitle: config.pointsSubtitle || '积分还能兑换礼品',
    actionText: config.pointsActionText || '登录查看更多精彩',
    link: config.pointsLink || '/pages/points_mall/index',
    image: config.pointsImage || '',
  };
  config.pointsSlides = (Array.isArray(config.pointsSlides) && config.pointsSlides.length ? config.pointsSlides : [legacyPointsSlide])
    .slice(0, 5)
    .map((item: Partial<AdSlide>) => ({
      title: item.title ?? legacyPointsSlide.title,
      subtitle: item.subtitle ?? legacyPointsSlide.subtitle,
      actionText: item.actionText ?? legacyPointsSlide.actionText,
      link: item.link ?? legacyPointsSlide.link,
      image: item.image ?? legacyPointsSlide.image,
    }));
  activePointsSlide.value = 0;
  config.serviceList = normalizeNavigationItems(config.serviceList, services);
  config.visibleServices = config.serviceList
    .filter((item) => item.enabled)
    .map((item) => item.value);
  config.bottomList = normalizeNavigationItems(config.bottomList, [
    { value: 0, label: '首页', icon: 'lucide:house', link: '/pages/index/index', enabled: true },
    { value: 1, label: '分类', icon: 'lucide:layout-grid', link: '/pages/category/index', enabled: true },
    { value: 2, label: '购物车', icon: 'lucide:shopping-cart', link: '/pages/cart/index', enabled: true },
    { value: 3, label: '我的', icon: 'lucide:circle-user-round', link: '/pages/user/index', enabled: true },
  ]);
  if (!Array.isArray(config.headerStats) || !config.headerStats.length) {
    config.headerStats = [
      { value: 0, label: '我的收藏', enabled: true },
      { value: 0, label: '关注店铺', enabled: true },
      { value: 0, label: '浏览记录', enabled: true },
      { value: 0, label: '优惠券', enabled: true },
    ];
  }
  if (!Array.isArray(config.orderItems) || !config.orderItems.length) config.orderItems = [...orderItems];
  for (const key of sections.map((item) => item.key)) {
    config.sectionStyles[key] = { ...createComponentStyle(), ...(saved.sectionStyles?.[key] ?? config.sectionStyles[key] ?? {}) };
  }
  if (Array.isArray(saved.genericComponents)) {
    genericComponents.value = saved.genericComponents
      .filter((item: GenericComponent) => item?.type && item?.label && item?.icon)
      .map((item: Partial<GenericComponent>, index: number) => ({
        id: item.id ?? `${item.type}-${index}`,
        type: item.type!,
        label: item.label!,
        icon: item.icon!,
        title: item.title ?? item.label!,
        description: item.description ?? '通用组件内容，点击可编辑',
        link: item.link ?? '',
        enabled: item.enabled ?? true,
        bannerImages: item.type === 'banner'
          ? (item.bannerImages?.length ? item.bannerImages : [{ image: '', link: '' }])
          : undefined,
        productSource: item.productSource ?? 'all',
        productCount: item.productCount ?? 6,
        productColumns: item.productColumns ?? 2,
        productLayout: item.productLayout ?? 'waterfall',
        style: { ...createComponentStyle(), ...(item.style ?? {}) },
      }));
  }
}

function addPointsSlide() {
  if (config.pointsSlides.length >= 5) {
    ElMessage.warning('广告组最多添加 5 项');
    return;
  }
  config.pointsSlides.push({
    title: '广告标题',
    subtitle: '广告副标题',
    actionText: '立即查看',
    link: '',
    image: '',
  });
  activePointsSlide.value = config.pointsSlides.length - 1;
}

function removePointsSlide(index: number) {
  if (config.pointsSlides.length <= 1) {
    ElMessage.warning('广告组至少保留 1 项');
    return;
  }
  config.pointsSlides.splice(index, 1);
  activePointsSlide.value = Math.min(activePointsSlide.value, config.pointsSlides.length - 1);
}

async function load() {
  try {
    const result = await getPersonalDecorationApi();
    applySavedConfig(result?.config ?? {});
  } catch {
    const draft = window.localStorage.getItem(personalDraftStorageKey);
    if (!draft) return;
    try { applySavedConfig(JSON.parse(draft)); } catch { window.localStorage.removeItem(personalDraftStorageKey); }
  }
}

function selectSection(key: SectionKey) {
  activeSection.value = key;
  activeGenericId.value = undefined;
  activeConfigTab.value = 'content';
  configScrollRef.value?.scrollTo({ top: 0, behavior: 'smooth' });
}

function addLibraryItem(item: LibraryItem) {
  if (item.section) {
    selectSection(item.section);
    return;
  }
  const existing = genericComponents.value.find((component) => component.type === item.type);
  if (existing) {
    selectGenericComponent(existing.id);
    return;
  }
  const component: GenericComponent = {
    id: `${item.type}-${Date.now()}`,
    type: item.type,
    label: item.label,
    icon: item.icon,
    title: item.label,
    description: item.type === 'banner' ? '请添加轮播图片' : item.type === 'product' ? '从商品库选择商品后展示' : '通用组件内容，点击可编辑',
    link: '',
    enabled: true,
    bannerImages: item.type === 'banner' ? [{ image: '', link: '' }] : undefined,
    productSource: 'all',
    productCount: 6,
    productColumns: 2,
    productLayout: 'waterfall',
    style: createComponentStyle(),
  };
  genericComponents.value.push(component);
  selectGenericComponent(component.id);
}

function addBannerImage(component: GenericComponent) {
  component.bannerImages ??= [];
  if (component.bannerImages.length >= 5) return ElMessage.warning('图片轮播最多添加5张图片');
  component.bannerImages.push({ image: '', link: '' });
}

function removeBannerImage(component: GenericComponent, index: number) {
  component.bannerImages?.splice(index, 1);
  if (!component.bannerImages?.length) component.bannerImages = [{ image: '', link: '' }];
}

function selectGenericComponent(id: string) {
  activeSection.value = 'generic';
  activeGenericId.value = id;
  activeConfigTab.value = 'content';
  configScrollRef.value?.scrollTo({ top: 0, behavior: 'smooth' });
}

function isLibraryItemActive(item: LibraryItem) {
  return item.section
    ? activeSection.value === item.section
    : activeSection.value === 'generic' && activeGenericComponent.value?.type === item.type;
}

function addService() {
  if (config.serviceList.length >= 13) {
    ElMessage.warning('我的服务最多添加13项');
    return;
  }
  const value = Math.max(-1, ...config.serviceList.map((item) => item.value)) + 1;
  config.serviceList.push({ value, label: '服务名称', icon: 'lucide:circle-dot', link: '', enabled: true });
  // 新增入口必须可见，不能被旧的“显示数量”截断而只存在于右侧配置。
  config.serviceCount = Math.max(config.serviceCount, config.serviceList.length);
  config.visibleServices = config.serviceList.filter((item) => item.enabled).map((item) => item.value);
}

function removeService(index: number) {
  const [removed] = config.serviceList.splice(index, 1);
  if (!removed) return;
  config.visibleServices = config.serviceList.filter((item) => item.enabled).map((item) => item.value);
}

function resetServiceStyle(key: keyof typeof config) {
  const defaults: Partial<typeof config> = {
    serviceBackground: 'rgba(255, 255, 255, 0)',
    serviceBackgroundEnd: 'rgba(255, 255, 255, 0)',
    serviceBottomBackground: 'rgba(255, 0, 0, 0)',
    serviceTextColor: '#333333',
    servicePaddingTop: 0,
    servicePaddingBottom: 10,
    servicePaddingLeft: 0,
    serviceMarginTop: 0,
  };
  const fallback = defaults[key];
  if (fallback !== undefined) config[key] = fallback as never;
}

function dragStart(list: 'service' | 'bottom', index: number) {
  if (list === 'service') draggingService.value = index;
  else draggingBottom.value = index;
}

function dropItem(list: 'service' | 'bottom', target: number) {
  const source = list === 'service' ? draggingService.value : draggingBottom.value;
  if (list === 'service') draggingService.value = undefined;
  else draggingBottom.value = undefined;
  if (source === undefined || source === target) return;
  const targetList = list === 'service' ? config.serviceList : config.bottomList;
  const [item] = targetList.splice(source, 1);
  if (item) targetList.splice(target, 0, item);
}

function dragOrderStart(index: number) {
  draggingOrder.value = index;
}

function dropOrder(target: number) {
  const source = draggingOrder.value;
  draggingOrder.value = undefined;
  if (source === undefined || source === target) return;
  const [item] = config.orderItems.splice(source, 1);
  if (item) config.orderItems.splice(target, 0, item);
}

async function save() {
  saving.value = true;
  const payload = { ...JSON.parse(JSON.stringify(config)), genericComponents: genericComponents.value };
  try {
    await savePersonalDecorationApi(payload);
    window.localStorage.removeItem(personalDraftStorageKey);
    ElMessage.success('我的装修已保存');
  } catch {
    window.localStorage.setItem(personalDraftStorageKey, JSON.stringify(payload));
    ElMessage.warning('接口服务尚未更新，当前配置已保存为本地草稿');
  } finally {
    saving.value = false;
  }
}

onMounted(load);
</script>

<template>
  <Page auto-content-height content-class="personal-decoration-page !p-0">
    <div class="personal-editor">
      <aside class="personal-editor__library">
        <section v-for="group in libraryGroups" :key="group.label" class="personal-library__group">
          <h3>{{ group.label }}</h3>
          <div class="personal-library__grid">
            <button v-for="item in group.items" :key="item.type" type="button" :class="{ 'is-active': isLibraryItemActive(item) }" @click="addLibraryItem(item)">
              <IconifyIcon :icon="item.icon" />
              <span>{{ item.label }}</span>
            </button>
          </div>
        </section>
      </aside>

      <main class="personal-editor__canvas">
        <section class="personal-phone" aria-label="我的页面实时预览">
          <div class="personal-phone__title">我的</div>
          <div class="personal-phone__screen">
            <section v-if="config.showHeader === 1" class="personal-preview__header personal-preview__module" :class="{ 'is-active': activeSection === 'header' }" :style="moduleCardStyle(config.sectionStyles.header)" @click="selectSection('header')">
              <button class="personal-preview__tip" type="button">会员头部</button>
              <div class="personal-preview__profile" :style="moduleHeaderStyle(config.sectionStyles.header)"><span class="personal-preview__avatar"><img v-if="config.headerAvatar" :src="config.headerAvatar" alt="会员头像" /><IconifyIcon v-else icon="lucide:user-round" /></span><button type="button" :style="{ color: config.sectionStyles.header.headerTextColor }">{{ config.headerLoginTitle }}</button><IconifyIcon v-if="config.headerChatEnabled === 1" class="personal-preview__chat" icon="lucide:message-circle" /></div>
              <div class="personal-preview__stats" :style="moduleContentStyle(config.sectionStyles.header)"><span v-for="item in visibleHeaderStats" :key="item.label"><b>{{ item.value }}</b>{{ item.label }}</span></div>
              <button v-if="config.showVip === 1" class="personal-preview__vip" type="button"><IconifyIcon icon="lucide:badge-dollar-sign" /><span>开通享六大特权，省钱又省心</span><b>立即开通</b></button>
            </section>

            <section v-if="config.showOrders === 1" class="personal-preview__orders personal-preview__module" :class="{ 'is-active': activeSection === 'orders' }" :style="moduleCardStyle(config.sectionStyles.orders)" @click="selectSection('orders')">
              <button class="personal-preview__tip" type="button">我的订单</button>
              <header :style="moduleHeaderStyle(config.sectionStyles.orders)"><b>{{ config.orderTitle }}</b><button type="button">{{ config.orderMoreText }} <IconifyIcon icon="lucide:chevron-right" /></button></header>
              <div :style="moduleContentStyle(config.sectionStyles.orders)"><button v-for="item in visibleOrderItems" :key="item.label" type="button" :style="{ color: config.sectionStyles.orders.contentTextColor }"><IconifyIcon :icon="item.icon" />{{ item.label }}</button></div>
            </section>

            <section v-if="config.showPointsBanner === 1" class="personal-preview__points personal-preview__module" :class="{ 'is-active': activeSection === 'points' }" :style="{ ...pointsPreviewStyle, ...moduleCardStyle(config.sectionStyles.points) }" @click="selectSection('points')">
              <button class="personal-preview__tip" type="button">广告组</button>
              <template v-if="currentPointsSlide">
                <b :style="moduleHeaderStyle(config.sectionStyles.points)">{{ currentPointsSlide.title }}</b>
                <small :style="moduleContentStyle(config.sectionStyles.points)">{{ currentPointsSlide.subtitle }}</small>
                <button class="personal-preview__points-action" type="button">{{ currentPointsSlide.actionText }}</button>
              </template>
              <div v-if="config.pointsSlides.length > 1" class="personal-preview__points-dots" aria-label="广告轮播切换">
                <button v-for="(_, index) in config.pointsSlides" :key="index" :class="{ active: index === activePointsSlide }" type="button" :aria-label="`切换到广告 ${index + 1}`" />
              </div>
            </section>

            <section v-if="config.showServices === 1" class="personal-preview__services personal-preview__module" :class="{ 'is-active': activeSection === 'services' }" :style="{ ...servicePreviewStyle, ...moduleCardStyle(config.sectionStyles.services) }" @click="selectSection('services')">
              <button class="personal-preview__tip" type="button">我的服务</button>
              <h3 :style="moduleHeaderStyle(config.sectionStyles.services)">我的服务</h3>
              <div :style="{ ...moduleContentStyle(config.sectionStyles.services), gridTemplateColumns: `repeat(${config.serviceRows}, minmax(0, 1fr))` }"><button v-for="item in visibleServices" :key="item.value" type="button" :style="{ color: config.sectionStyles.services.contentTextColor }"><img v-if="item.image && config.serviceNavigationType !== 'text'" class="personal-preview__service-image" :src="item.image" :style="serviceIconStyle" :alt="item.label" /><IconifyIcon v-else-if="config.serviceNavigationType !== 'text'" :icon="item.icon" :style="serviceIconStyle" /><span v-if="config.serviceNavigationType !== 'icon'">{{ item.label }}</span></button></div>
            </section>

            <section v-for="component in genericComponents.filter((item) => item.enabled)" :key="component.id" class="personal-preview__generic personal-preview__module" :class="[{ 'is-active': activeSection === 'generic' && activeGenericId === component.id }, `personal-preview__generic--${component.type}`]" :style="moduleCardStyle(component.style)" @click="selectGenericComponent(component.id)">
              <button class="personal-preview__tip" type="button">{{ component.label }}</button>
              <template v-if="component.type === 'banner'">
                <div class="personal-preview__banner" :style="moduleContentStyle(component.style)">
                  <img v-if="component.bannerImages?.[0]?.image" :src="component.bannerImages[0].image" alt="轮播图片" />
                  <div v-else class="personal-preview__banner-placeholder"><IconifyIcon icon="lucide:images" /><span>点击右侧添加轮播图片</span></div>
                  <div v-if="component.bannerImages && component.bannerImages.length > 1" class="personal-preview__dots"><i v-for="(_, index) in component.bannerImages" :key="index" :class="{ active: index === 0 }" /></div>
                </div>
              </template>
              <template v-else-if="component.type === 'product'">
                <header :style="moduleHeaderStyle(component.style)"><b>{{ component.title }}</b><span>{{ component.productSource === 'hot' ? '热门推荐' : component.productSource === 'new' ? '新品上架' : '全部商品' }}</span></header>
                <div class="personal-preview__product-grid" :style="{ ...moduleContentStyle(component.style), gridTemplateColumns: `repeat(${component.productColumns ?? 2}, minmax(0, 1fr))` }">
                  <article v-for="index in Math.min(component.productCount ?? 6, 6)" :key="index"><div class="personal-preview__product-image"><IconifyIcon icon="lucide:package" /></div><b>精选商品名称</b><em>¥99.00</em></article>
                </div>
              </template>
              <template v-else>
                <header :style="moduleHeaderStyle(component.style)"><IconifyIcon :icon="component.icon" /><b>{{ component.title }}</b></header>
                <div :style="moduleContentStyle(component.style)"><IconifyIcon :icon="component.icon" /><span :style="{ color: component.style.contentTextColor }">{{ component.description }}</span><IconifyIcon class="personal-preview__generic-add" icon="lucide:settings-2" /></div>
              </template>
            </section>

            <nav class="personal-preview__bottom personal-preview__module" :class="{ 'is-active': activeSection === 'bottom' }" :style="moduleCardStyle(config.sectionStyles.bottom)" @click="selectSection('bottom')">
              <button class="personal-preview__tip" type="button">底部导航</button>
              <button v-for="item in visibleBottom" :key="item.value" type="button" :class="{ 'is-current': item.label === '我的' }"><IconifyIcon :icon="item.icon" /><span>{{ item.label }}</span></button>
            </nav>
          </div>
        </section>
      </main>

      <aside class="personal-editor__config">
        <div class="personal-config__title"><span>{{ activeTitle }}</span><button type="button" :class="{ active: activeConfigTab === 'content' }" @click="activeConfigTab = 'content'">内容</button><button type="button" :class="{ active: activeConfigTab === 'style' }" @click="activeConfigTab = 'style'">样式</button></div>
        <div ref="configScrollRef" class="personal-config__scroll">
          <template v-if="activeConfigTab === 'style' && activeStyle">
            <section class="personal-config__section personal-config__section--settings">
              <h3>头部样式</h3>
              <div class="personal-config__color-row"><span>背景颜色</span><ElColorPicker v-model="activeStyle.headerBackground" show-alpha /><ElInput v-model="activeStyle.headerBackground" /><button type="button" @click="resetActiveStyle('headerBackground')">重置</button></div>
              <div class="personal-config__color-row"><span>文字颜色</span><ElColorPicker v-model="activeStyle.headerTextColor" show-alpha /><ElInput v-model="activeStyle.headerTextColor" /><button type="button" @click="resetActiveStyle('headerTextColor')">重置</button></div>
              <div class="personal-config__row personal-config__row--slider"><span>标题字号</span><ElSlider v-model="activeStyle.headerFontSize" :min="12" :max="28" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.headerFontSize" :min="12" :max="28" controls-position="right" /></div>
              <div class="personal-config__row personal-config__row--slider"><span>上内边距</span><ElSlider v-model="activeStyle.headerPaddingTop" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.headerPaddingTop" :min="0" :max="48" controls-position="right" /></div>
              <div class="personal-config__row personal-config__row--slider"><span>下内边距</span><ElSlider v-model="activeStyle.headerPaddingBottom" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.headerPaddingBottom" :min="0" :max="48" controls-position="right" /></div>
              <div class="personal-config__color-row"><span>底部边框</span><ElColorPicker v-model="activeStyle.headerBorderColor" show-alpha /><ElInput v-model="activeStyle.headerBorderColor" /><button type="button" @click="resetActiveStyle('headerBorderColor')">重置</button></div>
              <div class="personal-config__row personal-config__row--slider"><span>边框宽度</span><ElSlider v-model="activeStyle.headerBorderWidth" :max="6" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.headerBorderWidth" :min="0" :max="6" controls-position="right" /></div>
            </section>
            <section class="personal-config__section personal-config__section--settings">
              <h3>{{ activeSection === 'generic' && activeGenericComponent?.type === 'product' ? '商品样式' : '内容样式' }}</h3>
              <div class="personal-config__color-row"><span>内容背景</span><ElColorPicker v-model="activeStyle.contentBackground" show-alpha /><ElInput v-model="activeStyle.contentBackground" /><button type="button" @click="resetActiveStyle('contentBackground')">重置</button></div>
              <div class="personal-config__color-row"><span>文字颜色</span><ElColorPicker v-model="activeStyle.contentTextColor" show-alpha /><ElInput v-model="activeStyle.contentTextColor" /><button type="button" @click="resetActiveStyle('contentTextColor')">重置</button></div>
              <div class="personal-config__row personal-config__row--slider"><span>文字字号</span><ElSlider v-model="activeStyle.contentFontSize" :min="12" :max="24" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.contentFontSize" :min="12" :max="24" controls-position="right" /></div>
              <div class="personal-config__row personal-config__row--slider"><span>上内边距</span><ElSlider v-model="activeStyle.contentPaddingTop" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.contentPaddingTop" :min="0" :max="48" controls-position="right" /></div>
              <div class="personal-config__row personal-config__row--slider"><span>下内边距</span><ElSlider v-model="activeStyle.contentPaddingBottom" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.contentPaddingBottom" :min="0" :max="48" controls-position="right" /></div>
              <div class="personal-config__row personal-config__row--slider"><span>左右内边距</span><ElSlider v-model="activeStyle.contentPaddingHorizontal" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.contentPaddingHorizontal" :min="0" :max="48" controls-position="right" /></div>
              <div class="personal-config__color-row"><span>内容边框</span><ElColorPicker v-model="activeStyle.contentBorderColor" show-alpha /><ElInput v-model="activeStyle.contentBorderColor" /><button type="button" @click="resetActiveStyle('contentBorderColor')">重置</button></div>
              <div class="personal-config__row personal-config__row--slider"><span>圆角大小</span><ElSlider v-model="activeStyle.contentRadius" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.contentRadius" :min="0" :max="48" controls-position="right" /></div>
              <div class="personal-config__row"><span>开启阴影</span><ElRadioGroup v-model="activeStyle.contentShadow"><ElRadio value="off">关闭</ElRadio><ElRadio value="on">开启</ElRadio></ElRadioGroup></div>
              <div v-if="activeStyle.contentShadow === 'on'" class="personal-config__shadow-settings"><div class="personal-config__color-row"><span>阴影颜色</span><ElColorPicker v-model="activeStyle.contentShadowColor" show-alpha /><ElInput v-model="activeStyle.contentShadowColor" /><button type="button" @click="resetActiveStyle('contentShadowColor')">重置</button></div><div class="personal-config__row personal-config__row--slider"><span>横轴</span><ElSlider v-model="activeStyle.contentShadowX" :min="-48" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.contentShadowX" :min="-48" :max="48" controls-position="right" /></div><div class="personal-config__row personal-config__row--slider"><span>纵轴</span><ElSlider v-model="activeStyle.contentShadowY" :min="-48" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.contentShadowY" :min="-48" :max="48" controls-position="right" /></div><div class="personal-config__row personal-config__row--slider"><span>宽度</span><ElSlider v-model="activeStyle.contentShadowBlur" :max="100" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.contentShadowBlur" :min="0" :max="100" controls-position="right" /></div><div class="personal-config__row personal-config__row--slider"><span>扩散</span><ElSlider v-model="activeStyle.contentShadowSpread" :min="-48" :max="100" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.contentShadowSpread" :min="-48" :max="100" controls-position="right" /></div></div>
            </section>
            <section class="personal-config__section personal-config__section--settings">
              <h3>卡片样式</h3>
              <div class="personal-config__row personal-config__row--slider"><span>组件上浮</span><ElSlider v-model="activeStyle.cardFloat" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.cardFloat" :min="0" :max="48" controls-position="right" /></div>
              <div class="personal-config__color-row"><span>组件背景</span><ElColorPicker v-model="activeStyle.cardBackground" show-alpha /><ElInput v-model="activeStyle.cardBackground" /><button type="button" @click="resetActiveStyle('cardBackground')">重置</button></div>
              <div class="personal-config__color-row"><span>背景渐变</span><ElColorPicker v-model="activeStyle.cardBackgroundEnd" show-alpha /><ElInput v-model="activeStyle.cardBackgroundEnd" /><button type="button" @click="resetActiveStyle('cardBackgroundEnd')">重置</button></div>
              <div class="personal-config__color-row"><span>卡片边框</span><ElColorPicker v-model="activeStyle.cardBorderColor" show-alpha /><ElInput v-model="activeStyle.cardBorderColor" /><button type="button" @click="resetActiveStyle('cardBorderColor')">重置</button></div>
              <div class="personal-config__row personal-config__row--slider"><span>边框宽度</span><ElSlider v-model="activeStyle.cardBorderWidth" :max="6" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.cardBorderWidth" :min="0" :max="6" controls-position="right" /></div>
              <div class="personal-config__row personal-config__row--slider"><span>上内边距</span><ElSlider v-model="activeStyle.cardPaddingTop" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.cardPaddingTop" :min="0" :max="48" controls-position="right" /></div>
              <div class="personal-config__row personal-config__row--slider"><span>下内边距</span><ElSlider v-model="activeStyle.cardPaddingBottom" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.cardPaddingBottom" :min="0" :max="48" controls-position="right" /></div>
              <div class="personal-config__row personal-config__row--slider"><span>左右内边距</span><ElSlider v-model="activeStyle.cardPaddingHorizontal" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.cardPaddingHorizontal" :min="0" :max="48" controls-position="right" /></div>
              <div class="personal-config__row personal-config__row--slider"><span>页面上间距</span><ElSlider v-model="activeStyle.cardMarginTop" :max="96" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.cardMarginTop" :min="0" :max="96" controls-position="right" /></div>
              <div class="personal-config__row personal-config__row--slider"><span>背景圆角</span><ElSlider v-model="activeStyle.cardRadius" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.cardRadius" :min="0" :max="48" controls-position="right" /></div>
              <div class="personal-config__row"><span>开启阴影</span><ElRadioGroup v-model="activeStyle.cardShadow"><ElRadio value="off">关闭</ElRadio><ElRadio value="on">开启</ElRadio></ElRadioGroup></div>
              <div v-if="activeStyle.cardShadow === 'on'" class="personal-config__shadow-settings"><div class="personal-config__color-row"><span>阴影颜色</span><ElColorPicker v-model="activeStyle.cardShadowColor" show-alpha /><ElInput v-model="activeStyle.cardShadowColor" /><button type="button" @click="resetActiveStyle('cardShadowColor')">重置</button></div><div class="personal-config__row personal-config__row--slider"><span>横轴</span><ElSlider v-model="activeStyle.cardShadowX" :min="-48" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.cardShadowX" :min="-48" :max="48" controls-position="right" /></div><div class="personal-config__row personal-config__row--slider"><span>纵轴</span><ElSlider v-model="activeStyle.cardShadowY" :min="-48" :max="48" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.cardShadowY" :min="-48" :max="48" controls-position="right" /></div><div class="personal-config__row personal-config__row--slider"><span>宽度</span><ElSlider v-model="activeStyle.cardShadowBlur" :max="100" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.cardShadowBlur" :min="0" :max="100" controls-position="right" /></div><div class="personal-config__row personal-config__row--slider"><span>扩散</span><ElSlider v-model="activeStyle.cardShadowSpread" :min="-48" :max="100" :show-tooltip="false" /><ElInputNumber v-model="activeStyle.cardShadowSpread" :min="-48" :max="100" controls-position="right" /></div></div>
            </section>
          </template>
          <template v-else-if="activeSection === 'header'">
            <section class="personal-config__section"><h3>页面设置</h3><div class="personal-config__row"><span>页面名称</span><ElInput v-model="config.pageName" /></div><div class="personal-config__row"><span>分享标题</span><ElInput v-model="config.shareTitle" /></div></section>
            <section class="personal-config__section"><h3>会员头部</h3><div class="personal-config__row"><span>是否显示</span><ElRadioGroup v-model="config.showHeader"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div><div class="personal-config__row"><span>会员头像</span><button class="personal-config__image-select" type="button" @click="openImagePicker('headerAvatar')"><img v-if="config.headerAvatar" :src="config.headerAvatar" alt="会员头像" /><IconifyIcon v-else icon="lucide:image-plus" /><span>{{ config.headerAvatar ? '更换图片' : '选择图片' }}</span></button><button v-if="config.headerAvatar" class="personal-config__clear-image" type="button" @click="config.headerAvatar = ''">清除</button></div><div class="personal-config__row"><span>登录文案</span><ElInput v-model="config.headerLoginTitle" maxlength="12" show-word-limit /></div><div class="personal-config__row"><span>消息入口</span><ElRadioGroup v-model="config.headerChatEnabled"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div><div class="personal-config__row"><span>会员入口</span><ElRadioGroup v-model="config.showVip"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div></section>
            <section class="personal-config__section"><h3>会员数据</h3><p class="personal-config__hint">勾选控制展示，文案和数值会立即同步到预览。</p><div class="personal-config__simple-list"><label v-for="item in config.headerStats" :key="item.label"><ElSwitch v-model="item.enabled" /><ElInput v-model="item.label" maxlength="6" /><ElInputNumber v-model="item.value" :min="0" :max="99999" controls-position="right" /></label></div></section>
          </template>
          <template v-else-if="activeSection === 'orders'">
            <section class="personal-config__section"><h3>我的订单</h3><div class="personal-config__row"><span>是否显示</span><ElRadioGroup v-model="config.showOrders"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div><div class="personal-config__row"><span>模块标题</span><ElInput v-model="config.orderTitle" maxlength="8" /></div><div class="personal-config__row"><span>右侧文案</span><ElInput v-model="config.orderMoreText" maxlength="8" /></div></section>
            <section class="personal-config__section"><h3>订单状态</h3><p class="personal-config__hint">拖动可改变顺序；关闭后不在订单模块中展示。</p><div class="personal-config__sortable"><div v-for="(item, index) in config.orderItems" :key="item.label" draggable="true" @dragstart="dragOrderStart(index)" @dragover.prevent @drop="dropOrder(index)"><IconifyIcon icon="lucide:grip-vertical" /><ElSwitch v-model="item.enabled" /><IconifyIcon :icon="item.icon" /><ElInput v-model="item.label" maxlength="8" /></div></div></section>
          </template>
          <template v-else-if="activeSection === 'points'">
            <section class="personal-config__section">
              <h3>广告组</h3>
              <div class="personal-config__row">
                <span>是否显示</span>
                <ElRadioGroup v-model="config.showPointsBanner"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup>
              </div>
              <p class="personal-config__hint">最多添加 5 个广告项；预览区域展示为可切换的轮播。</p>
              <article v-for="(slide, index) in config.pointsSlides" :key="`${index}-${slide.title}`" class="personal-config__ad-card" :class="{ 'is-active': index === activePointsSlide }" @click="activePointsSlide = index">
                <header>
                  <span>广告 {{ index + 1 }}</span>
                  <button v-if="config.pointsSlides.length > 1" type="button" @click.stop="removePointsSlide(index)">删除</button>
                </header>
                <div class="personal-config__ad-card__image-row">
                  <button class="personal-config__ad-image" type="button" @click.stop="openImagePicker('pointsSlide', index)">
                    <img v-if="slide.image" :src="slide.image" :alt="`广告 ${index + 1}`" />
                    <IconifyIcon v-else icon="lucide:image-plus" />
                    <span>{{ slide.image ? '更换图片' : '选择图片' }}</span>
                  </button>
                  <button v-if="slide.image" class="personal-config__clear-image" type="button" @click.stop="slide.image = ''">清除</button>
                </div>
                <label><span>主标题</span><ElInput v-model="slide.title" maxlength="12" /></label>
                <label><span>副标题</span><ElInput v-model="slide.subtitle" maxlength="18" /></label>
                <label><span>按钮文案</span><ElInput v-model="slide.actionText" maxlength="14" /></label>
                <label><span>跳转链接</span><ElInput v-model="slide.link" /></label>
              </article>
              <ElButton plain :disabled="config.pointsSlides.length >= 5" @click="addPointsSlide">添加广告</ElButton>
              <div class="personal-config__row"><span>背景颜色</span><ElInput v-model="config.pointsBackground" /></div>
            </section>
          </template>
          <template v-else-if="activeSection === 'services'">
            <template v-if="activeConfigTab === 'content'">
              <section class="personal-config__section personal-config__section--settings"><h3>展示设置</h3><div class="personal-config__row"><span>导航样式</span><ElRadioGroup v-model="config.serviceNavigationType"><ElRadio value="icon-text">图标加文字</ElRadio><ElRadio value="icon">图标</ElRadio><ElRadio value="text">文字</ElRadio></ElRadioGroup></div><div class="personal-config__row"><span>单行显示</span><ElRadioGroup v-model="config.serviceRows"><ElRadio :value="3">3个</ElRadio><ElRadio :value="4">4个</ElRadio><ElRadio :value="5">5个</ElRadio></ElRadioGroup></div><div class="personal-config__row"><span>展示样式</span><ElRadioGroup v-model="config.serviceDisplayMode"><ElRadio value="fixed">固定显示</ElRadio><ElRadio value="page">分页滑动</ElRadio></ElRadioGroup></div><div class="personal-config__row personal-config__row--slider"><span>显示数量</span><ElSlider v-model="config.serviceCount" :min="4" :max="13" :show-tooltip="false" /><ElInputNumber v-model="config.serviceCount" :min="4" :max="13" controls-position="right" /></div></section>
              <section class="personal-config__section"><h3>内容设置</h3><p class="personal-config__hint">最多可添加13个服务入口；新增、编辑、启用和排序都会立即同步到预览</p><div class="personal-config__cards"><article v-for="(item, index) in config.serviceList" :key="item.value" class="personal-config__card" draggable="true" @dragstart="dragStart('service', index)" @dragover.prevent @drop="dropItem('service', index)"><IconifyIcon class="personal-config__drag" icon="lucide:grip-vertical" /><button class="personal-config__remove" type="button" @click="removeService(index)"><IconifyIcon icon="lucide:x" /></button><button class="personal-config__icon personal-config__icon--picker" type="button" @click="openImagePicker('serviceIcon', index)"><img v-if="item.image" :src="item.image" :alt="item.label" /><IconifyIcon v-else :icon="item.icon" /><small>选择图片</small></button><div class="personal-config__fields"><label>标题<ElInput v-model="item.label" maxlength="6" /></label><label>链接<ElInput v-model="item.link"><template #suffix><IconifyIcon icon="lucide:link" /></template></ElInput></label><label>状态<ElSwitch v-model="item.enabled" /></label></div></article></div><ElButton plain type="primary" @click="addService">+ 添加</ElButton></section>
            </template>
            <template v-else>
              <section class="personal-config__section personal-config__section--settings"><h3>图标样式</h3><div class="personal-config__row"><span>背景圆角</span><ElRadioGroup v-model="config.serviceIconRadiusMode"><ElRadio value="all">全部</ElRadio><ElRadio value="individual">单个</ElRadio></ElRadioGroup></div><template v-if="config.serviceIconRadiusMode === 'all'"><div class="personal-config__row personal-config__row--slider"><span>圆角值</span><ElSlider v-model="config.serviceIconRadius" :max="48" :show-tooltip="false" /><ElInputNumber v-model="config.serviceIconRadius" :min="0" :max="48" controls-position="right" /></div></template><div v-else class="personal-config__corner-grid"><label>左上<ElInputNumber v-model="config.serviceIconTopLeftRadius" :min="0" :max="48" /></label><label>右上<ElInputNumber v-model="config.serviceIconTopRightRadius" :min="0" :max="48" /></label><label>右下<ElInputNumber v-model="config.serviceIconBottomRightRadius" :min="0" :max="48" /></label><label>左下<ElInputNumber v-model="config.serviceIconBottomLeftRadius" :min="0" :max="48" /></label></div><div class="personal-config__row"><span>开启阴影</span><ElRadioGroup v-model="config.serviceIconShadow"><ElRadio value="off">关闭</ElRadio><ElRadio value="on">开启</ElRadio></ElRadioGroup></div></section>
              <section class="personal-config__section personal-config__section--settings"><h3>卡片样式</h3><div class="personal-config__row personal-config__row--slider"><span>组件上浮</span><ElSlider v-model="config.serviceFloat" :max="48" :show-tooltip="false" /><ElInputNumber v-model="config.serviceFloat" :min="0" :max="48" controls-position="right" /></div><div class="personal-config__color-row"><span>组件背景</span><ElInput v-model="config.serviceBackground" /><button type="button" @click="resetServiceStyle('serviceBackground')">重置</button></div><div class="personal-config__color-row"><span>背景渐变</span><ElInput v-model="config.serviceBackgroundEnd" /><button type="button" @click="resetServiceStyle('serviceBackgroundEnd')">重置</button></div><div class="personal-config__color-row"><span>底部背景</span><ElInput v-model="config.serviceBottomBackground" /><button type="button" @click="resetServiceStyle('serviceBottomBackground')">重置</button></div><div class="personal-config__color-row"><span>文字颜色</span><ElInput v-model="config.serviceTextColor" /><button type="button" @click="resetServiceStyle('serviceTextColor')">重置</button></div><div class="personal-config__row personal-config__row--slider"><span>上边距</span><ElSlider v-model="config.servicePaddingTop" :max="48" :show-tooltip="false" /><ElInputNumber v-model="config.servicePaddingTop" :min="0" :max="48" controls-position="right" /></div><div class="personal-config__row personal-config__row--slider"><span>下边距</span><ElSlider v-model="config.servicePaddingBottom" :max="48" :show-tooltip="false" /><ElInputNumber v-model="config.servicePaddingBottom" :min="0" :max="48" controls-position="right" /></div><div class="personal-config__row personal-config__row--slider"><span>左右边距</span><ElSlider v-model="config.servicePaddingLeft" :max="48" :show-tooltip="false" /><ElInputNumber v-model="config.servicePaddingLeft" :min="0" :max="48" controls-position="right" /></div><div class="personal-config__row personal-config__row--slider"><span>页面上间距</span><ElSlider v-model="config.serviceMarginTop" :max="96" :show-tooltip="false" /><ElInputNumber v-model="config.serviceMarginTop" :min="0" :max="96" controls-position="right" /></div><div class="personal-config__row"><span>背景圆角</span><ElRadioGroup v-model="config.serviceRadiusMode"><ElRadio value="all">全部</ElRadio><ElRadio value="individual">单个</ElRadio></ElRadioGroup></div><template v-if="config.serviceRadiusMode === 'all'"><div class="personal-config__row personal-config__row--slider"><span>圆角值</span><ElSlider v-model="config.serviceRadius" :max="48" :show-tooltip="false" /><ElInputNumber v-model="config.serviceRadius" :min="0" :max="48" controls-position="right" /></div></template><div v-else class="personal-config__corner-grid"><label>左上<ElInputNumber v-model="config.serviceTopLeftRadius" :min="0" :max="48" /></label><label>右上<ElInputNumber v-model="config.serviceTopRightRadius" :min="0" :max="48" /></label><label>右下<ElInputNumber v-model="config.serviceBottomRightRadius" :min="0" :max="48" /></label><label>左下<ElInputNumber v-model="config.serviceBottomLeftRadius" :min="0" :max="48" /></label></div><div class="personal-config__row"><span>开启阴影</span><ElRadioGroup v-model="config.serviceShadow"><ElRadio value="off">关闭</ElRadio><ElRadio value="on">开启</ElRadio></ElRadioGroup></div></section>
            </template>
          </template>
          <template v-else-if="activeSection === 'generic'">
            <section v-if="activeGenericComponent" class="personal-config__section personal-config__section--settings">
              <h3>{{ activeGenericComponent.label }}内容</h3>
              <div class="personal-config__row"><span>是否显示</span><ElRadioGroup v-model="activeGenericComponent.enabled"><ElRadio :value="true">显示</ElRadio><ElRadio :value="false">隐藏</ElRadio></ElRadioGroup></div>
              <template v-if="activeGenericComponent.type === 'banner'">
                <p class="personal-config__hint">最多添加 5 张图片，可拖动排序；图片点击后按配置链接跳转。</p>
                <article v-for="(banner, index) in activeGenericComponent.bannerImages" :key="index" class="personal-config__banner-card">
                  <IconifyIcon icon="lucide:grip-vertical" /><button type="button" class="personal-config__banner-image" @click="openImagePicker('genericBanner', undefined, activeGenericComponent.id, index)"><img v-if="banner.image" :src="banner.image" alt="轮播图片" /><IconifyIcon v-else icon="lucide:image-plus" /><span>{{ banner.image ? '更换图片' : '选择图片' }}</span></button><ElInput v-model="banner.link" placeholder="点击设置跳转链接" /><button type="button" class="personal-config__clear-image" @click="removeBannerImage(activeGenericComponent, index)">删除</button>
                </article>
                <ElButton plain type="primary" @click="addBannerImage(activeGenericComponent)">添加图片</ElButton>
              </template>
              <template v-else-if="activeGenericComponent.type === 'product'">
                <div class="personal-config__row"><span>模块标题</span><ElInput v-model="activeGenericComponent.title" maxlength="16" /></div>
                <div class="personal-config__row"><span>商品来源</span><ElRadioGroup v-model="activeGenericComponent.productSource"><ElRadio value="all">全部商品</ElRadio><ElRadio value="hot">热门商品</ElRadio><ElRadio value="new">新品商品</ElRadio></ElRadioGroup></div>
                <div class="personal-config__row"><span>展示样式</span><ElRadioGroup v-model="activeGenericComponent.productLayout"><ElRadio value="waterfall">双列</ElRadio><ElRadio value="list">列表</ElRadio><ElRadio value="big">大图</ElRadio></ElRadioGroup></div>
                <div class="personal-config__row personal-config__row--slider"><span>显示数量</span><ElSlider v-model="activeGenericComponent.productCount" :min="1" :max="20" :show-tooltip="false" /><ElInputNumber v-model="activeGenericComponent.productCount" :min="1" :max="20" controls-position="right" /></div>
                <div class="personal-config__row"><span>每行数量</span><ElRadioGroup v-model="activeGenericComponent.productColumns"><ElRadio :value="1">1 个</ElRadio><ElRadio :value="2">2 个</ElRadio><ElRadio :value="3">3 个</ElRadio></ElRadioGroup></div>
              </template>
              <template v-else><div class="personal-config__row"><span>模块标题</span><ElInput v-model="activeGenericComponent.title" maxlength="16" /></div><div class="personal-config__row"><span>内容文案</span><ElInput v-model="activeGenericComponent.description" type="textarea" :rows="3" maxlength="50" show-word-limit /></div><div class="personal-config__row"><span>跳转链接</span><ElInput v-model="activeGenericComponent.link" placeholder="/pages/..." /></div></template>
            </section>
          </template>
          <template v-else>
            <section class="personal-config__section"><h3>底部菜单</h3><p class="personal-config__hint">最多显示4项，拖动排序，勾选控制是否显示</p><div class="personal-config__cards personal-config__cards--compact"><article v-for="(item, index) in config.bottomList" :key="item.value" class="personal-config__card" draggable="true" @dragstart="dragStart('bottom', index)" @dragover.prevent @drop="dropItem('bottom', index)"><IconifyIcon class="personal-config__drag" icon="lucide:grip-vertical" /><div class="personal-config__icon"><IconifyIcon :icon="item.icon" /></div><div class="personal-config__fields"><label>标题<ElInput v-model="item.label" maxlength="4" /></label><label>链接<ElInput v-model="item.link"><template #suffix><IconifyIcon icon="lucide:link" /></template></ElInput></label><label>状态<ElSwitch v-model="item.enabled" /></label></div></article></div></section>
          </template>
        </div>
      </aside>
    </div>
    <ImagePickerDialog
      v-model:open="imagePickerOpen"
      default-library="system"
      @select="selectImageAsset"
    />
    <footer class="personal-editor__footer"><ElButton :loading="saving" type="primary" size="large" @click="save">保存</ElButton></footer>
  </Page>
</template>

<style scoped>
/* 预览区仅用于选中组件；内部控件不响应操作事件。 */
.personal-preview__module button { pointer-events:none; }
.personal-decoration-page { position:relative; height:calc(100vh - 130px); min-height:640px; overflow:hidden; font-family:'PingFang SC','Microsoft YaHei',sans-serif; }
.personal-editor { display:grid; grid-template-columns:270px minmax(0,1fr) 450px; height:calc(100% - 72px); min-height:568px; background:#f1f3f6; }
.personal-editor__library { min-width:0; padding:0 18px; border-right:1px solid #e7e9ee; background:#fff; overflow-y:auto; }.personal-library__header { padding:28px 8px 18px; border-bottom:1px solid #eef0f4; }.personal-library__header h2 { margin:0 0 26px; color:#303133; font-size:22px; font-weight:600; line-height:1; }.personal-library__header span { color:#8b9098; font-size:15px; }.personal-library__modules { display:grid; gap:10px; padding:18px 0; }.personal-library__modules button { display:flex; gap:16px; align-items:center; width:100%; height:68px; padding:0 18px; border:0; border-radius:8px; color:#515a6e; background:#fff; font-size:18px; text-align:left; cursor:pointer; }.personal-library__modules button:hover,.personal-library__modules button.is-active { color:#2468f2; background:#eaf2ff; }.personal-library__icon { display:grid; width:24px; height:24px; place-items:center; }.personal-library__icon svg { width:23px; height:23px; }
.personal-editor__canvas { min-width:0; overflow:auto; padding:24px 72px 38px; }.personal-phone { width:375px; margin:0 auto; overflow:visible; background:#f5f5f5; box-shadow:0 10px 28px rgb(30 44 74 / 10%); }.personal-phone__title { height:52px; display:grid; place-items:center; background:#fff; color:#222; font-size:19px; font-weight:600; }.personal-phone__screen { position:relative; min-height:667px; padding-bottom:66px; background:#f5f5f5; }.personal-preview__module { position:relative; border:2px solid transparent; cursor:pointer; }.personal-preview__module.is-active { border-color:#2f80ed; }.personal-preview__tip { position:absolute; z-index:8; top:12px; left:-102px; min-width:75px; padding:8px 10px; border:0; border-radius:4px 0 0 4px; color:#54637d; background:#fff; font-size:12px; white-space:nowrap; cursor:pointer; }.personal-preview__tip::after { position:absolute; top:50%; right:-8px; border-top:8px solid transparent; border-bottom:8px solid transparent; border-left:8px solid #fff; content:''; transform:translateY(-50%); }.is-active > .personal-preview__tip { color:#fff; background:#4073fa; }.is-active > .personal-preview__tip::after { border-left-color:#4073fa; }.personal-preview__header { min-height:196px; padding:20px 18px 0; color:#fff; background:#ff4a24; }.personal-preview__profile { display:flex; align-items:center; }.personal-preview__avatar { display:grid; width:54px; height:54px; place-items:center; border:2px solid rgb(255 255 255 / 50%); border-radius:50%; background:#fff1df; color:#555; }.personal-preview__avatar svg { width:31px; height:31px; }.personal-preview__profile button { padding:0; margin-left:12px; border:0; color:#fff; background:transparent; font-size:18px; font-weight:600; }.personal-preview__chat { margin-left:auto; font-size:25px; }.personal-preview__stats { display:flex; justify-content:space-between; margin-top:17px; text-align:center; }.personal-preview__stats span { display:flex; flex:1; flex-direction:column; gap:4px; font-size:11px; }.personal-preview__stats b { font-size:20px; font-weight:500; }.personal-preview__vip { display:flex; align-items:center; gap:8px; width:100%; height:44px; margin-top:20px; padding:0 12px; border:0; border-radius:17px 17px 0 0; color:#9c650f; background:#f7d788; font-size:13px; cursor:pointer; }.personal-preview__vip svg { color:#292929; font-size:23px; }.personal-preview__vip b { margin-left:auto; padding:7px 11px; border-radius:15px; color:#fff; background:#2d2d2d; font-size:11px; }.personal-preview__orders,.personal-preview__services { margin:0 16px 10px; border-radius:14px; background:#fff; overflow:visible; }.personal-preview__orders header { display:flex; align-items:center; justify-content:space-between; height:52px; padding:0 14px; border-bottom:1px dashed #e6e6e6; }.personal-preview__orders header b,.personal-preview__services h3 { font-size:16px; }.personal-preview__orders header button { display:flex; align-items:center; gap:2px; padding:0; border:0; color:#777; background:transparent; font-size:13px; }.personal-preview__orders > div { display:grid; grid-template-columns:repeat(5,1fr); padding:17px 7px; }.personal-preview__orders > div button,.personal-preview__services > div button { display:flex; flex-direction:column; gap:7px; align-items:center; padding:0; border:0; color:#333; background:transparent; font-size:11px; cursor:pointer; }.personal-preview__orders > div svg { color:#f34231; font-size:26px; }.personal-preview__points { display:flex; height:74px; flex-direction:column; justify-content:center; margin:0 16px 10px; padding:0 20px; border:0; border-radius:14px; color:#fff; background:#ff5b29; text-align:left; overflow:hidden; }.personal-preview__points b { font-size:20px; }.personal-preview__points small { margin-top:2px; font-size:12px; }.personal-preview__points em { width:max-content; margin-top:5px; padding:3px 12px; border-radius:10px; color:#e74c21; background:#ffe18a; font-size:9px; font-style:normal; }.personal-preview__points .personal-preview__tip { top:11px; }.personal-preview__services { padding-bottom:14px; }.personal-preview__services h3 { height:53px; margin:0; padding:18px 14px; border-bottom:1px dashed #e6e6e6; }.personal-preview__services > div { display:grid; grid-template-columns:repeat(4,1fr); gap:20px 6px; padding:20px 12px 4px; }.personal-preview__services > div svg { width:34px; height:34px; padding:8px; border-radius:50%; color:#fff; background:#5d92f7; }.personal-preview__services > div button:nth-child(4n + 2) svg { background:#8974f4; }.personal-preview__services > div button:nth-child(4n + 3) svg { background:#f15d5c; }.personal-preview__services > div button:nth-child(4n) svg { background:#f8ae27; }.personal-preview__bottom { position:absolute; right:0; bottom:0; left:0; z-index:3; display:flex; height:66px; align-items:center; border-top:1px solid #eee; background:#fff; }.personal-preview__bottom > button:not(.personal-preview__tip) { display:flex; flex:1; flex-direction:column; gap:4px; align-items:center; justify-content:center; padding:0; border:0; color:#343434; background:transparent; font-size:11px; cursor:pointer; }.personal-preview__bottom > button svg { width:23px; height:23px; }.personal-preview__bottom .is-current { color:#f14027; }.personal-preview__notice { position:absolute; z-index:15; right:30px; bottom:76px; left:30px; padding:8px; border-radius:18px; color:#fff; background:rgb(0 0 0 / 72%); font-size:12px; text-align:center; }.personal-notice-enter-active,.personal-notice-leave-active { transition:all .18s; }.personal-notice-enter-from,.personal-notice-leave-to { opacity:0; transform:translateY(6px); }
.personal-editor__config { min-width:0; overflow:hidden; border-left:1px solid #e7e9ee; background:#fff; }.personal-config__title { display:flex; height:60px; align-items:center; padding:0 20px; border-bottom:1px solid #ebeef5; color:#303133; font-size:16px; font-weight:600; }.personal-config__title > span { margin-right:auto; }.personal-config__title > button { min-width:66px; height:34px; border:0; color:#303133; background:#f5f6f8; font-size:14px; cursor:pointer; }.personal-config__title > button:first-of-type { border-radius:18px 0 0 18px; }.personal-config__title > button:last-child { border-radius:0 18px 18px 0; }.personal-config__title > button.active { color:#fff; background:#4073fa; }.personal-config__scroll { height:calc(100% - 60px); overflow-y:auto; background:#f0f2f5; }.personal-config__section { margin-bottom:4px; padding:22px 20px; background:#fff; }.personal-config__section h3 { margin:0 0 18px; color:#303133; font-size:15px; font-weight:500; }.personal-config__section--settings { padding-bottom:18px; }.personal-config__row { display:flex; gap:16px; align-items:center; margin:14px 0; }.personal-config__row > span { flex:0 0 70px; color:#999; font-size:14px; }.personal-config__row :deep(.el-input) { flex:1; }.personal-config__row :deep(.el-radio) { margin-right:18px; }.personal-config__row :deep(.el-radio__label) { font-size:13px; }.personal-config__row--slider :deep(.el-slider) { flex:1; }.personal-config__row--slider :deep(.el-input-number) { width:92px; }.personal-config__list { display:grid; flex:1; gap:8px; width:100%; }.personal-config__item { display:flex; height:40px; align-items:center; gap:9px; padding:0 10px; border-radius:4px; color:#303133; background:#fafafa; cursor:grab; }.personal-config__item:active { cursor:grabbing; }.personal-config__item > svg { color:#d5dbe3; font-size:18px; }.personal-config__item :deep(.el-checkbox) { margin:0; }.personal-config__hint { margin:0 0 14px; color:#9ca3af; font-size:12px; line-height:18px; }.personal-config__cards { display:grid; gap:12px; margin-bottom:16px; }.personal-config__card { position:relative; display:flex; min-height:128px; align-items:center; padding:18px 22px 18px 34px; border-radius:6px; background:#fafafa; cursor:grab; }.personal-config__card:active { cursor:grabbing; }.personal-config__cards--compact .personal-config__card { min-height:112px; }.personal-config__drag { position:absolute; left:8px; top:50%; color:#c7cfdb; font-size:20px; transform:translateY(-50%); }.personal-config__remove { position:absolute; top:-8px; right:-8px; display:grid; width:24px; height:24px; place-items:center; padding:0; border:0; border-radius:50%; color:#fff; background:#c8ccd4; cursor:pointer; }.personal-config__icon { display:grid; width:64px; height:64px; flex:0 0 64px; place-items:center; margin-right:16px; border-radius:50%; color:#fff; background:#5d92f7; }.personal-config__icon svg { width:32px; height:32px; }.personal-config__card:nth-child(4n + 2) .personal-config__icon { background:#8974f4; }.personal-config__card:nth-child(4n + 3) .personal-config__icon { background:#f15d5c; }.personal-config__card:nth-child(4n) .personal-config__icon { background:#f8ae27; }.personal-config__fields { display:grid; flex:1; gap:8px; min-width:0; }.personal-config__fields label { display:grid; grid-template-columns:42px minmax(0,1fr); gap:8px; align-items:center; color:#999; font-size:12px; }.personal-config__fields label:last-child { grid-template-columns:42px auto; }.personal-config__fields :deep(.el-input__wrapper) { min-height:28px; padding:1px 8px; }.personal-editor__footer { position:absolute; right:0; bottom:0; left:0; z-index:50; display:flex; height:72px; align-items:center; justify-content:center; border-top:1px solid #e7e9ed; background:#fff; box-shadow:0 -5px 16px rgb(30 45 70 / 4%); }.personal-editor__footer :deep(.el-button) { min-width:112px; }
.personal-config__color-row { display:grid; grid-template-columns:70px minmax(0,1fr) auto; gap:12px; align-items:center; margin:14px 0; }.personal-config__color-row > span { color:#999; font-size:14px; }.personal-config__color-row button { border:0; padding:0; color:#4073fa; background:transparent; font-size:13px; cursor:pointer; }.personal-config__corner-grid { display:grid; grid-template-columns:repeat(2, minmax(0, 1fr)); gap:10px; margin:12px 0 16px 86px; }.personal-config__corner-grid label { display:grid; grid-template-columns:32px minmax(0,1fr); align-items:center; gap:5px; color:#999; font-size:12px; }.personal-config__corner-grid :deep(.el-input-number) { width:100%; }.personal-editor__footer { position:absolute; right:0; bottom:0; left:0; z-index:50; display:flex; height:72px; align-items:center; justify-content:center; border-top:1px solid #e7e9ed; background:#fff; box-shadow:0 -5px 16px rgb(30 45 70 / 4%); }.personal-editor__footer :deep(.el-button) { min-width:112px; }
.personal-preview__avatar img { width:100%; height:100%; object-fit:cover; border-radius:50%; }
.personal-preview__service-image { width:34px; height:34px; padding:0 !important; object-fit:cover; }
.personal-config__image-select { display:flex; width:100%; min-height:38px; flex:1; align-items:center; gap:8px; padding:4px 10px; border:1px solid #dcdfe6; border-radius:4px; color:#606266; background:#fff; cursor:pointer; }.personal-config__image-select:hover { border-color:#4073fa; color:#4073fa; }.personal-config__image-select img { width:28px; height:28px; border-radius:4px; object-fit:cover; }.personal-config__image-select svg { width:20px; height:20px; }.personal-config__clear-image { flex:0 0 auto; padding:0; border:0; color:#4073fa; background:transparent; cursor:pointer; }.personal-config__icon--picker { position:relative; padding:0; border:0; cursor:pointer; overflow:visible; }.personal-config__icon--picker img { width:100%; height:100%; border-radius:inherit; object-fit:cover; }.personal-config__icon--picker small { position:absolute; right:-5px; bottom:-15px; width:74px; color:#4073fa; font-size:10px; font-weight:400; }
.personal-library__header { padding:24px 8px 14px; }
.personal-library__header h2 { margin-bottom:7px; }
.personal-library__group { padding:16px 0 2px; border-bottom:1px solid #eef0f4; }
.personal-library__group:last-child { border-bottom:0; }
.personal-library__group h3 { margin:0 0 10px 8px; color:#303133; font-size:15px; font-weight:600; }
.personal-library__grid { display:grid; grid-template-columns:repeat(3, minmax(0,1fr)); gap:8px; }
.personal-library__grid button { display:flex; height:82px; min-width:0; flex-direction:column; gap:8px; align-items:center; justify-content:center; padding:6px 2px; border:1px solid #e4e7ed; border-radius:8px; color:#343a46; background:#fff; font-size:12px; line-height:18px; cursor:pointer; }
.personal-library__grid button:hover,.personal-library__grid button.is-active { border-color:#4b7cff; color:#2468f2; background:#eef4ff; }
.personal-library__grid button svg { width:25px; height:25px; color:#2468f2; }
.personal-library__grid button span { max-width:100%; overflow:hidden; text-align:center; text-overflow:ellipsis; white-space:nowrap; }
.personal-preview__generic { min-height:96px; margin:0 16px 10px; border-radius:12px; color:#3a4354; background:#fff; overflow:visible; }
.personal-preview__generic > header { display:flex; min-height:42px; align-items:center; gap:8px; padding:0 13px; border-bottom:1px solid #edf0f4; }
.personal-preview__generic > header svg { width:19px; height:19px; color:#4073fa; }
.personal-preview__generic > header b { font-size:14px; }
.personal-preview__generic > div { display:flex; min-height:52px; align-items:center; gap:10px; padding:0 13px; }
.personal-preview__generic > div > svg:first-child { width:25px; height:25px; color:#4073fa; }
.personal-preview__generic span { flex:1; overflow:hidden; color:#9aa3af; font-size:11px; line-height:17px; text-overflow:ellipsis; white-space:nowrap; }
.personal-preview__generic .personal-preview__generic-add { width:18px; height:18px; color:#8b96a8; }
.personal-preview__generic--banner { min-height:0; overflow:visible; }.personal-preview__banner { position:relative; height:144px; overflow:hidden; border-radius:10px; background:#edf1f7; }.personal-preview__banner > img { width:100%; height:100%; object-fit:cover; }.personal-preview__banner-placeholder { display:flex; height:100%; flex-direction:column; gap:8px; align-items:center; justify-content:center; color:#7b8aa0; font-size:12px; }.personal-preview__banner-placeholder svg { width:30px; height:30px; color:#4073fa; }.personal-preview__dots { position:absolute; right:0; bottom:8px; left:0; display:flex; gap:5px; justify-content:center; }.personal-preview__dots i { width:5px; height:5px; border-radius:50%; background:rgb(255 255 255 / 65%); }.personal-preview__dots i.active { width:14px; border-radius:4px; background:#fff; }
.personal-preview__generic--product { overflow:visible; }.personal-preview__generic--product > header { justify-content:space-between; }.personal-preview__generic--product > header span { color:#8c96a6; font-size:11px; }.personal-preview__product-grid { display:grid; gap:8px; padding:10px 12px 14px; }.personal-preview__product-grid article { min-width:0; overflow:hidden; }.personal-preview__product-image { display:grid; width:100%; aspect-ratio:1; place-items:center; border-radius:7px; color:#a2adbc; background:linear-gradient(145deg,#f3f5f8,#e1e5eb); }.personal-preview__product-image svg { width:28px; height:28px; }.personal-preview__product-grid b { display:block; overflow:hidden; margin-top:6px; color:#303642; font-size:11px; line-height:16px; text-overflow:ellipsis; white-space:nowrap; }.personal-preview__product-grid em { display:block; margin-top:2px; color:#f04438; font-size:12px; font-style:normal; font-weight:600; }
.personal-config__simple-list { display:grid; gap:9px; }.personal-config__simple-list label { display:grid; grid-template-columns:auto minmax(0,1fr) 110px; gap:9px; align-items:center; padding:10px; border-radius:6px; background:#fafafa; }.personal-config__simple-list :deep(.el-input-number) { width:110px; }.personal-config__sortable { display:grid; gap:9px; }.personal-config__sortable > div { display:grid; grid-template-columns:20px auto 20px minmax(0,1fr); gap:10px; align-items:center; padding:10px; border-radius:6px; background:#fafafa; cursor:grab; }.personal-config__sortable > div:active { cursor:grabbing; }.personal-config__sortable > div > svg:first-child { color:#c7cfdb; font-size:18px; }.personal-config__sortable > div > svg:nth-child(3) { color:#5f7bc3; font-size:18px; }
.personal-config__color-row { grid-template-columns:70px 36px minmax(0,1fr) auto; }.personal-config__color-row :deep(.el-color-picker) { width:36px; }.personal-config__color-row :deep(.el-color-picker__trigger) { width:36px; height:36px; padding:4px; border-color:#dcdfe6; }.personal-config__shadow-settings { margin:8px 0 2px; padding:2px 0 0; border-top:1px solid #f0f2f5; }
.personal-config__banner-card { display:grid; grid-template-columns:18px 82px minmax(0,1fr) auto; gap:8px; align-items:center; margin:10px 0; padding:9px; border-radius:6px; background:#fafafa; }.personal-config__banner-card > svg { color:#c7cfdb; }.personal-config__banner-card :deep(.el-input) { min-width:0; }.personal-config__banner-image { position:relative; display:grid; width:82px; height:52px; place-items:center; overflow:hidden; border:1px dashed #c9d4e6; border-radius:5px; color:#4073fa; background:#f3f7ff; font-size:10px; cursor:pointer; }.personal-config__banner-image img { width:100%; height:100%; object-fit:cover; }.personal-config__banner-image span { position:absolute; right:0; bottom:0; left:0; padding:2px; color:#fff; background:rgb(0 0 0 / 45%); }
.personal-preview__points { position:relative; min-height:104px; height:auto; padding:12px 20px; overflow:visible; isolation:isolate; }.personal-preview__points::before { position:absolute; z-index:-1; inset:0; border-radius:inherit; background:linear-gradient(90deg,rgb(12 20 38 / 26%),transparent 72%); content:''; }.personal-preview__points b { overflow:hidden; font-size:20px; line-height:26px; text-overflow:ellipsis; white-space:nowrap; }.personal-preview__points small { overflow:hidden; margin-top:2px; font-size:12px; line-height:18px; text-overflow:ellipsis; white-space:nowrap; }.personal-preview__points-action { width:max-content; max-width:170px; margin-top:7px; padding:4px 12px; border:0; border-radius:11px; color:#e74c21; background:#ffe18a; font-size:10px; line-height:14px; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; cursor:pointer; }.personal-preview__points .personal-preview__tip { z-index:12; top:14px; }.personal-preview__points-dots { position:absolute; z-index:4; right:13px; bottom:10px; display:flex; gap:5px; align-items:center; }.personal-preview__points-dots button { width:5px; height:5px; padding:0; border:0; border-radius:50%; background:rgb(255 255 255 / 60%); cursor:pointer; }.personal-preview__points-dots button.active { width:14px; border-radius:4px; background:#fff; }
.personal-config__ad-card { display:grid; gap:9px; margin:12px 0; padding:12px; border:1px solid #e9edf3; border-radius:8px; background:#fafbfc; cursor:pointer; }.personal-config__ad-card.is-active { border-color:#4073fa; background:#f6f9ff; }.personal-config__ad-card header { display:flex; align-items:center; justify-content:space-between; color:#303133; font-size:13px; font-weight:600; }.personal-config__ad-card header button { padding:0; border:0; color:#f56c6c; background:transparent; cursor:pointer; }.personal-config__ad-card label { display:grid; grid-template-columns:58px minmax(0,1fr); gap:8px; align-items:center; color:#999; font-size:12px; }.personal-config__ad-card__image-row { display:flex; gap:8px; align-items:center; }.personal-config__ad-image { position:relative; display:grid; width:100%; height:72px; place-items:center; overflow:hidden; border:1px dashed #b9c9e6; border-radius:6px; color:#4073fa; background:#f3f7ff; font-size:12px; cursor:pointer; }.personal-config__ad-image img { position:absolute; inset:0; width:100%; height:100%; object-fit:cover; }.personal-config__ad-image span { position:relative; z-index:1; padding:3px 7px; border-radius:3px; color:#fff; background:rgb(0 0 0 / 48%); }.personal-config__ad-image > svg + span { color:#4073fa; background:transparent; }.personal-config__ad-image svg { width:21px; height:21px; }
@media (max-width:1200px) { .personal-editor { grid-template-columns:210px minmax(0,1fr) 370px; }.personal-editor__canvas { padding-right:28px; padding-left:28px; } }
</style>
