<template>
  <div class="diy-discount-group-host" :style="outerStyle" @click.stop="diyEditer(index)">
    <div class="drag optional" :class="{ selected: index === selectedIndex }">
      <section class="diy-discount-group" :style="cardStyle">
        <header class="discount-group-header">
          <div class="discount-group-title">
            <strong>{{ params.title }}</strong>
            <span>{{ params.promotion }}</span>
          </div>
          <div class="discount-group-slogan">
            <img v-if="params.iconImage" v-img-url="params.iconImage" alt="活动图标" />
            <span v-else class="discount-group-icon">★</span>
            <span>{{ params.slogan }}</span>
          </div>
        </header>

        <div class="discount-group-items">
          <article v-for="(entry, entryIndex) in items" v-show="entry.enabled" :key="entryIndex" class="discount-group-item">
            <div class="discount-group-cover">
              <img v-if="entry.image" v-img-url="entry.image" :alt="entry.title" />
              <span v-else>{{ entry.title.slice(0, 2) }}</span>
            </div>
            <strong>{{ entry.title }}</strong>
            <em v-if="entry.price">¥{{ entry.price }}</em>
          </article>
        </div>
      </section>
      <div class="btn-edit-del">
        <div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div>
      </div>
    </div>
  </div>
</template>

<script>
const defaultItems = () => [
  { enabled: true, image: '', price: '', productId: 0, productName: '', title: '全球美妆' },
  { enabled: true, image: '', price: '', productId: 0, productName: '', title: '大牌鞋包' },
  { enabled: true, image: '', price: '', productId: 0, productName: '', title: '数码产品' },
  { enabled: true, image: '', price: '', productId: 0, productName: '', title: '精品腕表' },
];

export default {
  inject: ['diyModel'],
  props: ['item', 'index', 'selectedIndex'],
  computed: {
    params() {
      const saved = this.item?.params || {};
      const defaults = defaultItems();
      const items = Array.isArray(saved.items) ? saved.items : [];
      return {
        iconImage: '',
        promotion: '券后低至7.3折',
        slogan: '真低价 放心买',
        title: '心动购物季',
        ...saved,
        items: defaults.map((entry, itemIndex) => ({ ...entry, ...(items[itemIndex] || {}) })),
      };
    },
    items() {
      return this.params.items;
    },
    outerStyle() {
      const style = this.item?.style || {};
      return {
        background: style.background || 'transparent',
        marginTop: `${Number(style.marginTop || 0)}px`,
        paddingBottom: `${Number(style.paddingBottom || 0)}px`,
        paddingLeft: `${Number(style.paddingLeft || 10)}px`,
        paddingRight: `${Number(style.paddingRight ?? style.paddingLeft ?? 10)}px`,
        paddingTop: `${Number(style.paddingTop || 9)}px`,
        transform: `translateY(-${Number(style.float || 0)}px)`,
      };
    },
    cardStyle() {
      const style = this.item?.style || {};
      const radius = Number(style.cardRadius || 0);
      const shadow = style.cardShadow === 'on' ? '0 6px 18px rgba(15, 23, 42, 0.12)' : 'none';
      if (style.radiusMode === 'individual') {
        return {
          background: style.bgcolor || '#f5f5f5',
          borderBottomLeftRadius: `${Number(style.bottomLeftRadio || 0)}px`,
          borderBottomRightRadius: `${Number(style.bottomRightRadio || 0)}px`,
          borderTopLeftRadius: `${Number(style.topLeftRadio || 0)}px`,
          borderTopRightRadius: `${Number(style.topRightRadio || 0)}px`,
          boxShadow: shadow,
        };
      }
      return {
        background: style.bgcolor || '#f5f5f5',
        borderRadius: `${radius}px`,
        boxShadow: shadow,
      };
    },
  },
  methods: {
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
.diy-discount-group-host {
  box-sizing: border-box;
}

.diy-discount-group {
  box-sizing: border-box;
  overflow: hidden;
  padding: 5px;
}

.discount-group-header {
  display: flex;
  height: 31px;
  align-items: center;
  justify-content: space-between;
  gap: 6px;
  padding: 0 8px;
  color: #fff;
  background: linear-gradient(90deg, #ff2d2d 0%, #ff4d22 100%);
}

.discount-group-title,
.discount-group-slogan {
  display: flex;
  min-width: 0;
  align-items: center;
  gap: 4px;
  overflow: hidden;
  white-space: nowrap;
}

.discount-group-title strong {
  overflow: hidden;
  font-size: 13px;
  font-weight: 700;
  text-overflow: ellipsis;
}

.discount-group-title span,
.discount-group-slogan span {
  overflow: hidden;
  font-size: 9px;
  text-overflow: ellipsis;
}

.discount-group-slogan {
  flex: none;

  img {
    width: 13px;
    height: 13px;
    border-radius: 50%;
    object-fit: cover;
  }
}

.discount-group-icon {
  color: #fff7c9;
}

.discount-group-items {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 3px;
  padding-top: 3px;
}

.discount-group-item {
  display: flex;
  min-width: 0;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  overflow: hidden;
  padding: 4px 2px 5px;
  border-radius: 3px;
  background: #fff;
}

.discount-group-cover {
  display: flex;
  width: 100%;
  height: 55px;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-radius: 2px;
  color: #f45b3f;
  background: #fff3ed;
  font-size: 10px;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.discount-group-item strong {
  width: 100%;
  overflow: hidden;
  color: #ef3d2d;
  font-size: 9px;
  font-weight: 600;
  line-height: 1.2;
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.discount-group-item em {
  color: #f33;
  font-size: 10px;
  font-style: normal;
  font-weight: 700;
}
</style>
