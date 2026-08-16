<template>
  <div
    class="diy-option"
    :class="{ selected: index === selectedIndex, 'diy-option--sticky': params.topUp === '1' }"
    :style="hostStyle"
    @click.stop="diyEditer(index)"
  >
    <div class="diy-option__list" :class="`diy-option__list--${displayStyle}`" :style="listStyle">
      <button
        v-for="(tab, tabIndex) in tabs"
        :key="tab._diyTabId ?? `${tab.text}-${tabIndex}`"
        class="diy-option__item"
        :class="{ 'diy-option__item--active': activeIndex === tabIndex }"
        :style="tabStyle(tabIndex)"
        type="button"
        @click.stop="setActive(tabIndex)"
      >
        <span>{{ tab.text || '选项卡' }}</span>
        <i v-if="displayStyle === '3' && activeIndex === tabIndex" class="diy-option__arc" :style="arcStyle"></i>
      </button>
    </div>
    <div class="btn-edit-del">
      <div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div>
    </div>
  </div>
</template>

<script>
export default {
  inject: ['diyModel'],
  props: ['item', 'index', 'selectedIndex'],
  data() {
    return { activeIndex: 0 };
  },
  computed: {
    style() {
      return this.item?.style || {};
    },
    params() {
      return this.item?.params || {};
    },
    tabs() {
      const data = Array.isArray(this.item?.data) ? this.item.data : [];
      return data.length > 0 ? data : [{ text: '首页' }];
    },
    displayStyle() {
      return String(this.params.type || '2');
    },
    activeColor() {
      return this.style.themeType === 'custom'
        ? this.style.activeColor || '#ff4d7d'
        : this.style.activeColor || '#ff4d7d';
    },
    activeText() {
      return this.style.themeType === 'custom' ? this.style.activeText || '#ffffff' : '#ffffff';
    },
    radius() {
      const value = Number(this.style.radius || 0);
      if (this.style.radiusMode !== 'individual') return `${value}px`;
      return [
        this.style.topLeftRadius ?? value,
        this.style.topRightRadius ?? value,
        this.style.bottomRightRadius ?? value,
        this.style.bottomLeftRadius ?? value,
      ]
        .map((item) => `${Number(item)}px`)
        .join(' ');
    },
    hostStyle() {
      const paddingLeft = Number(this.style.paddingLeft || 0);
      const float = Number(this.style.float || 0);
      return {
        background: this.style.background || 'transparent',
        marginTop: `${Number(this.style.marginTop || 0)}px`,
        paddingTop: `${Number(this.style.paddingTop || 0)}px`,
        paddingRight: `${Number(this.style.paddingRight ?? paddingLeft)}px`,
        paddingBottom: `${Number(this.style.paddingBottom || 0)}px`,
        paddingLeft: `${paddingLeft}px`,
        transform: float ? `translateY(-${float}px)` : '',
        zIndex: this.params.topUp === '1' || float ? 3 : '',
      };
    },
    listStyle() {
      return {
        background: this.style.bgcolor || '#ffffff',
        borderRadius: this.radius,
        boxShadow: this.style.shadow === 'on' ? '0 4px 14px rgba(15, 23, 42, 0.12)' : 'none',
      };
    },
    arcStyle() {
      return { borderColor: this.activeColor };
    },
  },
  watch: {
    tabs: {
      deep: true,
      handler(tabs) {
        if (this.activeIndex >= tabs.length) this.activeIndex = 0;
      },
    },
  },
  methods: {
    diyEditer(index) {
      this.diyModel?.onEditer(index);
    },
    diyDeleteItem(index) {
      this.diyModel?.onDeleleItem(index);
    },
    setActive(index) {
      this.activeIndex = index;
    },
    tabStyle(index) {
      const active = this.activeIndex === index;
      if (!active) return {};
      if (this.displayStyle === '1') {
        return { borderBottomColor: this.activeColor, color: this.activeColor };
      }
      if (this.displayStyle === '2') {
        return { background: this.activeColor, color: this.activeText };
      }
      return { color: this.activeColor };
    },
  },
};
</script>

<style lang="scss" scoped>
.diy-option {
  position: relative;
  transition: transform 0.2s ease;

  &--sticky {
    position: sticky;
    top: 0;
  }

  &__list {
    align-items: center;
    display: flex;
    gap: 4px;
    min-height: 43px;
    overflow-x: auto;
    padding: 4px 12px;
    scrollbar-width: none;

    &::-webkit-scrollbar {
      display: none;
    }
  }

  &__item {
    align-items: center;
    background: transparent;
    border: 0;
    color: #333;
    cursor: pointer;
    display: inline-flex;
    flex: none;
    font-size: 14px;
    height: 30px;
    justify-content: center;
    line-height: 1;
    padding: 0 12px;
    position: relative;
    white-space: nowrap;
  }

  &__list--1 {
    gap: 14px;

    .diy-option__item {
      border-bottom: 2px solid transparent;
      border-radius: 0;
      font-size: 15px;
      font-weight: 500;
      padding: 0 2px;
    }
  }

  &__list--2 {
    .diy-option__item {
      border-radius: 999px;
    }
  }

  &__list--3 {
    .diy-option__item {
      font-size: 15px;
      font-weight: 500;
    }
  }

  &__arc {
    border: 2px solid;
    border-left-color: transparent !important;
    border-radius: 50%;
    border-right-color: transparent !important;
    border-top-color: transparent !important;
    bottom: 0;
    height: 10px;
    left: 50%;
    position: absolute;
    transform: translateX(-50%);
    width: 24px;
  }
}
</style>
