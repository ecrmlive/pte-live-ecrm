<template>
  <view class="diy-root">
    <view v-for="(item, idx) in items" :key="itemKey(item, idx)" class="diy-block">
      <!-- 轮播 -->
      <swiper
        v-if="item.type === 'banner' || item.type === 'topMerge'"
        class="diy-banner"
        circular
        autoplay
        interval="4000"
        indicator-dots
      >
        <swiper-item
          v-for="(row, i) in bannerRows(item)"
          :key="i"
          @click="go(row.linkUrl || row.link_url)"
        >
          <image
            v-if="row.imgUrl || row.img_url"
            class="diy-banner-img"
            :src="String(row.imgUrl || row.img_url)"
            mode="aspectFill"
          />
          <view v-else class="diy-banner-fallback">
            <text>{{ row.imgName || row.title || "轮播图" }}</text>
          </view>
        </swiper-item>
      </swiper>

      <!-- 搜索 -->
      <view
        v-else-if="item.type === 'search'"
        class="diy-search"
        @click="go('/pages/search/index')"
      >
        <text class="diy-search-text">{{
          item.params?.searchText || item.params?.placeholder || "搜索商品"
        }}</text>
      </view>

      <!-- 导航组 -->
      <view v-else-if="item.type === 'navBar' || item.type === 'option'" class="diy-nav">
        <view
          v-for="(row, i) in dataRows(item)"
          :key="i"
          class="diy-nav-item"
          @click="go(row.linkUrl || row.link_url)"
        >
          <image
            v-if="row.imgUrl || row.img_url"
            class="diy-nav-icon"
            :src="String(row.imgUrl || row.img_url)"
            mode="aspectFill"
          />
          <text class="diy-nav-text">{{ row.text || row.title || row.name }}</text>
        </view>
      </view>

      <!-- 标题 -->
      <view v-else-if="item.type === 'title'" class="diy-title">
        <text class="diy-title-text">{{ item.params?.title || item.name || "标题" }}</text>
      </view>

      <!-- 公告 -->
      <view v-else-if="item.type === 'notice'" class="diy-notice">
        <text class="diy-notice-label">公告</text>
        <text class="diy-notice-text">{{
          item.params?.text || item.params?.title || "暂无公告"
        }}</text>
      </view>

      <!-- 空白 / 辅助线 -->
      <view
        v-else-if="item.type === 'blank' || item.type === 'guide'"
        class="diy-blank"
        :style="{ height: blankH(item) + 'px' }"
      />

      <!-- 富文本 -->
      <rich-text
        v-else-if="item.type === 'richText'"
        class="diy-rich"
        :nodes="String(item.params?.content || item.params?.text || '')"
      />

      <!-- 单图 -->
      <image
        v-else-if="item.type === 'imageSingle'"
        class="diy-single"
        :src="String(dataRows(item)[0]?.imgUrl || item.params?.imgUrl || '')"
        mode="widthFix"
        @click="go(dataRows(item)[0]?.linkUrl)"
      />

      <!-- 商品列表（占位卡片，点进商品列表） -->
      <view
        v-else-if="item.type === 'product' || item.type === 'previewProduct'"
        class="diy-section"
        @click="go('/pages/product/list')"
      >
        <text class="diy-section-title">{{ item.name || "商品列表" }}</text>
        <text class="diy-section-more">查看全部</text>
      </view>

      <!-- 优惠券 / 秒杀 / 拼团 -->
      <view
        v-else-if="
          item.type === 'coupon' ||
          item.type === 'seckillProduct' ||
          item.type === 'assembleProduct' ||
          item.type === 'bargainProduct'
        "
        class="diy-section"
        @click="goMarketing(item.type)"
      >
        <text class="diy-section-title">{{ item.name || item.type }}</text>
        <text class="diy-section-more">进入</text>
      </view>

      <!-- 热区 -->
      <view v-else-if="item.type === 'hotspot'" class="diy-hotspot">
        <image
          v-if="item.params?.imgUrl"
          class="diy-single"
          :src="String(item.params.imgUrl)"
          mode="widthFix"
        />
      </view>

      <!-- 其他未识别组件：显示名称 -->
      <view v-else class="diy-unknown">
        <text>{{ item.name || item.type }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import type { DiyItem } from "@/api/diy";

const props = defineProps<{ items: DiyItem[] }>();

function itemKey(item: DiyItem, idx: number) {
  return String(item._diyUid || `${item.type}-${idx}`);
}

function dataRows(item: DiyItem): Array<Record<string, any>> {
  const d = item.data || item.images || [];
  return Array.isArray(d) ? (d as Array<Record<string, any>>) : [];
}

function bannerRows(item: DiyItem) {
  const rows = dataRows(item);
  return rows.length ? rows : [{ imgName: item.name || "轮播" }];
}

function blankH(item: DiyItem) {
  const h = item.style?.height ?? item.params?.height ?? 20;
  return Number(h) || 20;
}

function go(url: unknown) {
  const path = String(url || "").trim();
  if (!path) return;
  if (path.startsWith("/pages")) {
    uni.navigateTo({ url: path }).catch(() => {
      uni.switchTab({ url: path }).catch(() => undefined);
    });
    return;
  }
  if (path.startsWith("http")) {
    // #ifdef H5
    window.open(path, "_blank");
    // #endif
  }
}

function goMarketing(type: string) {
  if (type === "coupon") go("/pages/coupon/center");
  else if (type === "seckillProduct") go("/pages/seckill/list");
  else if (type === "assembleProduct") go("/pages/combination/list");
  else go("/pages/product/list");
}
</script>

<style scoped>
.diy-banner {
  height: 320rpx;
  margin-bottom: 16rpx;
}
.diy-banner-img,
.diy-banner-fallback {
  width: 100%;
  height: 320rpx;
}
.diy-banner-fallback {
  background: linear-gradient(135deg, #1a6dff, #5b8cff);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 32rpx;
}
.diy-search {
  margin: 16rpx 24rpx;
  padding: 18rpx 24rpx;
  background: #f5f6f8;
  border-radius: 999rpx;
}
.diy-search-text {
  color: #999;
  font-size: 26rpx;
}
.diy-nav {
  display: flex;
  flex-wrap: wrap;
  padding: 8rpx 12rpx 20rpx;
  background: #fff;
}
.diy-nav-item {
  width: 25%;
  padding: 16rpx 8rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}
.diy-nav-icon {
  width: 72rpx;
  height: 72rpx;
  border-radius: 16rpx;
  margin-bottom: 8rpx;
  background: #f0f2f5;
}
.diy-nav-text {
  font-size: 24rpx;
  color: #333;
}
.diy-title {
  padding: 20rpx 24rpx 8rpx;
}
.diy-title-text {
  font-size: 32rpx;
  font-weight: 600;
}
.diy-notice {
  margin: 12rpx 24rpx;
  padding: 16rpx 20rpx;
  background: #fff7e6;
  border-radius: 12rpx;
  display: flex;
  gap: 12rpx;
}
.diy-notice-label {
  color: #d48806;
  font-size: 24rpx;
}
.diy-notice-text {
  flex: 1;
  font-size: 24rpx;
  color: #666;
}
.diy-blank {
  width: 100%;
}
.diy-rich,
.diy-single {
  width: 100%;
  display: block;
}
.diy-section {
  margin: 16rpx 24rpx;
  padding: 24rpx;
  background: #fff;
  border-radius: 16rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.diy-section-title {
  font-size: 30rpx;
  font-weight: 600;
}
.diy-section-more {
  font-size: 24rpx;
  color: #1a6dff;
}
.diy-unknown {
  margin: 8rpx 24rpx;
  padding: 16rpx;
  background: #fafafa;
  color: #999;
  font-size: 24rpx;
  border-radius: 8rpx;
}
.diy-hotspot {
  width: 100%;
}
</style>
