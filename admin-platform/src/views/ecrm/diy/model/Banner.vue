<template>
  <div
    class="drag optional diy-banner-host"
    :class="{ selected: index === selectedIndex }"
    :style="hostStyle"
    @click.stop="diyEditer(index)"
  >
    <div class="diy-banner" :class="`banner--${bannerStyle}`" :style="bannerStyleObject">
      <div class="img-list pr">
        <img
          v-if="activeBanner"
          v-img-url="activeBanner.imgUrl"
          class="banner-img"
          :style="imageStyle"
          alt="轮播图"
        />
        <div v-else class="banner-empty" :style="imageStyle">请添加轮播图片</div>

        <div class="dots" :class="[`dots--${indicatorPosition}`, `dots--${indicatorStyle}`]">
          <button
            v-for="(_banner, bannerIndex) in banners"
            :key="bannerIndex"
            type="button"
            :class="{ active: bannerIndex === activeIndex }"
            :style="indicatorDotStyle(bannerIndex)"
            :aria-label="`切换至第${bannerIndex + 1}张`"
            @click.stop="setActive(bannerIndex)"
          ></button>
        </div>
      </div>
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
    return {
      activeIndex: 0,
      timer: null,
    };
  },
  computed: {
    style() {
      return this.item?.style || {};
    },
    params() {
      return this.item?.params || {};
    },
    banners() {
      return Array.isArray(this.item?.data) ? this.item.data.filter((banner) => banner?.imgUrl) : [];
    },
    activeBanner() {
      return this.banners[this.activeIndex] || this.banners[0] || null;
    },
    bannerStyle() {
      return this.params.bannerStyle || 'style1';
    },
    indicatorStyle() {
      if (this.style.indicatorStyle) return this.style.indicatorStyle;
      if (this.style.imgShape === 'square') return 'style2';
      if (this.style.imgShape === 'rectangle') return 'style3';
      return 'style1';
    },
    indicatorPosition() {
      return this.style.indicatorPosition || this.style.btnShape || 'center';
    },
    indicatorColor() {
      return this.style.indicatorTone === 'custom'
        ? this.style.indicatorColor || '#ffffff'
        : this.style.btnColor || '#ffffff';
    },
    imageRadius() {
      const radius = Number(this.style.imageRadius ?? this.style.topRadio ?? 0);
      if (this.style.radiusMode !== 'individual') {
        return `${radius}px`;
      }
      return [
        Number(this.style.topLeftRadio ?? radius),
        Number(this.style.topRightRadio ?? radius),
        Number(this.style.bottomRightRadio ?? this.style.bottomRadio ?? radius),
        Number(this.style.bottomLeftRadio ?? this.style.bottomRadio ?? radius),
      ]
        .map((value) => `${value}px`)
        .join(' ');
    },
    hostStyle() {
      const paddingLeft = Number(this.style.paddingLeft || 0);
      const float = Number(this.style.float || 0);
      return {
        background: this.style.background || '#ffffff',
        marginTop: `${Number(this.style.marginTop || 0)}px`,
        paddingTop: `${Number(this.style.paddingTop || 0)}px`,
        paddingRight: `${Number(this.style.paddingRight ?? paddingLeft)}px`,
        paddingBottom: `${Number(this.style.paddingBottom || 0)}px`,
        paddingLeft: `${paddingLeft}px`,
        transform: float ? `translateY(-${float}px)` : '',
        zIndex: float ? 1 : '',
      };
    },
    bannerStyleObject() {
      return {
        height: `${Number(this.style.height || 340) * 0.5}px`,
        boxShadow:
          this.style.cardShadow === 'on' ? '0 4px 16px rgba(0, 0, 0, 0.16)' : 'none',
      };
    },
    imageStyle() {
      return {
        height: `${Number(this.style.height || 340) * 0.5}px`,
        borderRadius: this.imageRadius,
        boxShadow:
          this.style.imageShadow === 'on' ? '0 3px 12px rgba(0, 0, 0, 0.18)' : 'none',
      };
    },
  },
  watch: {
    banners: {
      handler() {
        if (this.activeIndex >= this.banners.length) this.activeIndex = 0;
        this.restartAutoPlay();
      },
      deep: true,
    },
  },
  mounted() {
    this.restartAutoPlay();
  },
  beforeUnmount() {
    this.stopAutoPlay();
  },
  methods: {
    diyEditer(index) {
      this.diyModel?.onEditer(index);
    },
    diyDeleteItem(index) {
      this.diyModel?.onDeleleItem(index);
    },
    indicatorDotStyle(bannerIndex) {
      return {
        background: this.indicatorColor,
        opacity: bannerIndex === this.activeIndex ? 1 : 0.45,
      };
    },
    setActive(index) {
      this.activeIndex = index;
      this.restartAutoPlay();
    },
    stopAutoPlay() {
      if (this.timer) {
        clearInterval(this.timer);
        this.timer = null;
      }
    },
    restartAutoPlay() {
      this.stopAutoPlay();
      if (this.banners.length < 2) return;
      this.timer = setInterval(() => {
        this.activeIndex = (this.activeIndex + 1) % this.banners.length;
      }, 3000);
    },
  },
};
</script>

<style lang="scss" scoped>
.diy-banner-host {
  position: relative;
  transition: transform 0.2s ease;
}

.diy-banner {
  overflow: hidden;
  position: relative;
  width: 100%;
}

.img-list {
  height: 100%;
}

.banner-img,
.banner-empty {
  display: block;
  width: 100%;
}

.banner-img {
  object-fit: cover;
}

.banner--style2 .banner-img {
  background: #fff;
  object-fit: contain;
}

.banner--style3 .banner-img {
  filter: saturate(1.06);
}

.banner-empty {
  align-items: center;
  background: #edf5ff;
  color: #9aa7b7;
  display: flex;
  justify-content: center;
}

.dots {
  align-items: center;
  bottom: 10px;
  display: flex;
  gap: 5px;
  position: absolute;
  z-index: 1;
}

.dots--left {
  left: 12px;
}

.dots--center {
  left: 50%;
  transform: translateX(-50%);
}

.dots--right {
  right: 12px;
}

.dots button {
  border: 0;
  cursor: pointer;
  padding: 0;
  transition: all 0.2s ease;
}

.dots--style1 button {
  border-radius: 50%;
  height: 7px;
  width: 7px;
}

.dots--style2 button {
  height: 8px;
  width: 8px;
}

.dots--style3 button {
  border-radius: 4px;
  height: 4px;
  width: 18px;
}

.dots--style4 button {
  border-radius: 4px;
  height: 3px;
  width: 12px;
}

.dots--style4 button.active {
  width: 24px;
}
</style>
