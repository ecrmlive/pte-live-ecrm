<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { fetchNotice, type Notice } from "@/api/content";

const route = useRoute();
const row = ref<Notice | null>(null);

onMounted(async () => {
  const id = Number(route.params.id || 0);
  if (!id) return;
  row.value = await fetchNotice(id);
});
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <h1>{{ row?.title || "公告" }}</h1>
      <pre class="body">{{ row?.content || "" }}</pre>
    </section>
  </div>
</template>

<style scoped>
.panel {
  background: #fff;
  border-radius: 12px;
  padding: 28px;
}
.body {
  white-space: pre-wrap;
  font-family: inherit;
  line-height: 1.7;
  color: #333;
}
</style>
