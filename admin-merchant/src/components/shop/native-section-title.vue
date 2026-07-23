<script setup lang="ts">
defineOptions({ name: 'NativeSectionTitle' });

defineProps<{
  hint?: string;
  title: string;
}>();
</script>

<template>
  <div class="native-section-title-block">
    <div class="native-section-title">{{ title }}</div>
    <div v-if="hint" class="native-section-title__hint">{{ hint }}</div>
  </div>
</template>
