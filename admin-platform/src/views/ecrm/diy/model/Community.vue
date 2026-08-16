<template>
  <div class="community-host" :style="hostStyle" @click.stop="diyEditer(index)">
    <section class="community-card" :style="cardStyle">
      <header class="community-head" :style="headStyle">
        <div class="community-title"><img v-if="params.titleType === 'image' && params.titleImage" v-img-url="params.titleImage" alt="种草社区标题" /><strong v-else>{{ params.title }}</strong></div>
        <span :style="{ color: style.buttonColor, fontSize: `${style.buttonSize}px` }">{{ params.more }} ›</span>
      </header>
      <div class="community-posts" :class="`community-posts--${params.layout}`" :style="{ gap: `${style.contentGap}px` }">
        <article v-for="(post, postIndex) in visiblePosts" :key="postIndex" class="community-post" :style="postStyle">
          <img v-if="post.image" v-img-url="post.image" alt="" />
          <div v-else class="community-cover">种草</div>
          <b v-if="params.showTitle">{{ post.title || '分享我的美好生活' }}</b>
          <footer v-if="params.showAvatar || params.showAuthor"><span v-if="params.showAvatar" class="avatar">{{ post.author?.slice(0, 1) || '好' }}</span><span v-if="params.showAuthor">{{ post.author || '好物分享官' }}</span><i v-if="post.images">{{ post.images }}图</i></footer>
        </article>
      </div>
    </section>
    <div class="btn-edit-del"><div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div></div>
  </div>
</template>

<script>
const fallback = () => [
  { title: '把春天装进生活里', author: '浅笑回眸', images: 8, image: '' },
  { title: '像我这种乐形身材又怎样', author: '国宝小熊猫', images: 6, image: '' },
  { title: '发现日常里的美好瞬间', author: '阿秋', images: 3, image: '' },
];
export default {
  inject: ['diyModel'], props: ['item', 'index', 'selectedIndex'],
  computed: {
    params() { return { title: '种草社区', more: '好物分享', titleType: 'text', titleImage: '', layout: 'scroll', showNum: 3, showTitle: true, showAvatar: true, showAuthor: true, ...(this.item?.params || {}) }; },
    style() { return { background: '#f5f5f5', cardBackground: '#ffffff', headStart: '#e93323', headEnd: '#ff7931', buttonColor: '#ffffff', buttonSize: 12, contentGap: 10, contentRadius: 8, contentShadow: 'off', paddingTop: 10, paddingBottom: 10, paddingLeft: 10, marginTop: 10, radius: 10, shadow: 'off', ...(this.item?.style || {}) }; },
    visiblePosts() { const list = Array.isArray(this.item?.data) && this.item.data.length ? this.item.data : fallback(); return list.slice(0, Number(this.params.showNum) || 3); },
    hostStyle() { return { background: this.style.background, padding: `${Number(this.style.paddingTop)}px ${Number(this.style.paddingLeft)}px ${Number(this.style.paddingBottom)}px`, marginTop: `${Number(this.style.marginTop)}px` }; },
    cardStyle() { return { background: this.style.cardBackground, borderRadius: `${Number(this.style.radius)}px`, boxShadow: this.style.shadow === 'on' ? '0 6px 18px rgba(15,23,42,.12)' : 'none' }; },
    headStyle() { return { background: `linear-gradient(100deg, ${this.style.headStart}, ${this.style.headEnd})` }; },
    postStyle() { return { borderRadius: `${Number(this.style.contentRadius)}px`, boxShadow: this.style.contentShadow === 'on' ? '0 4px 12px rgba(15,23,42,.12)' : 'none' }; },
  },
  methods: { diyEditer(i) { this.diyModel?.onEditer(i); }, diyDeleteItem(i) { this.diyModel?.onDeleleItem(i); } },
};
</script>

<style lang="scss" scoped>
.community-host { position: relative; box-sizing: border-box; width: 100%; }.community-card { overflow: hidden; padding-bottom: 10px; }.community-head { display: flex; align-items: center; justify-content: space-between; padding: 12px 13px; color: #fff; }.community-title strong { font-size: 19px; font-weight: 800; }.community-title img { display: block; max-width: 112px; max-height: 28px; object-fit: contain; }.community-head > span { white-space: nowrap; }.community-posts { display: grid; padding: 10px; gap: 10px; }.community-posts--scroll { grid-auto-columns: 145px; grid-auto-flow: column; overflow: hidden; }.community-posts--grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }.community-post { min-width: 0; overflow: hidden; background: #fff; }.community-post img, .community-cover { width: 100%; height: 144px; display: block; object-fit: cover; }.community-cover { display: flex; align-items: center; justify-content: center; color: #fff; background: linear-gradient(135deg, #93c5fd, #64748b); font-size: 18px; }.community-post b { display: block; overflow: hidden; padding: 7px 7px 2px; color: #303133; font-size: 12px; font-weight: 500; text-overflow: ellipsis; white-space: nowrap; }.community-post footer { display: flex; align-items: center; gap: 5px; padding: 4px 7px 8px; color: #7d8490; font-size: 10px; }.avatar { width: 17px; height: 17px; display: flex; align-items: center; justify-content: center; border-radius: 50%; color: #fff; background: #fb923c; font-size: 9px; }.community-post i { margin-left: auto; font-style: normal; }
</style>
