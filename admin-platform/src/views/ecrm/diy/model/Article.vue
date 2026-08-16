<template>
  <div class="article-host" :style="hostStyle" @click.stop="diyEditer(index)">
    <section class="article-card" :class="`article-card--${params.layout}`" :style="cardStyle">
      <article v-for="(article, articleIndex) in previewItems.slice(0, Number(params.showNum) || 3)" :key="articleIndex" class="article-item" :style="itemStyle">
        <img v-if="article.image" v-img-url="article.image" alt="" />
        <div v-else class="article-cover">文章</div>
        <div class="article-item__body"><b :style="{ color: style.titleColor, fontWeight: style.titleWeight === 'bold' ? 700 : 400 }">{{ article.article_title || '文章标题文章标题文章标题' }}</b><footer><span v-if="params.showDate" :style="{ color: style.metaColor }">{{ article.create_time || '2021-05-20 09:40' }}</span><span v-if="params.showViews" :style="{ color: style.viewColor }">◉ {{ article.views_num || 0 }}</span></footer></div>
      </article>
    </section>
    <div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
  </div>
</template>

<script>
import { resolveAutoSourcePreviewItems } from '../params/shared/marketing-helpers';

export default {
  inject: ['diyModel'],
	data() {
		return {};
	},
	props: ['item', 'index', 'selectedIndex'],
	computed: {
		params() { return { layout: 'list', showNum: 3, showDate: true, showViews: true, ...(this.item?.params || {}) }; },
		style() { return { bgcolor: '#f5f5f5', background: '#ffffff', paddingTop: 0, paddingBottom: 0, paddingLeft: 10, marginTop: 10, radius: 0, imageRadius: 0, titleColor: '#333333', metaColor: '#999999', viewColor: '#999999', titleWeight: 'normal', shadow: 'off', ...(this.item?.style || {}) }; },
		hostStyle() { return { background: this.style.bgcolor, padding: `${Number(this.style.paddingTop)}px ${Number(this.style.paddingLeft)}px ${Number(this.style.paddingBottom)}px`, marginTop: `${Number(this.style.marginTop)}px` }; },
		cardStyle() { return { background: this.style.background, borderRadius: `${Number(this.style.radius)}px`, boxShadow: this.style.shadow === 'on' ? '0 5px 16px rgba(15,23,42,.12)' : 'none' }; },
		itemStyle() { return { borderRadius: `${Number(this.style.imageRadius)}px` }; },
		previewItems() {
			return resolveAutoSourcePreviewItems(this.item, {
				article_title: '此处是文章标题',
				image: '',
				show_type: 10,
				views_num: 0,
			});
		},
	},
	methods: {
  diyEditer(index) {
    this.diyModel?.onEditer(index);
  },
  diyDeleteItem(index) {
    this.diyModel?.onDeleleItem(index);
  },
}
};
</script>

<style lang="scss" scoped>
.article-host { position: relative; box-sizing: border-box; width: 100%; }.article-card { overflow: hidden; }.article-card--list { display: grid; gap: 10px; padding: 10px; }.article-card--grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 10px; padding: 10px; }.article-card--scroll { display: grid; grid-auto-columns: 48%; grid-auto-flow: column; gap: 10px; overflow: hidden; padding: 10px; }.article-item { overflow: hidden; background: #fff; }.article-card--list .article-item { display: flex; gap: 10px; padding-bottom: 10px; border-bottom: 1px solid #f2f2f2; }.article-card--list .article-item:last-child { padding-bottom: 0; border-bottom: 0; }.article-item img, .article-cover { width: 100%; height: 108px; display: block; object-fit: cover; }.article-card--list .article-item img, .article-card--list .article-cover { width: 124px; height: 76px; flex: none; }.article-cover { display: flex; align-items: center; justify-content: center; color: #9eb8d8; background: #edf6ff; font-size: 12px; }.article-item__body { display: flex; min-width: 0; flex: 1; flex-direction: column; }.article-item__body b { display: -webkit-box; overflow: hidden; font-size: 14px; line-height: 20px; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }.article-card--grid .article-item__body, .article-card--scroll .article-item__body { padding: 7px; }.article-item footer { display: flex; justify-content: space-between; gap: 6px; margin-top: auto; padding-top: 7px; font-size: 10px; }.article-card--grid .article-item footer, .article-card--scroll .article-item footer { font-size: 9px; }
</style>
