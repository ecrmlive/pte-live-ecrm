<script setup lang="ts">
import { computed, nextTick, reactive, ref } from 'vue';

import { Icon as IconifyIcon } from '@iconify/vue';
import { Page } from '@vben/common-ui';
import { ElButton, ElCheckbox, ElCheckboxGroup, ElInputNumber, ElMessage, ElRadio, ElRadioGroup, ElSlider } from 'element-plus';
import { ArrowLeft, MoreFilled, RefreshRight } from '@element-plus/icons-vue';

import { saveProductDetailDecorationApi } from '#/api/core/diy';

type AnchorKey =
  | 'product'
  | 'ranking'
  | 'coupon'
  | 'params'
  | 'package'
  | 'review'
  | 'show'
  | 'store'
  | 'bottom';

const activeAnchor = ref<AnchorKey>('product');
const activeImageIndex = ref(0);
const menuOpen = ref(true);
const saving = ref(false);
const favorite = ref(false);
const cartCount = ref(0);
const previewNotice = ref('');
const configScrollRef = ref<HTMLElement>();
const draggedShareIndex = ref<number>();
const draggedServiceIndex = ref<number>();
const draggedBottomIndex = ref<number>();

// 与 CRMEB product_detail_diy 保持同一份字段结构；预览与右侧表单只读写此对象。
const config = reactive({
  navList: [0, 1, 2, 3, 4],
  openShare: 1,
  pictureConfig: 0,
  swiperHeight: 750,
  swiperDot: 1,
  shareList: [
    { disabled: true, label: '客服', value: 0 },
    { label: '收藏', value: 1 },
    { label: '分享', value: 2 },
  ],
  shareConfig: [1, 2],
  isOpen: [0, 1, 2],
  showSvip: 1,
  showCoupon: 1,
  showRank: 1,
  serviceList: [
    { label: '规格选择', previewLabel: '请选择：', previewValue: '蓝色，2件', value: 0 },
    { label: '运费说明', previewLabel: '运费：', previewValue: '免运费', value: 1 },
    { label: '服务保障', previewLabel: '保障：', previewValue: '假一赔四 极速退款 七天无理由退换', value: 2 },
    { label: '参数说明', previewLabel: '参数：', previewValue: '品牌 型号...', value: 3 },
  ],
  showService: [0, 1, 2, 3],
  showMatch: 1,
  matchNum: 3,
  showReply: 1,
  replyNum: 3,
  showCommunity: 1,
  communityNum: 3,
  showStore: 1,
  showRecommend: 1,
  recommendNum: 12,
  bottomList: [
    { icon: 'lucide:store', label: '店铺', value: 0 },
    { icon: 'lucide:headphones', label: '客服', value: 1 },
    { icon: 'lucide:shopping-cart', label: '购物车', value: 2 },
    { icon: 'lucide:star', label: '收藏', value: 3 },
    { icon: 'lucide:share-2', label: '分享', value: 4 },
    { icon: 'lucide:house', label: '首页', value: 5 },
  ],
  menuList: [0, 1, 2],
  showCart: 1,
});

const anchors: Array<{ key: AnchorKey; label: string }> = [
  { key: 'product', label: '商品信息' },
  { key: 'ranking', label: '排行榜' },
  { key: 'coupon', label: '优惠券' },
  { key: 'params', label: '商品参数' },
  { key: 'package', label: '优惠套餐' },
  { key: 'review', label: '商品评价' },
  { key: 'show', label: '种草秀' },
  { key: 'store', label: '店铺信息' },
  { key: 'bottom', label: '底部菜单' },
];

// CRMEB 参考页面使用的预览素材，保持预览内容与原页面一致。
const heroImage = 'https://mer.crmeb.net/system/img/product_diy_main.2b1ea5e7.webp';
const productImageOne = 'https://mer.crmeb.net/system/img/product_diy_menu.844ce44e.webp';
const productImageTwo = 'https://mer.crmeb.net/system/img/product_diy_reply.33cdffe1.webp';
const communityImage = 'https://mer.crmeb.net/system/img/product_diy_plant.8f0380b4.webp';
const storeLogoImage = 'https://mer.crmeb.net/system/img/product_diy_merlogo.e94df212.jpg';
const productImages = [
  productImageOne,
  productImageTwo,
  communityImage,
];
const visibleReviewImages = computed(() => Array.from({ length: config.replyNum }, () => productImageTwo));
const visibleCommunityImages = computed(() => Array.from({ length: config.communityNum }, () => communityImage));
const visibleMatchImages = computed(() => Array.from({ length: config.matchNum }, (_, index) => productImages[index % productImages.length]!));
const visibleRecommendImages = computed(() => Array.from({ length: config.recommendNum }, (_, index) => productImages[index % productImages.length]!));
const visibleServiceItems = computed(() => config.serviceList.filter((item) => config.showService.includes(item.value)));
const visibleBottomItems = computed(() => config.bottomList.filter((item) => config.menuList.includes(item.value)));

const heroImages = [heroImage, productImageOne, productImageTwo];
const activeHeroImage = computed(() => heroImages[activeImageIndex.value] ?? heroImage);
const panelTitle = computed(() => anchors.find((item) => item.key === activeAnchor.value)?.label ?? '商品信息');
const activeModuleVisible = computed({
  get() {
    switch (activeAnchor.value) {
      case 'ranking': return config.showRank;
      case 'coupon': return config.showCoupon;
      case 'package': return config.showMatch;
      case 'review': return config.showReply;
      case 'show': return config.showCommunity;
      case 'store': return config.showStore;
      default: return 1;
    }
  },
  set(value: number) {
    switch (activeAnchor.value) {
      case 'ranking': config.showRank = value; break;
      case 'coupon': config.showCoupon = value; break;
      case 'package': config.showMatch = value; break;
      case 'review': config.showReply = value; break;
      case 'show': config.showCommunity = value; break;
      case 'store': config.showStore = value; break;
    }
  },
});
const previewMenuItems = computed(() => [
  config.navList.includes(0) ? { icon: 'lucide:house', label: '首页' } : null,
  config.navList.includes(1) ? { icon: 'lucide:search', label: '搜索' } : null,
  config.navList.includes(2) ? { icon: 'lucide:shopping-cart', label: '购物车' } : null,
  config.navList.includes(3) ? { icon: 'lucide:star', label: '我的收藏' } : null,
  config.navList.includes(4) ? { icon: 'lucide:user-round', label: '个人中心' } : null,
].filter(Boolean) as Array<{ icon: string; label: string }>);

async function selectSection(key: AnchorKey) {
  activeAnchor.value = key;
  await nextTick();
  configScrollRef.value?.scrollTo({ top: 0, behavior: 'smooth' });
}

function startShareDrag(index: number) {
  if (config.shareList[index]?.disabled) return;
  draggedShareIndex.value = index;
}

function dropShareItem(targetIndex: number) {
  const sourceIndex = draggedShareIndex.value;
  draggedShareIndex.value = undefined;
  if (sourceIndex === undefined || sourceIndex === targetIndex || config.shareList[targetIndex]?.disabled) return;
  const [item] = config.shareList.splice(sourceIndex, 1);
  if (item) config.shareList.splice(targetIndex, 0, item);
}

function startServiceDrag(index: number) {
  draggedServiceIndex.value = index;
}

function dropServiceItem(targetIndex: number) {
  const sourceIndex = draggedServiceIndex.value;
  draggedServiceIndex.value = undefined;
  if (sourceIndex === undefined || sourceIndex === targetIndex) return;
  const [item] = config.serviceList.splice(sourceIndex, 1);
  if (item) config.serviceList.splice(targetIndex, 0, item);
}

function startBottomDrag(index: number) {
  draggedBottomIndex.value = index;
}

function dropBottomItem(targetIndex: number) {
  const sourceIndex = draggedBottomIndex.value;
  draggedBottomIndex.value = undefined;
  if (sourceIndex === undefined || sourceIndex === targetIndex) return;
  const [item] = config.bottomList.splice(sourceIndex, 1);
  if (item) config.bottomList.splice(targetIndex, 0, item);
}

function isBottomItemDisabled(value: number) {
  return !config.menuList.includes(value) && config.menuList.length >= 3;
}

function nextHeroImage() {
  activeImageIndex.value = (activeImageIndex.value + 1) % heroImages.length;
}

let previewNoticeTimer: ReturnType<typeof setTimeout> | undefined;

function showPreviewNotice(message: string) {
  previewNotice.value = message;
  if (previewNoticeTimer) clearTimeout(previewNoticeTimer);
  previewNoticeTimer = setTimeout(() => {
    previewNotice.value = '';
  }, 1600);
}

function toggleFavorite() {
  favorite.value = !favorite.value;
  showPreviewNotice(favorite.value ? '已收藏商品' : '已取消收藏');
}

function addToCart() {
  cartCount.value += 1;
  showPreviewNotice(`已加入购物车，共 ${cartCount.value} 件`);
}

function handlePreviewMenuAction(label: string) {
  menuOpen.value = false;
  showPreviewNotice(`已打开${label}`);
}

async function save() {
  saving.value = true;
  try {
    await saveProductDetailDecorationApi(JSON.parse(JSON.stringify(config)));
    ElMessage.success('商品详情装修已保存');
  } catch {
    ElMessage.error('商品详情装修保存失败，请稍后重试');
  } finally {
    saving.value = false;
  }
}

</script>

<template>
  <Page auto-content-height content-class="detail-decoration-page !p-0">
    <div class="detail-editor">
      <main class="detail-editor__canvas">
        <section class="detail-phone" aria-label="商品详情预览">
          <div class="detail-phone__title">商品详情</div>
          <div class="detail-phone__screen">
            <section data-anchor="product" class="detail-preview__product detail-preview__module" :class="{ 'is-active': activeAnchor === 'product' }" @click="selectSection('product')">
              <button class="detail-preview__module-tip" :class="{ 'is-active': activeAnchor === 'product' }" type="button" @click.stop="selectSection('product')">商品信息</button>
              <div class="detail-preview__hero-wrap">
                <img
                  class="detail-preview__hero"
                  :class="{ 'is-auto-height': config.pictureConfig === 1 }"
                  :style="config.pictureConfig === 0 ? { height: `${Math.round(config.swiperHeight / 2)}px` } : undefined"
                  :src="activeHeroImage"
                  alt="商品主图"
                  @click="nextHeroImage"
                />
                <div class="detail-preview__topbar">
                  <div class="detail-preview__nav-control"><span class="detail-preview__back" aria-hidden="true"><ArrowLeft /></span><button v-if="previewMenuItems.length" class="detail-preview__menu" type="button" @click.stop="menuOpen = !menuOpen"><IconifyIcon icon="lucide:menu" /></button></div>
                  <div class="detail-preview__top-actions">
                    <RefreshRight v-if="config.navList.includes(1)" class="detail-preview__refresh" aria-hidden="true" />
                    <button v-if="config.openShare === 1" class="detail-preview__icon-action" type="button" aria-label="分享" @click.stop="showPreviewNotice('已打开分享面板')"><MoreFilled class="detail-preview__more" aria-hidden="true" /></button>
                    <IconifyIcon class="detail-preview__mini-dot" icon="lucide:circle-dot" aria-hidden="true" />
                  </div>
                </div>
                <div v-if="config.swiperDot === 1" class="detail-preview__dots">
                  <button v-for="(_, index) in heroImages" :key="index" type="button" :class="{ 'is-current': index === activeImageIndex }" @click="activeImageIndex = index" />
                </div>
                <div v-if="menuOpen" class="detail-preview__menu-popover">
                  <button v-for="item in previewMenuItems" :key="item.label" type="button" @click="handlePreviewMenuAction(item.label)"><IconifyIcon :icon="item.icon" />{{ item.label }}</button>
                </div>
              </div>
              <div class="detail-preview__product-card">
                <div class="detail-preview__price-row">
                  <strong>¥199<small>.00</small></strong>
                  <span v-if="config.isOpen.includes(0)">¥234.00</span>
                  <em v-if="config.showSvip === 1">SVIP</em>
                  <div class="detail-preview__collect-actions">
                    <button v-if="config.shareConfig.includes(1)" :class="{ 'is-favorited': favorite }" type="button" @click.stop="toggleFavorite"><IconifyIcon icon="lucide:star" />{{ favorite ? '已收藏' : '收藏' }}</button>
                    <button v-if="config.shareConfig.includes(2)" type="button" @click.stop="showPreviewNotice('已打开分享面板')"><IconifyIcon icon="lucide:share-2" />分享</button>
                  </div>
                </div>
                <button v-if="config.showSvip === 1" class="detail-preview__vip" type="button" @click.stop="showPreviewNotice('已打开 SVIP 会员入口')">开通SVIP会员预计省立省29元 <b>立即开通 ›</b></button>
                <div class="detail-preview__badges"><span>积分最高可抵扣14.2元</span><span>包邮</span></div>
                <h2>年货节立即抢购兰蔻唇香礼盒 是我香水口红289节日套装礼品</h2>
                <div class="detail-preview__stock"><s v-if="config.isOpen.includes(0)">¥234.00</s><span v-if="config.isOpen.includes(2)">库存:1425件</span><span v-if="config.isOpen.includes(1)">已售:2399件</span></div>
              </div>
            </section>

            <section v-if="config.showRank === 1" data-anchor="ranking" class="detail-preview__card detail-preview__ranking detail-preview__module" :class="{ 'is-active': activeAnchor === 'ranking' }" @click="selectSection('ranking')">
              <button class="detail-preview__module-tip" :class="{ 'is-active': activeAnchor === 'ranking' }" type="button" @click.stop="selectSection('ranking')">排行榜</button>
              <b>TOP 榜单</b><span>紧致抗皱套装·第2名</span><IconifyIcon icon="lucide:chevron-right" />
            </section>

            <section v-if="config.showCoupon === 1" data-anchor="coupon" class="detail-preview__card detail-preview__coupon detail-preview__module" :class="{ 'is-active': activeAnchor === 'coupon' }" @click="selectSection('coupon')">
              <button class="detail-preview__module-tip" :class="{ 'is-active': activeAnchor === 'coupon' }" type="button" @click.stop="selectSection('coupon')">优惠券</button>
              <span>优惠券：</span><b>满100减10</b><IconifyIcon icon="lucide:chevron-right" />
            </section>

            <section data-anchor="params" class="detail-preview__card detail-preview__params detail-preview__module" :class="{ 'is-active': activeAnchor === 'params' }" @click="selectSection('params')">
              <button class="detail-preview__module-tip" :class="{ 'is-active': activeAnchor === 'params' }" type="button" @click.stop="selectSection('params')">商品参数</button>
              <p v-for="item in visibleServiceItems" :key="item.value"><span>{{ item.previewLabel }}</span><b>{{ item.previewValue }}</b><IconifyIcon icon="lucide:chevron-right" /></p>
            </section>

            <section v-if="config.showMatch === 1" data-anchor="package" class="detail-preview__card detail-preview__package detail-preview__module" :class="{ 'is-active': activeAnchor === 'package' }" @click="selectSection('package')">
              <button class="detail-preview__module-tip" :class="{ 'is-active': activeAnchor === 'package' }" type="button" @click.stop="selectSection('package')">优惠套餐</button>
              <h3>优惠套餐({{ config.matchNum }})<IconifyIcon icon="lucide:chevron-right" /></h3>
              <div><img v-for="(image, index) in visibleMatchImages" :key="`${image}-${index}`" :src="image" alt="套餐商品" /><span>共{{ config.matchNum }}件<br /><em>省 ¥32.00</em></span></div>
            </section>

            <section v-if="config.showReply === 1" data-anchor="review" class="detail-preview__card detail-preview__review detail-preview__module" :class="{ 'is-active': activeAnchor === 'review' }" @click="selectSection('review')">
              <button class="detail-preview__module-tip" :class="{ 'is-active': activeAnchor === 'review' }" type="button" @click.stop="selectSection('review')">商品评价</button>
              <h3>用户评价(20)<span>99%<i>好评率</i> ›</span></h3>
              <p>👨🏻　 一*兔　 <b>★★★★★</b></p>
              <small>2019.01.04 17:11 蓝色*1盒</small>
              <p class="detail-preview__review-text">质量挺不错的，用了一段时间也不起球，会推荐朋友一起回购的~~~</p>
              <div><img v-for="(image, index) in visibleReviewImages" :key="`${image}-${index}`" :src="image" alt="用户评价图片" /></div>
            </section>

            <section v-if="config.showCommunity === 1" data-anchor="show" class="detail-preview__card detail-preview__gallery-card detail-preview__module" :class="{ 'is-active': activeAnchor === 'show' }" @click="selectSection('show')">
              <button class="detail-preview__module-tip" :class="{ 'is-active': activeAnchor === 'show' }" type="button" @click.stop="selectSection('show')">种草秀</button>
              <h3>种草秀(3)<span>查看全部 ›</span></h3>
              <div><img v-for="(image, index) in visibleCommunityImages" :key="`${image}-${index}`" :src="image" alt="种草秀图片" /></div>
            </section>

            <section v-if="config.showStore === 1" data-anchor="store" class="detail-preview__card detail-preview__store detail-preview__module" :class="{ 'is-active': activeAnchor === 'store' }" @click="selectSection('store')">
              <button class="detail-preview__module-tip" :class="{ 'is-active': activeAnchor === 'store' }" type="button" @click.stop="selectSection('store')">店铺信息</button>
              <div class="detail-preview__store-head"><img :src="storeLogoImage" alt="店铺" /><p><b>爱花屋</b><span>41.6万人关注</span></p><button type="button" @click.stop="showPreviewNotice('已进入爱花屋店铺')">进店</button></div>
              <div class="detail-preview__scores"><span>商品描述 <b>5.0</b></span><span>卖家服务 <b>5.0</b></span><span>物流服务 <b>5.0</b></span></div>
              <template v-if="config.showRecommend === 1"><h3>店铺推荐</h3><div class="detail-preview__store-products"><article v-for="(image, index) in visibleRecommendImages" :key="`${image}-${index}`"><img :src="image" alt="推荐商品" /><p>商品名称商品</p><b>¥199.00</b></article></div></template>
            </section>

            <section data-anchor="bottom" class="detail-preview__bottom-nav detail-preview__module" :class="{ 'is-active': activeAnchor === 'bottom' }" @click="selectSection('bottom')">
              <button class="detail-preview__module-tip" :class="{ 'is-active': activeAnchor === 'bottom' }" type="button" @click.stop="selectSection('bottom')">底部菜单</button>
              <button v-for="item in visibleBottomItems" :key="item.value" class="detail-preview__bottom-entry" type="button" @click.stop="showPreviewNotice(`已打开${item.label}`)"><IconifyIcon :icon="item.icon" />{{ item.label }}</button>
              <button v-if="config.showCart === 1" class="detail-preview__buy-action detail-preview__buy-action--cart" @click.stop="addToCart">加入购物车</button><button class="detail-preview__buy-action detail-preview__buy-action--now" @click.stop="showPreviewNotice('已进入立即购买流程')">立即购买</button>
            </section>
            <Transition name="preview-notice"><div v-if="previewNotice" class="detail-preview__notice" role="status">{{ previewNotice }}</div></Transition>
          </div>
        </section>
      </main>

      <aside class="detail-editor__config">
        <div class="detail-config__title">{{ panelTitle }}</div>
        <div ref="configScrollRef" class="detail-config__scroll">
          <template v-if="activeAnchor === 'product'">
          <section class="detail-config__section detail-config__section--nav">
            <h3>顶部导航</h3>
            <div class="detail-config__row"><span>菜单内容</span><ElCheckboxGroup v-model="config.navList" class="detail-config__checks"><ElCheckbox :value="0">首页</ElCheckbox><ElCheckbox :value="1">搜索</ElCheckbox><ElCheckbox :value="2">购物车</ElCheckbox><ElCheckbox :value="3">我的收藏</ElCheckbox><ElCheckbox :value="4">个人中心</ElCheckbox></ElCheckboxGroup></div>
            <div class="detail-config__row"><span>开启分享</span><ElRadioGroup v-model="config.openShare"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div>
          </section>
          <section class="detail-config__section detail-config__section--image">
            <h3>商品主图</h3>
            <div class="detail-config__row"><span>商品主图</span><ElRadioGroup v-model="config.pictureConfig"><ElRadio :value="0">固定方图</ElRadio><ElRadio :value="1">高度自适应</ElRadio></ElRadioGroup></div>
            <div class="detail-config__row"><span>高度设置px</span><ElInputNumber v-model="config.swiperHeight" :min="750" :max="1500" controls-position="right" /></div>
            <div class="detail-config__row"><span>轮播点</span><ElRadioGroup v-model="config.swiperDot"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div>
          </section>
          <section class="detail-config__section">
            <h3>收藏/分享</h3>
            <div class="detail-config__row detail-config__row--list"><span>是否显示</span><ElCheckboxGroup v-model="config.shareConfig" class="detail-config__draggable-list"><div v-for="(item, index) in config.shareList" :key="item.value" class="detail-config__drag-row" :class="{ 'is-disabled': item.disabled }" :draggable="!item.disabled" @dragstart="startShareDrag(index)" @dragover.prevent @drop="dropShareItem(index)"><IconifyIcon icon="lucide:grip-vertical" /><ElCheckbox :value="item.value" :disabled="item.disabled">{{ item.label }}</ElCheckbox></div></ElCheckboxGroup></div>
          </section>
          <section class="detail-config__section">
            <h3>商品信息</h3>
            <div class="detail-config__row detail-config__row--compact"><span>是否开启</span><ElCheckboxGroup v-model="config.isOpen" class="detail-config__checks detail-config__checks--inline"><ElCheckbox :value="0">划线价</ElCheckbox><ElCheckbox :value="2">库存</ElCheckbox><ElCheckbox :value="1">已售</ElCheckbox></ElCheckboxGroup></div>
          </section>
          <section class="detail-config__section">
            <h3>svip会员</h3>
            <div class="detail-config__row"><span>会员入口</span><ElRadioGroup v-model="config.showSvip"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div>
          </section>
          </template>
          <section v-else-if="activeAnchor === 'coupon'" class="detail-config__section">
            <h3>显示状态</h3>
            <div class="detail-config__row"><span>是否显示</span><ElRadioGroup v-model="config.showCoupon"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div>
            <p class="detail-config__hint">仅普通、预售商品有优惠券</p>
          </section>
          <section v-else-if="activeAnchor === 'ranking'" class="detail-config__section">
            <h3>显示状态</h3>
            <div class="detail-config__row"><span>是否显示</span><ElRadioGroup v-model="config.showRank"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div>
          </section>
          <section v-else-if="activeAnchor === 'params'" class="detail-config__section">
            <h3>显示状态</h3>
            <div class="detail-config__row detail-config__row--list">
              <span>是否显示</span>
              <ElCheckboxGroup v-model="config.showService" class="detail-config__draggable-list">
                <div v-for="(item, index) in config.serviceList" :key="item.value" class="detail-config__drag-row" :class="{ 'is-dragging': draggedServiceIndex === index }" draggable="true" @dragstart="startServiceDrag(index)" @dragend="draggedServiceIndex = undefined" @dragover.prevent @drop="dropServiceItem(index)">
                  <IconifyIcon icon="lucide:grip-vertical" />
                  <ElCheckbox :value="item.value">{{ item.label }}</ElCheckbox>
                </div>
              </ElCheckboxGroup>
            </div>
          </section>
          <section v-else-if="activeAnchor === 'package'" class="detail-config__section">
            <h3>显示状态</h3>
            <div class="detail-config__row"><span>是否显示</span><ElRadioGroup v-model="config.showMatch"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div>
            <div class="detail-config__row detail-config__row--slider"><span>显示数量</span><ElSlider v-model="config.matchNum" :min="1" :max="3" :show-tooltip="false" /><ElInputNumber v-model="config.matchNum" :min="1" :max="3" controls-position="right" /></div>
          </section>
          <section v-else-if="activeAnchor === 'review'" class="detail-config__section">
            <h3>用户评价</h3>
            <div class="detail-config__row"><span>是否显示</span><ElRadioGroup v-model="config.showReply"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div>
            <div class="detail-config__row"><span>展示数量</span><ElSlider v-model="config.replyNum" :min="1" :max="3" :show-tooltip="false" /></div>
          </section>
          <section v-else-if="activeAnchor === 'show'" class="detail-config__section">
            <h3>种草秀</h3>
            <div class="detail-config__row"><span>是否显示</span><ElRadioGroup v-model="config.showCommunity"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div>
            <div class="detail-config__row"><span>展示数量</span><ElSlider v-model="config.communityNum" :min="1" :max="3" :show-tooltip="false" /></div>
          </section>
          <template v-else-if="activeAnchor === 'store'">
            <section class="detail-config__section">
              <h3>显示状态</h3>
              <div class="detail-config__row"><span>是否显示</span><ElRadioGroup v-model="config.showStore"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div>
            </section>
            <section class="detail-config__section">
              <h3>店铺推荐</h3>
              <div class="detail-config__row"><span>是否显示</span><ElRadioGroup v-model="config.showRecommend"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div>
              <div class="detail-config__row detail-config__row--slider"><span>显示数量</span><ElSlider v-model="config.recommendNum" :min="3" :max="12" :step="3" :show-tooltip="false" /><ElInputNumber v-model="config.recommendNum" :min="3" :max="12" :step="3" controls-position="right" /></div>
            </section>
          </template>
          <section v-else-if="activeAnchor === 'bottom'" class="detail-config__section">
            <div class="detail-config__row detail-config__row--list">
              <span>是否显示</span>
              <ElCheckboxGroup v-model="config.menuList" :max="3" class="detail-config__draggable-list">
                <div v-for="(item, index) in config.bottomList" :key="item.value" class="detail-config__drag-row" :class="{ 'is-disabled': isBottomItemDisabled(item.value), 'is-dragging': draggedBottomIndex === index }" draggable="true" @dragstart="startBottomDrag(index)" @dragend="draggedBottomIndex = undefined" @dragover.prevent @drop="dropBottomItem(index)">
                  <IconifyIcon icon="lucide:grip-vertical" />
                  <ElCheckbox :value="item.value" :disabled="isBottomItemDisabled(item.value)">{{ item.label }}</ElCheckbox>
                </div>
              </ElCheckboxGroup>
            </div>
            <div class="detail-config__row"><span>购物车按钮</span><ElRadioGroup v-model="config.showCart"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div>
            <p class="detail-config__hint">活动商品不可加入购物车，此按钮不生效</p>
          </section>
          <section v-else class="detail-config__section">
            <h3>显示状态</h3>
            <div class="detail-config__row"><span>是否显示</span><ElRadioGroup v-model="activeModuleVisible"><ElRadio :value="1">显示</ElRadio><ElRadio :value="0">隐藏</ElRadio></ElRadioGroup></div>
          </section>
        </div>
      </aside>
    </div>
    <footer class="detail-editor__footer"><ElButton :loading="saving" type="primary" size="large" @click="save">保存</ElButton></footer>
  </Page>
</template>

<style scoped>
.detail-decoration-page { height: 100%; min-height: 0; overflow: hidden; }
.detail-editor { display: flex; min-height: 600px; height: 100%; padding: 0 400px 72px 0; background: #f0f2f5; font-family: 'PingFang SC', 'Microsoft YaHei', Arial, sans-serif; }
.detail-editor__canvas { flex: 1; min-width: 0; overflow: auto; padding: 0 190px 36px; }
.detail-phone { width: 375px; margin: 20px auto 0; overflow: visible; background: #f5f5f5; }
.detail-phone__title { height: 58px; display: grid; place-items: center; background: #fff; color: #222; font-size: 19px; font-weight: 600; }
.detail-phone__screen { overflow: visible; background: #f5f5f5; }
.detail-preview__module { position: relative; box-sizing: border-box; border: 2px solid #f5f5f5; cursor: pointer; }
.detail-preview__module.is-active { border-color: #1890ff; }
.detail-preview__module-tip { position: absolute; z-index: 8; top: 0; left: -100px; min-width: 72px; padding: 8px 12px; border: 0; border-radius: 4px 0 0 4px; background: #fff; color: #515a6e; font-size: 12px; line-height: 18px; text-align: center; white-space: nowrap; cursor: pointer; }
.detail-preview__module-tip::after { content: ''; position: absolute; top: 50%; right: -8px; border-top: 8px solid transparent; border-bottom: 8px solid transparent; border-left: 8px solid #fff; transform: translateY(-50%); }
.detail-preview__module-tip.is-active { background: #4073fa; color: #fff; }.detail-preview__module-tip.is-active::after { border-left-color: #4073fa; }
.detail-preview__hero-wrap { position: relative; overflow: hidden; }.detail-preview__hero { display: block; width: 100%; height: 375px; object-fit: cover; cursor: pointer; }.detail-preview__hero.is-auto-height { height: auto; min-height: 240px; }
.detail-preview__topbar { position: absolute; top: 12px; right: 12px; left: 12px; display: flex; align-items: center; justify-content: space-between; color: #111; }.detail-preview__nav-control,.detail-preview__top-actions { display: flex; align-items: center; height: 34px; border-radius: 20px; background: rgb(255 255 255 / 74%); }.detail-preview__back { display: flex; width: 32px; align-items: center; justify-content: center; color: #111; }.detail-preview__back :deep(svg) { width: 21px; height: 21px; stroke-width: 2.6; }.detail-preview__menu { display: flex; width: 38px; height: 22px; padding: 0; border: 0; border-left: 1px solid rgb(0 0 0 / 20%); background: transparent; align-items: center; justify-content: center; color: inherit; cursor: pointer; }.detail-preview__menu :deep(svg) { width: 21px; height: 21px; stroke-width: 2.6; }.detail-preview__top-actions { gap: 12px; padding: 0 11px; }.detail-preview__refresh { width: 20px; height: 20px; }.detail-preview__more { width: 19px; height: 19px; }.detail-preview__mini-dot { width: 19px; height: 19px; color: #111; }
.detail-preview__dots { position: absolute; right: 18px; bottom: 10px; left: 18px; display: flex; gap: 4px; justify-content: center; }.detail-preview__dots button { height: 3px; flex: 1; padding: 0; border: 0; background: rgb(255 255 255 / 45%); }.detail-preview__dots .is-current { background: #fff; }
.detail-preview__menu-popover { position: absolute; z-index: 9; top: 52px; left: 16px; width: 140px; overflow: hidden; border-radius: 8px; background: #fff; box-shadow: 0 8px 22px rgb(34 44 60 / 24%); }.detail-preview__menu-popover button { display: flex; gap: 10px; align-items: center; width: 100%; padding: 11px 14px; border: 0; border-bottom: 1px solid #f0f0f0; background: #fff; color: #343434; font-size: 13px; text-align: left; cursor: pointer; }
.detail-preview__product-card { position: relative; padding: 16px; border-radius: 16px 16px 0 0; background: #fff; }.detail-preview__price-row { display: flex; gap: 8px; align-items: center; }.detail-preview__price-row strong { color: #ef3e2e; font-size: 27px; line-height: 1; }.detail-preview__price-row strong small { font-size: 16px; }.detail-preview__price-row > span { color: #444; text-decoration: line-through; }.detail-preview__price-row > em { padding: 2px 5px; border-radius: 9px; background: #e8c77c; color: #4f3b0c; font-size: 10px; font-style: normal; }.detail-preview__collect-actions { display: flex; gap: 10px; margin-left: auto; }.detail-preview__collect-actions button { display: flex; flex-direction: column; gap: 2px; align-items: center; padding: 0; border: 0; background: transparent; color: #555; font-size: 10px; cursor: pointer; }.detail-preview__collect-actions button.is-favorited { color: #ef3e2e; }.detail-preview__vip { display: flex; justify-content: space-between; width: 100%; margin-top: 12px; padding: 10px; border: 0; border-radius: 7px; background: #fff0cf; color: #c77b05; font-size: 13px; cursor: pointer; }.detail-preview__badges { display: flex; gap: 8px; margin-top: 10px; }.detail-preview__badges span { padding: 4px 9px; border-radius: 13px; background: #fff3df; color: #ef8d13; font-size: 11px; }.detail-preview__product-card h2 { margin: 10px 0; font-size: 16px; line-height: 1.35; }.detail-preview__stock { display: flex; justify-content: space-between; color: #a0a0a0; font-size: 11px; }
.detail-preview__card { margin: 10px 10px 0; padding: 14px; border-radius: 14px; background: #fff; }.detail-preview__ranking { display: flex; gap: 8px; align-items: center; }.detail-preview__ranking b { padding: 5px 8px; border-radius: 3px; background: #232323; color: #f4d993; font-size: 12px; }.detail-preview__ranking svg,.detail-preview__coupon svg { margin-left: auto; color: #777; }.detail-preview__coupon { display: flex; align-items: center; color: #999; }.detail-preview__coupon b { margin-left: 8px; color: #f24937; font-weight: 500; }.detail-preview__params p { display: flex; align-items: center; margin: 0; padding: 12px 0; border-bottom: 1px solid #f2f2f2; font-size: 13px; }.detail-preview__params p:last-child { border-bottom: 0; }.detail-preview__params p span { width: 55px; color: #999; }.detail-preview__params p b { font-weight: 400; }.detail-preview__params p svg { margin-left: auto; color: #777; }.detail-preview__card h3 { display: flex; align-items: center; justify-content: space-between; margin: 0 0 12px; font-size: 18px; }.detail-preview__card h3 span { color: #999; font-size: 14px; font-weight: 400; }
.detail-preview__package > div,.detail-preview__review > div,.detail-preview__gallery-card > div { display: flex; gap: 8px; align-items: center; }.detail-preview__package img { width: 82px; height: 82px; border-radius: 8px; object-fit: cover; }.detail-preview__package > div span { margin-left: auto; color: #333; font-size: 13px; line-height: 1.8; }.detail-preview__package em { color: #f13a2e; font-size: 16px; font-style: normal; }.detail-preview__review p { margin: 8px 0; font-size: 14px; }.detail-preview__review p b { color: #f54435; letter-spacing: 1px; }.detail-preview__review small { color: #a0a0a0; }.detail-preview__review-text { line-height: 1.55; }.detail-preview__review img { width: calc((100% - 16px) / 3); aspect-ratio: 1; border-radius: 6px; object-fit: cover; }.detail-preview__gallery-card img { width: calc((100% - 16px) / 3); height: 108px; border-radius: 7px; object-fit: cover; }
.detail-preview__store-head { display: flex; gap: 10px; align-items: center; }.detail-preview__store-head img { width: 50px; height: 50px; border-radius: 3px; object-fit: cover; }.detail-preview__store-head p { display: flex; flex-direction: column; gap: 4px; margin: 0; }.detail-preview__store-head p b { font-size: 17px; }.detail-preview__store-head p span { color: #777; font-size: 12px; }.detail-preview__store-head button { margin-left: auto; padding: 8px 18px; border: 0; border-radius: 18px; background: #f44a25; color: #fff; }.detail-preview__scores { display: flex; justify-content: space-between; margin: 16px 0; color: #999; font-size: 12px; }.detail-preview__scores b,.detail-preview__store-products b { color: #f24130; }.detail-preview__store > h3 { padding-top: 12px; border-top: 1px solid #eee; }.detail-preview__store-products { display: grid; grid-template-columns: repeat(3, 1fr); gap: 9px; }.detail-preview__store-products img { width: 100%; aspect-ratio: 1; border-radius: 7px; object-fit: cover; }.detail-preview__store-products p { overflow: hidden; margin: 5px 0 2px; font-size: 12px; white-space: nowrap; text-overflow: ellipsis; }.detail-preview__store-products b { font-size: 13px; font-weight: 500; }
.detail-preview__bottom-nav { display: flex; min-height: 66px; align-items: center; margin-top: 10px; padding: 8px 8px 8px 0; background: #fff; }.detail-preview__bottom-entry { display: flex; min-width: 44px; flex: 1; flex-direction: column; gap: 3px; align-items: center; justify-content: center; padding: 0; border: 0; background: transparent; color: #555; font-size: 10px; line-height: 14px; white-space: nowrap; cursor: pointer; }.detail-preview__bottom-entry svg { width: 21px; height: 21px; font-size: 19px; }.detail-preview__buy-action { height: 42px; padding: 0 16px; border: 0; color: #fff; font-size: 14px; line-height: 42px; white-space: nowrap; cursor: pointer; }.detail-preview__buy-action--cart { border-radius: 24px 0 0 24px; background: #ff9e09; }.detail-preview__buy-action--now { border-radius: 0 24px 24px 0; background: #f43b28; }.detail-preview__icon-action { display: grid; place-items: center; padding: 0; border: 0; background: transparent; color: inherit; cursor: pointer; }.detail-preview__notice { position: sticky; z-index: 12; bottom: 10px; width: max-content; max-width: calc(100% - 48px); margin: -42px auto 12px; padding: 8px 13px; border-radius: 18px; background: rgb(0 0 0 / 74%); color: #fff; font-size: 12px; text-align: center; }.preview-notice-enter-active,.preview-notice-leave-active { transition: opacity .18s ease, transform .18s ease; }.preview-notice-enter-from,.preview-notice-leave-to { opacity: 0; transform: translateY(5px); }
.detail-editor__config { position: fixed; top: 80px; right: 0; bottom: 0; z-index: 20; width: 400px; overflow: hidden; border-left: 1px solid #ebeef3; background: #fff; }.detail-editor__config,.detail-editor__config * { box-sizing: border-box; }.detail-config__title { height: 72px; padding: 22px 30px; border-bottom: 1px solid #e9ebef; color: #2b2f35; font-size: 20px; font-weight: 600; }.detail-config__scroll { height: calc(100% - 72px); overflow-x: clip; overflow-y: auto; background: #f0f2f5; }.detail-config__section { width: 100%; margin-bottom: 10px; padding: 24px 30px; background: #fff; }.detail-config__section h3 { margin: 0 0 18px; color: #343940; font-size: 17px; font-weight: 600; }.detail-config__row { display: flex; gap: 14px; align-items: flex-start; width: 100%; margin-top: 16px; color: #999; }.detail-config__row > span { flex: 0 0 84px; padding-top: 3px; white-space: nowrap; }.detail-config__row :deep(.el-radio-group) { min-width: 0; }.detail-config__row :deep(.el-slider) { flex: 1; margin: 10px 8px 0 0; }.detail-config__checks { display: flex; flex: 1; flex-wrap: wrap; gap: 12px 18px; min-width: 0; }.detail-config__checks :deep(.el-checkbox) { margin-right: 0; }.detail-config__row--list { align-items: flex-start; }.detail-config__draggable-list { display: grid; flex: 1; gap: 8px; min-width: 0; }.detail-config__drag-row { display: flex; align-items: center; gap: 10px; min-height: 48px; padding: 0 12px; border-radius: 4px; background: #fafafa; color: #4c5158; }.detail-config__drag-row > svg { color: #d7dce4; font-size: 19px; cursor: grab; }.detail-config__drag-row :deep(.el-checkbox) { margin-right: 0; }.detail-config__drag-row.is-disabled { color: #b4bbc5; }.detail-config__toggle { display: flex; align-items: center; justify-content: space-between; min-height: 48px; padding: 0 16px; margin-top: 8px; border-radius: 5px; background: #fafafa; color: #50545a; }.detail-config__hint { margin: 18px 0 0 98px; color: #999; font-size: 12px; }.detail-editor__footer { position: fixed; right: 0; bottom: 0; left: 0; z-index: 50; display: flex; height: 72px; align-items: center; justify-content: center; border-top: 1px solid #e7e9ed; background: #fff; box-shadow: 0 -5px 16px rgb(30 45 70 / 4%); }.detail-editor__footer :deep(.el-button) { min-width: 112px; }
.detail-config__title { height: 57px; margin-top: 6px; padding: 20px 15px; border-bottom: 0; color: #303133; font-size: 16px; font-weight: 500; line-height: 16px; }
.detail-config__scroll { height: calc(100% - 63px); }
.detail-config__section { margin-bottom: 3px; padding: 20px 15px; }
.detail-config__section h3 { margin: 0 0 20px; color: #303133; font-size: 14px; font-weight: 400; line-height: 14px; }
.detail-config__row { gap: 0; margin-top: 0; margin-bottom: 20px; color: #999; }
.detail-config__row:last-child { margin-bottom: 0; }
.detail-config__row > span { flex: 0 0 56px; margin-right: 46px; padding-top: 0; font-size: 14px; line-height: 17px; }
.detail-config__checks { display: grid; grid-template-columns: repeat(3, 84px); gap: 14px 0; }
.detail-config__checks :deep(.el-checkbox) { width: 84px; margin: 0; font-size: 12px; line-height: 17px; }
.detail-config__checks--inline { display: flex; flex-wrap: nowrap; gap: 0; }
.detail-config__checks--inline :deep(.el-checkbox) { width: auto; margin-right: 32px; }
.detail-config__row--compact { margin-bottom: 0; }
.detail-config__row--slider { align-items: center; }
.detail-config__row--slider :deep(.el-slider) { flex: 1; margin: 0 12px 0 0; }
.detail-config__row--slider :deep(.el-input-number) { width: 94px; flex: 0 0 94px; }
.detail-config__row :deep(.el-radio) { margin-right: 43px; font-size: 13px; }
.detail-config__drag-row { min-height: 36px; height: 36px; gap: 8px; padding: 0 10px; border-radius: 3px; cursor: grab; font-size: 13px; line-height: 36px; }
.detail-config__drag-row > svg { margin-right: 10px; font-size: 18px; }
.detail-config__draggable-list { gap: 6px; }
.detail-config__drag-row:active { cursor: grabbing; }
.detail-config__drag-row > svg { cursor: inherit; }
.detail-config__drag-row.is-disabled { cursor: not-allowed; }
.detail-config__toggle { min-height: 36px; margin-top: 6px; padding: 0 12px; font-size: 14px; line-height: 20px; }
.detail-config__toggle > span { color: #303133; font-size: 14px; line-height: 20px; }
.detail-config__toggle :deep(.el-checkbox) { margin-right: 0; font-size: 14px; }
.detail-config__section--nav { padding-top: 18px; padding-bottom: 18px; }
.detail-config__section--nav h3 { margin-bottom: 16px; }
.detail-config__section--nav .detail-config__row { margin-bottom: 12px; }
.detail-config__section--nav .detail-config__row:last-child { margin-bottom: 0; }
.detail-config__section--nav .detail-config__checks { row-gap: 10px; }
.detail-config__section--nav .detail-config__checks :deep(.el-checkbox) { height: 18px; line-height: 18px; }
.detail-config__section--image { padding-top: 18px; padding-bottom: 18px; }
.detail-config__section--image h3 { margin-bottom: 16px; }
.detail-config__section--image .detail-config__row { margin-bottom: 12px; }
.detail-config__section--image .detail-config__row:last-child { margin-bottom: 0; }
@media (max-width: 1240px) { .detail-editor { padding-right: 360px; }.detail-editor__config { width: 360px; }.detail-editor__canvas { padding-right: 145px; padding-left: 145px; }.detail-preview__module-tip { left: -106px; min-width: 72px; padding: 8px; font-size: 12px; }.detail-config__section { padding-right: 20px; padding-left: 20px; }.detail-config__row > span { flex-basis: 70px; } }
</style>
