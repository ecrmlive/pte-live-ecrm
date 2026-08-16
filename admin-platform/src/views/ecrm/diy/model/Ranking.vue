<template>
  <div class="ranking-host" :style="hostStyle" @click.stop="diyEditer(index)">
    <section class="ranking-card" :style="cardStyle">
      <header class="ranking-head">
        <div class="ranking-title">
          <img v-if="params.titleType === 'image' && params.titleImage" v-img-url="params.titleImage" alt="排行榜标题" />
          <strong v-else>{{ params.title }}</strong>
        </div>
        <span class="ranking-more" :style="{ color: style.moreColor }">{{ params.more }} ›</span>
      </header>
      <div class="ranking-boards" :class="`ranking-boards--${boards.length}`">
        <article v-for="(board, boardIndex) in boards" :key="boardIndex" class="ranking-board" :style="boardStyle">
          <h4 :style="{ color: style.boardTitleColor }">{{ board.icon || '🔥' }} {{ board.title }}</h4>
          <div class="ranking-products">
            <div v-for="(product, productIndex) in board.products" :key="productIndex" class="ranking-product">
              <span class="ranking-number" :class="`ranking-number--${productIndex + 1}`">{{ productIndex + 1 }}</span>
              <img v-if="product.image" v-img-url="product.image" alt="" />
              <div v-else class="ranking-cover">商品</div>
              <div class="ranking-product__info">
                <b>{{ product.name || '幸运美物' }}</b>
                <em :style="{ color: style.priceColor }">¥{{ product.price || '350.00' }}</em>
              </div>
            </div>
          </div>
        </article>
      </div>
    </section>
    <div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
  </div>
</template>

<script>
const fallbackBoards = () => [
  { icon: '🔥', title: '销量榜', products: [{ name: '幸运美物', image: '', price: '350.00' }, { name: '幸运美物', image: '', price: '350.00' }, { name: '幸运美物', image: '', price: '350.00' }] },
  { icon: '🔥', title: '好评榜', products: [{ name: '幸运美物', image: '', price: '350.00' }, { name: '幸运美物', image: '', price: '350.00' }, { name: '幸运美物', image: '', price: '350.00' }] },
];

export default {
  inject: ['diyModel'],
  props: ['item', 'index', 'selectedIndex'],
  computed: {
    params() { return { title: '排行榜', more: '更多', titleType: 'text', titleImage: '', ...(this.item?.params || {}) }; },
    style() { return { background: '#f5f5f5', cardBackground: '#ffffff', boardBackground: '#fceae9', boardTitleColor: '#ff4c8d', moreColor: '#999999', priceColor: '#ff4c8d', paddingTop: 10, paddingBottom: 10, paddingLeft: 10, marginTop: 10, radius: 10, boardRadius: 8, shadow: 'off', ...(this.item?.style || {}) }; },
    boards() { return Array.isArray(this.item?.data) && this.item.data.length ? this.item.data : fallbackBoards(); },
    hostStyle() { return { background: this.style.background, padding: `${Number(this.style.paddingTop)}px ${Number(this.style.paddingLeft)}px ${Number(this.style.paddingBottom)}px`, marginTop: `${Number(this.style.marginTop)}px` }; },
    cardStyle() { return { background: this.style.cardBackground, borderRadius: `${Number(this.style.radius)}px`, boxShadow: this.style.shadow === 'on' ? '0 6px 18px rgba(15, 23, 42, .12)' : 'none' }; },
    boardStyle() { return { background: this.style.boardBackground, borderRadius: `${Number(this.style.boardRadius)}px` }; },
  },
  methods: { diyEditer(i) { this.diyModel?.onEditer(i); }, diyDeleteItem(i) { this.diyModel?.onDeleleItem(i); } },
};
</script>

<style lang="scss" scoped>
.ranking-host { position: relative; box-sizing: border-box; width: 100%; }
.ranking-card { overflow: hidden; padding: 10px; }
.ranking-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 9px; }
.ranking-title { display: flex; min-width: 0; align-items: center; }.ranking-title strong { color: #20232a; font-size: 19px; font-weight: 800; }.ranking-title img { display: block; max-width: 100px; max-height: 28px; object-fit: contain; }
.ranking-more { font-size: 12px; }.ranking-boards { display: grid; gap: 8px; }.ranking-boards--1 { grid-template-columns: 1fr; }.ranking-boards--2 { grid-template-columns: repeat(2, minmax(0, 1fr)); }.ranking-boards--3 { grid-template-columns: repeat(3, minmax(0, 1fr)); }
.ranking-board { min-width: 0; padding: 9px 7px; }.ranking-board h4 { margin: 0 0 8px; overflow: hidden; font-size: 13px; font-weight: 700; text-overflow: ellipsis; white-space: nowrap; }
.ranking-products { display: grid; gap: 5px; }.ranking-product { position: relative; display: flex; min-width: 0; align-items: center; gap: 5px; padding: 5px; border-radius: 6px; background: #fff; }.ranking-product img, .ranking-cover { width: 42px; height: 42px; flex: none; border-radius: 5px; object-fit: cover; }.ranking-cover { display: flex; align-items: center; justify-content: center; color: #9eb8d8; background: #edf6ff; font-size: 9px; }
.ranking-number { position: absolute; top: 4px; left: 4px; z-index: 1; min-width: 12px; padding: 1px 3px; color: #fff; background: #aeb8cf; font-size: 10px; line-height: 12px; text-align: center; }.ranking-number--1 { background: #ff6741; }.ranking-number--2 { background: #f5bd36; }.ranking-product__info { display: flex; min-width: 0; flex: 1; flex-direction: column; gap: 3px; }.ranking-product__info b { overflow: hidden; color: #333; font-size: 11px; font-weight: 500; text-overflow: ellipsis; white-space: nowrap; }.ranking-product__info em { font-size: 12px; font-style: normal; font-weight: 700; }
.ranking-boards--3 .ranking-product { flex-direction: column; align-items: stretch; }.ranking-boards--3 .ranking-product img, .ranking-boards--3 .ranking-cover { width: 100%; height: 46px; }.ranking-boards--3 .ranking-product__info { gap: 1px; }
</style>
