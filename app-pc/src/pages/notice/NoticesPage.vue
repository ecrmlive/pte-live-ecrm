<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { fetchNotices, type Notice } from "@/api/content";

const router = useRouter();
const list = ref<Notice[]>([]);
const loading = ref(false);

onMounted(async () => {
  loading.value = true;
  try {
    const data = await fetchNotices(1, 50);
    list.value = data.list || [];
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <header>
        <h1>平台公告</h1>
        <p>来自平台运营的通知与说明。</p>
      </header>
      <p v-if="loading" class="muted">加载中…</p>
      <p v-else-if="!list.length" class="muted">暂无公告</p>
      <ul v-else class="list">
        <li v-for="n in list" :key="n.notice_id" @click="router.push(`/notices/${n.notice_id}`)">
          <strong>{{ n.title }}</strong>
          <span class="muted">{{ n.create_time || "" }}</span>
        </li>
      </ul>
    </section>
  </div>
</template>

<style scoped>
.panel {
  background: #fff;
  border-radius: 12px;
  padding: 28px;
}
header h1 {
  margin: 0 0 8px;
}
header p,
.muted {
  color: #888;
}
.list {
  list-style: none;
  padding: 0;
  margin: 16px 0 0;
}
.list li {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 0;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
}
.list li:hover strong {
  color: #e23030;
}
</style>
