<template>
  <section class="home-feature drag optional" :class="{ selected: index === selectedIndex }" @click.stop="diyEditer(index)">
    <header>
      <div><strong>{{ params.title || '热门专区' }}</strong><span>{{ params.subtitle }}</span></div>
      <em>{{ params.action || '更多 ›' }}</em>
    </header>
    <div class="home-feature__grid" :class="`columns-${columns}`">
      <article v-for="(entry, entryIndex) in entries" :key="entryIndex" :style="{ background: entry.background || '#fff4ee' }">
        <i :style="{ background: entry.iconBg || '#ff7a59' }">{{ entry.icon || '荐' }}</i>
        <div><b>{{ entry.title }}</b><small>{{ entry.subtitle }}</small></div>
        <mark v-if="entry.badge">{{ entry.badge }}</mark>
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
    columns() { return Math.min(4, Math.max(2, Number(this.params.columns) || 2)); },
  },
  methods: {
    diyEditer(index) { this.diyModel?.onEditer(index); },
    diyDeleteItem(index) { this.diyModel?.onDeleleItem(index); },
  },
};
</script>

<style lang="scss" scoped>
.home-feature { position: relative; margin: 0 10px 10px; padding: 12px; border-radius: 12px; background: #fff; }
.home-feature header { display: flex; align-items: baseline; justify-content: space-between; margin-bottom: 9px; }
.home-feature header div { display: flex; align-items: baseline; gap: 6px; }
.home-feature strong { color: #202636; font-size: 15px; }.home-feature span, .home-feature em { color: #8b95a7; font-size: 11px; font-style: normal; }
.home-feature__grid { display: grid; gap: 7px; }.home-feature__grid.columns-2 { grid-template-columns: repeat(2, 1fr); }.home-feature__grid.columns-3 { grid-template-columns: repeat(3, 1fr); }.home-feature__grid.columns-4 { grid-template-columns: repeat(4, 1fr); }
.home-feature article { position: relative; display: flex; min-height: 58px; align-items: center; gap: 7px; overflow: hidden; border-radius: 9px; padding: 7px; }.home-feature i { display: flex; width: 27px; height: 27px; flex: none; align-items: center; justify-content: center; border-radius: 50%; color: #fff; font-size: 12px; font-style: normal; font-weight: 700; }.home-feature b, .home-feature small { display: block; }.home-feature b { color: #2b3140; font-size: 12px; }.home-feature small { margin-top: 3px; color: #9199a8; font-size: 10px; }.home-feature mark { position: absolute; right: 3px; top: 3px; border-radius: 4px; background: #fff; color: #fa5e5e; font-size: 9px; padding: 1px 3px; }
</style>
