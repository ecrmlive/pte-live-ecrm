<script setup lang="ts">
import { computed, ref, watch } from 'vue';

import { Icon as IconifyIcon } from '@iconify/vue';

import { resolveCosMediaUrl } from '#/utils/live/cosMediaUrl.js';

export type ProductDetailDiyPreviewData = {
  image?: string;
  mer_type_name?: string;
  ot_price?: number;
  price?: number;
  sales?: number;
  slider_image?: string[];
  stock?: number;
  title?: string;
};

const props = defineProps<{
  product?: ProductDetailDiyPreviewData | null;
}>();

const activeIndex = ref(0);

const images = computed(() => {
  const list: string[] = [];
  const cover = resolveCosMediaUrl(String(props.product?.image || '').trim());
  if (cover) list.push(cover);
  for (const raw of props.product?.slider_image || []) {
    const url = resolveCosMediaUrl(String(raw || '').trim());
    if (url && !list.includes(url)) list.push(url);
  }
  return list;
});

const thumbImages = computed(() => images.value.slice(0, 6));

watch(
  () => props.product,
  () => {
    activeIndex.value = 0;
  },
);

const priceText = computed(() =>
  Number(props.product?.price || 0).toFixed(2),
);

const otPriceText = computed(() =>
  Number(props.product?.ot_price || 0).toFixed(2),
);

const showOtPrice = computed(() => Number(props.product?.ot_price || 0) > 0);

const stockText = computed(() => Number(props.product?.stock || 0));

const salesText = computed(() => Number(props.product?.sales || 0));

const storeBadge = computed(() => {
  const name = String(props.product?.mer_type_name || '').trim();
  return name || '';
});

const titleText = computed(
  () => String(props.product?.title || '').trim() || '商品标题',
);

function selectImage(index: number) {
  activeIndex.value = index;
}

function onCarouselClick() {
  if (images.value.length <= 1) return;
  activeIndex.value = (activeIndex.value + 1) % images.value.length;
}
</script>

<template>
  <div class="goods-diy">
    <div class="goods-diy__gallery" @click="onCarouselClick">
      <img
        v-if="images[activeIndex]"
        class="goods-diy__hero"
        :src="images[activeIndex]"
        alt=""
      />
      <div v-else class="goods-diy__hero goods-diy__hero--empty">暂无图片</div>
      <div v-if="images.length > 1" class="goods-diy__indicators">
        <span
          v-for="(_, idx) in images"
          :key="idx"
          class="goods-diy__indicator"
          :class="{ 'is-active': idx === activeIndex }"
        />
      </div>
    </div>

    <div class="goods-diy__price-row">
      <div class="goods-diy__thumbs">
        <button
          v-for="(img, idx) in thumbImages"
          :key="`${img}-${idx}`"
          type="button"
          class="goods-diy__thumb"
          :class="{ 'is-active': idx === activeIndex }"
          @click.stop="selectImage(idx)"
        >
          <img :src="img" alt="" />
        </button>
        <IconifyIcon
          v-if="images.length > 2"
          class="goods-diy__thumb-more"
          icon="ant-design:right-outlined"
        />
      </div>
      <div class="goods-diy__price">
        <span class="goods-diy__yen">¥</span>
        <span class="goods-diy__amount">{{ priceText }}</span>
      </div>
      <div class="goods-diy__actions">
        <div class="goods-diy__action">
          <IconifyIcon icon="ant-design:star-outlined" />
          <span>收藏</span>
        </div>
        <div class="goods-diy__action">
          <IconifyIcon icon="ant-design:share-alt-outlined" />
          <span>分享</span>
        </div>
      </div>
    </div>

    <div class="goods-diy__svip">
      <div class="goods-diy__svip-left">
        <span class="goods-diy__svip-mark">◆</span>
        <span>开通 SVIP会员 省钱多多，权益多多</span>
      </div>
      <span class="goods-diy__svip-cta">立即开通 ›</span>
    </div>

    <div class="goods-diy__info">
      <div class="goods-diy__title-line">
        <span v-if="storeBadge" class="goods-diy__badge">{{ storeBadge }}</span>
        <span class="goods-diy__title">{{ titleText }}</span>
      </div>
      <div class="goods-diy__tags">
        <span class="goods-diy__coupon">领券</span>
      </div>
      <div class="goods-diy__stats">
        <span v-if="showOtPrice" class="goods-diy__ot">¥{{ otPriceText }}</span>
        <span v-else class="goods-diy__ot goods-diy__ot--placeholder" />
        <span>库存: {{ stockText }}件</span>
        <span>已售: {{ salesText }}件</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.goods-diy {
  display: flex;
  flex-direction: column;
  min-height: 100%;
  background: #f5f5f5;
  color: #222;
  font-family:
    -apple-system,
    BlinkMacSystemFont,
    'PingFang SC',
    'Helvetica Neue',
    sans-serif;
  user-select: none;
}

.goods-diy__gallery {
  position: relative;
  width: 100%;
  aspect-ratio: 1 / 1;
  overflow: hidden;
  background: #1a2744;
  cursor: pointer;
}

.goods-diy__hero {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.goods-diy__hero--empty {
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.65);
  font-size: 14px;
}

.goods-diy__indicators {
  position: absolute;
  right: 12px;
  bottom: 10px;
  left: 12px;
  display: flex;
  gap: 4px;
  justify-content: center;
}

.goods-diy__indicator {
  flex: 1;
  max-width: 36px;
  height: 2px;
  border-radius: 1px;
  background: rgba(255, 255, 255, 0.35);
}

.goods-diy__indicator.is-active {
  background: rgba(255, 255, 255, 0.95);
}

.goods-diy__price-row {
  display: flex;
  gap: 8px;
  align-items: center;
  padding: 10px 12px 8px;
  background: #fff;
}

.goods-diy__thumbs {
  display: flex;
  flex-shrink: 0;
  gap: 6px;
  align-items: center;
}

.goods-diy__thumb {
  width: 36px;
  height: 36px;
  padding: 0;
  overflow: hidden;
  border: 1px solid #eee;
  border-radius: 4px;
  background: #fafafa;
  cursor: pointer;
}

.goods-diy__thumb.is-active {
  border-color: #e93323;
}

.goods-diy__thumb img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.goods-diy__thumb-more {
  color: #bbb;
  font-size: 12px;
}

.goods-diy__price {
  display: flex;
  flex: 1;
  gap: 2px;
  align-items: baseline;
  justify-content: center;
  min-width: 0;
  color: #e93323;
  font-weight: 700;
  line-height: 1;
}

.goods-diy__yen {
  font-size: 14px;
}

.goods-diy__amount {
  font-size: 26px;
  letter-spacing: -0.5px;
}

.goods-diy__actions {
  display: flex;
  flex-shrink: 0;
  gap: 10px;
  align-items: flex-start;
}

.goods-diy__action {
  display: flex;
  flex-direction: column;
  gap: 2px;
  align-items: center;
  color: #666;
  font-size: 11px;
  line-height: 1.2;
}

.goods-diy__action :deep(svg) {
  font-size: 18px;
}

.goods-diy__svip {
  display: flex;
  gap: 8px;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: linear-gradient(90deg, #f8ecd0 0%, #f3e2b8 100%);
  color: #6b4e1e;
  font-size: 11px;
  line-height: 1.3;
}

.goods-diy__svip-left {
  display: flex;
  gap: 4px;
  align-items: center;
  min-width: 0;
}

.goods-diy__svip-mark {
  color: #1a1a1a;
  font-size: 10px;
}

.goods-diy__svip-cta {
  flex-shrink: 0;
  color: #8a6420;
  font-weight: 600;
}

.goods-diy__info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  background: #fff;
}

.goods-diy__title-line {
  display: -webkit-box;
  overflow: hidden;
  color: #222;
  font-size: 14px;
  font-weight: 500;
  line-height: 1.45;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
}

.goods-diy__badge {
  display: inline-block;
  margin-right: 6px;
  padding: 1px 4px;
  border-radius: 2px;
  background: #e93323;
  color: #fff;
  font-size: 10px;
  font-weight: 600;
  line-height: 16px;
  vertical-align: 1px;
}

.goods-diy__title {
  word-break: break-word;
}

.goods-diy__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.goods-diy__coupon {
  display: inline-flex;
  align-items: center;
  padding: 0 6px;
  border: 1px solid #e93323;
  border-radius: 2px;
  color: #e93323;
  font-size: 11px;
  line-height: 18px;
}

.goods-diy__stats {
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: #999;
  font-size: 12px;
  line-height: 1.4;
}

.goods-diy__ot {
  text-decoration: line-through;
}

.goods-diy__ot--placeholder {
  min-width: 56px;
  text-decoration: none;
}
</style>
