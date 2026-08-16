<script lang="ts" setup>
import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { Icon as IconifyIcon } from '@iconify/vue';
import { ElButton, ElMessage } from 'element-plus';

import featuredBanner from '../../../../../app-uni/static/reference/crmeb-banner-03.png';

import {
  getCategoryDecorationApi,
  saveCategoryDecorationApi,
  type CategoryDecorationLayout,
} from '#/api/core/diy';

type LayoutOption = {
  label: string;
  value: CategoryDecorationLayout;
};

const layouts: LayoutOption[] = [
  { value: 'list', label: '商品列表' },
  { value: 'card', label: '大图卡片' },
  { value: 'grid', label: '分类宫格' },
];

const selectedLayout = ref<CategoryDecorationLayout>('list');
const saving = ref(false);
const expandedLayout = ref<CategoryDecorationLayout | null>(null);
const primaryCategories = ['美妆个护', '休闲零食', '食品茶饮', '家居生活', '服饰鞋包', '数码家电'];
const secondaryCategories = ['乳品推荐', '休闲零食', '香水', '特惠专区', '大闸蟹', '精选礼盒'];

async function load() {
  try {
    const data = await getCategoryDecorationApi();
    if (layouts.some((item) => item.value === data.layout)) {
      selectedLayout.value = data.layout;
    }
  } catch {
    // 分类装修未设置时使用第一种布局；不阻断页面编辑。
  }
}

async function save() {
  saving.value = true;
  try {
    await saveCategoryDecorationApi(selectedLayout.value);
    ElMessage.success('分类装修已保存，C 端分类页将使用此布局');
  } finally {
    saving.value = false;
  }
}

onMounted(load);
</script>

<template>
  <Page auto-content-height class="category-decoration-page">
    <section class="phone-stage">
      <div
        v-for="option in layouts"
        :key="option.value"
        class="phone-card"
        :class="{ 'phone-card--active': selectedLayout === option.value }"
        role="button"
        tabindex="0"
        @click="selectedLayout = option.value"
        @keyup.enter="selectedLayout = option.value"
      >
        <span class="phone-card__title">{{ option.label }}</span>
        <span class="phone-preview" :class="`phone-preview--${option.value}`">
          <span class="phone-preview__status"><b>9:41</b><i>▮▮▮ ◔ ▰</i></span>
          <span class="phone-preview__search">⌕&nbsp;&nbsp;请输入商品名称 <b>•••　◉</b></span>

          <span v-if="option.value === 'list' || option.value === 'card'" class="phone-preview__roots phone-preview__roots--image">
            <span v-for="index in 4" :key="index" :class="{ 'is-active': index === 1 }"><i :class="`product-tone product-tone--${index}`"></i>{{ primaryCategories[index - 1] }}</span>
            <button class="phone-preview__all-trigger" type="button" @click.stop="expandedLayout = expandedLayout === option.value ? null : option.value">全部<br />分类</button>
          </span>
          <span v-else class="phone-preview__roots phone-preview__roots--text">
            <span class="is-active">推荐</span><span v-for="category in primaryCategories" :key="category">{{ category }}</span>
          </span>
          <span v-if="expandedLayout === option.value && (option.value === 'list' || option.value === 'card')" class="phone-preview__primary-panel">
            <button v-for="(category, index) in primaryCategories" :key="category" type="button" @click.stop="expandedLayout = null"><i :class="`product-tone product-tone--${(index % 5) + 1}`"></i><b>{{ category }}</b></button>
          </span>

          <span class="phone-preview__body">
            <span class="phone-preview__aside"><i v-for="(category, index) in option.value === 'grid' ? primaryCategories : secondaryCategories" :key="category" :class="{ 'is-active': index === 0 }">{{ category }}</i></span>
            <span v-if="option.value === 'list'" class="phone-preview__main phone-preview__list">
              <span v-for="index in 5" :key="index" class="phone-preview__list-item"><i :class="`product-tone product-tone--${(index % 5) + 1}`"></i><b>橙中爱马仕黑标新骑士晚...</b><small v-if="index === 1">满减活动　包邮</small><em>¥32.00</em><strong><IconifyIcon icon="lucide:shopping-cart" /></strong></span>
            </span>
            <span v-else-if="option.value === 'card'" class="phone-preview__main phone-preview__card-main">
              <span v-for="index in 2" :key="index" class="phone-preview__large-card">
                <span class="phone-preview__banner"><img :src="featuredBanner" alt="精选好物" /></span>
                <b>Haier/海尔 时尚静音冰箱三开门出租家用...</b>
                <span class="phone-preview__large-card-bottom"><em>¥32.00　<small>VIP ¥26.00</small></em><button aria-label="加入购物车"><IconifyIcon icon="lucide:shopping-cart" /></button></span>
              </span>
            </span>
            <span v-else class="phone-preview__main phone-preview__grid">
              <span v-for="index in 9" :key="index" class="phone-preview__grid-item"><i :class="`product-tone product-tone--${(index % 5) + 1}`"></i><em>{{ ['全部商品', '糕点饼干', '水果蔬菜'][index % 3] }}</em></span>
            </span>
          </span>
          <span v-if="option.value === 'list' || option.value === 'card'" class="phone-preview__cartbar"><b><IconifyIcon icon="lucide:shopping-cart" /><small>2</small></b><strong><small>已选 12 件</small>¥999<sup>.00</sup></strong><span>查看明细⌃</span><em>去结算</em></span>
          <span class="phone-preview__tabbar"><i>⌂<small>首页</small></i><i class="is-active">▦<small>分类</small></i><i>♧<small>购物车</small></i><i>♙<small>我的</small></i></span>
        </span>
      </div>
    </section>

    <footer class="category-decoration-page__footer">
      <ElButton :loading="saving" size="large" type="primary" @click="save">保存</ElButton>
    </footer>
  </Page>
</template>

<style scoped>
.category-decoration-page :deep(.vben-page-content) { height:100%; min-height:0; display:flex; padding:0; }
.phone-stage { min-height:0; height:calc(100% - 64px); flex:0 0 calc(100% - 64px); display:flex; align-items:center; justify-content:center; gap:32px; padding:12px 28px; box-sizing:border-box; overflow-x:auto; border-radius:8px; background:#f3f5f8; }
.phone-card { position:relative; flex:0 0 278px; width:278px; padding:0; border:2px solid transparent; border-radius:18px; color:inherit; background:#fff; cursor:pointer; overflow:hidden; box-shadow:0 8px 22px rgb(15 23 42 / 8%); transition:.18s ease; }
.phone-card--active { border-color:#2f80ed; box-shadow:0 8px 28px rgb(47 128 237 / 20%); }
.phone-card__title { position:absolute; z-index:2; top:10px; left:12px; padding:3px 7px; border-radius:10px; color:#fff; background:rgb(31 41 55 / 72%); font-size:11px; }
.phone-preview { position:relative; height:548px; display:flex; flex-direction:column; padding-bottom:50px; box-sizing:border-box; color:#3f3f46; background:#fff; font-size:9px; text-align:left; overflow:hidden; }
.phone-preview__status { height:40px; display:flex; align-items:center; justify-content:space-between; padding:0 16px; color:#101010; font-size:9px; }
.phone-preview__status b { font-size:15px; }.phone-preview__status i { font-style:normal; letter-spacing:1px; }
.phone-preview__search { height:24px; display:flex; align-items:center; justify-content:space-between; margin:0 12px 9px; padding:0 9px; border-radius:13px; color:#adb1b8; background:#f5f5f5; font-size:8px; }.phone-preview__search b { padding-left:10px; color:#252525; border-left:1px solid #e5e5e5; }
.phone-preview__roots { height:69px; display:flex; flex-shrink:0; gap:9px; padding:3px 9px 0; border-bottom:1px solid #f1f1f1; overflow:hidden; }.phone-preview__roots > span { width:42px; flex:0 0 42px; display:flex; flex-direction:column; align-items:center; gap:3px; color:#555; white-space:nowrap; }.phone-preview__roots--image i { width:31px; height:31px; border-radius:50%; }.phone-preview__roots .is-active { color:#ff5147; font-weight:700; }.phone-preview__roots--image .is-active i { outline:1px solid #ff5147; outline-offset:2px; }.phone-preview__all-trigger { width:27px; height:45px; flex:0 0 27px; margin-left:auto; border:0; border-left:1px solid #f0f0f0; color:#555; background:#fff; font-size:8px; line-height:13px; }.phone-preview__roots--text { height:36px; align-items:center; gap:8px; padding:0 11px; overflow-x:auto; scrollbar-width:none; }.phone-preview__roots--text > span { width:auto; flex:0 0 auto; display:block; padding:5px 7px; border-radius:11px; color:#555; background:#f7f7f7; white-space:nowrap; }.phone-preview__roots--text .is-active { border-left:2px solid #ff5147; border-radius:0; color:#ff5147; background:transparent; }.phone-preview__primary-panel { display:grid; grid-template-columns:repeat(3, 1fr); flex:0 0 48px; gap:4px; padding:5px 9px; border-bottom:1px solid #f0f0f0; background:#fff; }.phone-preview__primary-panel button { overflow:hidden; border:0; border-radius:6px; color:#555; background:#f7f7f7; font-size:7px; text-overflow:ellipsis; white-space:nowrap; }
.phone-preview__body { min-height:0; flex:1; display:flex; }.phone-preview__aside { width:61px; flex:0 0 61px; display:flex; flex-direction:column; background:#f8f8f8; }.phone-preview__aside i { min-height:40px; display:flex; align-items:center; justify-content:center; color:#888; font-style:normal; }.phone-preview__aside .is-active { position:relative; color:#ff5147; background:#fff; font-weight:700; }.phone-preview__aside .is-active::before { position:absolute; left:0; width:2px; height:17px; border-radius:1px; background:#ff5147; content:''; }
.phone-preview__main { min-width:0; flex:1; padding:12px 10px; box-sizing:border-box; overflow:hidden; }.phone-preview__grid { display:grid; grid-template-columns:repeat(3, 1fr); align-content:start; gap:12px 4px; }.phone-preview__grid > b { grid-column:1 / -1; font-size:12px; }.phone-preview__grid-item { display:flex; flex-direction:column; align-items:center; gap:5px; }.phone-preview__grid-item i { width:38px; height:38px; }.phone-preview__grid-item em { color:#555; font-style:normal; text-align:center; }
.phone-preview__list { padding:8px; }.phone-preview__list-item { position:relative; min-height:54px; display:grid; grid-template-columns:54px 1fr; grid-template-rows:1fr 17px; column-gap:7px; padding:5px 0; border-bottom:1px solid #f3f3f3; }.phone-preview__list-item i { grid-row:1 / 3; width:54px; height:54px; }.phone-preview__list-item b { overflow:hidden; padding-top:2px; font-size:10px; text-overflow:ellipsis; white-space:nowrap; }.phone-preview__list-item > small { position:absolute; top:21px; left:61px; padding:1px 3px; border:1px solid #ff8a80; border-radius:3px; color:#ff5147; font-size:6px; }.phone-preview__list-item em { align-self:end; color:#ff5147; font-size:12px; font-style:normal; font-weight:700; }.phone-preview__list-item strong { position:absolute; right:0; bottom:6px; width:22px; height:22px; display:flex; align-items:center; justify-content:center; border-radius:50%; color:#fff; background:#ff5147; }.phone-preview__list-item strong svg { width:13px; height:13px; stroke-width:2.5; }
.phone-preview__card-main { display:flex; flex-direction:column; gap:6px; font-size:10px; }.phone-preview__large-card { display:flex; flex-direction:column; gap:5px; padding:7px; border-radius:8px; background:#fafafa; }.phone-preview__banner { height:76px; display:block; border-radius:7px; background:#eef1f5; overflow:hidden; }.phone-preview__banner img { width:100%; height:100%; display:block; object-fit:cover; }.phone-preview__large-card > b { overflow:hidden; font-size:9px; text-overflow:ellipsis; white-space:nowrap; }.phone-preview__large-card-bottom { display:flex; align-items:center; justify-content:space-between; }.phone-preview__large-card-bottom em { color:#ff5147; font-size:11px; font-style:normal; font-weight:700; }.phone-preview__card-main small { padding:2px; border-radius:4px; color:#fff; background:#555; font-size:7px; }.phone-preview__card-main button { padding:5px 7px; border:0; border-radius:12px; color:#fff; background:#ff5147; font-size:8px; }
.product-tone { display:block; border-radius:8px; background:linear-gradient(145deg, #eceff4, #c7ceda); }.product-tone--1 { background:linear-gradient(145deg, #f6e7db, #c9a68c); }.product-tone--2 { background:linear-gradient(145deg, #f4f4f1, #d4d7d2); }.product-tone--3 { background:linear-gradient(145deg, #f6d3ba, #d8956d); }.product-tone--4 { background:linear-gradient(145deg, #ddd6ce, #af9f8e); }.product-tone--5 { background:linear-gradient(145deg, #dcead4, #88a37b); }
.phone-preview__tabbar { position:absolute; z-index:3; right:0; bottom:0; left:0; height:50px; display:flex; align-items:center; border-top:1px solid #eee; background:#fff; }.phone-preview__tabbar i { flex:1; display:flex; flex-direction:column; align-items:center; color:#555; font-size:17px; font-style:normal; line-height:17px; }.phone-preview__tabbar small { margin-top:3px; font-size:8px; }.phone-preview__tabbar .is-active { color:#ff5147; }
.phone-preview__cartbar { position:absolute; z-index:2; right:10px; bottom:59px; left:10px; height:38px; display:flex; align-items:center; gap:7px; padding:0 8px 0 4px; border:1px solid rgb(255 255 255 / 12%); border-radius:19px; color:#fff; background:#252a31; box-shadow:0 6px 14px rgb(28 34 42 / 24%); }.phone-preview__cartbar > b { width:32px; height:32px; display:flex; align-items:center; justify-content:center; border-radius:50%; color:#fff6ed; background:#ff704d; box-shadow:0 2px 6px rgb(255 112 77 / 36%); font-size:15px; }.phone-preview__cartbar strong { flex:1; color:#fff; font-size:14px; letter-spacing:.2px; }.phone-preview__cartbar strong small { color:#b7c0ca; font-size:7px; font-weight:400; }.phone-preview__cartbar em { padding:7px 9px; border-radius:13px; color:#fff; background:#ff513b; box-shadow:0 2px 5px rgb(255 81 59 / 32%); font-size:8px; font-style:normal; }
.category-decoration-page__footer { position:fixed; z-index:30; right:0; bottom:0; left:0; height:64px; display:flex; align-items:center; justify-content:center; box-sizing:border-box; border-top:1px solid #e5e7eb; background:#fff; box-shadow:0 -8px 18px rgb(15 23 42 / 5%); }.category-decoration-page__footer .el-button { min-width:140px; }
.phone-preview__roots--text { padding-left:0; }.phone-preview__roots--text .is-active { position:sticky; z-index:1; left:0; background:#fff; }.phone-preview__primary-panel { position:absolute; z-index:5; top:142px; right:0; left:0; display:grid; grid-template-columns:repeat(4, 1fr); gap:7px 4px; padding:9px 8px; border-top:1px solid #f0f0f0; border-bottom:1px solid #e8e8e8; background:#fff; box-shadow:0 10px 18px rgb(15 23 42 / 12%); }.phone-preview__primary-panel button { display:flex; flex-direction:column; align-items:center; gap:3px; overflow:hidden; padding:0; border:0; color:#555; background:transparent; font-size:7px; }.phone-preview__primary-panel i { width:27px; height:27px; border-radius:50%; }.phone-preview__primary-panel b { overflow:hidden; max-width:100%; font-size:7px; font-weight:400; text-overflow:ellipsis; white-space:nowrap; }.phone-preview__cartbar { height:40px; gap:6px; padding:0 7px 0 4px; border-color:#3c4654; border-radius:20px; background:#242a33; }.phone-preview__cartbar > b { position:relative; width:33px; height:33px; background:#ff704d; }.phone-preview__cartbar > b small { position:absolute; top:-3px; right:-3px; min-width:11px; height:11px; border:1px solid #fff; border-radius:50%; color:#fff; background:#e94d3d; font-size:7px; line-height:11px; text-align:center; }.phone-preview__cartbar strong { display:flex; flex:1; flex-direction:column; line-height:13px; }.phone-preview__cartbar strong small { color:#aeb8c4; }.phone-preview__cartbar strong sup { color:#c9d0d8; font-size:7px; }.phone-preview__cartbar > span { color:#aeb8c4; font-size:7px; white-space:nowrap; }.phone-preview__cartbar em { background:#ff563e; }
.phone-preview__cartbar > b svg { width:17px; height:17px; stroke-width:2.5; }
.phone-preview__card-main button { width:26px; height:26px; display:flex; align-items:center; justify-content:center; padding:0; border-radius:50%; }
.phone-preview__card-main button svg { width:14px; height:14px; stroke-width:2.5; }
@media (max-width: 900px) { .phone-stage { justify-content:flex-start; gap:16px; padding:12px; }.phone-card { flex-basis:250px; width:250px; } }
</style>
