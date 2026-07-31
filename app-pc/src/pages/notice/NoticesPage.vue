<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { fetchNotices, type Notice } from "@/api/content";

const router = useRouter();
const list = ref<Notice[]>([]);
const loading = ref(false);
const latest = computed(() => list.value.slice(0, 6));

function preview(content: string) {
  return content.replace(/\s+/g, " ").slice(0, 66) || "七禧商城官方资讯";
}

function go(row: Notice) { void router.push({ name: "notice-detail", params: { id: row.notice_id } }); }

onMounted(async () => {
  loading.value = true;
  try { const data = await fetchNotices(1, 20); list.value = data.list || []; }
  finally { loading.value = false; }
});
</script>

<template>
  <div class="news-page">
    <div class="pc-container news-layout">
      <section class="news-panel">
        <header class="news-title"><strong>生活</strong></header>
        <p v-if="loading" class="muted">资讯加载中…</p>
        <p v-else-if="!list.length" class="muted empty">暂无已发布资讯</p>
        <article v-for="row in list" :key="row.notice_id" class="news-row" role="link" tabindex="0" @click="go(row)" @keydown.enter="go(row)">
          <img v-if="row.cover_url" :src="row.cover_url" :alt="row.title" />
          <div v-else class="news-cover-fallback" aria-hidden="true">七禧资讯</div>
          <div class="news-body"><h1>{{ row.title }}</h1><p>{{ preview(row.content) }}</p><time>◷ {{ row.create_time }}</time></div>
        </article>
        <nav v-if="list.length" class="pagination"><button type="button" disabled>‹</button><button class="active" type="button">1</button><button type="button" disabled>›</button></nav>
      </section>
      <aside class="latest-panel"><h2>最新资讯</h2><ol><li v-for="(row, index) in latest" :key="row.notice_id" @click="go(row)"><i>{{ index + 1 }}</i><span>{{ row.title }}</span></li></ol></aside>
    </div>
  </div>
</template>

<style scoped>
.news-page { min-height: 650px; padding: 22px 0 52px; background: #f7f7f7; }.news-layout { display: grid; grid-template-columns: minmax(0, 1fr) 310px; gap: 10px; }.news-panel, .latest-panel { min-height: 600px; background: #fff; border-radius: 8px; }.news-title { padding: 26px 40px 20px; border-bottom: 1px solid #eee; }.news-title strong { position: relative; color: #f13728; font-size: 17px; font-weight: 500; }.news-title strong::after { position: absolute; bottom: -21px; left: 0; width: 68px; height: 2px; background: #f13728; content: ""; }.news-row { display: grid; grid-template-columns: 282px minmax(0, 1fr); gap: 30px; margin: 0 32px; padding: 42px 0; border-bottom: 1px solid #eee; cursor: pointer; }.news-row:hover h1 { color: #f13728; }.news-row img, .news-cover-fallback { width: 282px; height: 174px; object-fit: cover; background: #f5f5f5; }.news-cover-fallback { display: grid; place-items: center; color: #bbb; font-size: 14px; }.news-body { display: flex; min-width: 0; flex-direction: column; }.news-body h1 { margin: 5px 0 18px; color: #333; font-size: 20px; transition: color .15s ease; }.news-body p { margin: 0; color: #555; line-height: 1.8; }.news-body time { margin-top: auto; color: #aaa; font-size: 14px; }.latest-panel { padding: 22px 28px; }.latest-panel h2 { margin: 0; padding-bottom: 20px; border-bottom: 1px dashed #eee; color: #333; font-size: 18px; text-align: center; }.latest-panel ol { margin: 18px 0; padding: 0; list-style: none; }.latest-panel li { display: flex; gap: 12px; align-items: center; padding: 11px 0; cursor: pointer; }.latest-panel li:hover span { color: #f13728; }.latest-panel i { display: grid; width: 20px; height: 20px; flex: 0 0 auto; place-items: center; background: #f5f5f5; color: #999; font-size: 12px; font-style: normal; }.latest-panel li:first-child i { background: #f13728; color: #fff; }.latest-panel span { overflow: hidden; color: #666; font-size: 14px; text-overflow: ellipsis; white-space: nowrap; }.pagination { display: flex; justify-content: center; gap: 10px; padding: 30px 0; }.pagination button { width: 40px; height: 40px; border: 0; background: #f4f4f4; color: #bbb; }.pagination button.active { background: #4399f2; color: #fff; }.muted { padding: 42px; color: #999; text-align: center; }.empty { min-height: 380px; }
@media (max-width: 860px) { .news-layout { grid-template-columns: 1fr; }.latest-panel { min-height: 0; }.news-row { grid-template-columns: 180px 1fr; gap: 18px; margin: 0 18px; }.news-row img,.news-cover-fallback { width: 180px; height: 124px; } }
</style>
