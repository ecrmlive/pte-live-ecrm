<template>
  <div class="diy-hotspot-host" :style="outerStyle" @click.stop="diyEditer(index)">
    <div class="drag optional" :class="{ selected: index === selectedIndex }">
      <div class="diy-hotspot" :style="cardStyle">
        <section v-if="modules.seckill.enabled" class="hotspot-item hotspot-products hotspot-seckill">
          <div class="hotspot-head">
            <strong>{{ modules.seckill.title }}</strong>
            <span>{{ modules.seckill.subtitle }}</span>
          </div>
          <div class="hotspot-product-grid">
            <div v-for="(product, productIndex) in modules.seckill.products" :key="productIndex" class="hotspot-product">
              <div class="hotspot-product-image">
                <img v-if="product.image" v-img-url="product.image" :alt="product.name" />
                <span v-else>秒杀</span>
              </div>
              <em v-if="product.price">¥{{ product.price }}</em>
            </div>
          </div>
        </section>

        <section v-if="modules.groupbuy.enabled" class="hotspot-item hotspot-products hotspot-groupbuy">
          <div class="hotspot-head">
            <strong>{{ modules.groupbuy.title }}</strong>
            <span>{{ modules.groupbuy.subtitle }}</span>
          </div>
          <div class="hotspot-product-grid">
            <div v-for="(product, productIndex) in modules.groupbuy.products" :key="productIndex" class="hotspot-product">
              <div class="hotspot-product-image">
                <img v-if="product.image" v-img-url="product.image" :alt="product.name" />
                <span v-else>团购</span>
              </div>
              <em v-if="product.price">¥{{ product.price }}</em>
            </div>
          </div>
        </section>

        <section v-if="modules.category.enabled" class="hotspot-item hotspot-promo hotspot-category">
          <div class="hotspot-promo-copy">
            <strong>{{ modules.category.title }}</strong>
            <span>{{ modules.category.subtitle }}</span>
          </div>
          <div class="hotspot-promo-image">
            <img v-if="modules.category.image" v-img-url="modules.category.image" :alt="modules.category.title" />
            <span v-else>分类</span>
          </div>
        </section>

        <section v-if="modules.newcomer.enabled" class="hotspot-item hotspot-promo hotspot-newcomer">
          <div class="hotspot-promo-copy">
            <strong>{{ modules.newcomer.title }}</strong>
            <span>{{ modules.newcomer.subtitle }}</span>
          </div>
          <div class="hotspot-promo-image">
            <img v-if="modules.newcomer.image" v-img-url="modules.newcomer.image" :alt="modules.newcomer.title" />
            <span v-else>新人</span>
          </div>
        </section>
      </div>
      <div class="btn-edit-del">
        <div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div>
      </div>
    </div>
  </div>
</template>

<script>
const defaultProduct = () => ({ image: '', name: '', price: '' });

const defaultModules = () => ({
  category: { enabled: true, image: '', subtitle: '全屋｜厨卫｜保洁', title: '日常保洁' },
  groupbuy: {
    enabled: true,
    products: [defaultProduct(), defaultProduct()],
    subtitle: '超低价 随时退',
    title: '热销团购',
  },
  newcomer: { enabled: true, image: '', subtitle: '下单咨询即领取45元券', title: '新人专享' },
  seckill: {
    enabled: true,
    products: [defaultProduct(), defaultProduct()],
    subtitle: '4折秒杀',
    title: '限时秒杀',
  },
});

export default {
  inject: ['diyModel'],
  props: ['item', 'index', 'selectedIndex'],
  computed: {
    modules() {
      const defaults = defaultModules();
      const saved = this.item?.params?.modules || {};
      return Object.keys(defaults).reduce((result, key) => {
        const savedModule = saved[key] || {};
        const module = { ...defaults[key], ...savedModule };
        if (defaults[key].products) {
          const savedProducts = Array.isArray(savedModule.products) ? savedModule.products : [];
          module.products = defaults[key].products.map((product, index) => ({
            ...product,
            ...(savedProducts[index] || {}),
          }));
        }
        result[key] = module;
        return result;
      }, {});
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
          background: style.bgcolor || '#ffffff',
          borderBottomLeftRadius: `${Number(style.bottomLeftRadio || 0)}px`,
          borderBottomRightRadius: `${Number(style.bottomRightRadio || 0)}px`,
          borderTopLeftRadius: `${Number(style.topLeftRadio || 0)}px`,
          borderTopRightRadius: `${Number(style.topRightRadio || 0)}px`,
          boxShadow: shadow,
        };
      }
      return {
        background: style.bgcolor || '#ffffff',
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
.diy-hotspot-host {
  box-sizing: border-box;
}

.diy-hotspot {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 6px;
  box-sizing: border-box;
  padding: 7px;
}

.hotspot-item {
  min-width: 0;
  overflow: hidden;
  border-radius: 7px;
}

.hotspot-products {
  min-height: 112px;
  padding: 8px 8px 6px;
  background: #fff9f7;
}

.hotspot-seckill {
  background: linear-gradient(140deg, #fff7f3 0%, #fff0ea 100%);
}

.hotspot-groupbuy {
  background: linear-gradient(140deg, #f9fbff 0%, #f0f6ff 100%);
}

.hotspot-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 4px;
  margin-bottom: 6px;

  strong {
    overflow: hidden;
    color: #252525;
    font-size: 14px;
    font-weight: 700;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  span {
    overflow: hidden;
    color: #f65f4b;
    font-size: 9px;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.hotspot-groupbuy .hotspot-head span {
  color: #6b7280;
}

.hotspot-product-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 5px;
}

.hotspot-product {
  position: relative;
  min-width: 0;
}

.hotspot-product-image {
  display: flex;
  height: 64px;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.7);
  color: #c7a9a0;
  font-size: 10px;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.hotspot-product em {
  display: block;
  margin-top: 3px;
  color: #f04438;
  font-size: 11px;
  font-style: normal;
  font-weight: 700;
  text-align: center;
}

.hotspot-promo {
  display: flex;
  min-height: 54px;
  align-items: center;
  justify-content: space-between;
  padding: 8px 10px;
}

.hotspot-category {
  background: linear-gradient(120deg, #f5fcf9 0%, #ebfaf3 100%);
}

.hotspot-newcomer {
  background: linear-gradient(120deg, #fff8f1 0%, #fff0e6 100%);
}

.hotspot-promo-copy {
  min-width: 0;

  strong,
  span {
    display: block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  strong {
    color: #333;
    font-size: 13px;
  }

  span {
    margin-top: 4px;
    color: #87909d;
    font-size: 9px;
  }
}

.hotspot-promo-image {
  display: flex;
  width: 42px;
  height: 42px;
  flex: none;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  margin-left: 6px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.68);
  color: #9ba7b4;
  font-size: 10px;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}
</style>
