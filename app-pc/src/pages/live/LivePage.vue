<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { fetchLiveRooms, type LiveRoom } from "@/api/live";

const router = useRouter();
const loading = ref(false);
const list = ref<LiveRoom[]>([]);

onMounted(async () => {
  loading.value = true;
  try {
    const data = await fetchLiveRooms();
    list.value = data.list || [];
  } finally {
    loading.value = false;
  }
});

function statusText(s: number) {
  if (s === 101) return "直播中";
  if (s === 103) return "已结束";
  return "未开始";
}

function open(id: number) {
  router.push(`/live/${id}`);
}
</script>

<template>
  <div class="page">
    <header class="head">
      <h1>直播专区</h1>
      <p>演示直播间 · 无真实推流 · 可查看挂货商品</p>
    </header>
    <p v-if="loading" class="hint">加载中…</p>
    <p v-else-if="!list.length" class="hint">暂无直播</p>
    <div v-else class="grid">
      <article v-for="r in list" :key="r.broadcast_room_id" class="card">
        <div class="row">
          <strong>{{ r.name }}</strong>
          <span class="badge" :class="{ on: r.live_status === 101 }">
            {{ statusText(r.live_status) }}
          </span>
        </div>
        <p class="mer">{{ r.anchor_name }} · {{ r.mer_name }}</p>
        <button type="button" class="btn" @click="open(r.broadcast_room_id)">进入直播间</button>
      </article>
    </div>
  </div>
</template>

<style scoped>
.page {
  max-width: 1080px;
  margin: 0 auto;
  padding: 32px 24px 64px;
}
.head h1 {
  margin: 0 0 8px;
  font-size: 28px;
}
.head p,
.hint {
  color: #888;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 16px;
  margin-top: 24px;
}
.card {
  background: #fff;
  border: 1px solid #eee;
  border-radius: 12px;
  padding: 20px;
}
.row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}
.badge {
  font-size: 12px;
  color: #999;
}
.badge.on {
  color: #c45c26;
}
.mer {
  margin: 8px 0 16px;
  font-size: 13px;
  color: #999;
}
.btn {
  width: 100%;
  height: 40px;
  border: 0;
  border-radius: 8px;
  background: #1a1a1a;
  color: #fff;
  cursor: pointer;
}
</style>
