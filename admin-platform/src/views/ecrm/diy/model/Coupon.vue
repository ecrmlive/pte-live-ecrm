<template>
  <div class="drag optional" :class="{ selected: index === selectedIndex }" @click.stop="diyEditer(index)">
    <div class="diy-coupon" :style="componentStyle">
      <div class="coupon-wrapper-box" :style="wrapperStyle">
        <img v-if="Number(style.bgtype) === 2 && style.bgimage" class="bg-couponbg" :src="style.bgimage" alt="" />
        <div class="coupon-rail" :style="{ gap: `${couponSpacing}px` }">
          <article
            v-for="couponIndex in previewCount"
            :key="`coupon-${couponIndex}`"
            class="coupon-card"
            :class="`coupon-style-${couponStyleType}`"
            :style="couponCardStyle"
          >
            <template v-if="couponStyleType === 2">
              <span class="coupon-kind">通用券</span>
              <strong>¥70</strong>
              <small>满500元可用</small>
              <button type="button">{{ buttonText }}</button>
            </template>
            <template v-else-if="couponStyleType === 3">
              <div class="coupon-main"><strong>¥70</strong><small>满500元可用</small></div>
              <button type="button">{{ buttonText }}</button>
            </template>
            <template v-else-if="couponStyleType === 4">
              <div class="coupon-main"><span>通用券</span><strong>¥70</strong><small>满500元可用</small></div>
              <button type="button">{{ buttonText }}</button>
            </template>
            <template v-else-if="couponStyleType === 5">
              <div class="coupon-main"><strong>¥70</strong><small>满500元可用</small></div>
              <button type="button">领<br />取</button>
            </template>
            <template v-else>
              <div class="coupon-main"><strong>¥70</strong><small>满500元可用</small><em>通用券</em></div>
              <button type="button">{{ buttonText }}</button>
            </template>
          </article>
        </div>
      </div>
    </div>
    <div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
  </div>
</template>

<script>
import { resolveCouponPreviewCount } from '../params/shared/marketing-helpers';

export default {
  inject: ['diyModel'],
  props: ['item', 'index', 'selectedIndex'],
  computed: {
    style() { return this.item?.style || {}; },
    params() { return this.item?.params || {}; },
    previewCount() { return resolveCouponPreviewCount(this.item || {}); },
    couponStyleType() {
      const value = Number(this.style.type || 1);
      return value >= 1 && value <= 5 ? value : 1;
    },
    couponSpacing() {
      const value = Number(this.style.couponSpacing ?? 6);
      return Number.isFinite(value) ? Math.min(24, Math.max(0, value)) : 6;
    },
    useCustomColors() { return this.style.colorMode === 'custom'; },
    accentColor() { return this.useCustomColors ? this.style.btncolor || '#ff4c01' : '#ff4c01'; },
    priceColor() { return this.useCustomColors ? this.style.pricecolor || this.accentColor : this.accentColor; },
    buttonText() { return this.params.btntext || '立即领取'; },
    componentStyle() {
      return {
        background: this.style.bgcolor || '#f5f5f5',
        paddingLeft: `${Number(this.style.paddingLeft ?? 10)}px`,
        paddingRight: `${Number(this.style.paddingLeft ?? 10)}px`,
        paddingTop: `${Number(this.style.paddingTop ?? 0)}px`,
        paddingBottom: `${Number(this.style.paddingBottom ?? 0)}px`,
      };
    },
    wrapperStyle() {
      return {
        background: Number(this.style.bgtype) === 1 ? this.style.background || '#ffffff' : 'transparent',
        borderRadius: `${Number(this.style.topRadio ?? 8)}px ${Number(this.style.topRadio ?? 8)}px ${Number(this.style.bottomRadio ?? 8)}px ${Number(this.style.bottomRadio ?? 8)}px`,
      };
    },
    couponCardStyle() {
      return {
        '--coupon-accent': this.accentColor,
        '--coupon-price': this.priceColor,
        '--coupon-text': this.style.descolor || '#666666',
        '--coupon-subtext': this.style.cillcolor || '#999999',
        '--coupon-button-text': this.style.btnTxtcolor || '#ffffff',
        borderRadius: `${Number(this.style.btnRadio ?? 12)}px`,
        boxShadow: this.style.shadow === 'on' ? '0 4px 12px rgb(0 0 0 / 12%)' : 'none',
      };
    },
  },
  methods: {
    diyEditer(index) { this.diyModel?.onEditer(index); },
    diyDeleteItem(index) { this.diyModel?.onDeleleItem(index); },
  },
};
</script>

<style lang="scss" scoped>
.diy-coupon { position: relative; width: 100%; box-sizing: border-box; }
.coupon-wrapper-box { position: relative; box-sizing: border-box; min-height: 126px; overflow: hidden; padding: 12px; }
.bg-couponbg { position: absolute; inset: 0; width: 100%; height: 100%; object-fit: cover; }
.coupon-rail { position: relative; z-index: 1; display: flex; overflow-x: auto; padding-bottom: 2px; scrollbar-width: none; }.coupon-rail::-webkit-scrollbar { display: none; }
.coupon-card { position: relative; display: flex; flex: 0 0 104px; align-items: stretch; min-height: 104px; overflow: hidden; color: var(--coupon-text); background: #fff; }
.coupon-main { display: flex; flex: 1; flex-direction: column; align-items: center; justify-content: center; min-width: 0; padding: 8px 4px; text-align: center; }
.coupon-card strong { color: var(--coupon-price); font-size: 23px; line-height: 1; }.coupon-card small { margin-top: 5px; color: var(--coupon-subtext); font-size: 9px; }.coupon-card em,.coupon-kind { margin-top: 5px; color: var(--coupon-text); font-size: 9px; font-style: normal; }
.coupon-card button { min-width: 26px; padding: 4px; color: var(--coupon-button-text); font-size: 10px; line-height: 1.1; background: var(--coupon-accent); border: 0; cursor: pointer; }
.coupon-style-1 { flex-direction: column; border: 1px solid var(--coupon-accent); }.coupon-style-1 .coupon-main { border-bottom: 1px dashed var(--coupon-accent); }.coupon-style-1 button { min-height: 28px; border-radius: 0; }
.coupon-style-2 { flex-direction: column; padding: 8px; color: #fff; background: linear-gradient(135deg, var(--coupon-accent), #ff8a22); }.coupon-style-2 .coupon-kind,.coupon-style-2 strong,.coupon-style-2 small { color: #fff; }.coupon-style-2 strong { margin-top: 3px; }.coupon-style-2 button { width: 68px; min-height: 21px; margin: 6px auto 0; background: rgb(255 255 255 / 22%); border: 1px solid rgb(255 255 255 / 45%); border-radius: 12px; }
.coupon-style-3 { min-height: 84px; border: 1px solid var(--coupon-accent); }.coupon-style-3 .coupon-main { align-items: flex-start; padding-left: 10px; }.coupon-style-3 button { border-left: 1px dashed rgb(255 255 255 / 70%); }
.coupon-style-4 { min-height: 90px; background: linear-gradient(135deg, var(--coupon-accent), #ff7940); }.coupon-style-4 .coupon-main { color: #fff; }.coupon-style-4 strong,.coupon-style-4 small,.coupon-style-4 .coupon-main span { color: #fff; }.coupon-style-4 button { margin: 31px 7px 31px 0; border-radius: 14px; }
.coupon-style-5 { min-height: 70px; border: 1px solid var(--coupon-accent); }.coupon-style-5 .coupon-main { align-items: flex-start; padding-left: 10px; }.coupon-style-5 button { display: flex; align-items: center; justify-content: center; min-width: 22px; padding: 2px; }
</style>
