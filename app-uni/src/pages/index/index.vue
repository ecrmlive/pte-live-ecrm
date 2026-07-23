<template>
  <view class="page">
    <view class="nav" :style="{ paddingTop: statusBar + 'px' }">
      <view class="nav-inner">
        <text class="brand">{{ homeTitle || "栖息商城" }}</text>
        <view class="search" @click="goSearch">
          <text class="search-text">搜索商品 / 店铺</text>
        </view>
      </view>
    </view>

    <scroll-view scroll-y class="body" :style="{ paddingTop: navPad + 'px' }">
      <DiyRenderer v-if="diyItems.length" :items="diyItems" />

      <template v-else>
        <view class="banner">
          <swiper circular autoplay interval="4000" indicator-dots class="banner-swiper">
            <swiper-item v-for="b in banners" :key="b.id" @click="goUrl(b.url)">
              <view class="banner-item">
                <text class="banner-title">{{ b.title }}</text>
              </view>
            </swiper-item>
          </swiper>
        </view>

        <view v-if="menus.length" class="menus">
          <view v-for="m in menus" :key="m.id" class="menu-item" @click="goUrl(m.url)">
            <text class="menu-name">{{ m.name }}</text>
          </view>
        </view>
      </template>

      <view class="section">
        <view class="section-head">
          <text class="section-title">限时秒杀</text>
          <text class="section-more" @click="goSeckill">全部</text>
        </view>
        <view v-if="!seckill.length" class="empty">暂无秒杀</view>
        <scroll-view v-else scroll-x class="seckill-row">
          <view
            v-for="a in seckill"
            :key="a.seckill_active_id"
            class="seckill-item"
            @click="goDetail(a.product_id)"
          >
            <text class="s-name">{{ a.store_name || a.name }}</text>
            <text class="s-price">¥{{ a.seckill_price }}</text>
            <text class="s-badge">{{ a.in_window ? "抢购中" : "未开场" }}</text>
          </view>
        </scroll-view>
      </view>

      <view class="section">
        <view class="section-head">
          <text class="section-title">拼团好物</text>
          <text class="section-more" @click="goCombination">全部</text>
        </view>
        <view v-if="!groups.length" class="empty">暂无拼团</view>
        <scroll-view v-else scroll-x class="seckill-row">
          <view
            v-for="g in groups"
            :key="g.product_group_id"
            class="seckill-item"
            @click="goCombinationDetail(g.product_group_id)"
          >
            <text class="s-name">{{ g.store_name || "拼团商品" }}</text>
            <text class="s-price">¥{{ g.price }}</text>
            <text class="s-badge">{{ g.buying_count_num }}人团</text>
          </view>
        </scroll-view>
      </view>

      <view class="section">
        <view class="section-head">
          <text class="section-title">预约 / 预售 / 直播 / 社区 / 助力</text>
        </view>
        <view class="menus" style="margin-top: 0">
          <view class="menu-item" @click="goReservation">
            <text class="menu-name">到店预约</text>
          </view>
          <view class="menu-item" @click="goPresell">
            <text class="menu-name">全款预售</text>
          </view>
          <view class="menu-item" @click="goLive">
            <text class="menu-name">直播带货</text>
          </view>
          <view class="menu-item" @click="goAssist">
            <text class="menu-name">邀请助力</text>
          </view>
          <view class="menu-item" @click="goCommunity">
            <text class="menu-name">社区种草</text>
          </view>
        </view>
      </view>

      <view v-if="notices.length" class="section">
        <view class="section-head">
          <text class="section-title">平台公告</text>
          <text class="section-more" @click="goNotices">全部</text>
        </view>
        <view
          v-for="n in notices"
          :key="n.notice_id"
          class="notice-row"
          @click="goNotice(n.notice_id)"
        >
          <text class="notice-title">{{ n.title }}</text>
        </view>
      </view>

      <view class="section">
        <view class="section-head">
          <text class="section-title">热销好物</text>
          <text class="section-more" @click="goList()">全部</text>
        </view>
        <view v-if="!hot.length" class="empty">暂无可售商品</view>
        <view class="goods-grid">
          <view v-for="p in hot" :key="p.id" class="goods-card" @click="goDetail(p.id)">
            <view class="goods-cover">
              <text class="cover-text">{{ p.store_name.slice(0, 1) }}</text>
            </view>
            <view class="goods-meta">
              <text class="goods-name">{{ p.store_name }}</text>
              <text class="goods-mer">{{ p.mer_name }}</text>
              <view class="price-row">
                <text class="price">¥{{ p.price }}</text>
                <text v-if="p.ot_price" class="ot">¥{{ p.ot_price }}</text>
              </view>
            </view>
          </view>
        </view>
      </view>

      <view class="foot-tip">DIY 首页 · 秒杀价在结算时生效</view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { fetchHome, type ProductItem } from "@/api/catalog";
import { fetchNotices, type Notice } from "@/api/content";
import { fetchDiyHome, type DiyBanner, type DiyItem, type DiyMenu } from "@/api/diy";
import DiyRenderer from "@/components/diy/DiyRenderer.vue";
import { fetchSeckillList, type SeckillActive } from "@/api/seckill";
import { fetchGroups, type ProductGroup } from "@/api/combination";

const statusBar = ref(20);
const navPad = computed(() => statusBar.value + 44);
const homeTitle = ref("");
const banners = ref<DiyBanner[]>([{ id: 1, title: "栖息商城" }]);
const menus = ref<DiyMenu[]>([]);
const diyItems = ref<DiyItem[]>([]);
const hot = ref<ProductItem[]>([]);
const seckill = ref<SeckillActive[]>([]);
const groups = ref<ProductGroup[]>([]);
const notices = ref<Notice[]>([]);

onShow(async () => {
  const sys = uni.getSystemInfoSync();
  statusBar.value = sys.statusBarHeight || 20;
  try {
    const [diy, home, sk, combo, noticeRes] = await Promise.all([
      fetchDiyHome().catch(() => null),
      fetchHome(),
      fetchSeckillList(1, 10).catch(() => ({ list: [] as SeckillActive[] })),
      fetchGroups(1, 10).catch(() => ({ list: [] as ProductGroup[] })),
      fetchNotices(1, 3).catch(() => ({ list: [] as Notice[] })),
    ]);
    if (diy) {
      homeTitle.value = diy.title || diy.name || "";
      diyItems.value = Array.isArray(diy.items) ? diy.items : [];
      if (diy.banners?.length) banners.value = diy.banners;
      menus.value = diy.menus || [];
    } else if (home.banners?.length) {
      banners.value = home.banners.map((b, i) => ({
        id: b.id || i + 1,
        title: b.title,
        image: b.image,
      }));
    }
    hot.value = home.hot || [];
    seckill.value = sk.list || [];
    groups.value = combo.list || [];
    notices.value = noticeRes.list || [];
  } catch (e) {
    uni.showToast({ title: (e as Error).message || "首页加载失败", icon: "none" });
  }
});

function goSearch() {
  uni.navigateTo({ url: "/pages/goods/list" });
}
function goList() {
  uni.navigateTo({ url: "/pages/goods/list" });
}
function goSeckill() {
  uni.navigateTo({ url: "/pages/seckill/list" });
}
function goCombination() {
  uni.navigateTo({ url: "/pages/combination/list" });
}
function goCombinationDetail(id: number) {
  uni.navigateTo({ url: `/pages/combination/detail?id=${id}` });
}
function goReservation() {
  uni.navigateTo({ url: "/pages/reservation/list" });
}
function goPresell() {
  uni.navigateTo({ url: "/pages/presell/list" });
}

function goLive() {
  uni.navigateTo({ url: "/pages/live/list" });
}
function goCommunity() {
  uni.navigateTo({ url: "/pages/community/list" });
}
function goAssist() {
  uni.navigateTo({ url: "/pages/assist/list" });
}
function goNotices() {
  uni.navigateTo({ url: "/pages/notice/list" });
}
function goNotice(id: number) {
  uni.navigateTo({ url: `/pages/notice/detail?id=${id}` });
}
function goDetail(id: number) {
  uni.navigateTo({ url: `/pages/goods/detail?id=${id}` });
}
function goUrl(url?: string) {
  if (!url) return;
  uni.navigateTo({ url });
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: var(--qx-bg);
}
.nav {
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  z-index: 20;
  background: linear-gradient(180deg, #fff 0%, #fff8f8 100%);
}
.nav-inner {
  height: 44px;
  padding: 0 24rpx;
  display: flex;
  align-items: center;
  gap: 16rpx;
}
.brand {
  font-size: 32rpx;
  font-weight: 700;
  color: #222;
  flex-shrink: 0;
}
.search {
  flex: 1;
  height: 64rpx;
  border-radius: 32rpx;
  background: #f3f3f3;
  display: flex;
  align-items: center;
  padding: 0 24rpx;
}
.search-text {
  color: #999;
  font-size: 24rpx;
}
.body {
  height: 100vh;
  box-sizing: border-box;
}
.banner {
  margin: 24rpx;
}
.banner-swiper {
  height: 220rpx;
  border-radius: 16rpx;
  overflow: hidden;
}
.banner-item {
  height: 220rpx;
  background: linear-gradient(135deg, #ff6b4a, #e23030);
  display: flex;
  align-items: flex-end;
  padding: 28rpx;
}
.banner-title {
  color: #fff;
  font-size: 36rpx;
  font-weight: 700;
}
.menus {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
  padding: 0 24rpx 8rpx;
}
.menu-item {
  width: calc(25% - 12rpx);
  background: #fff;
  border-radius: 12rpx;
  padding: 20rpx 0;
  text-align: center;
}
.menu-name {
  font-size: 24rpx;
  color: #333;
}
.section {
  margin: 8rpx 24rpx 24rpx;
}
.section-head {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16rpx;
}
.section-title {
  font-size: 32rpx;
  font-weight: 700;
}
.section-more {
  color: #999;
  font-size: 24rpx;
}
.notice-row {
  background: #fff;
  border-radius: 12rpx;
  padding: 20rpx 24rpx;
  margin-bottom: 12rpx;
}
.notice-title {
  font-size: 26rpx;
  color: #333;
}
.empty {
  color: #999;
  font-size: 24rpx;
  padding: 24rpx 0;
}
.seckill-row {
  white-space: nowrap;
}
.seckill-item {
  display: inline-flex;
  flex-direction: column;
  width: 200rpx;
  margin-right: 16rpx;
  padding: 20rpx;
  background: #fff;
  border-radius: 12rpx;
  vertical-align: top;
}
.s-name {
  font-size: 24rpx;
  overflow: hidden;
  text-overflow: ellipsis;
}
.s-price {
  margin-top: 8rpx;
  color: #e23030;
  font-weight: 700;
}
.s-badge {
  margin-top: 6rpx;
  font-size: 20rpx;
  color: #999;
}
.goods-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}
.goods-card {
  width: calc(50% - 8rpx);
  background: #fff;
  border-radius: 12rpx;
  overflow: hidden;
}
.goods-cover {
  height: 220rpx;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
}
.cover-text {
  font-size: 48rpx;
  color: #ccc;
}
.goods-meta {
  padding: 16rpx;
}
.goods-name {
  font-size: 26rpx;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.goods-mer {
  display: block;
  margin-top: 6rpx;
  font-size: 22rpx;
  color: #999;
}
.price-row {
  margin-top: 8rpx;
  display: flex;
  gap: 8rpx;
  align-items: baseline;
}
.price {
  color: #e23030;
  font-weight: 700;
}
.ot {
  color: #bbb;
  text-decoration: line-through;
  font-size: 22rpx;
}
.foot-tip {
  text-align: center;
  color: #bbb;
  font-size: 22rpx;
  padding: 32rpx 0 80rpx;
}
</style>
