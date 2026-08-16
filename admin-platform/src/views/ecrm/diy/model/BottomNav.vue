<template>
  <div class="bottom-nav-wrap" :class="{ 'bottom-nav-wrap--fixed': isFixed }" :style="wrapStyle" @click.stop="diyEditer(index)">
    <div class="drag optional" :class="{ selected: index === selectedIndex }">
      <nav class="bottom-nav" :class="[`bottom-nav--${navigationType}`, { 'bottom-nav--floating': isFloating }]" :style="navStyle">
        <button
          v-for="(entry, entryIndex) in visibleItems"
          :key="entryIndex"
          type="button"
          class="bottom-nav__item"
          :class="{ 'is-active': entryIndex === activeIndex }"
        >
          <span v-if="navigationType !== 'text'" class="bottom-nav__icon">
            <img v-if="iconUrl(entry, entryIndex)" v-img-url="iconUrl(entry, entryIndex)" alt="" />
            <span v-else>{{ entry.icon || '◌' }}</span>
          </span>
          <span v-if="navigationType !== 'icon'" class="bottom-nav__label" :style="{ color: entryIndex === activeIndex ? activeColor : textColor }">{{ entry.text || '导航' }}</span>
        </button>
      </nav>
      <div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
    </div>
  </div>
</template>

<script>
export default {
  inject: ['diyModel'],
  props: ['item', 'index', 'selectedIndex'],
  computed: {
    style() { return this.item?.style || {}; },
    visibleItems() { return (this.item?.data || []).filter((entry) => !entry.hide); },
    navigationType() { return ['icon-text', 'icon', 'text'].includes(this.style.navigationType) ? this.style.navigationType : 'icon-text'; },
    isFloating() { return this.style.positionType === 'float'; },
    isFixed() { return this.style.positionType !== 'float'; },
    activeIndex() {
      const active = Number(this.item?.params?.activeIndex ?? 0);
      return active >= 0 && active < this.visibleItems.length ? active : 0;
    },
    activeColor() { return this.style.themeMode === 'custom' ? this.style.activeColor || '#f62c2c' : '#f62c2c'; },
    textColor() { return this.style.themeMode === 'custom' ? this.style.textColor || '#282828' : '#282828'; },
    wrapStyle() { return { padding: `0 ${Number(this.style.pagePadding || 0)}px ${Number(this.style.bottomSpacing || 0)}px` }; },
    navStyle() {
      const radius = Number(this.style.radius || 0);
      return {
        background: this.style.background || 'rgba(255,255,255,0.96)',
        borderRadius: this.isFloating ? `${radius}px` : '0',
        boxShadow: this.isFloating ? '0 5px 18px rgba(24, 39, 75, .14)' : '0 -1px 0 rgba(0,0,0,.05)',
        paddingTop: `${Number(this.style.paddingTop || 0)}px`,
        paddingBottom: `${Number(this.style.paddingBottom || 0)}px`,
      };
    },
  },
  methods: {
    iconUrl(entry, entryIndex) {
      return entryIndex === this.activeIndex ? entry.selectedImgUrl || entry.imgUrl : entry.unselectedImgUrl || entry.imgUrl;
    },
    diyEditer(index) { this.diyModel?.onEditer(index); },
    diyDeleteItem(index) { this.diyModel?.onDeleleItem(index); },
  },
};
</script>

<style lang="scss" scoped>
.bottom-nav-wrap { width: 100%; z-index: 4; }
.bottom-nav-wrap--fixed { position: sticky; bottom: 0; margin-top: auto; }
.bottom-nav { display: flex; min-height: 58px; align-items: stretch; }
.bottom-nav__item { display: flex; min-width: 0; flex: 1; padding: 5px 2px; align-items: center; justify-content: center; border: 0; background: transparent; flex-direction: column; }
.bottom-nav__icon { display: flex; width: 25px; height: 25px; align-items: center; justify-content: center; color: #7b8495; font-size: 22px; line-height: 1; }
.bottom-nav__icon img { width: 100%; height: 100%; object-fit: contain; }
.bottom-nav__label { margin-top: 2px; overflow: hidden; font-size: 11px; line-height: 16px; text-overflow: ellipsis; white-space: nowrap; }
.bottom-nav--icon .bottom-nav__item { padding: 10px 2px; }
.bottom-nav--text .bottom-nav__item { padding: 12px 2px; }
.bottom-nav--text .bottom-nav__label { margin-top: 0; font-size: 13px; }
</style>
