<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { fetchLiveRoom, type LiveRoom } from "@/api/live";

const route = useRoute();
const router = useRouter();
const room = ref<LiveRoom | null>(null);

onMounted(async () => {
  const id = Number(route.params.id || 0);
  if (!id) return;
  try {
    room.value = await fetchLiveRoom(id);
  } catch (e) {
    alert((e as Error).message || "加载失败");
  }
});

function goGoods(id: number) {
  router.push(`/goods/${id}`);
}
</script>

<template>
  <div class="page" v-if="room">
    <header class="head">
      <h1>{{ room.name }}</h1>
      <p>{{ room.anchor_name }} · {{ room.mer_name }}</p>
    </header>
    <div class="stage">
      <video
        v-if="room.live_status === 101 && room.play_url"
        :poster="room.cover_img"
        :src="room.play_url"
        autoplay
        controls
        playsinline
      />
      <span v-else>{{ room.live_status === 101 ? '直播播放地址暂未配置' : '直播尚未开始或已结束' }}</span>
    </div>
    <h2>讲解商品</h2>
    <p v-if="!(room.goods || []).length" class="hint">暂无挂货</p>
    <ul v-else class="goods">
      <li v-for="g in room.goods" :key="g.product_id">
        <button type="button" @click="goGoods(g.product_id)">
          {{ g.store_name || `商品 #${g.product_id}` }} · ¥{{ Number(g.price || 0).toFixed(2) }}
        </button>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.page {
  max-width: 800px;
  margin: 0 auto;
  padding: 32px 24px 64px;
}
.head h1 {
  margin: 0 0 8px;
}
.head p,
.hint {
  color: #888;
}
.stage {
  margin: 24px 0;
  height: 280px;
  border-radius: 12px;
  background: linear-gradient(160deg, #2a2a2a, #111);
  color: #bbb;
  display: flex;
  align-items: center;
  justify-content: center;
}
.stage video {
  width: 100%;
  height: 100%;
  border-radius: inherit;
  background: #000;
  object-fit: contain;
}
.goods {
  list-style: none;
  padding: 0;
  margin: 0;
}
.goods li {
  margin-bottom: 10px;
}
.goods button {
  width: 100%;
  text-align: left;
  padding: 14px 16px;
  border: 1px solid #eee;
  border-radius: 8px;
  background: #fff;
  cursor: pointer;
}
</style>
