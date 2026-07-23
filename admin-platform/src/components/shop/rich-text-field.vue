<script setup lang="ts">
import { ElInput } from 'element-plus';
const model = defineModel<string>({ default: '' });
</script>
<template>
  <ElInput v-model="model" type="textarea" :rows="8" placeholder="富文本 HTML" />
</template>
