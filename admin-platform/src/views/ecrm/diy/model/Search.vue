<template>
  <div
    class="drag optional diy-search-host"
    :class="{ selected: index === selectedIndex, 'diy-search-host--fixed': isFixed }"
    :style="hostStyle"
    @click.stop="diyEditer(index)"
  >
    <div class="diy-search navigation d-s-c" :style="cardStyle">
      <img v-if="leadingMode === 'logo'" class="logo-img" v-img-url="params.toplogo" alt="" />
      <div v-else-if="leadingMode === 'title'" class="search-leading search-leading--title" :style="titleStyle">
        {{ params.title }}
      </div>
      <div v-else-if="leadingMode === 'location'" class="search-leading search-leading--location" :style="titleStyle">
        <el-icon><Location /></el-icon>
        <span>{{ params.locationText }}</span>
        <el-icon class="location-arrow"><ArrowRight /></el-icon>
      </div>
      <div class="phone-top-search-box d-s-c" :style="searchBoxStyle">
        <el-icon class="search-icon" :color="searchIconColor"><Search /></el-icon>
        <span class="search-keyword" :style="activeHotWord ? hotWordStyle : placeholderStyle">
          {{ activeHotWord || searchText }}
        </span>
      </div>
    </div>
    <div class="btn-edit-del">
      <div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div>
    </div>
  </div>
</template>

<script>
import { ArrowRight, Location, Search } from '@element-plus/icons-vue';

export default {
  components: { ArrowRight, Location, Search },
  inject: ['diyModel'],
  props: ['item', 'index', 'selectedIndex'],
  data() {
    return {
      hotWordIndex: 0,
      hotWordTimer: null,
    };
  },
  computed: {
    params() {
      return this.item?.params || {};
    },
    style() {
      return this.item?.style || {};
    },
    isFixed() {
      return this.params.display_mode === 'fixed';
    },
    leadingMode() {
      if (this.params.style_type === 'logo') return 'logo';
      if (this.params.style_type === 'location') return 'location';
      if (this.params.style_type === 'title') return 'title';
      if (this.params.title_type === 'image') return 'logo';
      if (this.params.title_type === 'title' || this.params.title_type === 'text') return 'title';
      return 'location';
    },
    searchText() {
      return this.params.searchText || '搜索商品';
    },
    hotWords() {
      return Array.isArray(this.params.hotWords)
        ? this.params.hotWords.filter((word) => typeof word === 'string' && word.trim())
        : [];
    },
    activeHotWord() {
      return this.hotWords[this.hotWordIndex % Math.max(this.hotWords.length, 1)] || '';
    },
    hostStyle() {
      const float = Number(this.style.float || 0);
      const paddingRight = this.style.paddingRight ?? this.style.paddingLeft ?? 0;
      return {
        background: this.style.background || '#fff1f0',
        paddingLeft: `${Number(this.style.paddingLeft || 0)}px`,
        paddingRight: `${Number(paddingRight || 0)}px`,
        paddingTop: `${Number(this.style.paddingTop || 0)}px`,
        paddingBottom: `${Number(this.style.paddingBottom || 0)}px`,
        transform: float > 0 ? `translateY(-${float}px)` : undefined,
        marginBottom: float > 0 ? `-${float}px` : undefined,
      };
    },
    cardStyle() {
      const allRadius = Number(this.style.topRadio || 0);
      const individual = this.style.radiusMode === 'individual';
      return {
        background: this.style.bgcolor || '#fff1f0',
        borderTopLeftRadius: `${individual ? Number(this.style.topLeftRadio || 0) : allRadius}px`,
        borderTopRightRadius: `${individual ? Number(this.style.topRightRadio || 0) : allRadius}px`,
        borderBottomLeftRadius: `${individual ? Number(this.style.bottomLeftRadio || 0) : Number(this.style.bottomRadio ?? allRadius)}px`,
        borderBottomRightRadius: `${individual ? Number(this.style.bottomRightRadio || 0) : Number(this.style.bottomRadio ?? allRadius)}px`,
        boxShadow: this.style.shadow === 'on' ? '0 6px 18px rgba(0, 0, 0, 0.12)' : 'none',
      };
    },
    titleStyle() {
      const textStyle = this.style.titleTextStyle || 'normal';
      return {
        color: this.style.titleTextColor || '#333333',
        fontSize: `${Number(this.style.titleTextSize || 14)}px`,
        fontStyle: textStyle === 'italic' ? 'italic' : 'normal',
        fontWeight: textStyle === 'bold' ? 700 : 400,
        justifyContent: this.style.titleAlign === 'right' ? 'flex-end' : this.style.titleAlign === 'center' ? 'center' : 'flex-start',
      };
    },
    searchBoxStyle() {
      return {
        background: this.style.searchBackGround || '#ffffff',
        color: this.style.searchColor || '#cccccc',
      };
    },
    placeholderStyle() {
      return { color: this.style.searchColor || '#cccccc' };
    },
    hotWordStyle() {
      return { color: this.style.hotWordColor || '#666666' };
    },
    searchIconColor() {
      return this.style.searchColor || '#999999';
    },
  },
  watch: {
    hotWords: {
      deep: true,
      handler() {
        this.hotWordIndex = 0;
        this.startHotWordTimer();
      },
    },
    'params.hotWordInterval'() {
      this.startHotWordTimer();
    },
  },
  mounted() {
    this.startHotWordTimer();
  },
  beforeUnmount() {
    this.stopHotWordTimer();
  },
  methods: {
    startHotWordTimer() {
      this.stopHotWordTimer();
      if (this.hotWords.length < 2) return;
      const seconds = Math.min(Math.max(Number(this.params.hotWordInterval || 3), 1), 10);
      this.hotWordTimer = window.setInterval(() => {
        this.hotWordIndex = (this.hotWordIndex + 1) % this.hotWords.length;
      }, seconds * 1000);
    },
    stopHotWordTimer() {
      if (this.hotWordTimer) {
        window.clearInterval(this.hotWordTimer);
        this.hotWordTimer = null;
      }
    },
    diyEditer(index) {
      this.diyModel?.onEditer(index);
    },
    diyDeleteItem(index) {
      this.diyModel?.onDeleleItem(index);
    },
  },
};
</script>

<style lang="scss" scoped>
.diy-search-host--fixed {
  position: sticky;
  top: 0;
  z-index: 20;
}

.diy-search {
  min-height: 58px;
  padding: 10px 12px;

  .logo-img {
    display: block;
    width: 39px;
    height: 32px;
    margin-right: 10px;
    object-fit: contain;
  }

  .search-leading {
    display: flex;
    align-items: center;
    flex: 0 0 auto;
    min-width: 0;
    margin-right: 8px;
    overflow: hidden;
    white-space: nowrap;
  }

  .search-leading--title {
    max-width: 34%;
  }

  .search-leading--location {
    max-width: 42%;
    gap: 4px;

    span {
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .location-arrow {
      flex: 0 0 auto;
      margin-left: -2px;
      font-size: 14px;
    }
  }

  .phone-top-search-box {
    flex: 1;
    min-width: 0;
    height: 38px;
    border-radius: 22px;
    font-size: 14px;
    line-height: 38px;
    padding: 0 14px;
  }

  .search-icon {
    flex: 0 0 auto;
    margin-right: 10px;
    font-size: 18px;
  }

  .search-keyword {
    display: block;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
</style>
