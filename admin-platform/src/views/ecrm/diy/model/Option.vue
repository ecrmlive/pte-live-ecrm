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
        :class="{ 'diy-option__item--active': tabIndex === 0 }"
        :style="tabStyle(tabIndex)"
        type="button"
      >
        <span>{{ tab.text || '选项卡' }}</span>
        <i v-if="displayStyle === '2' && tabIndex === 0" class="diy-option__arc" :style="arcStyle"></i>
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
  computed: {
    style() {
      return this.item?.style || {};
    },
    params() {
      return this.item?.params || {};
    },
    tabs() {
      const data = Array.isArray(this.item?.data) ? this.item.data : [];
      return data.length > 0
        ? data
        : [{ text: '首页' }, { text: '果蔬生鲜' }, { text: '健康医疗' }, { text: '非遗文创' }, { text: '优选茶叶' }];
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
        background: this.style.bgcolor || '#fff0f3',
        borderRadius: this.radius,
        boxShadow: this.style.shadow === 'on' ? '0 4px 14px rgba(15, 23, 42, 0.12)' : 'none',
      };
    },
    arcStyle() {
      return { borderColor: this.activeColor };
    },
  },
  methods: {
    diyEditer(index) {
      this.diyModel?.onEditer(index);
    },
    diyDeleteItem(index) {
      this.diyModel?.onDeleleItem(index);
    },
    tabStyle(index) {
      if (index !== 0) return {};
      if (this.displayStyle === '1') {
        return { borderBottomColor: this.activeColor, color: this.activeColor };
      }
      if (this.displayStyle === '2') {
        return { color: this.activeColor };
      }
      return { background: this.activeColor, color: this.activeText };
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
    gap: 28px;
    min-height: 56px;
    overflow-x: auto;
    padding: 7px 16px;
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
    cursor: default;
    display: inline-flex;
    flex: none;
    font-size: 16px;
    height: 42px;
    justify-content: center;
    line-height: 1;
    padding: 0;
    pointer-events: none;
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
    gap: 28px;

    .diy-option__item {
      font-size: 17px;
      font-weight: 500;
    }

    .diy-option__item--active {
      font-weight: 600;
    }
  }

  &__list--3 {
    gap: 10px;

    .diy-option__item {
      border-radius: 999px;
      font-size: 15px;
      font-weight: 500;
      padding: 0 14px;
    }
  }

  &__arc {
    border: 0;
    border-bottom: 3px solid;
    border-radius: 0 0 50% 50%;
    bottom: 1px;
    height: 8px;
    left: 18%;
    position: absolute;
    transform: none;
    width: 64%;
  }
}
</style>
