<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { fetchAgreement, type Agreement } from "@/api/content";

const route = useRoute();
const row = ref<Agreement | null>(null);

async function load(key: string) {
  row.value = await fetchAgreement(key || "sys_user_agree");
}

onMounted(() => {
  void load(String(route.params.key || "sys_user_agree"));
});

watch(
  () => route.params.key,
  (k) => {
    void load(String(k || "sys_user_agree"));
  },
);
</script>

<template>
  <div class="pc-container">
    <section class="panel">
      <h1>{{ row?.label || "协议" }}</h1>
      <pre class="body">{{ row?.content || "暂无内容" }}</pre>
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
