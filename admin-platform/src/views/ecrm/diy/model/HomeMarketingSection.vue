<template>
  <section class="home-marketing drag optional" :class="{ selected: index === selectedIndex }" @click.stop="diyEditer(index)">
    <header :style="{ background: params.tone || '#1677ff' }"><strong>{{ params.title || '推荐商品' }}</strong><span>{{ params.subtitle || '精选好物' }}</span><em>{{ params.action || '查看全部 ›' }}</em></header>
    <div class="home-marketing__products">
      <article v-for="(entry, entryIndex) in entries" :key="entryIndex">
        <div class="product-icon" :style="{ background: entry.background || '#eef4ff' }">{{ entry.icon || '好物' }}</div>
        <b>{{ entry.title }}</b><small>{{ entry.subtitle }}</small><strong>{{ entry.price }}</strong>
      </article>
    </div>
    <div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
  </section>
</template>

<script>
export default {
  inject: ['diyModel'],
  props: ['item', 'index', 'selectedIndex'],
  computed: {
    params() { return this.item?.params || {}; },
    entries() { return this.item?.data || []; },
  },
  methods: {
    diyEditer(index) { this.diyModel?.onEditer(index); },
    diyDeleteItem(index) { this.diyModel?.onDeleleItem(index); },
  },
};
</script>

<style lang="scss" scoped>
.home-marketing { position: relative; margin: 0 10px 10px; overflow: hidden; border-radius: 12px; background: #fff; }.home-marketing header { display: flex; align-items: baseline; gap: 7px; padding: 9px 11px; color: #fff; }.home-marketing header strong { font-size: 15px; }.home-marketing header span { font-size: 11px; opacity: .85; }.home-marketing header em { margin-left: auto; font-size: 11px; font-style: normal; opacity: .9; }.home-marketing__products { display: grid; grid-template-columns: repeat(2, 1fr); gap: 8px; padding: 10px; }.home-marketing article { min-width: 0; overflow: hidden; border-radius: 8px; background: #f8f9fc; padding: 7px; }.product-icon { display: flex; height: 48px; align-items: center; justify-content: center; border-radius: 6px; color: #52627a; font-size: 12px; font-weight: 700; }.home-marketing b, .home-marketing small, .home-marketing article > strong { display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }.home-marketing b { margin-top: 6px; color: #293042; font-size: 12px; }.home-marketing small { margin-top: 3px; color: #8f98a8; font-size: 10px; }.home-marketing article > strong { margin-top: 4px; color: #f04d4d; font-size: 13px; }
</style>
