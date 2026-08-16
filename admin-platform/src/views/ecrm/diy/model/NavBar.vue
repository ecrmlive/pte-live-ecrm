<template>
  <div @click.stop="diyEditer(index)" :style="outerStyle">
    <div class="drag optional" :class="{ selected: index === selectedIndex }">
      <div class="diy-navBar" :style="cardStyle">
        <div v-if="item.params && item.params.title" class="nav-title">{{ item.params.title }}</div>
        <ul
          class="list"
          :class="[`column-${rowsNum}`, `mode-${navigationType}`]"
          @touchend="onTouchEnd"
          @touchstart="onTouchStart"
        >
          <li v-for="(navBar, navIndex) in currentItems" :key="navIndex" class="item">
            <div v-if="navigationType !== 'text'" class="item-image" :style="iconStyle">
              <span
                v-if="navBar.icon"
                class="item-icon"
                :style="{ background: navBar.iconBg || '#eef4ff', color: navBar.iconColor || '#1677ff' }"
              >{{ navBar.icon }}</span>
              <img v-else-if="navBar.imgUrl" v-img-url="navBar.imgUrl" alt="" />
              <span v-else class="item-image__empty">图标</span>
            </div>
            <div
              v-if="navigationType !== 'icon'"
              class="item-text text-ellipsis"
              :style="{ color: navBar.color || textColor }"
            >{{ navBar.text }}</div>
          </li>
        </ul>
        <div v-if="isPaged && totalPages > 1" class="nav-pagination">
          <span
            v-for="page in totalPages"
            :key="page"
            :class="{ active: page - 1 === currentPage }"
            @click.stop="currentPage = page - 1"
          ></span>
        </div>
      </div>
      <div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
    </div>
  </div>
</template>

<script>
export default {
  inject: ['diyModel'],
  props: ['item', 'index', 'selectedIndex'],
  data() {
    return { currentPage: 0, touchStartX: 0 };
  },
  computed: {
    style() {
      return this.item?.style || {};
    },
    rowsNum() {
      const num = Number(this.style.rowsNum);
      return Number.isFinite(num) && num >= 3 && num <= 5 ? num : 5;
    },
    navigationType() {
      return ['icon-text', 'icon', 'text'].includes(this.style.navigationType)
        ? this.style.navigationType
        : 'icon-text';
    },
    isPaged() {
      return this.style.displayMode === 'page';
    },
    pageSize() {
      return this.rowsNum * 2;
    },
    visibleItems() {
      return (this.item?.data || []).filter((entry) => !entry.hide);
    },
    totalPages() {
      if (!this.isPaged) return 1;
      return Math.max(1, Math.ceil(this.visibleItems.length / this.pageSize));
    },
    currentItems() {
      if (!this.isPaged) return this.visibleItems;
      const page = Math.min(this.currentPage, this.totalPages - 1);
      return this.visibleItems.slice(page * this.pageSize, (page + 1) * this.pageSize);
    },
    textColor() {
      return this.style.textColor || '#333333';
    },
    outerStyle() {
      return {
        background: this.style.background || 'transparent',
        marginTop: `${this.style.marginTop || 0}px`,
        paddingLeft: `${this.style.paddingLeft || 0}px`,
        paddingRight: `${this.style.paddingRight ?? this.style.paddingLeft ?? 0}px`,
        paddingTop: `${this.style.paddingTop || 0}px`,
        paddingBottom: `${this.style.paddingBottom || 0}px`,
      };
    },
    cardStyle() {
      const start = this.style.bgcolor || 'rgba(255, 255, 255, 0)';
      const end = this.style.bgcolorEnd || start;
      const radius = Number(this.style.cardRadius || 0);
      const all = this.style.cardRadiusMode !== 'individual';
      return {
        background: start === end ? start : `linear-gradient(180deg, ${start} 0%, ${end} 100%)`,
        borderTopLeftRadius: `${all ? radius : this.style.cardTopLeftRadius || 0}px`,
        borderTopRightRadius: `${all ? radius : this.style.cardTopRightRadius || 0}px`,
        borderBottomRightRadius: `${all ? radius : this.style.cardBottomRightRadius || 0}px`,
        borderBottomLeftRadius: `${all ? radius : this.style.cardBottomLeftRadius || 0}px`,
        boxShadow: this.style.cardShadow === 'on' ? '0 8px 18px rgba(20, 37, 63, 0.12)' : 'none',
        transform: `translateY(-${this.style.float || 0}px)`,
      };
    },
    iconStyle() {
      const radius = Number(this.style.iconRadius ?? this.style.topRadio ?? 8);
      const all = this.style.iconRadiusMode !== 'individual';
      return {
        borderTopLeftRadius: `${all ? radius : this.style.iconTopLeftRadius || 0}px`,
        borderTopRightRadius: `${all ? radius : this.style.iconTopRightRadius || 0}px`,
        borderBottomRightRadius: `${all ? radius : this.style.iconBottomRightRadius || 0}px`,
        borderBottomLeftRadius: `${all ? radius : this.style.iconBottomLeftRadius || 0}px`,
        boxShadow: this.style.iconShadow === 'on' ? '0 4px 10px rgba(20, 37, 63, 0.16)' : 'none',
      };
    },
  },
  watch: {
    visibleItems() {
      if (this.currentPage >= this.totalPages) this.currentPage = 0;
    },
  },
  methods: {
    diyEditer(index) {
      this.diyModel?.onEditer(index);
    },
    diyDeleteItem(index) {
      this.diyModel?.onDeleleItem(index);
    },
    previousPage() {
      this.currentPage = (this.currentPage - 1 + this.totalPages) % this.totalPages;
    },
    nextPage() {
      this.currentPage = (this.currentPage + 1) % this.totalPages;
    },
    onTouchStart(event) {
      this.touchStartX = event.changedTouches?.[0]?.clientX || 0;
    },
    onTouchEnd(event) {
      if (!this.isPaged) return;
      const endX = event.changedTouches?.[0]?.clientX || this.touchStartX;
      const delta = endX - this.touchStartX;
      if (Math.abs(delta) < 32 || this.totalPages <= 1) return;
      delta > 0 ? this.previousPage() : this.nextPage();
    },
  },
};
</script>

<style lang="scss" scoped>
.diy-navBar .list {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
}

.nav-title {
  padding: 10px 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #252b3a;
}

.diy-navBar .list .item {
  display: flex;
  min-width: 0;
  padding: 10px 0;
  align-items: center;
  justify-content: center;
  flex-direction: column;
}

.diy-navBar .list.column-3 .item { width: 33.333333%; }
.diy-navBar .list.column-4 .item { width: 25%; }
.diy-navBar .list.column-5 .item { width: 20%; }

.diy-navBar .list .item-image {
  display: flex;
  width: min(60%, 58px);
  aspect-ratio: 1;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.diy-navBar .list .item-image img,
.diy-navBar .list .item-icon {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.diy-navBar .list .item-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: inherit;
  font-size: 20px;
  font-weight: 700;
}

.item-image__empty {
  display: flex;
  width: 100%;
  height: 100%;
  align-items: center;
  justify-content: center;
  background: #f1f5fb;
  color: #a9b4c6;
  font-size: 10px;
}

.diy-navBar .list .item-text {
  width: 100%;
  padding: 5px 4px 0;
  overflow: hidden;
  font-size: 12px;
  line-height: 18px;
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.diy-navBar .list.mode-text .item { padding: 9px 2px; }
.diy-navBar .list.mode-text .item-text { padding-top: 0; }
.diy-navBar .list.mode-icon .item { padding-bottom: 8px; }

.nav-pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 0 0 10px;
}

.nav-pagination span {
  width: 5px;
  height: 5px;
  cursor: pointer;
  border-radius: 999px;
  background: #d8deea;
}

.nav-pagination span.active {
  width: 13px;
  background: #1677ff;
}
</style>
